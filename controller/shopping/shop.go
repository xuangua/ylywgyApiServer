package shopping

import (
	"fmt"
	"net/http"
	"strconv"
	// "unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/xuangua/ylywgyApiServer/controller/common"
	"github.com/xuangua/ylywgyApiServer/model"
	// "github.com/xuangua/ylywgyApiServer/utils"
)

// // AddShopActivities
// type AddShopActivities struct {
// 	// ID       bson.ObjectId `bson:"_id"`
// 	description   string `json:"description" bson:"description"`
// 	icon_color     string `json:"icon_color" bson:"icon_color"`
// 	icon_name      string `json:"icon_name" bson:"icon_name"`
// 	id            int    `json:"id" bson:"id"`
// 	name          string `json:"name" bson:"name"`
// 	ranking_weight int    `json:"ranking_weight" bson:"ranking_weight"`
// }
// activities                     []AddShopActivities `bson:"activities"`

type AddShopInfo struct {
	//must have
	// name      string  `json:"name"`
	Name                       string `json:"name" binding:"required"`
	Address                    string `json:"address"`
	Phone                      int    `json:"phone"`
	Latitude                   string `json:"latitude"`
	Longitude                  string `json:"longitude"`
	Category                   string `json:"category"`
	Image_path                 string `json:"image_path"`
	Float_delivery_fee         int    `json:"float_delivery_fee"`
	Float_minimum_order_amount int    `json:"float_minimum_order_amount"`

	//optional
	Description                    string                 `json:"description"`
	Promotion_info                 string                 `json:"promotion_info"`
	Is_premium                     bool                   `json:"is_premium"`
	Delivery_mode                  bool                   `json:"delivery_mode"`
	Is_new                         bool                   `json:"new"`
	Bao                            bool                   `json:"bao"`
	Zhun                           bool                   `json:"zhun"`
	Piao                           bool                   `json:"piao"`
	StartTime                      string                 `json:"startTime"`
	EndTime                        string                 `json:"endTime"`
	Business_license_image         string                 `json:"business_license_image"`
	Catering_service_license_image string                 `json:"catering_service_license_image"`
	Activities                     []model.ShopActivities `json:"activities"`

	OwnerId          int    `json:"owner_id" bson:"owner_id"`
	SchoolName       string `json:"school_name" bson:"school_name"`
	SchoolCampusName string `json:"school_campus_name" bson:"school_campus_name"`
	SchoolDormName   string `json:"school_dorm_name" bson:"school_dorm_name"`

	// not exist in Shop DB model
	// order_lead_time            string                `bson:"order_lead_time"`
	// distance                   string                `bson:"distance"`
	// location                   []float64             `bson:"location"`
	// identification             ShopIdentification    `bson:"identification"`
	// license                    ShopLicense           `bson:"license"`
	// opening_hours              string                `bson:"opening_hours"`
	// piecewise_agent_fee        ShopPiecewiseAgentFee `bson:"piecewise_agent_fee"`
	// rating                     int                   `bson:"rating"`
	// rating_count               int                   `bson:"rating_count"`
	// recent_order_num           int                   `bson:"recent_order_num"`
	// status                     int                   `bson:"status"`
}

type GetShopListParm struct {
	//must have
	latitude   float64 `json:"latitude"`
	longitude  float64 `json:"longitude"`
	offset     int     `json:"offset"`
	limit      int     `json:"limit"`
	category   string  `json:"category"`
	image_path string  `json:"image_path"`
	order_by   int     `json:"order_by"`

	OwnerId            uint   `json:"owner_id" bson:"owner_id"`
	school_name        string `json:"school_name" bson:"school_name"`
	school_campus_name string `json:"school_campus_name" bson:"school_campus_name"`
	school_dorm_name   string `json:"school_dorm_name" bson:"school_dorm_name"`
	// |latitude      |Y       |string  |纬度|
	// |longitude      |Y       |string  |经度|
	// |offset      |N       |int |跳过多少条数据，默认0|
	// |limit      |N      |int |请求数据的数量，默认20|
	// |restaurant_category_id      |N      |int |餐馆分类id|
	// |order_by      |N       |int |排序方式id： 1：起送价、2：配送速度、3:评分、4: 智能排序(默认)、5:距离最近、6:销量最高|
	// |delivery_mode      |N      |array |配送方式id|
	// |support_ids      |N      |array |餐馆支持特权的id|
	// |restaurant_category_ids      |N      |array |餐馆分类id|
}

// Save 保存    //添加商铺（创建或更新）
func Save(c *gin.Context, isEdit bool) {

}

// Create 创建商店
func AddShop(c *gin.Context) {
	SendErrJSON := common.SendErrJSON

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

	fmt.Printf("addShopInfo: %v\n", addShopInfo)
	if addShopInfo.Name == "" {
		SendErrJSON("店铺名称不能为空", c)
		return
	}

	// TODO: 是否需要 Verify the getShopListParm.OwnerId 就是访问者 userId (从cookie/token获取)？
	// 还是直接信赖cookie/token？

	// get a new mgo session
	mgoConn := model.GetMgoSession()
	defer mgoConn.Close()

	// verify if this 店铺 已經存在，已存在則提示錯誤
	var shopInfo model.Shop
	err := mgoConn.WithLog().FindOne("shops", &shopInfo, bson.M{"name": addShopInfo.Name}, nil)
	if err == nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"status":  0,
			"type":    "RESTURANT_EXISTS",
			"message": "店铺已存在，请尝试其他店铺名称",
			"data": gin.H{
				"shopInfo": shopInfo,
			},
		})
		return
	}

	// 从ids表中获取最新shop_id
	shopId := model.GetId("restaurant_id")
	if shopId == 0 {
		SendErrJSON("获取商店id失败", c)
		return
	}

	// required fileds
	shopInfo.Id = shopId
	shopInfo.Name = addShopInfo.Name
	shopInfo.Address = addShopInfo.Address
	shopInfo.Phone = addShopInfo.Phone
	shopInfo.Latitude = addShopInfo.Latitude
	shopInfo.Longitude = addShopInfo.Longitude
	shopInfo.Category = addShopInfo.Category
	shopInfo.ImagePath = addShopInfo.Image_path
	shopInfo.FloatDeliveryFee = addShopInfo.Float_delivery_fee
	shopInfo.FloatMinimumOrderAmount = addShopInfo.Float_minimum_order_amount

	//optional fileds
	shopInfo.Description = addShopInfo.Description
	// 促销信息
	shopInfo.PromotionInfo = addShopInfo.Promotion_info
	// is 品牌店
	shopInfo.IsPremium = addShopInfo.Is_premium
	// shopInfo.DeliveryMode = addShopInfo.delivery_mode
	//配送方式
	if addShopInfo.Delivery_mode {
		shopInfo.DeliveryMode = model.ShopDeliveryMode{
			Color:   "57A9FF",
			Id:      1,
			Text:    "蜂鸟专送",
			IsSolid: true}
	}
	// is 新店开业
	shopInfo.IsNew = addShopInfo.Is_new

	// is 支持外卖保
	if addShopInfo.Bao {
		support := model.ShopSupport{
			Description: "已加入“外卖保”计划，食品安全有保障",
			IconColor:   "999999",
			IconName:    "保",
			Id:          7,
			Name:        "外卖保"}
		shopInfo.Supports = append(shopInfo.Supports, support)
	}
	// is 支持准时达
	if addShopInfo.Zhun {
		support := model.ShopSupport{
			Description: "准时必达，超时秒赔",
			IconColor:   "57A9FF",
			IconName:    "准",
			Id:          9,
			Name:        "准时达"}
		shopInfo.Supports = append(shopInfo.Supports, support)
	}
	// is 支持开发票
	if addShopInfo.Piao {
		support := model.ShopSupport{
			Description: "该商家支持开发票，请在下单时填写好发票抬头",
			IconColor:   "999999",
			IconName:    "票",
			Id:          4,
			Name:        "开发票"}
		shopInfo.Supports = append(shopInfo.Supports, support)
	}

	shopInfo.StartTime = addShopInfo.StartTime
	shopInfo.EndTime = addShopInfo.EndTime

	shopInfo.OpeningHours = "8:30/20:30"
	// if addShopInfo.startTime != "" && addShopInfo.endTime != "" {
	//     shopInfo.OpeningHours = addShopInfo.startTime + "/" + addShopInfo.endTime;
	// }

	var shopLicense model.ShopLicense
	shopLicense.BusinessLicenseImage = addShopInfo.Business_license_image
	shopLicense.CateringServiceLicenseImage = addShopInfo.Catering_service_license_image
	shopInfo.License = shopLicense

	//商店支持的活动
	index := 0
	for _, activitiy := range addShopInfo.Activities {
		switch activitiy.IconName {
		case "减":
			activitiy.IconColor = "f07373"
			activitiy.Id = index + 1
			activitiy.Name = "减"
			activitiy.Description = "减"
			activitiy.RankingWeight = 1
			break
		case "特":
			activitiy.IconColor = "EDC123"
			activitiy.Id = index + 1
			activitiy.Name = "特"
			activitiy.Description = "特"
			activitiy.RankingWeight = 2
			break
		case "新":
			activitiy.IconColor = "70bc46"
			activitiy.Id = index + 1
			activitiy.Name = "新"
			activitiy.Description = "新用户优惠"
			activitiy.RankingWeight = 3
			break
		case "领":
			activitiy.IconColor = "E3EE0D"
			activitiy.Id = index + 1
			activitiy.Name = "领"
			activitiy.Description = "领"
			activitiy.RankingWeight = 4
			break
		}
		shopInfo.Activities = append(shopInfo.Activities, activitiy)
	}

	//     //保存数据，并增加对应食品种类的数量 ？？？
	//     CategoryHandle.addCategory(fields.category)
	//     Rating.initData(restaurant_id);
	//     Food.initData(restaurant_id);

	userInter, _ := c.Get("user")
	user := userInter.(model.User)
	// 店铺所有者用户名
	shopInfo.OwnerName = user.Name
	// 店铺所有者用户ID
	shopInfo.OwnerId = int(user.ID)
	shopInfo.SchoolName = addShopInfo.SchoolName
	shopInfo.SchoolCampusName = addShopInfo.SchoolCampusName
	shopInfo.SchoolDormName = addShopInfo.SchoolDormName

	// 创建店铺
	// if err := model.MongoDB.C("shops").Insert(&shopInfo); err != nil {
	// 	fmt.Println(err)
	// 	SendErrJSON("创建店铺数据库失败.", c)
	// 	return
	// }
	if err := mgoConn.WithLog().Insert("shops", &shopInfo); err != nil {
		fmt.Println(err)
		SendErrJSON("创建店铺数据库失败.", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"shopInfo": shopInfo,
		},
	})
}

// Update 更新商店
func UpdateShop(c *gin.Context) {
	SendErrJSON := common.SendErrJSON

	var shopInfo model.Shop
	if err := c.ShouldBindJSON(&shopInfo); err != nil {
		SendErrJSON("参数无效", c)
		return
	}

	// TODO: 是否需要 Verify the getShopListParm.OwnerId 就是访问者 userId (从cookie/token获取)？
	// 还是直接信赖cookie/token？
	// userInter, _ := c.Get("user")
	// user := userInter.(model.User)
	// if uint(shopInfo.OwnerId) != user.ID {
	// 	SendErrJSON("店铺 OwnerId 参数无效", c)
	// 	return
	// }

	// get a new mgo session
	mgoConn := model.GetMgoSession()
	defer mgoConn.Close()

	selector := bson.M{"id": shopInfo.Id}
	update := bson.M{"$set": bson.M{
		"name":    shopInfo.Name,
		"address": shopInfo.Address,
		"phone":   shopInfo.Phone}}
	err := mgoConn.WithLog().UpdateOne("shops", selector, update)
	// if err != nil && err != mgo.ErrNotFound {
	if err != nil {
		fmt.Println(err)
		SendErrJSON("更新店铺数据库失败.", c)
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"shopInfo": shopInfo,
		},
	})
}

// Delete 删除商店
func DeleteShop(c *gin.Context) {
	SendErrJSON := common.SendErrJSON

	// var shopInfo model.Shop
	// if err := c.ShouldBindJSON(&shopInfo); err != nil {
	// 	SendErrJSON("参数无效", c)
	// 	return
	// }
	shopId, err := strconv.Atoi(c.Param("shopId"))
	if err != nil {
		SendErrJSON("错误的shopId", c)
		return
	}
	fmt.Printf("shopId: %v\n", shopId)
	// TODO: 是否需要 Verify the getShopListParm.OwnerId 就是访问者 userId (从cookie/token获取)？
	// 还是直接信赖cookie/token？
	// userInter, _ := c.Get("user")
	// user := userInter.(model.User)
	// if uint(shopInfo.OwnerId) != user.ID {
	// 	SendErrJSON("店铺 OwnerId 参数无效", c)
	// 	return
	// }

	// get a new mgo session
	mgoConn := model.GetMgoSession()
	defer mgoConn.Close()

	selector := bson.M{"id": shopId}
	err = mgoConn.WithLog().DeleteOneReal("shops", selector)
	if err != nil && err != mgo.ErrNotFound {
		fmt.Println(err)
		SendErrJSON("删除店铺数据库失败.", c)
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		// "data": gin.H{
		// 	"id": book.ID,
		// },
	})
}

// ### 6、获取商铺列表
// #### 请求URL：
// ```
// https://elm.cangdu.org/shopping/restaurants
// ```
// #### 示例：
// [https://elm.cangdu.org/shopping/restaurants?latitude=31.22967&longitude=121.4762](https://elm.cangdu.org/shopping/restaurants?latitude=31.22967&longitude=121.4762)
// #### 请求方式：
// ```
// GET
// ```
// #### 参数类型：query
// |参数|是否必选|类型|说明|
// |:-----|:-------:|:-----|:-----|
// |latitude      |Y       |string  |纬度|
// |longitude      |Y       |string  |经度|
// |OwnerId      |Y       |string  |店铺所有者|
// |offset      |N       |int |跳过多少条数据，默认0|
// |limit      |N      |int |请求数据的数量，默认20|
// |restaurant_category_id      |N      |int |餐馆分类id|
// |order_by      |N       |int |排序方式id： 1：起送价、2：配送速度、3:评分、4: 智能排序(默认)、5:距离最近、6:销量最高|
// |delivery_mode      |N      |array |配送方式id|
// |support_ids      |N      |array |餐馆支持特权的id|
// |restaurant_category_ids      |N      |array |餐馆分类id|

// GetCustomerShopList 获取顾客可见店铺商店列表
func GetCustomerShopList(c *gin.Context) {
	SendErrJSON := common.SendErrJSON
	var shopList []model.Shop

	var getShopListParm GetShopListParm
	if err := c.ShouldBindJSON(&getShopListParm); err != nil {
		SendErrJSON("参数无效", c)
		return
	}

	mgoConn := model.GetMgoSession()
	defer mgoConn.Close()

	// err := mgoConn.WithLog().FindAll("shopCategories", &categories, bson.M{}, nil, 0, 0, "-_id")
	selector := bson.M{
		"school_name":        getShopListParm.school_name,
		"school_campus_name": getShopListParm.school_campus_name,
		"school_dorm_name":   getShopListParm.school_dorm_name}
	err := mgoConn.WithLog().FindAll("shops", &shopList, selector, nil, getShopListParm.offset, getShopListParm.limit, "$natural")
	if err != nil {
		fmt.Println(err)
		SendErrJSON("error", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"shopList": shopList,
		},
	})
}

// GetAdminUserShopList 获取店铺管理员商店列表
func GetAdminUserShopList(c *gin.Context) {
	SendErrJSON := common.SendErrJSON
	var shopList []model.Shop

	var err error

	var offset int
	if offset, err = strconv.Atoi(c.Query("offset")); err != nil {
		SendErrJSON("参数offset无效", c)
		return
	}

	var limit int
	if limit, err = strconv.Atoi(c.Query("limit")); err != nil {
		SendErrJSON("参数limit无效", c)
		return
	}

	// TODO: 是否需要 Verify the getShopListParm.OwnerId 就是访问者 userId (从cookie/token获取)？
	// 还是直接信赖cookie/token？

	userInter, _ := c.Get("user")
	user := userInter.(model.User)
	// if uint(ownerId) != user.ID {
	// 	SendErrJSON("店铺 OwnerId 参数无效", c)
	// 	return
	// }

	mgoConn := model.GetMgoSession()
	defer mgoConn.Close()

	// err := mgoConn.WithLog().FindAll("shopCategories", &categories, bson.M{}, nil, 0, 0, "-_id")
	err = mgoConn.WithLog().FindAll("shops", &shopList, bson.M{"owner_id": user.ID}, nil, offset, limit, "$natural")
	if err != nil {
		fmt.Println(err)
		SendErrJSON("error", c)
		return
	}

	// c.JSON(http.StatusOK, gin.H{
	// 	"errNo": model.ErrorCode.SUCCESS,
	// 	"msg":   "success",
	// 	"data": gin.H{
	// 		"shopList": shopList,
	// 	},
	// })
	c.JSON(http.StatusOK, shopList)
}

//获取店铺管理员-商店详细信息
// GetAdminUserShopDetail 获取店铺管理员商店列表
func GetAdminUserShopDetail(c *gin.Context) {
}

//根据商店名,搜索店铺
// GetSearchedShopListWithShopName 获取店铺管理员商店列表
func GetSearchedShopListWithShopName(c *gin.Context) {
}

//获取顾客可见店铺总数
// GetCustomerTotalShopCount 获取顾客可见店铺总数
func GetCustomerTotalShopCount(c *gin.Context) {
	SendErrJSON := common.SendErrJSON

	var getShopListParm GetShopListParm
	if err := c.ShouldBindJSON(&getShopListParm); err != nil {
		SendErrJSON("参数无效", c)
		return
	}

	mgoConn := model.GetMgoSession()
	defer mgoConn.Close()

	// err := mgoConn.WithLog().FindAll("shopCategories", &categories, bson.M{}, nil, 0, 0, "-_id")
	selector := bson.M{
		"school_name":        getShopListParm.school_name,
		"school_campus_name": getShopListParm.school_campus_name,
		"school_dorm_name":   getShopListParm.school_dorm_name}
	count, err := mgoConn.WithLog().CountAll("shops", selector)
	if err != nil {
		fmt.Println(err)
		SendErrJSON("error", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"count": count,
		},
	})
}

//获取店铺管理员所持有店铺总数
// GetAdminUserTotalShopCount 获取管理员所持有店铺总数
func GetAdminUserTotalShopCount(c *gin.Context) {
	SendErrJSON := common.SendErrJSON

	var ownerIdErr error
	var ownerId int
	if ownerId, ownerIdErr = strconv.Atoi(c.Query("owner_id")); ownerIdErr != nil {
		SendErrJSON("参数无效", c)
		return
	}

	// TODO: 是否需要 Verify the getShopListParm.OwnerId 就是访问者 userId (从cookie/token获取)？
	// 还是直接信赖cookie/token？

	userInter, _ := c.Get("user")
	user := userInter.(model.User)
	if uint(ownerId) != user.ID {
		SendErrJSON("店铺 OwnerId 参数无效", c)
		return
	}

	mgoConn := model.GetMgoSession()
	defer mgoConn.Close()

	selector := bson.M{"owner_id": ownerId}
	count, err := mgoConn.WithLog().CountAll("shops", selector)
	if err != nil {
		fmt.Println(err)
		SendErrJSON("error", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"count": count,
		},
	})
}
