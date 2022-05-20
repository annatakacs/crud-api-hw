package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"

	_ "github.com/lib/pq"
)

var db *sql.DB

var (
	host     string
	port     int
	user     string
	password string
	name     string
	schema   string
	table    string
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

func DbConn(psqlInfo string) *sql.DB {
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
	query := fmt.Sprintf("SELECT * FROM %s.%s", schema, table)
	rows, _ := db.Query(query)
	defer rows.Close()
	//defer db.Close()
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
	query := fmt.Sprintf("SELECT * FROM %s.%s WHERE ID=%s", schema, table, queryId)
	row := db.QueryRow(query)
	var name, ingredients, description string
	var price float64
	var kcal, id int
	var spicy, vegan, glutenFree bool
	//defer db.Close()
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

func InsertMeal(meal Meal) (error, int) {
	price := fmt.Sprintf("%f", meal.Price)
	lastInsertId := 0
	query := fmt.Sprintf("INSERT INTO %s.%s (name, price, ingredients, spicy, vegan, gluten_free, description, kcal) VALUES ('%s', %s, '%s', %s, %s, %s, '%s', %s) RETURNING id", schema, table, meal.Name, price, meal.Ingredients, strconv.FormatBool(meal.Spicy), strconv.FormatBool(meal.Vegan), strconv.FormatBool(meal.GlutenFree), meal.Description, strconv.Itoa(meal.Kcal))
	_, err := db.Exec(query)
	_ = db.QueryRow(query).Scan(&lastInsertId)
	//defer db.Close()
	if err != nil {
		return err, -1
	} else {
		return nil, lastInsertId
	}
}

func UpdateMeal(queryId string, meal Meal) error {
	price := fmt.Sprintf("%f", meal.Price)
	query := fmt.Sprintf("UPDATE %s.%s SET name='%s', price=%s, ingredients='%s', spicy=%s, vegan=%s, gluten_free=%s, description='%s', kcal=%s WHERE id=%s", schema, table, meal.Name, price, meal.Ingredients, strconv.FormatBool(meal.Spicy), strconv.FormatBool(meal.Vegan), strconv.FormatBool(meal.GlutenFree), meal.Description, strconv.Itoa(meal.Kcal), queryId)
	_, err := db.Exec(query)
	//defer db.Close()
	if err != nil {
		return err
	} else {
		return nil
	}
}

func DeleteMeal(queryId string) error {
	query := fmt.Sprintf("DELETE FROM %s.%s WHERE id =%s", schema, table, queryId)
	_, err := db.Exec(query)
	if err != nil {
		return err
	} else {
		return nil
	}
}

// DB initialization functions
func ConfigureDb(dbHost string, dbPort int, dbUser string, dbPassword string, dbName string, dbSchema string, dbTable string, mainDB *sql.DB) {
	host = dbHost
	port = dbPort
	user = dbUser
	password = dbPassword
	name = dbName
	schema = dbSchema
	table = dbTable
	db = mainDB
}

func InitializeDb() {
	if db == nil {
		log.Println("Database doesn't exist. Creation in progress...")
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, "postgres")
		postgresDb, _ := sql.Open("postgres", psqlInfo)
		createDbQuery := fmt.Sprintf("CREATE DATABASE %s;", name)
		log.Println(createDbQuery)
		_, postgresErr := postgresDb.Exec(createDbQuery)
		if postgresErr != nil {
			log.Println("Failed to create database")
			log.Fatal(postgresErr)
		}
		createTable()
	}

	query := fmt.Sprintf("SELECT * FROM %s.%s", schema, table)
	_, table_check := db.Query(query)
	if table_check != nil {
		log.Println("Table is missing. Creating...")
		createTable()
	}
	log.Println("Successfully initialized database")
}

func createTable() {
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

func CleanUpDb() {
	db.Close()
}
