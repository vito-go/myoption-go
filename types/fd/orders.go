package fd

import (
	"myoption/internal/dao/mtype"
)

type OrdersBinaryOption struct {
	TransId        string             `json:"transId"`
	SymbolCode     string             `json:"symbolCode"`
	CountryCode    mtype.CountryCode  `json:"country_code"`
	SymbolName     string             `json:"symbolName"`
	StrikePrice    float64            `json:"strikePrice"`
	Option         mtype.Option       `json:"option"`
	BetMoney       int64              `json:"betMoney"`
	OrderTime      string             `json:"orderTime"`
	SessionTimeMin mtype.TimeMin      `json:"sessionTimMin"`
	Session        mtype.Session      `json:"session"`
	SettlePrice    float64            `json:"settlePrice"`
	SettleResult   mtype.SettleResult `json:"settleResult"`
	OrderStatus    mtype.OrderStatus  `json:"orderStatus"`
	ProfitLoss     int64              `json:"profitLoss"`
}

type OrdersBinaryOptionString struct {
	TransId string `json:"transId"`
	//SymbolCode     string `json:"symbolCode"`
	CountryCode string `json:"countryCode"`
	//SymbolName     string `json:"symbolName"`
	SymbolCodeName string `json:"symbolCodeName"`
	StrikePrice    string `json:"strikePrice"`
	Option         string `json:"option"`
	BetMoney       string `json:"betMoney"`
	OrderTime      string `json:"orderTime"`
	SessionTimeMin string `json:"sessionTimMin"`
	Session        string `json:"session"`
	SettlePrice    string `json:"settlePrice"`
	SettleResult   string `json:"settleResult"`
	OrderStatus    string `json:"orderStatus"`
	ProfitLoss     string `json:"profitLoss"`
}
