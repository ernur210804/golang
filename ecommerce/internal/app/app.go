package app

import (
	"encoding/json"
	"fmt"
	"goProject/internal/model"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

// App represents the e-commerce application
type App struct {
	Users     map[string]model.User
	Products  map[string]model.Product
	Carts     map[string]model.ShoppingCart
	UserMutex sync.RWMutex
	CartMutex sync.RWMutex
}

// NewApp initializes a new App instance
func NewApp() *App {
	return &App{
		Users:    make(map[string]model.User),
		Products: make(map[string]model.Product),
		Carts:    make(map[string]model.ShoppingCart),
	}
}

// RegisterUser handles user registration
func (a *App) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var newUser model.User
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

	newUser.ID = fmt.Sprintf("%d", len(a.Users)+1)
	a.Users[newUser.Username] = newUser

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

// LoginUser handles user login
func (a *App) LoginUser(w http.ResponseWriter, r *http.Request) {
	var loginUser model.User
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

	products := make([]model.Product, 0, len(a.Products))
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
		cart = model.ShoppingCart{UserID: userID, Products: make([]model.Product, 0)}
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
		cart = model.ShoppingCart{UserID: userID, Products: make([]model.Product, 0)}
	}

	cart.Products = append(cart.Products, product)
	a.Carts[userID] = cart

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cart)
}

// AddProduct adds a product to the list of products
func (a *App) AddProduct(product model.Product) {
	a.Products[product.ID] = product
}
