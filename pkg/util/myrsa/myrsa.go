package myrsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"
)

// GenRasKeyPKCS8PriPKIXPub 生成特定编码的key.
func GenRasKeyPKCS8PriPKIXPub(bits int) (priKey, pubKey string, err error) {
	// 生成私钥
	keyPair, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", err
	}
	// X509序列化
	//  x509.MarshalPKCS1PrivateKey
	// This kind of key is commonly encoded in PEM blocks of type "RSA PRIVATE KEY".
	// der := x509.MarshalPKCS1PrivateKey(keyPair)
	//  x509.MarshalPKCS8PrivateKey:
	// This kind of key is commonly encoded in PEM blocks of type "PRIVATE KEY".
	der, err := x509.MarshalPKCS8PrivateKey(keyPair)
	if err != nil {
		return "", "", err
	}
	// pem 块
	block := pem.Block{
		Type:    "PRIVATE KEY",
		Headers: nil, //
		Bytes:   der,
	}
	var priKeyBuilder strings.Builder
	var pubKeyBuilder strings.Builder
	// pem 编码
	err = pem.Encode(&priKeyBuilder, &block)
	if err != nil {
		return "", "", err
	}
	// x509序列化
	// x509.MarshalPKIXPublicKe:flutter user this.
	//This kind of key is commonly encoded in PEM blocks of type "PUBLIC KEY".
	//
	der, err = x509.MarshalPKIXPublicKey(&keyPair.PublicKey)
	////  x509.MarshalPKCS1PublicKey：
	//// This kind of key is commonly encoded in PEM blocks of type "RSA PUBLIC KEY".
	//der = x509.MarshalPKCS1PublicKey(&keyPair.GroupPubKey)
	//if err != nil {
	//	return "", "", err
	//}
	block = pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil, //
		Bytes:   der,
	}
	err = pem.Encode(&pubKeyBuilder, &block)
	if err != nil {
		return "", "", err
	}
	priKey = priKeyBuilder.String()
	pubKey = pubKeyBuilder.String()
	return
}

func EncWithPKIXPublicKey(publicKey string, src []byte) ([]byte, error) {
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return nil, errors.New("nil block")
	}
	//pubKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	// 可以与flutter交互的公钥使用 x509.MarshalPKIXPublicKey序列化，
	// 因此这里用 x509.ParsePKIXPublicKey反序列化
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pubKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("publicKey type error. it is not *rsa.GroupPubKey")
	}
	result, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, src)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func CheckPKIXPublicKey(publicKey string) bool {
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return false

	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false
	}
	if _, ok := pub.(*rsa.PublicKey); !ok {
		return false

	}
	return true
}
func CheckPKCS8PrivateKey(privateKey string) bool {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return false
	}
	priKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return false
	}
	_, ok := priKey.(*rsa.PrivateKey)
	if !ok {
		return false
	}
	return true
}

// DecPKCS8PrivateKey flutter发来的用这个进行解密.
func DecPKCS8PrivateKey(privateKey string, encData []byte) ([]byte, error) {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, errors.New("nil block")
	}
	priKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pKey, ok := priKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("not *rsa.PrivateKey")
	}
	//result, err := rsa.DecryptPKCS1v15(rand.Reader, pKey, encData)
	result, err := rsa.DecryptPKCS1v15(rand.Reader, pKey, encData)
	if err != nil {
		panic(err)
		return nil, err
	}
	return result, nil
}

func encWithPKCS1PublicKey(publicKey string, src []byte) ([]byte, error) {
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return nil, errors.New("nil block")
	}
	pubKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	result, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, src)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func decWithPKCS1PrivateKey(privateKey string, encData []byte) ([]byte, error) {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, errors.New("nil block")
	}
	priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	result, err := rsa.DecryptPKCS1v15(rand.Reader, priKey, encData)
	if err != nil {
		panic(err)
		return nil, err
	}
	return result, nil
}
