package dao

import (
	"context"
	"gorm.io/gorm"
	"myoption/internal/dao/model"
	"myoption/internal/dao/mtype"
	"time"
)

type stockPrice struct {
	gdb *gorm.DB
}

func (u *stockPrice) TableName() string {
	return "stock_price"
}

func (u *stockPrice) CreateTXs(ctx context.Context, items ...model.StockPrice) error {
	if len(items) == 0 {
		return nil
	}
	return u.gdb.WithContext(ctx).Table(u.TableName()).Create(&items).Error
}

// CountBy  .
func (u *stockPrice) CountBy(ctx context.Context, countryCode mtype.CountryCode, symbolCode string, today int64) (int64, error) {
	var count int64
	tx := u.gdb.WithContext(ctx).Table(u.TableName()).Select([]string{"count(*)"}).Where("country_code=? AND symbol_code=? AND today=?", countryCode, symbolCode, today).Count(&count)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return count, nil
}

// UpdateOrInsert   price | volume | avg_price | amount | status .
func (u *stockPrice) UpdateOrInsert(TX *gorm.DB, mo *model.StockPrice) error {
	// update
	var countryCode = mo.CountryCode
	var symbolCode = mo.SymbolCode
	var today = mo.Today
	var timeMin = mo.TimeMin
	updates := map[string]interface{}{
		"price":       mo.Price,
		"volume":      mo.Volume,
		"avg_price":   mo.AvgPrice,
		"amount":      mo.Amount,
		"status":      mo.Status,
		"update_time": time.Now(),
	}
	tx := TX.Table(u.TableName()).Where("country_code=? AND symbol_code=? AND today=? AND time_min=?",
		countryCode, symbolCode, today, timeMin).Updates(updates)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		// insert
		mo.CreateTime = time.Now()
		mo.UpdateTime = time.Now()
		return TX.Table(u.TableName()).Create(mo).Error
	}
	return nil
}
