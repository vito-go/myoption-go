package repo

import (
	"context"
	"fmt"
	"github.com/vito-go/mylog"
	"myoption/iface"
	"myoption/internal/cache"
	"myoption/internal/dao"
	"myoption/internal/dao/model"
	"myoption/internal/dao/mtype"
	"strconv"
	"time"
)

type OrdersBinaryOption struct {
	cache  *cache.Cache
	allDao *dao.AllDao
}

var _ iface.OrdersBinaryOptionIface = (*OrdersBinaryOption)(nil)

func (w *OrdersBinaryOption) SubmitOrder(ctx context.Context, userId string, symbolCode string, session mtype.Session,
	timeMin mtype.TimeMin, strikePrice float64, option mtype.Option, betMoney int64) (orderId string, err error) {

	TX := w.allDao.Begin(ctx, nil)

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
	rowAffect, err = w.allDao.Wallet.SubAmount(TX, userId, betMoney)
	if err != nil {
		return "", err
	}
	if rowAffect == 0 {
		return "", fmt.Errorf("wallet not found, user_id: %s", userId)
	}
	// 查询余额
	walletInfo, err := w.allDao.Wallet.ItemByUserId(ctx, userId)
	if err != nil {
		return "", err
	}
	balance := walletInfo.Balance
	transId := NewOrderID(userId)
	remark := fmt.Sprintf("%s%s%d %dmin 下单", symbolCode, option.Name(), int64(strikePrice*100), session)
	detail := &model.WalletDetail{
		TransId:       transId,
		TransType:     mtype.TransTypeOrder,
		UserId:        userId,
		Amount:        -betMoney,
		Status:        mtype.TransStatusSuc,
		Remark:        remark,
		SourceKind:    mtype.TransSourceNil,
		SourceTransId: "",
		FromAccount:   "",
		ToAccount:     "",
		Balance:       balance,
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
	}
	if err = w.allDao.WalletDetail.CreateTX(TX, detail); err != nil {
		return "", err
	}
	now := time.Now()
	today, _ := strconv.ParseInt(now.Format("20060102"), 10, 64)
	orderOption := &model.OrdersBinaryOption{
		ID:             0,
		TransId:        transId,
		UserId:         userId,
		SymbolCode:     symbolCode,
		CountryCode:    mtype.CountryCodeCN,
		StrikePrice:    strikePrice,
		Option:         option,
		BetMoney:       betMoney,
		OrderTime:      now,
		Session:        session,
		Today:          today,
		SessionTimeMin: timeMin,
		SettleTime:     0,
		SettlePrice:    0,
		SettleResult:   0,
		ProfitLoss:     0,
		OrderStatus:    mtype.OrderStatusInit,
		CreateTime:     time.Now(),
		UpdateTime:     time.Now(),
	}
	if err = w.allDao.OrdersBinaryOption.CreateTX(TX, orderOption); err != nil {
		return "", err
	}
	return transId, nil
}
