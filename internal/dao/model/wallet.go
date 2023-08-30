package model

import (
	"myoption/internal/dao/mtype"
	"myoption/types/fd"
	"strconv"
	"time"
)

type WalletInfo struct {
	ID           int64     `json:"id,omitempty"`
	UserId       string    `json:"user_id"`
	Balance      int64     `json:"balance"`
	FrozenAmount int64     `json:"frozen_amount"`
	TotalAmount  int64     `json:"total_amount"`
	CreateTime   time.Time `json:"create_time"`
	UpdateTime   time.Time `json:"update_time"` // 访问的时间 一般只记录第一次
}

func (w *WalletInfo) ToFD() *fd.WalletInfo {
	return &fd.WalletInfo{
		Balance:      w.Balance,
		FrozenAmount: w.FrozenAmount,
		TotalAmount:  w.TotalAmount,
	}
}

// id：钱包 ID，自动递增生成。
//user_id：与用户表关联的外键，用于标识哪个用户拥有该钱包。
//total_amount：钱包中的总金额。
//frozen_amount：冻结的金额，无法使用。
//balance：可用余额，等于总金额减去冻结金额。
//created_at：创建时间戳，在插入数据时自动设置为当前时间。
//updated_at：更新时间戳，在修改数据时自动更新为当前时间。

type WalletDetail struct {
	TransId       string            `gorm:"column:trans_id" json:"trans_id"`
	TransType     mtype.TransType   `gorm:"column:trans_type" json:"trans_type"`
	UserId        string            `gorm:"column:user_id" json:"user_id"`
	Amount        int64             `gorm:"column:amount" json:"amount"`
	Status        mtype.TransStatus `gorm:"column:status" json:"status"`
	Remark        string            `gorm:"column:remark" json:"remark"`
	SourceKind    mtype.SourceKind  `gorm:"column:source_kind" json:"source_kind"`
	SourceTransId string            `gorm:"column:source_trans_id" json:"source_trans_id"`
	FromAccount   string            `gorm:"column:from_account" json:"from_account"`
	ToAccount     string            `gorm:"column:to_account" json:"to_account"`
	Balance       int64             `gorm:"column:balance" json:"balance"`
	CreateTime    time.Time         `json:"create_time"`
	UpdateTime    time.Time         `json:"update_time"` // 访问的时间 一般只记录第一次
}

func (*WalletDetail) TableName() string {
	return "wallet_detail"
}

func (w *WalletDetail) ToFd() *fd.WalletDetail {
	return &fd.WalletDetail{
		TransId:       w.TransId,
		TransType:     w.TransType.Name(),
		UserId:        w.UserId,
		Amount:        strconv.FormatInt(w.Amount, 10),
		Status:        w.Status.Name(),
		Remark:        w.Remark,
		SourceKind:    w.SourceKind.Name(),
		SourceTransId: w.SourceTransId,
		FromAccount:   w.FromAccount,
		ToAccount:     w.ToAccount,
		Balance:       strconv.FormatInt(w.Balance, 10),
		CreateTime:    w.CreateTime.Format("2006-01-02 15:04:05"),
		UpdateTime:    w.UpdateTime.Format("2006-01-02 15:04:05"),
	}
}

//设计一个高级的postgresql钱包明细表，要求支持转帐，发送红包，充值，消费，记录转帐目标用户，充值来源，消费去向，并被给sql语句

//表结构：
//
//CREATE TABLE wallet_detail (
//   id SERIAL PRIMARY KEY,
//   user_id INT NOT NULL,
//   type VARCHAR(32) NOT NULL,   -- 支付类型：充值、消费、转账、红包
//   detail_type VARCHAR(32) NOT NULL,  -- 具体类型：充值方式、消费去向、交易状态、红包改变类型等
//   amount NUMERIC(18,2) NOT NULL DEFAULT 0,  -- 交易金额
//   transaction_id VARCHAR(64),    -- 支付平台交易 ID（如：微信支付订单号，支付宝订单号等）
//   transaction_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,   -- 交易时间
//   remark TEXT,        -- 交易备注
//   source_user_id INT,      -- 转出用户 ID（在转账业务中使用）
//   target_user_id INT,      -- 转入用户 ID（在转账、红包业务中使用）
//   red_packets_type VARCHAR(16),    -- 红包变更类型（在红包业务中使用）
//   transaction_status VARCHAR(16),  -- 交易状态，如“成功”、“失败”、“待处理”等
//   create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,   -- 数据创建时间
//   update_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP   -- 数据更新时间
//);
//
//
//表结构说明：
//
//- id：明细记录唯一标识 ID。
//- user_id：钱包所属用户 ID。
//- type：交易类型。可以是“充值”、“消费”、“转账”、“红包”等，需要根据不同业务需求设置对应的类型。
//- detail_type：交易详细类型，可以用来描述充值方式、消费去向、交易状态等具体的信息。
//- amount：交易金额（使用 NUMERIC 数据类型存储）。
//- transaction_id：支付平台交易 ID，例如微信支付订单号、支付宝订单号等。
//- transaction_time：交易时间。
//- remark：交易备注，可用于记录交易详情。
//- source_user_id：转出用户 ID，在转账业务中使用。
//- target_user_id：转入用户 ID，可以在转账、红包等业务中使用。
//- red_packets_type：红包变更类型，在发送红包业务中使用，例如“发出”、“领取”、“退回”等类型。
//- transaction_status：交易状态，例如“成功”、“失败”、“待处理”等。
//- create_time：数据创建时间。用于记录数据的创建时间，方便进行数据追溯。
//- update_time：数据更新时间。用于记录数据的最后更新时间。
//
//这个表结构可以满足转账、发送红包、充值和消费等业务需求。具体业务在插入数据时通过type和detail_type来标记交易类型和详细类型。
//
//比如说，对于充值业务，你可以插入以下数据：
//
//INSERT INTO wallet_detail(user_id, type, detail_type, amount, transaction_id, remark)
//VALUES(123, '充值', '支付宝', 100.00, '123456789', '支付宝充值 100 元');
//
//
//对于转账业务，你可以插入以下数据：
//
//`sql
//INSERT INTO wallet_detail(user_id, type, detail_type, amount, source_user_id, target_user_id, remark)
//VALUES
//
//The number of tokens used this time: 4097

//Pending - PEND （待处理）：交易已被提交但尚未完成处理。
//Processing - PROC （处理中）：交易正在处理过程中，可能需要一些时间才能完成。
//Completed - COMP （已完成）：交易已经成功完成。
//Failed - FAIL （失败）：交易因某种原因未能成功完成。
//Refunded - REFND （已退款）：交易金额已被全额或部分退还给买家。
//Reversed - REVRS （已撤销）：交易金额已被取消或撤销。
//Chargeback - CBK （退单）：买家在未经授权的情况下向其支付的银行或信用卡提供商发起了争议/纠纷，从而导致付款被返还。
//Hold - HOLD （暂停）：交易金额已被暂时保留，直到特定条件得到满足或问题得到解决。
//Authorized - AUTH （已授权）：交易金额已被授权并将在未来的某个时间点进行支付。
//Settled - SETL （已结算）：交易金额已被成功结算并转移到接收方账户中。
