package model

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// type User struct {
//     Name    string        `json:"name,omitempty" bson:"name,omitempty"`
//     Secret  string        `json:"-,omitempty" bson:"secret,omitempty"`
//   }
// user_id                                      int      `json:"user_id" form:"user_id" gorm:"column:user_id"`

// ShopCategory Shop的分类 ---------- superuser 管理
type ShopCategory struct {
	ID            bson.ObjectId     `json:"_id" bson:"_id"`
	Count         int               `json:"count" bson:"count"`
	Dd            int               `json:"id" bson:"id"`
	Ids           []int             `json:"ids" bson:"ids"`
	ImageUrl      string            `json:"image_url" bson:"image_url"`
	Level         int               `json:"level" bson:"level"`
	Name          string            `json:"name" bson:"name"`
	SubCategories []SubShopCategory `json:"sub_categories" bson:"sub_categories"`
}

// FoodCategory Food的分类 ---------- superuser 管理
type FoodCategory struct {
	Description string `json:"description" bson:"description"`
	Id          int    `json:"id" bson:"id"`
	Name        string `json:"name" bson:"name"`
	Type        string `json:"type" bson:"type"`
	IconUrl     string `json:"icon_url" bson:"icon_url"`
}

// ShopCategory Shop的子类 ---------- superuser 管理
type SubShopCategory struct {
	ID       bson.ObjectId `json:"_id" bson:"_id"`
	Count    int           `json:"count" bson:"count"`
	Dd       int           `json:"id" bson:"id"`
	ImageUrl string        `json:"image_url" bson:"image_url"`
	Name     string        `json:"name" bson:"name"`
}

// ShopActivities ---------- superuser 管理
type ShopActivities struct {
	// ID       bson.ObjectId `bson:"_id"`
	Description   string `json:"description" bson:"description"`
	IconColor     string `json:"icon_color" bson:"icon_color"`
	IconName      string `json:"icon_name" bson:"icon_name"`
	Id            int    `json:"id" bson:"id"`
	Name          string `json:"name" bson:"name"`
	RankingWeight int    `json:"ranking_weight" bson:"ranking_weight"`
}

// ShopDeliveryMode---------- superuser 管理
type ShopDeliveryMode struct {
	// ID       bson.ObjectId `bson:"_id"`
	Color   string `json:"color" bson:"color"`
	IsSolid bool   `json:"is_solid" bson:"is_solid"`
	Text    string `json:"text" bson:"text"`
	Id      int    `json:"id" bson:"id"`
}

// ShopDeliveryMode---------- superuser 管理
type ShopSupport struct {
	Description string `json:"description" bson:"description"`
	IconColor   string `json:"icon_color" bson:"icon_color"`
	IconName    string `json:"icon_name" bson:"icon_name"`
	Id          int    `json:"id" bson:"id"`
	Name        string `json:"name" bson:"name"`
}

type ShopIdentification struct {
	// ID       bson.ObjectId `bson:"_id"`
	CompanyName        string    `json:"company_name" bson:"company_name"`
	IdentificateAgency string    `json:"identificate_agency" bson:"identificate_agency"`
	IdentificateDate   time.Time `json:"identificate_date" bson:"identificate_date"`
	LegalPerson        string    `json:"legal_person" bson:"legal_person"`
	LicensesDate       string    `json:"licenses_date" bson:"licenses_date"`
	LicensesNumber     string    `json:"licenses_number" bson:"licenses_number"`
	LicensesScope      string    `json:"licenses_scope" bson:"licenses_scope"`

	OperationPeriod             string `json:"operation_period" bson:"operation_period"`
	RegisteredAddress           string `json:"registered_address" bson:"registered_address"`
	RegisteredNumber            string `json:"registered_number" bson:"registered_number"`
	BusinessLicenseImage        string `json:"business_license_image" bson:"business_license_image"`
	CateringServiceLicenseImage string `json:"catering_service_license_image" bson:"catering_service_license_image"`
}

type ShopLicense struct {
	BusinessLicenseImage        string `json:"business_license_image" bson:"business_license_image"`
	CateringServiceLicenseImage string `json:"catering_service_license_image" bson:"catering_service_license_image"`
}

type ShopPiecewiseAgentFee struct {
	Tips string `json:"tips" bson:"tips"`
}

type ShopMenu struct {
	ShopFoodCategory []FoodCategory `json:"ShopFoodCategory" bson:"ShopFoodCategory"`
}

// Shop
type Shop struct {
	Activities              []ShopActivities      `json:"activities" bson:"activities"`
	Address                 string                `json:"address" bson:"address"`
	SchoolName              string                `json:"school_name" bson:"school_name"`
	SchoolCampusName        string                `json:"school_campus_name" bson:"school_campus_name"`
	SchoolDormName          string                `json:"school_dorm_name" bson:"school_dorm_name"`
	DeliveryMode            ShopDeliveryMode      `json:"delivery_mode" bson:"delivery_mode"`
	Description             string                `json:"description" bson:"description"`
	Id                      int                   `json:"id" bson:"id"`
	Name                    string                `json:"name" bson:"name"`
	OwnerName               string                `json:"owner_name" bson:"owner_name"`
	OwnerId                 int                   `json:"owner_id" bson:"owner_id"`
	OrderLeadTime           string                `json:"order_lead_time" bson:"order_lead_time"`
	Distance                string                `json:"distance" bson:"distance"`
	Location                []float64             `json:"location" bson:"location"`
	FloatDeliveryFee        int                   `json:"float_delivery_fee" bson:"float_delivery_fee"`
	FloatMinimumOrderAmount int                   `json:"float_minimum_order_amount" bson:"float_minimum_order_amount"`
	Category                string                `json:"category" bson:"category"`
	Menu                    ShopMenu              `json:"menu" bson:"menu"`
	Identification          ShopIdentification    `json:"identification" bson:"identification"`
	ImagePath               string                `json:"image_path" bson:"image_path"`
	IsPremium               bool                  `json:"is_premium" bson:"is_premium"`
	IsNew                   bool                  `json:"is_new" bson:"is_new"`
	Latitude                string               `json:"latitude" bson:"latitude"`
	Longitude               string               `json:"longitude" bson:"longitude"`
	License                 ShopLicense           `json:"license" bson:"license"`
	OpeningHours            string                `json:"opening_hours" bson:"opening_hours"`
	StartTime               string                `json:"startTime" bson:"startTime"`
	EndTime                 string                `json:"endTime" bson:"endTime"`
	Phone                   int                   `json:"phone" bson:"phone"`
	PiecewiseAgentFee       ShopPiecewiseAgentFee `json:"piecewise_agent_fee" bson:"piecewise_agent_fee"`
	PromotionInfo           string                `json:"promotion_info" bson:"promotion_info"`
	Rating                  int                   `json:"rating" bson:"rating"`
	RatingCount             int                   `json:"rating_count" bson:"rating_count"`
	RecentOrderNum          int                   `json:"recent_order_num" bson:"recent_order_num"`
	Status                  int                   `json:"status" bson:"status"`
	Supports                []ShopSupport         `json:"supports" bson:"supports"`
}

// type AddShopInfo struct {
// 	//must have
//     Name                       string                `json:"name" bson:"name"`
//     Address                    string                `json:"address" bson:"address"`
//     Phone                      string                `json:"phone" bson:"phone"`
//     Latitude                   float64               `json:"latitude" bson:"latitude"`
//     Longitude                  float64               `json:"longitude" bson:"longitude"`
//     Category                   string                `json:"category" bson:"category"`
//     ImagePath                 string                `json:"image_path" bson:"image_path"`
//     FloatDeliveryFee         int                   `json:"float_delivery_fee" bson:"float_delivery_fee"`
//     FloatMinimumOrderAmount int                   `json:"float_minimum_order_amount" bson:"float_minimum_order_amount"`

// 	//optional
// 	description                    string           `bson:"description"`
// 	promotion_info                 string           `bson:"promotion_info"`
// 	is_premium                     bool             `bson:"is_premium"`
//     delivery_mode                  bool             `bson:"delivery_mode"`

//     Description                string                `json:"description" bson:"description"`
//     PromotionInfo             string                `json:"promotion_info" bson:"promotion_info"`
//     IsPremium                 bool                  `json:"is_premium" bson:"is_premium"`
//     DeliveryMode              bool      `json:"delivery_mode" bson:"delivery_mode"`
//     IsNew                     bool                  `json:"is_new" bson:"is_new"`

// 	bao                            bool             `bson:"bao"`
// 	zhun                           bool             `bson:"zhun"`
// 	piao                           bool             `bson:"piao"`
// 	startTime                      time.Time        `bson:"startTime"`
// 	endTime                        time.Time        `bson:"endTime"`
// 	business_license_image         string           `bson:"business_license_image"`
// 	catering_service_license_image string           `bson:"catering_service_license_image"`
// 	activities                     []ShopActivities `bson:"activities"`

// 	// not exist in Shop DB model
// 	// order_lead_time            string                `bson:"order_lead_time"`
// 	// distance                   string                `bson:"distance"`
// 	// location                   []float64             `bson:"location"`
// 	// identification             ShopIdentification    `bson:"identification"`
// 	// license                    ShopLicense           `bson:"license"`
// 	// opening_hours              string                `bson:"opening_hours"`
// 	// piecewise_agent_fee        ShopPiecewiseAgentFee `bson:"piecewise_agent_fee"`
// 	// rating                     int                   `bson:"rating"`
// 	// rating_count               int                   `bson:"rating_count"`
// 	// recent_order_num           int                   `bson:"recent_order_num"`
// 	// status                     int                   `bson:"status"`
// }
