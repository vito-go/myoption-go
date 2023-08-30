package user

import (
	"context"
	"encoding/json"
	"github.com/vito-go/mylog"
	"myoption/iface/myerr"
	"myoption/internal/httpsrv/handler"
	"myoption/internal/httpsrv/handler/httperr"
	"myoption/internal/repo"
	"myoption/pkg/resp"
	"myoption/types"
	"net/http"
)

type LogIn struct {
	RepoClient *repo.Client
}

func (h *LogIn) GetParam(r *http.Request) (*handler.ReqParam, error) {
	return handler.GetParamBodyEnc(r)
}

type logInData struct {
	LoinToken        string `json:"loinToken"` // 同步数据的时候不用返回
	Balance          int64  `json:"balance"`
	X25519PubKey     string `json:"x25519PubKey"`  //
	Ed25519PubKey    string `json:"ed25519PubKey"` // Ed25519公钥
	Salt             string `json:"salt"`
	X25519PriEncKey  string `json:"x25519PriEncKey"`  // X25519私钥加密密钥
	Ed25519PriEncKey string `json:"ed25519PriEncKey"` // Ed25519私钥加密密钥
}

func (h *LogIn) GetRespBody(ctx context.Context, p *handler.ReqParam) *resp.HTTPBody {
	var param registerParam
	err := json.Unmarshal(p.Body, &param)
	if err != nil {
		return resp.Err(ctx, httperr.ErrParam.Error())
	}
	if param.Salt == "" {
		return resp.Err(ctx, "salt is empty")
	}
	if err = checkUserId(param.UserId); err != nil {
		return resp.Err(ctx, err.Error())
	}
	userId := param.UserId
	userInfo, err := h.RepoClient.UserCli.GetUserInfoByUserId(ctx, userId)
	if err != nil {
		if err == myerr.DataNotFound {
			return resp.Err(ctx, "用户不存在")
		}
		mylog.Ctx(ctx).WithField("userId", userId).Error(err.Error())
		return resp.Err(ctx, httperr.ErrData.Error())
	}

	userKey, err := h.RepoClient.UserCli.GetUserKeyByUserId(ctx, userId)
	if err != nil {
		if err == myerr.DataNotFound {
			return resp.Err(ctx, "用户不存在")
		}
		mylog.Ctx(ctx).WithField("userId", userId).Error(err.Error())
		return resp.Err(ctx, httperr.ErrData.Error())
	}
	if param.Password != userKey.Password {
		// 记录登录失败
		return resp.Err(ctx, "密码错误")
	}
	walletInfo, err := h.RepoClient.Wallet.GetWalletInfoByUserId(ctx, userId)
	if err != nil {
		mylog.Ctx(ctx).WithField("userId", userId).Error(err.Error())
		return resp.Err(ctx, httperr.ErrData.Error())
	}

	loinToken := types.NewLoginInfo(userId, p.DeviceId, p.UA)
	respData := logInData{
		LoinToken:        loinToken.LoginToken,
		Balance:          walletInfo.Balance,
		X25519PubKey:     userInfo.X25519PubKey,
		Ed25519PubKey:    userInfo.Ed25519PubKey,
		Salt:             userKey.Salt,
		X25519PriEncKey:  userKey.X25519PriEncKey,
		Ed25519PriEncKey: userKey.Ed25519PriEncKey,
	}
	return resp.DataOK(ctx, respData)
}

func (h *LogIn) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler.ServeHTTP(w, r, h)
}
