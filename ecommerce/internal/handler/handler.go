package handler

import (
	"encoding/json"
	"fmt"
	"goProject/internal/model"
	"net/http"

	"github.com/gorilla/mux"
)

// AppHandler handles HTTP requests
type AppHandler struct {
	App *model.App
}

// NewAppHandler creates a new instance of AppHandler
func NewAppHandler(app *model.App) *AppHandler {
	return &AppHandler{App: app}
}

// RegisterUser handles user registration
func (h *AppHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var newUser model.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.App.UserMutex.Lock()
	defer h.App.UserMutex.Unlock()

	if _, exists := h.App.Users[newUser.Username]; exists {
		http.Error(w, "Username already exists", http.StatusBadRequest)
		return
	}

	// Simulate a simple ID generation (you might use a library or database-generated ID)
	newUser.ID = fmt.Sprintf("%d", len(h.App.Users)+1)

	// In a real-world scenario, you'd hash the password before saving it
	h.App.Users[newUser.Username] = newUser

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

// LoginUser handles user login
func (h *AppHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var loginUser model.User
	err := json.NewDecoder(r.Body).Decode(&loginUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.App.UserMutex.RLock()
	defer h.App.UserMutex.RUnlock()

	user, exists := h.App.Users[loginUser.Username]
	if !exists || user.Password != loginUser.Password {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// GetProducts returns the list of products
func (h *AppHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	h.App.UserMutex.RLock()
	defer h.App.UserMutex.RUnlock()

	products := make([]model.Product, 0, len(h.App.Products))
	for _, product := range h.App.Products {
		products = append(products, product)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// GetShoppingCart returns the user's shopping cart
func (h *AppHandler) GetShoppingCart(w http.ResponseWriter, r *http.Request) {
	userID := "1" // In a real-world scenario, get the user ID from the authentication token

	h.App.CartMutex.RLock()
	defer h.App.CartMutex.RUnlock()

	cart, exists := h.App.Carts[userID]
	if !exists {
		cart = model.ShoppingCart{UserID: userID, Products: make([]model.Product, 0)}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

// AddToCart adds a product to the user's shopping cart
func (h *AppHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	userID := "1" // In a real-world scenario, get the user ID from the authentication token

	params := mux.Vars(r)
	productID := params["productID"]

	h.App.CartMutex.Lock()
	defer h.App.CartMutex.Unlock()

	product, exists := h.App.Products[productID]
	if !exists {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	cart, exists := h.App.Carts[userID]
	if !exists {
		cart = model.ShoppingCart{UserID: userID, Products: make([]model.Product, 0)}
	}

	cart.Products = append(cart.Products, product)
	h.App.Carts[userID] = cart

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cart)
}

// AddProduct adds a product to the list of products
func (h *AppHandler) AddProduct(product model.Product) {
	h.App.Products[product.ID] = product
}
