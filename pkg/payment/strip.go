package payment

import (
	"errors"
	"fmt"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/paymentintent"
	"log"
)

type PaymentClient interface {
	CreatePayment(amount float64, userId uint, orderId string) (*stripe.PaymentIntent, error)
	GetPaymentStatus(pId string) (*stripe.PaymentIntent, error)
}

type payment struct {
	stripeSecretKey string
	successUrl      string
	cancelUrl       string
}

func (p payment) CreatePayment(amount float64, userId uint, orderId string) (*stripe.PaymentIntent, error) {
	stripe.Key = p.stripeSecretKey
	amountCents := amount * 100

	params := &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(int64(amountCents)),
		Currency:           stripe.String(string(stripe.CurrencyBRL)),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
	}

	//params := &stripe.CheckoutSessionParams{
	//	PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
	//	LineItems: []*stripe.CheckoutSessionLineItemParams{
	//		{
	//			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
	//				UnitAmount: stripe.Int64(int64(amountinCents)),
	//				Currency:   stripe.String("brl"),
	//				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
	//					Name: stripe.String("Eletronics"),
	//				},
	//			},
	//			Quantity: stripe.Int64(1),
	//		},
	//	},
	//	Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
	//	SuccessURL: stripe.String(p.successUrl),
	//	CancelURL:  stripe.String(p.cancelUrl),
	//}
	//checkoutSession, err := session.New(params)

	params.AddMetadata("order_id", fmt.Sprintf("%d", orderId))
	params.AddMetadata("user_id", fmt.Sprintf("%d", userId))
	pi, err := paymentintent.New(params)

	if err != nil {
		log.Printf("Error creating payment: %v", err)
		return nil, errors.New("payment create session failure")
	}

	return pi, nil
}

func (p payment) GetPaymentStatus(pId string) (*stripe.PaymentIntent, error) {
	stripe.Key = p.stripeSecretKey
	params := &stripe.PaymentIntentParams{}
	result, err := paymentintent.Get(pId, params)

	if err != nil {
		log.Printf("Error getting payment status: %v", err)
		return nil, errors.New("payment status failure")
	}

	return result, nil
}

func NewPaymentClent(stripeSecretKey string, successUrl string, cancelUrl string) PaymentClient {
	return &payment{
		stripeSecretKey: stripeSecretKey,
		successUrl:      successUrl,
		cancelUrl:       cancelUrl,
	}
}
