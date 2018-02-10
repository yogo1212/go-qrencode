package qrencode

/*
#cgo CFLAGS: -I/usr/include/
#cgo LDFLAGS: -L/usr/lib/ -lqrencode
#include <errno.h>
#include <stdlib.h>
#include <string.h>
#include <qrencode.h>

const char *strerrno(void) {
	return strerror(errno);
}
*/
import "C"

import (
	"errors"
	"unsafe"
)

type QRCode struct {
	Version int
	Lines [][]bool
}

type QRECLevel C.QRecLevel

const (
	QR_ECLEVEL_L QRECLevel = C.QR_ECLEVEL_L
	QR_ECLEVEL_M           = C.QR_ECLEVEL_M
	QR_ECLEVEL_Q           = C.QR_ECLEVEL_Q
	QR_ECLEVEL_H           = C.QR_ECLEVEL_H
)

type QRMode C.QRencodeMode

const (
  QR_MODE_NUM QRMode = C.QR_MODE_NUM
  QR_MODE_AN         = C.QR_MODE_AN
  QR_MODE_8          = C.QR_MODE_8
  QR_MODE_KANJI      = C.QR_MODE_KANJI
  QR_MODE_STRUCTURE  = C.QR_MODE_STRUCTURE
  QR_MODE_ECI        = C.QR_MODE_ECI
  QR_MODE_FNC1FIRST  = C.QR_MODE_FNC1FIRST
  QR_MODE_FNC1SECOND = C.QR_MODE_FNC1SECOND
)

var qrCaseSensitive = map[bool]C.int {
	false: 0,
	true: 1,
}

func qrmeh(q *C.QRcode) (qr QRCode) {
	qr.Version = int(q.version)
	qr.Lines = make([][]bool, q.width)

	b := C.GoBytes(unsafe.Pointer(q.data), q.width*q.width)

	for i, _ := range qr.Lines {
		for _, x := range b[0:q.width] {
			qr.Lines[i] = append(qr.Lines[i], (x & 1) == 1)
		}
		b = b[q.width:]
	}

	return
}

func QRCodeEncodeString(what string, version int, level QRECLevel, hint QRMode, caseSensitive bool) (qr QRCode, err error) {
	s := C.CString(what)
	defer C.free(unsafe.Pointer(s))
	q := C.QRcode_encodeString(s, C.int(version), C.QRecLevel(level), C.QRencodeMode(hint), qrCaseSensitive[caseSensitive])
	if q == nil {
		err = errors.New(C.GoString(C.strerrno()))
		return
	}
	defer C.QRcode_free(q)

	qr = qrmeh(q)
	return
}

func QRCodeEncodeStringMQR(what string, version int, level QRECLevel, hint QRMode, caseSensitive bool) (qr QRCode, err error) {
	s := C.CString(what)
	defer C.free(unsafe.Pointer(s))
	q := C.QRcode_encodeStringMQR(s, C.int(version), C.QRecLevel(level), C.QRencodeMode(hint), qrCaseSensitive[caseSensitive])
	if q == nil {
		err = errors.New(C.GoString(C.strerrno()))
		return
	}
	defer C.QRcode_free(q)

	qr = qrmeh(q)
	return
}