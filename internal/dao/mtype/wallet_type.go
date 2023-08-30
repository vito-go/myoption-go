package mtype

// PaymentKind 实时转帐，普通转帐，红包类型（私人红包，群红包）
type PaymentKind int

const (
	PaymentKindTransfer      = PaymentKind(1)
	PaymentKindPriRedPacket  = PaymentKind(10)
	PaymentKinGroupRedPacket = PaymentKind(20)
)

type TransStatus int

// 以下是一些可能的交易状态及其英文缩写和中文解释：
//
// - 成功（Success，SUC）：交易已经成功完成。
// - 失败（Failure，FAL）：交易失败。
// - 处理中（Processing，PEN）：交易正在处理中，还未完成。
// - 等待中（Waiting，WAI）：交易已被提交，但还未处理。
// - 未处理（Untreated，UNT）：交易已提交，但系统尚未开始处理。
// - 已取消（Canceled，CAN）：交易被取消，但尚未退款。
// - 已退款（Refunded，REF）：交易已被取消并退款。
// - 已完成（Completed，COM）：交易已完成并已被处理。
// - 已关闭（Closed，CLO）：交易已彻底关闭。
//
// 需要根据具体业务需求，添加、删除或修改相关状态的定义及缩写。
//
// The number of tokens used this time: 2289
const (
	TransStatusWaiting = TransStatus(1)   // 等待进一步的操作 比如发起转帐通知
	TransStatusPending = TransStatus(2)   // 处理中
	TransStatusClose   = TransStatus(100) // 服务端主动关闭（比如转帐通知超时）

	TransStatusSuc = TransStatus(200)
)

func (t TransStatus) Name() string {
	// 使用switch语句实现，完整的给出代码如下所示 :

	switch t {
	case TransStatusWaiting:
		return "等待处理"
	case TransStatusPending:
		return "处理中"
	case TransStatusClose:
		return "已关闭"
	case TransStatusSuc:
		return "已完成"
	default:
		return "unknown"
	}
}

// TransType 充值deposit，提现 withdraw，下单， 盈利
type TransType int

func (t TransType) Name() string {

	switch t {
	case TransTypeOrder:
		return "下单"
	case TransTypeWithdraw:
		return "提现"
	case TransTypeDeposit:
		return "充值"
	case TransTypeSettle:
		return "结算"
	default:
		return "unknown"
	}
}

const (
	// 100以下为减
	TransTypeOrder    = TransType(1)
	TransTypeWithdraw = TransType(11)
	// 100以上为加
	TransTypeDeposit = TransType(100)
	TransTypeSettle  = TransType(101)
)

type TransactionType int

const (
	TransactionTypeRecharge = TransactionType(1)
	TransactionTypeConsume  = TransactionType(2)
)

//type TransSourceType int

// const (
//
//	TransSourceNil        = TransSourceType(0) //非充值
//	TransSourceDefault    = TransSourceType(1)
//	TransSourceTypeAliPay = TransSourceType(10)
//	TransSourceTypeWechat = TransSourceType(11)
//	TransSourceTypePaypal = TransSourceType(12)
//	TransSourceTypeBank   = TransSourceType(100)
//
// )
type SourceKind int

const (
	TransSourceNil        = SourceKind(0) //　no
	TransSourceDefault    = SourceKind(1)
	TransSourceTypeAliPay = SourceKind(10)
	TransSourceTypeWechat = SourceKind(11)
	TransSourceTypePaypal = SourceKind(12)
	TransSourceTypeBank   = SourceKind(100)
)

// Name 用switch实现SourceKind的Name方法
func (s SourceKind) Name() string {
	switch s {
	case TransSourceNil:
		return ""
	case TransSourceDefault:
		return "default"
	case TransSourceTypeAliPay:
		return "支付宝"
	case TransSourceTypeWechat:
		return "微信"
	case TransSourceTypePaypal:
		return "Paypal"
	case TransSourceTypeBank:
		return "银行"
	default:
		return ""
	}
}
