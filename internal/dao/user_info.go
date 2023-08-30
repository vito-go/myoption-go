package dao

import (
	"context"
	"gorm.io/gorm"
	"myoption/iface/myerr"
	"myoption/internal/dao/model"
	"myoption/internal/dao/mtype"
	"time"
)

type userInfo struct {
	gdb *gorm.DB
}

func (u *userInfo) TableName() string {
	return "user_info"
}

func (u *userInfo) UpdateStatusByUserId(gdb *gorm.DB, userId string, status mtype.UserStatus) error {
	tx := gdb.Table(u.TableName()).Where("user_id=?", userId).Updates(map[string]interface{}{
		"status":      status,
		"update_time": time.Now(),
	})
	if tx.RowsAffected <= 0 {
		return myerr.DataNotFound
	}
	return tx.Error
}

// ItemByUserId  .
func (u *userInfo) ItemByUserId(ctx context.Context, userId string) (*model.UserInfo, error) {
	var m model.UserInfo
	tx := u.gdb.WithContext(ctx).Table(u.TableName()).Where("user_id=?", userId).Limit(1).Find(&m)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, myerr.DataNotFound
	}
	return &m, nil
}

func (u *userInfo) ItemsByOffsetLimit(ctx context.Context, offset, limit int) ([]model.UserInfo, error) {
	var result []model.UserInfo
	err := u.gdb.WithContext(ctx).Table(u.TableName()).Order("id DESC").
		Offset(offset).Limit(limit).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *userInfo) AllUserIds(ctx context.Context) ([]string, error) {
	var result []string
	err := u.gdb.WithContext(ctx).Table(u.TableName()).Select([]string{"user_id"}).Where("id>0").Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *userInfo) CountAll(ctx context.Context) (int64, error) {
	var result int64
	tx := u.gdb.WithContext(ctx).Table(u.TableName()).Select([]string{"count(*)"}).Find(&result)
	return result, tx.Error
}

func (u *userInfo) CreateTX(TX *gorm.DB, m *model.UserInfo) error {
	return TX.Table(u.TableName()).Create(m).Error
}
