package model

import (
    // "time"

    // "github.com/globalsign/mgo/bson"
)

type Payments struct {
    Id                      int                     `json:"id" bson:"id"`
    Name                    string                  `json:"name" bson:"name"`
    Description             string                  `json:"description" bson:"description"`

    IsOnlinePayment         bool                    `json:"is_online_payment" bson:"is_online_payment"`
    DisabledReason          string                  `json:"disabled_reason" bson:"disabled_reason"`
    SelectState             int                     `json:"select_state" bson:"select_state"`
    promotion               []PromotionStruct       `json:"promotion" bson:"promotion"`
}
