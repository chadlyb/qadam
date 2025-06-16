package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func isPrintable(b byte) bool {
	return b >= 0x20 && b <= 0x7E
}

// ISO CP 862
const charset = " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~ ÇüéâäůćçłëŐőîŹÄĆÉĹĺôöĽľŚśÖÜŤťŁ×čáíóúĄąŽžĘę¬źČş«»░▒▓│┤ÁÂĚŞ╣║╗╝Żż┐└┴┬├─┼Ăă╚╔╩╦╠═╬¤đĐĎËďŇÍÎě┘┌█▄ŢŮ▀ÓßÔŃńňŠšŔÚŕŰýÝţ´-˝˛ˇ˘§÷¸°¨˙űŘř■ "

var tableRunes = []rune(charset)

func printByte(b byte) {
	switch {
	case b == '\n':
		fmt.Printf("\\n")
	case b == '\t':
		fmt.Printf("\\t")
	case b == '\\':
		fmt.Printf("\\\\")
	case b == ' ':
		fmt.Printf(" ")
	case b == '"':
		fmt.Printf("\\\"")
	case b < 0x20 || tableRunes[b-0x20] == ' ':
		fmt.Printf("\\x%02X", b)
	default:
		fmt.Printf("%c", tableRunes[b-0x20])
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: printable <input file>")
		os.Exit(1)
	}

	inputFile := os.Args[1]

	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read file: %v\n", err)
		os.Exit(1)
	}

	numEntries := int(data[0])
	offsets := make([]int, numEntries)
	for i := 0; i != numEntries; i++ {
		offsets[i] = int(data[1+i*3+2])<<16 + int(data[1+i*3+1])<<8 + int(data[1+i*3+0])
	}

	binaryMode := true
	begin := 1 + (numEntries+1)*3
	for at := begin; at < len(data); at++ {
		if binaryMode {
			for i := 0; i != numEntries; i++ {
				if offsets[i] == at {
					fmt.Printf("SECTION %v\n", i)
					break
				}
			}

			if data[at] == 0xFF || data[at] == 0x30 {
				fmt.Printf("[%02X]\n", data[at])
			} else {
				fmt.Printf("[%02X %02X %02X %02X %02X] \"", data[at], data[at+1], data[at+2], data[at+3], data[at+4])
				at += 4
				binaryMode = false
			}
			continue
		}

		if data[at] == 0x00 {
			fmt.Printf("\"\n")
			binaryMode = true
			continue
		}

		b := data[at] - 0x31

		printByte(b)
	}

	fmt.Printf("\n")
}
