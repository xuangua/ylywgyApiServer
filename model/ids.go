package model

import (
	// "time"
	"fmt"

	"github.com/globalsign/mgo/bson"
)

// ShopCategory Shop的分类
type Ids struct {
	_id              bson.ObjectId `bson:"_id"`
	Restaurant_id    int           `bson:"restaurant_id"`
	Food_id          int           `bson:"food_id"`
	Order_id         int           `bson:"order_id"`
	User_id          int           `bson:"user_id"`
	Address_id       int           `bson:"address_id"`
	Cart_id          int           `bson:"cart_id"`
	Img_id           int           `bson:"img_id"`
	Shop_category_id int           `bson:"shop_category_id"`
	Item_id          int           `bson:"item_id"`
	Sku_id           int           `bson:"sku_id"`
	Admin_id         int           `bson:"admin_id"`
	Statis_id        int           `bson:"statis_id"`
}

//获取id列表
func GetId(idType string) int {
	// idList := []string{"restaurant_id", "food_id", "order_id", "user_id", "address_id", "cart_id",
	// "img_id", "shop_category_id", "item_id", "sku_id", "admin_id", "statis_id"}
	// i := sort.SearchStrings(idList, idType)
	// if i == len(idList) {
	//     fmt.Println("id类型错误")
	//     panic("id类型错误")
	// }
	var idData Ids

	mgoConn := GetMgoSession()
	defer mgoConn.Close()

	err := mgoConn.WithLog().FindOne("ids", &idData, bson.M{}, nil)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	var idValue int
	switch idType {
	case "restaurant_id":
		idValue = idData.Restaurant_id + 1
	case "food_id":
		idValue = idData.Food_id + 1
	case "order_id":
		idValue = idData.Order_id + 1
	case "user_id":
		idValue = idData.User_id + 1
	case "address_id":
		idValue = idData.Address_id + 1
	case "cart_id":
		idValue = idData.Cart_id + 1
	case "img_id":
		idValue = idData.Img_id + 1
	case "shop_category_id":
		idValue = idData.Shop_category_id + 1
	case "item_id":
		idValue = idData.Item_id + 1
	case "sku_id":
		idValue = idData.Sku_id + 1
	case "admin_id":
		idValue = idData.Admin_id + 1
	case "statis_id":
		idValue = idData.Statis_id + 1
	default:
		fmt.Println("id类型错误")
		panic("id类型错误")
	}
	fmt.Printf("idData: %v\n", idData)

	selector := bson.M{}
	update := bson.M{"$set": bson.M{idType: idValue}}
	err = mgoConn.WithLog().UpdateOne("ids", selector, update)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return idValue
}
