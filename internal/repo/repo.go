package repo

import (
	"context"
	"myoption/configs"
	model "myoption/internal/dao/model"
	"myoption/internal/dao/mtype"
	"sync"
	"time"

	"myoption/iface"
	"myoption/iface/myerr"
	"myoption/internal/connector"
	"myoption/internal/dao"
	"myoption/pkg/util/slice"

	"github.com/go-redis/redis/v8"
	"github.com/vito-go/mylog"

	"myoption/internal/cache"
)

type Client struct {
	_       sync.Mutex
	allDao  *dao.AllDao
	cache   *cache.Cache
	UserCli iface.UserAPI
	//
	Wallet             iface.WalletIface
	StockData          iface.StockDataIface
	OrdersBinaryOption iface.OrdersBinaryOptionIface
}

func (c *Client) AllDao() *dao.AllDao {
	return c.allDao
}

func (c *Client) Cache() *cache.Cache {
	return c.cache
}

type userCli struct {
	cache *cache.Cache
	dao   *dao.AllDao
}

func NewClient(conn *connector.Connector) *Client {
	redisCli := conn.RedisCli
	allDao := dao.NewAllDao(conn.GDB)
	c := cache.New(redisCli, allDao)
	return &Client{
		allDao:  allDao,
		cache:   c,
		UserCli: &userCli{cache: c, dao: allDao},

		Wallet:             &Wallet{cache: c, allDao: allDao},
		StockData:          NewStockOrigin(allDao, c, configs.Symbols),
		OrdersBinaryOption: &OrdersBinaryOption{cache: c, allDao: allDao},
	}
}

var _ iface.UserAPI = (*userCli)(nil)

func (c *userCli) CreateUser(ctx context.Context, createUserInfo *model.UserInfo, userKey *model.UserKey) (info *model.UserInfo, err error) {
	userId := createUserInfo.UserId
	tx := c.dao.GDB().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		if err = tx.Commit().Error; err != nil {
			return
		}
		c.cache.User.DelUserInfo(ctx, userId)
		c.cache.Wallet.DelWalletInfo(ctx, userId)
	}()

	mo := model.UserInfo{
		ID:            0,
		UserId:        createUserInfo.UserId,
		Nick:          createUserInfo.Nick,
		Status:        mtype.UserStatusNormal,
		X25519PubKey:  createUserInfo.X25519PubKey,
		Ed25519PubKey: createUserInfo.Ed25519PubKey,
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
	}
	if err = c.dao.UserInfo.CreateTX(tx, &mo); err != nil {
		return nil, err
	}

	userKeyMo := model.UserKey{
		ID:               0,
		UserId:           userKey.UserId,
		Password:         userKey.Password,
		Salt:             userKey.Salt,
		X25519PriEncKey:  userKey.X25519PriEncKey,
		Ed25519PriEncKey: userKey.Ed25519PriEncKey,
		CreateTime:       time.Now(),
		UpdateTime:       time.Now(),
	}
	if err = c.dao.UserKey.CreateTX(tx, &userKeyMo); err != nil {
		return nil, err
	}
	walletMo := model.WalletInfo{
		ID:     0,
		UserId: userId,
		//Balance:      0,
		Balance:      10000,
		FrozenAmount: 0,
		//TotalAmount:  0,
		TotalAmount: 10000,
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	}
	err = c.dao.Wallet.CreateTX(tx, &walletMo)
	if err != nil {
		return nil, err
	}
	return &mo, nil
}

func (c *userCli) GetUserInfoByUserId(ctx context.Context, userId string) (*model.UserInfo, error) {
	return c.getUserInfoByUserId(ctx, userId)
}

func (c *userCli) GetUserKeyByUserId(ctx context.Context, userId string) (*model.UserKey, error) {
	return c.getUserKeyByUserId(ctx, userId)
}

func (c *userCli) GetUserInfoMapByUserIds(ctx context.Context, userIds ...string) (map[string]model.UserInfo, error) {
	userIds = slice.FilterStr(userIds)
	result := make(map[string]model.UserInfo, len(userIds))
	for _, userId := range userIds {
		info, err := c.getUserInfoByUserId(ctx, userId)
		if err != nil {
			// 为了保持一致性，除非用户不存在，否则有一个出错就返回
			if err == myerr.DataNotFound {
				continue
			}
			mylog.Ctx(ctx).Error(err)
			return nil, err
		}
		result[userId] = *info
	}
	return result, nil
}
func (c *userCli) GetUserInfosByUserIds(ctx context.Context, userIds ...string) ([]*model.UserInfo, error) {
	userIds = slice.FilterStr(userIds)
	result := make([]*model.UserInfo, 0, len(userIds))
	for _, userId := range userIds {
		info, err := c.getUserInfoByUserId(ctx, userId)
		if err != nil {
			// 为了保持一致性，除非用户不存在，否则有一个出错就返回
			if err == myerr.DataNotFound {
				continue
			}
			mylog.Ctx(ctx).Error(err)
			return nil, err
		}
		result = append(result, info)
	}
	return result, nil
}

// getUserInfoByUserId 最核心的底层方法
func (c *userCli) getUserInfoByUserId(ctx context.Context, userId string) (*model.UserInfo, error) {
	result, err := c.cache.User.GetUserInfo(ctx, userId)
	if err == nil {
		return result, nil
	}
	if err != redis.Nil {
		return nil, err
	}
	item, err := c.dao.UserInfo.ItemByUserId(ctx, userId)
	if err == myerr.DataNotFound {
		//c.cache.SetNilKey(ctx, k)
		return nil, myerr.DataNotFound
	}
	if err != nil {
		return nil, err
	}
	mylog.Ctx(ctx).WithField("userId", userId).Info("获取用户信息，并添加缓存")
	c.cache.User.SetUserInfo(ctx, userId, item)
	return item, nil
}

// getUserInfoByUserId 最核心的底层方法
func (c *userCli) getUserKeyByUserId(ctx context.Context, userId string) (*model.UserKey, error) {
	result, err := c.cache.User.GetUserKey(ctx, userId)
	if err == nil {
		return result, nil
	}
	if err != redis.Nil {
		return nil, err
	}
	item, err := c.dao.UserKey.ItemByUserId(ctx, userId)
	if err == myerr.DataNotFound {
		//c.cache.SetNilKey(ctx, k)
		return nil, myerr.DataNotFound
	}
	if err != nil {
		return nil, err
	}
	mylog.Ctx(ctx).WithField("userId", userId).Info("并添加缓存")
	c.cache.User.SetUserKey(ctx, userId, item)
	return item, nil
}
