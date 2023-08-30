package web

import (
	"context"
	"myoption/internal/httpsrv/handler"
	"myoption/internal/repo"
	"myoption/pkg/resp"
	"myoption/types/fd"
	"net/http"
	"time"
)

type ProductList struct {
	RepoClient *repo.Client
}
type productListData struct {
	Items      []fd.SymbolPrice `json:"items,omitempty"`
	UpdateTime int64            `json:"updateTime,omitempty"`
}

func (h *ProductList) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler.ServeHTTP(w, r, h)

}

func (h *ProductList) GetParam(r *http.Request) (*handler.ReqParam, error) {
	return handler.GetParam(r)

}

func (h *ProductList) GetRespBody(ctx context.Context, p *handler.ReqParam) *resp.HTTPBody {
	lastPrices := h.RepoClient.StockData.LastTodayPrices()
	return resp.DataOK(ctx, productListData{Items: lastPrices, UpdateTime: time.Now().UnixMilli()})

}
