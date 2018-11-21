package shopping

import (
    "fmt"
    "net/http"
    // "strconv"
    // "unicode/utf8"

    "github.com/gin-gonic/gin"
    "github.com/gin-gonic/gin/binding"

    "github.com/xuangua/ylywgyApiServer/controller/common"
    "github.com/xuangua/ylywgyApiServer/model"
    // "github.com/xuangua/ylywgyApiServer/utils"
)

// type AddShopInfo struct {
//  //must have
//  // name      string  `json:"name"`
//  name      string  `json:"name" binding:"required"`
//  address   string  `json:"address"`
//  phone     string  `json:"phone"`
//  latitude  float64 `json:"latitude"`
//  longitude float64 `json:"longitude"`
//  // category string `json:"category"`
//  // image_path                 string  `json:"image_path"`
//  float_delivery_fee         int `json:"float_delivery_fee"`
//  float_minimum_order_amount int `json:"float_minimum_order_amount"`

//  // //optional
//  // description                    string                 `json:"description"`
//  // promotion_info                 string                 `json:"promotion_info"`
//  // is_premium                     bool                   `json:"is_premium"`
//  // delivery_mode                  bool                   `json:"delivery_mode"`
//  // is_new                         bool                   `json:"new"`
//  // bao                            bool                   `json:"bao"`
//  // zhun                           bool                   `json:"zhun"`
//  // piao                           bool                   `json:"piao"`
//  // startTime                      time.Time              `json:"startTime"`
//  // endTime                        time.Time              `json:"endTime"`
//  // business_license_image         string                 `json:"business_license_image"`
//  // catering_service_license_image string                 `json:"catering_service_license_image"`
//  // activities                     []model.ShopActivities `json:"activities"`

//  // OwnerId            int    `json:"owner_id" bson:"owner_id"`
//  // school_name        string `json:"school_name" bson:"school_name"`
//  // school_campus_name string `json:"school_campus_name" bson:"school_campus_name"`
//  // school_dorm_name   string `json:"school_dorm_name" bson:"school_dorm_name"`

//  // not exist in Shop DB model
//  // order_lead_time            string                `bson:"order_lead_time"`
//  // distance                   string                `bson:"distance"`
//  // location                   []float64             `bson:"location"`
//  // identification             ShopIdentification    `bson:"identification"`
//  // license                    ShopLicense           `bson:"license"`
//  // opening_hours              string                `bson:"opening_hours"`
//  // piecewise_agent_fee        ShopPiecewiseAgentFee `bson:"piecewise_agent_fee"`
//  // rating                     int                   `bson:"rating"`
//  // rating_count               int                   `bson:"rating_count"`
//  // recent_order_num           int                   `bson:"recent_order_num"`
//  // status                     int                   `bson:"status"`
// }

// Create 创建商店
func AddShop2(c *gin.Context) {
    SendErrJSON := common.SendErrJSON

    type AddShopInfo struct {
        name    string `json:"name" binding:"required"`
        address string `json:"address" binding:"required"`
        // phone     string  `json:"phone"`
        // latitude  float64 `json:"latitude"`
        // longitude float64 `json:"longitude"`
    }

    var addShopInfo AddShopInfo
    if err := c.ShouldBindWith(&addShopInfo, binding.JSON); err != nil {
        SendErrJSON("参数无效", c)
        return
    }

    // TODO: Verify addShopInfo
    // if addShopInfo.ContentType != model.ContentTypeMarkdown && bookData.ContentType != model.ContentTypeHTML {
    //     SendErrJSON("无效的图书格式", c)
    //     return
    // }
    // TODO: avoidXSS and TrimSpace
    // bookData.Name = utils.AvoidXSS(bookData.Name)
    // bookData.Name = strings.TrimSpace(bookData.Name)

    // if utf8.RuneCountInString(addShopInfo.name) > model.MaxNameLen {
    //     msg := "店铺名称不能超过" + strconv.Itoa(model.MaxNameLen) + "个字符"
    //     SendErrJSON(msg, c)
    //     return
    // }

    fmt.Println("addShopInfo: ", addShopInfo)
    if addShopInfo.name == "" {
        SendErrJSON("店铺名称不能为空", c)
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "errNo": model.ErrorCode.SUCCESS,
        "msg":   "success",
        "data": gin.H{
            "shopInfo": addShopInfo,
        },
    })
}
