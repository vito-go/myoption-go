package web

import (
	"fmt"
	"log"
	"math"
	"myoption/types/fd"
	"strconv"
	"time"
)

// LineValues .
type LineValues struct {
	Tip     string       `json:"tip"`
	FlSpots [][2]float64 `json:"flSpots"`
}

type Lines struct {
	Title          string            `json:"title"`
	SubTitle       string            `json:"subTitle"`
	XName          string            `json:"xName"`
	YName          string            `json:"yName"`
	DotDataShow    bool              `json:"dotDataShow"`
	XTitleMap      map[string]string `json:"xTitleMap,omitempty"`
	XTitleIndexMap map[string]string `json:"xTitleIndexMap,omitempty"`
	YTitleMap      map[string]string `json:"yTitleMap,omitempty"`
	LineValues     []LineValues      `json:"lineValues"`

	BaselineY float64  `json:"baselineY,omitempty"`
	BaselineX *float64 `json:"baselineX,omitempty"`
	MaxX      *float64 `json:"maxX,omitempty"`
	MaxY      *float64 `json:"maxY,omitempty"`
	MinX      *float64 `json:"minX,omitempty"`
	MinY      *float64 `json:"minY,omitempty"`
}
type Sse struct {
	Code      string      `json:"code"`
	PrevClose float64     `json:"prev_close"`
	Highest   float64     `json:"highest"`
	Lowest    float64     `json:"lowest"`
	Date      int         `json:"date"`
	Time      int         `json:"time"`
	Total     int         `json:"total"`
	Begin     int         `json:"begin"`
	End       int         `json:"end"`
	Line      [][]float64 `json:"line"` // time,price,volume,avg_price,amount,highest,lowest
}

//    [
//      150000, // 时间
//      3226.8912, // 价格
//      3452829, 	// 成交量
//      3243.1779,// 均价
//      3934035868, //成交额
//      null, 		// 最高价
//      null 		// 最低价
//    ]

func (sse *Sse) buildLines() Lines {
	var flSpots [][2]float64
	xTitleMap := make(map[string]string)
	xTitleIndexMap := make(map[string]string) // 真实值与索引值的映射

	mins := fd.GetTimeMinS()
	prevClose := sse.PrevClose
	var maxGap float64
	var lastPrice float64
	for x, min := range mins {
		t, _ := time.ParseInLocation("20060102 150405", fmt.Sprintf("20060102 %06d", min), time.Local)
		xTitleMap[strconv.FormatInt(int64(x), 10)] = t.Format("15:04")
		xTitleIndexMap[strconv.FormatInt(int64(min), 10)] = strconv.FormatInt(int64(x), 10)
		var exist bool
		for i := 0; i < len(sse.Line); i++ {
			lineData := sse.Line[i]
			if int64(lineData[0]) == int64(min) {
				if abs := math.Abs(prevClose - lineData[1]); abs > maxGap {
					maxGap = abs
				}
				v, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", lineData[1]), 64)
				flSpots = append(flSpots, [2]float64{float64(x), v})
				lastPrice = v
				exist = true
				break
			}
		}
		if !exist {
			flSpots = append(flSpots, [2]float64{float64(x), lastPrice})
			if abs := math.Abs(prevClose - lastPrice); abs > maxGap {
				maxGap = abs
			}
		}

	}
	lv := LineValues{
		Tip:     "",
		FlSpots: flSpots,
	}

	maxAmplitude := math.Abs(sse.Highest/sse.PrevClose - 1)
	if m := math.Abs(sse.Lowest/sse.PrevClose - 1); m > maxAmplitude {
		maxAmplitude = m
	}
	_maxY := toFix2(sse.PrevClose * (1 + maxGap/prevClose + 0.003))
	_minY := toFix2(sse.PrevClose * (1 - (maxGap/prevClose + 0.003)))

	minY := &_minY
	maxY := &_maxY
	const fix = 0.002
	log.Println(fmt.Sprintf("_maxY: %+v, baselineY: %+v, minY:%+v", _maxY, toFix2(sse.PrevClose), minY))
	//yTitleMap[fmt.Sprintf("%+v", _minY)] = fmt.Sprintf("-%.2f%%", (maxAmplitude+fix)*100)
	//yTitleMap[fmt.Sprintf("%+v", _maxY)] = fmt.Sprintf("%.2f%%", (maxAmplitude+fix)*100)
	return Lines{
		//Title:          "上证指数(000001)",
		//SubTitle:       time.Now().Format("2006-01-02"),
		XName:          "时间",
		YName:          "价格",
		DotDataShow:    false,
		XTitleMap:      xTitleMap,
		XTitleIndexMap: xTitleIndexMap,
		//YTitleMap:      yTitleMap,
		LineValues: []LineValues{lv},
		BaselineY:  toFix2(sse.PrevClose),
		BaselineX:  nil,
		MaxX:       nil,
		MaxY:       maxY,
		MinX:       nil,
		MinY:       minY,
	}
}

func toFix(f float64, prec int) float64 {
	v, _ := strconv.ParseFloat(strconv.FormatFloat(f, 'f', prec, 64), 64)
	return v
}
func toFix2(f float64) float64 {
	v, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", f), 64)
	return v
}
