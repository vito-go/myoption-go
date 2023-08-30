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
	"net/http"
)

type MyBalance struct {
	RepoClient *repo.Client
}

func (h *MyBalance) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler.ServeHTTP(w, r, h)

}

func (h *MyBalance) GetParam(r *http.Request) (*handler.ReqParam, error) {
	return handler.GetParam1(r)
}

func (h *MyBalance) GetRespBody(ctx context.Context, p *handler.ReqParam) *resp.HTTPBody {

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

	result, err := h.RepoClient.Wallet.GetWalletInfoByUserId(ctx, userId)
	if err != nil {
		mylog.Ctx(ctx).Error(err.Error())
		return resp.Err(ctx, httperr.ErrData.Error())
	}
	return resp.DataOK(ctx, result.ToFD())
}
