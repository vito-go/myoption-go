package cache

import (
	"github.com/go-redis/redis/v8"

	"myoption/internal/dao"
)

const null = "null"

type Cache struct {
	redisCli redis.Cmdable
	allDao   *dao.AllDao
	User     *user

	Wallet *wallet
}

func (c *Cache) RedisCli() redis.Cmdable {
	return c.redisCli
}

func New(redisCli redis.Cmdable, allDao *dao.AllDao) *Cache {
	return &Cache{
		redisCli: redisCli,
		allDao:   allDao,
		User:     &user{redisCli: redisCli},
		Wallet:   &wallet{redisCli: redisCli},
	}
}
