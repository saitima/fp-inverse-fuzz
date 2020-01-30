package main

import (
	"C"
	"encoding/hex"
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

func byteArrToCChar(in []byte, out *C.char) {
	inStr := hex.EncodeToString(in)
	inCharArr := C.CString(inStr)
	inPtr := uintptr(unsafe.Pointer(inCharArr))
	outPtr := uintptr(unsafe.Pointer(out))
	for i := 0; i < len(inStr); i++ {
		inElem := (*C.char)(unsafe.Pointer(inPtr))
		outElem := (*C.char)(unsafe.Pointer(outPtr))
		*outElem = *inElem
		inPtr++
		outPtr++
	}
}

//export c_perform_inverse
func c_perform_inverse(in *C.char, i_len C.int, o *C.char, o_len *C.int) {
	debug := os.Getenv("DEBUG")
	_err := func(err error) {
		if debug == "ok" {
			fmt.Println(err)
		}
		*o_len = C.int(0)
		byteArrToCChar([]byte{0x00}, o)
		return
	}
	buf := C.GoBytes(unsafe.Pointer(in), C.int(i_len))
	modLen := int(buf[0])
	modulusBuf, rest, err := split(buf[1:], modLen)
	if err != nil {
		_err(err)
		return
	}
	elemLen := int(rest[0])
	elemBuf, rest, err := split(rest[1:], elemLen)
	if err != nil {
		_err(err)
		return
	}
	if debug == "on" {
		fmt.Printf("modLen: %x\n", modLen)
		fmt.Printf("modulusBuf: %x\n", modulusBuf)
		fmt.Printf("elemLen: %x\n", elemLen)
		fmt.Printf("elemBuf: %x\n", elemBuf)
	}
	if len(modulusBuf) < 8 {
		modulusBuf = padBytes(modulusBuf, 8)
	}
	if len(elemBuf) < 8 {
		elemBuf = padBytes(elemBuf, 8)
	}

	f, err := fp.NewField(modulusBuf)
	if err != nil {
		_err(err)
		return
	}

	elem, err := f.NewFieldElementFromBytes(elemBuf)
	if err != nil {
		_err(err)
		return
	}
	inv := f.NewFieldElement()
	if ok := f.Inverse(inv, elem); !ok {
		_err(fmt.Errorf("element has no inverse"))
		return
	}
	res := f.ToBytes(inv)
	*o_len = C.int(len(res))
	byteArrToCChar(res, o)
}

func main() {}
