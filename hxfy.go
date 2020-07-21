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
	reset    = "\033[0m"
	black    = "\033[0;31m%s" + reset
	darkGray = "\033[1;30m%s" + reset
	green    = "\033[1;32m%s" + reset
	yellow   = "\033[0;33m%s" + reset
	cyan     = "\033[1;36m%s" + reset
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
		color = darkGray
		letter = fmt.Sprintf(color, "0")
	} else if b[0] > 126 { // nonascii
		color = yellow
		letter = fmt.Sprintf(color, ".")
	} else if b[0] < 33 { // ascii non-printable
		color = green
		letter = fmt.Sprintf(color, "_")
	} else { // ascii printable
		color = cyan
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
		hexStr.WriteString(fmt.Sprintf(black, "  "))
		hexStr.WriteString(" ")
	}

	fmt.Printf("%s ", hexStr.String())
	fmt.Printf("%s\n", str.String())
	fmt.Println(reset)
}
