package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/vito-go/mylog"
	"myoption/internal/cache/cachekey"
	"myoption/internal/dao/model"
	"time"
)

type wallet struct {
	redisCli redis.Cmdable
}

func (c *wallet) GetWalletInfo(ctx context.Context, userId string) (*model.WalletInfo, error) {
	var result model.WalletInfo
	k := fmt.Sprintf(cachekey.WalletInfo, userId)
	if err := getWithUnmarshal(ctx, c.redisCli, k, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *wallet) SetWalletInfo(ctx context.Context, userId string, item *model.WalletInfo) {
	if item == nil {
		mylog.Ctx(ctx).Error("忽略缓存设置: nil User")
		return
	}
	b, _ := json.Marshal(item)
	k := fmt.Sprintf(cachekey.WalletInfo, userId)
	var keyExpire = time.Hour * 12
	set(ctx, c.redisCli, k, b, keyExpire)
}

func (c *wallet) DelWalletInfo(ctx context.Context, userId string) {
	k := fmt.Sprintf(cachekey.WalletInfo, userId)
	del(ctx, c.redisCli, k)
}
