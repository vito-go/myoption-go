package handler

import (
	"bytes"
	"context"
	"crypto/ecdh"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"myoption/configs"
	"myoption/pkg/util/myaes"
	"myoption/types"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/vito-go/mylog"

	"myoption/pkg/resp"
)

// Handler 实现Controller接口即可添加路由
// 除了实现接口的四个函数，顺序写在每一个具体的controller文件中。 、
// 其他功能函数，例如组装结果composeUserInfo等一律不可导出使用小写。
type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)            //  提供添加路由的接口函数
	GetParam(r *http.Request) (*ReqParam, error)                 // GetParam 校验以及获取参数
	GetRespBody(ctx context.Context, p *ReqParam) *resp.HTTPBody // GetRespBody 获取需要响应的httpBody. 重点聚焦在这个函数的实现
}

//     "X-IV": base64Encode(IV),
//    "X-PubKey-Number": serverPubKeyNoMix,
//    "X-Client-TimeStamp": xTime,
//    "X-Client-PubKey": base64Encode(alicePubKey.bytes),
//    "X-Sign": base64Encode(signature.bytes)

// ReqParam 请求的参数接口.
type ReqParam struct {
	QueryForm url.Values // QueryForm 存储从前端获取到的一些参数
	Header    http.Header
	Body      []byte
	// 这里可以添加其他项目可能需要的字段，尽管加 向前兼容
	// Auth Auth
	UserId string // from header X-User
	UA     *types.UA
	XIV    []byte
	//XPubKeyNumber string
	ServePriKey     *ecdh.PrivateKey
	ClientPubKey    *ecdh.PublicKey
	DeviceId        string
	SignPubKeyBytes ed25519.PublicKey
	xClientTime     int64
	Sign            []byte
}

// Get 获取参数值.
func (r *ReqParam) Get(key string) string {
	return r.QueryForm.Get(key)
}

// Set 设定参数值.
func (r *ReqParam) Set(key string, value string) {
	if r.QueryForm == nil {
		r.QueryForm = make(url.Values)
	}
	r.QueryForm.Set(key, value)
}

// ServeHTTP 提供添加路由的接口函数.一个完整的路由请求函数。
// 响应时间、哨兵监控
// postWithAes  {"code":0,"message":"","data":{}}
// http code==200,请求正确 返回内容加密，解密密钥 pwdSha1 [0:20]
// http code !=200, 请求错误，返回内容明文。
func ServeHTTP(w http.ResponseWriter, r *http.Request, h Handler) {
	tid := r.Context().Value("tid")
	ctx := context.WithValue(context.Background(), "tid", tid)
	startTime := time.Now()
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false) //

	reqParam, err := h.GetParam(r)
	if err != nil {
		mylog.Ctx(ctx).Error(err)
		respBody := resp.ErrParam(ctx)
		w.Header().Set("content-type", "application/json")
		_ = encoder.Encode(respBody)
		_, _ = w.Write(buf.Bytes())
		return
	}
	respBody := h.GetRespBody(ctx, reqParam)
	w.Header().Set("content-type", "application/json")
	_ = encoder.Encode(respBody)
	_, _ = w.Write(buf.Bytes())
	mylog.Ctx(ctx).WithFields("RT", time.Since(startTime).String(),
		"remoteAddr", r.RemoteAddr, "method", r.Method,
		"path", r.URL.Path, "header", r.Header,
		"query", reqParam.Body, "respBody", respBody).Info("")
}

// GetParam 支持 get， post（www/urlEncode）.
func GetParam(r *http.Request, keys ...string) (*ReqParam, error) {
	return checkAndGetParam(r, nil, keys...)
}
func GetParam1(r *http.Request, keys ...string) (*ReqParam, error) {
	param, err := getParamEnc(r)
	if err != nil {
		return nil, err
	}
	for _, key := range keys {
		if param.QueryForm.Get(key) == "" {
			return nil, errors.New("param not found: " + key)
		}
	}
	return param, err
}

func GetParamWithDefault1(r *http.Request, defaultParamMap DefaultParamMap, keys ...string) (*ReqParam, error) {
	param, err := getParamEnc(r)
	if err != nil {
		return nil, err
	}
	for k, v := range defaultParamMap {
		if !param.QueryForm.Has(k) {
			param.QueryForm.Set(k, v)
		}
	}
	for _, key := range keys {
		if param.QueryForm.Get(key) == "" {
			return nil, errors.New("param not found: " + key)
		}
	}
	return param, err
}

// DefaultParamMap 非必须的参数 key为必传参数名称，value为默认参数值
type DefaultParamMap = map[string]string

// GetParamWithDefault 通过 ctx.Query 方法获取参数 defaultParamMap 传入非必须的参数 keys为必传参数
func GetParamWithDefault(r *http.Request, d DefaultParamMap, keys ...string) (*ReqParam, error) {
	return getParamWithDefaultParam(r, d, keys...)
}

//     "X-IV": base64Encode(IV),
//    "X-PubKey-Number": serverPubKeyNoMix,
//    "X-Client-TimeStamp": xTime,
//    "X-Client-PubKey": base64Encode(alicePubKey.bytes),
//    "X-Sign": base64Encode(signature.bytes)

// GetParamBody 通过 ctx.Query 方法获取参数 defaultParamMap 传入非必须的参数 keys为必传参数
func GetParamBody(r *http.Request) (*ReqParam, error) {
	encBody, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	_ = encBody
	return nil, nil
}

var regUA = regexp.MustCompile(`(.+?)/(\d+\.\d+\.\d+) (.+?)/(.+?) (.+?)/\((.+?)\)`)

func getUAFromUserAgent(ua string) (*types.UA, bool) {
	result := regUA.FindAllStringSubmatch(ua, -1)
	mylog.Ctx(context.Background()).Infof("%+v", result)
	if len(result) > 0 {
		osName := types.Platform(strings.ToLower(result[0][3]))
		if !osName.Check() {
			return nil, false
		}
		return &types.UA{
			AppName:    result[0][1],
			Version:    result[0][2],
			OsName:     osName,
			OsVersion:  result[0][4],
			DeviceName: result[0][5],
			DeviceInfo: result[0][6],
		}, true

	}
	return nil, false
}

// GetParamBodyEnc .
// Deprecated please use GetParamBody
// List the field in the header:
//
//	X-IV, X-PubKey-Number, X-Client-TimeStamp, X-Client-PubKey, X-Sign, X-User-Agent, X-User
func GetParamBodyEnc(r *http.Request) (*ReqParam, error) {
	return getParamBodyEnc(r)
}
func getParamBodyEnc(r *http.Request) (*ReqParam, error) {
	encBody, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	iv, err := base64.StdEncoding.DecodeString(r.Header.Get("X-IV"))
	if err != nil {
		return nil, err
	}
	if len(iv) != 16 {
		return nil, errors.New("iv length error")
	}
	xTime, err := strconv.ParseInt(r.Header.Get("X-Client-TimeStamp"), 10, 64)
	if err != nil {
		return nil, err
	}
	xPubKeyNumber, err := strconv.ParseInt(r.Header.Get("X-PubKey-Number"), 10, 64)
	if err != nil {
		return nil, err
	}
	xPubKeyNumberReal := uint32(xTime ^ xPubKeyNumber)
	servePriKey, ok := configs.PrivateKeyByPubKey(xPubKeyNumberReal)
	if !ok {
		return nil, errors.New("xPubKeyNumber not found")
	}
	clientPubKey := r.Header.Get("X-Client-PubKey")
	if len(clientPubKey) == 0 {
		return nil, errors.New("X-Client-PubKey not found")
	}
	clientPubKeyBytes, err := base64.StdEncoding.DecodeString(clientPubKey)
	if err != nil {
		return nil, err
	}
	clientSignPubKey := r.Header.Get("X-Client-SignPubKey")
	if len(clientPubKey) == 0 {
		return nil, errors.New("X-Client-PubKey not found")
	}
	signPubKeyBytes, err := base64.StdEncoding.DecodeString(clientSignPubKey)
	if err != nil {
		return nil, err
	}
	curve := ecdh.X25519()
	clientPublicKey, err := curve.NewPublicKey(clientPubKeyBytes)
	if err != nil {
		return nil, err
	}
	sharedKey, err := servePriKey.ECDH(clientPublicKey)
	if err != nil {
		return nil, err
	}
	xUser := r.Header.Get("X-User-U")
	var userId string
	if xUser == "" {
		return nil, errors.New("X-User not found")
	}
	userIdBytes, err := base64.StdEncoding.DecodeString(xUser)
	if err != nil {
		return nil, err
	}
	userId = string(myaes.EncOrDecWithKeyAndIV(userIdBytes, sharedKey, iv))
	bodyEncBytes, err := base64.StdEncoding.DecodeString(string(encBody))
	if err != nil {
		return nil, err
	}
	bodyDec := myaes.EncOrDecWithKeyAndIV(bodyEncBytes, sharedKey, iv)
	signature, err := base64.StdEncoding.DecodeString(r.Header.Get("X-Sign"))
	if err != nil {
		return nil, err
	}
	verified := ed25519.Verify(signPubKeyBytes, bodyDec, signature)
	if !verified {
		return nil, errors.New("signature error")
	}
	xUserAgent := r.Header.Get("X-User-Agent")
	var UA *types.UA
	if xUserAgent != "" {
		UA, ok = getUAFromUserAgent(xUserAgent)
		if !ok {
			return nil, errors.New("user agent error")
		}
	} else {
		userAgent := r.Header.Get("User-Agent")
		if userAgent == "" {
			return nil, errors.New("empty user-agent")
		}
		UA = &types.UA{
			UserAgent: userAgent,
		}
	}
	err = r.ParseForm()
	if err != nil {
		return nil, err
	}
	deviceId := r.Header.Get("X-User-D")
	if len(deviceId) == 0 {
		return nil, errors.New("X-Device-Id error")
	}
	reqParam := ReqParam{
		QueryForm:       r.Form,
		Header:          r.Header,
		UA:              UA,
		Body:            bodyDec,
		XIV:             iv,
		ServePriKey:     servePriKey,
		xClientTime:     xTime,
		DeviceId:        deviceId,
		ClientPubKey:    clientPublicKey,
		SignPubKeyBytes: signPubKeyBytes,
		Sign:            signature,
		UserId:          userId,
	}
	return &reqParam, nil
}
func getParamEnc(r *http.Request) (*ReqParam, error) {

	iv, err := base64.StdEncoding.DecodeString(r.Header.Get("X-IV"))
	if err != nil {
		return nil, err
	}
	if len(iv) != 16 {
		return nil, errors.New("iv length error")
	}
	xTime, err := strconv.ParseInt(r.Header.Get("X-Client-TimeStamp"), 10, 64)
	if err != nil {
		return nil, err
	}
	xPubKeyNumber, err := strconv.ParseInt(r.Header.Get("X-PubKey-Number"), 10, 64)
	if err != nil {
		return nil, err
	}
	xPubKeyNumberReal := uint32(xTime ^ xPubKeyNumber)
	servePriKey, ok := configs.PrivateKeyByPubKey(xPubKeyNumberReal)
	if !ok {
		return nil, errors.New("xPubKeyNumber not found")
	}
	clientPubKey := r.Header.Get("X-Client-PubKey")
	if len(clientPubKey) == 0 {
		return nil, errors.New("X-Client-PubKey not found")
	}
	clientPubKeyBytes, err := base64.StdEncoding.DecodeString(clientPubKey)
	if err != nil {
		return nil, err
	}
	clientSignPubKey := r.Header.Get("X-Client-SignPubKey")
	if len(clientPubKey) == 0 {
		return nil, errors.New("X-Client-PubKey not found")
	}
	signPubKeyBytes, err := base64.StdEncoding.DecodeString(clientSignPubKey)
	if err != nil {
		return nil, err
	}
	curve := ecdh.X25519()
	clientPublicKey, err := curve.NewPublicKey(clientPubKeyBytes)
	if err != nil {
		return nil, err
	}
	sharedKey, err := servePriKey.ECDH(clientPublicKey)
	if err != nil {
		return nil, err
	}
	xUser := r.Header.Get("X-User-U")
	var userId string
	if xUser == "" {
		return nil, errors.New("X-User not found")
	}
	userIdBytes, err := base64.StdEncoding.DecodeString(xUser)
	if err != nil {
		return nil, err
	}
	userId = string(myaes.EncOrDecWithKeyAndIV(userIdBytes, sharedKey, iv))
	xUserAgent := r.Header.Get("X-User-Agent")
	var UA *types.UA
	if xUserAgent != "" {
		UA, ok = getUAFromUserAgent(xUserAgent)
		if !ok {
			return nil, errors.New("user agent error")
		}
	} else {
		userAgent := r.Header.Get("User-Agent")
		if userAgent == "" {
			return nil, errors.New("empty user-agent")
		}
		UA = &types.UA{
			UserAgent: userAgent,
		}
	}
	err = r.ParseForm()
	if err != nil {
		return nil, err
	}
	deviceId := r.Header.Get("X-User-D")
	if len(deviceId) == 0 {
		return nil, errors.New("X-Device-Id error")
	}
	reqParam := ReqParam{
		QueryForm:       r.Form,
		Header:          r.Header,
		UA:              UA,
		Body:            nil,
		XIV:             iv,
		ServePriKey:     servePriKey,
		xClientTime:     xTime,
		DeviceId:        deviceId,
		ClientPubKey:    clientPublicKey,
		SignPubKeyBytes: signPubKeyBytes,
		Sign:            nil,
		UserId:          userId,
	}
	return &reqParam, nil
}

func getParamWithDefaultParam(r *http.Request, defaultParamMap map[string]string, keys ...string) (*ReqParam, error) {
	return checkAndGetParam(r, defaultParamMap, keys...)

}
func checkAndGetParam(r *http.Request, defaultParamMap map[string]string, keys ...string) (*ReqParam, error) {
	_ = r.ParseMultipartForm(4096)
	var queryMap = make(map[string]string)
	for k, v := range defaultParamMap {
		queryMap[k] = v
		if value := r.FormValue(k); value != "" {
			queryMap[k] = value
		}
	}
	for _, k := range keys {
		var v string
		v = r.FormValue(k)
		if v == "" {
			return nil, fmt.Errorf("%s 不能为空", k)
		}
		queryMap[k] = v
	}
	reqParam := &ReqParam{
		QueryForm:    r.Form,
		Header:       nil,
		Body:         nil,
		XIV:          nil,
		ServePriKey:  nil,
		ClientPubKey: nil,
		xClientTime:  0,
		Sign:         nil,
	}
	return reqParam, nil
}
