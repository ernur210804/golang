// model/app.go

package model

import "sync"

// User represents a user in the system
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Product represents a product in the e-commerce system
type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// ShoppingCart represents a user's shopping cart
type ShoppingCart struct {
	UserID   string    `json:"user_id"`
	Products []Product `json:"products"`
}

// App represents the e-commerce application
type App struct {
	Users     map[string]User
	Products  map[string]Product
	Carts     map[string]ShoppingCart
	UserMutex sync.RWMutex
	CartMutex sync.RWMutex
}

// NewApp initializes a new instance of the App
func NewApp() *App {
	return &App{
		Users:    make(map[string]User),
		Products: make(map[string]Product),
		Carts:    make(map[string]ShoppingCart),
	}
}
