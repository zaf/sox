/*
	Sound file conversion using sox

	usage: goconvert -o output.format [input-files]

	output.format can be one of the following:
	alaw, ulaw, gsm, wav, wav8, ogg, mp3, flac

	Copyright (C) 2016, Lefteris Zafiris <zaf@fastmail.com>

	This program is free software, distributed under the terms of
	the BSD 3-Clause License. See the LICENSE file
	at the top of the source tree.

*/

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/zaf/sox"
)

var out = flag.String("o", "", "Output format")

// Read a sound file from disk and convert it to another format.
func main() {
	flag.Parse()
	if *out == "" {
		log.Fatalln("No output format specified")
	}
	outFormat := *out

	if len(flag.Args()) < 1 {
		log.Fatalln("No input files specified")
	}

	err := sox.Init()
	if err != nil {
		log.Fatal(err)
	}

	err = sox.FormatInit()
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range flag.Args() {
		var err error
		var input, output []byte
		input, err = ioutil.ReadFile(file)
		if err != nil {
			log.Println(err)
			continue
		}
		inFormat := filepath.Ext(file)
		if len(inFormat) > 1 {
			inFormat = inFormat[1:]
		} else {
			log.Printf("not able to determine input format from file extension for %s\n", file)
			continue
		}
		output, err = sox.Convert(input, inFormat, outFormat)
		if err != nil {
			log.Println(err)
			continue
		}
		outputFile := strings.TrimSuffix(file, filepath.Ext(file)) + "." + outFormat
		err = ioutil.WriteFile(outputFile, output, 0644)
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Printf("Saved %s\n", outputFile)
	}
	sox.FormatQuit()
	sox.Quit()
}
