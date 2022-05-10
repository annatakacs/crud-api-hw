package routes

import (
	"github.com/annatakacs/go-crud/pkg/controllers"

	"github.com/gorilla/mux"
)

func Routing(router *mux.Router) {
	router.HandleFunc("/api/v1/get/meals", controllers.RetrieveMeals).Methods("GET")
	router.HandleFunc("/api/v1/get/meals/{id}", controllers.RetrieveMeal).Methods("GET")
	router.HandleFunc("/api/v1/post/meals", controllers.CreateMeal).Methods("POST")
	router.HandleFunc("/api/v1/put/meals/{id}", controllers.ModifyMeal).Methods("PUT")
	router.HandleFunc("/api/v1/delete/meals/{id}", controllers.RemoveMeal).Methods("DELETE")
}
