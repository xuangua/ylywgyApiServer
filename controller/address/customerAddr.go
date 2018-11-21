package address

import (
    "fmt"
    "net/http"
    "strings"
    // "strconv"
    // "unicode/utf8"

    "github.com/gin-gonic/gin"
    // "github.com/gin-gonic/gin/binding"
    "github.com/globalsign/mgo"
    "github.com/globalsign/mgo/bson"
    "github.com/xuangua/ylywgyApiServer/controller/common"
    "github.com/xuangua/ylywgyApiServer/model"
    "github.com/xuangua/ylywgyApiServer/utils"
)

// 根据客户坐标，获取QQ LBS pois 的详细地址
func GetQQPoisAddress(c *gin.Context) {
    SendErrJSON := common.SendErrJSON

    geoHash := c.Param("geoHash")
    if geoHash == "" || false == strings.Contains(geoHash, ",") {
        SendErrJSON("错误的 geoHash", c)
        return
    }
    fmt.Printf("geoHash: %v\n", geoHash)

    // var response []byte
    var addrResult utils.AddressResult
    response, addrResult, err := utils.GetQQPois(geoHash)
    if err != nil {
        SendErrJSON("获取 pois 错误", c)
        fmt.Println(err)
        return
    }
    fmt.Printf("GetQQPois() addrResult: %v\n", addrResult)

    // TODO : convert addrResult to 省/市/区/学校名/校区名
    // 1. 检查 addrResult.result.formatted_addresses.recommend 是否包含 “大学/学院” 关键词
    //    是则 go 3
    //    否则 进入2
    // 2. 查看poisDetail是否有 “title”包含 “大学/学院” 关键词, 且 "category": "教育学校:大学",
    //    是则 go 3
    //    否则 go 5
    // 3. 使用该模糊 校名去查找 “系统已有店铺的” 学校及地址。 TBD---“使用测距功能获取校区名（校区名在学校名后面的括号里面）。”
    //    找到匹配的准确学校名/校区名：  go 4
    //    未找到：go 5
    // 4. 使用准确的 学校名/校区名 去标记该用户的shop地址并显示，并去校区的shopAddress里面获取数据shopList
    // 5. 显示用户的实际actually地址, 并标记shop地址为null （手机端提示用户选择新的收货地址）。
    //    TBD---“使用测距获取到已有店铺的学校的距离，并选择最近的学校作为shopAddress“。并去校区的shopAddress里面获取数据shopList

    var province = ""
    var city = ""
    var district = ""
    var realAddrName string
    var shopAddrName = ""
    var roughSchoolName = ""
    var distance = 0.0
    var decidedSchoolName = ""
    var decidedSchoolCampusName = ""
    var decidedSchoolData = model.SchoolData{}

    // 1
    realAddrName = addrResult.ResultDetail.FormattedAddresses.Recommend
    if strings.Contains(realAddrName, "大学") || strings.Contains(realAddrName, "学院") {
        roughSchoolName = realAddrName
    } else {
        for _, piosDetail := range addrResult.ResultDetail.PiosDetail {
            if strings.Contains(piosDetail.Title, "大学") || strings.Contains(piosDetail.Title, "学院") {
                roughSchoolName = piosDetail.Title
                distance = piosDetail.Distance
                break;
            }
        }
    }
    province = addrResult.ResultDetail.AddressesComponent.Province;
    city = addrResult.ResultDetail.AddressesComponent.City;
    district = addrResult.ResultDetail.AddressesComponent.District;

    // get a new mgo session
    mgoConn := model.GetMgoSession()
    defer mgoConn.Close()

    // 3
    if 0 != strings.Compare(roughSchoolName, "") && distance < 1000.00 /*within 1000 meter */ {
        var selectedSchoolList []model.SchoolData
        selector := bson.M{"province": province, "city": city, "district": district,}
        // err = mgoConn.WithLog().FindAll("schoolList", &selectedSchoolList, selector, nil, 0, 0, "$natural")
        err = mgoConn.FindAll("schoolList", &selectedSchoolList, selector, nil, 0, 0, "$natural")
        if err != nil && err != mgo.ErrNotFound {
            fmt.Println(err)
            SendErrJSON("pois 匹配数据库失败.", c)
        }

        for _, tmpSchoolData := range selectedSchoolList {
            if true == tmpSchoolData.IsActive && 
               strings.Contains(roughSchoolName, tmpSchoolData.SchoolName) && 
               strings.Contains(roughSchoolName, tmpSchoolData.SchoolCampusName) {
                decidedSchoolData = tmpSchoolData
                decidedSchoolName = tmpSchoolData.SchoolName
                if 0==strings.Compare(decidedSchoolCampusName,"") && 
                    0!=strings.Compare(tmpSchoolData.SchoolCampusName,"") {
                    // 匹配到精确校区名，退出
                    decidedSchoolData = tmpSchoolData
                    decidedSchoolCampusName = tmpSchoolData.SchoolCampusName
                    break;
                } // 否则，继续匹配可能的校区
            }
        }
    }

    // 4
    if 0!=strings.Compare(decidedSchoolName,"") {
        shopAddrName = decidedSchoolName
        if 0!=strings.Compare(decidedSchoolCampusName,"") {
            shopAddrName = decidedSchoolName + "(" + decidedSchoolCampusName + ")"
        }
    }

    // 5
    // shopAddrName == ""
    // decidedSchoolData == {}

    c.JSON(http.StatusOK, gin.H{
        "errNo": model.ErrorCode.SUCCESS,
        "msg":   "success",
        "realAddrName":   realAddrName,
        "roughSchoolName":   roughSchoolName,
        "shopAddrName":   shopAddrName,
        "decidedSchoolData":   decidedSchoolData,
        "realaddressDetail": addrResult.ResultDetail,
        "pios":  response})
}


type SupportedSchoolString struct {
    Text        string      `json:"text"` 
    Value       int         `json:"value"`
}

// 获取 以支持学校列表
func GetSupportedSchoolList(c *gin.Context) {
    SendErrJSON := common.SendErrJSON

    var province string
    if province = c.Query("province"); province == "" {
        SendErrJSON("参数 province 无效", c)
        return
    }

    var city string
    if city = c.Query("city"); city == "" {
        SendErrJSON("参数 city 无效", c)
        return
    }

    // get a new mgo session
    mgoConn := model.GetMgoSession()
    defer mgoConn.Close()

    var selectedSchoolList []model.SchoolData
    selector := bson.M{"province": province, "city": city,}
    err := mgoConn.FindAll("schoolList", &selectedSchoolList, selector, nil, 0, 0, "$natural")
    if err != nil && err != mgo.ErrNotFound {
        fmt.Println(err)
        SendErrJSON("GetSupportedSchoolList 获取 以支持学校列表数据库失败.", c)
    }

    var supportedSchoolStringList []SupportedSchoolString
    for _, tmpSchoolData := range selectedSchoolList {
        // TODO: check school is actived
        // if true == tmpSchoolData.IsActive && 
            var tmpSchoolName SupportedSchoolString
            tmpSchoolName.Text = tmpSchoolData.SchoolName + tmpSchoolData.SchoolCampusName
            tmpSchoolName.Value = tmpSchoolData.Id
            supportedSchoolStringList = append(supportedSchoolStringList, tmpSchoolName)
        // }
    }

    c.JSON(http.StatusOK, gin.H{
        "errNo": model.ErrorCode.SUCCESS,
        "supportedSchoolList":   supportedSchoolStringList,
        "supportedSchoolDetailList":   selectedSchoolList})
}

// 根据客户坐标，获取QQ LBS建议的详细地址
func GetQQSuggestedAddress(c *gin.Context) {
    SendErrJSON := common.SendErrJSON

    geoHash := c.Param("geoHash")
    if geoHash == "" || false == strings.Contains(geoHash, ",") {
        SendErrJSON("错误的 geoHash", c)
        return
    }
    fmt.Printf("geoHash: %v\n", geoHash)

    c.JSON(http.StatusOK, gin.H{
        "errNo": model.ErrorCode.SUCCESS,
        "msg":   "success",
        // "data": gin.H{
        //  "id": book.ID,
        // },
    })
}

func isContainsSchoolCampusName(schoolData model.SchoolData) (string, string, bool) {
    isContains :=  false
    schoolName :=  ""
    schoolCampusName :=  ""

    if true==strings.Contains(schoolData.Title, "(") && 
        true==strings.Contains(schoolData.Title, ")") && 
        1==strings.Count(schoolData.Title, "(") && 
        1==strings.Count(schoolData.Title, ")") && 
        (len(schoolData.Title)-1)==strings.Index(schoolData.Title, ")") {

        isContains =  true
        s1 := strings.Split(schoolData.Title, "(")
        schoolName, schoolCampusName = s1[0],s1[1]

        s2 := strings.Split(schoolCampusName, ")")
        schoolCampusName = s2[0]
    }
    return schoolName, schoolCampusName, isContains
}