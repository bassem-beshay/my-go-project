# Go E-Commerce API

This project is a Go-based backend for an e-commerce platform. It provides RESTful APIs for managing products, users, orders, and integrates Stripe for payment processing. The project uses PostgreSQL as the database, and Docker is used for both the application and the PostgreSQL database setup.

## Requirements

- **Go programming language** 
- **PostgreSQL** for the database
- **GorillaMux** for implementing REST APIs
- **Docker** for both the application and PostgreSQL setup
- **Postman collection** for API endpoints
- **No ORM**: Uses plain SQL for database interactions

## Features

## Endpoints

### Admin Endpoints
- **Create Product**  
  `POST /products`  
  Create a new product.
  
- **Update Product**  
  `PUT /products/{id}`  
  Update an existing product by ID.
  
- **Delete Product**  
  `DELETE /products/{id}`  
  Delete an existing product by ID.

- **Get Product Sales**  
  `GET /products/admin/{username}`  


### User Endpoints
- **Get All Users**  
  `GET /users`  
  Retrieve a list of all users.

- **Register User**  
  `POST /users/register`  
  Register a new user.

- **Login User**  
  `POST /users/login`  
  Login an existing user and get a JWT token.

- **Add Product to Cart**  
  `POST /users/cart`  
  Add a product to the user's cart.

- **Get Cart Items**  
  `GET /users/cart`  
  Retrieve all products in the user's cart.

- **Add Order Product**  
  `POST /users/addOrder-product`  
  Add products to an order.

- **Add Order**  
  `POST /users/addOrder`  
  Complete an order by adding the products to the user's order history.

- **Add Credit Card**  
  `POST /users/add-credit`  
  Add a credit card to the user's account.

- **Delete Credit Card**  
  `DELETE /users/credit/{card_id}`  
  Delete a specific credit card by ID.

- **Get Order History**  
  `GET /users/history`  
  Retrieve the user's order history.

### Product in Cart
- **Add Product to Cart**  
  `POST /addProduct-cart`  
  Add a product to the cart (this endpoint is user-specific).
## Setup and Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/bassem-beshay/my-go-project.git
   cd my-go-project
