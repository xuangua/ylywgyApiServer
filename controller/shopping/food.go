package shopping

import (
	"fmt"
	"net/http"
	"strconv"
	// "strconv"
	// "strings"

	// "unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/xuangua/ylywgyApiServer/controller/common"
	"github.com/xuangua/ylywgyApiServer/model"
	// "github.com/xuangua/ylywgyApiServer/utils"
)

// foodForm: {
// 	name: '',
// 	description: '',
// 	image_path: '',
// 	activity: '',
// 	attributes: [],
// 	specs: [{
// 		specs: '默认',
// 		  packing_fee: 0,
// 		  price: 20,
// 	}],
// },
// const params = {
// 	...this.foodForm,
// 	category_id: this.selectValue.id,
// 	restaurant_id: this.restaurant_id,
// }

// #### 参数类型：query

// |参数|是否必选|类型|说明|
// |:-----|:-------:|:-----|:-----|
// |restaurant_id      |Y       |int   | 餐馆ID |
// |category_id      |Y       |int   | 分类ID |
// |name      |Y       |string   | 食品名称 |
// |image_path      |Y       |string   | 图片地址 |
// |specs      |Y       |array   | 规格： [{specs: '默认',packing_fee: 0,price: 20,}]|
// |description      |N       |string   |描述 |
// |activity      |N      |string   |活动 |
// |attributes      |N       |array   |特点：[{value: '新',label: '新品'}] |

type AddFoodInfo struct {
	//must have
	ShopId         int     `json:"restaurant_id"`
	FoodCategoryId int     `json:"category_id"`
	Name           string  `json:"name"`
	ImagePath      string  `json:"image_path"`
	Prize          float64 `json:"prize"`

	//optional
	Description string                 `json:"description"`
	activities  []model.FoodActivities `json:"activities"`
}

type GetFoodListParm struct {
	//must have
	offset int `json:"offset"`
	limit  int `json:"limit"`
	shopId int `json:"restaurant_id"`
}

// // FoodTable
// type Foods struct {
// 	ShopSummary ShopSummary `json:"activities" bson:"activities"`
// 	Description string      `json:"description" bson:"description"`
// 	Id          int         `json:"id" bson:"id"`
// 	Name        string      `json:"name" bson:"name"`
// 	ImagePath   string      `json:"image_path" bson:"image_path"`
// 	SkuId       int         `json:"sku_id" bson:"sku_id"`
// 	Prize       float64     `json:"prize" bson:"prize"`
// 	Tips        string      `json:"tips" bson:"tips"`
// 	Rating      float64     `json:"rating" bson:"rating"`

// 	ShopId         int              `json:"shop_id" bson:"shop_id"`
// 	FoodCategoryId int              `json:"food_category_id" bson:"food_category_id"`
// 	Activities     []FoodActivities `json:"activities" bson:"activities"`
// }

// Save 保存    //添加商铺（创建或更新）
func SaveFood(c *gin.Context, isEdit bool) {

}

// 添加商品
func AddFood(c *gin.Context) {
	SendErrJSON := common.SendErrJSON

	var addFoodInfo AddFoodInfo
	if err := c.ShouldBindJSON(&addFoodInfo); err != nil {
		SendErrJSON("参数无效", c)
		return
	}

	if addFoodInfo.Name == "" {
		SendErrJSON("店铺名称不能为空", c)
		return
	}

	// get a new mgo session
	mgoConn := model.GetMgoSession()
	defer mgoConn.Close()

	// verify if this 店铺 已經存在，不存在則提示錯誤
	var shopInfo model.Shop
	err := mgoConn.WithLog().FindOne("shops", &shopInfo, bson.M{"id": addFoodInfo.ShopId}, nil)
	if err != nil || err == mgo.ErrNotFound {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"status":  0,
			"type":    "RESTURANT_NOT_EXISTS",
			"message": "店铺不存在，请先创建店铺",
			"data": gin.H{
				"shopInfo": shopInfo,
			},
		})
		return
	}

	// TODO: enable this
	// verify if this 店铺 包含指定FoodCategoryId，不存在則提示錯誤
	// for _, foodCategory := range shopInfo.Menu.ShopFoodCategory {
	// i := 0
	// for ; i < len(shopInfo.Menu.ShopFoodCategory); i++ {
	// 	if addFoodInfo.FoodCategoryId == shopInfo.Menu.ShopFoodCategory[i].Id {
	// 		break
	// 	}
	// }
	// if i >= len(shopInfo.Menu.ShopFoodCategory) {
	// 	fmt.Println(err)
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"status":  0,
	// 		"type":    "FOOD_CATEGORY_EXISTS",
	// 		"message": "商品分类不存在，请先联系管理员创建商品分类",
	// 		"data": gin.H{
	// 			"shopInfo": shopInfo,
	// 		},
	// 	})
	// 	return
	// }

	// 从ids表中获取最新 food_id
	foodId := model.GetId("food_id")
	if foodId == 0 {
		SendErrJSON("获取商品id失败", c)
		return
	}

	var foodInfo model.Foods
	// required fileds
	foodInfo.Name = addFoodInfo.Name
	foodInfo.Description = addFoodInfo.Description
	foodInfo.ShopId = addFoodInfo.ShopId
	foodInfo.FoodCategoryId = addFoodInfo.FoodCategoryId
	foodInfo.ImagePath = addFoodInfo.ImagePath
	foodInfo.Prize = addFoodInfo.Prize
	foodInfo.ShopSummary = shopInfo

	//     //保存数据，并增加对应食品种类的数量 ？？？
	//     CategoryHandle.addCategory(fields.category)
	//     Rating.initData(restaurant_id);
	//     Food.initData(restaurant_id);

	// userInter, _ := c.Get("user")
	// user := userInter.(model.User)

	// 新增商品
	if err := mgoConn.WithLog().Insert("foods", &foodInfo); err != nil {
		fmt.Println(err)
		SendErrJSON("创建商品数据库失败.", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"foodInfo": foodInfo,
		},
	})
}

// Update 更新商品
func UpdateFood(c *gin.Context) {
	// 	SendErrJSON := common.SendErrJSON

	// 	var shopInfo model.Shop
	// 	if err := c.ShouldBindJSON(&shopInfo); err != nil {
	// 		SendErrJSON("参数无效", c)
	// 		return
	// 	}

	// 	// get a new mgo session
	// 	mgoConn := model.GetMgoSession()
	// 	defer mgoConn.Close()

	// 	shopInfo.Name = "test"

	// 	selector := bson.M{"id": shopInfo.Id}
	// 	update := bson.M{"$set": bson.M{"name": "test"}}
	// 	err := mgoConn.WithLog().UpdateOne("shops", selector, update)
	// 	if err != nil && err != mgo.ErrNotFound {
	// 		fmt.Println(err)
	// 		SendErrJSON("更新店铺数据库失败.", c)
	// 	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		// "data": gin.H{
		// 	"shopInfo": shopInfo,
		// },
	})
}

// Delete 删除商品
func DeleteFood(c *gin.Context) {
	SendErrJSON := common.SendErrJSON

	// foodIDStr := c.Query("foodId")
	// if foodIDStr == "" {
	// 	SendErrJSON("FoodID不正确", c)
	// 	return
	// } else {
	// 	var err error
	// 	if foodID, err = strconv.Atoi(foodIDStr); err != nil {
	// 		fmt.Println(err.Error())
	// 		SendErrJSON("FoodID不正确", c)
	// 		return
	// 	}
	// }

	// id, err := strconv.Atoi(c.Param("bookID"))
	foodID, err := strconv.Atoi(c.Query("foodId"))
	if err != nil {
		SendErrJSON("FoodID不正确", c)
		return
	}

	// get a new mgo session
	mgoConn := model.GetMgoSession()
	defer mgoConn.Close()

	selector := bson.M{"id": foodID}
	err = mgoConn.WithLog().DeleteOne("foods", selector)
	if err != nil && err != mgo.ErrNotFound {
		fmt.Println(err)
		SendErrJSON("删除商品数据库失败.", c)
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		// "data": gin.H{
		// 	"id": book.ID,
		// },
	})
}

// ### 52、获取食品列表
// |参数|是否必选|类型|说明|
// |:-----|:-------:|:-----|:-----|
// |limit      |Y       |int | 获取数据数量，默认 20 |
// |offset      |Y       |int | 跳过数据条数 默认 0 |
// |restaurant_id      |Y       |int | 餐馆id |

// 店铺持有人获取店铺内商品 列表/包含总数
func GetAdminUserShopFoodList(c *gin.Context) {
	SendErrJSON := common.SendErrJSON
	var foodList []model.Foods

	userIDStr := c.Query("userId")
	if userIDStr == "" {
		SendErrJSON("UserID不正确", c)
		return
	}

	var getFoodListParm GetFoodListParm
	if err := c.ShouldBindJSON(&getFoodListParm); err != nil {
		SendErrJSON("参数无效", c)
		return
	}

	mgoConn := model.GetMgoSession()
	defer mgoConn.Close()

	selector := bson.M{"id": getFoodListParm.shopId}
	var shop model.Shop
	err := mgoConn.WithLog().FindOne("shops", &shop, bson.M{"id": getFoodListParm.shopId}, nil)
	if err != nil {
		fmt.Println(err)
		SendErrJSON("店铺不存在", c)
		return
	}

	// TODO: Verify the shop.OwnerId 就是访问者 userId

	selector = bson.M{"shop_id": getFoodListParm.shopId}
	err = mgoConn.WithLog().FindAll("foods", &foodList, selector, nil, getFoodListParm.offset, getFoodListParm.limit, "$natural")
	if err != nil {
		fmt.Println(err)
		SendErrJSON("获取店铺内商品信息失败", c)
		return
	}

	count, err := mgoConn.WithLog().CountAll("foods", selector)
	if err != nil {
		fmt.Println(err)
		SendErrJSON("获取店铺内商品总数失败", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"count": count,
		"data": gin.H{
			"foodList": foodList,
		},
	})
}

// 顾客获取店铺内商品 列表/包含总数 --- should always use this API
func GetCustomerShopFoodList(c *gin.Context) {
	SendErrJSON := common.SendErrJSON
	var foodList []model.Foods

	var getFoodListParm GetFoodListParm
	if err := c.ShouldBindJSON(&getFoodListParm); err != nil {
		SendErrJSON("参数无效", c)
		return
	}

	mgoConn := model.GetMgoSession()
	defer mgoConn.Close()

	selector := bson.M{"shop_id": getFoodListParm.shopId}
	err := mgoConn.WithLog().FindAll("foods", &foodList, selector, nil, getFoodListParm.offset, getFoodListParm.limit, "$natural")
	if err != nil {
		fmt.Println(err)
		SendErrJSON("获取店铺内商品信息失败", c)
		return
	}

	count, err := mgoConn.WithLog().CountAll("foods", selector)
	if err != nil {
		fmt.Println(err)
		SendErrJSON("获取店铺内商品总数失败", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"count": count,
		"data": gin.H{
			"foodList": foodList,
		},
	})
}
