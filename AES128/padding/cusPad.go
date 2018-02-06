package padding

import "bytes"

type cusPadder struct {
	padChr   []byte
	blockLen int
}

func NewCusPadder(blen int, pch []byte) *cusPadder {
	return &cusPadder{padChr: pch, blockLen: blen}
}

func (se *cusPadder) GoPad(src []byte) []byte {
	return Padding(src, se.padChr, se.blockLen)
}
func (se *cusPadder) UnPad(src []byte) []byte {
	dd := bytes.SplitAfter(src, se.padChr)
	return dd[0]
}
