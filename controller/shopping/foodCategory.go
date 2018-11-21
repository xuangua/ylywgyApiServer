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

// SaveFoodCategory 保存食品分类（）
func SaveFoodCategory(c *gin.Context, isEdit bool) {
    SendErrJSON := common.SendErrJSON

    var category model.FoodCategory
    if err := c.ShouldBindJSON(&category); err != nil {
        SendErrJSON("参数无效", c)
        return
    }

    category.Name = bluemonday.UGCPolicy().Sanitize(category.Name)
    category.Name = strings.TrimSpace(category.Name)

    if category.Name == "" {
        SendErrJSON("分类名称不能为空", c)
        return
    }

    if utf8.RuneCountInString(category.Name) > model.MaxNameLen {
        msg := "分类名称不能超过" + strconv.Itoa(model.MaxNameLen) + "个字符"
        SendErrJSON(msg, c)
        return
    }

    var updatedCategory model.FoodCategory
    if !isEdit {
        //创建分类
        if err := model.MongoDB.C("foodCategories").Insert(&category); err != nil {
            fmt.Println(err)
            SendErrJSON("error.", c)
            return
        }
    }
    // } else {
    //  //更新分类
    //  if err := model.DB.First(&updatedCategory, category.ID).Error; err == nil {
    //      updateMap := make(map[string]interface{})
    //      updateMap["name"] = category.Name
    //      // updateMap["sequence"] = category.Sequence
    //      // updateMap["parent_id"] = category.ParentID
    //      if err := model.DB.Model(&updatedCategory).Updates(updateMap).Error; err != nil {
    //          fmt.Println(err.Error())
    //          SendErrJSON("error", c)
    //          return
    //      }
    //  } else {
    //      SendErrJSON("无效的分类id", c)
    //      return
    //  }

    var categoryJSON model.FoodCategory
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

// CreateFoodCategory 创建店铺分类
func CreateFoodCategory(c *gin.Context) {
    SaveFoodCategory(c, false)
}

// UpdateFoodCategory 更新店铺分类
func UpdateFoodCategory(c *gin.Context) {
    SaveFoodCategory(c, true)
}

// ShopCategoryList 获取店铺所有食品分类列表
func GetFoodCategoryList(c *gin.Context) {
    SendErrJSON := common.SendErrJSON
    var categories []model.FoodCategory

    mgoConn := model.GetMgoSession()
    defer mgoConn.Close()

    // err := mgoConn.FindAll("shopCategories", &categories, bson.M{}, nil, 0, 0, "-_id")
    err := mgoConn.FindAll("foodCategories", &categories, bson.M{}, nil, 0, 0, "$natural")
    if err != nil {
        fmt.Println(err)
        SendErrJSON("error", c)
        return
    }

    // if model.DB.Order("sequence asc").Find(&categories).Error != nil {
    //  SendErrJSON("error", c)
    //  return
    // }

    // c.JSON(http.StatusOK, gin.H{
    //  "errNo": model.ErrorCode.SUCCESS,
    //  "msg":   "success",
    //  "data": gin.H{
    //      "categories": categories,
    //  },
    // })
    c.JSON(http.StatusOK, categories)
}
