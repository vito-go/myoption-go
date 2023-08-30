package iface

import (
	"context"
	"myoption/internal/dao/model"
)

type UpdateUserInfoParam struct {
	Avatar     string `json:"avatar,omitempty"`
	Background string `json:"background,omitempty"`
	Nick       string `json:"nick,omitempty"`
	Region     string `json:"region,omitempty"`
	Info       string `json:"info,omitempty"`
}

type UserAPI interface {
	// 最新可用---------------------------

	// Register 注册.
	//Register(ctx context.Context, userId, nick, pwd, privateKeyByAes, publicKey string) (*model.UserInfo, error)

	CreateUser(ctx context.Context, createUserInfo *model.UserInfo, userKey *model.UserKey) (info *model.UserInfo, err error)
	// GetUserInfoByUserId .
	GetUserInfoByUserId(ctx context.Context, userId string) (*model.UserInfo, error)
	GetUserKeyByUserId(ctx context.Context, userId string) (*model.UserKey, error)
	// GetUserInfoMapByUserIds .// 为了保持一致性，除非用户不存在，否则有一个出错就返回
	GetUserInfoMapByUserIds(ctx context.Context, userIds ...string) (map[string]model.UserInfo, error)
	// GetUserInfosByUserIds .
	GetUserInfosByUserIds(ctx context.Context, userIds ...string) ([]*model.UserInfo, error)
}

type UserLastOnlineInfo struct {
	Time      int64  `json:"time,omitempty"`
	IpAddress string `json:"ipAddress,omitempty"`
}
