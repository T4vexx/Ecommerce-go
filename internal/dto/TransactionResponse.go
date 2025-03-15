package dto

import "github.com/stripe/stripe-go/v81"

type SellerOrderDetails struct {
	OrderRefNumber  int     `json:"order_ref_number"`
	OrderStatus     int     `json:"order_status"`
	CreatedAt       string  `json:"created_at"`
	OrderItemId     uint    `json:"order_item_id"`
	ProductId       uint    `json:"product_id"`
	Name            string  `json:"name"`
	ImageUrl        string  `json:"image_url"`
	Price           float64 `json:"price"`
	Qty             uint    `json:"qty"`
	CustomerName    string  `json:"customer_name"`
	CustomerEmail   string  `json:"customer_email"`
	CustomerPhone   string  `json:"customer_phone"`
	CustomerAddress string  `json:"customer_address"`
}

type MakePaymentSuccess struct {
	Message string `json:"message" example:"make payment"`
	PubKey  string `json:"pubKey" example:"exemplo_de_chave_privada_asd2asd2"`
	Secret  string `json:"secret" example:"exemplo_de_segredo_stripe_asd2asd2"`
}

type VerifyPaymentSuccess struct {
	Message  string                `json:"message" example:"Create payment"`
	Response *stripe.PaymentIntent `json:"response"`
}
