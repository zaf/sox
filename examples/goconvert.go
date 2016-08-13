/*
	Sound file conversion using sox

	usage: goconvert input.file output.format

	output.format can be one of the following:
	alaw, ulaw, gsm, wav, wav8, ogg, mp3, flac

	Copyright (C) 2016, Lefteris Zafiris <zaf@fastmail.com>

	This program is free software, distributed under the terms of
	the BSD 3-Clause License. See the LICENSE file
	at the top of the source tree.

*/

package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/zaf/sox"
)

// Read a sound file from disk and convert it to another format.
func main() {
	if len(os.Args) < 3 {
		log.Fatalln("not enough parameters")
	}
	var err error
	var input, output []byte
	inputFile := os.Args[1]
	outFormat := os.Args[2]

	input, err = ioutil.ReadFile(inputFile)
	if err != nil {
		log.Fatalln(err)
	}

	inFormat := filepath.Ext(inputFile)
	if len(inFormat) > 1 {
		inFormat = inFormat[1:]
	} else {
		log.Fatalln("not able to determine input format from file extension")
	}
	outputFile := strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + "." + outFormat

	err = sox.Init()
	if err != nil {
		log.Fatal(err)
	}

	err = sox.FormatInit()
	if err != nil {
		log.Fatal(err)
	}

	// libsox is not thread safe, never use it in go-rourunes or similar context.
	output, err = sox.Convert(input, inFormat, outFormat)
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
