package user

type FieldName struct {
	Field string `json:"field"`
	Name  string `json:"name"`
}

var fieldNamesOrderList = []FieldName{
	{Field: "orderTime", Name: "订单时间"},
	{Field: "orderStatus", Name: "订单状态"},
	{Field: "sessionTimMin", Name: "到期时间"},
	{Field: "symbolCodeName", Name: "名称代码"},
	{Field: "strikePrice", Name: "行权价"},
	{Field: "option", Name: "看涨/看跌"},
	{Field: "betMoney", Name: "下注金额"},
	{Field: "session", Name: "场次"},
	{Field: "settlePrice", Name: "结算价"},
	{Field: "settleResult", Name: "交割结果"},
	{Field: "profitLoss", Name: "损益"},
}

var fieldWalletDetails = []FieldName{
	{Field: "transId", Name: "交易ID"},
	{Field: "transType", Name: "交易类型"},
	//{Field: "userId", Name: "用户ID"},
	{Field: "amount", Name: "金额"},
	{Field: "status", Name: "状态"},
	{Field: "remark", Name: "说明"},
	//{Field: "sourceKind", Name: "来源类型"},
	{Field: "sourceTransId", Name: "来源交易ID"},
	//{Field: "fromAccount", Name: "转出账户"},
	//{Field: "toAccount", Name: "转入账户"},
	{Field: "balance", Name: "余额"},
	{Field: "createTime", Name: "创建时间"},
	{Field: "updateTime", Name: "更新时间"},
}
