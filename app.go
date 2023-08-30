package myoption

import (
	"myoption/conf"
	"myoption/internal/connector"
	"myoption/internal/httpsrv"
)

type APP struct {
	Cfg     *conf.Cfg
	HTTPSrv *httpsrv.Server
}
type APPName string

func NewAPP(cfg *conf.Cfg) (*APP, error) {

	c, err := connector.New(cfg)
	if err != nil {
		return nil, err
	}
	var isOnline = true
	if cfg.Environment == "test" {
		isOnline = false
	} else if cfg.Environment == "online" {
		isOnline = true
	} else {
		panic("unknown environment")
	}

	httpSrv := httpsrv.NewServer(isOnline, cfg.HTTPServer, cfg.ConstantKey, c)
	return &APP{
		Cfg:     cfg,
		HTTPSrv: httpSrv,
	}, nil
}
