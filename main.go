package main

import (
	"Avito_go/db"
	"Avito_go/pkg/api"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func main() {
	db := db.CreateTables()
	http.HandleFunc("/income", api.UpdateBalance(db))
	http.HandleFunc("/balance", api.GetBalance(db))
	http.HandleFunc("/order", api.AccountOrder(db))
	log.Println("Server is working") // pseudo test
	if err := http.ListenAndServe(":8080", nil); err != http.ErrServerClosed {
		log.Println("ListenAndServe error")
		panic(err)
	}
}
