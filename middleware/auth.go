package middleware

import (
	"errors"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/xuangua/ylywgyApiServer/config"
	"github.com/xuangua/ylywgyApiServer/controller/common"
	"github.com/xuangua/ylywgyApiServer/model"
)

func getUser(c *gin.Context) (model.User, error) {
	var user model.User
	tokenString, cookieErr := c.Cookie("token")

	if cookieErr != nil {
		return user, errors.New("未登录")
	}

	token, tokenErr := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.ServerConfig.TokenSecret), nil
	})

	if tokenErr != nil {
		return user, errors.New("未登录")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// fmt.Printf("claims: %v\n", claims)
		userID := int(claims["id"].(float64))
		var err error
		user, err = model.UserFromRedis(userID)
		if err != nil {
			return user, errors.New("未登录")
		}
		return user, nil
	}
	return user, errors.New("未登录")
}

// SetContextUser 给 context 设置 user
func SetContextUser(c *gin.Context) {
	var user model.User
	var err error
	if user, err = getUser(c); err != nil {
		c.Set("user", nil)
		c.Next()
		return
	}
	c.Set("user", user)
	c.Next()
}

// SigninRequired 必须是登录用户
func SigninRequired(c *gin.Context) {
	SendErrJSON := common.SendErrJSON
	var user model.User
	var err error
	if user, err = getUser(c); err != nil {
		SendErrJSON("未登录", model.ErrorCode.LoginTimeout, c)
		return
	}
	c.Set("user", user)
	c.Next()
}

// EditorRequired 必须是网站编辑
func EditorRequired(c *gin.Context) {
	SendErrJSON := common.SendErrJSON
	var user model.User
	var err error
	if user, err = getUser(c); err != nil {
		SendErrJSON("未登录", model.ErrorCode.LoginTimeout, c)
		return
	}
	if user.Role == model.UserRoleEditor || user.Role == model.UserRoleAdmin || user.Role == model.UserRoleCrawler || user.Role == model.UserRoleSuperAdmin {
		c.Set("user", user)
		c.Next()
	} else {
		SendErrJSON("没有权限", c)
	}
}

// AdminRequired 必须是管理员
func AdminRequired(c *gin.Context) {
	SendErrJSON := common.SendErrJSON
	var user model.User
	var err error
	if user, err = getUser(c); err != nil {
		SendErrJSON("未登录", model.ErrorCode.LoginTimeout, c)
		return
	}
	if user.Role == model.UserRoleAdmin || user.Role == model.UserRoleCrawler || user.Role == model.UserRoleSuperAdmin {
		c.Set("user", user)
		c.Next()
	} else {
		SendErrJSON("没有权限", c)
	}
}

func getUserWithOpenId(c *gin.Context) (model.User, string, map[string]interface{}, error) {
	var user model.User
	var nullSessionData map[string]interface{}
	sid, cookieErr := c.Cookie("sid")

	if cookieErr != nil {
		fmt.Println("cookie sid doesn't exist", cookieErr.Error())
		return user, sid, nullSessionData, errors.New("no sid in Cookie")
	}

	var openID string
	var sessionKey string
	session := sessions.Default(c)
	userWxSessionData := session.Get(sid)
	wxSessionData, ok := userWxSessionData.(map[string]interface{})
	if ok {
		openID = wxSessionData["openid"].(string)
		sessionKey = wxSessionData["session_key"].(string)
		if (openID == "") || (sessionKey == "") {
			fmt.Println("local session doesn't include openid and session_key")
			return user, sid, wxSessionData, errors.New("session error")
		}
		var err error
		user, err = model.UserFromRedisWithOpenId(openID)
		if err != nil {
			fmt.Println(err.Error(), "try get from DB next")
			user, err = model.UserFromDBWithOpenId(openID)
			if err != nil {
				return user, sid, wxSessionData, errors.New("can't find user from DB with openID " + openID)
			}
		}
		return user, sid, wxSessionData, nil
	} else {
		return user, sid, wxSessionData, errors.New("session error")
	}

	return user, sid, wxSessionData, errors.New("未登录")
}

// WxSessionRequired 必须是登录用户
func WxSessionRequired(c *gin.Context) {
	SendErrJSON := common.SendErrJSON
	var user model.User
	var sessionData map[string]interface{}
	var sid string
	var err error
	if user, sid, sessionData, err = getUserWithOpenId(c); err != nil {
		fmt.Println(err.Error(), "get user failed")
		SendErrJSON("未登录", model.ErrorCode.LoginTimeout, c)
		return
	}
	c.Set("user", user)
	c.Set("sid", sid)
	c.Set("sessionData", sessionData)
	c.Next()
}
