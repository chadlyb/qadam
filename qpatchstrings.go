package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// ISO CP 852
// We use unicode 2400 "symbol for NUL" for NUL (0), so it is printable
// We use unicode 2423 "Open Box" for NBSP (255), so it is printable
// We use unicode 00A5 "Yen" aka the paragraph sign for the section symbol (u00A7), so there aren't two.
// We use unicode 00AF "Macron" for the soft hyphen, so there aren't two.
const ctrlCharacters = "\u2400\u263a\u263b\u2665\u2666\u2663\u2660\u2022\u25D8\u25CB\u25D9\u2642\u2640\u266A\u266B\u263C\u25BA\u25C4\u2195\u203C\u00B6\u00A5\u25AC\u21A8\u2191\u2193\u2192\u2190\u221F\u2194\u25B2\u25BC"
const charsetString = ctrlCharacters + " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~⌂ÇüéâäůćçłëŐőîŹÄĆÉĹĺôöĽľŚśÖÜŤťŁ×čáíóúĄąŽžĘę¬źČş«»░▒▓│┤ÁÂĚŞ╣║╗╝Żż┐└┴┬├─┼Ăă╚╔╩╦╠═╬¤đĐĎËďŇÍÎě┘┌█▄ŢŮ▀ÓßÔŃńňŠšŔÚŕŰýÝţ´\u00AF˝˛ˇ˘§÷¸°¨˙űŘř■\u2423"

var charsetRunes = []rune(charsetString)

var charset = map[rune]byte{}

func initCharset() {
	for k, v := range charsetRunes {
		_, has := charset[v]
		if has {
			fmt.Printf("%c redundant (%x vs %x)\n", v, k, charset[v])

		}
		if !has {
			charset[v] = byte(k)
		}
	}

	//charset[' '] = ' '
	charset['\r'] = '\r'
	charset['\n'] = '\n'
	charset['\t'] = '\t'
	// Unicode madness
	charset['–'] = charset['-']
	charset['—'] = charset['-']
	charset['‘'] = charset['\'']
	charset['’'] = charset['\'']
}

// Unescape string supporting \n, \t, \\, \x##
func unescapeString(s string) (string, error) {
	var out strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' {
			i++
			if i >= len(s) {
				return "", errors.New("Trailing backslash in string")
			}
			switch s[i] {
			case 'n':
				out.WriteByte('\n')
			case 't':
				out.WriteByte('\t')
			case 'r':
				out.WriteByte('\r')
			case '\\':
				out.WriteByte('\\')
			case '"':
				out.WriteByte('"')
			case 'x':
				if i+2 >= len(s) {
					return "", errors.New("Invalid \\x escape in string")
				}
				b, err := strconv.ParseUint(s[i+1:i+3], 16, 8)
				if err != nil {
					return "", fmt.Errorf("Invalid hex in \\x escape: %v", err)
				}
				out.WriteByte(byte(b))
				i += 2
			default:
				return "", fmt.Errorf("Unknown escape sequence: \\%c", s[i])
			}
		} else {
			out.WriteByte(s[i])
		}
	}
	return out.String(), nil
}

func fromString(str string) ([]byte, error) {
	unescaped, err := unescapeString(str)
	if err != nil {
		return nil, err
	}

	result := make([]byte, 0, len(unescaped))
	for _, v := range unescaped {
		n, ok := charset[v]
		if !ok {
			return nil, fmt.Errorf("Unrecognized character '%c' (%v)", v, int(v))
		}
		result = append(result, n)
	}

	return result, nil
}

// BEGIN-END:"STRING" ; COMMENT
const lineRegexSrc = `^\s*(?:0x)?(?P<begin>[0-9a-fA-F]+)\s*-\s*(?:0x)?(?P<end>[0-9a-fA-F]+)\s*:\s*"(?P<string>(?:[^"\\]|\\"|\\n|\\\\|\\t|\\r)*)"\s*(?:;.*)?$`

var lineRegex = regexp.MustCompile(lineRegexSrc)

func handleLine(data []byte, line string) error {
	matches := lineRegex.FindStringSubmatch(line)
	//fmt.Println(matches)
	if len(matches) != 4 {
		return fmt.Errorf("Line didn't match expected format")
	}

	patchBegin, err := strconv.ParseUint(matches[1], 16, 64)
	if err != nil {
		return fmt.Errorf("Couldn't parse begin offset (%w)", err)
	}
	patchEnd, err := strconv.ParseUint(matches[2], 16, 64)
	if err != nil {
		return fmt.Errorf("Couldn't parse end offset (%w)", err)
	}
	patchBytes, err := fromString(matches[3])
	if err != nil {
		return fmt.Errorf("Couldn't translate string (%w)", err)
	}

	len := uint64(len(patchBytes))
	if len+1 > patchEnd-patchBegin {
		return fmt.Errorf("Ignoring too-long string (%v > %v)", len+1, patchEnd-patchBegin)
	}

	for i := uint64(0); i != len; i++ {
		data[patchBegin+i] = patchBytes[i]
	}
	data[patchBegin+len] = 0
	return nil
}

func realMain() error {
	initCharset()

	if len(os.Args) != 4 {
		return fmt.Errorf("Usage: %s <source file> <dest file> <patch strings file>\n", os.Args[0])
	}

	srcPath := os.Args[1]
	destPath := os.Args[2]
	patchPath := os.Args[3]

	data, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("Error reading source file %s: %w", srcPath, err)
	}

	// Parse patch path, line by line.
	patchFile, err := os.Open(patchPath)
	if err != nil {
		return fmt.Errorf("Error opening patch file %s: %w", patchPath, err)
	}
	defer patchFile.Close()

	scanner := bufio.NewScanner(patchFile)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		err := handleLine(data, line)
		if err != nil {
			fmt.Printf("Ignored line %v due to error: %v\n", lineNum, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error while scanning patch file %s: %w", patchPath, err)
	}

	// Write to output file
	err = os.WriteFile(destPath, data, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", destPath, err)
		os.Exit(1)
	}

	fmt.Printf("Patched %s -> %s", srcPath, destPath)
	return nil
}

func main() {
	err := realMain()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
}
