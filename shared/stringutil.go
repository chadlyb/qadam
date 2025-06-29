package shared

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
const CharsetString = ctrlCharacters + " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~⌂ÇüéâäůćçłëŐőîŹÄĆÉĹĺôöĽľŚśÖÜŤťŁ×čáíóúĄąŽžĘę¬źČş«»░▒▓│┤ÁÂĚŞ╣║╗╝Żż┐└┴┬├─┼Ăă╚╔╩╦╠═╬¤đĐĎËďŇÍÎě┘┌█▄ŢŮ▀ÓßÔŃńňŠšŔÚŕŰýÝţ´\u00AF˝˛ˇ˘§÷¸°¨˙űŘř■\u2423"

var CharsetRunes = []rune(CharsetString)

var CharsetMapToByte = map[rune]byte{}

func init() {
	for k, v := range CharsetRunes {
		_, has := CharsetMapToByte[v]
		if has {
			fmt.Printf("%c redundant (%x vs %x)\n", v, k, CharsetMapToByte[v])

		}
		if !has {
			CharsetMapToByte[v] = byte(k)
		}
	}

	CharsetMapToByte['\r'] = '\r'
	CharsetMapToByte['\n'] = '\n'
	CharsetMapToByte['\t'] = '\t'

	// Unicode madness - convert from commonly used unicode characters to ASCII equivalents
	CharsetMapToByte['\u2013'] = CharsetMapToByte['-']  // en dash (–)
	CharsetMapToByte['\u2014'] = CharsetMapToByte['-']  // em dash (—)
	CharsetMapToByte['\u2018'] = CharsetMapToByte['\''] // left single quote (')
	CharsetMapToByte['\u2019'] = CharsetMapToByte['\''] // right single quote (')
	CharsetMapToByte['\u201A'] = CharsetMapToByte['"']  // left double quote (")
	CharsetMapToByte['\u201B'] = CharsetMapToByte['"']  // right double quote (")
}

func UnescapeString(s string) (string, error) {
	var out strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' {
			i++
			if i >= len(s) {
				return "", errors.New("trailing backslash in string")
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
					return "", errors.New("invalid \\x escape in string")
				}
				b, err := strconv.ParseUint(s[i+1:i+3], 16, 8)
				if err != nil {
					return "", fmt.Errorf("invalid hex in \\x escape: %v", err)
				}
				out.WriteByte(byte(b))
				i += 2
			default:
				return "", fmt.Errorf("unknown escape sequence: \\%c", s[i])
			}
		} else {
			out.WriteByte(s[i])
		}
	}
	return out.String(), nil
}

func ToString(bytes []byte) string {
	var b strings.Builder
	for _, v := range bytes {
		switch v {
		case '"':
			b.WriteString("\\\"")
		case '\n':
			b.WriteString("\\n")
		case '\t':
			b.WriteString("\\t")
		case '\\':
			b.WriteString("\\\\")
		default:
			b.WriteRune(CharsetRunes[v])
		}
	}
	return b.String()
}

func FromString(str string) ([]byte, error) {
	unescaped, err := UnescapeString(str)
	if err != nil {
		return nil, err
	}

	result := make([]byte, 0, len(unescaped))
	for _, v := range unescaped {
		n, ok := CharsetMapToByte[v]
		if !ok {
			return nil, fmt.Errorf("unrecognized character '%c' (%d)", v, int(v))
		}
		result = append(result, n)
	}

	return result, nil
}
