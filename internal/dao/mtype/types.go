package mtype

type UserStatus int

const (
	UserStatusNormal = UserStatus(1)
	UserStatusFreeze = UserStatus(2)
)

func (u UserStatus) Check() bool {
	switch u {
	case UserStatusNormal, UserStatusFreeze:
		return true
	default:
		return false
	}
}
