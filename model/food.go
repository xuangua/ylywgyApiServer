package model

import (
	"time"
)

// FoodActivities ---------- superuser 管理
type FoodActivities struct {
	// ID       bson.ObjectId `bson:"_id"`
	Description   string `json:"description" bson:"description"`
	IconColor     string `json:"icon_color" bson:"icon_color"`
	IconName      string `json:"icon_name" bson:"icon_name"`
	Id            int    `json:"id" bson:"id"`
	Name          string `json:"name" bson:"name"`
	RankingWeight int    `json:"ranking_weight" bson:"ranking_weight"`
}

// ShopSummary
type ShopSummary struct {
	Activities              []ShopActivities      `json:"activities" bson:"activities"`
	Address                 string                `json:"address" bson:"address"`
	SchoolName              string                `json:"school_name" bson:"school_name"`
	SchoolCampusName        string                `json:"school_campus_name" bson:"school_campus_name"`
	SchoolDormName          string                `json:"school_dorm_name" bson:"school_dorm_name"`
	DeliveryMode            ShopDeliveryMode      `json:"delivery_mode" bson:"delivery_mode"`
	Description             string                `json:"description" bson:"description"`
	Id                      int                   `json:"id" bson:"id"`
	Name                    string                `json:"name" bson:"name"`
	OwnerName               string                `json:"OwnerName" bson:"OwnerName"`
	OrderLeadTime           string                `json:"order_lead_time" bson:"order_lead_time"`
	Distance                string                `json:"distance" bson:"distance"`
	Location                []float64             `json:"location" bson:"location"`
	FloatDeliveryFee        int                   `json:"float_delivery_fee" bson:"float_delivery_fee"`
	FloatMinimumOrderAmount int                   `json:"float_minimum_order_amount" bson:"float_minimum_order_amount"`
	Category                string                `json:"category" bson:"category"`
	Identification          ShopIdentification    `json:"identification" bson:"identification"`
	ImagePath               string                `json:"image_path" bson:"image_path"`
	IsPremium               bool                  `json:"is_premium" bson:"is_premium"`
	IsNew                   bool                  `json:"is_new" bson:"is_new"`
	Latitude                float64               `json:"latitude" bson:"latitude"`
	Longitude               float64               `json:"longitude" bson:"longitude"`
	License                 ShopLicense           `json:"license" bson:"license"`
	OpeningHours            string                `json:"opening_hours" bson:"opening_hours"`
	StartTime               time.Time             `json:"startTime" bson:"startTime"`
	EndTime                 time.Time             `json:"endTime" bson:"endTime"`
	Phone                   string                `json:"phone" bson:"phone"`
	PiecewiseAgentFee       ShopPiecewiseAgentFee `json:"piecewise_agent_fee" bson:"piecewise_agent_fee"`
	PromotionInfo           string                `json:"promotion_info" bson:"promotion_info"`
	Rating                  int                   `json:"rating" bson:"rating"`
	RatingCount             int                   `json:"rating_count" bson:"rating_count"`
	RecentOrderNum          int                   `json:"recent_order_num" bson:"recent_order_num"`
	Status                  int                   `json:"status" bson:"status"`
	Supports                []ShopSupport         `json:"supports" bson:"supports"`
}

// FoodTable
type Foods struct {
	ShopSummary Shop `json:"shop" bson:"shop"`
	Description string      `json:"description" bson:"description"`
	Id          int         `json:"id" bson:"id"`
	Name        string      `json:"name" bson:"name"`
	ImagePath   string      `json:"image_path" bson:"image_path"`
	SkuId       int         `json:"sku_id" bson:"sku_id"`
	Prize       float64     `json:"prize" bson:"prize"`
	Tips        string      `json:"tips" bson:"tips"`
	Rating      float64     `json:"rating" bson:"rating"`

	ShopId         int              `json:"shop_id" bson:"shop_id"`
	FoodCategoryId int              `json:"food_category_id" bson:"food_category_id"`
	Activities     []FoodActivities `json:"activities" bson:"activities"`
}
