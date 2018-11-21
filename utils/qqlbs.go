package utils

import (
    "encoding/json"
    // "errors"
    "fmt"
    // "io/ioutil"
    // "net/http"
    "net/url"
    "github.com/silenceper/wechat/util"
)

const (
    qqLbsApiKey     = "SVZBZ-AL7WV-IMZPP-UBSMC-VHYT7-7NF3J" //QQ 84506525
    getPoisURL      = "https://apis.map.qq.com/ws/geocoder/v1/?location=%s&poi_options=policy=%d&poi_options=category=%s&key=%s&get_poi=1"
)

type Location struct {
    Lat      float64    `json:"lat"`
    Lng      float64    `json:"lng"`
}

type FormattedAddr struct {
    Recommend       string    `json:"recommend"`
    Rough           string    `json:"rough"`
}

type AddrComp struct {
    Nation          string    `json:"nation"`
    Province        string    `json:"province"`
    City            string    `json:"city"`
    District        string    `json:"district"`
    Street          string    `json:"street"`
    StreetNumber    string    `json:"street_number"`
}

type AddrInfo struct {
    NationCode      string    `json:"nation_code"`
    AdCode          string    `json:"adcode"`
    CityCode        string    `json:"city_code"`
    Name            string    `json:"name"`
    AdLocation      Location  `json:"location"`
    Nation          string    `json:"nation"`
    Province        string    `json:"province"`
    City            string    `json:"city"`
    District        string    `json:"district"`
}

type AddrRefer struct {
// TODO: may not need now
}

type PoiDetail struct {
    Id              string    `json:"id"`
    Title           string    `json:"title"`
    Address         string    `json:"address"`
    Category        string    `json:"category"`
    AdLocation      Location  `json:"location"`
    AddressesInfo   json.RawMessage       `json:"ad_info"`
    Distance        float64    `json:"_distance"`
    DirDesc         string    `json:"_dir_desc"`
}

type AddrResultDetail struct {
    GeoLocation         Location        `json:"location"`
    Address             string          `json:"address"`
    FormattedAddresses  FormattedAddr   `json:"formatted_addresses"`
    AddressesComponent  AddrComp        `json:"address_component"`
    AddressesInfo       AddrInfo        `json:"ad_info"`
    // AdRefer             interface{}       `json:"address_reference"`
    AdRefer             json.RawMessage       `json:"address_reference"`

    PoiCount    int  `json:"poi_count"`
    PiosDetail []PoiDetail `json:"pois"`
}

// ResAccessToken 获取用户授权access_token的返回结果
type AddressResult struct {
    Status       int    `json:"status"`
    Message      string `json:"message"`
    // Count        int    `json:"count"`
    RequestId      string `json:"request_id"`

    ResultDetail   AddrResultDetail `json:"result"`
}

func GetQQPois(geoHash string) (response []byte, result AddressResult, err error) {
    // reqData := make(url.Values)
	// reqData["location"] = geoHash
    // reqData["poi_options=policy"] = 2
    // reqData["poi_options=category"] = "教育学校:大学"
    // reqData["key"] = qqLbsApiKey
    // reqData["get_poi"] = 1
    // body := reqData.Encode()
    // fmt.Println(body)
    categoryStr := url.QueryEscape("教育学校:大学")

    urlStr := fmt.Sprintf(getPoisURL, geoHash, 2, categoryStr, qqLbsApiKey)
    // var response []byte
    response, err = util.HTTPGet(urlStr)
    if err != nil {
        return
    }

    if result.Status != 0 {
        err = fmt.Errorf("GetQQPois error : errcode=%v , errmsg=%v", result.Status, result.Message)
        return
    }

    err = json.Unmarshal(response, &result)
    if err != nil {
        return
    }

    return
}
