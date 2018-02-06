package padding

import (
	"bytes"
	"fmt"
)

type Padder interface {
	GoPad(src []byte) []byte
	UnPad(src []byte) []byte
}

func Padding(src []byte, padchr []byte, blocklen int) []byte {
	srclen := len(src)
	padlen := blocklen - (srclen % blocklen)
	dst := bytes.NewBuffer(src)
	dst.Write(bytes.Repeat(padchr, padlen))
	fmt.Printf("%v %s\n", len(dst.Bytes()), dst.Bytes())
	return dst.Bytes()
}
