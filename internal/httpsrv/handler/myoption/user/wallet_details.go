package user

import (
	"context"
	"encoding/base64"
	"github.com/vito-go/mylog"
	"myoption/iface/myerr"
	"myoption/internal/httpsrv/handler"
	"myoption/internal/httpsrv/handler/httperr"
	"myoption/internal/repo"
	"myoption/pkg/resp"
	"myoption/types/fd"
	"net/http"
	"strconv"
)

type WalletDetails struct {
	RepoClient *repo.Client
}

func (h *WalletDetails) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler.ServeHTTP(w, r, h)

}

func (h *WalletDetails) GetParam(r *http.Request) (*handler.ReqParam, error) {
	return handler.GetParam1(r, "offset", "limit")
}

type walletDetailsData struct {
	Items      []*fd.WalletDetail `json:"items"`
	FieldNames []FieldName        `json:"fieldNames"`
}

func (h *WalletDetails) GetRespBody(ctx context.Context, p *handler.ReqParam) *resp.HTTPBody {
	offset, err := strconv.ParseInt(p.QueryForm.Get("offset"), 10, 64)
	if err != nil {
		return resp.Err(ctx, "offset error")
	}
	limit, err := strconv.ParseInt(p.QueryForm.Get("limit"), 10, 64)
	if err != nil {
		return resp.Err(ctx, "limit error")
	}
	if offset < 0 || limit < 0 {
		return resp.Err(ctx, "offset or limit error")
	}
	// check size
	if limit > 100 {
		return resp.Err(ctx, "limit error")
	}
	// check offset
	if offset > 1000 {
		return resp.Err(ctx, "offset error")
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

	result, err := h.RepoClient.AllDao().WalletDetail.ItemsByUserIdOffsetLimit(ctx, userId, int(offset), int(limit))
	if err != nil {
		mylog.Ctx(ctx).Error(err.Error())
		return resp.Err(ctx, httperr.ErrData.Error())
	}
	items := make([]*fd.WalletDetail, 0, len(result))
	for _, option := range result {
		items = append(items, option.ToFd())
	}
	return resp.DataOK(ctx, walletDetailsData{Items: items,
		FieldNames: fieldWalletDetails})
}
