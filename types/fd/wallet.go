package fd

type WalletInfo struct {
	Balance      int64 `json:"balance"`
	FrozenAmount int64 `json:"frozenAmount"`
	TotalAmount  int64 `json:"totalAmount"`
}

type WalletDetail struct {
	TransId       string `json:"transId"`
	TransType     string `json:"transType"`
	UserId        string `json:"userId"`
	Amount        string `json:"amount"`
	Status        string `json:"status"`
	Remark        string `json:"remark"`
	SourceKind    string `json:"sourceKind"`
	SourceTransId string `json:"sourceTransId"`
	FromAccount   string `json:"fromAccount"`
	ToAccount     string `json:"toAccount"`
	Balance       string `json:"balance"`
	CreateTime    string `json:"createTime"`
	UpdateTime    string `json:"updateTime"`
}
