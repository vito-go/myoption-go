package repo

import (
	"context"
	"github.com/vito-go/mylog"
	"myoption/internal/dao/model"
	"myoption/internal/dao/mtype"
	"strconv"
	"time"
)

func (s *stockOrigin) initToadyTodayPrices(symbolCodes ...string) {
	ctx := context.Background()
	ifSet, err := s.cache.DistributeDoOnce("initToadySymbols")
	if err != nil {
		panic(err)
	}
	lock := s.cache.NewDLock("stockOrigin.initToadySymbols")
	err = lock.Lock(ctx)
	if err != nil {
		panic(err)
	}
	defer lock.UnLock(ctx)
	if !ifSet {
		mylog.Ctx(ctx).Warn("未获取初始化锁，节点忽略初始化")
		return
	}
	mylog.Ctx(ctx).WithField("symbolCodes", symbolCodes).Info("开始进行代码价格初始化")
	for _, code := range symbolCodes {
		now := time.Now()
		todayInt, _ := strconv.ParseInt(now.Format("20060102"), 10, 64)
		var count int64
		countryCode := mtype.CountryCodeCN
		count, err = s.allDao.StockPrice.CountBy(ctx, countryCode, code, todayInt)
		if err != nil {
			panic(err)
		}
		if count > 0 {
			mylog.Ctx(ctx).WithFields("countryCode", countryCode, "code", code, "todayInt", todayInt).Warn("已初始化，忽略初始化")
			continue
		}
		amStartTime := time.Date(now.Year(), now.Month(), now.Day(), 9, 30, 0, 0, now.Location())
		amEndTime := time.Date(now.Year(), now.Month(), now.Day(), 11, 30, 0, 0, now.Location())
		pmStartTime := time.Date(now.Year(), now.Month(), now.Day(), 13, 0, 0, 0, now.Location())
		pmEndTime := time.Date(now.Year(), now.Month(), now.Day(), 15, 0, 0, 0, now.Location())
		var models []model.StockPrice
		// am
		for i := 0; true; i++ {
			t := amStartTime.Add(time.Minute * time.Duration(i))
			if t.After(amEndTime) {
				break
			}
			min := mtype.GetTimeMin(t)
			mo := model.StockPrice{
				CountryCode: mtype.CountryCodeCN,
				SymbolCode:  code,
				Today:       todayInt,
				TimeMin:     min,
				Price:       0,
				Volume:      0,
				AvgPrice:    0,
				Amount:      0,
				Status:      0,
				CreateTime:  now,
				UpdateTime:  now,
			}
			models = append(models, mo)

		}
		// pm
		// 下午的价格时间从13:00开始
		for i := 1; true; i++ {
			t := pmStartTime.Add(time.Minute * time.Duration(i))
			if t.After(pmEndTime) {
				break
			}
			min := mtype.GetTimeMin(t)
			mo := model.StockPrice{
				CountryCode: mtype.CountryCodeCN,
				SymbolCode:  code,
				Today:       todayInt,
				TimeMin:     min,
				Price:       0,
				Volume:      0,
				AvgPrice:    0,
				Amount:      0,
				Status:      1,
				CreateTime:  now,
				UpdateTime:  now,
			}
			models = append(models, mo)
		}
		err = s.allDao.StockPrice.CreateTXs(context.Background(), models...)
		if err != nil {
			panic(err)
		}
		mylog.Ctx(ctx).WithFields("countryCode", countryCode, "code", code, "todayInt", todayInt).Info("已初始化完成")
	}
	mylog.Ctx(ctx).WithField("symbolCodes", symbolCodes).Info("代码价格初始化完毕")

}
