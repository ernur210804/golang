package main

import (
	"fmt"
	"goProject/internal/app"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	appInstance := app.NewApp()
	r := mux.NewRouter()

	r.HandleFunc("/register", appInstance.RegisterUser).Methods("POST")
	r.HandleFunc("/login", appInstance.LoginUser).Methods("POST")
	r.HandleFunc("/products", appInstance.GetProducts).Methods("GET")
	r.HandleFunc("/cart", appInstance.GetShoppingCart).Methods("GET")
	r.HandleFunc("/cart/add/{productID}", appInstance.AddToCart).Methods("POST")

	port := ":8080"
	fmt.Printf("Server is running on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
