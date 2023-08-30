package dao

import (
	"context"
	"gorm.io/gorm"
	"myoption/iface/myerr"
	"myoption/internal/dao/model"
	"myoption/internal/dao/mtype"
	"time"
)

type walletDetail struct {
	gdb *gorm.DB
}

func (*walletDetail) TableName() string {
	return "wallet_detail"
}

func (u *walletDetail) CreateTX(TX *gorm.DB, m *model.WalletDetail) error {
	return TX.Table(u.TableName()).Create(m).Error
}

// ItemsByUserIdOffsetLimit  .
func (u *walletDetail) ItemsByUserIdOffsetLimit(ctx context.Context, userId string, offset int, limit int) ([]model.WalletDetail, error) {
	var m []model.WalletDetail
	tx := u.gdb.WithContext(ctx).Table(u.TableName()).Where("user_id=?", userId).Order("create_time DESC").Offset(offset).Limit(limit).Find(&m)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return m, nil
}

func (u *walletDetail) UpdateStatusSucByTX(TX *gorm.DB, transactionId string, status mtype.TransStatus, amount int64) (rowsAffected int64, err error) {
	tx := TX.Table(u.TableName()).Where("transaction_id=?  and  status<? and amount=?", transactionId, mtype.TransStatusSuc, amount).Updates(map[string]interface{}{
		"status":      status,
		"update_time": time.Now(),
	})
	return tx.RowsAffected, tx.Error
}

// ItemByTransactionId .
func (u *walletDetail) ItemByTransactionId(ctx context.Context, transactionId string) (*model.WalletDetail, error) {
	var m model.WalletDetail
	tx := u.gdb.WithContext(ctx).Table(u.TableName()).Where("transaction_id=?", transactionId).Limit(1).Find(&m)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, myerr.DataNotFound
	}
	return &m, nil
}
