package controllers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/annatakacs/go-crud/pkg/database"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
}

var (
	dbHost     = "localhost"
	dbPort     = 5432
	dbUser     = "postgres"
	dbPassword = "password"
	dbName     = "meals"
	dbSchema   = "meals_schema"
	dbTable    = "test_table"
)

func (suite *TestSuite) SetupSuite() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db := database.DbConn(psqlInfo)
	database.ConfigureDb(dbHost, dbPort, dbUser, dbPassword, dbName, dbSchema, dbTable, db)
	db.Exec("CREATE TABLE meals_schema.test_table (id SERIAL PRIMARY KEY, name varchar(30) NOT NULL, price numeric NOT NULL, ingredients varchar(50) NOT NULL, spicy boolean NOT NULL, vegan boolean NOT NULL, gluten_free boolean  NOT NULL, description varchar(50) NOT NULL, kcal int NOT NULL);")
}

func (suite *TestSuite) TearDownSuite() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db := database.DbConn(psqlInfo)
	db.Exec("DROP TABLE meals_schema.test_table;")
}

func (suite *TestSuite) SetupTest() {
}

func (suite *TestSuite) TearDownTest() {
}

func (suite *TestSuite) TestCreateMealStatusCode() {
	input := []byte(`{"Id":1,"name":"TEST Pizza","price":4.25,"ingredients":"mozarella, tomato sauce, basil","spicy":false,"vegan":false,"glutenFree":false,"description":"Delicious classic pizza","kcal":700}`)
	expected := 201
	r := httptest.NewRequest(http.MethodPost, "/api/v1/get/meals", bytes.NewBuffer(input))
	w := httptest.NewRecorder()
	CreateMeal(w, r)
	res := w.Result().StatusCode
	message := fmt.Sprintf("Expected status code: %s got %s instead", strconv.Itoa(expected), strconv.Itoa(res))
	suite.Equal(expected, res, message)
}

func (suite *TestSuite) TestRetrieveMeals() {
	r := httptest.NewRequest(http.MethodGet, "/api/v1/get/meals", nil)
	w := httptest.NewRecorder()
	RetrieveMeals(w, r)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		suite.Fail("Expected error to be nil got %v", err)
	}
	expected := `[{"Id":1,"name":"TEST Pizza","price":4.25,"ingredients":"mozarella, tomato sauce, basil","spicy":false,"vegan":false,"glutenFree":false,"description":"Delicious classic pizza","kcal":700},{"Id":2,"name":"TEST Pizza","price":4.25,"ingredients":"mozarella, tomato sauce, basil","spicy":false,"vegan":false,"glutenFree":false,"description":"Delicious classic pizza","kcal":700}]
`
	suite.Equal(expected, string(data))
}

func (suite *TestSuite) TestRetrievMealById() {
	expected := `{"Id":1,"name":"TEST Pizza","price":4.25,"ingredients":"mozarella, tomato sauce, basil","spicy":false,"vegan":false,"glutenFree":false,"description":"Delicious classic pizza","kcal":700}
`
	r := httptest.NewRequest(http.MethodGet, "/api/v1/get/meals/1", nil)
	w := httptest.NewRecorder()
	vars := map[string]string{
		"id": "1",
	}
	r = mux.SetURLVars(r, vars)
	RetrieveMeal(w, r)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		suite.Fail("Expected error to be nil got %s", err)
	}
	suite.Equal(expected, string(data))
}

func (suite *TestSuite) TestStatusCodeNonExistentID() {
	expected := 404

	r := httptest.NewRequest(http.MethodGet, "/api/v1/get/meals/0", nil)
	w := httptest.NewRecorder()
	vars := map[string]string{
		"id": "0",
	}
	r = mux.SetURLVars(r, vars)
	RetrieveMeal(w, r)
	res := w.Result().StatusCode
	message := fmt.Sprintf("Expected status code: %s got %s instead", strconv.Itoa(expected), strconv.Itoa(res))
	suite.Equal(expected, res, message)
}

func (suite *TestSuite) TestStatusCodeExistingID() {
	expected := 200
	r := httptest.NewRequest(http.MethodGet, "/api/v1/get/meals/1", nil)
	w := httptest.NewRecorder()
	vars := map[string]string{
		"id": "1",
	}
	r = mux.SetURLVars(r, vars)
	RetrieveMeal(w, r)
	res := w.Result().StatusCode
	message := fmt.Sprintf("Expected status code: %s got %s instead", strconv.Itoa(expected), strconv.Itoa(res))
	suite.Equal(expected, res, message)
}

func (suite *TestSuite) TestYUpdateMeal() {
	expected := 200
	input := []byte(`{"Id":1,"name":"TEST Pizza light","price":8.25,"ingredients":"mozarella, tomato sauce, basil","spicy":false,"vegan":false,"glutenFree":false,"description":"Delicious classic pizza","kcal":450}`)
	r := httptest.NewRequest(http.MethodGet, "/api/v1/put/meals/1", bytes.NewBuffer(input))
	w := httptest.NewRecorder()
	vars := map[string]string{
		"id": "1",
	}
	r = mux.SetURLVars(r, vars)
	ModifyMeal(w, r)
	res := w.Result().StatusCode
	message := fmt.Sprintf("Expected status code: %s got %s instead", strconv.Itoa(expected), strconv.Itoa(res))
	suite.Equal(expected, res, message)
}

func (suite *TestSuite) TestYYretrieveMeal() {
	expected := `{"Id":1,"name":"TEST Pizza light","price":8.25,"ingredients":"mozarella, tomato sauce, basil","spicy":false,"vegan":false,"glutenFree":false,"description":"Delicious classic pizza","kcal":450}
`
	r := httptest.NewRequest(http.MethodGet, "/api/v1/get/meals/1", nil)
	w := httptest.NewRecorder()
	vars := map[string]string{
		"id": "1",
	}
	r = mux.SetURLVars(r, vars)
	RetrieveMeal(w, r)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		suite.Fail("Expected error to be nil got %s", err)
	}
	suite.Equal(expected, string(data))
}

func (suite *TestSuite) TestZDeleteMeal() {
	expected := 200
	r := httptest.NewRequest(http.MethodGet, "/api/v1/delete/meals/1", nil)
	w := httptest.NewRecorder()
	vars := map[string]string{
		"id": "1",
	}
	r = mux.SetURLVars(r, vars)
	RemoveMeal(w, r)
	res := w.Result().StatusCode
	message := fmt.Sprintf("Expected status code: %s got %s instead", strconv.Itoa(expected), strconv.Itoa(res))
	suite.Equal(expected, res, message)
}

func TestCalculatorTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
