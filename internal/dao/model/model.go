package model

import (
	"myoption/internal/dao/mtype"
	"myoption/types/fd"
	"time"
)

// UserInfo 表示用户信息
type UserInfo struct {
	ID            int64            `json:"id"`              // 用户信息表主键ID
	UserId        string           `json:"user_id"`         // 用户ID
	Nick          string           `json:"nick"`            // 用户昵称
	X25519PubKey  string           `json:"x25519_pub_key"`  //
	Ed25519PubKey string           `json:"ed25519_pub_key"` // Ed25519公钥
	Status        mtype.UserStatus `json:"status"`          // 用户状态，1表示激活，其他值表示禁用或其他状态
	CreateTime    time.Time        `json:"create_time"`     // 用户信息创建时间
	UpdateTime    time.Time        `json:"update_time"`     // 用户信息最后更新时间，默认为当前时间
}

// UserKey 表示用户密钥
type UserKey struct {
	ID               int64     `json:"id"`       // 用户信息表主键ID
	UserId           string    `json:"user_id"`  // 用户ID
	Password         string    `json:"password"` // 用户密码 PBKDF2-SHA256
	Salt             string    `json:"salt"`
	X25519PriEncKey  string    `json:"x25519_pri_enc_key"`  // X25519私钥加密密钥
	Ed25519PriEncKey string    `json:"ed25519_pri_enc_key"` // Ed25519私钥加密密钥
	CreateTime       time.Time `json:"create_time"`         // 用户信息创建时间
	UpdateTime       time.Time `json:"update_time"`         // 用户信息最后更新时间，默认为当前时间
}

func (u *UserInfo) ToFD() *fd.UserInfo {
	if u == nil {
		return nil
	}
	return &fd.UserInfo{
		ID:               u.ID,
		UserID:           u.UserId,
		Nick:             u.Nick,
		X25519PublicKey:  u.X25519PubKey,
		Ed25519PublicKey: u.Ed25519PubKey,
		CreateTime:       u.CreateTime,
		UpdateTime:       u.UpdateTime,
	}

}
