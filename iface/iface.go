package iface

import (
	"context"
	"myoption/internal/dao/model"
)

type UserAPI interface {
	// CreateUser .
	CreateUser(ctx context.Context, createUserInfo *model.UserInfo, userKey *model.UserKey) (info *model.UserInfo, err error)
	// GetUserInfoByUserId .
	GetUserInfoByUserId(ctx context.Context, userId string) (*model.UserInfo, error)
	GetUserKeyByUserId(ctx context.Context, userId string) (*model.UserKey, error)
	// GetUserInfoMapByUserIds .// in order to keep consistency, return an error if any error occurs
	GetUserInfoMapByUserIds(ctx context.Context, userIds ...string) (map[string]model.UserInfo, error)
	// GetUserInfosByUserIds .
	GetUserInfosByUserIds(ctx context.Context, userIds ...string) ([]*model.UserInfo, error)
}
