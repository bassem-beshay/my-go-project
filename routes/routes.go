package routes

import (
	"github.com/gorilla/mux"
	"my-go-project/handlers"
	"my-go-project/middlewares"
	"net/http"
)

func SetupRoutes(productHandler *handlers.ProductHandler, userHandler *handlers.UserHandler) *mux.Router {
	r := mux.NewRouter()

	// ✅ endpoint for admin
	productRoutes := r.PathPrefix("/products").Subrouter()
	productRoutes.HandleFunc("", productHandler.GetProductsHandler).Methods("GET")
	productRoutes.HandleFunc("", productHandler.CreateProductHandler).Methods("POST")
	productRoutes.HandleFunc("/{id}", productHandler.UpdateProductHandler).Methods("PUT")
	productRoutes.HandleFunc("/{id}", productHandler.DeleteProductHandler).Methods("DELETE")
	productRoutes.HandleFunc("/admin/{username}", productHandler.GetProductSalesHandler).Methods("GET")

	// ✅ endpoint for users
	userRoutes := r.PathPrefix("/users").Subrouter()
	userRoutes.HandleFunc("", userHandler.GetUsersHandler).Methods("GET")
	userRoutes.HandleFunc("/register", userHandler.RegisterHandler).Methods("POST")
	userRoutes.HandleFunc("/login", userHandler.LoginHandler).Methods("POST")
	userRoutes.Handle("/cart", middlewares.JWTMiddleware(http.HandlerFunc(userHandler.AddCartHandler))).Methods("POST")
	userRoutes.Handle("/cart", middlewares.JWTMiddleware(http.HandlerFunc(userHandler.GetCartHandler))).Methods("GET")
	userRoutes.HandleFunc("/addOrder-product", userHandler.AddOrderProductHandler).Methods("POST")
	userRoutes.HandleFunc("/addOrder", userHandler.AddOrderHandler).Methods("POST")
	// ✅ product in  cart
	r.HandleFunc("/addProduct-cart", userHandler.AddCartProductHandler).Methods("POST")

	// ✅ endpoint for users after login
	r.Handle("/add-credit", middlewares.JWTMiddleware(http.HandlerFunc(userHandler.AddCreditCardHandler))).Methods("POST")
	r.HandleFunc("/credit/{card_id}", userHandler.DeleteCreditCardHandler).Methods("DELETE")

	// ✅ endpoint history order
	r.Handle("/history", middlewares.JWTMiddleware(http.HandlerFunc(userHandler.GetHistoryOrder))).Methods("GET")

	return r
}
