package order

import (
    "fmt"
    "net/http"
    "strconv"
    "time"
    // "unicode/utf8"

    "github.com/gin-gonic/gin"
    "github.com/gin-gonic/gin/binding"
    // "github.com/globalsign/mgo"
    "github.com/globalsign/mgo/bson"
    "github.com/xuangua/ylywgyApiServer/controller/common"
    "github.com/xuangua/ylywgyApiServer/model"
    // "github.com/xuangua/ylywgyApiServer/utils"
)

type AddOrderInfo struct {
    ComeFrom                   string       `json:"come_from"`
    AddressId                     int          `json:"address_id"`
    DeliverTime                     string          `json:"deliver_time"`
    Description                     string          `json:"description"`
    Sig                     string          `json:"sig"`
    PaymethodId                     int          `json:"paymethod_id"`
    Geohash                    string       `json:"geohash"`
    Latitude                   string       `json:"latitude"`
    Longitude                  string       `json:"longitude"`
    Entities                   []model.CartGroupStruct     `json:"entities" bson:"entities"`
}

// 添加顾客 订单 表项
func CustomerPlaceOrderHandler(c *gin.Context) {
    SendErrJSON := common.SendErrJSON
    var err error

    userInter, _ := c.Get("user")
    user := userInter.(model.User)

    user_id := c.Param("user_id")

    // fmt.Printf("user_id: %v\n", user_id)
    // fmt.Printf("user.OpenId: %v\n", user.OpenId)

    if user_id != user.OpenId {
        SendErrJSON("user_id 不匹配", c)
        return
    }

    cartId, err := strconv.Atoi(c.Param("cart_id"))
    if err != nil {
        SendErrJSON("错误的 cart_id", c)
        return
    }
    fmt.Printf("cartId: %v\n", cartId)

    var AddOrderData AddOrderInfo
    if err := c.ShouldBindWith(&AddOrderData, binding.JSON); err != nil {
        fmt.Printf("get addCartInfo err: %v\n", err)
        SendErrJSON("参数无效", c)
        return
    }
    // fmt.Printf("AddOrderData: %v\n", AddOrderData)

    userShipAddrId := AddOrderData.AddressId

    newOrder := model.NewOrderTable()

    // 从ids表中获取最新 order_id
    newOrderId := model.GetId("order_id")
    if newOrderId == 0 {
        SendErrJSON("获取 newOrderId 失败", c)
        return
    }

    // get a new mgo session
    mgoConn := model.GetMgoSession()
    defer mgoConn.Close()

    // 从 carts 表中获取 指定 cart
    selector := bson.M{"id": cartId}
    var CartData model.CartTable
    err = mgoConn.FindOne("carts", &CartData, selector, nil)
    if err != nil {
        fmt.Println(err)
        SendErrJSON(" cart 不存在", c)
        return
    }

    // TODO: check if there are enough food for this order
    // if yes, reserve it. Once order timeout, release this order and release the reserve.
    // if no, return "xxx is not enough"

    newOrder.Id = newOrderId
    newOrder.UniqueId = newOrderId
    newOrder.UserId = user_id
    newOrder.AddressId = userShipAddrId
    newOrder.TotalAmount = CartData.Cart.Total
    newOrder.TotalQuantity = len(AddOrderData.Entities)

    newOrder.ShopId = CartData.Cart.ShopId
    newOrder.ShopName = CartData.Cart.ShopInfo.Name
    newOrder.ShopImageUrl = CartData.Cart.ShopInfo.ImagePath

    newOrder.Basket.DeliverFee.Price = CartData.Cart.DeliverAmount

    newOrder.Basket.PackingFee.Name = CartData.Cart.Extra[0].Name
    newOrder.Basket.PackingFee.Price = CartData.Cart.Extra[0].Price
    newOrder.Basket.PackingFee.Quantity = CartData.Cart.Extra[0].Quantity
    newOrder.Basket.Groups = AddOrderData.Entities

    if err := mgoConn.Insert("orders", &newOrder); err != nil {
        fmt.Println(err)
        SendErrJSON("创建 orders 数据库失败.", c)
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "errNo": model.ErrorCode.SUCCESS,
        "status":   "1",
        "order_id":   newOrderId,
        "success":   "下单成功，请及时付款",
        "need_validation": false,
        "data": gin.H{
            "newOrder": newOrder,
        },
    })
}


// 获取顾客 指定订单（merchantOrderNo） 详细信息
func GetCustomerOrderDetail(c *gin.Context) {
    SendErrJSON := common.SendErrJSON
    var err error

    userInter, _ := c.Get("user")
    user := userInter.(model.User)

    user_id := c.Query("user_id")

    // fmt.Printf("user_id: %v\n", user_id)
    // fmt.Printf("user.OpenId: %v\n", user.OpenId)

    if user_id != user.OpenId {
        SendErrJSON("user_id 不匹配", c)
        return
    }

    orderId, err := strconv.Atoi(c.Query("merchantOrderNo"))
    if err != nil {
        SendErrJSON("错误的 merchantOrderNo", c)
        return
    }
    fmt.Printf("orderId: %v\n", orderId)

    // get a new mgo session
    mgoConn := model.GetMgoSession()
    defer mgoConn.Close()

    // 从 orders 表中获取 指定 order
    selector := bson.M{"id": orderId}
    var OrderData model.OrderTable
    err = mgoConn.FindOne("orders", &OrderData, selector, nil)
    if err != nil {
        fmt.Println(err)
        SendErrJSON(" order 不存在", c)
        return
    }

    if user_id != OrderData.UserId {
        SendErrJSON("user_id 与订单不匹配", c)
        return
    }

    // 从 customerShipAddr 表中获取 用户收货地址
    selector = bson.M{"user_id": user_id}
    var userShipAddrList model.CustomerShipAddr
    err = mgoConn.FindOne("customerShipAddr", &userShipAddrList, selector, nil)
    if err != nil {
        fmt.Println(err)
        SendErrJSON(" Customer ShipAddre 不存在", c)
        return
    }

    var userShipAddressDetail model.ShipAddr
    for _, userShipAddr := range userShipAddrList.ShipAddrArray {
        if userShipAddr.Id == OrderData.AddressId {
            userShipAddressDetail = userShipAddr
            break;
        }
    }

    c.JSON(http.StatusOK, gin.H{
        "errNo": model.ErrorCode.SUCCESS,
        "userShipAddressDetail": userShipAddressDetail,
        "OrderData": OrderData,
    })
}


// 获取顾客 所有订单 列表
func GetCustomerOrderList(c *gin.Context) {
    SendErrJSON := common.SendErrJSON
    var err error

    userInter, _ := c.Get("user")
    user := userInter.(model.User)

    user_id := c.Query("user_id")

    // fmt.Printf("user_id: %v\n", user_id)
    // fmt.Printf("user.OpenId: %v\n", user.OpenId)

    if user_id != user.OpenId {
        SendErrJSON("user_id 不匹配", c)
        return
    }

    limitStr := c.Query("limit")
    fmt.Printf("limitStr: %v\n", limitStr)
    limit, err := strconv.Atoi(limitStr)
    if err != nil {
        SendErrJSON("错误的 limit", c)
        return
    }
    fmt.Printf("limit: %v\n", limit)

    offset, err := strconv.Atoi(c.Query("offset"))
    if err != nil {
        SendErrJSON("错误的 offset", c)
        return
    }
    fmt.Printf("offset: %v\n", offset)

    // get a new mgo session
    mgoConn := model.GetMgoSession()
    defer mgoConn.Close()

    selector := bson.M{"user_id": user_id}

    sortBy := []string{"-created_date"}

    var OrderList []model.OrderTable
    // 从 orders 表中获取 所有用户 orders
    err = mgoConn.FindAll("orders", &OrderList, selector, nil, offset, limit, sortBy...)
    if err != nil {
        fmt.Println(err)
        SendErrJSON("error", c)
        return
    }

    userOrderRes := OrderList[:0]
    for _, userOrder := range OrderList {
        if (time.Now().Unix() - userOrder.OrderTime) < 900000 {
            userOrder.StatusBar.Title = "等待支付"
        } else {
            userOrder.StatusBar.Title = "支付超时"
        }
        userOrder.TimePass = (time.Now().Unix() - userOrder.OrderTime)/1000

        selector := bson.M{"id": userOrder.Id}
        update := bson.M{"$set": bson.M{
            "time_pass":    userOrder.TimePass,
            "status_bar": bson.M{
                "title": userOrder.StatusBar.Title},
            }}
        err := mgoConn.UpdateOne("orders", selector, update)
        // if err != nil && err != mgo.ErrNotFound {
        if err != nil {
            fmt.Println(err)
            SendErrJSON("更新 user 订单时间 失败.", c)
        }
        userOrderRes = append(userOrderRes, userOrder)
    }

    c.JSON(http.StatusOK, gin.H{
        "errNo": model.ErrorCode.SUCCESS,
        "userOrderList": userOrderRes,
    })
}
