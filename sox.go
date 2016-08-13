/*
	Sound file conversion using soxlib

	Copyright (C) 2016, Lefteris Zafiris <zaf@fastmail.com>

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
func Convert(inputData []byte, inputFormat string, outputFormat string) (outData []byte, err error) {
	var sndIn, sndOut C.snd_file
	inFormat := C.CString(inputFormat)
	outFormat := C.CString(outputFormat)
	sndIn.size = (C.size_t)(len(inputData))
	sndIn.buff = C.malloc(sndIn.size)
	cBuf := (*[1 << 30]byte)(sndIn.buff)
	//cBuf := (*[1 << 30]byte)(sndIn.buff)[:len(inputData)]
	//sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&cBuf))
	//sliceHeader.Cap = len(inputData)
	copy(cBuf[:], inputData)

	sndOut = C.convert_snd(&sndIn, inFormat, outFormat)
	if (uint)(sndOut.size) == 0 {
		err = fmt.Errorf("Failed to convert sound data")
	}
	outData = C.GoBytes(sndOut.buff, (C.int)(sndOut.size))
	C.free(unsafe.Pointer(sndIn.buff))
	C.free(unsafe.Pointer(sndOut.buff))
	C.free(unsafe.Pointer(inFormat))
	C.free(unsafe.Pointer(outFormat))
	return
}
