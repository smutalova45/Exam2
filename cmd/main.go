package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"main.go/api"
	"main.go/config"
	"main.go/storage/postgres"
)

func main() {
	cfg := config.Load()

	store, err := postgres.New(cfg)
	if err != nil {
		log.Fatalln("error while connecting to db err:", err.Error())
		return
	}
	defer store.CloseDB()

	api.New(store)
	fmt.Println("listening at port :8081")
	if err = http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalln("Server has stopped!", err.Error())
	}
}
