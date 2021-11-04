//Name:Fiawornu Nobel
//Position:Intern@npontu technologies
//School:KNUST
//Email:denoblesnobility2@gmail.com
//Rest API with golang using MUX router and Mysql DB--(G-ORM)..

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//initialising golang-object-relational-Model..
 
var DB *gorm.DB
var err error

const DNS = "root@tcp(127.0.0.1:3306)/goapi?charset=utf8mb4&parseTime=True"

func InitialMigration() {
	DB, err = gorm.Open(mysql.Open(DNS), &gorm.Config{}) //connecting to mysql db
	if err != nil {
		fmt.Println(err.Error())
		panic("Cannot connect to DB")
	}

	//migrating book and author tables.
	//DB.AutoMigrate(&Book{})
	//DB.AutoMigrate(&Author{})

}

// Book Struct (Model)

type Book struct {
	gorm.Model
	ID     string `json:"ID"`
	Isbn   string `json:"isbn"`
	Title  string `json:"title"`
	Author Author `json:"author"  gorm:"foreignKey:ID;references:ID"` //foreign key dependency.
}

//Author Struct
type Author struct {
	gorm.Model
	ID        string `json:"ID"`
	Firstname string `json:"firstname" `
	Lastname  string `json:"lastname"`
}

//Initial DB Migration

//Get All Books
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var query []Book

	var authors []Author
	DB.Order("created_at").Find(&query)
	DB.Order("created_at").Find(&authors)

	//for loop to get a book and its corresponding author
	//because of two tables linked by foreign key dependency.
	for i, b := range query {

		complete := Book{ID: b.ID, Isbn: b.Isbn, Title: b.Title, Author: authors[i]}
		json.NewEncoder(w).Encode(complete)
		fmt.Println(complete)

	}

}

//Get Single Book

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) //Get params
	var book Book
	var author1 Author
	DB.First(&book, params["id"])
	DB.First(&author1, params["id"])
	complete := Book{ID: book.ID, Isbn: book.Isbn, Title: book.Title, Author: author1}
	json.NewEncoder(w).Encode(complete)
	json.NewEncoder(w).Encode(author1)
	fmt.Println(book, author1)

}

//Create a New Book

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book1 Book
	json.NewDecoder(r.Body).Decode(&book1)
	DB.Create(&book1) //insert new book into db.
	json.NewEncoder(w).Encode(book1)
	fmt.Println(book1)

}

//Update

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var book Book
	var author Author
	DB.First(&book, params["id"])
	DB.First(&author, params["id"])
	json.NewDecoder(r.Body).Decode(&book)
	json.NewDecoder(r.Body).Decode(&author)
	DB.Save(&book)
	json.NewEncoder(w).Encode(book)
}

//Delete Book

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var book Book
	var author Author
	DB.First(&book, params["id"])
	DB.First(&author, params["id"])
	DB.Delete(&book, params["id"])
	DB.Delete(&author, params["id"])
	json.NewEncoder(w).Encode("The Book was deleted successfully")
	fmt.Println(book)
}

func InitialiseRouter() {
	//init router
	r := mux.NewRouter()

	//Route Handlers /Endpoints
	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", r))

}

func main() {

	InitialMigration()
	InitialiseRouter()

}
