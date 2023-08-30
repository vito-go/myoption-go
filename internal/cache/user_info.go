package cache

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"myoption/internal/cache/cachekey"
	"myoption/internal/dao/model"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/vito-go/mylog"
)

type user struct {
	redisCli redis.Cmdable
}

func (c *user) GetUserInfo(ctx context.Context, userId string) (*model.UserInfo, error) {
	var result model.UserInfo
	k := fmt.Sprintf(cachekey.UserInfo, userId)
	if err := getWithUnmarshal(ctx, c.redisCli, k, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *user) SetUserInfo(ctx context.Context, userId string, item *model.UserInfo) {
	if item == nil {
		mylog.Ctx(ctx).Error("忽略缓存设置: nil User")
		return
	}
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(item)
	b := buf.Bytes()
	k := fmt.Sprintf(cachekey.UserInfo, userId)
	var keyExpire = time.Hour * 24
	set(ctx, c.redisCli, k, b, keyExpire)
}

func (c *user) DelUserInfo(ctx context.Context, userId string) {
	k := fmt.Sprintf(cachekey.UserInfo, userId)
	del(ctx, c.redisCli, k)
}

func (c *user) GetUserKey(ctx context.Context, userId string) (*model.UserKey, error) {
	var result model.UserKey
	k := fmt.Sprintf(cachekey.UserKey, userId)
	if err := getWithUnmarshal(ctx, c.redisCli, k, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *user) SetUserKey(ctx context.Context, userId string, item *model.UserKey) {
	if item == nil {
		mylog.Ctx(ctx).Error("忽略缓存设置: nil User")
		return
	}
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(item)
	b := buf.Bytes()
	k := fmt.Sprintf(cachekey.UserKey, userId)
	var keyExpire = time.Hour * 12
	set(ctx, c.redisCli, k, b, keyExpire)
}

func (c *user) DelUserKey(ctx context.Context, userId string) {
	k := fmt.Sprintf(cachekey.UserKey, userId)
	del(ctx, c.redisCli, k)
}
