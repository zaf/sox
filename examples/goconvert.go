/*
	Sound file conversion using sox

	usage: goconvert input.file output.file output.format

	output.format can be one of the following:
	alaw, ulaw, gsm, wav, wav8, ogg, mp3, flac

	Copyright (C) 2016, Lefteris Zafiris <zaf.000@gmail.com>

	This program is free software, distributed under the terms of
	the BSD 3-Clause License. See the LICENSE file
	at the top of the source tree.

*/

package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/zaf/sox"
)

// Read a wav file from disk and convert it to another format.
func main() {
	if len(os.Args) < 4 {
		log.Fatalln("not enough parameters")
	}
	inputFile := os.Args[1]
	outputFile := os.Args[2]
	outFormat := os.Args[3]
	var err error
	var input, output []byte

	input, err = ioutil.ReadFile(inputFile)
	if err != nil {
		log.Fatalln(err)
	}

	err = sox.Init()
	if err != nil {
		log.Fatal(err)
	}

	err = sox.FormatInit()
	if err != nil {
		log.Fatal(err)
	}

	// libsox is not threaad safe, never use it in go-rourunes or similar context.
	output, err = sox.Convert(input, outFormat)
	if err != nil {
		log.Fatalln(err)
	}

	sox.FormatQuit()
	sox.Quit()

	err = ioutil.WriteFile(outputFile, output, 0644)
	if err != nil {
		log.Fatalln(err)
	}
}
