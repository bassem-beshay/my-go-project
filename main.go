package main

import (
	"fmt"
	"log"
	"my-go-project/db"
	"my-go-project/handlers"
	"my-go-project/repository"
	"my-go-project/routes"
	"net/http"
	"github.com/rs/cors"
)

func main() {
	//****************************connection**********************
	dbPool, err := db.ConnectDB()
	if err != nil {
		log.Fatal("failed to connect with db", err)
	}
	defer dbPool.Close()

	//****************************repository**********************
	productRepo := repository.NewProductRepository(dbPool)
	userRepo := repository.NewUserRepository(dbPool)

	//****************************handlers**********************
	productHandler := handlers.NewProductHandler(productRepo)
	userHandler := handlers.NewUserHandler(userRepo)

	//****************************routes**********************
	r := routes.SetupRoutes(productHandler, userHandler)
	//****************************allows**********************
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, 
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	//****************************run server**********************
	port := "8080"
	fmt.Printf("✅ Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
