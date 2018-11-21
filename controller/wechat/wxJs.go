package wechat

import (
    // "encoding/json"
    "fmt"
    "net/http"
    // "strings"
    // "time"

    // "github.com/dgrijalva/jwt-go"
    // "github.com/gin-contrib/sessions"
    "github.com/gin-gonic/gin"
    // "github.com/gin-gonic/gin/binding"
    // "github.com/satori/go.uuid"
    // "github.com/xuangua/ylywgyApiServer/config"
    "github.com/xuangua/ylywgyApiServer/controller/common"
    // "github.com/xuangua/ylywgyApiServer/model"
    "github.com/xuangua/ylywgyApiServer/utils"
)

// WeJS auth 获取auth Config
func GetWxAuthConfig(ctx *gin.Context) {
    SendErrJSON := common.SendErrJSON
    // accessurl := ctx.FormValue("accessurl")
    accessurl := ctx.Query("accessurl")
    if accessurl == "" {
        SendErrJSON("accessurl不能为空", ctx)
        return
    }

    wechatJs := utils.WechatInst.GetJs()
    jsAuthConfig, err := wechatJs.GetConfig(accessurl)
    if err != nil {
        fmt.Println(err.Error())
        SendErrJSON("get JsAuth Config error", ctx)
        return
    }

    // ctx.JSON(http.StatusOK, gin.H{
    //     "errNo": model.ErrorCode.SUCCESS,
    //     "msg":   "success",
    //     "data":  resData,
    // })
    ctx.JSON(http.StatusOK, jsAuthConfig)
}
