package shopping

import (
    "fmt"
    "net/http"
    "strconv"
    "strings"
    "unicode/utf8"

    "github.com/gin-gonic/gin"
    "github.com/globalsign/mgo/bson"
    "github.com/microcosm-cc/bluemonday"

    "github.com/xuangua/ylywgyApiServer/controller/common"
    "github.com/xuangua/ylywgyApiServer/model"
)

// SaveShopCategory 保存店铺分类（创建或更新）
func SaveShopCategory(c *gin.Context, isEdit bool) {
    SendErrJSON := common.SendErrJSON

    // minOrder := model.MinOrder
    // maxOrder := model.MaxOrder

    var category model.ShopCategory
    if err := c.ShouldBindJSON(&category); err != nil {
        SendErrJSON("参数无效", c)
        return
    }

    category.Name = bluemonday.UGCPolicy().Sanitize(category.Name)
    category.Name = strings.TrimSpace(category.Name)

    // categoryArr := strings.Split(category.Name, "/")
    // for i := len(categoryArr) - 1; i >= 0; i-- {
    // }

    if category.Name == "" {
        SendErrJSON("分类名称不能为空", c)
        return
    }

    if utf8.RuneCountInString(category.Name) > model.MaxNameLen {
        msg := "分类名称不能超过" + strconv.Itoa(model.MaxNameLen) + "个字符"
        SendErrJSON(msg, c)
        return
    }

    // if category.Sequence < minOrder || category.Sequence > maxOrder {
    // 	msg := "分类的排序要在" + strconv.Itoa(minOrder) + "到" + strconv.Itoa(maxOrder) + "之间"
    // 	SendErrJSON(msg, c)
    // 	return
    // }

    // if category.ParentID != 0 {
    // 	var parentCate model.ShopCategory
    // 	if err := model.DB.First(&parentCate, category.ParentID).Error; err != nil {
    // 		SendErrJSON("无效的父分类", c)
    // 		return
    // 	}
    // }

    var updatedCategory model.ShopCategory
    if !isEdit {
        //创建分类
        if err := model.MongoDB.C("shopCategories").Insert(&category); err != nil {
            fmt.Println(err)
            SendErrJSON("error.", c)
            return
        }
    }
    // } else {
    // 	//更新分类
    // 	if err := model.DB.First(&updatedCategory, category.ID).Error; err == nil {
    // 		updateMap := make(map[string]interface{})
    // 		updateMap["name"] = category.Name
    // 		// updateMap["sequence"] = category.Sequence
    // 		// updateMap["parent_id"] = category.ParentID
    // 		if err := model.DB.Model(&updatedCategory).Updates(updateMap).Error; err != nil {
    // 			fmt.Println(err.Error())
    // 			SendErrJSON("error", c)
    // 			return
    // 		}
    // 	} else {
    // 		SendErrJSON("无效的分类id", c)
    // 		return
    // 	}

    var categoryJSON model.ShopCategory
    if isEdit {
        categoryJSON = updatedCategory
    } else {
        categoryJSON = category
    }

    c.JSON(http.StatusOK, gin.H{
        "errNo": model.ErrorCode.SUCCESS,
        "msg":   "success",
        "data": gin.H{
            "category": categoryJSON,
        },
    })
}

// CreateShopCategory 创建店铺分类
func CreateShopCategory(c *gin.Context) {
    SaveShopCategory(c, false)
}

// UpdateShopCategory 更新店铺分类
func UpdateShopCategory(c *gin.Context) {
    SaveShopCategory(c, true)
}

// ShopCategoryList 获取所有店铺分类列表
func GetShopCategoryList(c *gin.Context) {
    SendErrJSON := common.SendErrJSON
    var categories []model.ShopCategory

    mgoConn := model.GetMgoSession()
    defer mgoConn.Close()

    // err := mgoConn.WithLog().FindAll("shopCategories", &categories, bson.M{}, nil, 0, 0, "-_id")
    err := mgoConn.WithLog().FindAll("shopCategories", &categories, bson.M{}, nil, 0, 0, "$natural")
    if err != nil {
        fmt.Println(err)
        SendErrJSON("error", c)
        return
    }

    // if model.DB.Order("sequence asc").Find(&categories).Error != nil {
    // 	SendErrJSON("error", c)
    // 	return
    // }

    // c.JSON(http.StatusOK, gin.H{
    // 	"errNo": model.ErrorCode.SUCCESS,
    // 	"msg":   "success",
    // 	"data": gin.H{
    // 		"categories": categories,
    // 	},
    // })
    c.JSON(http.StatusOK, categories)
}
