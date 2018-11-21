package model

import (
    "strings"
    "time"
)


type BasketAbandonedExtraStruct struct {
    Name                    string              `json:"name" bson:"name"`               // {type: String, default: '餐盒'},
    // Description             string              `json:"description" bson:"description"`

    Price                   float64             `json:"price" bson:"price"`             // {type: Number, default: 0},
    Quantity                int                 `json:"quantity" bson:"quantity"`       // {type: Number, default: 0},
    FoodCategoryId          int                 `json:"category_id" bson:"category_id"`               // {type: Number, default: 0},
}

type BasketExtraStruct struct {
    Name                    string              `json:"name" bson:"name"`               // {type: String, default: '餐盒'},
    Description             string              `json:"description" bson:"description"`

    Price                   float64             `json:"price" bson:"price"`             // {type: Number, default: 0},
    Quantity                int                 `json:"quantity" bson:"quantity"`       // {type: Number, default: 0},
    Type                    int                 `json:"type" bson:"type"`               // {type: Number, default: 0},
}

type BasketDeliverFeeStruct struct {
    Name                    string              `json:"name" bson:"name"`               // {type: String, default: '餐盒'},
    // Description             string              `json:"description" bson:"description"`

    Price                   float64             `json:"price" bson:"price"`             // {type: Number, default: 0},
    Quantity                int                 `json:"quantity" bson:"quantity"`       // {type: Number, default: 0},
    CategoryId          int                 `json:"category_id" bson:"category_id"`               // {type: Number, default: 0},
}

func NewBasketDeliverFee() BasketDeliverFeeStruct {
    return BasketDeliverFeeStruct{
        Name: "配送费",
        Price: 2,
        CategoryId: 2,
        Quantity: 1,
    }
}

type BasketPackingFeeStruct struct {
    Name                    string              `json:"name" bson:"name"`               // {type: String, default: '餐盒'},
    // Description             string              `json:"description" bson:"description"`

    Price                   float64             `json:"price" bson:"price"`             // {type: Number, default: 0},
    Quantity                int                 `json:"quantity" bson:"quantity"`       // {type: Number, default: 0},
    CategoryId              int                 `json:"category_id" bson:"category_id"`               // {type: Number, default: 0},
}

func NewBasketPackingFee() BasketPackingFeeStruct {
    return BasketPackingFeeStruct{
        Name: "餐盒",
        CategoryId: 1,
    }
}

type BasketPindanMapStruct struct {
}


type BasketStruct struct {

    AbandonedExtra  []BasketAbandonedExtraStruct  `json:"abandoned_extra" bson:"abandoned_extra"`
    DeliverFee      BasketDeliverFeeStruct  `json:"deliver_fee" bson:"deliver_fee"`
    Extra           []BasketExtraStruct  `json:"extra" bson:"extra"`
    Groups          []CartGroupStruct     `json:"groups" bson:"groups"`

    PackingFee      BasketPackingFeeStruct     `json:"packing_fee" bson:"packing_fee"`
    PindanMap       []BasketPindanMapStruct     `json:"pindan_map" bson:"pindan_map"`
}

func NewBasket() BasketStruct {
    return BasketStruct{
        DeliverFee: NewBasketDeliverFee(),
        PackingFee: NewBasketPackingFee(),
    }
}

type StatusBarStruct struct {
    Color            string     `json:"color" bson:"color"`
    ImageType            string     `json:"image_type" bson:"image_type"`
    SubTitle            string     `json:"sub_title" bson:"sub_title"`
    Title            string     `json:"title" bson:"title"`
}
func NewStatusBar() StatusBarStruct {
    return StatusBarStruct{
        Color: "f60",
        ImageType: "",
        SubTitle: "15分钟内支付",
        Title: "",
    }
}

type TimelineNodeStruct struct {
    Actions            []string     `json:"actions" bson:"actions"`
    InProcessing                      int                     `json:"in_processing" bson:"in_processing"`
    Description            string     `json:"description" bson:"description"`
    SubDescription            string     `json:"sub_description" bson:"sub_description"`
    Title            string     `json:"title" bson:"title"`
}

func NewTimelineNode() TimelineNodeStruct {
    return TimelineNodeStruct{
        InProcessing: 0,
    }
}


// Orders
type OrderTable struct {
    Id                      int                     `json:"id" bson:"id"`

    CreatedDate             time.Time               `json:"created_date" bson:"created_date"`
    UpdatedDate             time.Time               `json:"updated_date" bson:"updated_date"`
    DeletedDate             *time.Time               `json:"deleted_date" bson:"deleted_date"`

    Basket                 BasketStruct           `json:"basket" bson:"basket"`
    FormattedCreatedAt          string                  `json:"formatted_created_at" bson:"formatted_created_at"`
    OrderTime          int64                     `json:"order_time" bson:"order_time"`
    TimePass          int64                     `json:"time_pass" bson:"time_pass"`

    IsBrand          int                     `json:"is_brand" bson:"is_brand"`
    IsDeletable          int                     `json:"is_deletable" bson:"is_deletable"`
    IsNewPay          int                     `json:"is_new_pay" bson:"is_new_pay"`
    IsPindan          int                     `json:"is_pindan" bson:"is_pindan"`
    OperationConfirm          int                     `json:"operation_confirm" bson:"operation_confirm"`
    OperationPay          int                     `json:"operation_pay" bson:"operation_pay"`
    OperationRate          int                     `json:"operation_rate" bson:"operation_rate"`
    OperationRebuy          int                     `json:"operation_rebuy" bson:"operation_rebuy"`
    OperationUploadPhoto          int                     `json:"operation_upload_photo" bson:"operation_upload_photo"`
    PayRemainSeconds          int                     `json:"pay_remain_seconds" bson:"pay_remain_seconds"`
    RatedPoint          int                     `json:"rated_point" bson:"rated_point"`
    RemindReplyCount          int                     `json:"remind_reply_count" bson:"remind_reply_count"`

    ShopId            int     `json:"shop_id" bson:"shop_id"`
    ShopImageHash            string     `json:"shop_image_hash" bson:"shop_image_hash"`
    ShopImageUrl            string     `json:"shop_image_url" bson:"shop_image_url"`
    ShopName            string     `json:"shop_name" bson:"shop_name"`
    ShopType            int     `json:"shop_type" bson:"shop_type"`

    StatusBar                 StatusBarStruct           `json:"status_bar" bson:"status_bar"`
    StatusCode            int     `json:"status_code" bson:"status_code"`

    TimelineNode                 TimelineNodeStruct           `json:"timeline_node" bson:"timeline_node"`
    
    TopShow          int                     `json:"top_show" bson:"top_show"`
    TotalAmount          float64                     `json:"total_amount" bson:"total_amount"`
    TotalQuantity          int                     `json:"total_quantity" bson:"total_quantity"`

    UniqueId          int                     `json:"unique_id" bson:"unique_id"`
    UserId          string                     `json:"user_id" bson:"user_id"`
    AddressId          int                     `json:"address_id" bson:"address_id"`
}

func NewOrderTable() OrderTable {
    return OrderTable{
        CreatedDate: time.Now(),
        FormattedCreatedAt: time.Now().Format("2006-01-02 15:04:05"),
        OrderTime: time.Now().Unix(),
        TimePass: 900,

        Basket: NewBasket(),
        TimelineNode: NewTimelineNode(),

        IsBrand: 0,
        IsDeletable: 1,
        IsNewPay: 1,
        IsPindan: 0,
        OperationConfirm: 0,
        OperationPay: 0,
        OperationRate: 0,
        OperationRebuy: 2,
        OperationUploadPhoto: 0,
        PayRemainSeconds: 0,
        RatedPoint: 0,
        RemindReplyCount: 0,

        ShopType: 0,
        StatusBar: NewStatusBar(),
        StatusCode: 0,
        TopShow: 0,
    }
}

// Order 订单
type Order struct {
    ID             uint       `gorm:"primary_key" json:"id"`
    CreatedAt      time.Time  `json:"createdAt"`
    UpdatedAt      time.Time  `json:"updatedAt"`
    DeletedAt      *time.Time `sql:"index" json:"deletedAt"`
    UserID         uint       `json:"userId"`
    TotalPrice     float64    `json:"totalPrice"`
    Payment        float64    `json:"payment"`
    Freight        float64    `json:"freight"`
    Remark         string     `json:"remark"`      
    Discount       int        `json:"discount"`
    DeliverStart   time.Time  `json:"deliverStart"`
    DeliverEnd     time.Time  `json:"deliverEnd"`
    Status         int        `json:"status"`
    PayAt          time.Time  `json:"payAt"`
}

// Total 总的订单数
func (order Order) Total() int {
    count := 0
    if DB.Model(&Order{}).Count(&count).Error != nil {
        count = 0   
    }
    return count
}

// TotalSale 总的销售额
func (order Order) TotalSale() float64 {
    result := new(struct{
        TotalSale float64 `gorm:"column:totalPay"` 
    })

    var err = DB.Table("orders").Select("sum(payment) as totalPay").Where("status = ?",
        OrderStatusPaid).Scan(&result).Error

    if err != nil {
        return 0
    }
    return result.TotalSale
}

// CountByDate 指定日期的订单数
func (order Order) CountByDate(date time.Time) int {
    startTime    := date
    startSec     := startTime.Unix();
    tomorrowSec  := startSec + 24 * 60 * 60;
    tomorrowTime := time.Unix(tomorrowSec, 0)
    startYMD     := startTime.Format("2006-01-02")
    tomorrowYMD  := tomorrowTime.Format("2006-01-02")

    var count int
    var err = DB.Model(&Order{}).Where("created_at >= ? AND created_at < ?", 
        startYMD, tomorrowYMD).Count(&count).Error
    if err != nil {
        return 0
    }
    return count
}

// TotalSaleByDate 指定日期的销售额
func (order Order) TotalSaleByDate(date time.Time) float64 {
    startTime    := date
    startSec     := startTime.Unix();
    tomorrowSec  := startSec + 24 * 60 * 60;
    tomorrowTime := time.Unix(tomorrowSec, 0)
    startStr     := startTime.Format("2006-01-02")
    tomorrowStr  := tomorrowTime.Format("2006-01-02")

    result := new(struct{
        TotalPay float64 `gorm:"column:totalPay"` 
    })

    var err = DB.Table("orders").Select("sum(payment) as totalPay").Where("pay_at >= ? AND pay_at < ? AND status = ?",
        startStr, tomorrowStr, OrderStatusPaid).Scan(&result).Error

    if err != nil {
        return 0
    }
    return result.TotalPay
}

const (
    // OrderStatusPending 未支付
    OrderStatusPending  = 0

    // OrderStatusPaid 已支付
    OrderStatusPaid = 1 
)

// OrderPerDay 每天的订单数
type OrderPerDay []struct {
    Count        int    `json:"count"`
    CreatedAt    string `gorm:"column:createdAt" json:"createdAt"` 
}

// Latest30Day 近30天，每天的订单数
func (orders OrderPerDay) Latest30Day() (OrderPerDay) {
    now          := time.Now()
    year         := now.Year()
    month        := now.Month()
    date         := now.Day()
    today        := time.Date(year, month, date, 0, 0, 0, 0, time.Local)

    before29     := today.Unix() - 29 * 24 * 60 * 60; //29天前（秒）
    before29Date := time.Unix(before29, 0)

    sqlData      := before29Date.Format("2006-01-02")
    sqlArr       := []string{
        "SELECT count(id) as count, DATE_FORMAT(created_at,'%Y-%m-%d') as createdAt",
        "FROM `orders`",
        "WHERE created_at > ?",
        "GROUP BY DATE_FORMAT(created_at,'%Y-%m-%d');",
    }
    sql := strings.Join(sqlArr, " ")
    var result OrderPerDay
    var err = DB.Raw(sql, sqlData).Scan(&result).Error
    if err != nil {
        return nil
    }
    return result
}

// AmountPerDay 每天的销售额
type AmountPerDay []struct {
    Amount   float64 `json:"amount"`
    PayAt    string  `gorm:"column:payAt" json:"payAt"` 
}

// AmountLatest30Day 近30天，每天的销售额
func (amount AmountPerDay) AmountLatest30Day() (AmountPerDay) {
    now          := time.Now()
    year         := now.Year()
    month        := now.Month()
    date         := now.Day()
    today        := time.Date(year, month, date, 0, 0, 0, 0, time.Local)

    before29     := today.Unix() - 29 * 24 * 60 * 60; //29天前（秒）
    before29Date := time.Unix(before29, 0)

    sqlData      := before29Date.Format("2006-01-02")
    sqlArr       := []string{
        "SELECT sum(payment) as amount, DATE_FORMAT(pay_at,'%Y-%m-%d') as payAt",
        "FROM `orders`",
        "WHERE pay_at > ? and status = ?",
        "GROUP BY DATE_FORMAT(pay_at,'%Y-%m-%d');",
    };

    sql := strings.Join(sqlArr, " ")
    var result AmountPerDay
    var err = DB.Raw(sql, sqlData, OrderStatusPaid).Scan(&result).Error
    if err != nil {
        return nil
    }
    return result
}