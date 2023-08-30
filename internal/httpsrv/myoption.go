package httpsrv

import (
	"context"
	"github.com/vito-go/mylog"
	"myoption/internal/httpsrv/handler/myoption/user"
	"myoption/internal/httpsrv/handler/myoption/web"
	"net/http"

	"myoption/internal/repo"
)

type myOption struct {
	server *Server
}

func NewMyOption(server *Server) *myOption {
	return &myOption{server: server}
}

func (s *myOption) Route() {
	mux := s.server.serverMux
	repoClient := s.server.repoCli
	s.routeUser(mux, repoClient)
	s.routeOrder(mux, repoClient)
	s.routerWeb(mux, repoClient)
	// post

}
func (s *myOption) routeUser(mux *http.ServeMux, repoClient *repo.Client) {
	route(mux, post, "/myoption/api/v1/user/register", &user.Register{RepoClient: repoClient})
	route(mux, post, "/myoption/api/v1/user/logIn", &user.LogIn{RepoClient: repoClient})
}
func (s *myOption) routeOrder(mux *http.ServeMux, repoClient *repo.Client) {
	// ----------------- user -----------------
	route(mux, get, "/myoption/api/v1/user/order/orderList", &user.OrderList{RepoClient: repoClient})
	route(mux, post, "/myoption/api/v1/user/order/submit", &user.SubmitOrder{RepoClient: repoClient})
	route(mux, get, "/myoption/api/v1/user/wallet/walletDetails", &user.WalletDetails{RepoClient: repoClient})
	route(mux, get, "/myoption/api/v1/user/wallet/myBalance", &user.MyBalance{RepoClient: repoClient})
}
func (s *myOption) routerWeb(mux *http.ServeMux, repoClient *repo.Client) {
	routeWithCorsNoLogin(s.server, get, "/myoption/ws/v1/price/productList", &web.ProductList{RepoClient: repoClient})
	routeWithCorsNoLogin(s.server, get, "/myoption/sse/v1/price/last", &web.SseLastPrice{RepoClient: repoClient})
	routeWithCorsNoLogin(s.server, get, "/myoption/web/v1/user/generalConfig", &user.GeneralConfig{RepoClient: repoClient})
	routeWithCorsNoLogin(s.server, get, "/myoption/web/v1/user/lineChart", &web.LineChart{RepoClient: repoClient})
	if s.server.isOnline {
		// todo: you can make link to your static file which is in your project such as flutter web
		mux.Handle("/web/", http.FileServer(http.Dir("./www/online/")))
		mylog.Ctx(context.Background()).Info("online mode")
	} else {
		mylog.Ctx(context.Background()).Info("test mode")
		// todo: test mode
		mux.Handle("/web/", http.FileServer(http.Dir("./www/online/")))
	}
}
