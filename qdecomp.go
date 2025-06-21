package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func isPrintable(b byte) bool {
	return b >= 0x20 && b <= 0x7E
}

// ISO CP 852
// We use unicode 2400 "symbol for NUL" for NUL (0), so it is printable
// We use unicode 2423 "Open Box" for NBSP (255), so it is printable
// We use unicode 00A5 "Yen" aka the paragraph sign for the section symbol (u00A7), so there aren't two.
// We use unicode 00AF "Macron" for the soft hyphen, so there aren't two.
const ctrlCharacters = "\u2400\u263a\u263b\u2665\u2666\u2663\u2660\u2022\u25D8\u25CB\u25D9\u2642\u2640\u266A\u266B\u263C\u25BA\u25C4\u2195\u203C\u00B6\u00A5\u25AC\u21A8\u2191\u2193\u2192\u2190\u221F\u2194\u25B2\u25BC"
const charsetString = ctrlCharacters + " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~⌂ÇüéâäůćçłëŐőîŹÄĆÉĹĺôöĽľŚśÖÜŤťŁ×čáíóúĄąŽžĘę¬źČş«»░▒▓│┤ÁÂĚŞ╣║╗╝Żż┐└┴┬├─┼Ăă╚╔╩╦╠═╬¤đĐĎËďŇÍÎě┘┌█▄ŢŮ▀ÓßÔŃńňŠšŔÚŕŰýÝţ´\u00AF˝˛ˇ˘§÷¸°¨˙űŘř■\u2423"

var tableRunes = []rune(charsetString)

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

func usage() {
	fmt.Println("Usage: qdecomp [--hex] <input file>")
	os.Exit(1)
}

func main() {
	hexMode := false
	files := []string{}

	for k, v := range os.Args {
		if v[1] == '-' {
			if v == "--hex" || v == "-h" {
				hexMode = true
			} else {
				usage()
			}
		} else if k > 0 {
			files = append(files, v)
		}
	}

	_ = hexMode
	if len(files) != 1 {
		usage()
	}

	inputFile := files[0]

	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read file: %v\n", err)
		os.Exit(1)
	}

	numEntries := int(data[0])
	offsets := make([]int, numEntries+1)
	for i := 0; i != numEntries+1; i++ {
		offsets[i] = int(data[1+i*3+2])<<16 + int(data[1+i*3+1])<<8 + int(data[1+i*3+0])
	}

	// Todo: Validate offsets[numEntries] == len(data)

	//binaryMode := true
	hexRun := 0
	stringRun := 0

	const HEX_RUN_LIMIT = 5
	begin := 1 + (numEntries+1)*3
	sectionEnd := begin
	for at := begin; ; at++ {
		if at == sectionEnd {
			if stringRun > 0 {
				// This is bad...
				fmt.Printf("\" NO_NUL\n") // I guess we just hack this and can handle it in recompiler...
			} else if hexRun > 0 {
				fmt.Printf("]\n")
			}
			for i := 0; i != numEntries; i++ {
				if offsets[i] == at {
					sectionEnd = offsets[i+1]
					fmt.Printf("SECTION %v\n", i)
				}
			}
			if at == len(data) {
				break
			}
			hexRun = 0
			stringRun = 0
		}

		if stringRun > 0 || hexRun == HEX_RUN_LIMIT {
			if stringRun == 0 && hexRun == HEX_RUN_LIMIT {
				fmt.Printf("] \"")
				hexRun = 0
			}
			if data[at] == 0x00 {
				fmt.Printf("\"\n")
				stringRun = 0
			} else {
				printByte(data[at] - 0x31)
				stringRun++
			}
		} else {
			if hexRun == 0 {
				fmt.Printf("[")
			} else {
				fmt.Printf(" ")
			}
			fmt.Printf("%02X", data[at])
			hexRun++
		}
	}

	fmt.Printf("\n")
}
