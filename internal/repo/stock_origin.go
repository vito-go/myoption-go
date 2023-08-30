package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/vito-go/mylog"
	"hash/crc32"
	"io"
	"math"
	"math/rand"
	"myoption/configs"
	"myoption/internal/cache"
	"myoption/internal/dao"
	"myoption/internal/dao/model"
	"myoption/internal/dao/mtype"
	"myoption/types/fd"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type stockOrigin struct {
	cache       *cache.Cache
	allDao      *dao.AllDao
	symbolCodes []string

	//updateChannel  sync.Map    // symbolCode: chan *StockOriginResult
	sharedCount    uint32      //sharedCount 的长度与syncMaps长度保持一致
	syncMaps       []*sync.Map // symbolCode: *StockOriginResult
	symbolCodeLast sync.Map    // symbolCode: []float64
}

func (s *stockOrigin) GetStockOriginResult(symbolCode string) (*StockOriginResult, bool) {
	return s.getStockOriginResult(symbolCode)
}

func (s *stockOrigin) LastPriceExist(symbolCode string, price float64) bool {

	data, ok := s.symbolCodeLast.Load(symbolCode)
	if !ok {
		return false
	}
	lastPrices := data.([]float64)
	for _, lastPrice := range lastPrices {
		if lastPrice == price {
			return true
		}
	}
	return false

}
func (s *stockOrigin) LastTodayPrices() []fd.SymbolPrice {
	var result []fd.SymbolPrice
	_, status := fd.GetMarketStatus()
	for _, symbolCode := range s.symbolCodes {
		noSymbolPrice := fd.SymbolPrice{
			MarketStatus: status,
			SymbolCode:   symbolCode,
			SymbolName:   configs.SymbolNameByCode(symbolCode),
			Exist:        false,
			Price:        0,
			TimeMin:      0,
		}
		idx := s.getIdxBySymbolCode(symbolCode)
		data, ok := s.syncMaps[idx].Load(symbolCode)
		if !ok {
			result = append(result, noSymbolPrice)
			continue
		}
		stockOriginResult := data.(*StockOriginResult)
		line := stockOriginResult.Line
		if len(line) == 0 {
			result = append(result, noSymbolPrice)
			continue
		}
		lastLineItem := line[len(line)-1]
		//if len(lastLineItem) < 2 {
		//	result = append(result, noSymbolPrice)
		//	continue
		//}

		if strconv.FormatInt(stockOriginResult.Date, 10) != time.Now().Format("20060102") {
			result = append(result, noSymbolPrice)
			continue
		}
		result = append(result, fd.SymbolPrice{
			MarketStatus: status,
			SymbolCode:   symbolCode,
			SymbolName:   configs.SymbolNameByCode(symbolCode),
			Exist:        true,
			Price:        lastLineItem.Price(),
			Day:          stockOriginResult.Date,
			TimeMin:      lastLineItem.TimeMin(),
		})
	}
	return result
}
func NewStockOrigin(allDao *dao.AllDao, ca *cache.Cache, symbolCodes []string) *stockOrigin {
	if len(symbolCodes) == 0 {
		panic("no symbolCodes")
	}
	codes := make([]string, len(symbolCodes))
	copy(codes, symbolCodes)
	sharedCount := len(symbolCodes)/100 + 1 // 根据　symbolCodes　的数量进行动态调整
	syncMaps := make([]*sync.Map, sharedCount)
	for i := 0; i < sharedCount; i++ {
		syncMaps[i] = &sync.Map{}
	}
	s := &stockOrigin{
		cache:          ca,
		allDao:         allDao,
		symbolCodes:    codes,
		sharedCount:    uint32(sharedCount),
		syncMaps:       syncMaps,
		symbolCodeLast: sync.Map{},
		//updateChannel:  sync.Map{},
	}
	go s.init()
	return s
}
func (s *stockOrigin) parseAndStoreData(symbolCode string, data *StockOriginResult) {
	idx := s.getIdxBySymbolCode(symbolCode)
	line := data.Line
	if len(line) > 0 {
		lastLineItem := line[len(line)-1]
		lastPrices, ok := s.symbolCodeLast.Load(symbolCode)
		if ok {
			s.symbolCodeLast.Store(symbolCode, append(lastPrices.([]float64), lastLineItem.Price()))
		} else {
			s.symbolCodeLast.Store(symbolCode, append([]float64{}, lastLineItem.Price()))
		}
	}
	s.syncMaps[idx].Store(symbolCode, data)
}
func (s *stockOrigin) updateSymbolCodeData(symbolCode string) {
	//golang用time.Sleep实现在一个for循环中随机休眠1到3秒。
	//要求休眠在1到2秒与2-3秒内的比例为: ４比1,并且要求休眠在1.1秒内,1.2秒内,1.3秒内,1.4秒内,1.5秒内,1.6秒内等一直到3.0秒的概率逐渐降低。

	//在这个代码中，我们首先生成了一个1到3秒的随机休眠时间，然后根据随机生成的数字在1到5之间判断应该休眠在哪个时间区间，
	//即1到2秒的四分之三和2到3秒的四分之一。接着，我们再生成一个随机的休眠时间，
	//使其在1到2秒或2到3秒之间。最后，我们根据指数分布的概率密度函数，等比例调整休眠时间，实现了在1到2秒与2到3秒的比例为4比1
	for i := 0; true; i++ {
		data, err := s._getStockOriginResult(symbolCode)
		if err != nil {
			time.Sleep(time.Second * 3)
			continue
		}
		sleep, status := fd.GetMarketStatus()
		_ = status
		//mylog.Ctx(context.Background()).WithFields("i", i, "symbolCode", symbolCode, "status", status.ToString(), "sleep", sleep.String()).Info("更新data")
		s.parseAndStoreData(symbolCode, data)
		time.Sleep(sleep)
		if sleep == 0 {
			sleepWithAdjustedSleepTime()
		}
	}
}
func sleepWithAdjustedSleepTime() {
	// 生成1到3秒的随机休眠时间
	sleepTime := time.Duration(rand.Intn(2000)+1000) * time.Millisecond
	// 判断应该休眠在哪个时间区间
	var adjustedSleepTime time.Duration
	switch rand.Intn(5) {
	case 0, 1, 2, 3:
		// 休眠在1到2秒的区间
		adjustedSleepTime = time.Duration(rand.Intn(1000)+1000) * time.Millisecond
	case 4:
		// 休眠在2到3秒的区间
		adjustedSleepTime = time.Duration(rand.Intn(1000)+2000) * time.Millisecond
	}
	// 计算指数分布的概率密度函数
	p := 1 - math.Exp(-float64(sleepTime)/2000)
	// 根据概率密度函数等比例调整休眠时间
	adjustedSleepTime = time.Duration(float64(adjustedSleepTime) * p)
	time.Sleep(adjustedSleepTime)
}

func (s *stockOrigin) init() {
	for _, code := range s.symbolCodes {
		go s.updateSymbolCodeData(code)
		go s.updateDB(mtype.CountryCodeCN, code)
	}
	go s.initSymbolsPrice()
	for {
		s.settleToady()
		exchangeTime := fd.GetExchangeTime(mtype.CountryCodeCN)
		sleep := exchangeTime.AMStart.AddDate(0, 0, 1).Sub(time.Now())
		mylog.Ctx(context.Background()).WithField("sleep", sleep.String()).Info("定时结算: 等待到下一个交易日")
		time.Sleep(sleep)
		sleep, status := fd.GetMarketStatus()
		mylog.Ctx(context.Background()).WithFields("sleep", sleep.String(), "status", status).Info("定时结算: 获取市场状态")
		time.Sleep(sleep)
		time.Sleep(time.Second)
	}
}

const rate = 0.8

func getTimeMinS() []mtype.TimeMin {
	exchangeTime := fd.GetExchangeTime(mtype.CountryCodeCN)
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
func (s *stockOrigin) settleToady() {
	mins := getTimeMinS()
	var minMap = make(map[mtype.TimeMin]bool)
	for {
		today, _ := strconv.ParseInt(time.Now().Format("20060102"), 10, 64)
		if len(minMap) == len(mins) {
			mylog.Ctx(context.Background()).WithFields("today", today, "len(minMap)", len(minMap), "len(mins)", len(mins)).Info("已经结算完毕")
			break
		}
		mylog.Ctx(context.Background()).WithFields("today", today, "len(minMap)", len(minMap), "len(mins)", len(mins)).Info("开始结算今天的数据")
		for i, min := range mins {
			if minMap[min] {
				continue
			}
			now := time.Now()
			ss := fmt.Sprintf("%d %06d", today, min)
			t, _ := time.ParseInLocation("20060102 150405", ss, time.Local)
			lateMin := now.Add(-time.Minute - time.Second*10)
			mylog.Ctx(context.Background()).WithField("lateMin", lateMin).Info("没有带结算")
			if lateMin.Before(t) {
				mylog.Ctx(context.Background()).WithFields("lateMin", lateMin, "t", t, "sleep", t.Sub(lateMin).String()).Info("没有带结算，准备休眠")
				time.Sleep(t.Sub(lateMin))
			}
			mylog.Ctx(context.Background()).WithFields("i", i, "today", today, "min", min)
			err := s.toSettle(mtype.CountryCodeCN, today, min)
			if err != nil {
				mylog.Ctx(context.Background()).WithFields("i", i, "today", today, "min", min).Errorf("结算失败 %s", err.Error())
				time.Sleep(time.Second)
				continue
			}
			minMap[min] = true

		}
	}
}

func (s *stockOrigin) initSymbolsPrice() {
	for {
		go s.initToadyTodayPrices(s.symbolCodes...)
		time.Sleep(time.Hour * 2)
	}
}
func (s *stockOrigin) toSettle(countryCode mtype.CountryCode, today int64, sessionTimeMin mtype.TimeMin) error {
	ctx := context.WithValue(context.Background(), "tid", time.Now().UnixNano())
	mylog.Ctx(ctx).WithFields("today", today, "sessionTimeMin", sessionTimeMin).Info("准备结算: ")
	for {
		//exchangeTime := fd.GetExchangeTime(countryCode)
		_, status := fd.GetMarketStatus()
		var items []model.OrdersBinaryOption
		var err error
		const limit = 1000
		//if status == fd.MarketStatusNormal {
		//	items, err = s.allDao.OrdersBinaryOption.ItemsBy(ctx, today, mtype.OrderStatusInit, sessionTimeMin, limit)
		//} else if status == fd.MarketStatusClose {
		//	if time.Now().Add(-time.Minute).Sub(exchangeTime.PMEnd) < time.Minute*15 {
		//		// 收盘十分钟仍然保持开奖
		//		items, err = s.allDao.OrdersBinaryOption.ItemsBy(ctx, today, mtype.OrderStatusInit, sessionTimeMin, limit)
		//	}
		//}
		items, err = s.allDao.OrdersBinaryOption.ItemsBy(ctx, today, mtype.OrderStatusInit, sessionTimeMin, limit)
		if err != nil {
			mylog.Ctx(ctx).WithField("status", status).Error(err.Error())
			return err
		}
		if len(items) == 0 {
			mylog.Ctx(ctx).WithFields("today", today, "sessionTimeMin", sessionTimeMin).Infof("没有待结算的订单")
			return nil
		}
		symbolCodeMap := make(map[string]*StockOriginResult)
		mylog.Ctx(ctx).WithFields("today", today, "sessionTimeMin", sessionTimeMin).Infof("一共%d条待计结算订单，准备开始结算", len(items))

		for _, item := range items {
			if _, ok := symbolCodeMap[item.SymbolCode]; ok {
				continue
			}
			data, ok := s.getStockOriginResult(item.SymbolCode)
			if !ok {
				mylog.Ctx(ctx).Infof("开奖失败，行情数据不存在,symbolCode: %s", item.SymbolCode)
				continue
			}
			symbolCodeMap[item.SymbolCode] = data
		}
		for _, item := range items {
			data, ok := symbolCodeMap[item.SymbolCode]
			if !ok {
				continue
			}
			for _, lineData := range data.Line {
				if lineData.TimeMin() == item.SessionTimeMin {
					settlePrice := lineData.Price()
					var settleResult mtype.SettleResult
					var profitLoss int64
					switch item.Option {
					case mtype.OptionCALL:
						if settlePrice > item.StrikePrice {
							settleResult = mtype.SettleResultUserWin
							profitLoss = int64(float64(item.BetMoney) * rate)
						} else {
							settleResult = mtype.SettleResultUserLost
							profitLoss = -item.BetMoney
						}
					case mtype.OptionPUT:
						if settlePrice < item.StrikePrice {
							settleResult = mtype.SettleResultUserWin
							profitLoss = int64(float64(item.BetMoney) * rate)

						} else {
							settleResult = mtype.SettleResultUserLost
							profitLoss = -item.BetMoney
						}
					default:
						mylog.Ctx(ctx).WithField("item", item).Errorf("unknown Option: %+v", item.Option)
						continue
					}
					param := settleParam{
						UserID:       item.UserId,
						BetMoney:     item.BetMoney,
						TransID:      item.TransId,
						SettlePrice:  settlePrice,
						SettleResult: settleResult,
						ProfitLoss:   profitLoss,
						SymbolCode:   item.SymbolCode,
						Option:       item.Option,
						StrikePrice:  item.StrikePrice,
						Session:      item.Session,
					}
					err = s.settle(ctx, &param)
					if err != nil {
						mylog.Ctx(ctx).WithField("param", param).Error(err.Error())
						time.Sleep(time.Second * 3)
					}
				}
			}
		}
	}

}

type settleParam struct {
	UserID       string
	BetMoney     int64
	TransID      string
	SettlePrice  float64
	SettleResult mtype.SettleResult
	ProfitLoss   int64
	SymbolCode   string
	Option       mtype.Option
	StrikePrice  float64
	Session      mtype.Session
}

func (s *stockOrigin) settle(ctx context.Context, param *settleParam) (err error) {
	userId := param.UserID
	betMoney := param.BetMoney
	transId := param.TransID
	settlePrice := param.SettlePrice
	settleResult := param.SettleResult
	profitLoss := param.ProfitLoss
	symbolCode := param.SymbolCode
	option := param.Option
	strikePrice := param.StrikePrice
	session := param.Session

	TX := s.allDao.Begin(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if e := TX.Rollback().Error; e != nil {
				mylog.Ctx(ctx).Error(e.Error())
			}
			return
		}
		if err = TX.Commit().Error; err != nil {
			return
		}
		s.cache.Wallet.DelWalletInfo(ctx, userId)
		//TODO 删除缓存
	}()
	if err = s.allDao.OrdersBinaryOption.Settle(TX, transId, settlePrice, settleResult, profitLoss, mtype.OrderStatusSuccess); err != nil {
		return err
	}
	var detailAmount int64
	if profitLoss > 0 {
		detailAmount = profitLoss + betMoney
		if _, err = s.allDao.Wallet.AddAmount(TX, userId, profitLoss+betMoney); err != nil {
			return err
		}
	} else {
		detailAmount = 0
		if _, err = s.allDao.Wallet.UpdateTimeWithLock(TX, userId); err != nil {
			return err
		}
	}
	// 查询余额
	walletInfo, err := s.allDao.Wallet.ItemByUserId(ctx, userId)
	if err != nil {
		return err
	}
	balance := walletInfo.Balance
	pl := "收益"
	if profitLoss < 0 {
		pl = "损失"
	}
	remark := fmt.Sprintf("%s%s%d %dmin %s", symbolCode, option.Name(), int64(strikePrice*100), session, pl)

	detail := &model.WalletDetail{
		TransId:       NewOrderID(userId),
		TransType:     mtype.TransTypeSettle,
		UserId:        userId,
		Amount:        detailAmount,
		Status:        mtype.TransStatusSuc,
		Remark:        remark,
		SourceKind:    mtype.TransSourceNil,
		SourceTransId: transId,
		FromAccount:   "",
		ToAccount:     "",
		Balance:       balance,
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
	}
	if err = s.allDao.WalletDetail.CreateTX(TX, detail); err != nil {
		return err
	}
	return nil
}

func (s *stockOrigin) updateStockPrice(ctx context.Context, mo *model.StockPrice) (err error) {
	TX := s.allDao.Begin(ctx, nil)
	defer func() {
		if err != nil {
			if e := TX.Rollback().Error; e != nil {
				mylog.Ctx(ctx).Error(e.Error())
			}
			return
		}
		if err = TX.Commit().Error; err != nil {
			return
		}
	}()
	if err = s.allDao.StockPrice.UpdateOrInsert(TX, mo); err != nil {
		return err
	}
	return nil
}
func (s *stockOrigin) updateDB(countryCode mtype.CountryCode, symbolCode string) {
	time.Sleep(time.Second * 10)
	ctx := context.WithValue(context.Background(), "tid", time.Now().UnixNano())
	mylog.Ctx(ctx).WithFields("countryCode", countryCode, "symbolCode", symbolCode).Info("更新db")
	for {
		idx := s.getIdxBySymbolCode(symbolCode)
		data, ok := s.syncMaps[idx].Load(symbolCode)
		if !ok {
			time.Sleep(time.Second * 10)
			continue
		}
		stockOriginResult := data.(*StockOriginResult)
		today := stockOriginResult.Date
		line := stockOriginResult.Line
		for _, priceItem := range line {
			mo := model.StockPrice{
				CountryCode: countryCode,
				SymbolCode:  symbolCode,
				Today:       today,
				TimeMin:     priceItem.TimeMin(),
				Price:       priceItem.Price(),
				Volume:      priceItem.Volume(),
				AvgPrice:    priceItem.AvgPrice(),
				Amount:      priceItem.Amount(),
				Status:      2,
				CreateTime:  time.Now(),
				UpdateTime:  time.Now(),
			}
			err := s.updateStockPrice(ctx, &mo)
			if err != nil {
				mylog.Ctx(ctx).WithField("mo", mo).Error(err.Error())
				time.Sleep(time.Second * 3)
			}
		}
		time.Sleep(time.Minute * 20)
	}
}

// GetStockOriginData .
func (s *stockOrigin) GetStockOriginData(symbolCode string) ([]byte, error) {
	client := &http.Client{}
	defer client.CloseIdleConnections()
	url := `http://yunhq.sse.com.cn:32041/v1/sh1/line/` + symbolCode + `?select=time%2Cprice%2Cvolume%2Cavg_price%2Camount%2Chighest%2Clowest&_=1678959699330`
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Add("Cache-Control", "max-age=0")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	res, err := client.Do(req)
	if err != nil {
		mylog.Ctx(context.Background()).WithField("symbolCode", symbolCode).Error(err.Error())
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		mylog.Ctx(context.Background()).WithField("symbolCode", symbolCode).Error(err.Error())
		return nil, err
	}
	return body, nil
}
func (s *stockOrigin) _getStockOriginResult(symbolCode string) (*StockOriginResult, error) {
	body, err := s.GetStockOriginData(symbolCode)
	if err != nil {
		return nil, err
	}
	var result StockOriginResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		mylog.Ctx(context.Background()).WithFields("symbolCode", symbolCode, "respBody", string(body)).Error(err.Error())
		return nil, err
	}
	return &result, nil
}

// [
// 93000,
// 3206.741,
// 2635571,
// 3206.741,
// 2059792827,
// null,
// null
// ]
// time,price,volume,avg_price,amount,highest,lowest&_=1678959699330
type StockOriginResult struct {
	Code      string     `json:"code"`
	PrevClose float64    `json:"prev_close"`
	Highest   float64    `json:"highest"`
	Lowest    float64    `json:"lowest"`
	Date      int64      `json:"date"`
	Time      int        `json:"time"`
	Total     int        `json:"total"`
	Begin     int        `json:"begin"`
	End       int        `json:"end"`
	Line      []LineData `json:"line"`
}

// LineData time,price,volume,avg_price,amount,highest,lowest
type LineData [7]float64

func (d LineData) TimeMin() mtype.TimeMin {
	return mtype.TimeMin(d[0])
}
func (d LineData) Price() float64 {
	return formatFloat64(d[1])
}
func (d LineData) Volume() int64 {
	return int64(d[2])
}
func (d LineData) AvgPrice() float64 {
	return formatFloat64(d[3])
}
func (d LineData) Amount() int64 {
	return int64(d[4])
}

func formatFloat64(f float64) float64 {
	value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", f), 64)
	return value
}

func (s *stockOrigin) getStockOriginResult(symbolCode string) (*StockOriginResult, bool) {

	idx := s.getIdxBySymbolCode(symbolCode)
	data, ok := s.syncMaps[idx].Load(symbolCode)
	if !ok {
		return nil, false
	}
	return data.(*StockOriginResult), true

}
func (s *stockOrigin) getIdxBySymbolCode(symbolCode string) uint32 {
	crc := crc32.ChecksumIEEE([]byte(symbolCode))
	idx := crc % s.sharedCount
	return idx
}
