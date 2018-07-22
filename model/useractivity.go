package model

import "time"

// Message 消息
type UserActivity struct {
	ID            uint       `gorm:"primary_key" json:"id"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	DeletedAt     *time.Time `sql:"index" json:"deletedAt"`
	Type          uint       `json:"type"` //1. 已报名；2. 已参加
	ActivityUrl   string     `json:"activityUrl"`
	ActivityId    string     `json:"activityId"`
	ActivityTitle string     `json:"activityTitle"`
}

const (
	// 1. 已报名
	UserActivityTypeRegistered = 1

	// 2. 已参加
	UserActivityTypeAttended = 2
)
