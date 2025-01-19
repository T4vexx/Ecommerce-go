package domain

import "time"

type Payment struct {
	ID            uint          `gorm:"primary_key" json:"id"`
	UserId        uint          `json:"user_id"`
	CaptureMethod string        `json:"capture_method"`
	Amount        float64       `json:"amount"`
	OrderId       string        `json:"order_id"`
	CustomerId    string        `json:"customer_id"`
	PaymentId     string        `json:"payment_id"`
	ClientSecret  string        `json:"client_secret"`
	Status        PaymentStatus `json:"status" gorm:"default:initial"` // initial, success, failed
	Response      string        `json:"response"`                      // response do gateway para ser usado caso o cliente fale que ele pagou mas nao recebeu o pedido
	CreatedAt     time.Time     `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt     time.Time     `json:"updated_at" gorm:"default:current_timestamp"`
}

type PaymentStatus string

const (
	PaymentStatusInitial PaymentStatus = "initial"
	PaymentStatusSuccess PaymentStatus = "success"
	PaymentStatusFailed  PaymentStatus = "failed"
	PaymentStatusPending PaymentStatus = "pending"
)
