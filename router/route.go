package router

import (
    "github.com/gin-gonic/gin"
    "github.com/xuangua/ylywgyApiServer/config"
    "github.com/xuangua/ylywgyApiServer/controller/article"
    "github.com/xuangua/ylywgyApiServer/controller/baidu"
    "github.com/xuangua/ylywgyApiServer/controller/book"
    "github.com/xuangua/ylywgyApiServer/controller/category"
    "github.com/xuangua/ylywgyApiServer/controller/collect"
    "github.com/xuangua/ylywgyApiServer/controller/comment"
    "github.com/xuangua/ylywgyApiServer/controller/common"
    "github.com/xuangua/ylywgyApiServer/controller/crawler"
    "github.com/xuangua/ylywgyApiServer/controller/haixun"
    "github.com/xuangua/ylywgyApiServer/controller/keyvalueconfig"
    "github.com/xuangua/ylywgyApiServer/controller/message"
    "github.com/xuangua/ylywgyApiServer/controller/shopping"
    "github.com/xuangua/ylywgyApiServer/controller/stats"
    "github.com/xuangua/ylywgyApiServer/controller/user"
    "github.com/xuangua/ylywgyApiServer/controller/vote"
    "github.com/xuangua/ylywgyApiServer/controller/wechat"
    "github.com/xuangua/ylywgyApiServer/controller/address"
    "github.com/xuangua/ylywgyApiServer/controller/cart"
    "github.com/xuangua/ylywgyApiServer/controller/order"
    "github.com/xuangua/ylywgyApiServer/controller/pay"
    "github.com/xuangua/ylywgyApiServer/middleware"
)

// Route 路由
func Route(router *gin.Engine) {
    apiPrefix := config.ServerConfig.APIPrefix

    api := router.Group(apiPrefix, middleware.RefreshTokenCookie)
    {
        api.GET("/weAppLogin", user.WeAppLogin)
        api.POST("/setWeAppUser", user.LoginWithUserInfo)

        api.GET("/siteinfo", common.SiteInfo)
        api.POST("/signin", user.Signin)
        api.POST("/signup", user.Signup)
        api.POST("/signout", middleware.SigninRequired,
            user.Signout)
        api.POST("/upload", middleware.SigninRequired,
            common.UploadHandler)
        api.POST("crawlnotsavecontent", middleware.EditorRequired,
            crawler.CrawlNotSaveContent)

        api.POST("/active/sendmail", user.ActiveSendMail)
        api.POST("/active/user/:id/:secret", user.ActiveAccount)

        api.POST("/reset/sendmail", user.ResetPasswordMail)
        api.GET("/reset/verify/:id/:secret", user.VerifyResetPasswordLink)
        api.POST("/reset/password/:id/:secret", user.ResetPassword)

        api.GET("/user/info", middleware.SigninRequired,
            user.SecretInfo)
        api.GET("/user/score/top10", user.Top10)
        api.GET("/user/score/top100", user.Top100)
        api.GET("/user/info/detail", middleware.SigninRequired,
            user.InfoDetail)
        api.GET("/user/info/public/:id", user.PublicInfo)
        api.POST("/user/uploadavatar", middleware.SigninRequired,
            user.UploadAvatar)
        api.POST("/user/career/add", middleware.SigninRequired,
            user.AddCareer)
        api.POST("/user/school/add", middleware.SigninRequired,
            user.AddSchool)
        api.PUT("/user/update/:field", middleware.SigninRequired,
            user.UpdateInfo)
        api.PUT("/user/password/update", middleware.SigninRequired,
            user.UpdatePassword)
        api.DELETE("/user/career/delete/:id", middleware.SigninRequired,
            user.DeleteCareer)
        api.DELETE("/user/school/delete/:id", middleware.SigninRequired,
            user.DeleteSchool)

        api.GET("/messages/unread", middleware.SigninRequired,
            message.Unread)
        api.GET("/messages/read/:id", middleware.SigninRequired,
            message.Read)

        api.GET("/categories", category.List)

        api.GET("/articles", article.List)
        api.GET("/articles/max/bycomment", article.ListMaxComment)
        api.GET("/articles/max/bybrowse", article.ListMaxBrowse)
        api.GET("/articles/top/global", article.Tops)
        api.GET("/articles/info/:id", article.Info)
        api.GET("/articles/user/:userID", article.UserArticleList)
        api.POST("/articles/create", middleware.SigninRequired,
            article.Create)
        api.POST("/articles/top/:id", middleware.EditorRequired,
            article.Top)
        api.PUT("/articles/update", middleware.SigninRequired,
            article.Update)
        api.DELETE("/articles/delete/:id", middleware.SigninRequired,
            article.Delete)
        api.DELETE("/articles/deltop/:id", middleware.EditorRequired,
            article.DeleteTop)

        api.GET("/collects", collect.Collects)
        api.GET("/collects/folders/withsource", middleware.SigninRequired,
            collect.FoldersWithSource)
        api.GET("/collects/user/:userID/folders", collect.Folders)
        api.POST("/collects/create", middleware.SigninRequired,
            collect.CreateCollect)
        api.POST("/collects/folder/create", middleware.SigninRequired,
            collect.CreateFolder)
        api.DELETE("/collects/delete/:id", middleware.SigninRequired,
            collect.DeleteCollect)

        api.GET("/comments/user/:userID", comment.UserCommentList)
        api.GET("/comments/source/:sourceName/:sourceID", comment.SourceComments)
        api.POST("/comments/create", middleware.SigninRequired,
            comment.Create)
        api.PUT("/comments/update", middleware.SigninRequired,
            comment.Update)
        api.DELETE("/comments/delete/:id", middleware.SigninRequired,
            comment.Delete)

        api.GET("/votes", vote.List)
        api.GET("/votes/info/:id", vote.Info)
        api.GET("/votes/max/bybrowse", vote.ListMaxBrowse)
        api.GET("/votes/max/bycomment", vote.ListMaxComment)
        api.GET("/votes/user/:userID", vote.UserVoteList)
        api.POST("/votes/create", middleware.EditorRequired,
            vote.Create)
        api.POST("/votes/item/create", middleware.EditorRequired,
            vote.CreateVoteItem)
        api.POST("/votes/uservote/:id", middleware.SigninRequired,
            vote.UserVoteVoteItem)
        api.PUT("/votes/update", middleware.EditorRequired,
            vote.Update)
        api.PUT("/votes/item/edit", middleware.EditorRequired,
            vote.EditVoteItem)
        api.DELETE("/votes/delete/:id", middleware.EditorRequired,
            vote.Delete)
        api.DELETE("/votes/item/delete/:id", middleware.EditorRequired,
            vote.DeleteItem)

        api.GET("/books", book.List)
        api.GET("/books/categories", category.BookCategoryList)
        api.GET("/books/my/:userID", middleware.SigninRequired,
            book.MyBooks)
        api.GET("/books/user/public/:userID", book.UserPublicBooks)
        api.GET("/books/info/:id", middleware.SetContextUser,
            book.Info)
        api.GET("/books/chapters/:bookID", middleware.SetContextUser,
            book.Chapters)
        api.GET("/books/chapter/:chapterID", middleware.SetContextUser,
            book.Chapter)
        api.POST("/books", middleware.SigninRequired,
            book.Create)
        api.POST("/books/chapters", middleware.SigninRequired,
            book.CreateChapter)
        api.PUT("/books/update", middleware.SigninRequired,
            book.Update)
        api.PUT("/books/updatename", middleware.SigninRequired,
            book.UpdateName)
        api.PUT("/books/publish/:bookID", middleware.SigninRequired,
            book.Publish)
        api.PUT("/books/chapters/content", middleware.SigninRequired,
            book.UpdateChapterContent)
        api.PUT("/books/chapters/updatename", middleware.SigninRequired,
            book.UpdateChapterName)
        api.DELETE("/books/delete/:id", middleware.SigninRequired,
            book.Delete)
        api.DELETE("/books/chapters/:chapterID", middleware.SigninRequired,
            book.DeleteChapter)

        api.GET("/stats/visit", stats.PV)
    }

    wxapi := router.Group(apiPrefix, middleware.RefreshTokenCookie)
    {
        wxapi.GET("/v3/login", user.WeAppLogin)
        wxapi.POST("/v3/login-with-user-info", middleware.WxSessionRequired, user.LoginWithUserInfo)
        wxapi.GET("/v3/user-index", middleware.WxSessionRequired, user.GetWxUserInfo)
    }

    adminAPI := router.Group(apiPrefix+"/admin", middleware.RefreshTokenCookie, middleware.AdminRequired)
    {
        adminAPI.POST("/keyvalueconfig", keyvalueconfig.SetKeyValue)

        adminAPI.GET("/users", user.AllList)

        adminAPI.GET("/books/categories", category.BookCategoryList)
        adminAPI.POST("/books/categories/create", category.CreateBookCategory)

        adminAPI.GET("/categories", category.List)
        adminAPI.POST("/categories/create", category.Create)
        adminAPI.PUT("/categories/update", category.Update)

        adminAPI.GET("/articles", article.AllList)
        adminAPI.PUT("/articles/status/update", article.UpdateStatus)

        adminAPI.GET("/comments", comment.Comments)
        adminAPI.PUT("/comments/update/status/:id", comment.UpdateStatus)

        adminAPI.GET("/crawl/account", crawler.CrawlAccount)
        adminAPI.POST("/crawl", crawler.Crawl)
        adminAPI.POST("/customcrawl", crawler.CustomCrawl)
        adminAPI.POST("/crawl/account", crawler.CreateAccount)

        adminAPI.POST("/pushBaiduLink", baidu.PushToBaidu)
    }

    xuawuAPI := router.Group(apiPrefix + "/xuewu")
    {
        xuawuAPI.POST("/admin/login", user.Signin)
        xuawuAPI.POST("/admin/signup", user.Signup)
        xuawuAPI.POST("/admin/register", user.Signup)
        xuawuAPI.POST("/admin/signout", middleware.SigninRequired, user.Signout)
        xuawuAPI.GET("/admin/shopping/shopCategories", middleware.AdminRequired, shopping.GetShopCategoryList)
        xuawuAPI.POST("/admin/shopping/addShop", middleware.AdminRequired, shopping.AddShop)
        xuawuAPI.POST("/admin/shopping/updateShop", middleware.AdminRequired, shopping.UpdateShop)
        xuawuAPI.DELETE("/admin/shopping/deleteShop/:shopId", middleware.AdminRequired, shopping.DeleteShop)
        xuawuAPI.GET("/admin/shopping/getAdminUserShopList", middleware.AdminRequired, shopping.GetAdminUserShopList)
        xuawuAPI.GET("/admin/shopping/getShopCount", middleware.AdminRequired, shopping.GetAdminUserTotalShopCount)

        xuawuAPI.POST("/admin/shopping/addFood", middleware.AdminRequired, shopping.AddFood)
        xuawuAPI.POST("/admin/shopping/updateFood", middleware.AdminRequired, shopping.UpdateFood)
        xuawuAPI.DELETE("/admin/shopping/deleteFood/:shopId", middleware.AdminRequired, shopping.DeleteFood)

        xuawuAPI.POST("/wechat/loginWithWxUserInfo", wechat.LoginWithUserInfo)
        xuawuAPI.GET("/wechat/getWxUserInfo", wechat.GetWxWebUserInfo)
        xuawuAPI.GET("/wechat/wxAuthConfig", wechat.GetWxAuthConfig)

        //customer
        xuawuAPI.GET("/customer/get_pois/:geoHash", middleware.SigninRequired, address.GetQQPoisAddress)
        xuawuAPI.GET("/customer/getSupportedSchoolList", middleware.SigninRequired, address.GetSupportedSchoolList)
        xuawuAPI.GET("/customer/getShipAddrList/:user_id", middleware.SigninRequired, user.GetCustomerShipAddrList)
        xuawuAPI.POST("/customer/addShipAddr/:user_id", middleware.SigninRequired, user.AddCustomerShipAddr)
        xuawuAPI.GET("/customer/getShopCategory/index_entry", middleware.SigninRequired, shopping.GetShopCategoryList)
        xuawuAPI.GET("/customer/shopping/getCustomerUserShopList", middleware.SigninRequired, shopping.GetCustomerShopList)
        xuawuAPI.GET("/customer/shopping/getShopMenuFoods", middleware.SigninRequired, shopping.GetCustomerShopMenuFoodsList)

        xuawuAPI.POST("/customer/shopping/carts/checkout", middleware.SigninRequired, cart.CustomerCheckoutCartHandler)
        xuawuAPI.POST("/customer/users/:user_id/carts/:cart_id/orders", middleware.SigninRequired, order.CustomerPlaceOrderHandler)
        xuawuAPI.GET("/customer/payapi/payment/queryOrder", middleware.SigninRequired, order.GetCustomerOrderDetail)
        xuawuAPI.GET("/customer/getUserOrderList", middleware.SigninRequired, order.GetCustomerOrderList)
        xuawuAPI.GET("/customer/payapi/getWxPayParameters", middleware.SigninRequired, pay.GetWxPayParameters)
        xuawuAPI.POST("/customer/payapi/WxPaySucNotify", pay.WxPaySucNotify)

        // xuawuAPI.POST("/shoping/categories/create", shopping.CreateBookCategory)
    }

    haixunErpAPI := router.Group(apiPrefix+"/apis/haixunerp", middleware.RefreshTokenCookie)
    {
        haixunErpAPI.POST("/apis/v1807g/account/login", haixun.Signin)
    }
}
