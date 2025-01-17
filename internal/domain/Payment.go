package domain

import "time"

type Payment struct {
	ID            uint      `gorm:"primary_key" json:"id"`
	UserId        uint      `json:"user_id"`
	CaptureMethod string    `json:"capture_method"`
	Amount        float64   `json:"amount"`
	TransactionId uint      `json:"transaction_id"`
	CustomerId    string    `json:"customer_id"`
	PaymentId     string    `json:"payment_id"`
	Status        string    `json:"status"`
	Response      string    `json:"response"` // response do gateway para ser usado caso o cliente fale que ele pagou mas nao recebeu o pedido
	CreatedAt     time.Time `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"default:current_timestamp"`
}
