package httpsrv

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/vito-go/mylog"

	"myoption/pkg/util/slice"
)

const httpMethodAny = "ANY"
const get = http.MethodGet
const post = http.MethodPost

type Router interface {
	Route()
}

func handle(mux *http.ServeMux, path string, h http.Handler) {}
func route(mux *http.ServeMux, method string, path string, h http.Handler) {
	if len(method) == 0 {
		panic(fmt.Sprintf("path: %s. no methods", path))
	}
	if len(path) == 0 {
		panic(fmt.Sprintf("path: %s. no methods", path))
	}
	mylog.Ctx(context.Background()).WithFields("method", method, "path", path).Info("ServeMux register router")
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		//defer mylog.Ctx(r.Context()).WithFields("remoteAddr", r.RemoteAddr, "method", r.Method, "path", path).Info("====")
		// app端暂不需考虑支持跨域
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Max-Age", strconv.FormatInt(int64(time.Second*60*60*24*3), 10))
			return
		}
		if r.Method != method {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		h.ServeHTTP(w, r)
	})

}
func routeWithMethods(mux *http.ServeMux, path string, h http.Handler, method ...string) {
	if len(method) == 0 {
		panic(fmt.Sprintf("path: %s. no methods", path))
	}
	if len(path) == 0 {
		panic(fmt.Sprintf("path: %s. no methods", path))
	}
	mylog.Ctx(context.Background()).WithFields("method", method, "path", path).Info("ServeMux register router")
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		//defer mylog.Ctx(r.Context()).WithFields("remoteAddr",r.RemoteAddr,"method",r.Method,"path",path).Info("====")
		// app端暂不需考虑支持跨域
		//w.Header1().Set("Access-Control-Allow-Origin", "*")
		//w.Header1().Set("Access-Control-Allow-Headers", "*")
		//w.Header1().Set("Access-Control-Allow-Methods", "*")
		//if r.Method == http.MethodOptions {
		//	w.Header1().Set("Access-Control-Max-Age", strconv.FormatInt(int64(time.Second*60*60*24*3), 10))
		//	return
		//}
		if !slice.IsInSlice(method, r.Method) {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		h.ServeHTTP(w, r)
	})

}

func routeWithCorsWithLogin(srv *Server, method string, path string, h http.Handler) {
	mux := srv.serverMux
	if len(method) == 0 {
		panic(fmt.Sprintf("path: %s. no methods", path))
	}
	if len(path) == 0 {
		panic(fmt.Sprintf("path: %s. no methods", path))
	}
	// access.log
	//mylog.Ctx(context.Background()).WithFields("method", method, "path", path).Info("ServeMux register router")
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		//defer mylog.Ctx(r.Context()).WithFields("remoteAddr", r.RemoteAddr, "method", r.Method, "path", path).Info("====")
		//------ 必须写在外面 否则本地无法调用，但是服务器可以的
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		//------ 必须写在外面
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Max-Age", strconv.FormatInt(int64(time.Second*60*60*24*3), 10))
			return
		}
		if r.Method != method {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		if r.Header.Get("test") != "test" {

		}
		if r.Method == http.MethodPost {
			// 关闭修改数据
			//w.Header1().Set("Access-Control-Max-Age", strconv.FormatInt(int64(time.Second*60*60*24*3), 10))
			//w.WriteHeader(503)
			//return
		}
		h.ServeHTTP(w, r)
	})

}
func routeWithCorsNoLogin(srv *Server, method string, path string, h http.Handler) {
	mux := srv.serverMux
	if len(method) == 0 {
		panic(fmt.Sprintf("path: %s. no methods", path))
	}
	if len(path) == 0 {
		panic(fmt.Sprintf("path: %s. no methods", path))
	}
	mylog.Ctx(context.Background()).WithFields("method", method, "path", path).Info("ServeMux register router")
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		//defer mylog.Ctx(r.Context()).WithFields("remoteAddr", r.RemoteAddr, "method", r.Method, "path", path).Info("====")
		//------ 必须写在外面 否则本地无法调用，但是服务器可以的
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Cache-Control", "no-cache")
		//------ 必须写在外面
		if r.Method == http.MethodOptions {
			//w.Header1().Set("Access-Control-Allow-Origin", "*")
			//w.Header1().Set("Access-Control-Allow-Headers", "*")
			//w.Header1().Set("Access-Control-Allow-Methods", "*")
			w.Header().Set("Access-Control-Max-Age", strconv.FormatInt(int64(time.Second*60*60*24*3), 10))
			return
		}
		if r.Method != method {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		h.ServeHTTP(w, r)
	})

}
