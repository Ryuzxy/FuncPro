package Komoditas 

import "time"

type Commodity struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Nama      string    `json:"nama"`
	Category  string    `json:"category"`
	Unit      string    `json:"unit"`
	CreatedAt time.Time `json:"created_at"`
}
