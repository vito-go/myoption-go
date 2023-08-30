package dao

import (
	"context"
	"gorm.io/gorm"
	"myoption/internal/dao/model"
	"myoption/internal/dao/mtype"
	"time"
)

type ordersBinaryOption struct {
	gdb *gorm.DB
}

func (u *ordersBinaryOption) TableName() string {
	return "orders_binary_option"
}

func (u *ordersBinaryOption) CreateTX(TX *gorm.DB, m *model.OrdersBinaryOption) error {
	return TX.Table(u.TableName()).Create(m).Error
}

// ItemsByUserIdOffsetLimit  .
func (u *ordersBinaryOption) ItemsByUserIdOffsetLimit(ctx context.Context, userId string, offset int, limit int) ([]model.OrdersBinaryOption, error) {
	var m []model.OrdersBinaryOption
	tx := u.gdb.WithContext(ctx).Table(u.TableName()).Where("user_id=?", userId).Order("create_time DESC").Offset(offset).Limit(limit).Find(&m)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return m, nil
}

// ItemsBy  .
func (u *ordersBinaryOption) ItemsBy(ctx context.Context, today int64, orderStatus mtype.OrderStatus, sessionTimeMin mtype.TimeMin, limit int) ([]model.OrdersBinaryOption, error) {
	var m []model.OrdersBinaryOption
	tx := u.gdb.WithContext(ctx).Table(u.TableName()).Where("today=? AND order_status=? AND session_time_min<=?", today, orderStatus, sessionTimeMin).Limit(limit).Find(&m)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return m, nil
}

// Settle  .
func (u *ordersBinaryOption) Settle(TX *gorm.DB, transId string, settlePrice float64, settleResult mtype.SettleResult, profitLoss int64, orderStatus mtype.OrderStatus) error {
	now := time.Now()
	tx := TX.Table(u.TableName()).Where("trans_id=?", transId).Updates(map[string]interface{}{
		"order_status":  orderStatus,
		"settle_price":  settlePrice,
		"settle_result": settleResult,
		"profit_loss":   profitLoss,
		"update_time":   now,
		"settle_time":   now.UnixMilli(),
	})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
