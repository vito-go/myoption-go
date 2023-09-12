package user

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/vito-go/mylog"
	"myoption/configs"
	"myoption/iface/myerr"
	"myoption/internal/dao/mtype"
	"myoption/internal/httpsrv/handler"
	"myoption/internal/httpsrv/handler/httperr"
	"myoption/internal/repo"
	"myoption/pkg/resp"
	"myoption/pkg/util/slice"
	"myoption/types/fd"
	"net/http"
	"time"
)

type SubmitOrder struct {
	RepoClient *repo.Client
}

func (h *SubmitOrder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler.ServeHTTP(w, r, h)

}

func (h *SubmitOrder) GetParam(r *http.Request) (*handler.ReqParam, error) {
	return handler.GetParamBodyEnc(r)
}

type submitOrderParam struct {
	SymbolCode  string        `json:"symbolCode,omitempty"`
	StrikePrice float64       `json:"strikePrice,omitempty"`
	Option      mtype.Option  `json:"option,omitempty"`
	BetMoney    int64         `json:"betMoney,omitempty"`
	Session     mtype.Session `json:"session,omitempty"`
}

func (h *SubmitOrder) GetRespBody(ctx context.Context, p *handler.ReqParam) *resp.HTTPBody {
	var param submitOrderParam
	err := json.Unmarshal(p.Body, &param)
	if err != nil {
		return resp.Err(ctx, httperr.ErrParam.Error())
	}

	userId := p.UserId
	userInfo, err := h.RepoClient.UserCli.GetUserInfoByUserId(ctx, userId)
	if err != nil {
		if err == myerr.DataNotFound {
			return resp.Err(ctx, httperr.ErrUserNotFound.Error())
		}
		mylog.Ctx(ctx).Error(err.Error())
		return resp.Err(ctx, httperr.ErrData.Error())
	}

	if userInfo.Ed25519PubKey != base64.StdEncoding.EncodeToString(p.SignPubKeyBytes) {
		return resp.Err(ctx, "sign pub key error")
	}
	if userInfo.X25519PubKey != base64.StdEncoding.EncodeToString(p.ClientPubKey.Bytes()) {
		return resp.Err(ctx, "sign pub key error")
	}
	_, status := fd.GetMarketStatus()
	if status != fd.MarketStatusNormal {
		return resp.Err(ctx, status.ToString())
	}

	if result := h.checkParam(ctx, &param); result != nil {
		return result
	}

	walletInfo, err := h.RepoClient.Wallet.GetWalletInfoByUserId(ctx, userId)
	if err != nil {
		mylog.Ctx(ctx).Error(err.Error())
		return resp.Err(ctx, httperr.ErrData.Error())
	}

	if walletInfo.Balance < param.BetMoney {
		return resp.Err(ctx, "余额不足")
	}
	now := time.Now()
	exchangeTime := fd.GetExchangeTime(mtype.CountryCodeCN)
	timeMinTime := now.Add(time.Minute * time.Duration(param.Session))
	if param.Session == mtype.Session0 {
		timeMinTime = exchangeTime.PMEnd
	}
	if !exchangeTime.InTradingTime(timeMinTime) {
		return resp.Err(ctx, "场次不在交易时间内")
	}
	timeMin := mtype.GetTimeMin(timeMinTime)
	orderId, err := h.RepoClient.OrdersBinaryOption.SubmitOrder(ctx, userId, param.SymbolCode, param.Session, timeMin, param.StrikePrice, param.Option, param.BetMoney)
	if err != nil {
		mylog.Ctx(ctx).Error(err.Error())
		return resp.Err(ctx, httperr.ErrData.Error())
	}
	return resp.DataOK(ctx, map[string]string{"orderId": orderId})
}

func (h *SubmitOrder) checkParam(ctx context.Context, param *submitOrderParam) *resp.HTTPBody {
	if param.StrikePrice <= 0 {
		return resp.Err(ctx, "下单价格非法")
	}
	ok := h.RepoClient.StockData.LastPriceExist(param.SymbolCode, param.StrikePrice)
	if !ok {
		return resp.Err(ctx, "下单价格不存在")
	}
	if param.SymbolCode == "" {
		return resp.Err(ctx, "SymbolCode 不能为空")
	}
	if !slice.IsInSlice(configs.Symbols, param.SymbolCode) {
		return resp.Err(ctx, "SymbolCode 不存在")
	}
	if !param.Option.Check() {
		return resp.Err(ctx, "option　错误")
	}
	if param.BetMoney < 10 {
		return resp.Err(ctx, "金额非法")
	}
	if param.BetMoney > 10000 {
		return resp.Err(ctx, "金额超限")
	}
	if !param.Session.Check() {
		return resp.Err(ctx, "Session　错误")
	}
	return nil
}
