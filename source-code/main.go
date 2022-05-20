package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/annatakacs/go-crud/pkg/database"
	"github.com/annatakacs/go-crud/pkg/routes"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var (
	dbHost     = os.Getenv("DBHOST")
	dbPort     = 5432
	dbUser     = os.Getenv("DBUSER")
	dbPassword = os.Getenv("DBPASSWORD")
	dbName     = os.Getenv("DBNAME")
	dbSchema   = os.Getenv("DBSCHEMA")
	dbTable    = os.Getenv("DBTABLE")
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	mainDB := database.DbConn(psqlInfo)
	database.ConfigureDb(dbHost, dbPort, dbUser, dbPassword, dbName, dbSchema, dbTable, mainDB)
	database.InitializeDb()
	defer database.CleanUpDb()
	router := mux.NewRouter()

	routes.Routing(router)
	log.Println("Server started up on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
