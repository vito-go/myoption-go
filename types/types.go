package types

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/sha3"
)

type Platform string

func (p Platform) Check() bool {
	switch p {
	case PlatformAndroid, PlatformWindows, PlatformIos, PlatformLinux:
		return true
	default:
		return false
	}
}

const (
	PlatformAndroid = "android"
	PlatformWindows = "windows"
	PlatformIos     = "ios"
	PlatformLinux   = "linux"
)

type LoginInfo struct {
	DeviceId   string `json:"platform"`
	Expire     int64  `json:"expire,omitempty"` // 過期毫秒時間戳
	UA         *UA    `json:"ua"`
	LoginToken string `json:"loginToken,omitempty"`

	LoginTime int64  `json:"loginTime"`
	IpAddress string `json:"ipAddress,omitempty"`
}

func NewLoginInfo(userId string, deviceId string, ua *UA) *LoginInfo {
	now := time.Now()
	expire := now.AddDate(0, 6, 0).UnixMilli()
	b, _ := json.Marshal(ua)
	rand.Seed(now.UnixNano())
	loginToken := fmt.Sprintf("%x",
		sha3.Sum256([]byte(fmt.Sprintf(
			"%s%d%s%d", userId, now.UnixNano(), string(b), rand.Int63()))))
	return &LoginInfo{
		LoginTime:  time.Now().UnixMilli(),
		Expire:     expire,
		UA:         ua,
		LoginToken: loginToken,
		DeviceId:   deviceId,
	}
}

// UA user-agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36
// myoption/<version> OsName/OsVersion deviceName/(DeviceInfo)
// MyChat/1.6.5 android/13 samsung/(SM-G9910)
type UA struct {
	UserAgent  string // 优先取UserAgent字段
	AppName    string
	Version    string
	OsName     Platform
	OsVersion  string
	DeviceName string
	DeviceInfo string
}
