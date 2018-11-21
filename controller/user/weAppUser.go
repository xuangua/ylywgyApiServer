package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	// "github.com/kataras/iris"
	"github.com/satori/go.uuid"
	"github.com/xuangua/ylywgyApiServer/config"
	"github.com/xuangua/ylywgyApiServer/controller/common"
	"github.com/xuangua/ylywgyApiServer/model"
	"github.com/xuangua/ylywgyApiServer/utils"
)

// WeAppLogin 微信小程序登录
// func WeAppLogin(ctx *iris.Context) {
func WeAppLogin(ctx *gin.Context) {
	SendErrJSON := common.SendErrJSON
	// code := ctx.FormValue("code")
	code := ctx.Query("code") //code, _ := ctx.Get("user")
	if code == "" {
		SendErrJSON("code不能为空", ctx)
		return
	}
	appID := config.WeAppConfig.AppID
	secret := config.WeAppConfig.AppSecret
	CodeToSessURL := config.WeAppConfig.CodeToSessURL
	CodeToSessURL = strings.Replace(CodeToSessURL, "{appid}", appID, -1)
	CodeToSessURL = strings.Replace(CodeToSessURL, "{secret}", secret, -1)
	CodeToSessURL = strings.Replace(CodeToSessURL, "{code}", code, -1)

	resp, err := http.Get(CodeToSessURL)
	if err != nil {
		fmt.Println(err.Error())
		SendErrJSON("error", ctx)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		SendErrJSON("error", ctx)
		return
	}

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Println(err.Error())
		SendErrJSON("error", ctx)
		return
	}

	if _, ok := data["session_key"]; !ok {
		fmt.Println("session_key 不存在")
		fmt.Println(data)
		SendErrJSON("error", ctx)
		return
	}

	sid := uuid.NewV4().String()
	var openID string
	var sessionKey string
	openID = data["openid"].(string)
	sessionKey = data["session_key"].(string)
	// expireTime = data["session_key"].(string)

	session := sessions.Default(ctx)
	session.Set(sid, data)
	// session.Set(openID, sessionKey)
	session.Save()
	fmt.Println("wx user access [weAppOpenID: %s, weAppSessionKey: %d ]has been saved to session\n", openID, sessionKey)

	var sql string
	var user model.User
	// var userOpenId uint
	sql = "open_id = ?"
	if err := model.DB.Where(sql, openID).First(&user).Error; err != nil {
		var newUser model.User
		newUser.Uid = uuid.NewV4().String()
		newUser.OpenId = openID
		newUser.Role = model.UserRoleNormal
		newUser.Status = model.UserStatusInActive

		if err := model.DB.Create(&newUser).Error; err != nil {
			SendErrJSON("error", ctx)
			return
		}
		// userOpenId = newUser.OpenId
	} /* else {
		userOpenId = user.OpenId
	}*/

	if err := model.UserToRedisWithOpenId(user); err != nil {
		SendErrJSON("内部错误.", ctx)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": openID,
	})
	tokenString, err := token.SignedString([]byte(config.ServerConfig.TokenSecret))
	if err != nil {
		fmt.Println(err.Error())
		SendErrJSON("内部错误", ctx)
		return
	}

	fmt.Println("wx user access [tokenString: %s, ]has been saved to session\n", tokenString)

	resData := make(map[string]interface{})
	resData[config.ServerConfig.SessionID] = sid
	ctx.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  resData,
	})
}

// LoginWithUserInfo 设置小程序用户加密信息
func LoginWithUserInfo(ctx *gin.Context) {
	SendErrJSON := common.SendErrJSON
	type EncryptedUser struct {
		EncryptedData string `json:"encryptedData"`
		IV            string `json:"iv"`
	}

	userInter, exists := ctx.Get("user")
	var user model.User
	if exists {
		user = userInter.(model.User)
	}

	sidInter, exists := ctx.Get("sid")
	var sid string
	if exists {
		sid = sidInter.(string)
	}

	sessionDataInter, exists := ctx.Get("sessionData")
	var wxSessionData map[string]interface{}
	// var openID string
	var sessionKey string
	if exists {
		wxSessionData = sessionDataInter.(map[string]interface{})
		sessionKey = wxSessionData["session_key"].(string)
	}
	// var openID string
	// var sessionKey string
	// session := sessions.Default(ctx)
	// userWxSessionData := session.Get(sid)
	// wxSessionData, ok := userWxSessionData.(map[string]interface{})
	// if ok {
	// 	/* act on str */
	// 	openID = wxSessionData["openid"].(string)
	// 	sessionKey = wxSessionData["session_key"].(string)
	// 	// expireTime = data["session_key"].(string)
	// 	if (openID == "") || (sessionKey == "") {
	// 		SendErrJSON("session error", ctx)
	// 		return
	// 	}
	// } else {
	// 	/* not string */
	// 	SendErrJSON("session error", ctx)
	// 	return
	// }

	var weAppUser EncryptedUser
	if err := ctx.ShouldBindWith(&weAppUser, binding.JSON); err != nil {
		SendErrJSON("invalid binding Json data", ctx)
		return
	}

	userInfoStr, err := utils.DecodeWeAppUserInfo(weAppUser.EncryptedData, sessionKey, weAppUser.IV)
	if err != nil {
		fmt.Println(err.Error())
		SendErrJSON("error", ctx)
		return
	}

	var wxAppUser model.WeAppUser
	if err := json.Unmarshal([]byte(userInfoStr), &wxAppUser); err != nil {
		SendErrJSON("error", ctx)
		return
	}

	updatedData := map[string]interface{}{
		"name":       wxAppUser.Nickname,
		"sex":        wxAppUser.Gender,
		"avatar_url": wxAppUser.AvatarURL,
	}
	if err := model.DB.Model(&user).Update(updatedData).Error; err != nil {
		SendErrJSON("update user wx data to database error", ctx)
		return
	}
	user.Name = wxAppUser.Nickname
	user.Sex = wxAppUser.Gender
	user.AvatarURL = wxAppUser.AvatarURL

	if model.UserToRedisWithOpenId(user) != nil {
		SendErrJSON("update user wx data to redis error", ctx)
		return
	}
	// var sql string
	// var user model.User
	// sql = "open_id = ?"
	// if err := model.DB.Where(sql, openID).First(&user).Error; err != nil {
	// 	SendErrJSON("can't find user with open_id " + openID, ctx)
	// 	return
	// } else {
	// 	user.Name = wxAppUser.Nickname
	// 	user.Sex = wxAppUser.Gender
	// 	user.AvatarURL = wxAppUser.AvatarURL
	// }

	// session.Set("weAppUser", wxAppUser)
	resData := make(map[string]interface{})
	resData["data"] = wxAppUser
	resData["token"] = sid
	ctx.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  resData,
	})
	return
}

// GetWxUserInfo 获取小程序用户yewu信息
func GetWxUserInfo(ctx *gin.Context) {
	// SendErrJSON := common.SendErrJSON
	// type RspUser struct {
	// 	Nickname    string `json:"nick_name"`
	// 	AvatarURL   string `json:"avatar_url"`
	// }

	userInter, exists := ctx.Get("user")
	var user model.User
	if exists {
		user = userInter.(model.User)
	}

	// build responsed user
	respUser := map[string]interface{}{
		"nick_name":  user.Name,
		"sex":        user.Sex,
		"avatar_url": user.AvatarURL,
	}

	// // build responsed user_data
	// if model.DB.Where("open_id = ? AND uid = ?", 0, user.ID).Order("created_at DESC").
	// 	Offset((pageNo-1)*pageSize).Limit(pageSize).Find(&messages).Error != nil {
	// 	SendErrJSON("error", c)
	// 	return
	// }

	// if err := model.DB.Where(sql, user.OpenId).First(&user).Error; err != nil {
	// 	SendErrJSON("can't find user with open_id " + openID, ctx)
	// 	return
	// }

	resData := make(map[string]interface{})
	resData["user"] = respUser
	// resData["user_data"] = sid
	ctx.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  resData,
	})
	return
}

// YesterdayRegisterUser 昨日注册的用户数
func YesterdayRegisterUser(ctx *gin.Context) {
	var user model.WxUser
	count := user.YesterdayRegisterUser()
	resData := make(map[string]interface{})
	resData["count"] = count
	ctx.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  resData,
	})
}

// TodayRegisterUser 今日注册的用户数
func TodayRegisterUser(ctx *gin.Context) {
	var user model.WxUser
	count := user.TodayRegisterUser()
	resData := make(map[string]interface{})
	resData["count"] = count
	ctx.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  resData,
	})
}

// Latest30Day 近30天，每天注册的新用户数
// func Latest30Day(ctx *gin.Context) {
// 	var users model.UserPerDay
// 	result := users.Latest30Day()
// 	var data iris.Map
// 	if result == nil {
// 		data = iris.Map{
// 			"users": [0]int{},
// 		}
// 	} else {
// 		data = iris.Map{
// 			"users": result,
// 		}
// 	}
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"errNo": model.ErrorCode.SUCCESS,
// 		"msg":   "success",
// 		"data":  data,
// 	})
// }

// Analyze 用户分析
func Analyze(ctx *gin.Context) {
	var user model.WxUser
	now := time.Now()
	nowSec := now.Unix()              //秒
	yesterdaySec := nowSec - 24*60*60 //秒
	yesterday := time.Unix(yesterdaySec, 0)

	yesterdayCount := user.PurchaseUserByDate(yesterday)
	todayCount := user.PurchaseUserByDate(now)
	yesterdayRegisterCount := user.YesterdayRegisterUser()
	todayRegisterCount := user.TodayRegisterUser()

	data := make(map[string]interface{})
	data["todayNewUser"] = todayRegisterCount
	data["yesterdayNewUser"] = yesterdayRegisterCount
	data["todayPurchaseUser"] = todayCount
	data["yesterdayPurchaseUser"] = yesterdayCount
	// data := iris.Map{
	// 	"todayNewUser":          todayRegisterCount,
	// 	"yesterdayNewUser":      yesterdayRegisterCount,
	// 	"todayPurchaseUser":     todayCount,
	// 	"yesterdayPurchaseUser": yesterdayCount,
	// }

	ctx.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  data,
	})
}

// WeAppLogin 微信小程序登录
// func WeAppLogin(ctx *iris.Context) {
func WeAppLogin1(ctx *gin.Context) {
	SendErrJSON := common.SendErrJSON
	// code := ctx.FormValue("code")
	code := ctx.Query("code") //code, _ := ctx.Get("user")
	if code == "" {
		SendErrJSON("code不能为空", ctx)
		return
	}

	code2 := ctx.Query("code2") //code, _ := ctx.Get("user")
	if code2 == "" {
		SendErrJSON("code2不能为空", ctx)
		return
	}

	session := sessions.Default(ctx)
	session.Set("weAppOpenID", code)
	session.Set("weAppSessionKey", code2)
	session.Save()
	fmt.Println("wx user access [weAppOpenID: %s, weAppSessionKey: %d ]has been saved to session\n", code, code2)

	resData := make(map[string]interface{})
	resData[config.ServerConfig.SessionID] = code
	ctx.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  resData,
	})
}

// LoginWithUserInfo 设置小程序用户加密信息
func LoginWithUserInfo1(ctx *gin.Context) {
	SendErrJSON := common.SendErrJSON
	session := sessions.Default(ctx)
	sessionKey := session.Get("weAppSessionKey")
	sessionKeyStr, ok := sessionKey.(string)
	if ok {
		/* act on str */
		if sessionKeyStr == "" {
			SendErrJSON("session error", ctx)
			return
		}
	} else {
		/* not string */
		SendErrJSON("session error", ctx)
		return
	}

	resData := make(map[string]interface{})
	// resData[config.ServerConfig.SessionID] = session.ID()
	ctx.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  resData,
	})
	return
}
