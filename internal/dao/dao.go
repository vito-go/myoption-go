package dao

import (
	"context"
	"database/sql"
	"gorm.io/gorm"
)

type AllDao struct {
	gdb      *gorm.DB
	UserInfo *userInfo
	UserKey  *userKey

	Wallet             *wallet
	WalletDetail       *walletDetail
	OrdersBinaryOption *ordersBinaryOption
	StockPrice         *stockPrice

	// 各种表
}

func (a *AllDao) GDB() *gorm.DB {
	return a.gdb
}

func (a *AllDao) DB() (*sql.DB, error) {
	return a.gdb.DB()
}

func NewAllDao(gdb *gorm.DB) *AllDao {
	return &AllDao{
		gdb:      gdb,
		UserInfo: &userInfo{gdb: gdb},
		UserKey:  &userKey{gdb: gdb},

		Wallet:             &wallet{gdb: gdb},
		WalletDetail:       &walletDetail{gdb: gdb},
		OrdersBinaryOption: &ordersBinaryOption{gdb: gdb},
		StockPrice:         &stockPrice{gdb: gdb},
	}
}

func (a *AllDao) Begin(ctx context.Context, options ...*sql.TxOptions) *gorm.DB {
	return a.gdb.WithContext(ctx).Begin(options...)
}
