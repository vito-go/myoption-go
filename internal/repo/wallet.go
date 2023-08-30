package repo

import (
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/vito-go/mylog"
	"myoption/iface/myerr"
	"myoption/internal/cache"
	"myoption/internal/cache/cachekey"
	"myoption/internal/dao"
	"myoption/internal/dao/model"
	"myoption/internal/dao/mtype"
	"strconv"
	"strings"
	"time"
)

type Wallet struct {
	cache  *cache.Cache
	allDao *dao.AllDao
}

//Recharge（充值）
//Consumption（消费）
//SendRedPacket（发红包）
//PurchaseExpressionPack（购买表情包）
//Balance（余额）
//Amount（金额）
//Transaction（交易）
//PaymentMethod（付款方式）
//Currency（货币）
//Refund（退款）
//Deduction（扣除）
//TopUp（充值）
//Cost（花费）
//Coupon（优惠券）
//ConversionRate（转换率）

func (w *Wallet) RechargeNotification(ctx context.Context, param *model.WalletDetail) (transactionId string, err error) {
	userId := param.UserId
	amount := param.Amount
	TX := w.allDao.Begin(ctx, nil)
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			if e := TX.Rollback().Error; e != nil {
				mylog.Ctx(ctx).Error(e.Error())
			}
			return
		}
		if err = TX.Commit().Error; err != nil {
			return
		}
	}()
	orderId := NewOrderID(userId)
	detail := model.WalletDetail{
		TransId:       "",
		TransType:     0,
		UserId:        param.UserId,
		Amount:        amount,
		Status:        mtype.TransStatusWaiting,
		Remark:        param.Remark,
		SourceKind:    0,
		SourceTransId: param.SourceTransId,
		FromAccount:   "",
		ToAccount:     "",
		Balance:       0,
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
	}
	err = w.allDao.WalletDetail.CreateTX(TX, &detail)
	if err != nil {
		return "", err
	}

	return orderId, nil
}
func (w *Wallet) AddWithSuccess(ctx context.Context, transactionId, userId string, amount int64, unFrozenUserId string) (err error) {
	TX := w.allDao.Begin(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if e := TX.Rollback().Error; e != nil {
				mylog.Ctx(ctx).Error(e.Error())
			}
			return
		}
		if err = TX.Commit().Error; err != nil {
			return
		}
		w.cache.Wallet.DelWalletInfo(ctx, userId)
	}()
	var rowAffect int64
	rowAffect, err = w.allDao.WalletDetail.UpdateStatusSucByTX(TX, transactionId, mtype.TransStatusSuc, amount)
	if err != nil {
		return err
	}
	if rowAffect == 0 {
		return fmt.Errorf("wallet detail: no recoSendRedPacketrd")
	}
	rowAffect, err = w.allDao.Wallet.AddAmount(TX, userId, amount)
	if err != nil {
		return err
	}
	if rowAffect == 0 {
		return fmt.Errorf("wallet Recharge: no user wallet")
	}
	if unFrozenUserId != "" {
		rowAffect, err = w.allDao.Wallet.UnFrozenAmount(TX, unFrozenUserId, amount)
		if err != nil {
			return err
		}
		if rowAffect == 0 {
			return fmt.Errorf("wallet Recharge: no user wallet")
		}
	}

	return nil
}

// ConsumerWallet 先校验余额度
func (w *Wallet) ConsumerWallet(ctx context.Context, param *model.WalletDetail, sub bool) (transactionId string, err error) {
	//5/11/2023, 1:46:47 AM
	//
	//在给对方用户发送红包的业务中，一般是先冻结账户的金额，待对方领取红包成功后再实际扣减冻结的金额。
	//
	//这种方式可以避免因为发送红包后对方未领取的情况下导致资金流动不稳定。如果你直接扣掉账户金额，会使用户账户变少，但是如果对方没接收到红包而需要退款，则操作比较复杂。将资金冻结并在对方实际领取时再进行扣减，可以更好地控制资金流转，并且能更灵活的处理退款等问题。
	detail := &model.WalletDetail{}
	*detail = *param
	TX := w.allDao.Begin(ctx, nil)
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			if e := TX.Rollback().Error; e != nil {
				mylog.Ctx(ctx).Error(e.Error())
			}
			return
		}
		if err = TX.Commit().Error; err != nil {
			return
		}
		w.cache.Wallet.DelWalletInfo(ctx, param.UserId)

	}()
	orderId := NewOrderID(detail.UserId)
	detail.TransId = orderId
	var rowAffect int64
	if sub {
		rowAffect, err = w.allDao.Wallet.SubAmount(TX, detail.UserId, detail.Amount)
		if err != nil {
			return "", err
		}
		if rowAffect == 0 {
			return "", fmt.Errorf("wallet : 余额不足？")
		}
	} else {
		rowAffect, err = w.allDao.Wallet.FreezeAmount(TX, detail.UserId, detail.Amount)
		if err != nil {
			return "", err
		}
		if rowAffect == 0 {
			return "", fmt.Errorf("wallet : 余额不足？")
		}
	}
	err = w.allDao.WalletDetail.CreateTX(TX, detail)
	if err != nil {
		return "", err
	}
	return orderId, nil
}

// NewOrderID 生成一个20位的订单ID
// 前4位:年份
// 中间8位: 用户标识的sh1摘要前7位
// 后8位: 随机时间戳的sha1摘要前7位
func NewOrderID(uid string) string {
	today, _ := strconv.ParseInt(time.Now().Format("20060102"), 10, 64)
	yearX := strings.ToUpper(strconv.FormatInt(today, 32))
	uidX := fmt.Sprintf("%X", sha1.Sum([]byte(uid)))[:8]
	nanoX := fmt.Sprintf("%X", sha1.Sum([]byte(strconv.FormatInt(time.Now().UnixNano(), 10))))[:8]
	result := fmt.Sprintf("%s%s%s", yearX, uidX, nanoX)
	return result
}

// GetWalletInfoByUserId 最核心的底层方法
func (w *Wallet) GetWalletInfoByUserId(ctx context.Context, userId string) (*model.WalletInfo, error) {

	return w.getWalletInfoByUserId(ctx, userId)
}
func (w *Wallet) getWalletInfoByUserId(ctx context.Context, userId string) (*model.WalletInfo, error) {
	k := fmt.Sprintf(cachekey.WalletInfo, userId)
	result, err := w.cache.Wallet.GetWalletInfo(ctx, userId)
	if err == nil {
		return result, nil
	}
	if err != redis.Nil {
		return nil, err
	}
	item, err := w.allDao.Wallet.ItemByUserId(ctx, userId)
	if err == myerr.DataNotFound {
		w.cache.SetNilKey(ctx, k)
		return nil, myerr.DataNotFound
	}
	if err != nil {
		return nil, err
	}
	mylog.Ctx(ctx).WithField("userId", userId).Info("并添加缓存")
	w.cache.Wallet.SetWalletInfo(ctx, userId, item)
	return item, nil
}
