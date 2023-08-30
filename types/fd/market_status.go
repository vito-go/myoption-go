package fd

import (
	"myoption/internal/dao/mtype"
	"time"
)

type MarketStatus int

func (m MarketStatus) ToString() string {
	switch m {
	case MarketStatusNormal:
		return "交易中"
	case MarketStatusClose:
		return "已收盘"
	case MarketWaitToOpen:
		return "待开盘"
	case MarketStatusPause:
		return "休市"
	case MarketStatusWeekend:
		return "周末休市"
	case MarketStatusHoliday:
		return "假日休市"
	default:
		return "-"
	}
}

const (
	MarketStatusNormal MarketStatus = iota + 1
	MarketStatusClose
	MarketWaitToOpen
	MarketStatusPause
	MarketStatusWeekend

	MarketStatusHoliday
)

type ExchangeTime struct {
	AMStart       time.Time
	AMEnd         time.Time
	PMStart       time.Time
	PMEnd         time.Time
	TomorrowStart time.Time
}

func (e *ExchangeTime) InTradingTime(t time.Time) bool {
	if t.UnixMilli() >= e.AMStart.UnixMilli() && t.UnixMilli() <= e.AMEnd.UnixMilli() {
		return true
	}
	if t.UnixMilli() >= e.PMStart.UnixMilli() && t.UnixMilli() <= e.PMEnd.UnixMilli() {
		return true
	}
	return false
}
func GetExchangeTime(countryCode mtype.CountryCode) *ExchangeTime {
	return getExchangeTime(countryCode)
}

func GetTimeMinS() []mtype.TimeMin {
	exchangeTime := GetExchangeTime(mtype.CountryCodeCN)
	amStartTime := exchangeTime.AMStart
	amEndTime := exchangeTime.AMEnd
	pmStartTime := exchangeTime.PMStart
	pmEndTime := exchangeTime.PMEnd
	var timeMinS []mtype.TimeMin
	// am
	for i := 0; true; i++ {
		t := amStartTime.Add(time.Minute * time.Duration(i))
		if t.After(amEndTime) {
			break
		}
		min := mtype.GetTimeMin(t)
		timeMinS = append(timeMinS, min)
	}
	// pm
	// 下午的价格时间从13:00开始
	for i := 1; true; i++ {
		t := pmStartTime.Add(time.Minute * time.Duration(i))
		if t.After(pmEndTime) {
			break
		}
		min := mtype.GetTimeMin(t)
		timeMinS = append(timeMinS, min)
	}
	return timeMinS
}
func getExchangeTime(countryCode mtype.CountryCode) *ExchangeTime {
	// todo
	now := time.Now()
	AMStart := time.Date(now.Year(), now.Month(), now.Day(), 9, 30, 0, 0, time.Local)
	AMEnd := time.Date(now.Year(), now.Month(), now.Day(), 11, 30, 0, 0, time.Local)
	PMStart := time.Date(now.Year(), now.Month(), now.Day(), 13, 0, 0, 0, time.Local)
	PMEnd := time.Date(now.Year(), now.Month(), now.Day(), 15, 0, 0, 0, time.Local)
	TomorrowStart := time.Date(now.Year(), now.Month(), now.Day()+1, 9, 30, 0, 0, time.Local)
	return &ExchangeTime{
		AMStart:       AMStart,
		AMEnd:         AMEnd,
		PMStart:       PMStart,
		PMEnd:         PMEnd,
		TomorrowStart: TomorrowStart,
	}
}

// GetMarketStatus 获取市场状态
func GetMarketStatus() (sleep time.Duration, status MarketStatus) {
	exchangeTime := getExchangeTime(mtype.CountryCodeCN)

	now := time.Now()
	amStart := exchangeTime.AMStart
	amEnd := exchangeTime.AMEnd
	pmStart := exchangeTime.PMStart
	pmEnd := exchangeTime.PMEnd
	tomorrowStart := exchangeTime.TomorrowStart
	// 如果是周末，休眠到周一
	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		// Sleep until Monday
		daysUntilMonday := (8 - int(now.Weekday())) % 7
		if daysUntilMonday == 0 {
			daysUntilMonday = 7
		}
		targetTime := time.Date(now.Year(), now.Month(), now.Day()+daysUntilMonday, 9, 30, 0, 0, time.Local)
		sleep = targetTime.Sub(now)
		return sleep, MarketStatusWeekend
	}

	//　FIXME　法定节假日还没算
	// 如果是法定节假日，休眠到下一个交易日
	if now.Before(amStart) {
		sleep = amStart.Sub(time.Now())
		return sleep, MarketWaitToOpen
	} else if now.After(amEnd) && now.Before(pmStart) {
		sleep = pmStart.Sub(time.Now())
		return sleep, MarketStatusPause
	} else if now.After(pmEnd) {
		sleep = tomorrowStart.Sub(time.Now())
		return sleep, MarketStatusClose
	}
	return 0, MarketStatusNormal
}
