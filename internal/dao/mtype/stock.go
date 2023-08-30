package mtype

import (
	"fmt"
	"strconv"
	"time"
)

type Option int

const (
	OptionCALL = Option(1)
	OptionPUT  = Option(2)
)

func (o Option) Check() bool {
	switch o {
	case OptionCALL, OptionPUT:
		return true

	default:
		return false

	}
}
func (o Option) Name() string {
	switch o {
	case OptionCALL:
		return "看涨"
	case OptionPUT:
		return "看跌"
	default:
		return ""
	}
}

type OrderStatus int

const (
	OrderStatusInit    = OrderStatus(0)
	OrderStatusSuccess = OrderStatus(200)
)

func (s OrderStatus) Name() string {
	switch s {
	case OrderStatusInit:
		return "待结算"
	case OrderStatusSuccess:
		return "已结算"
	default:
		return ""
	}
}

type SettleResult int

const (
	SettleResultInit     = SettleResult(0)
	SettleResultUserWin  = SettleResult(1)
	SettleResultUserLost = SettleResult(2)
)

func (s SettleResult) ToString() string {
	switch s {
	case SettleResultInit:
		return "初始化"
	case SettleResultUserWin:
		return "收益"
	case SettleResultUserLost:
		return "损失"
	default:
		return ""
	}
}

type Session int64

const (
	Session0  = Session(0)
	Session2  = Session(2)
	Session3  = Session(3)
	Session5  = Session(5)
	Session10 = Session(10)
	Session15 = Session(15)
	Session20 = Session(20)
	Session30 = Session(30)
	Session60 = Session(60)
)

func (o Session) Check() bool {
	switch o {
	case Session0, Session2, Session3, Session5, Session10, Session15, Session20, Session30, Session60:
		return true
	default:
		return false
	}
}
func (o Session) ToString() string {
	switch o {
	case Session0:
		return "全天"
	case Session2:
		return "2分钟"

	case Session3:
		return "3分钟"

	case Session5:
		return "5分钟"

	case Session10:
		return "10分钟"

	case Session15:
		return "15分钟"

	case Session20:
		return "20分钟"

	case Session30:
		return "30分钟"

	case Session60:
		return "60分钟"
	default:
		return "unknown"
	}
}

type CountryCode string

const (
	CountryCodeCN = CountryCode("CN")
)

type TimeMin int64

func (m TimeMin) FormatToTime() string {
	s := fmt.Sprintf("20060102 %06d", m)
	t, err := time.ParseInLocation("20060102 150405", s, time.Local)
	if err != nil {
		return ""
	}
	return t.Format("15:04")
}

func GetTimeMin(t time.Time) TimeMin {
	timeMin, _ := strconv.ParseInt(t.Format("1504")+"00", 10, 64)
	return TimeMin(timeMin)
}
