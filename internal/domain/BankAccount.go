package domain

import "time"

type BankAccount struct {
	ID          uint      `json:"id" gorm:"PrimaryKey"`
	UserId      uint      `json:"user_id"`
	BankAccount uint      `json:"bank_account"`
	SwiftCode   string    `json:"swift_code"`
	PaymentType string    `json:"payment_type"`
	CreatedAt   time.Time `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"default:current_timestamp"`
}
