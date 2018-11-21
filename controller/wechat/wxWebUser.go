package wechat

import (
    // "encoding/json"
    "fmt"
    "net/http"
    // "strings"
    // "time"

    "github.com/dgrijalva/jwt-go"
    // "github.com/gin-contrib/sessions"
    "github.com/gin-gonic/gin"
    "github.com/gin-gonic/gin/binding"
    "github.com/satori/go.uuid"
    "github.com/silenceper/wechat/oauth"
    "github.com/xuangua/ylywgyApiServer/controller/common"
    "github.com/xuangua/ylywgyApiServer/model"
    "github.com/xuangua/ylywgyApiServer/utils"
    "github.com/xuangua/ylywgyApiServer/config"
)

// func UserLoginWithWx(ctx *gin.Context) {
//     SendErrJSON := common.SendErrJSON
//     var err error

//     wechatOauth := utils.WechatInst.GetOauth()

//     code := ctx.Query("code")
//     if code == "" {
//         location, e := wechatOauth.GetRedirectURL("http://wemall.wxapp.xuangua.xyz/api/xuewu/wechat/userLoginWithWx", "snsapi_userinfo", "test")
//         if e != nil {
//             fmt.Println(e.Error())
//         }
//         ctx.Redirect(302, location)

//         // wechatOauth.Redirect(ctx.Writer, ctx.Request, "http://wemall.wxapp.xuangua.xyz/api/xuewu/wechat/userLoginWithWx", "snsapi_userinfo", "test")
//         // SendErrJSON("code 不能为空", ctx)
//         return
//     }

//     var resATResp oauth.ResAccessToken
//     resATResp, err = wechatOauth.GetUserAccessToken(code)
//     if err != nil {
//         fmt.Println(err.Error())
//         SendErrJSON("get User AccessToken error", ctx)
//         return
//     }

//     var wxUserInfo oauth.UserInfo
//     wxUserInfo, err = wechatOauth.GetUserInfo(resATResp.AccessToken, resATResp.OpenID)
//     if err != nil {
//         fmt.Println(err.Error())
//         SendErrJSON("get UserInfo error", ctx)
//         return
//     }

//     LoginWithUserInfo(wxUserInfo, ctx)

//     ctx.JSON(http.StatusOK, gin.H{
//         "errNo": model.ErrorCode.SUCCESS,
//         "msg":   "success",
//         "wxUserInfo":  wxUserInfo,
//     }) 
// }

// We Oauth 获取auth Config
func GetWxWebUserInfo(ctx *gin.Context) {
    SendErrJSON := common.SendErrJSON
    var err error

    code := ctx.Query("code")
    if code == "" {
        SendErrJSON("code 不能为空", ctx)
        return
    }

    wechatOauth := utils.WechatInst.GetOauth()
    var resATResp oauth.ResAccessToken
    resATResp, err = wechatOauth.GetUserAccessToken(code)
    if err != nil {
        fmt.Println(err.Error())
        SendErrJSON("get User AccessToken error", ctx)
        return
    }

    var wxUserInfo oauth.UserInfo
    wxUserInfo, err = wechatOauth.GetUserInfo(resATResp.AccessToken, resATResp.OpenID)
    if err != nil {
        fmt.Println(err.Error())
        SendErrJSON("get UserInfo error", ctx)
        return
    }

    // LoginWithUserInfo(wxUserInfo, ctx)

    ctx.JSON(http.StatusOK, gin.H{
        "errNo": model.ErrorCode.SUCCESS,
        "msg":   "success",
        "wxUserInfo":  wxUserInfo,
    })
    // ctx.JSON(http.StatusOK, wxUserInfo)
}

type LoginUserInfo struct {
    OpenID     string   `json:"openid"`
    Nickname   string   `json:"nickname"`
    Sex        int32    `json:"sex"`
    Province   string   `json:"province"`
    City       string   `json:"city"`
    Country    string   `json:"country"`
    HeadImgURL string   `json:"headimgurl"`
    // Privilege  []string `json:"privilege"`
    // Unionid    string   `json:"unionid"`
}

// Customer LoginWithUserInfo 设置 web微信用户 登陆信息与cookie
func LoginWithUserInfo(ctx *gin.Context) {
    SendErrJSON := common.SendErrJSON

    var WxUserInfo oauth.UserInfo
    if err := ctx.ShouldBindWith(&WxUserInfo, binding.JSON); err != nil {
        SendErrJSON("invalid binding wxUserInfo Json data", ctx)
        return
    }

    openID := WxUserInfo.OpenID
    // (1) setup wx web user info in mysql DB
    var sql string
    var user model.User
    // var userOpenId uint
    sql = "open_id = ?"
    if err := model.DB.Where(sql, openID).First(&user).Error; err != nil {
        user.Uid = uuid.NewV4().String()
        user.OpenId = openID
        user.Nickname = WxUserInfo.Nickname
        user.Province = WxUserInfo.Province
        user.City = WxUserInfo.City
        user.Country = WxUserInfo.Country
        user.HeadImgURL = WxUserInfo.HeadImgURL
        user.Unionid = WxUserInfo.Unionid
        user.Role = model.UserRoleNormal
        user.Status = model.UserStatusActived
        user.ShipAddr = ""
        user.Points = 0

        if err := model.DB.Create(&user).Error; err != nil {
            fmt.Println(err.Error())
            SendErrJSON("创建用户数据库表项错误", ctx)
            return
        }
        // userOpenId = newUser.OpenId
    } /* else {
        userOpenId = user.OpenId
    }*/

    // (2) setup cookie
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "id": user.ID,
    })
    tokenString, err := token.SignedString([]byte(config.ServerConfig.TokenSecret))
    if err != nil {
        fmt.Println(err.Error())
        SendErrJSON("内部 存储用户Token 错误", ctx)
        return
    }
    ctx.SetCookie("token", tokenString, config.ServerConfig.TokenMaxAge, "/", "", false, true)

    // (3) setup reis
    if err := model.UserToRedis(user); err != nil {
        fmt.Println(err.Error())
        SendErrJSON("内部 存储用户信息到 Redis 错误.", ctx)
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "errNo": model.ErrorCode.SUCCESS,
        "msg":   "success",
        "data": gin.H{
            "token": tokenString,
            "user":  user,
        },
    })

    // sid := uuid.NewV4().String()
    // openID := wxUserInfo.openID
    // data := map[string]interface{
    //     "openid": openID
    // }

    // data["expireTime"].(string) = xxx

    // session := sessions.Default(ctx)
    // session.Set(sid, data)
    // // session.Set(openID, sessionKey)
    // session.Save()
    // fmt.Println("wx user access [weAppOpenID: %s ]has been saved to session for cookie\n", openID)


    // if err := model.UserToRedisWithOpenId(user); err != nil {
    //     SendErrJSON("内部错误.", ctx)
    //     return
    // }

    // token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    //     "id": openID,
    // })
    // tokenString, err := token.SignedString([]byte(config.ServerConfig.TokenSecret))
    // if err != nil {
    //     fmt.Println(err.Error())
    //     SendErrJSON("内部错误", ctx)
    //     return
    // }

    // userInter, exists := ctx.Get("user")
    // var user model.User
    // if exists {
    //     user = userInter.(model.User)
    // }

    // sidInter, exists := ctx.Get("sid")
    // var sid string
    // if exists {
    //     sid = sidInter.(string)
    // }

    // sessionDataInter, exists := ctx.Get("sessionData")
    // var wxSessionData map[string]interface{}
    // // var openID string
    // var sessionKey string
    // if exists {
    //     wxSessionData = sessionDataInter.(map[string]interface{})
    //     sessionKey = wxSessionData["session_key"].(string)
    // }
    // // var openID string
    // // var sessionKey string
    // // session := sessions.Default(ctx)
    // // userWxSessionData := session.Get(sid)
    // // wxSessionData, ok := userWxSessionData.(map[string]interface{})
    // // if ok {
    // //  /* act on str */
    // //  openID = wxSessionData["openid"].(string)
    // //  sessionKey = wxSessionData["session_key"].(string)
    // //  // expireTime = data["session_key"].(string)
    // //  if (openID == "") || (sessionKey == "") {
    // //      SendErrJSON("session error", ctx)
    // //      return
    // //  }
    // // } else {
    // //  /* not string */
    // //  SendErrJSON("session error", ctx)
    // //  return
    // // }

    // var weAppUser EncryptedUser
    // if err := ctx.ShouldBindWith(&weAppUser, binding.JSON); err != nil {
    //     SendErrJSON("invalid binding Json data", ctx)
    //     return
    // }

    // userInfoStr, err := utils.DecodeWeAppUserInfo(weAppUser.EncryptedData, sessionKey, weAppUser.IV)
    // if err != nil {
    //     fmt.Println(err.Error())
    //     SendErrJSON("error", ctx)
    //     return
    // }

    // var wxAppUser model.WeAppUser
    // if err := json.Unmarshal([]byte(userInfoStr), &wxAppUser); err != nil {
    //     SendErrJSON("error", ctx)
    //     return
    // }

    // updatedData := map[string]interface{}{
    //     "name":       wxAppUser.Nickname,
    //     "sex":        wxAppUser.Gender,
    //     "avatar_url": wxAppUser.AvatarURL,
    // }
    // if err := model.DB.Model(&user).Update(updatedData).Error; err != nil {
    //     SendErrJSON("update user wx data to database error", ctx)
    //     return
    // }
    // user.Name = wxAppUser.Nickname
    // user.Sex = wxAppUser.Gender
    // user.AvatarURL = wxAppUser.AvatarURL

    // if model.UserToRedisWithOpenId(user) != nil {
    //     SendErrJSON("update user wx data to redis error", ctx)
    //     return
    // }
    // // var sql string
    // // var user model.User
    // // sql = "open_id = ?"
    // // if err := model.DB.Where(sql, openID).First(&user).Error; err != nil {
    // //  SendErrJSON("can't find user with open_id " + openID, ctx)
    // //  return
    // // } else {
    // //  user.Name = wxAppUser.Nickname
    // //  user.Sex = wxAppUser.Gender
    // //  user.AvatarURL = wxAppUser.AvatarURL
    // // }

    // // session.Set("weAppUser", wxAppUser)
    // resData := make(map[string]interface{})
    // resData["data"] = wxAppUser
    // resData["token"] = sid
    // ctx.JSON(http.StatusOK, gin.H{
    //     "errNo": model.ErrorCode.SUCCESS,
    //     "msg":   "success",
    //     "data":  resData,
    // })
    // return
}
