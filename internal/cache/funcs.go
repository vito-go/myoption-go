package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/vito-go/mylog"

	"myoption/iface/myerr"
)

func getList(ctx context.Context, redisCli redis.Cmdable, k string) ([]string, error) {
	items, err := redisCli.LRange(ctx, k, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, redis.Nil
	}
	if len(items) == 1 && items[0] == null {
		return nil, nil
	}
	result := make([]string, 0, len(items))
	for _, item := range items {
		if item == null {
			continue
		}
		result = append(result, item)
	}
	return result, nil
}

func (c *Cache) SetNilKey(ctx context.Context, k string) {
	var expireTime = time.Minute
	err := c.redisCli.SetEX(ctx, k, null, expireTime).Err()
	if err != nil {
		mylog.Ctx(ctx).WithField("key", k).Error("缓存 nilKey设置失败. error: ", err)
		return
	}
	mylog.Ctx(ctx).WithField("key", k).Info("缓存: nilKey设置 expireTime: ", expireTime)
}

func (c *Cache) SetHashNilKey(ctx context.Context, k string) {
	var expireTime = time.Second * (time.Duration(time.Now().UnixNano()%10) + 10)
	err := c.redisCli.HSet(ctx, k, null, null).Err()
	if err != nil {
		mylog.Ctx(ctx).WithField("key", k).Error("缓存 nilKey设置失败. error: ", err)
		return
	}
	c.redisCli.Expire(ctx, k, expireTime)
}
func (c *Cache) SetListNilKey(ctx context.Context, k string) {
	var expireTime = time.Second * (time.Duration(time.Now().UnixNano()%10) + 10)
	err := c.redisCli.LPush(ctx, k, null).Err()
	if err != nil {
		mylog.Ctx(ctx).WithField("key", k).Error("缓存 nilKey设置失败. error: ", err)
		return
	}
	c.redisCli.Expire(ctx, k, expireTime)
}

// set . valueData 应该是 Marshal 后的值.或者直接是一个字符串/[]byte， 要和get对应
func set(ctx context.Context, redisCli redis.Cmdable, k string, valueData []byte, keyExpire time.Duration) {
	// 首次设定最多等待缓存建立800毫秒
	const wait = time.Millisecond * 300
	ch := make(chan struct{}, 1)
	go func() {
		defer func() {
			ch <- struct{}{}
		}()
		if err := redisCli.SetEX(ctx, k, valueData, keyExpire).Err(); err != nil {
			mylog.Ctx(ctx).Error(err.Error())
			if err = redisCli.Del(ctx, k).Err(); err != nil {
				mylog.Ctx(ctx).Error(err.Error())
				return
			}
			return
		}
		mylog.Ctx(ctx).WithFields("key", k, "expire", keyExpire.String()).Info("缓存: 设置.")
	}()
	timer := time.NewTimer(wait)
	defer timer.Stop()
	select {
	case <-ch:
	case <-timer.C:
	}
}
func updateHashValue(ctx context.Context, redisCli redis.Cmdable, k string, field string, value string) {
	values := []interface{}{null, 1, field, value}
	num, err := redisCli.HSet(ctx, k, values...).Result()
	if err != nil {
		mylog.Ctx(ctx).WithFields("key", k, "field", field, "values", values).Error(err.Error())
		// 为了数据完整性， 添加群成员如果缓存失败的话，应该将整个key删除
		if err = redisCli.Del(ctx, k).Err(); err != nil {
			mylog.Ctx(ctx).Error(err.Error())
			return
		}
	}
	mylog.Ctx(ctx).WithFields("key", k, "field", field).Infof("resultNumber: %d", num)
	// 沒有重建緩存
	return
	switch num {
	case 0: // 两个field都存在. 如：更新成员
	case 1: // 正常情况，hashNullField 存在即key存在, 如: 新加入成员
	case 2: // hashNullField 值添加成功 key不存在，单单只添加成功一个成员. 所以要删除key
		if err = redisCli.Del(ctx, k).Err(); err != nil {
			mylog.Ctx(ctx).Error(err.Error())
		}
	}
}

// getWithUnmarshal valuePtr应该为一个指针.
func getWithUnmarshal(ctx context.Context, redisCli redis.Cmdable, k string, valuePtr interface{}) error {
	result, err := redisCli.Get(ctx, k).Result()
	if err != nil {
		return err
	}
	if result == null {
		mylog.Ctx(ctx).WithField("key", k).Infof("缓存: 获取. null value")
		return myerr.DataNotFound
	}
	mylog.Ctx(ctx).WithField("key", k).Infof("缓存: 获取")
	if len(result) == 0 {
		return nil
	}
	return json.Unmarshal([]byte(result), valuePtr)
}

// getValueData .
func getValueData(ctx context.Context, redisCli redis.Cmdable, k string) (string, error) {
	result, err := redisCli.Get(ctx, k).Result()
	if err != nil {
		return "", err
	}
	if result == null {
		return "", myerr.DataNotFound
	}
	mylog.Ctx(ctx).WithFields("key", k, "result", result).Infof("缓存: 获取")
	if len(result) == 0 {
		return "", nil
	}
	return result, nil
}

// del 删除缓存
func del(ctx context.Context, redisCli redis.Cmdable, k string) {
	i, err := redisCli.Del(ctx, k).Result()
	if err != nil {
		mylog.Ctx(ctx).WithField("key", k).Error("缓存删除失败. error: ", err)
		return
	}
	mylog.Ctx(ctx).WithField("key", k).Info("缓存: 删除. result: ", i)
}

// hGetAll hGetAll 不会返回 redis.Nil这种错误, 返回一个空map
func hGetAll(ctx context.Context, redisCli redis.Cmdable, k string) (map[string]string, error) {
	result, err := redisCli.HGetAll(ctx, k).Result()
	if err != nil {
		return nil, err
	}
	mylog.Ctx(ctx).WithFields("key", k, "result", result).Info("缓存: 获取: ", len(result))
	return result, nil
}

// hGet  k不存在或者 filed对应的值不存在就返回 redis.Nil
func hGet(ctx context.Context, redisCli redis.Cmdable, k, filed string) (string, error) {
	result, err := redisCli.HGet(ctx, k, filed).Result()
	if err != nil {
		return "", err
	}
	mylog.Ctx(ctx).WithFields("key", k, "result", result).Info("缓存: 获取: ")
	return result, nil
}

// hSet .
func hSet(ctx context.Context, redisCli redis.Cmdable, k string, keyExpire time.Duration, values ...interface{}) {
	result, err := redisCli.HSet(ctx, k, values...).Result()
	if err != nil {
		mylog.Ctx(ctx).Error(err.Error())
		return
	}
	if err = redisCli.Expire(ctx, k, keyExpire).Err(); err != nil {
		mylog.Ctx(ctx).Error(err.Error())
		return
	}
	mylog.Ctx(ctx).WithFields("key", k, "values", values).Infof("缓存: 设置: %+v", result)
}

// zrange  zset 不会返回 redis.Nil这种错误 返回一个空列表。
func zrange(ctx context.Context, redisCli redis.Cmdable, k string, start, stop int64) ([]string, error) {
	result, err := redisCli.ZRange(ctx, k, start, stop).Result()
	if err != nil {
		return nil, err
	}
	mylog.Ctx(ctx).WithFields("key", k, "start", start,
		"stop", stop, "result", result).Info("缓存: Zrange获取: ", len(result))
	return result, nil
}

func hdel(ctx context.Context, redisCli redis.Cmdable, k string, fields ...string) {
	result, err := redisCli.HDel(ctx, k, fields...).Result()
	if err != nil {
		redisCli.Del(ctx, k)
		mylog.Ctx(ctx).Error(err.Error())
	}
	mylog.Ctx(ctx).WithFields("key", k, "fields", fields).Infof("hash缓存: 删除. result: %+v err: %+v", result, err)
}
