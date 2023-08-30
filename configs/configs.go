package configs

import (
	"context"
	"crypto/ecdh"
	"crypto/sha1"
	_ "embed"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"github.com/vito-go/mylog"
	"gopkg.in/yaml.v3"
)

var x25519PubKeyNoPriKeyMap = make(map[uint32]string, 10)

func PrivateKeyByPubKey(pubKeyNo uint32) (*ecdh.PrivateKey, bool) {
	priKey, ok := x25519PubKeyNoPriKeyMap[pubKeyNo]
	if !ok {
		return nil, false
	}
	privateKeyBytes, _ := base64.StdEncoding.DecodeString(priKey)
	curve := ecdh.X25519()
	privateKeyA, _ := curve.NewPrivateKey(privateKeyBytes)
	return privateKeyA, true
}

func init() {
	var x25519KeyPairsMap map[string]string
	err := yaml.Unmarshal(x25519KeyPairsBytes, &x25519KeyPairsMap)
	if err != nil {
		panic(err)
	}
	for pubKey, priKey := range x25519KeyPairsMap {
		sum := sha1.Sum([]byte(pubKey))
		pubKeyNo := binary.BigEndian.Uint32(sum[:4])
		x25519PubKeyNoPriKeyMap[pubKeyNo] = priKey
	}
	mylog.Ctx(context.Background()).Info("x25519KeyPairsMap", x25519PubKeyNoPriKeyMap)

}

//go:embed symbol_codes.yaml
var symbolCodesBytes []byte

var Symbols []string
var symbolCodeNameMap = make(map[string]string, 10)

type symbolCodeName struct {
	SymbolCode string `yaml:"symbolCode"`
	SymbolName string `yaml:"symbolName"`
}

func SymbolNameByCode(symbolCode string) string {
	return symbolCodeNameMap[symbolCode]

}
func init() {
	var result []symbolCodeName
	err := yaml.Unmarshal(symbolCodesBytes, &result)
	if err != nil {
		panic(err)
	}
	for _, v := range result {
		if v.SymbolCode == "" {
			continue
		}
		Symbols = append(Symbols, v.SymbolCode)
		if symbolCodeNameMap[v.SymbolCode] != "" {
			panic(fmt.Sprintf("duplicate symbolCode: %s in symbol_codes.yaml", v.SymbolCode))
		}
		symbolCodeNameMap[v.SymbolCode] = v.SymbolName
	}
	if len(Symbols) == 0 {
		panic("no symbolCode in symbol_codes.yaml")
	}
	mylog.Ctx(context.Background()).Info("Symbols", Symbols)
	mylog.Ctx(context.Background()).Info("symbolCodeNameMap", symbolCodeNameMap)
}
