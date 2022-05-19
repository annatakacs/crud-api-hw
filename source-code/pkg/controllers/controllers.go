package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/annatakacs/go-crud/pkg/database"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func RetrieveMeal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	log.Println(vars)
	queryId := vars["id"]
	log.Println(queryId)
	meal, err := database.GetMeal(queryId)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		message := "404 Not found"
		json.NewEncoder(w).Encode(message)
	} else {
		json.NewEncoder(w).Encode(meal)
		w.WriteHeader(http.StatusOK)
		log.Println("Successfully returned meal")
	}
}

func RetrieveMeals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var mealSlice []database.Meal
	mealSlice, err := database.GetAllMeals()
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		message := "404 Not found"
		json.NewEncoder(w).Encode(message)
	} else {
		json.NewEncoder(w).Encode(mealSlice)
		w.WriteHeader(http.StatusOK)
		log.Println("Successfully returned all meals")
	}
	fmt.Println(mealSlice)

}

func CreateMeal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var meal database.Meal
	err := json.NewDecoder(r.Body).Decode(&meal)
	if err != nil {
		message := "Failed to fetch input data"
		json.NewEncoder(w).Encode(message)
		log.Println(message)
	}
	dbErr, id := database.InsertMeal(meal)
	if dbErr != nil {
		message := "Failed to insert data"
		json.NewEncoder(w).Encode(message)
		log.Println(message)
	} else {
		fmt.Printf("%d", id)
		message := fmt.Sprintf("Succesfully inserted data with id %s", strconv.Itoa(id))
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(message)
		log.Println(message)
	}
}

func ModifyMeal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	queryId := vars["id"]
	var meal database.Meal
	err := json.NewDecoder(r.Body).Decode(&meal)
	if err != nil {
		message := "Failed to fetch input data"
		json.NewEncoder(w).Encode(message)
		log.Println(message)
	}

	dbError := database.UpdateMeal(queryId, meal)
	if dbError != nil {
		w.WriteHeader(http.StatusNotModified)
		output := fmt.Sprintf("Failed to update row")
		log.Println(output)
		json.NewEncoder(w).Encode(output)
	} else {
		w.WriteHeader(http.StatusOK)
		output := fmt.Sprintf("Succesfully updated row")
		log.Println(output)
		json.NewEncoder(w).Encode(output)
	}
}

func RemoveMeal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	queryId := vars["id"]
	var meal database.Meal
	err := json.NewDecoder(r.Body).Decode(&meal)
	if err != nil {
		message := "Failed to fetch input data"
		json.NewEncoder(w).Encode(message)
		log.Println(message)
	}

	dbError := database.DeleteMeal(queryId)
	if dbError != nil {
		output := fmt.Sprintf("Failed to delete row")
		log.Println(output)
		json.NewEncoder(w).Encode(output)
	} else {
		w.WriteHeader(http.StatusOK)
		output := fmt.Sprintf("Succesfully deleted row")
		log.Println(output)
		json.NewEncoder(w).Encode(output)
	}
}
