package padding

import (
	"fmt"
	"strconv"
)

type p7Padder struct {
	blockLen int
}

func NewP7Padder(blen int) *p7Padder {
	return &p7Padder{blockLen: blen}
}

func (se *p7Padder) GoPad(src []byte) []byte {
	sl := len(src)
	ll := se.blockLen - (sl % se.blockLen)
	if ll == 16 {
		return Padding(src, []byte("g"), se.blockLen)
	}
	h := fmt.Sprintf("%x", ll)
	return Padding(src, []byte(h), se.blockLen)
}
func (se *p7Padder) UnPad(src []byte) []byte {
	if string(src[len(src)-1]) == "g" {
		return src[:len(src)-16]
	}
	nu, _ := strconv.ParseUint(string(src[len(src)-1]), 16, 64)
	return src[:len(src)-int(nu)]
}
