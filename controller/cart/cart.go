package cart

import (
    "fmt"
    "net/http"
    // "strconv"
    // "unicode/utf8"

    "github.com/gin-gonic/gin"
    "github.com/gin-gonic/gin/binding"
    // "github.com/globalsign/mgo"
    "github.com/globalsign/mgo/bson"
    "github.com/xuangua/ylywgyApiServer/controller/common"
    "github.com/xuangua/ylywgyApiServer/model"
    // "github.com/xuangua/ylywgyApiServer/utils"
)

type AddCartInfo struct {
    ComeFrom                   string       `json:"come_from"`
    ShopId                     int          `json:"shop_id"`
    Geohash                    string       `json:"geohash"`
    Latitude                   string       `json:"latitude"`
    Longitude                  string       `json:"longitude"`
    Entities                   []model.CartGroupStruct     `json:"entities" bson:"entities"`
}

// 添加顾客购物车表项
func CustomerCheckoutCartHandler(c *gin.Context) {
    SendErrJSON := common.SendErrJSON
    var err error

    var addCartInfo AddCartInfo
    if err := c.ShouldBindWith(&addCartInfo, binding.JSON); err != nil {
        fmt.Printf("get addCartInfo err: %v\n", err)
        SendErrJSON("参数无效", c)
        return
    }
    // fmt.Printf("addCartInfo: %v\n", addCartInfo)

    // var shopId int
    // shopId, err := strconv.Atoi(addCartInfo.ShopId)
    // if err != nil {
    //     SendErrJSON("错误的 ShopId", c)
    //     return
    // }
    shopId := addCartInfo.ShopId

    // get a new mgo session
    mgoConn := model.GetMgoSession()
    defer mgoConn.Close()

    // TODO: if need verify this Cart 已經存在? no needed now.
    newCart := model.NewCartTable()

    // 从ids表中获取最新shop_id
    shopcartId := model.GetId("shopcart_id")
    if shopcartId == 0 {
        SendErrJSON("获取 shopcartId 失败", c)
        return
    }

    var paymentList []model.Payments
    selector := bson.M{}
    err = mgoConn.FindAll("payments", &paymentList, selector, nil, 0, 0, "$natural")
    if err != nil {
        fmt.Println(err)
        SendErrJSON("获取 数据库-支付方式 信息失败", c)
        return
    }

    selector = bson.M{"id": shopId}
    var ShopData model.Shop
    err = mgoConn.FindOne("shops", &ShopData, selector, nil)
    if err != nil {
        fmt.Println(err)
        SendErrJSON("店铺不存在", c)
        return
    }

    // TODO: getDistance for delivery_reach_time
    // const to = restaurant.latitude+ ',' + restaurant.longitude;
    // deliver_time = await this.getDistance(from, to, 'tiemvalue');
    // let time = new Date().getTime() + deliver_time*1000;
    // let hour = ('0' + new Date(time).getHours()).substr(-2);
    // let minute = ('0' + new Date(time).getMinutes()).substr(-2);
    // delivery_reach_time = hour + ':' + minute;

    deliverAmount := 2.0 //deliver price
    price := float64(0) //食品价格
    cartExtra := []model.CartExtraStruct{{
        Description: "",
        Name: "餐盒",
        Price: float64(0),
        Quantity: 1,
        Type: 0}}

    for _, cartGroupEntity := range addCartInfo.Entities {
        // fmt.Printf("food cartGroupEntity : %v\n", cartGroupEntity)
        price += cartGroupEntity.Price * float64(cartGroupEntity.Quantity)
        if (cartGroupEntity.PackingFee != 0) {
            cartExtra[0].Price += cartGroupEntity.PackingFee*float64(cartGroupEntity.Quantity)
        }
        if (cartGroupEntity.Specs[0] != "") {
            cartGroupEntity.Name = cartGroupEntity.Name + "-" + cartGroupEntity.Specs[0]
        }
    }

    //食品总价格
    total := price + cartExtra[0].Price * float64(cartExtra[0].Quantity) + float64(deliverAmount);

    // userInter, _ := c.Get("user")
    // user := userInter.(model.User)

    // required fileds
    newCart.Id = shopcartId
    newCart.Cart = model.CartDetailStruct{
        Id: shopcartId,
        Groups: addCartInfo.Entities,
        Extra: cartExtra,
        DeliverAmount: float64(deliverAmount),
        OriginalTotal: total,
        Phone: ShopData.Phone,
        ShopId: ShopData.Id,
        ShopInfo: ShopData,
        RestaurantMinOrderAmount: ShopData.FloatMinimumOrderAmount,
        Total: total,
        UserId: ""}
    newCart.Payments = paymentList

    if err := mgoConn.Insert("carts", &newCart); err != nil {
        fmt.Println(err)
        SendErrJSON("创建 carts 数据库失败.", c)
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "errNo": model.ErrorCode.SUCCESS,
        "msg":   "success",
        "data": gin.H{
            "newCart": newCart,
        },
    })
}
