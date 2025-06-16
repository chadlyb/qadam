package main

/*
 Generate TEXTS.FIL
 - Read the file specified on the command-line. Use tabs, spaces, \r, and \n as token delimiters (except inside quotes.)
 - Ignore anything on a line past ';', and also spaces/blanks/empty lines
 - When we see [, there will be some number of hex bytes followed by ] that go verbatim into the file
 - When we see '"', we parse a string, read a NUL-terminated string goes straight into the file (look up in the charset table, and add 0x31 to obfuscate)
   - the string ends with "
   - \n \t \" \\ and \x## are supported escape sequences.
   - if a character is missing from the charset, this is a fatal error.
 - When we see SECTION N, take note of this spot, then this will be pointed to by the directory
   - Sections are numbered 0..N and sequential. Anything else is a fatal error.

 The file format is:
   - a byte indicating how many sections there are
   - three bytes per section indicating their index (from the beginning INCLUDING this directory) into the file
   - then three bytes containing the filesize.
   - the data from above, not otherwise modified.

 Write it out to TEXTS.FIL
*/

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

const charsetStartingAt0x20 = " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~ ÇüéâäůćçłëŐőîŹÄĆÉĹĺôöĽľŚśÖÜŤťŁ×čáíóúĄąŽžĘę¬źČş«»░▒▓│┤ÁÂĚŞ╣║╗╝Żż┐└┴┬├─┼Ăă╚╔╩╦╠═╬¤đĐĎËďŇÍÎě┘┌█▄ŢŮ▀ÓßÔŃńňŠšŔÚŕŰýÝţ´-˝˛ˇ˘§÷¸°¨˙űŘř■ "

var charsetRunes = []rune(charsetStartingAt0x20)

var charset = map[rune]byte{}

type Section struct {
	index int // SECTION N
	pos   int // offset in output (after directory)
}

func main() {
	for k, v := range charsetRunes {
		_, has := charset[v]
		if !has {
			charset[v] = byte(k) + 0x20
		}
	}

	charset[' '] = ' '
	charset['\n'] = '\n'
	charset['\t'] = '\t'
	// Unicode madness
	charset['–'] = charset['-']
	charset['‘'] = charset['\'']
	charset['’'] = charset['\'']
	charset['—'] = charset['—']

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <inputfile>\n", os.Args[0])
		os.Exit(1)
	}
	infile := os.Args[1]
	f, err := os.Open(infile)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	output, err := processFile(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile("TEXTS.FIL", output, 0644); err != nil {
		panic(err)
	}
}

// The main parsing and file format logic
func processFile(r io.Reader) ([]byte, error) {
	scanner := bufio.NewScanner(r)
	var sections []Section
	var outData []byte
	expectedSection := 0

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Trim spaces
		line = strings.TrimSpace(line)
		if line == "" {
			continue // skip blank lines
		}

		// Tokenize line (preserving quoted strings as single tokens)
		tokens, err := tokenize(line)
		if err != nil {
			return nil, fmt.Errorf("Line %d: %v", lineNum, err)
		}

		i := 0
		for i < len(tokens) {
			token := tokens[i]
			switch {
			case strings.ToUpper(token) == "SECTION":
				if i+1 >= len(tokens) {
					return nil, fmt.Errorf("Line %d: SECTION missing argument", lineNum)
				}
				n, err := strconv.Atoi(tokens[i+1])
				if err != nil {
					return nil, fmt.Errorf("Line %d: Invalid SECTION number: %v", lineNum, err)
				}
				if n != expectedSection {
					return nil, fmt.Errorf("Line %d: Out-of-order SECTION, expected %d got %d", lineNum, expectedSection, n)
				}
				sections = append(sections, Section{index: n, pos: len(outData)})
				expectedSection++
				i += 2
			case strings.HasPrefix(token, "["):
				// Hex block
				if !strings.HasSuffix(token, "]") {
					// Possibly split across tokens
					hexstr := token[1:]
					for {
						i++
						if i >= len(tokens) {
							return nil, fmt.Errorf("Line %d: Missing closing ] for hex block", lineNum)
						}
						hexpart := tokens[i]
						if strings.HasSuffix(hexpart, "]") {
							hexstr += strings.TrimSuffix(hexpart, "]")
							break
						} else {
							hexstr += hexpart
						}
					}
					// Now hexstr is all hex digits
					bytes, err := hex.DecodeString(hexstr)
					if err != nil {
						return nil, fmt.Errorf("Line %d: Invalid hex: %v", lineNum, err)
					}
					outData = append(outData, bytes...)
					i++
				} else {
					// [010203]
					hexstr := token[1 : len(token)-1]
					bytes, err := hex.DecodeString(hexstr)
					if err != nil {
						return nil, fmt.Errorf("Line %d: Invalid hex: %v", lineNum, err)
					}
					outData = append(outData, bytes...)
					i++
				}
			case strings.HasPrefix(token, "\""):
				// Quoted string
				s, err := parseStringToken(token, tokens, &i)
				if err != nil {
					return nil, fmt.Errorf("Line %d: %v", lineNum, err)
				}
				const key = 0x31
				// Charset lookup and obfuscation
				for _, ch := range s {
					enc, ok := charset[ch]
					if !ok {
						return nil, fmt.Errorf("Line %d: Character %q missing from charset", lineNum, ch)
					}
					outData = append(outData, enc+key)
				}
				outData = append(outData, 0) // NUL-terminate
				i++
			default:
				// Ignore
				return nil, fmt.Errorf("Line %d: Unrecognized token '%v'", lineNum, token)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	// Write out file format
	numSections := len(sections)
	if numSections == 0 {
		return nil, errors.New("No sections found")
	}
	var dir []byte
	dirSize := 1 + 3*numSections + 3 // num, offsets, filesize
	dir = append(dir, byte(numSections))
	for _, s := range sections {
		// 3-byte offset, including directory size
		offs := s.pos + dirSize
		dir = append(dir, byte((offs)&0xFF), byte((offs>>8)&0xFF), byte(offs>>16&0xFF))
	}
	totalSize := len(outData) + dirSize
	dir = append(dir, byte((totalSize)&0xFF), byte((totalSize>>8)&0xFF), byte(totalSize>>16&0xFF))
	dir = append(dir, outData...)
	return dir, nil
}

// Tokenize a line using tabs, spaces, \r, \n as delimiters, except inside quotes
func tokenize(s string) ([]string, error) {
	var tokens []string
	var sb strings.Builder
	inQuote := false
	escapedQuote := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if inQuote {
			sb.WriteByte(c)

			if c == '"' && !escapedQuote {
				tokens = append(tokens, sb.String())
				sb.Reset()
				inQuote = false
			}

			escapedQuote = !escapedQuote && c == '\\'
		} else {
			if c == '"' {
				if sb.Len() > 0 {
					tokens = append(tokens, sb.String())
					sb.Reset()
				}
				inQuote = true
				sb.WriteByte(c)
			} else if unicode.IsSpace(rune(c)) || c == '\r' || c == '\n' || c == '\t' {
				if sb.Len() > 0 {
					tokens = append(tokens, sb.String())
					sb.Reset()
				}
			} else if c == ';' {
				break
			} else {
				sb.WriteByte(c)
			}
		}
	}
	if sb.Len() > 0 {
		if inQuote {
			return nil, errors.New("Unterminated quoted string")
		}
		tokens = append(tokens, sb.String())
	}
	return tokens, nil
}

// Parses quoted string token and returns decoded string, advances i as necessary
func parseStringToken(token string, tokens []string, i *int) (string, error) {
	// token starts with "
	s := token[1:]
	for !strings.HasSuffix(s, "\"") {
		*i++
		if *i >= len(tokens) {
			return "", errors.New("Unterminated quoted string")
		}
		s += " " + tokens[*i]
	}
	s = s[:len(s)-1] // remove ending "
	// Now unescape
	return unescapeString(s)
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
