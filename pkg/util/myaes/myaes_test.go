package myaes

import (
	"crypto/aes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
)

func TestEncOrDec1(t *testing.T) {
	var sharedKey = []byte{185, 238, 219, 21, 195, 251, 161, 4, 142, 120, 172, 90, 189, 232, 247, 79, 105, 86, 15, 139, 145, 55, 136, 49, 168, 253, 46, 234, 185, 46, 207, 76}
	iv := make([]byte, aes.BlockSize)
	io.ReadFull(rand.Reader, iv)

	result := encOrDecWithKeyAndIV([]byte("nihao"), sharedKey, sharedKey[0:16])
	fmt.Println(len(result))
	fmt.Println(result)
}
func TestEncOrDec(t *testing.T) {
	sha1bytes, err := hex.DecodeString("2a2f5c437c27ffb86a352e09726834947266d6f7")
	if err != nil {
		panic(err)
	}
	aesKey := sha1bytes[:16]
	cli := http.Client{}
	req, err := http.NewRequest("GET", `http://192.168.89.64:9070/resource/view/0a4b43dc5f4c21de6627bad7768ffe3091d31650`, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("range", "bytes=0-4365174")
	resp, err := cli.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println(resp.Header)
	bb, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(bb))
	os.WriteFile("/home/vito/go/src/vitogo.tpddns.cn/liushihao/myoption-go/aaa", EncOrDecByKey(bb, aesKey), 0644)
}
