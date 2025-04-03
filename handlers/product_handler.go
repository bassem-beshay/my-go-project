package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"my-go-project/models"
	"my-go-project/repository"
	"net/http"
	"strconv"
)

type ProductHandler struct {
	repo *repository.ProductRepository
}

type UserHandler struct {
	repo *repository.UserRepository
}

func NewProductHandler(repo *repository.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// ***********************get products*********************************************
func (h *ProductHandler) GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	products, err := h.repo.GetAllProducts()
	if err != nil {
		http.Error(w, "Failed to retrieve products", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// **********************add product **********************************************
func (h *ProductHandler) CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	var p models.Products
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if p.Product_name == "" || p.Price <= 0 {
		http.Error(w, "Invalid product data", http.StatusBadRequest)
		return
	}

	err := h.repo.CreateProduct(p)
	if err != nil {
		http.Error(w, "Failed to add product", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

// ***********************delete product *************************************
func (h *ProductHandler) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}
	err = h.repo.DeleteProduct(id)
	if err != nil {
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// **********************update product ***************************************
func (h *ProductHandler) UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	var p models.Products
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if p.Product_id == 0 || p.Product_name == "" || p.Price <= 0 {
		http.Error(w, "Invalid product data", http.StatusBadRequest)
		return
	}

	err := h.repo.UpdateProduct(p)
	if err != nil {
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// **************************get product  by  username****************************
func (h *ProductHandler) GetProductSalesHandler(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	sales, err := h.repo.GetProductSales(username)
	if err != nil {
		http.Error(w, "Failed to retrieve sales", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sales)
}

// *************************************************************************************************
// *******************get all users **********************************
func (h *UserHandler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := h.repo.GetAllUsers()
	if err != nil {
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

//**********************sign up ***********************************

func (h *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.Users
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if user.User_name == "" || user.Password == "" || user.Email == "" || user.Address == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	err := h.repo.CreateUser(user)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully", "username": user.User_name})
}

// ************************login*****************************
func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginData struct {
		Email    string
		Password string
	}

	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if loginData.Email == "" || loginData.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	token, err := h.repo.LoginUser(loginData.Email, loginData.Password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// give me token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// *********************add credit card *******************************
func (h *UserHandler) AddCreditCardHandler(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "User ID not found in token", http.StatusUnauthorized)
		return
	}

	var card models.CreditCard
	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if card.Card_id == 0 || card.User_id == 0 || card.Card_num == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	err := h.repo.AddCreditCard(userID, card)
	if err != nil {
		http.Error(w, "Failed to add credit card", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Credit card added successfully"})
}

// *****************************delete credit card ********************************
func (h *UserHandler) DeleteCreditCardHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cardID, err := strconv.Atoi(vars["card_id"])
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	err = h.repo.DeleteCreditCard(cardID)
	if err != nil {
		http.Error(w, "Failed to delete credit card", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Credit card deleted successfully"))
}

// ************************add product in cart ***************************************
func (h *UserHandler) AddCartProductHandler(w http.ResponseWriter, r *http.Request) {
	// decode
	var cartProduct models.CartProduct
	if err := json.NewDecoder(r.Body).Decode(&cartProduct); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if cartProduct.CP_id == 0 || cartProduct.Cart_id == 0 || cartProduct.Product_id == 0 || cartProduct.Quantity <= 0 {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	err := h.repo.AddCartProduct(cartProduct)
	if err != nil {
		http.Error(w, "Failed to add product to cart", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": " cartadded successfully"})
}

//***************add cart ***************************************

func (h *UserHandler) AddCartHandler(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "User ID not found in token", http.StatusUnauthorized)
		return
	}

	var cart models.Cart
	if err := json.NewDecoder(r.Body).Decode(&cart); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if cart.Cart_id == 0 {
		http.Error(w, "Cart ID cannot be zero", http.StatusBadRequest)
		return
	}

	err := h.repo.AddCart(userID, cart)
	if err != nil {
		http.Error(w, "Failed to add cart", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Cart added successfully"})
}

// *********************get cart *****************************************
func (h *UserHandler) GetCartHandler(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "User ID not found in token", http.StatusUnauthorized)
		return
	}

	carts, err := h.repo.GetAllCart(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve cart products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(carts)
}

// ********************add order***************************************
func (h *UserHandler) AddOrderHandler(w http.ResponseWriter, r *http.Request) {

	var order models.Orders
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if order.Order_id == 0 || order.TotalPrice == 0 || order.Status == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	err := h.repo.AddOrder(order)
	if err != nil {
		http.Error(w, "Failed to add order ", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": " ordered successfully"})
}

// ************************details of order****************************
func (h *UserHandler) AddOrderProductHandler(w http.ResponseWriter, r *http.Request) {

	var orderProduct models.OrderProduct
	if err := json.NewDecoder(r.Body).Decode(&orderProduct); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if orderProduct.OP_id == 0 || orderProduct.Order_id == 0 || orderProduct.Product_id == 0 || orderProduct.Quantity <= 0 || orderProduct.Price_update == 0 {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	err := h.repo.AddOrderProduct(orderProduct)
	if err != nil {
		http.Error(w, "Failed to add order to product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": " ordered product successfully"})
}

// ********************************get all orders********************************
func (h *UserHandler) GetHistoryOrder(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "User ID not found in token", http.StatusUnauthorized)
		return
	}

	orders, err := h.repo.GetHistory(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve cart products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}
