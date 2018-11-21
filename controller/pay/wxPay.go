package pay

import (
    "fmt"
    "net/http"
    "strconv"
    "time"
    "io/ioutil"
    "encoding/xml"
    "strings"
    // "unicode/utf8"
    "encoding/hex"
    "errors"
    "crypto/md5"
    "sort"
    "bytes"

    "github.com/gin-gonic/gin"
    // "github.com/gin-gonic/gin/binding"
    // "github.com/globalsign/mgo"
    "github.com/globalsign/mgo/bson"
    _wxpay "github.com/silenceper/wechat/pay"
    "github.com/xuangua/ylywgyApiServer/controller/common"
    "github.com/xuangua/ylywgyApiServer/model"
    "github.com/xuangua/ylywgyApiServer/utils"
    "github.com/zhufuyi/logger"
)

// 获取顾客 指定订单（merchantOrderNo） 微信支付参数信息
func GetWxPayParameters(c *gin.Context) {
    SendErrJSON := common.SendErrJSON
    var err error

    userInter, _ := c.Get("user")
    user := userInter.(model.User)

    user_id := c.Query("user_id")

    fmt.Printf("user_id: %v\n", user_id)
    fmt.Printf("user.OpenId: %v\n", user.OpenId)

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

    newOrderPay := model.NewOrderPayTable()

    // 从ids表中获取最新 order_id
    newOrderPayId := model.GetId("order_pay_id")
    if newOrderPayId == 0 {
        SendErrJSON("获取 newOrderPayId 失败", c)
        return
    }

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

    wechatPay := utils.WechatInst.GetPay()
    var wxPayParams _wxpay.Params
    // wxPayParams.TotalFee = strconv.Itoa(int(OrderData.TotalAmount*100))
    wxPayParams.TotalFee = strconv.Itoa(int(1))
    wxPayParams.CreateIP = c.ClientIP();
    body := user.Nickname + "-" + OrderData.ShopName + "-" + wxPayParams.TotalFee //
    if len(body) >= 128 {
        body = body[0:127]
    }
    wxPayParams.Body = body;
    wxPayParams.OutTradeNo = time.Now().Format("20060102150405") + "-" + strconv.Itoa(OrderData.Id);
    wxPayParams.OpenID = user_id;
    // fmt.Printf("wxPayParams: %v\n", wxPayParams)

    wxPayJsdkConfig, err := wechatPay.GetWxPayJsdkConfigParams(&wxPayParams)
    // fmt.Printf("wxPayJsdkConfig: %v\n", wxPayJsdkConfig)

    newOrderPay.Id = newOrderPayId
    newOrderPay.OutTradeNo = wxPayParams.OutTradeNo
    newOrderPay.OrderId = OrderData.Id
    newOrderPay.ShopId = OrderData.ShopId
    newOrderPay.ShopName = OrderData.ShopName
    newOrderPay.TotalFee = OrderData.TotalAmount
    newOrderPay.UserId = user_id

    if err := mgoConn.Insert("ordersPay", &newOrderPay); err != nil {
        fmt.Println(err)
        SendErrJSON("创建 orderPay 数据库失败.", c)
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "errNo": model.ErrorCode.SUCCESS,
        "wxPayJsdkConfig": wxPayJsdkConfig,
    })
}

// 微信支付 成功后的更新 orderPay 表
func WxPaySucHandler(mr WXPayNotifyReq) {
    // get a new mgo session
    mgoConn := model.GetMgoSession()
    defer mgoConn.Close()

    // 从 orders 表中获取 指定 order
    selector := bson.M{"out_trade_no": mr.Out_trade_no}
    var OrderPayData model.OrderPayTable
    err := mgoConn.FindOne("ordersPay", &OrderPayData, selector, nil)
    if err != nil {
        fmt.Println(err)
        fmt.Printf("获取 orderPay 数据库失败. %v\n", mr.Out_trade_no)
        return
    }

    selector = bson.M{"out_trade_no": mr.Out_trade_no}
    update := bson.M{"$set": bson.M{
        "return_code":    mr.Return_code,
        "return_msg":    mr.Return_msg,
        "result_code":    mr.Result_code,
        "trade_type":    mr.Trade_type,
        "bank_type":    mr.Bank_type,
        "transaction_id":    mr.Transaction_id,
        "time_end":    mr.Time_end,
        }}
    err = mgoConn.UpdateOne("ordersPay", selector, update)
    // if err != nil && err != mgo.ErrNotFound {
    if err != nil {
        fmt.Println(err)
        fmt.Printf("更新 orderPay 数据库失败. %v\n", mr.Out_trade_no)
    }
}

// 微信支付 成功后的回调
func WxPaySucNotify(c *gin.Context) {
    // SendErrJSON := common.SendErrJSON
    // var err error

    Out_trade_no,Result_code := WxpayCallback(c.Writer, c.Request, WxPaySucHandler, utils.WechatInst.Context.PayKey)
    // fmt.Printf("Out_trade_no: %v\n", Out_trade_no)
    // fmt.Printf("Result_code: %v\n", Result_code)
    logger.WithFields(
        // logger.Err(err),
        logger.String("Out_trade_no", Out_trade_no),
        logger.String("Result_code", Result_code),
        // logger.Any("content", a),
    ).Info("mongodb insert failed!")
}

type WXPayNotifyReq struct {
    Return_code    string `xml:"return_code"`
    Return_msg     string `xml:"return_msg"`
    Appid          string `xml:"appid"`
    Mch_id         string `xml:"mch_id"`
    Nonce          string `xml:"nonce_str"`
    Sign           string `xml:"sign"`
    Result_code    string `xml:"result_code"`
    Openid         string `xml:"openid"`
    Is_subscribe   string `xml:"is_subscribe"`
    Trade_type     string `xml:"trade_type"`
    Bank_type      string `xml:"bank_type"`
    Total_fee      int    `xml:"total_fee"`
    Fee_type       string `xml:"fee_type"`
    Cash_fee       int    `xml:"cash_fee"`
    Cash_fee_Type  string `xml:"cash_fee_type"`
    Transaction_id string `xml:"transaction_id"`
    Out_trade_no   string `xml:"out_trade_no"`
    Attach         string `xml:"attach"`
    Time_end       string `xml:"time_end"`
}

type WXPayNotifyResp struct {
    Return_code string `xml:"return_code"`
    Return_msg  string `xml:"return_msg"`
}

type GetOrderDetail struct{
    Appid string `xml:"appid"`
    Mch_id string `xml:"mch_id"`
    Out_trade_no string `xml:"out_trade_no"`
    Noncestr string `xml:"nonce_str"`
    Sign string `xml:"sign"`
    Sign_type string `xml:"sign_type"`
}

//具体的微信支付回调函数的范例
func WxpayCallback(w http.ResponseWriter, r *http.Request, f func(WXPayNotifyReq), key string)(string,string) {
    // body
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        fmt.Println("读取http body失败，原因!", err)
        http.Error(w.(http.ResponseWriter), http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
        return "","FAIL"
    }
    defer r.Body.Close()

    fmt.Println("微信支付异步通知，HTTP Body:", "成功")
    var mr WXPayNotifyReq
    err = xml.Unmarshal(body, &mr)
    if err != nil {
        fmt.Println("解析HTTP Body格式到xml失败，原因!", err)
        http.Error(w.(http.ResponseWriter), http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
        return "","FAIL"
    }

    var reqMap map[string]interface{}
    reqMap = make(map[string]interface{}, 0)

    reqMap["return_code"] = mr.Return_code
    reqMap["return_msg"] = mr.Return_msg
    reqMap["appid"] = mr.Appid
    reqMap["mch_id"] = mr.Mch_id
    reqMap["nonce_str"] = mr.Nonce
    reqMap["result_code"] = mr.Result_code
    reqMap["openid"] = mr.Openid
    reqMap["is_subscribe"] = mr.Is_subscribe
    reqMap["trade_type"] = mr.Trade_type
    reqMap["bank_type"] = mr.Bank_type
    reqMap["total_fee"] = mr.Total_fee
    reqMap["fee_type"] = mr.Fee_type
    reqMap["cash_fee"] = mr.Cash_fee
    reqMap["cash_fee_type"] = mr.Cash_fee_Type
    reqMap["transaction_id"] = mr.Transaction_id
    reqMap["out_trade_no"] = mr.Out_trade_no
    reqMap["attach"] = mr.Attach
    reqMap["time_end"] = mr.Time_end

    var resp WXPayNotifyResp
    //进行签名校验
    if wxpayVerifySign(reqMap, mr.Sign,key) {
        if f != nil {
            f(mr)
        }
        //这里就可以更新我们的后台数据库了，其他业务逻辑同理。
        resp.Return_code = "SUCCESS"
        resp.Return_msg = "OK"
    } else {
        resp.Return_code = "FAIL"
        resp.Return_msg = "failed to verify sign, please retry!"
    }

    //结果返回，微信要求如果成功需要返回return_code "SUCCESS"
    bytes, _err := xml.Marshal(resp)
    strResp := strings.Replace(string(bytes), "WXPayNotifyResp", "xml", -1)
    if _err != nil {
        fmt.Println("xml编码失败，原因：", _err)
        http.Error(w.(http.ResponseWriter), http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
        return "","FAIL"
    }
    w.(http.ResponseWriter).WriteHeader(http.StatusOK)
    fmt.Fprint(w.(http.ResponseWriter), strResp)
    return mr.Out_trade_no,mr.Result_code
}


//微信支付 下单签名
func wxpayCalcSign(mReq map[string]interface{}, key string) string {

    //fmt.Println("========STEP 1, 对key进行升序排序.========")
    //fmt.Println("微信支付签名计算, API KEY:", key)
    //STEP 1, 对key进行升序排序.
    sorted_keys := make([]string, 0)
    for k, _ := range mReq {
        sorted_keys = append(sorted_keys, k)
    }

    sort.Strings(sorted_keys)

    //fmt.Println("========STEP2, 对key=value的键值对用&连接起来，略过空值========")
    //STEP2, 对key=value的键值对用&连接起来，略过空值
    var signStrings string
    for _, k := range sorted_keys {
        //fmt.Printf("k=%v, v=%v\n", k, mReq[k])
        value := fmt.Sprintf("%v", mReq[k])
        if value != "" {
            signStrings = signStrings + k + "=" + value + "&"
        }
    }

    //fmt.Println("========STEP3, 在键值对的最后加上key=API_KEY========")
    //STEP3, 在键值对的最后加上key=API_KEY
    if key != "" {
        signStrings = signStrings + "key=" + key
    }

    //fmt.Println("========STEP4, 进行MD5签名并且将所有字符转为大写.========")
    //STEP4, 进行MD5签名并且将所有字符转为大写.
    md5Ctx := md5.New()
    md5Ctx.Write([]byte(signStrings))
    cipherStr := md5Ctx.Sum(nil)
    upperSign := strings.ToUpper(hex.EncodeToString(cipherStr))

    return upperSign
}

//微信支付签名验证函数
func wxpayVerifySign(needVerifyM map[string]interface{}, sign string,key string) bool {
    signCalc := wxpayCalcSign(needVerifyM , key)
    if sign == signCalc {
        fmt.Println("签名校验通过!")
        return true
    }

    fmt.Println("签名校验失败!")
    return false
}

//查询订单
func GetWeixinOrderInfo(appid string , mch_id string,out_trade_no string,noncestr string,_type string,key string)(*WXPayNotifyReq,error){
    sendData := GetOrderDetail{
        Appid:appid,
        Mch_id:mch_id,
        Out_trade_no:out_trade_no,
        Noncestr:noncestr,
        Sign_type:_type,
    }
    var m map[string]interface{}
    m = make(map[string]interface{}, 0)
    m["appid"] = sendData.Appid
    m["mch_id"] = sendData.Mch_id
    m["nonce_str"] = sendData.Noncestr
    m["out_trade_no"] = sendData.Out_trade_no
    m["sign_type"] = sendData.Sign_type
    sendData.Sign = wxpayCalcSign(m, key)
    bytes_req, err := xml.Marshal(sendData)     
    str_req := strings.Replace(string(bytes_req), "UnifyOrderReq", "xml", -1)
    //fmt.Println("转换为xml--------", str_req)
    bytes_req = []byte(str_req)

     //发送unified order请求.
     req, err := http.NewRequest("POST", "https://api.mch.weixin.qq.com/pay/orderquery", bytes.NewReader(bytes_req))
     if err != nil {
         fmt.Println("New Http Request发生错误，原因:", err)
         return nil,errors.New("Http Request发生错误")

     }
     req.Header.Set("Accept", "application/xml")
     //这里的http header的设置是必须设置的.
     req.Header.Set("Content-Type", "application/xml;charset=utf-8")

     client := http.Client{}
     resp, _err := client.Do(req)
     if _err != nil {
         fmt.Println("请求微信支付统一下单接口发送错误, 原因:", _err)
         return nil,errors.New("请求微信支付统一下单接口发送错误")
     }
    respBytes, err := ioutil.ReadAll(resp.Body)
     if err != nil {
         fmt.Println("解析返回body错误", err)
         return nil,errors.New("解析返回body错误")
     }
     xmlResp := WXPayNotifyReq{}
     _err = xml.Unmarshal(respBytes, &xmlResp)
     //处理return code.
     if xmlResp.Return_code == "FAIL" {
         fmt.Println("微信支付查询订单不成功，原因:", xmlResp.Return_msg, " str_req-->", str_req)
         return nil,errors.New("不成功失败原因:"+xmlResp.Return_msg)
     }else{
        return &xmlResp,nil
     }
}

    // /***
    //  * 付款成功微信异步回调处理
    //  *
    //  * @throws IOException
    //  * @throws JDOMException
    //  */
    // @Action("notifysuccess")
    // public void notifysuccess() throws IOException, JDOMException {
    //     HttpServletRequest request=Struts2Utils.getRequest();
    //     HttpServletResponse response=Struts2Utils.getResponse();
    //     log.info("微信支付成功调用回调URL");
    //     InputStream inStream = request.getInputStream();
    //     ByteArrayOutputStream outSteam = new ByteArrayOutputStream();
    //     byte[] buffer = new byte[1024];
    //     int len = 0;
    //     while ((len = inStream.read(buffer)) != -1) {
    //         outSteam.write(buffer, 0, len);
    //     }
    //     log.info("~~~~~~~~~~~~~~~~付款成功~~~~~~~~~");
    //     outSteam.close();
    //     inStream.close();

    //     /** 支付成功后，微信回调返回的信息 */
    //     String result = new String(outSteam.toByteArray(), "utf-8");
    //     log.info("微信返回的订单支付信息:" + result);
    //     Map<Object, Object> map = XMLUtil.doXMLParse(result);

    //     // 用于验签
    //     SortedMap<Object, Object> parameters = new TreeMap<Object, Object>();
    //     for (Object keyValue : map.keySet()) {
    //         /** 输出返回的订单支付信息 */
    //         log.info(keyValue + "=" + map.get(keyValue));
    //         if (!"sign".equals(keyValue)) {
    //             parameters.put(keyValue, map.get(keyValue));
    //         }
    //     }
    //     if (map.get("result_code").toString().equalsIgnoreCase("SUCCESS")) {
    //         // 先进行校验，是否是微信服务器返回的信息
    //         String checkSign = createSign("UTF-8", parameters);
    //         log.info("对服务器返回的结果进行签名：" + checkSign);
    //         log.info("服务器返回的结果签名：" + map.get("sign"));
    //         if (checkSign.equals(map.get("sign"))) {// 如果签名和服务器返回的签名一致，说明数据没有被篡改过
    //             log.info("签名校验成功，信息合法，未被篡改过");


    //             //告诉微信服务器，我收到信息了，不要再调用回调方法了

    //             /**如果不返回SUCCESS的信息给微信服务器，则微信服务器会在一定时间内，多次调用该回调方法，如果最终还未收到回馈，微信默认该订单支付失败*/
    //             /** 微信默认会调用8次该回调地址 */
    //             /*
    //              * 【需要注意】：
    //              *      后期很多朋友都反应说，微信会一直请求这个回调地址，也都是使用的是 response.getWriter().write(setXML("SUCCESS", ""));
    //              *      百思不得其解，最终发现处理办法。其实不能直接使用response.getWriter()返回结果，这样微信是接收不到的。
    //              *      只能使用OutputStream流的方式返回结果给微信。
    //              *      切记！！！！
    //              * */
    //             OutputStream outputStream = null;
    //             try {
    //                 outputStream = response.getOutputStream();
    //                 outputStream.flush();
    //                 outputStream.write(setXML("SUCCESS", "").getBytes());
    //             } catch (Exception e) {
    //                 e.printStackTrace();
    //             }finally {
    //                 try {
    //                     outputStream.close();
    //                 } catch (Exception e) {
    //                     e.printStackTrace();
    //                 }
    //             }

    //             //TODO 增加用户投票
    //             String orderId = map.get("out_trade_no").toString();
    //             String openid = map.get("openid").toString();
    //             String appid = map.get("appid").toString();
    //             String mchid = map.get("mch_id").toString();
    //             String transactionId = map.get("transaction_id").toString();
    //             String transactionTimeEnd = map.get("time_end").toString();
    //             // `is_subscribe` varchar(32) DEFAULT NULL,
    //             String bankType = map.get("bank_type").toString();
    //             String tradeType = map.get("trade_type").toString();

    //             ToupiaoOrder toupiaoOrder = toupiaoOrderManager.getToupiaoOrderWithOpenIdandOrderId(openid, orderId);
    //             if (toupiaoOrder != null) {
    //                 toupiaoOrder.setAppid(appid);
    //                 toupiaoOrder.setMchid(mchid);
    //                 toupiaoOrder.setTransactionId(transactionId);
    //                 toupiaoOrder.setTransactionTimeEnd(transactionTimeEnd);
    //                 toupiaoOrder.setResultCode("SUCCESS");
    //                 toupiaoOrder.setBankType("bankType");
    //                 toupiaoOrder.setTradeType("tradeType");
    //                 toupiaoOrder.setOrderState(3);
    //                 toupiaoOrderManager.save(toupiaoOrder);
    //                 try {
    //                     toupiaoDetailManager.zuanshiToupiao(toupiaoOrder.getDirId(), toupiaoOrder.getXuanshouId(), toupiaoOrder.getZuanshiNumber());
    //                 } catch (Exception e) {
    //                     e.printStackTrace();
    //                 }
    //             } else {
    //                 log.error("cannot find ToupiaoOrder.");
    //             }
    //         }
    //     }
    // }
