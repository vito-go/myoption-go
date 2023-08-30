package dao

import (
	"context"
	"gorm.io/gorm"
	"myoption/iface/myerr"
	"myoption/internal/dao/model"
	"time"
)

type wallet struct {
	gdb *gorm.DB
}

func (*wallet) TableName() string {
	return "wallet_info"
}

func (u *wallet) CreateTX(TX *gorm.DB, m *model.WalletInfo) error {
	return TX.Table(u.TableName()).Create(m).Error
}
func (u *wallet) ItemByUserId(ctx context.Context, userId string) (*model.WalletInfo, error) {
	var m model.WalletInfo
	tx := u.gdb.WithContext(ctx).Table(u.TableName()).Where("user_id=?  ", userId).Find(&m)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, myerr.DataNotFound
	}
	return &m, nil
}
func (u *wallet) AddAmount(TX *gorm.DB, userId string, amount int64) (rowsAffected int64, err error) {
	tx := TX.Table(u.TableName()).Where("user_id=?  ", userId).Updates(map[string]interface{}{
		"total_amount": gorm.Expr("total_amount+?", amount),
		"balance":      gorm.Expr("balance+?", amount),
		"update_time":  time.Now(),
	})
	return tx.RowsAffected, tx.Error
}

func (u *wallet) UpdateTimeWithLock(TX *gorm.DB, userId string) (rowsAffected int64, err error) {
	tx := TX.Table(u.TableName()).Where("user_id=?", userId).Updates(map[string]interface{}{
		"update_time": gorm.Expr("update_time"),
	})
	return tx.RowsAffected, tx.Error
}

func (u *wallet) FreezeAmount(TX *gorm.DB, userId string, amount int64) (rowsAffected int64, err error) {
	tx := TX.Table(u.TableName()).Where("user_id=?", userId).Updates(map[string]interface{}{
		"balance":       gorm.Expr("balance-?", amount),
		"frozen_amount": gorm.Expr("frozen_amount+?", amount),
		"update_time":   time.Now(),
	})
	return tx.RowsAffected, tx.Error
}

func (u *wallet) UnFrozenAmount(TX *gorm.DB, userId string, amount int64) (rowsAffected int64, err error) {
	tx := TX.Table(u.TableName()).Where("user_id=?", userId).Updates(map[string]interface{}{
		"total_amount":  gorm.Expr("total_amount-?", amount),
		"frozen_amount": gorm.Expr("frozen_amount-?", amount),
		"update_time":   time.Now(),
	})
	return tx.RowsAffected, tx.Error
}

func (u *wallet) SubAmount(TX *gorm.DB, userId string, amount int64) (rowsAffected int64, err error) {
	tx := TX.Table(u.TableName()).Where("user_id=?  ", userId).Updates(map[string]interface{}{
		"total_amount": gorm.Expr("total_amount-?", amount),
		"balance":      gorm.Expr("balance-?", amount),
		"update_time":  time.Now(),
	})
	return tx.RowsAffected, tx.Error
}
