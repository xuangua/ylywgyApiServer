package model

import (
    // "strings"
    "time"
)

// Orders
type OrderPayTable struct {
    Id                  int             `json:"id" bson:"id"`

    CreatedDate         time.Time       `json:"created_date" bson:"created_date"`
    UpdatedDate         time.Time       `json:"updated_date" bson:"updated_date"`
    DeletedDate         *time.Time      `json:"deleted_date" bson:"deleted_date"`

    OutTradeNo          string          `json:"out_trade_no" bson:"out_trade_no"`
    OrderId             int             `json:"order_id" bson:"order_id"` // id of OrderTable
    ShopId              int             `json:"shop_id" bson:"shop_id"`
    ShopName            string          `json:"shop_name" bson:"shop_name"`
    StatusDescription   string          `json:"status_description" bson:"status_description"`
    StatusCode          int             `json:"status_code" bson:"status_code"` // 1:"提交支付"; 2:"支付成功"; 3:"支付失败"
    TotalFee           float64          `json:"total_fee" bson:"total_fee"`
    UserId              string          `json:"user_id" bson:"user_id"`

    ReturnCode         string          `json:"return_code" bson:"return_code"`
    ReturnMsg          string          `json:"return_msg" bson:"return_msg"`
    ResultCode         string          `json:"result_code" bson:"result_code"`
    TradeType          string          `json:"trade_type" bson:"trade_type"`
    BankType           string          `json:"bank_type" bson:"bank_type"`
    TransactionId      string          `json:"transaction_id" bson:"transaction_id"`
    TimeEnd            string          `json:"time_end" bson:"time_end"`
}

const (
    // 1:"提交支付";
    PAY_STATUS_COMMITED = 1

    // 2:"支付成功";
    PAY_STATUS_SUCCESS = 2

    // 3:"支付失败";
    PAY_STATUS_FAIL = 3
)


const (
    // 1:"提交支付";
    PAY_STATUS_COMMITED_STR = "提交支付"

    // 2:"支付成功";
    PAY_STATUS_SUCCESS_STR = "支付成功"

    // 3:"支付失败";
    PAY_STATUS_FAIL_STR = "支付失败"
)

func NewOrderPayTable() OrderPayTable {
    return OrderPayTable{
        CreatedDate: time.Now(),
        StatusDescription: PAY_STATUS_COMMITED_STR,
        StatusCode: PAY_STATUS_COMMITED,
    }
}
