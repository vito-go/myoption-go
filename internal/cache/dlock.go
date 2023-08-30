package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/vito-go/mylog"
	"myoption/internal/cache/cachekey"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type DLock struct {
	redisCli redis.Cmdable
}

// DistributeDoOnce 分布式执行一次。默认要求10秒内完成.
func (c *Cache) DistributeDoOnce(flag string) (bool, error) {
	return c.redisCli.SetNX(context.Background(), fmt.Sprintf(cachekey.KeyGeneralDistributeDoOnce, flag),
		1, time.Second*15).Result()
}

// DistributeDoOnceWithTime 若干时间内分布式执行一次。
func (c *Cache) DistributeDoOnceWithTime(flag string, t time.Duration) (bool, error) {
	return c.redisCli.SetNX(context.Background(), fmt.Sprintf(cachekey.KeyGeneralDistributeDoOnce, flag),
		1, t).Result()
}

// DistributeDoOnceDel 若干时间内分布式执行一次。
func (c *Cache) DistributeDoOnceDel(ctx context.Context, flag string) {
	key := fmt.Sprintf(cachekey.KeyGeneralDistributeDoOnce, flag)
	_, err := c.redisCli.Del(context.Background(), key).Result()
	if err != nil {
		mylog.Ctx(ctx).Error(err.Error())
	}
}

type dlock struct {
	redisCli redis.Cmdable
	key      string
	nano     int64
}

func (c *Cache) NewDLock(flag string) *dlock {
	return &dlock{redisCli: c.redisCli, key: fmt.Sprintf(cachekey.KeyGeneralDistributeLock, flag)}
}

// Lock 阻塞式获取锁.
func (d *dlock) Lock(ctx context.Context) error {
	return d.lock(ctx, time.Second*15)
}

// UnLock 释放锁。
func (d *dlock) UnLock(ctx context.Context) {
	result, err := d.redisCli.Get(ctx, d.key).Result()
	if err != nil {
		// 不能return，继续尝试解锁删除.
		mylog.Ctx(ctx).Error(err.Error())
	}
	if result != strconv.FormatInt(d.nano, 10) {
		mylog.Ctx(ctx).Error("unlock failed: 执行任务超时.")
	}
	_, err = d.redisCli.Del(ctx, d.key).Result()
	if err != nil {
		mylog.Ctx(ctx).Error("unlock failed: ", err.Error())
	}
}

func (d *dlock) lock(ctx context.Context, expire time.Duration) error {
	start := time.Now()
	nowNano := start.UnixNano()
	d.nano = nowNano
	ifSet, err := d.redisCli.SetNX(ctx, d.key, nowNano, expire).Result()
	if err != nil {
		return fmt.Errorf("DLock: lock error: %w", err)
	}
	if ifSet {
		return nil
	}
	const maxRetry = 42 // 5s左右
	sleep := time.Millisecond * 10
	for i := 0; true; i++ {
		if time.Since(start) > expire || i >= maxRetry {
			return errors.New("DLock: lock error: timeout")
		}
		if ifSet, err = d.redisCli.SetNX(ctx, d.key, 1, expire).Result(); err != nil {
			return fmt.Errorf("DLock: lock error: %w", err)
		}
		if ifSet {
			return nil
		}
		time.Sleep(sleep)
		if i <= 10 {
			sleep = sleep + time.Millisecond*10
		} else {
			sleep = sleep + time.Duration(i-10)*time.Millisecond*10
		}
	}
	return nil
}
