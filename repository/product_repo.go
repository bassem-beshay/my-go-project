package repository

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"log"
	"my-go-project/models"
	"time"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

type UserRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db}
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// ******************************get all product*************************************
func (r *ProductRepository) GetAllProducts() ([]models.Products, error) {
	query := `SELECT product_id, product_name, description, price, img_url FROM products`
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		log.Println("err failed excute sql", err)
		return nil, err
	}
	defer rows.Close()

	var products []models.Products
	for rows.Next() {
		var p models.Products
		if err := rows.Scan(&p.Product_id, &p.Product_name, &p.Description, &p.Price, &p.Img_url); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

// ******************************add product*************************************
func (r *ProductRepository) CreateProduct(p models.Products) error {
	query := `INSERT INTO products (product_id, product_name, description, price, img_url) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(context.Background(), query, p.Product_id, p.Product_name, p.Description, p.Price, p.Img_url)
	return err
}

// *****************************delete product**************************************
func (r *ProductRepository) DeleteProduct(id int) error {
	query := `DELETE FROM products WHERE product_id = $1`
	_, err := r.db.Exec(context.Background(), query, id)
	return err
}

// *****************************update product****************************************
func (r *ProductRepository) UpdateProduct(p models.Products) error {
	query := `UPDATE products SET product_name = $1, description = $2, price = $3 WHERE product_id = $4`
	_, err := r.db.Exec(context.Background(), query, p.Product_name, p.Description, p.Price, p.Product_id)
	return err
}

// *************************** Get product sales with filtration using user name***************************
func (r *ProductRepository) GetProductSales(username string) ([]models.Orders, error) {
	query := `SELECT o.order_id, o.create_at, u.user_name, p.product_name, op.quantity, (op.quantity * op.price_update) AS total_price
    FROM orders o
    JOIN order_product op ON o.order_id = op.order_id
    JOIN users u ON op.user_id = u.user_id
    JOIN products p ON op.product_id = p.product_id
    WHERE u.user_name = $1`

	rows, err := r.db.Query(context.Background(), query, username)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var sales []models.Orders
	for rows.Next() {
		var sale models.Orders
		if err := rows.Scan(&sale.Order_id, &sale.CreatedAt, &sale.Username, &sale.ProductName, &sale.Quantity, &sale.TotalPrice); err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		sales = append(sales, sale)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating rows:", err)
		return nil, err
	}

	return sales, nil
}

// **********************************************************************************
// *****************************get all user*******************************************
func (r *UserRepository) GetAllUsers() ([]models.Users, error) {
	query := `SELECT user_id, user_name FROM users`
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		log.Println("err in sql", err)
		return nil, err
	}
	defer rows.Close()

	var users []models.Users
	for rows.Next() {
		var u models.Users
		if err := rows.Scan(&u.User_id, &u.User_name); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// **************************sign up***********************************************
func (r *UserRepository) CreateUser(user models.Users) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `INSERT INTO users (user_id, user_name, password, email, phone, address, create_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err = r.db.Exec(context.Background(), query,
		user.User_id, user.User_name, string(hashedPassword),
		user.Email, user.Phone, user.Address, user.CreatedAt,
	)

	return err
}

// ************************login***************************************************
func (r *UserRepository) LoginUser(email, password string) (string, error) {
	var user models.Users

	query := `SELECT user_id, user_name, password FROM users WHERE email = $1`
	err := r.db.QueryRow(context.Background(), query, email).Scan(&user.User_id, &user.User_name, &user.Password)
	if err != nil {
		return "", errors.New("user is not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("password incorrect")
	}

	// jwt token
	claims := models.JWTClaims{
		UserID:   user.User_id,
		UserName: user.User_name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := []byte("mysecretkey")

	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", errors.New("failed to grnerate token")
	}

	return signedToken, nil
}

// ***************************add credit card*********************************
func (r *UserRepository) AddCreditCard(userID int, card models.CreditCard) error {
	query := `INSERT INTO credit_card (user_id, card_id, card_num) 
              VALUES ($1, $2, $3)`

	_, err := r.db.Exec(context.Background(), query,
		userID, card.Card_id, card.Card_num,
	)

	if err != nil {
		log.Println("Error inserting credit card:", err)
		return err
	}

	log.Println("Credit card added successfully for user:", userID)
	return nil
}

// ***********************delete credit card**********************************
func (r *UserRepository) DeleteCreditCard(cardID int) error {
	query := `DELETE FROM credit_card WHERE card_id = $1`
	_, err := r.db.Exec(context.Background(), query, cardID)

	if err != nil {
		log.Println("Error deleting credit card:", err)
		return err
	}

	log.Println("Credit card deleted successfully, Card ID:", cardID)
	return nil
}

// *********************add product in cart **********************************
func (r *UserRepository) AddCartProduct(cp models.CartProduct) error {
	query := `INSERT INTO cart_product (cp_id, cart_id, product_id, quantity) 
              VALUES ($1, $2, $3, $4)`

	_, err := r.db.Exec(context.Background(), query, cp.CP_id, cp.Cart_id, cp.Product_id, cp.Quantity)
	if err != nil {
		log.Println("Error inserting product into cart:", err)
		return err
	}

	log.Println("Product added to cart successfully")
	return nil
}
func (r *UserRepository) AddCart(userID int, cart models.Cart) error {
	query := `INSERT INTO cart (user_id, cart_id) VALUES ($1, $2)`
	_, err := r.db.Exec(context.Background(), query, userID, cart.Cart_id)
	if err != nil {
		log.Println("Error inserting cart:", err)
		return err
	}

	log.Println("Cart added successfully for user:", userID)
	return nil
}

// ***************************get cart *********************************
func (r *UserRepository) GetAllCart(userID int) ([]models.CartProduct, error) {

	query := `SELECT cp.cp_id, cp.cart_id, cp.product_id, cp.quantity 
			  FROM cart_product cp 
			  JOIN cart c ON cp.cart_id = c.cart_id 
			  WHERE c.user_id = $1`

	rows, err := r.db.Query(context.Background(), query, userID)
	if err != nil {
		log.Println("err failed to excute sql", err)
		return nil, err
	}
	defer rows.Close()
	var carts []models.CartProduct

	for rows.Next() {
		var p models.CartProduct

		if err := rows.Scan(&p.CP_id, &p.Cart_id, &p.Product_id, &p.Quantity); err != nil {
			log.Println("err , failed scan ", err)
			return nil, err
		}
		carts = append(carts, p)
	}

	if err := rows.Err(); err != nil {
		log.Println("err , failed read this rows", err)
		return nil, err
	}

	return carts, nil
}

// **********************add order *************************************
func (r *UserRepository) AddOrder(order models.Orders) error {
	query := `INSERT INTO orders (order_id, total_price ,status , create_at) VALUES ($1, $2 ,$3,$4)`
	_, err := r.db.Exec(context.Background(), query, order.Order_id, order.TotalPrice, order.Status, order.CreatedAt)
	if err != nil {
		log.Println("Error inserting order:", err)
		return err
	}

	log.Println("order added successfully ")
	return nil
}

// *******************details of order ************************************
func (r *UserRepository) AddOrderProduct(op models.OrderProduct) error {
	query := `INSERT INTO order_product (op_id , order_id , product_id ,quantity , price_update )
				VALUES($1,$2,$3,$4,$5)`

	_, err := r.db.Exec(context.Background(), query, op.OP_id, op.Order_id, op.Product_id, op.Quantity, op.Price_update)
	if err != nil {
		log.Println("Error inserting order_product:", err)
		return err
	}

	log.Println("order_product added successfully")
	return nil

}

// ********************get all order ************************************
func (r *UserRepository) GetHistory(userID int) ([]models.OrderProduct, error) {
	query := `SELECT op.op_id , op.product_id ,op.order_id ,op.quantity 
       ,p.product_name
		FROM order_product op
		JOIN products p ON p.product_id=op.product_id
		JOIN orders o ON o.order_id = op.order_id
		WHERE op.user_id = $1 
		AND o.status = 'completed';`

	rows, err := r.db.Query(context.Background(), query, userID)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var orderHistory []models.OrderProduct

	for rows.Next() {
		var order models.OrderProduct
		if err := rows.Scan(&order.OP_id, &order.Product_id, &order.Order_id, &order.Quantity, &order.ProductName); err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		orderHistory = append(orderHistory, order)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating rows:", err)
		return nil, err
	}

	return orderHistory, nil
}
