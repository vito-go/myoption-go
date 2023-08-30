package model

import (
	"fmt"
	"myoption/configs"
	"myoption/internal/dao/mtype"
	"myoption/types/fd"
	"strconv"
	"time"
)

type OrdersBinaryOption struct {
	ID             int64              `json:"id,omitempty"`
	TransId        string             `json:"trans_id,omitempty"`
	UserId         string             `json:"user_id,omitempty"`
	SymbolCode     string             `json:"symbol_code,omitempty"`
	CountryCode    mtype.CountryCode  `json:"country_code,omitempty"`
	StrikePrice    float64            `json:"strike_price,omitempty"`
	Option         mtype.Option       `json:"option,omitempty"`
	Today          int64              `json:"today,omitempty"`
	BetMoney       int64              `json:"bet_money,omitempty"`
	OrderTime      time.Time          `json:"order_time,omitempty"`
	Session        mtype.Session      `json:"session,omitempty"`
	SessionTimeMin mtype.TimeMin      `json:"session_time_min"`
	SettleTime     int64              `json:"settle_time,omitempty"`
	SettlePrice    float64            `json:"settle_price,omitempty"`
	SettleResult   mtype.SettleResult `json:"settle_result,omitempty"`
	OrderStatus    mtype.OrderStatus  `json:"order_status,omitempty"`
	ProfitLoss     int64              `json:"profit_loss,omitempty"`
	CreateTime     time.Time          `json:"create_time,omitempty"`
	UpdateTime     time.Time          `json:"update_time,omitempty"`
}

func (o *OrdersBinaryOption) ToFd() *fd.OrdersBinaryOptionString {
	sessionTimeMin := o.SessionTimeMin.FormatToTime()
	var pl = strconv.FormatInt(o.ProfitLoss, 10)
	if o.ProfitLoss > 0 {
		pl = "+" + strconv.FormatInt(o.ProfitLoss, 10)
	}

	return &fd.OrdersBinaryOptionString{
		TransId:     o.TransId,
		CountryCode: string(o.CountryCode),
		//SymbolCode:     o.SymbolCode,
		//SymbolName:     configs.SymbolNameByCode(o.SymbolCode),
		SymbolCodeName: fmt.Sprintf("%s\n%s", o.SymbolCode, configs.SymbolNameByCode(o.SymbolCode)),
		StrikePrice:    fmt.Sprintf("%.2f", o.StrikePrice),
		Option:         o.Option.Name(),
		BetMoney:       strconv.FormatInt(o.BetMoney, 10),
		OrderTime:      o.OrderTime.Format("2006-01-02\n15:04:05"),
		SessionTimeMin: sessionTimeMin,
		Session:        o.Session.ToString(),
		SettlePrice:    fmt.Sprintf("%.2f", o.SettlePrice),
		SettleResult:   o.SettleResult.ToString(),
		OrderStatus:    o.OrderStatus.Name(),
		ProfitLoss:     pl,
	}
}

type StockPrice struct {
	CountryCode mtype.CountryCode `json:"country_code"`
	SymbolCode  string            `json:"symbol_code"`
	Today       int64             `json:"today"`
	TimeMin     mtype.TimeMin     `json:"time_min"`
	Price       float64           `json:"price"`
	Volume      int64             `json:"volume"`
	AvgPrice    float64           `json:"avg_price"`
	Amount      int64             `json:"amount"`
	Status      int64             `json:"status"` // 1: 正常 2: 已经更新
	CreateTime  time.Time         `json:"create_time,omitempty"`
	UpdateTime  time.Time         `json:"update_time,omitempty"`
}

func (StockPrice) TableName() string {
	return "stock_price"
}
