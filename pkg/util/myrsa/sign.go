package myrsa

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

// RsaVerySignWithSha256 验证
func RsaVerySignWithSha256(publicKey string, data, signData []byte) (bool, error) {
	return rsaVerySignWithHash(publicKey, data, signData, crypto.SHA256)
}

func rsaVerySignWithHash(publicKey string, data, signData []byte, h crypto.Hash) (bool, error) {
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return false, errors.New("nil block")
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, err
	}
	H := h.New()
	H.Write(data)
	hashed := H.Sum(nil)

	err = rsa.VerifyPKCS1v15(pubKey.(*rsa.PublicKey), h, hashed, signData)
	if err != nil {
		return false, nil
	}
	return true, nil
}
