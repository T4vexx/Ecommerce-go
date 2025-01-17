package domain

import "time"

const (
	SELLER = "seller"
	BUYER  = "buyer"
)

type User struct {
	ID          uint        `json:"id" gorm:"PrimaryKey"`
	FirstName   string      `json:"first_name"`
	LastName    string      `json:"last_name"`
	Email       string      `json:"email" gorm:"index;unique;not null"`
	Phone       string      `json:"phone"`
	Password    string      `json:"password"`
	Address     Address     `json:"address"`
	BankAccount BankAccount `json:"bank_account"`
	Cart        []Cart      `json:"cart"`
	Order       []Order     `json:"order"`
	Payment     []Payment   `json:"payment"`
	Code        string      `json:"code"`
	Expiry      time.Time   `json:"expiry"`
	Verified    bool        `json:"verified" gorm:"default:false"`
	UserType    string      `json:"user_type" gorm:"default:buyer"`
	CreatedAt   time.Time   `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt   time.Time   `json:"updated_at" gorm:"default:current_timestamp"`
}
