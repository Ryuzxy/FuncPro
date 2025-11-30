package price

import (
    "time"
    
    "gorm.io/gorm"
)

type Price struct {
    ID           uint           `gorm:"primarykey" json:"id"`
    KomoditasID  uint           `gorm:"not null;index" json:"komoditas_id"`
    Value        float64        `gorm:"type:decimal(10,2);not null" json:"value"`
    Date         time.Time      `gorm:"type:date;not null" json:"date"`
    Market       string         `gorm:"size:100" json:"market"`
    CreatedAt    time.Time      `json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
    DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type PriceAnalysis struct {
    Current    float64 `json:"current"`
    Previous   float64 `json:"previous"`
    Change     float64 `json:"change"`
    ChangePct  float64 `json:"change_percentage"`
    Trend      string  `json:"trend"`
    Volatility float64 `json:"volatility"`
}