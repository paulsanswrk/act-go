package bngx

import "time"

type BotPlan struct {
	Id           uint      `gorm:"primaryKey;column:id"`
	BotID        string    `gorm:"type:text;column:bot_id"`
	OpenOrderID  string    `gorm:"type:text;column:open_order_id"`
	CloseOrderID string    `gorm:"type:text;column:close_order_id"`
	Description  string    `gorm:"type:text;column:description"`
	Date         time.Time `gorm:"default:now();column:date"`
}

func (BotPlan) TableName() string {
	return "bingx.bot_plan"
}
