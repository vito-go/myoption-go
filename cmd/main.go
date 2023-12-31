package main

import (
	"bytes"
	"context"
	"crypto/sha1"
	"flag"
	"fmt"
	"io"
	"log"
	"myoption"
	"myoption/internal/httpsrv"
	"net"
	"net/http"
	_ "net/http/pprof" // 性能、goroutine监控服务
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	_ "github.com/vito-go/daemon"
	"github.com/vito-go/mylog"

	"myoption/conf"
)

const (
	myOptionAPI = "myoption"
)

var mid string // go build -X "main.mid=`cat /etc/machine-id`"
var git string // go build -X "main.git=`git rev-parse HEAD`"

var noCheck string

func init() {
	// 可以在这里做授权校验，也可以在main函数中做授权校验
	return
	func() {
		if noCheck == "true" {
			return
		}
		b, err := os.ReadFile("/etc/machine-id")
		if err != nil {
			fmt.Println("program auth check error. please contact liushihao888@gmail.com")
			os.Exit(1)
		}
		sum := fmt.Sprintf("%x", sha1.Sum(bytes.TrimSpace(b)))
		checkSum := fmt.Sprintf("%x", sha1.Sum([]byte(mid)))
		if sum != checkSum {
			fmt.Println("program auth failed. please contact liushihao888@gmail.com")
			os.Exit(1)
		}
	}()
}

func main() {
	ctx := context.WithValue(context.Background(), "tid", time.Now().UnixNano())
	envPath := flag.String("env", "configs/myoption/test.yaml", "specify the configuration")
	out := flag.Bool("out", true, "only print in os.StdOut, usually for the local running")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "version: %s\nIf have any question, feel free to contact me. email: liushihao888@gmail.com\nUsage of %s:\n", git, os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	// //////////////////
	cfg, err := conf.NewCfg(conf.Env(*envPath))
	if err != nil {
		// 提交代码前需使用make acp 经过配置文件检查
		panic(err)
	}
	chanExit := make(chan struct{}, 1)
	app, err := myoption.NewAPP(cfg)
	if err != nil {
		panic(err)
	}
	pporfPort := app.Cfg.PprofPort

	initLog(cfg.AppName, cfg.Environment, cfg.LogDir, *out)
	mylog.Ctx(ctx).Infof("git version: %+v", git)
	//bash: git rev-parse HEAD

	mylog.Ctx(ctx).Infof("ready to start http server, HTTPServer configs,  %s,  %s ", cfg.AppName, app.Cfg.HTTPServer)
	var router httpsrv.Router

	switch cfg.AppName {
	case myOptionAPI:
		router = httpsrv.NewMyOption(app.HTTPSrv)
	default:
		panic(fmt.Sprintf("unknown app name: %s", cfg.AppName))
	}
	go func() {
		if err := app.HTTPSrv.Start(router); err != nil {
			mylog.Ctx(ctx).Errorf(err.Error())
			chanExit <- struct{}{}
		}
	}()
	mylog.Ctx(ctx).Info("envPath:", *envPath)
	mylog.Ctx(ctx).WithField("configure", cfg).Info()
	mylog.Ctx(ctx).Info("out:", *out)
	mylog.Ctx(ctx).Info("pid:", os.Getpid())
	readyToExit := make(chan struct{})
	if pporfPort > 0 {
		go goStartPProf(pporfPort, readyToExit)
	}
	safeExit(ctx, chanExit, app.HTTPSrv, readyToExit)
}
func goStartPProf(pporfPort uint16, readyToExit <-chan struct{}) {
	ctx := context.Background()
	address := fmt.Sprintf(":%d", pporfPort)
	mylog.Ctx(ctx).Infof("正在启动pprof性能监控服务,addr: [%s]", address)
	ln, err := net.Listen("tcp", address)
	if err != nil {
		mylog.Ctx(ctx).Errorf("pprof服务启动失败！[%s] err: %s", address, err.Error())
		return
	}
	srv := &http.Server{Addr: address, Handler: nil}
	go func() {
		if err := srv.Serve(ln); err != nil {
			mylog.Ctx(ctx).Warnf("pprof服务结束！[%s] err: %s", address, err.Error())
			// it should not be fatal when pprof service is down
			// chanExit <- struct{}{}
		}
	}()
	select {
	case <-readyToExit:
		_ = ln.Close()
	}
}

// safeExit exit elegantly. exit current program will not affect next program start, but will process unfinished connection in a set time.
func safeExit(ctx context.Context, chanExit chan struct{}, srv *httpsrv.Server, readyToExit chan struct{}) {
	c := make(chan os.Signal, 1)
	// If no signals are provided, all incoming signals will be relayed to c.
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT) // 监听键盘终止，以及 kill-15 的信号。注意无法获取kill -9的信号
	select {
	case <-chanExit:
		os.Exit(1)
	case sig := <-c:
		mylog.Ctx(ctx).Warnf("收到进程退出信号: %s", sig.String())
	}
	close(readyToExit)
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	gracefulStopChan := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		srv.Stop(ctx)
	}()
	go func() {
		wg.Wait()
		gracefulStopChan <- struct{}{}
	}()
	select {
	case <-ctx.Done():
		log.Println("exit with timeout")
		os.Exit(1)
	case <-gracefulStopChan:
		log.Println("exit gracefully!")
		os.Exit(0)
	}
}

func initLog(appName string, environment string, logDir string, out bool) {
	if out {
		mylog.Init(os.Stdout, os.Stderr, os.Stderr, "tid")
		return
	}
	if err := os.MkdirAll(filepath.Dir(logDir), 0755); err != nil {
		panic(err)
	}

	infoPath := filepath.Join(logDir, fmt.Sprintf("%s-%s.log", appName, environment))
	infoErrPath := filepath.Join(logDir, fmt.Sprintf("%s-%s.error.log", appName, environment))
	fInfo, err := os.Create(infoPath)
	if err != nil {
		panic(err)
	}
	fErr, err := os.Create(infoErrPath)
	if err != nil {
		panic(err)
	}
	mylog.Init(fInfo, io.MultiWriter(fInfo, fErr), io.MultiWriter(fInfo, fErr), "tid")

}
