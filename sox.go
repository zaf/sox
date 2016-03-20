/*
	Sound file conversion using soxlib

	Copyright (C) 2016, Lefteris Zafiris <zaf.000@gmail.com>

	This program is free software, distributed under the terms of
	the BSD 3-Clause License. See the LICENSE file
	at the top of the source tree.

*/

package sox

/*
#cgo CFLAGS: -O2 -march=native
#cgo LDFLAGS: -lsox

#include "convert.h"
*/
import "C"

import (
	"fmt"
	"unsafe"
)

//Init Initialise sox
func Init() (err error) {
	if int(C.sox_init()) != 0 {
		err = fmt.Errorf("Failed to initialise sox")
	}
	return
}

//FormatInit Initialise sox_formats
func FormatInit() (err error) {
	if int(C.sox_format_init()) != 0 {
		err = fmt.Errorf("Failed to initialise sox formats")
	}
	return
}

//Quit sox
func Quit() (err error) {
	if int(C.sox_quit()) != 0 {
		err = fmt.Errorf("Failed to quit sox")
	}
	return
}

//FormatQuit quits sox_formats
func FormatQuit() {
	C.sox_format_quit()
}

//Convert wav files
func Convert(inputData []byte, outputFormat string) (outData []byte, err error) {
	var sndIn, sndOut C.snd_file
	format := C.CString(outputFormat)
	sndIn.buff = unsafe.Pointer(&inputData[0])
	sndIn.size = (C.size_t)(len(inputData))

	sndOut = C.convert_snd(&sndIn, format)
	if (uint)(sndOut.size) == 0 {
		err = fmt.Errorf("Failed to convert sound data")
	}
	outData = C.GoBytes(sndOut.buff, (C.int)(sndOut.size))
	C.free(unsafe.Pointer(sndOut.buff))
	C.free(unsafe.Pointer(format))
	return
}
