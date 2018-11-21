package main

import (
    "fmt"
    "io"
    "os"
    "time"

    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"
    "github.com/gin-gonic/gin"
    "github.com/silenceper/wechat"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    "github.com/xuangua/ylywgyApiServer/config"
    "github.com/xuangua/ylywgyApiServer/cron"
    "github.com/xuangua/ylywgyApiServer/middleware"
    "github.com/xuangua/ylywgyApiServer/model"
    "github.com/xuangua/ylywgyApiServer/router"
    "github.com/xuangua/ylywgyApiServer/utils"
)

func main() {
    fmt.Println("gin.Version: ", gin.Version)
    if config.ServerConfig.Env != model.DevelopmentMode {
        gin.SetMode(gin.ReleaseMode)
        // Disable Console Color, you don't need console color when writing the logs to file.
        gin.DisableConsoleColor()
        // Logging to a file.
        logFile, err := os.OpenFile(config.ServerConfig.LogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
        if err != nil {
            fmt.Printf(err.Error())
            os.Exit(-1)
        }
        gin.DefaultWriter = io.MultiWriter(logFile)
    }

    // Creates a router without any middleware by default
    app := gin.New()

    // Set a lower memory limit for multipart forms (default is 32 MiB)
    maxSize := int64(config.ServerConfig.MaxMultipartMemory)
    app.MaxMultipartMemory = maxSize << 20 // 3 MiB

    // Global middleware
    // Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
    // By default gin.DefaultWriter = os.Stdout
    app.Use(gin.Logger())

    // Recovery middleware recovers from any panics and writes a 500 if there was one.
    app.Use(gin.Recovery())

    app.Use(middleware.APIStatsD())

    // Setup gin Session
    store := cookie.NewStore([]byte("secret"))
    // store := sessions.NewCookieStore([]byte("secret"))
    store.Options(sessions.Options{
        MaxAge: int(30 * time.Minute), //30min
        Path:   "/",
    })
    app.Use(sessions.Sessions("mysession", store))

    //使用memcache/redis保存access_token
    redis:=model.GetRedis()

    //配置微信参数
    wxConfig := &wechat.Config{
        AppID:          config.WeAppConfig.AppID,
        AppSecret:      config.WeAppConfig.AppSecret,
        Token:          config.WeAppConfig.Token,
        EncodingAESKey: config.WeAppConfig.EncodingAESKey,
        PayMchID:       config.WeAppConfig.PayMchID,
        PayNotifyURL:   config.WeAppConfig.PayNotifyURL,
        PayKey:         config.WeAppConfig.PayKey,
        Cache: redis}
    utils.WechatInst = wechat.NewWechat(wxConfig)

    router.Route(app)

    if config.ServerConfig.StatsEnabled {
        cron.New().Start()
    }

    app.Run(":" + fmt.Sprintf("%d", config.ServerConfig.Port))
}
