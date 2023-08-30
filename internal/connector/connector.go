package connector

import (
	"context"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/apache/pulsar-client-go/pulsar/log"
	"github.com/go-redis/redis/v8"
	"github.com/vito-go/mylog"
	"gorm.io/gorm"
	"time"

	"myoption/conf"
)

type Connector struct {
	RedisCli *redis.Client
	GDB      *gorm.DB

	PulsarClient pulsar.Client
}

func New(cfg *conf.Cfg) (*Connector, error) {

	gdb, err := OpenGromDB(cfg.Database)
	if err != nil {
		return nil, err
	}
	redisCli, err := NewRedisClient(cfg.Redis)
	if err != nil {
		return nil, err
	}
	var pulsarCli pulsar.Client
	if cfg.Pulsar.ServiceURL != "" {
		pulsarCli, err = pulsar.NewClient(pulsar.ClientOptions{
			URL:               cfg.Pulsar.ServiceURL,
			OperationTimeout:  30 * time.Second,
			ConnectionTimeout: 30 * time.Second,
			Logger:            log.DefaultNopLogger(),
		})
		if err != nil {
			return nil, err
		}
	} else {
		mylog.Ctx(context.Background()).Warn("pulsar配置为空，忽略pulsar client初始化")
	}
	return &Connector{
		GDB:          gdb,
		RedisCli:     redisCli,
		PulsarClient: pulsarCli}, nil
}
func (c *Connector) Close(ctx context.Context) {
	if c.PulsarClient != nil {
		c.PulsarClient.Close()
	}
	db, err := c.GDB.DB()
	if err != nil {
		mylog.Ctx(ctx).Error(err.Error())
		return
	}
	err = db.Close()
	if err != nil {
		mylog.Ctx(ctx).Error(err.Error())
	}
}

// NewRedisClient generate a Redis client representing a pool of zero or more
// underlying connections. It's safe for concurrent use by multiple goroutines.
func NewRedisClient(cfg conf.RedisConf) (*redis.Client, error) {
	redisCli := redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     cfg.Addr,
		Username: cfg.UserName,
		Password: cfg.Password,
		DB:       cfg.DB,
		// 可以在配置中添加更多需要的配置
	})
	if err := redisCli.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return redisCli, nil
}
