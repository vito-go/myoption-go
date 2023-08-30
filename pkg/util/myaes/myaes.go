package myaes

import (
	"crypto/aes"
	"crypto/cipher"
)

var key = []byte("despotic traitor")
var iv = append(make([]byte, aes.BlockSize-1), 5)

func EncOrDec(data []byte) []byte {
	return encOrDecWithKeyAndIV(data, key, iv)
}
func EncOrDecByKey(data []byte, k []byte) []byte {
	return encOrDecWithKeyAndIV(data, k, iv)
}

func encOrDecWithKeyAndIV(data []byte, key []byte, iv []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	ctr := cipher.NewCTR(block, iv)
	ctr.XORKeyStream(data, data)
	return data
}
func EncOrDecWithKeyAndIV(data []byte, key []byte, iv []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		// FIXME
		panic(err)
	}
	ctr := cipher.NewCTR(block, iv)
	ctr.XORKeyStream(data, data)
	return data
}
