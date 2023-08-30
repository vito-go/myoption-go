package dao

import (
	"context"
	"gorm.io/gorm"
	"myoption/iface/myerr"
	"myoption/internal/dao/model"
)

type userKey struct {
	gdb *gorm.DB
}

func (u *userKey) TableName() string {
	return "user_key"
}

// ItemByUserId  .
func (u *userKey) ItemByUserId(ctx context.Context, userId string) (*model.UserKey, error) {
	var m model.UserKey
	tx := u.gdb.WithContext(ctx).Table(u.TableName()).Where("user_id=?", userId).Limit(1).Find(&m)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, myerr.DataNotFound
	}
	return &m, nil
}

func (u *userKey) CountAll(ctx context.Context) (int64, error) {
	var result int64
	tx := u.gdb.WithContext(ctx).Table(u.TableName()).Select([]string{"count(*)"}).Find(&result)
	return result, tx.Error
}

func (u *userKey) CreateTX(TX *gorm.DB, m *model.UserKey) error {
	return TX.Table(u.TableName()).Create(m).Error
}
