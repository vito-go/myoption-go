package resp

import (
	"context"
	"fmt"
	"strconv"
)

const (
	generalErrCode   = 100000
	errCodeParam     = 100100
	errCodeToast     = 400400
	errGroupNotFound = 100404
	errNotInGroup    = 100403
)

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrParse = Error("系统内部数据解析错误")
)

type HTTPBody struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Tid     string      `json:"tid,omitempty"` // Tid traceID 全链路追踪tid
}

// DataOK data should, can be marshaled.
func DataOK(ctx context.Context, data interface{}) *HTTPBody {
	tid, _ := ctx.Value("tid").(int64)
	return &HTTPBody{
		Data: data,
		Tid:  strconv.FormatInt(tid, 10),
	}
}

// ErrToast .
func ErrToast(ctx context.Context, errMsg string) *HTTPBody {
	tid, _ := ctx.Value("tid").(int64)
	return &HTTPBody{
		Code:    errCodeToast,
		Message: errMsg,
		Tid:     strconv.FormatInt(tid, 10),
	}
}
func ErrCodeMsg(ctx context.Context, errCode int, errMsg string) *HTTPBody {
	tid, _ := ctx.Value("tid").(int64)
	return &HTTPBody{
		Code:    errCode,
		Message: errMsg,
		Tid:     strconv.FormatInt(tid, 10),
	}
}

func Err(ctx context.Context, message string) *HTTPBody {
	tid, _ := ctx.Value("tid").(int64)
	return &HTTPBody{
		Code:    generalErrCode,
		Message: message,
		Tid:     strconv.FormatInt(tid, 10),
	}
}

func ErrGroupNotFound(ctx context.Context) *HTTPBody {
	tid, _ := ctx.Value("tid").(int64)
	return &HTTPBody{
		Code:    errGroupNotFound,
		Message: "群不存在或已解散",
		Tid:     strconv.FormatInt(tid, 10),
	}
}
func ErrNotInGroup(ctx context.Context) *HTTPBody {
	tid, _ := ctx.Value("tid").(int64)
	return &HTTPBody{
		Code:    errNotInGroup,
		Message: "您已不在该群组中",
		Tid:     strconv.FormatInt(tid, 10),
	}
}

func ErrParam(ctx context.Context) *HTTPBody {
	tid, _ := ctx.Value("tid").(int64)
	return &HTTPBody{
		Code:    errCodeParam,
		Message: "请求参数不合法",
		Tid:     strconv.FormatInt(tid, 10),
	}
}

func Errf(ctx context.Context, format string, args ...interface{}) *HTTPBody {
	tid, _ := ctx.Value("tid").(int64)
	return &HTTPBody{
		Code:    generalErrCode,
		Message: fmt.Sprintf(format, args...),
		Tid:     strconv.FormatInt(tid, 10),
	}
}
