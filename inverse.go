package main

import (
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
func c_perform_inverse(in *C.char, i_len C.int, o *C.char, o_len *C.int) {
	debug := os.Getenv("DEBUG")
	_err := func(err error) {
		if debug == "ok" {
			fmt.Println(err)
		}
		*o_len = C.int(0)
		// byteArrToCChar([]byte{0x00}, o)
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
	var resPadded [32]byte
	res := f.ToBytes(inv)
	copy(resPadded[32-len(res):], res[:])
	*o_len = C.int(len(res))
	outBuf := (*[32]byte)(unsafe.Pointer(o))
	length := copy(outBuf[:], resPadded[:])
	if length != len(resPadded) {
		fmt.Printf("copy was not successfull %d %d", length, len(resPadded))
	}
	// fmt.Printf("[go result]: %x\n", resPadded)
}

func main() {}
