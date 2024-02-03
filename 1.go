// main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

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

func main() {
	app := &App{
		Users:    make(map[string]User),
		Products: make(map[string]Product),
		Carts:    make(map[string]ShoppingCart),
	}

	// Add some example products
	app.AddProduct(Product{ID: "1", Name: "Laptop", Price: 999.99})
	app.AddProduct(Product{ID: "2", Name: "Smartphone", Price: 599.99})
	app.AddProduct(Product{ID: "3", Name: "Headphones", Price: 79.99})

	r := mux.NewRouter()

	// Routes for user authentication
	r.HandleFunc("/register", app.RegisterUser).Methods("POST")
	r.HandleFunc("/login", app.LoginUser).Methods("POST")

	// Routes for products
	r.HandleFunc("/products", app.GetProducts).Methods("GET")

	// Routes for shopping cart
	r.HandleFunc("/cart", app.GetShoppingCart).Methods("GET")
	r.HandleFunc("/cart/add/{productID}", app.AddToCart).Methods("POST")

	port := ":8080"
	fmt.Printf("Server is running on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}

// RegisterUser handles user registration
func (a *App) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	a.UserMutex.Lock()
	defer a.UserMutex.Unlock()

	if _, exists := a.Users[newUser.Username]; exists {
		http.Error(w, "Username already exists", http.StatusBadRequest)
		return
	}

	// Simulate a simple ID generation (you might use a library or database-generated ID)
	newUser.ID = fmt.Sprintf("%d", len(a.Users)+1)

	// In a real-world scenario, you'd hash the password before saving it
	a.Users[newUser.Username] = newUser

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

// LoginUser handles user login
func (a *App) LoginUser(w http.ResponseWriter, r *http.Request) {
	var loginUser User
	err := json.NewDecoder(r.Body).Decode(&loginUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	a.UserMutex.RLock()
	defer a.UserMutex.RUnlock()

	user, exists := a.Users[loginUser.Username]
	if !exists || user.Password != loginUser.Password {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// GetProducts returns the list of products
func (a *App) GetProducts(w http.ResponseWriter, r *http.Request) {
	a.UserMutex.RLock()
	defer a.UserMutex.RUnlock()

	products := make([]Product, 0, len(a.Products))
	for _, product := range a.Products {
		products = append(products, product)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// GetShoppingCart returns the user's shopping cart
func (a *App) GetShoppingCart(w http.ResponseWriter, r *http.Request) {
	userID := "1" // In a real-world scenario, get the user ID from the authentication token

	a.CartMutex.RLock()
	defer a.CartMutex.RUnlock()

	cart, exists := a.Carts[userID]
	if !exists {
		cart = ShoppingCart{UserID: userID, Products: make([]Product, 0)}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

// AddToCart adds a product to the user's shopping cart
func (a *App) AddToCart(w http.ResponseWriter, r *http.Request) {
	userID := "1" // In a real-world scenario, get the user ID from the authentication token

	params := mux.Vars(r)
	productID := params["productID"]

	a.CartMutex.Lock()
	defer a.CartMutex.Unlock()

	product, exists := a.Products[productID]
	if !exists {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	cart, exists := a.Carts[userID]
	if !exists {
		cart = ShoppingCart{UserID: userID, Products: make([]Product, 0)}
	}

	cart.Products = append(cart.Products, product)
	a.Carts[userID] = cart

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cart)
}

// AddProduct adds a product to the list of products
func (a *App) AddProduct(product Product) {
	a.Products[product.ID] = product
}
