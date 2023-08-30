package myerr

// 为了可能的数据层微服务拆分做准备,

// Err 可以直接对端上提示的Err。
type Err string

func IsErr(err error) bool {
	_, ok := err.(Err)
	return ok
}

func (e Err) Error() string {
	return string(e)
}

const (
	DataNotFound      = Err("data not found")
	ErrUserNotLogin   = Err("用户未登录")
	ErrUserLoginAgain = Err("请重新登录")
)
