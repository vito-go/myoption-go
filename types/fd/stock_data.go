package fd

import "myoption/internal/dao/mtype"

type SymbolPrice struct {
	SymbolCode   string        `json:"symbolCode"`
	SymbolName   string        `json:"symbolName"`
	Exist        bool          `json:"exist"`
	Price        float64       `json:"price"`
	Day          int64         `json:"day"`
	TimeMin      mtype.TimeMin `json:"timeMin"`
	MarketStatus MarketStatus  `json:"marketStatus"`
}
