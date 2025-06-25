package main

import (
	"errors"
	"fmt"
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

//const charsetStartingAt0x20 = " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~ ÇüéâäůćçłëŐőîŹÄĆÉĹĺôöĽľŚśÖÜŤťŁ×čáíóúĄąŽžĘę¬źČş«»░▒▓│┤ÁÂĚŞ╣║╗╝Żż┐└┴┬├─┼Ăă╚╔╩╦╠═╬¤đĐĎËďŇÍÎě┘┌█▄ŢŮ▀ÓßÔŃńňŠšŔÚŕŰýÝţ´-˝˛ˇ˘§÷¸°¨˙űŘř■ "

var charsetRunes = []rune(charsetString)

var charset = map[rune]byte{}

func init() {
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

func printByte(b byte) {
	switch {
	case b == '\n':
		fmt.Printf("\\n")
	case b == '\t':
		fmt.Printf("\\t")
	case b == '\\':
		fmt.Printf("\\\\")
	case b == '"':
		fmt.Printf("\\\"")
	default:
		fmt.Printf("%c", tableRunes[b])
	}
}

func toString(bytes []byte) string {
	var b strings.Builder
	for _, v := range bytes {
		if v == '"' {
			fmt.Fprintf(&b, "\\\"")
		} else if v == '\n' {
			fmt.Fprintf(&b, "\\n")
		} else if v == '\t' {
			fmt.Fprintf(&b, "\\t")
		} else {
			fmt.Fprintf(&b, "%c", tableRunes[v])
		}
	}
	return b.String()
}

func heuristicIsHumanString(bytes []byte) bool {
	if len(bytes) < 2 {
		return false
	}
	spaceCount := 0
	ctrlCount := 0
	for _, v := range bytes {
		if v == ' ' {
			spaceCount++
		}
		if v < 32 {
			ctrlCount++
		}
	}
	if len(bytes) < 4 && ctrlCount > 0 {
		return false
	}

	if len(bytes) > 16 && spaceCount == 0 {
		return false
	}

	return true
}
