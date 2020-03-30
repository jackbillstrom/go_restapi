package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
)

type Contract struct {
	gorm.Model
	Title       string    `json:"title"`
	CompanyName string    `json:"company_name"`
	CompanyID   int       `json:"company_id"`
	Period      string    `json:"period"`
	Reciever    Reciever  `gorm:"embedded" json:"reciever`
	Sender      Sender    `gorm:"embedded" json:"sender`
	Description string    `gorm:"type:text" json:"description"`
	Products    []Product `gorm:"many2many:contract_products"  json:"products"`
}

type Product struct {
	gorm.Model
	Name        string `json:"name"`
	Comment     string `json:"comment"`
	Value       string `json:"value"`
	Period      string `json:"period"`
	Description string `json:"description"`
}

type Reciever struct {
	RecieverUserID string `json:"user_id"`
	RecieverName   string `json:"name"`
	RecieverPhone  string `json:"phone"`
	RecieverEmail  string `json:"email"`
}

type Sender struct {
	SenderUserID string `json:"user_id"`
	SenderName   string `json:"name"`
	SenderPhone  string `json:"phone"`
	SenderEmail  string `json:"email"`
}

func Migration() {
	db, err := gorm.Open("mssql", "sqlserver://jack:Q29ndavr@Gby+Ve@ampilio.database.windows.net:1433?database=ampilio-hive")
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to  connect")
	}
	defer db.Close()

	db.AutoMigrate(&Contract{}, &Product{})
}

func All(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	db, err := gorm.Open("mssql", "sqlserver://jack:Q29ndavr@Gby+Ve@ampilio.database.windows.net:1433?database=ampilio-hive")
	if err != nil {
		panic("Could not connect to DB")
	}
	defer db.Close()

	var contracts []Contract
	db.Find(&contracts)
	// products := []Product{}
	// db.Preload("Contracts").Find(&)
	// db.Model(&contract).Association("Products")
	json.NewEncoder(w).Encode(contracts)
}

func Create(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers for the main request.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var contract Contract
	decoder := json.NewDecoder(r.Body).Decode(&contract)
	db, err := gorm.Open("mssql", "sqlserver://jack:Q29ndavr@Gby+Ve@ampilio.database.windows.net:1433?database=ampilio-hive")
	if err != nil {
		panic("Could not connect to DB")
	}
	defer db.Close()

	fmt.Println(decoder)

	// var contract Contract
	// err = decoder.Decode(&contract)
	// if err != nil {
	// 	panic(err)
	// }

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	defer r.Body.Close()
	json.Unmarshal(data, &contract)
	db.Create(&contract)
	fmt.Fprintf(w, "Kontraktet skapades")
}

func handleRequests() {
	_router := mux.NewRouter().StrictSlash(true)
	// Default
	_router.HandleFunc("/", All).Methods("GET")
	_router.HandleFunc("/", Create).Methods("POST")
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(_router)))
}

func main() {
	fmt.Println("= API initieraSsss  🐍  på = " + os.Getenv("PORT"))

	fmt.Println("= Migration structs into tables")
	Migration()

	handleRequests()
}
