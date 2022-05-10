package main

import (
	"log"
	"net/http"

	"github.com/annatakacs/go-crud/pkg/database"
	"github.com/annatakacs/go-crud/pkg/routes"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	database.InitializeDb()
	router := mux.NewRouter()

	routes.Routing(router)
	log.Println("Server started up on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
