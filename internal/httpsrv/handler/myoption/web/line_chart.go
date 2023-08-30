package web

import (
	"context"
	"encoding/json"
	"github.com/vito-go/mylog"
	"myoption/internal/httpsrv/handler"
	"myoption/internal/repo"
	"myoption/pkg/resp"
	"net/http"
)

type LineChart struct {
	RepoClient *repo.Client
}

func (h *LineChart) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler.ServeHTTP(w, r, h)

}

func (h *LineChart) GetParam(r *http.Request) (*handler.ReqParam, error) {
	return handler.GetParam(r, "symbolCode")

}

func (h *LineChart) GetRespBody(ctx context.Context, p *handler.ReqParam) *resp.HTTPBody {
	symbolCode := p.Get("symbolCode")
	bytes, err := h.RepoClient.StockData.GetStockOriginData(symbolCode)
	if err != nil {
		return resp.Err(ctx, err.Error())
	}
	var sse Sse
	err = json.Unmarshal(bytes, &sse)
	if err != nil {
		mylog.Ctx(context.Background()).Error(err)
	}
	line := sse.buildLines()
	return resp.DataOK(ctx, line)
}
