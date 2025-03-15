package dto

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserSignup struct {
	UserLogin
	Phone string `json:"phone"`
}

type VerificationCodeInput struct {
	Code string `json:"code"`
}

type SellerInput struct {
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	PhoneNumber       string `json:"phone_number"`
	BankAccountNumber uint   `json:"bank_account_number"`
	SwiftCode         string `json:"swift_code"`
	PaymentType       string `json:"payment_type"`
}

type AddressInput struct {
	AddressLine1 string `json:"address_line_1"`
	AddressLine2 string `json:"address_line_2"`
	City         string `json:"city"`
	PostalCode   int32  `json:"postal_code"`
	Country      string `json:"country"`
}

type ProfileInput struct {
	FirstName    string       `json:"first_name"`
	LastName     string       `json:"last_name"`
	AddressInput AddressInput `json:"address_input"`
}

type UserSignupResponse struct {
	Message string `json:"message" example:"login/register"`
	Token   string `json:"token" example:"token JWT"`
}

type ErrorResponse struct {
	Message string `json:"message" example:"Please provide valid inputs"`
	Reason  string `json:"reason,omitempty" example:"detalhes do erro se houver"`
}
