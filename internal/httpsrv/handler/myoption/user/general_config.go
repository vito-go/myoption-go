package user

import (
	"context"
	"myoption/internal/httpsrv/handler"
	"myoption/internal/repo"
	"myoption/pkg/resp"
	"net/http"
	"time"
)

type GeneralConfig struct {
	RepoClient *repo.Client
}
type generalConfigData struct {
	UpdateTime int64  `json:"updateTime"`
	RuleInfo   string `json:"ruleInfo"`
}

const ruleInfo = "产品介绍：\n本产品以以上证指数或股票为标的，用于投注上证指数的价格。用户可选择下注金币数量（10-2000金币），以及不同场次（2分钟、3分钟、5分钟、10分钟、20分钟、30分钟、60分钟、全天）。投注方向有两个选项：看涨和看跌。\n\n举例说明：当前时间为10:01:00，实时价格指数为3230.64。用户选择5分钟场次，看涨，并下注10金币。\n\n到10:06:00，如果价格高于3230.64（如3233.56），用户盈利10金币；如果价格小于等于3230.64，用户亏损10金币。"

type ExchangeTime struct {
}

func (h *GeneralConfig) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler.ServeHTTP(w, r, h)

}

func (h *GeneralConfig) GetParam(r *http.Request) (*handler.ReqParam, error) {
	return handler.GetParam(r)
}

func (h *GeneralConfig) GetRespBody(ctx context.Context, p *handler.ReqParam) *resp.HTTPBody {
	respData := &generalConfigData{UpdateTime: time.Now().UnixMilli(), RuleInfo: ruleInfo}
	return resp.DataOK(ctx, respData)
}
