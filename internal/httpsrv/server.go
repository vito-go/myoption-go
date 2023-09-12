package httpsrv

import (
	"context"
	"fmt"
	"github.com/vito-go/mylog"
	"myoption/internal/cache"
	"myoption/internal/dao"
	"myoption/internal/repo"
	"net"
	"net/http"
	"sync"
	"time"

	"myoption/conf"
	"myoption/internal/connector"
)

// Server 启动http服务.
type Server struct {
	HTTPServerConfigs []*HTTPServerConfig
	isOnline          bool
	constantKey       conf.ConstantKey
	serverMux         *http.ServeMux //相当于gin.Engine
	connector         *connector.Connector
	repoCli           *repo.Client
	cache             *cache.Cache
}
type HTTPServerConfig struct {
	Server         *http.Server
	HttpServerConf conf.HttpServerConf
}

func (s *Server) Connector() *connector.Connector {
	return s.connector
}
func (s *Server) IsOnline() bool {
	return s.isOnline
}
func (s *Server) Cache() *cache.Cache {
	return s.cache
}

type TidCxtKey struct {
	name string
}

func NewServer(isOnline bool, cfgs []conf.HttpServerConf, constantKey conf.ConstantKey, c *connector.Connector) *Server {
	serverMux := http.NewServeMux()
	httpServerConfigs := make([]*HTTPServerConfig, 0, len(cfgs))
	for _, cfg := range cfgs {
		srv := &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Port),
			Handler: serverMux,
			//ReadTimeout: time.Millisecond * time.Duration(cfg.ReadTimeout),
			// 千万不要配置WriteTimeout啊，影响大数据的传输，应该log日志的推送！！ 会超时断开
			// WriteTimeout: time.Millisecond * time.Duration(cfg.WriteTimeout),
			ConnContext: func(ctx context.Context, c net.Conn) context.Context {
				// 中间件也可以使用这里的context
				return context.WithValue(ctx, "tid", time.Now().UnixNano())
			},
		}
		httpServerConfigs = append(httpServerConfigs, &HTTPServerConfig{
			Server:         srv,
			HttpServerConf: cfg,
		})
	}
	s := &Server{
		HTTPServerConfigs: httpServerConfigs,
		connector:         c,
		isOnline:          isOnline,
		serverMux:         serverMux,
		constantKey:       constantKey, repoCli: repo.NewClient(c), cache: cache.New(c.RedisCli, dao.NewAllDao(c.GDB))}
	return s
}

func (s *Server) Start(r Router) error {
	// 所有的路由路径都显式写在这里，如ugo超过64行可以封装函数
	r.Route()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, httpServerConfig := range s.HTTPServerConfigs {
		cfg := httpServerConfig.HttpServerConf
		server := httpServerConfig.Server
		if cfg.KeyFile != "" && cfg.CertFile != "" {
			mylog.Ctx(ctx).WithField("cfg", cfg).Infof("----- HTTPS ServerStart: [:%d] -----", cfg.Port)
			go func(cfg conf.HttpServerConf, server *http.Server) {
				if err := server.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile); err != nil {
					mylog.Ctx(ctx).WithField("cfg", cfg).Info(err.Error())
					cancel()
				}
			}(cfg, server)
		} else {
			mylog.Ctx(ctx).WithField("cfg", cfg).Infof("----- HTTP  ServerStart: [:%d] -----", cfg.Port)
			go func(cfg conf.HttpServerConf, server *http.Server) {
				if err := server.ListenAndServe(); err != nil {
					mylog.Ctx(ctx).WithField("cfg", cfg).Info(err.Error())
					cancel()
				}
			}(cfg, server)
		}
	}
	<-ctx.Done()
	return http.ErrServerClosed
}
func (s *Server) Stop(ctx context.Context) {
	var wg sync.WaitGroup
	for _, httpServerConfig := range s.HTTPServerConfigs {
		cfg := httpServerConfig.HttpServerConf
		server := httpServerConfig.Server
		mylog.Ctx(ctx).WithField("cfg", cfg).Info("shutdown httpServer")
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := server.Shutdown(ctx)
			if err != nil {
				mylog.Ctx(ctx).Error(err)
			}
		}()
	}
	wg.Wait()
	s.connector.Close(ctx)

}
