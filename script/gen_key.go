package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	err := os.Chdir(filepath.Dir(file))
	if err != nil {
		panic(err)
	}
}

//go:generate go run gen_key.go
func main() {
	log.Println("gen key which is used for https server")
	keysDir := filepath.Dir("../configs/keys/")
	if err := os.MkdirAll(keysDir, 0755); err != nil {
		panic(err)
	}
	err := genHTTPSKeyPair(filepath.Join(keysDir, "server_cert.pem"), filepath.Join(keysDir, "server_key.pem"), 2048)
	if err != nil {
		panic(err)
	}
	log.Println("gen key success,keys directory: ", keysDir)
}

//x509.Certificate结构体中的字段代表了证书的各种信息，下面是各个字段的含义：
//SerialNumber：证书序列号，应该是唯一的。
//Subject：证书的主题，即证书所代表的实体的信息。
//Issuer：证书的颁发者，即颁发该证书的实体的信息。
//NotBefore：证书的生效时间。
//NotAfter：证书的过期时间。
//KeyUsage：证书的使用场景，例如加密、数字签名等。
//ExtKeyUsage：扩展密钥用法，例如TLS网站身份验证、电子邮件身份验证等。
//BasicConstraintsValid：一个布尔值，表示是否启用基本约束。
//IsCA：一个布尔值，表示该证书是否是根证书。
//MaxPathLen：最大路径长度限制。
//SubjectKeyId：用于在CA层次结构中查找此证书的标识符。
//AuthorityKeyId：证书颁发者的标识符。
//DNSNames：一个包含该证书所代表的实体的DNS名称的列表。
//IPAddresses：一个包含该证书所代表的实体的IP地址的列表。

// genHTTPSKeyPair generates a key pair and a self-signed X.509 certificate for an HTTPS server.
func genHTTPSKeyPair(certFileName, keyFileName string, bits int) error {
	// Generate a new RSA key
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}

	// Define certificate fields
	subject := pkix.Name{
		CommonName:         "example.com",
		Country:            []string{"US"},
		Province:           []string{"California"},
		Locality:           []string{"San Francisco"},
		Organization:       []string{"Example, Inc."},
		OrganizationalUnit: []string{"IT"},
	}
	notBefore := time.Now()
	notAfter := notBefore.AddDate(1, 0, 0)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return err
	}

	// Create a new X.509 certificate template
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      subject,
		NotBefore:    notBefore,
		NotAfter:     notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		//DNSNames:              []string{"example.com", "www.example.com"},
		//IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
	}
	// Generate a self-signed X.509 certificate
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return err
	}
	// Encode the certificate and private key to PEM format
	certOut, err := os.Create(certFileName)
	if err != nil {
		return err
	}
	defer certOut.Close()
	if err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return err
	}
	keyOut, err := os.Create(keyFileName)
	if err != nil {
		panic(err)
	}
	defer keyOut.Close()
	if err = pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}); err != nil {
		return err
	}
	return nil
}
