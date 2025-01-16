package domain

import "time"

type Address struct {
	ID           uint      `gorm:"primary_key" json:"id"`
	AddressLine1 string    `json:"address_line_1"`
	AddressLine2 string    `json:"address_line_2"`
	City         string    `json:"city"`
	PostalCode   int32     `json:"postal_code"`
	Country      string    `json:"country"`
	UserID       uint      `json:"user_id"`
	CreatedAt    time.Time `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt    time.Time `gorm:"default:current_timestamp" json:"updated_at"`
}
