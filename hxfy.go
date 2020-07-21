/*
Copyright 2020, aeter

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

/*
Black        0;30     Dark Gray     1;30
Red          0;31     Light Red     1;31
Green        0;32     Light Green   1;32
Brown/Orange 0;33     Yellow        1;33
Blue         0;34     Light Blue    1;34
Purple       0;35     Light Purple  1;35
Cyan         0;36     Light Cyan    1;36
Light Gray   0;37     White         1;37
*/

const (
	Reset    = "\033[0m"
	Black    = "\033[0;31m%s" + Reset
	DarkGray = "\033[1;30m%s" + Reset
	Green    = "\033[1;32m%s" + Reset
	Yellow   = "\033[0;33m%s" + Reset
	Cyan     = "\033[1;36m%s" + Reset
)

func usage() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <filename>\n", os.Args[0])
		os.Exit(1)
	}
}

func byteScanner() *bufio.Scanner {
	filename, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sc := bufio.NewScanner(filename)
	sc.Split(bufio.ScanBytes)
	return sc
}

func identify(b []byte) (color, letter string) {
	if b[0] == 0 { // NULL byte
		color = DarkGray
		letter = fmt.Sprintf(color, "0")
	} else if b[0] > 126 { // nonascii
		color = Yellow
		letter = fmt.Sprintf(color, ".")
	} else if b[0] < 33 { // ascii non-printable
		color = Green
		letter = fmt.Sprintf(color, "_")
	} else { // ascii printable
		color = Cyan
		letter = fmt.Sprintf(color, b)
	}
	return
}

func main() {
	usage()

	sc := byteScanner()

	i := 0
	// for efficiency we use strings.Builder
	var str strings.Builder
	var hexStr strings.Builder
	for sc.Scan() {
		color, letter := identify(sc.Bytes())
		hexStr.WriteString(fmt.Sprintf(color, hex.EncodeToString(sc.Bytes())))
		hexStr.WriteString(" ")
		str.WriteString(letter)

		i++

		if i%16 == 0 {
			fmt.Printf("%s ", hexStr.String())
			hexStr.Reset()

			fmt.Printf("%s\n", str.String())
			str.Reset()
		}
	}

	// equalizing columns at the last line
	for ; i%16 != 0; i++ {
		hexStr.WriteString(fmt.Sprintf(Black, "  "))
		hexStr.WriteString(" ")
	}

	fmt.Printf("%s ", hexStr.String())
	fmt.Printf("%s\n", str.String())
	fmt.Println(Reset)
}
