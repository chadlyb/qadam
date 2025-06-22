package main

import (
	"fmt"
	"os"
	"strings"
)

// ISO CP 852
// We use unicode 2400 "symbol for NUL" for NUL (0), so it is printable
// We use unicode 2423 "Open Box" for NBSP (255), so it is printable
// We use unicode 00A5 "Yen" aka the paragraph sign for the section symbol (u00A7), so there aren't two.
// We use unicode 00AF "Macron" for the soft hyphen, so there aren't two.
const ctrlCharacters = "\u2400\u263a\u263b\u2665\u2666\u2663\u2660\u2022\u25D8\u25CB\u25D9\u2642\u2640\u266A\u266B\u263C\u25BA\u25C4\u2195\u203C\u00B6\u00A5\u25AC\u21A8\u2191\u2193\u2192\u2190\u221F\u2194\u25B2\u25BC"
const charsetString = ctrlCharacters + " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~⌂ÇüéâäůćçłëŐőîŹÄĆÉĹĺôöĽľŚśÖÜŤťŁ×čáíóúĄąŽžĘę¬źČş«»░▒▓│┤ÁÂĚŞ╣║╗╝Żż┐└┴┬├─┼Ăă╚╔╩╦╠═╬¤đĐĎËďŇÍÎě┘┌█▄ŢŮ▀ÓßÔŃńňŠšŔÚŕŰýÝţ´\u00AF˝˛ˇ˘§÷¸°¨˙űŘř■\u2423"

var tableRunes = []rune(charsetString)

func toString(bytes []byte) string {
	var b strings.Builder
	for _, v := range bytes {
		if v == '\n' {
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

func realMain() error {
	if len(os.Args) != 2 {
		return fmt.Errorf("Usage: %s <source file>\n", os.Args[0])
	}

	srcPath := os.Args[1]

	data, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("Error reading %s: %w", srcPath, err)
	}

	// Crawl for strings, output OFFSET LENGTH "<string>"

	stringBegin := 0
	at := 0
	end := len(data)
	for at != end {
		if data[at] == 0 {
			if heuristicIsHumanString(data[stringBegin:at]) {
				fmt.Printf("%08x-%08x: \"%v\"\n", stringBegin, at+1, toString(data[stringBegin:at]))
			}
			stringBegin = at + 1
		}

		at++
	}
	if stringBegin != at {
		if heuristicIsHumanString(data[stringBegin:at]) {
			fmt.Printf("%08x-%08x: \"%v\" NO_NUL\n", stringBegin, at, toString(data[stringBegin:at]))
		}
	}
	return nil
}

func main() {
	err := realMain()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
}
