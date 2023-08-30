package httperr

type Err string

func (e Err) Error() string {
	return string(e)
}

const (
	ErrInternal       = Err("system error")
	ErrParam          = Err("param not valid")
	ErrAuth           = Err("auth failed")
	ErrData           = Err("data error")
	ErrUA             = Err("未识别的客户端")
	ErrIllegalCall    = Err("非法调用")
	ErrSecTimeOut     = Err("安全超时")
	ErrDuplicate      = Err("重复操作")
	ErrUserNotLogin   = Err("用户未登录")
	ErrUserNotFound   = Err("用户不存在")
	ErrGroupIsFull    = Err("the group is full")
	ErrInBlack        = Err("对方已把你拉黑")
	ErrMomentNotExist = Err("moment不存在")
)
