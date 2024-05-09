package entities

type Price struct {
	Symbol        string  `gorm:"primaryKey;column:symbol;type:varchar(50);not null"`
	Time          int64   `gorm:"primaryKey;column:time;not null;autoIncrement:false"`
	Open          float64 `gorm:"column:open;type:numeric(20,10);not null"`
	High          float64 `gorm:"column:high;type:numeric(20,10);not null"`
	Low           float64 `gorm:"column:low;type:numeric(20,10);not null"`
	Close         float64 `gorm:"column:close;type:numeric(20,10);not null"`
	Volume        float64 `gorm:"column:volume;type:numeric(20,10);not null"`
	WeightedPrice float64 `gorm:"column:weighted_price;type:numeric(20,10);not null"`

	// Specify the composite primary key
	// GORM will automatically consider both Symbol and Time as primary keys
	// You can also use gorm.Model if you need default fields like ID, CreatedAt, UpdatedAt, DeletedAt
}

// TableName specifies the table name for the Price model
func (Price) TableName() string {
	return "binance.prices"
}
