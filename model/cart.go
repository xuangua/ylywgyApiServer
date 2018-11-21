package model

import (
    "time"

    // "github.com/globalsign/mgo/bson"
)

type CartGroupAttrsStruct struct {
}

type CartGroupExtraStruct struct {
}

type CartGroupNewSpecStruct struct {
}

type CartGroupStruct struct {
    Id                      int                     `json:"id" bson:"id"`
    Name                    string                  `json:"name" bson:"name"`
    
    Attrs               []CartGroupAttrsStruct     `json:"attrs" bson:"attrs"`
    Extra               CartGroupExtraStruct     `json:"extra" bson:"extra"`
    NewSpecs               []CartGroupNewSpecStruct     `json:"new_specs" bson:"new_specs"`
    Price                   float64                 `json:"price" bson:"price"`             // {type: Number, default: 0},
    Quantity                int                     `json:"quantity" bson:"quantity"`       // {type: Number, default: 0},
    Specs            []string     `json:"specs" bson:"specs"`
    PackingFee          float64     `json:"packing_fee" bson:"packing_fee"`
    SkuId               int         `json:"sku_id" bson:"sku_id"`
    Stock               int         `json:"stock" bson:"stock"`
}

type CartExtraStruct struct {
    Name                    string              `json:"name" bson:"name"`               // {type: String, default: '餐盒'},
    Description             string              `json:"description" bson:"description"`

    Price                   float64             `json:"price" bson:"price"`             // {type: Number, default: 0},
    Quantity                int                 `json:"quantity" bson:"quantity"`       // {type: Number, default: 0},
    Type                    int                 `json:"type" bson:"type"`               // {type: Number, default: 0},
}

type CartDetailStruct struct {
    Id                      int                     `json:"id" bson:"id"`
    Groups                  []CartGroupStruct     `json:"groups" bson:"groups"`
    Extra                   []CartExtraStruct       `json:"extra" bson:"extra"`

    DeliverAmount            float64     `json:"deliver_amount" bson:"deliver_amount"` // deliver Price
    DeliverTime            string     `json:"deliver_time" bson:"deliver_time"`
    DiscountAmount            string     `json:"discount_amount" bson:"discount_amount"`
    DistInfo            string     `json:"dist_info" bson:"dist_info"`
    IsAddressTooFar            bool     `json:"is_address_too_far" bson:"is_address_too_far"` //  {type: Boolean, default: false},
    IsDeliverByFengniao            bool     `json:"is_deliver_by_fengniao" bson:"is_deliver_by_fengniao"` //  {type: Boolean, default: false},
    IsOnlinePaid            int     `json:"is_online_paid" bson:"is_online_paid"` //  {type: Number, default: 1},
    IsOntimeAvailable            int     `json:"is_ontime_available" bson:"is_ontime_available"` //  {type: Number, default: 0},
    MustNewUser            int     `json:"must_new_user" bson:"must_new_user"` //  {type: Number, default: 0},
    MustPayOnline            int     `json:"must_pay_online" bson:"must_pay_online"` //  {type: Number, default: 0},
    OntimeStatus            int     `json:"ontime_status" bson:"ontime_status"` //  {type: Number, default: 0},
    OntimeUnavailableReason            string     `json:"ontime_unavailable_reason" bson:"ontime_unavailable_reason"`
    OriginalTotal            float64     `json:"original_total" bson:"original_total"`
    Phone            int     `json:"phone" bson:"phone"`
    PromiseDeliveryTime            int     `json:"promise_delivery_time" bson:"promise_delivery_time"` //{type: Number, default: 0},
    ShopId            int     `json:"shop_id" bson:"shop_id"`
    ShopInfo            Shop     `json:"shop_info" bson:"shop_info"`
    RestaurantMinOrderAmount            int     `json:"restaurant_minimum_order_amount" bson:"restaurant_minimum_order_amount"`
    RestaurantNameForUrl            string     `json:"restaurant_name_for_url" bson:"restaurant_name_for_url"`
    RestaurantStatus            int     `json:"restaurant_status" bson:"restaurant_status"` // {type: Number, default: 1},
    ServiceFeeExplanation            int     `json:"service_fee_explanation" bson:"service_fee_explanation"` // {type: Number, default: 0},
    Total            float64     `json:"total" bson:"total"`
    UserId            string     `json:"user_id" bson:"user_id"`
}

type InvoiceStruct struct {
    IsAvailable         bool                    `json:"is_available" bson:"is_available"`
    StatusText          string                  `json:"status_text" bson:"status_text"`
}

type ShipAddressStruct struct {
}

type PromotionStruct struct {
}

type DeliverTimeStruct struct {
}

type MerchantCouponStruct struct {
}

type NumberMealsStruct struct {
}

type DiscountRuleStruct struct {
}

type HongbaoInfoStruct struct {
}

// Carts
type CartTable struct {
    Id                      int                     `json:"id" bson:"id"`
    CustomerUserId          string                  `json:"customer_user_id" bson:"customer_user_id"`
    CreatedDate             time.Time               `json:"created_date" bson:"created_date"`
    UpdatedDate             time.Time               `json:"updated_date" bson:"updated_date"`
    DeletedDate             *time.Time               `json:"deleted_date" bson:"deleted_date"`

    Cart                    CartDetailStruct        `json:"cart" bson:"cart"`
    DeliveryReachTime       string                  `json:"delivery_reach_time" bson:"delivery_reach_time"`
    Invoice                 InvoiceStruct           `json:"invoice" bson:"invoice"`
    Sig                     string                  `json:"sig" bson:"sig"`
    CurrentAddress          ShipAddressStruct       `json:"current_address" bson:"current_address"`
    Payments                []Payments               `json:"payments" bson:"payments"`
    DeliverTimes            []DeliverTimeStruct     `json:"deliver_times" bson:"deliver_times"`
    DeliverTimesV2          []DeliverTimeStruct     `json:"deliver_times_v2" bson:"deliver_times_v2"`
    MerchantCouponInfo      MerchantCouponStruct    `json:"merchant_coupon_info" bson:"merchant_coupon_info"`
    NumberMeals             NumberMealsStruct       `json:"number_of_meals" bson:"number_of_meals"`
    DiscountRule            DiscountRuleStruct      `json:"discount_rule" bson:"discount_rule"`
    HongbaoInfo             HongbaoInfoStruct       `json:"hongbao_info" bson:"hongbao_info"`
    IsSupportCoupon         bool                    `json:"is_support_coupon" bson:"is_support_coupon"`
    IsSupportNinja          int                     `json:"is_support_ninja" bson:"is_support_ninja"`
}

func NewCartTable() CartTable {
    return CartTable{
        CreatedDate: time.Now(),
        Cart: CartDetailStruct{
            IsAddressTooFar: false,
            IsDeliverByFengniao: false,
            IsOnlinePaid: 1,
            IsOntimeAvailable: 0,
            MustNewUser: 0,
            MustPayOnline: 0,
            OntimeStatus: 0,
            OntimeUnavailableReason: "",
            PromiseDeliveryTime: 0,
            RestaurantStatus: 1,
            ServiceFeeExplanation: 0},
        DeliveryReachTime: "",
        Invoice: InvoiceStruct{
            IsAvailable: false,
            StatusText: "商家不支持开发票"},
        Sig: "",
        Payments: []Payments{{
            IsOnlinePayment: true}},
        DeliverTimes: []DeliverTimeStruct{},
        DeliverTimesV2: []DeliverTimeStruct{},
        MerchantCouponInfo: MerchantCouponStruct{},
        NumberMeals: NumberMealsStruct{},
        DiscountRule: DiscountRuleStruct{},
        HongbaoInfo: HongbaoInfoStruct{},
        IsSupportCoupon: false,
        IsSupportNinja: 1}
}