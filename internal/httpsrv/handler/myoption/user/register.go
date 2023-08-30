package user

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/vito-go/mylog"
	"myoption/internal/dao/model"
	"myoption/internal/dao/mtype"
	"myoption/internal/httpsrv/handler"
	"myoption/internal/httpsrv/handler/httperr"
	"myoption/types"
	"myoption/types/fd"
	"net/http"
	"time"
	"unicode"

	"myoption/internal/repo"
	"myoption/pkg/util/slice"

	"myoption/pkg/resp"
)

type Register struct {
	RepoClient *repo.Client
}

func (h *Register) GetParam(r *http.Request) (*handler.ReqParam, error) {
	return handler.GetParamBodyEnc(r)
}

type registerParam struct {
	UserId   string `json:"userId,omitempty"`
	Password string `json:"password,omitempty"`
	Salt     string `json:"salt,omitempty"` // 16个字节的随机言
	// 新增字段
	X25519PriEncKey  string `json:"x25519PriEncKey"`
	Ed25519PriEncKey string `json:"ed25519PriEncKey"`
}
type loginData struct {
	UserInfo  *fd.UserInfo `json:"userInfo"`
	LoinToken string       `json:"loinToken"` // 同步数据的时候不用返回
	Balance   int64        `json:"balance"`
}

func (h *Register) GetRespBody(ctx context.Context, p *handler.ReqParam) *resp.HTTPBody {
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
	_, err = h.RepoClient.UserCli.GetUserKeyByUserId(ctx, userId)
	if err == nil {
		return resp.Err(ctx, "用户已存在")
	}
	loinToken := types.NewLoginInfo(userId, p.DeviceId, p.UA)
	userInfo := &model.UserInfo{
		ID:            0,
		UserId:        param.UserId,
		Nick:          param.UserId,
		Status:        mtype.UserStatusNormal,
		X25519PubKey:  base64.StdEncoding.EncodeToString(p.ClientPubKey.Bytes()),
		Ed25519PubKey: base64.StdEncoding.EncodeToString(p.SignPubKeyBytes),
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
	}
	userKey := &model.UserKey{
		ID:               0,
		UserId:           param.UserId,
		Password:         param.Password,
		Salt:             param.Salt,
		X25519PriEncKey:  param.X25519PriEncKey,
		Ed25519PriEncKey: param.Ed25519PriEncKey,
		CreateTime:       time.Now(),
		UpdateTime:       time.Now(),
	}
	info, err := h.RepoClient.UserCli.CreateUser(ctx, userInfo, userKey)
	if err != nil {
		mylog.Ctx(ctx).Error(err.Error())
		return resp.Err(ctx, httperr.ErrData.Error())
	}
	respData := loginData{
		UserInfo:  info.ToFD(),
		LoinToken: loinToken.LoginToken,
		Balance:   10000, // todo 仅供测试用
	}
	return resp.DataOK(ctx, respData)
}

var preUserIds = []string{
	"system", "myoption", "google", "apple", "microsoft", "twitter", "facebook", "android",
}

func (h *Register) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler.ServeHTTP(w, r, h)
}

// checkPassword 检验是否是手机号，并返回手机号数字
// 6-20位 只能是字母数字下划线 至少一位数字和字母
func checkPassword(pwd string) error {
	if len(pwd) == 0 {
		return errors.New("密码不能为为空")
	}
	if len(pwd) < 6 {
		return errors.New("密码长度不能小于6位")
	}
	if len(pwd) > 20 {
		return errors.New("密码长度不能大于20位")
	}

	if len(pwd) != len([]rune(pwd)) {
		return errors.New("密码包含非法字符")
	}
	return nil
}

// checkUserID 用户id： 字母开头-数字-下划线组成。 最小6位，最大18位长度, 下划线不能连续
func checkUserId(userID string) error {
	if slice.IsInSlice(preUserIds, userID) {
		return errors.New("不能使用系统预置用户名")
	}
	if len(userID) == 0 {
		return errors.New("用户名不能为为空")
	}
	if len(userID) < 6 {
		return errors.New("用户名长度不能小于6位")
	}
	if len(userID) > 20 {
		return errors.New("用户名长度不能大于20位")
	}
	if !unicode.IsLetter(rune(userID[0])) {
		return errors.New("用户名必须以小写字母开头")
	}
	if !unicode.IsLower(rune(userID[0])) {
		return errors.New("用户名必须以小字母开头")
	}
	for n, s := range userID {
		if !unicode.IsDigit(s) && !unicode.IsLetter(s) && s != '_' && !unicode.IsLower(s) {
			return errors.New("用户名只能为小写字母、数字或下划线")
		}
		// 下划线不能连续
		if n > 0 && s == '_' && userID[n-1] == '_' {
			return errors.New("用户名不能包含连续的下划线")
		}
	}
	return nil
}
