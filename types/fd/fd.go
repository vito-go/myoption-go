// Package fd 定义一些与前端交互的数据结构
package fd

import (
	"time"
)

// UserInfo 用户信息.
type UserInfo struct {
	ID               int64     `json:"id"`              // 用户信息表主键ID
	UserID           string    `json:"user_id"`         // 用户ID
	Nick             string    `json:"nick"`            // 用户昵称
	X25519PublicKey  string    `json:"x25519_pub_key"`  // X25519公钥
	Ed25519PublicKey string    `json:"ed25519_pub_key"` // Ed25519公钥
	CreateTime       time.Time `json:"create_time"`     // 用户信息创建时间
	UpdateTime       time.Time `json:"update_time"`     // 用户信息最后更新时间，默认为当前时间
}
