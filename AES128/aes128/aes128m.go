package aes128

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	"../padding"
)

func NewAES128Cry(kk []byte, vv []byte) *AES128Cry {
	return &AES128Cry{key: kk, iv: vv}
}

type AES128Cry struct {
	key []byte
	iv  []byte
}

func (aa *AES128Cry) encrypt(plaintext []byte) []byte {

	// CBC mode works on blocks so plaintexts may need to be padded to the
	// next whole block. For an example of such padding, see
	// https://tools.ietf.org/html/rfc5246#section-6.2.3.2. Here we'll
	// assume that the plaintext is already of the correct length.
	if len(plaintext)%aes.BlockSize != 0 {
		panic("plaintext is not a multiple of the block size")
	}

	block, err := aes.NewCipher(aa.key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	copy(ciphertext[:aes.BlockSize], aa.iv)
	// iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, aa.iv); err != nil {
		panic(err)
	}

	mode := cipher.NewCBCEncrypter(block, aa.iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.

	fmt.Printf("%x\n", ciphertext)
	return ciphertext
}

func (aa *AES128Cry) decrypt(ciphertext []byte) []byte {
	// key := []byte("example key 1234")
	// ciphertext, _ := hex.DecodeString("f363f3ccdcb12bb883abf484ba77d9cd7d32b5baecb3d4b1b3e0e4beffdb3ded")

	block, err := aes.NewCipher(aa.key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	copy(ciphertext[:aes.BlockSize], aa.iv)
	// iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// CBC mode always works in whole blocks.
	if len(ciphertext)%aes.BlockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, aa.iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(ciphertext, ciphertext)

	// If the original plaintext lengths are not a multiple of the block
	// size, padding would have to be added when encrypting, which would be
	// removed at this point. For an example, see
	// https://tools.ietf.org/html/rfc5246#section-6.2.3.2. However, it's
	// critical to note that ciphertexts must be authenticated (i.e. by
	// using crypto/hmac) before being decrypted in order to avoid creating
	// a padding oracle.

	fmt.Printf("%s\n", ciphertext)
	return ciphertext
}

func main() {
	key := []byte("1234567890098765")
	iv := []byte("qwertyuioppoiuyt")
	// cc := padding.NewCusPadder(aes.BlockSize, []byte("="))
	cc := padding.NewP7Padder(aes.BlockSize)
	aa := NewAES128Cry(key, iv)
	plain := cc.GoPad([]byte("I have a dream"))
	en := aa.encrypt(plain)
	de := aa.decrypt(en)
	back := cc.UnPad(de)
	fmt.Printf("%s\n", back)

}
