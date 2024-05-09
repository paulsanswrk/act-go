package entities

import (
	"time"
)

const (
	LogPost = iota
	LogWarning
	LogError
)

type Log struct {
	Id         uint      `gorm:"primaryKey"`
	Category   int       `gorm:"not null;default:0"`
	Message    string    `gorm:"type:text"`
	Tag        string    `gorm:"type:text"`
	Request    string    `gorm:"type:text"`
	Response   string    `gorm:"type:text"`
	Module     string    `gorm:"type:text"`
	StackTrace string    `gorm:"type:text;column:stacktrace"`
	Date       time.Time `gorm:"default:now()"`
}

func (Log) TableName() string {
	return "app.log"
}
