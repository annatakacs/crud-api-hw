package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	_ "github.com/lib/pq"
)

type Meal struct {
	Id          int     `json:id`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Ingredients string  `json:"ingredients"`
	Spicy       bool    `json:"spicy"`
	Vegan       bool    `json:"vegan"`
	GlutenFree  bool    `json:"glutenFree"`
	Description string  `json:"description"`
	Kcal        int     `json:"kcal"`
}

type User struct {
	Id       int    `json:id`
	Name     string `json:name`
	Email    string `json:email`
	Password []byte `json:-`
}

var (
	dbHost     = os.Getenv("DBHOST")
	dbPort     = 5432
	dbUser     = os.Getenv("DBUSER")
	dbPassword = os.Getenv("DBPASSWORD")
	dbName     = os.Getenv("DBNAME")
	dbSchema   = os.Getenv("DBSCHEMA")
	dbTable    = os.Getenv("DBTABLE")
)

func DbConn() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println(err)
		return nil
	}
	err = db.Ping()
	if err != nil {
		log.Println(err)
		return nil
	}
	log.Println("Successfully connected to database")
	return db
}

func GetAllMeals() ([]Meal, error) {
	var mealSlice []Meal
	var meal Meal
	db := DbConn()
	query := fmt.Sprintf("SELECT * FROM %s.%s", dbSchema, dbTable)
	rows, _ := db.Query(query)
	defer rows.Close()
	defer db.Close()
	for rows.Next() {
		var name, ingredients, description string
		var price float64
		var kcal, id int
		var spicy, vegan, glutenFree bool
		err := rows.Scan(&id, &name, &price, &ingredients, &spicy, &vegan, &glutenFree, &description, &kcal)
		if err == sql.ErrNoRows {
			return nil, err
		} else {
			meal.Id = id
			meal.Name = name
			meal.Price = price
			meal.Ingredients = ingredients
			meal.Spicy = spicy
			meal.Vegan = vegan
			meal.GlutenFree = glutenFree
			meal.Description = description
			meal.Kcal = kcal
			mealSlice = append(mealSlice, meal)
		}
	}
	return mealSlice, nil
}

func GetMeal(queryId string) (Meal, error) {
	var meal Meal
	db := DbConn()
	query := fmt.Sprintf("SELECT * FROM %s.%s WHERE ID=%s", dbSchema, dbTable, queryId)
	row := db.QueryRow(query)
	var name, ingredients, description string
	var price float64
	var kcal, id int
	var spicy, vegan, glutenFree bool
	defer db.Close()
	err := row.Scan(&id, &name, &price, &ingredients, &spicy, &vegan, &glutenFree, &description, &kcal)
	if err == sql.ErrNoRows {
		return meal, err
	} else {
		meal.Id = id
		meal.Name = name
		meal.Price = price
		meal.Ingredients = ingredients
		meal.Spicy = spicy
		meal.Vegan = vegan
		meal.GlutenFree = glutenFree
		meal.Description = description
		meal.Kcal = kcal
		return meal, nil
	}
}

func InsertMeal(meal Meal) error {
	db := DbConn()
	price := fmt.Sprintf("%f", meal.Price)
	query := fmt.Sprintf("INSERT INTO %s.%s (name, price, ingredients, spicy, vegan, gluten_free, description, kcal) VALUES ('%s', %s, '%s', %s, %s, %s, '%s', %s)", dbSchema, dbTable, meal.Name, price, meal.Ingredients, strconv.FormatBool(meal.Spicy), strconv.FormatBool(meal.Vegan), strconv.FormatBool(meal.GlutenFree), meal.Description, strconv.Itoa(meal.Kcal))
	_, err := db.Exec(query)
	defer db.Close()
	if err != nil {
		return err
	} else {
		return nil
	}
}

func UpdateMeal(queryId string, meal Meal) error {
	db := DbConn()
	price := fmt.Sprintf("%f", meal.Price)
	query := fmt.Sprintf("UPDATE %s.%s SET name='%s', price=%s, ingredients='%s', spicy=%s, vegan=%s, gluten_free=%s, description='%s', kcal=%s WHERE id=%s", dbSchema, dbTable, meal.Name, price, meal.Ingredients, strconv.FormatBool(meal.Spicy), strconv.FormatBool(meal.Vegan), strconv.FormatBool(meal.GlutenFree), meal.Description, strconv.Itoa(meal.Kcal), queryId)
	_, err := db.Exec(query)
	defer db.Close()
	if err != nil {
		return err
	} else {
		return nil
	}
}

func DeleteMeal(queryId string) error {
	db := DbConn()
	query := fmt.Sprintf("DELETE FROM %s.%s WHERE id =%s", dbSchema, dbTable, queryId)
	_, err := db.Exec(query)
	if err != nil {
		return err
	} else {
		return nil
	}
}

// DB initialization functions

func InitializeDb() {
	db := DbConn()
	if db == nil {
		log.Println("Database doesn't exist. Creation in progress...")
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, "postgres")
		postgresDb, _ := sql.Open("postgres", psqlInfo)
		createDbQuery := fmt.Sprintf("CREATE DATABASE %s;", dbName)
		log.Println(createDbQuery)
		_, postgresErr := postgresDb.Exec(createDbQuery)
		if postgresErr != nil {
			log.Println("Failed to create database")
			log.Fatal(postgresErr)
		}
		createTable()
	}

	query := fmt.Sprintf("SELECT * FROM %s.%s", dbSchema, dbTable)
	_, table_check := db.Query(query)
	if table_check != nil {
		log.Println("Table is missing. Creating...")
		createTable()
	}
	log.Println("Successfully initialized database")
}

func createTable() {
	db := DbConn()
	path := filepath.Join("./", "meals.sql")

	c, ioErr := ioutil.ReadFile(path)
	if ioErr != nil {
		log.Println("Failed to read sql file")
		log.Fatal(ioErr)
	}
	sql := string(c)
	_, execErr := db.Exec(sql)
	if execErr != nil {
		log.Println("Failed to exec db migration script")
		log.Fatal(execErr)
	}
}
