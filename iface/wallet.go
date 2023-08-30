package iface

import (
	"context"
	"myoption/internal/dao/model"
	"myoption/internal/dao/mtype"
	"myoption/types/fd"
)

type WalletIface interface {
	RechargeNotification(ctx context.Context, param *model.WalletDetail) (transactionId string, err error)
	AddWithSuccess(ctx context.Context, transactionId, userId string, amount int64, unFrozenUserId string) (err error)
	ConsumerWallet(ctx context.Context, detail *model.WalletDetail, sub bool) (transactionId string, err error)
	GetWalletInfoByUserId(ctx context.Context, userId string) (*model.WalletInfo, error)
}

type StockDataIface interface {
	LastTodayPrices() []fd.SymbolPrice
	LastPriceExist(symbolCode string, price float64) bool
	GetStockOriginData(symbolCode string) ([]byte, error)
}
type OrdersBinaryOptionIface interface {
	SubmitOrder(ctx context.Context, userId string, symbolCode string, session mtype.Session, timeMin mtype.TimeMin,
		strikePrice float64, option mtype.Option, betMoney int64) (orderId string, err error)
}
