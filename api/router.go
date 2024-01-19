package api

import (
	"net/http"

	"main.go/api/handler"
	"main.go/storage"
)

func New(store storage.IStorage) {
	h := handler.New(store)
	http.HandleFunc("/city", h.City)
	http.HandleFunc("/customer", h.Customer)
	http.HandleFunc("/driver", h.Driver)
	http.HandleFunc("/car", h.Car)
	http.HandleFunc("/trip", h.Trip)
}
