package price

import (
	"time"
)

type Price struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	CommodityID  uint      `json:"commodity_id"`
	Value        float64   `json:"value"`
	DateRecorded time.Time `json:"date_recorded"`
}
