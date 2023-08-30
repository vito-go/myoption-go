package web

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/vito-go/mylog"
	"myoption/internal/repo"
	"myoption/types/fd"
	"net/http"
	"time"
)

type SseLastPrice struct {
	RepoClient *repo.Client
}

func (sse *SseLastPrice) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	flusher, ok := rw.(http.Flusher)
	if !ok {
		http.Error(rw, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	dataType := req.URL.Query().Get("dataType")
	switch dataType {
	case "":
		dataType = "json"
	case "json", "text":
	default:
		http.Error(rw, "dataType unsupported!", http.StatusBadRequest)
	}
	rw.Header().Set("Content-Type", "text/event-stream;charset=utf-8")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.WriteHeader(200)
	_, err := rw.Write([]byte(sse.getSSEDataByDataType(0, dataType)))
	if err != nil {
		mylog.Ctx(context.Background()).Error(err)
		return
	}
	mylog.Ctx(context.Background()).Info(0, "<-------- SSE 推送 ----------->", req.RemoteAddr)
	flusher.Flush()

	for i := 1; true; i++ {
		delay := time.Millisecond * 1000
		sleep, status := fd.GetMarketStatus()
		if sleep > 0 {
			delay = sleep
		}
		select {
		case <-req.Context().Done():
			mylog.Ctx(context.Background()).Warn("%s: req done...", req.RemoteAddr)
			return
		default:
			if status != fd.MarketStatusNormal {
				time.Sleep(time.Millisecond * 1000)
				continue
			}
			// 返回数据包含id、event(非必须)、data，结尾必须使用\n\n
			data := []byte(sse.getSSEDataByDataType(i, dataType))
			if _, err = rw.Write(data); err != nil {
				mylog.Ctx(context.Background()).Error(err)
				return
			}
			mylog.Ctx(context.Background()).Infof("%d <-------- SSE 推送 %s-----------> %s", i, dataType, req.RemoteAddr)
			flusher.Flush()
			time.Sleep(delay)
		}

	}
}

func (sse *SseLastPrice) getSSEDataByDataType(id int, dataType string) string {
	lastTodayPrices := sse.RepoClient.StockData.LastTodayPrices()
	var dataMap = map[string]interface{}{
		"items":      lastTodayPrices,
		"updateTime": time.Now().UnixMilli(),
	}
	b, _ := json.Marshal(dataMap)
	switch dataType {
	case "json":
		return string(b) + "\n\n"
	case "", "text":
		return fmt.Sprintf("id: %d\nevent: %s\ndata: %s\n\n", id, "", string(b))
		return fmt.Sprintf("id: %d\nevent: %s\nretry: %d\ndata: %s\n\n", id, "", 3000, string(b))
	default:
		return ""
	}
}
