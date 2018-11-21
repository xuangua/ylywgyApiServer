
package model

import (
    // "time"

    // "github.com/globalsign/mgo/bson"
    "github.com/xuangua/ylywgyApiServer/utils"
)

type SchoolData struct {
    Id                  int             `bson:"id" json:"id"`
    IsActive            bool            `bson:"is_active" json:"is_active"`
    SchoolName          string              `bson:"school_name" json:"school_name"`
    SchoolCampusName    string              `bson:"school_campus_name" json:"school_campus_name"`

    QQLbsId             string          `bson:"qqlbs_id" json:"qqlbs_id"`
    Title               string          `bson:"title" json:"title"`
    Address             string          `bson:"address" json:"address"`
    Category            string          `bson:"category" json:"category"`
    Type                int             `bson:"type" json:"type"`

    AdLocation          utils.Location        `bson:"location" json:"location"`

    AdCode              int              `bson:"adcode" json:"adcode"`
    Province            string              `bson:"province" json:"province"`
    City                string              `bson:"city" json:"city"`
    District            string              `bson:"district" json:"district"`
}

type ShipAddr struct {
    Id                  int             `bson:"id" json:"id"`

    Name            string           `bson:"name" json:"name"`
    Sex                int           `bson:"sex" json:"sex"`
    Phone            string           `bson:"phone" json:"phone"`
    Tag            string           `bson:"tag" json:"tag"`

    AddressDetail             string          `bson:"address_detail" json:"address_detail"`

    AdLocation          utils.Location        `bson:"location" json:"location"`

    Province            string           `bson:"province" json:"province"`
    City                string           `bson:"city" json:"city"`
    District            string           `bson:"district" json:"district"`
    SchoolName          string              `bson:"school_name" json:"school_name"`
    SchoolCampusName    string              `bson:"school_campus_name" json:"school_campus_name"`
    SchoolDormAreaName  string                `json:"school_dormarea_name" bson:"school_dormarea_name"`
    SchoolDormName      string                `json:"school_dorm_name" bson:"school_dorm_name"`
    SchoolDormRoomName  string                `json:"school_dorm_room_name" bson:"school_dorm_room_name"`

}

type CustomerShipAddr struct {
    UserId              string          `bson:"user_id" json:"user_id"`
    ShipAddrArray       []ShipAddr      `bson:"ship_addr_array" json:"ship_addr_array"`
}
