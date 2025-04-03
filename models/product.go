package models

import "time"
import "github.com/golang-jwt/jwt/v5"

type Products struct {
	Product_id   int
	Product_name string
	Description  string
	Price        int
	Img_url      string
}

type Cart struct {
	Cart_id int
	User_id int
}

type CartProduct struct {
	CP_id      int
	Cart_id    int
	Product_id int
	Quantity   int
}

type CreditCard struct {
	Card_id  int
	User_id  int
	Card_num string
}

type OrderProduct struct {
	OP_id        int
	Order_id     int
	Product_id   int
	Quantity     int
	Price_update int

	ProductName string
}

type Orders struct {
	Order_id   int
	TotalPrice int
	Status     string
	CreatedAt  time.Time
	ProductName string
	Username string
	Quantity int
}

type Payment struct {
	Payment_id int
	Order_id   int
	Card_id    int
	Method     string
	Amount     int
	Currency   string
	Status     string
	CreatedAt  time.Time
}

type Users struct {
	User_id   int
	User_name string
	Email     string
	Password  string
	Phone     string
	Address   string
	CreatedAt time.Time
}

type JWTClaims struct {
	UserID   int
	UserName string
	jwt.RegisteredClaims
}
