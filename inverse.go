package main

import (
	"encoding/hex"
	//#include "stdlib.h"
	"C"
	"fmt"
	"os"
	"unsafe"

	fp "github.com/saitima/eip-fp"
)

func padBytes(in []byte, size int) []byte {
	out := make([]byte, size)
	if len(in) > size {
		panic("bad input for padding")
	}
	copy(out[size-len(in):], in)
	return out
}

func split(in []byte, offset int) ([]byte, []byte, error) {
	if len(in) < offset {
		return nil, nil, fmt.Errorf("cant split at given offset %d", offset)
	}
	return in[:offset], in[offset:], nil
}

//export c_perform_inverse
func c_perform_inverse(in *C.char, i_len C.int, o_len *C.int) *C.char {
	debug := os.Getenv("DEBUG")
	_err := func(err error) *C.char {
		if debug == "ok" {
			fmt.Println(err)
		}
		*o_len = C.int(0)
		return C.CString("")
	}
	buf := C.GoBytes(unsafe.Pointer(in), C.int(i_len))
	modLen := int(buf[0])
	modulusBuf, rest, err := split(buf[1:], modLen)
	if err != nil {
		return _err(err)
	}
	elemLen := int(rest[0])
	elemBuf, rest, err := split(rest[1:], elemLen)
	if err != nil {
		return _err(err)
	}
	// if debug == "on" {
	// 	fmt.Printf("modLen: %x\n", modLen)
	// 	fmt.Printf("modulusBuf: %x\n", modulusBuf)
	// 	fmt.Printf("elemLen: %x\n", elemLen)
	// 	fmt.Printf("elemBuf: %x\n", elemBuf)
	// }
	if len(modulusBuf) < 8 {
		modulusBuf = padBytes(modulusBuf, 8)
	}
	if len(elemBuf) < 8 {
		elemBuf = padBytes(elemBuf, 8)
	}

	f, err := fp.NewField(modulusBuf)
	if err != nil {
		return _err(err)
	}

	elem, err := f.NewFieldElementFromBytes(elemBuf)
	if err != nil {
		return _err(err)
	}
	inv := f.NewFieldElement()
	if ok := f.Inverse(inv, elem); !ok {
		return _err(fmt.Errorf("element has no inverse"))
	}
	res := f.ToBytes(inv)
	*o_len = C.int(len(res))
	return C.CString(hex.EncodeToString(res))
}

func main() {}
