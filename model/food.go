package model

import (
    "time"
)

// FoodActivity ---------- superuser 管理 (TODO: not used now, TBD)
type FoodActivity struct {
    // ID       bson.ObjectId `bson:"_id"`
    Description   string `json:"description" bson:"description"`
    IconColor     string `json:"icon_color" bson:"icon_color"`
    IconName      string `json:"icon_name" bson:"icon_name"`
    Id            int    `json:"id" bson:"id"`
    Name          string `json:"name" bson:"name"`
    RankingWeight int    `json:"ranking_weight" bson:"ranking_weight"`
}

// FoodAttrubution ---------- superuser 管理
type FoodAttrubution struct {
    // ID       bson.ObjectId `bson:"_id"`
    Description   string `json:"description" bson:"description"`
    IconColor     string `json:"icon_color" bson:"icon_color"`
    IconName      string `json:"icon_name" bson:"icon_name"`
    Id            int    `json:"id" bson:"id"`
    Name          string `json:"name" bson:"name"`
}

// FoodSpecSummary ---------- 管理 后台添加
type FoodSpecSummary struct {
    Name            string      `json:"name" bson:"name"`       // vaule （应该）一直是 ‘规格’？
    SpecNames       []string    `json:"values" bson:"values"`   // 各个 规格 的名称
}

// 食品规格
type FoodSpec struct {
    SpecName            string      `json:"specs_name" bson:"specs_name"`   // 规格 的名称, same as FoodSpecSummarySpecNames[x]
    FoodId              int         `json:"food_id" bson:"food_id"`         // 每个 规格 占用一个 FoodId
    ItemId              int         `json:"item_id" bson:"item_id"`         // Same as Food.ItemId
    Name                string      `json:"name" bson:"name"`               // Same as Food.Name
    ShopId              int         `json:"shop_id" bson:"shop_id"`         // Same as Food.ShopId

    SkuId               int         `json:"sku_id" bson:"sku_id"`
    Stock               int         `json:"stock" bson:"stock"`
    CheckoutMode        int         `json:"checkout_mode" bson:"checkout_mode"`
    IsEssential         bool        `json:"is_essential" bson:"is_essential"`
    RecentPopularity    int         `json:"recent_popularity" bson:"recent_popularity"`
    IsSoldOut           bool        `json:"sold_out" bson:"sold_out"`
    Price               float64     `json:"price" bson:"price"`
    PromotionStock      int         `json:"promotion_stock" bson:"promotion_stock"`
    RecentRating        float64     `json:"recent_rating" bson:"recent_rating"`
    PackingFee          float64     `json:"packing_fee" bson:"packing_fee"`
    PinyinName          string      `json:"pinyin_name" bson:"pinyin_name"`  // Keep ""
    OriginalPrice       float64     `json:"original_price" bson:"original_price"`               // Same as Food.Name
}

// FoodTable
type Foods struct {
    // ShopSummary Shop `json:"shop" bson:"shop"`
    Description string      `json:"description" bson:"description"`
    ItemId      int         `json:"item_id" bson:"item_id"`             // 每一个通过管理后台添加的商品有一个 ItemId.  
                                                                        // 如果单规格， 该Item只占用一个FoodId, 
                                                                        // 如果多规格，该Item会占用多个FoodId 
                                                                        // 分配Id时,需保证每个添加的商品的ItemId和它第一个规格的FoodId相同
                                                                        // 比如第4个商品的ItemId=5，而第4个商品有三个规格（FoodId = 5，6，7）
                                                                        // 在添加第5个商品时，它的ItemId应该为8，且第一个规格的FoodId 也应该为8
    Name        string      `json:"name" bson:"name"`
    ImagePath   string      `json:"image_path" bson:"image_path"`
    FoodCategoryId  int                 `json:"food_category_id" bson:"food_category_id"`
    ShopId          int                 `json:"shop_id" bson:"shop_id"`

    // SkuId       int         `json:"sku_id" bson:"sku_id"`
    Price       float64     `json:"price" bson:"price"`         // min prize of specs
    Tips        string      `json:"tips" bson:"tips"`

    Rating      float64     `json:"rating" bson:"rating"`
    RatingCount      int     `json:"rating_count" bson:"rating_count"`
    MonthSales      int     `json:"month_sales" bson:"month_sales"`
    SatisfyRate      int     `json:"satisfy_rate" bson:"satisfy_rating"`
    SatisfyCount      int     `json:"satisfy_count" bson:"satisfy_rate"`

    IsFeatured          bool        `json:"is_featured" bson:"is_featured"`
    IsEssential         bool        `json:"is_essential" bson:"is_essential"`
    ServerUtc           time.Time             `json:"server_utc" bson:"server_utc"`

    Activity      []FoodActivity    `json:"activity" bson:"activity"`
    Attributions    []FoodAttrubution   `json:"attributes" bson:"attributes"`
    SpecSummary     []FoodSpecSummary   `json:"specifications" bson:"specifications"`       //（应该）只有一个元素？
    SpecDetail      []FoodSpec          `json:"specfoods" bson:"specfoods"`                 // 每个规格 占用一个元素
}

// db checked