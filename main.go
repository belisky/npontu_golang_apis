//Name:Fiawornu Nobel and Paapa Quansah
//Position:Intern@npontu technologies
//School:KNUST
//Email:denoblesnobility2@gmail.com
//Rest API with golang using MUX router and Mysql DB--(G-ORM)..

package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//initialising golang-object-relational-Model..
var tpl *template.Template
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

//Author Struct (Model)
type Author struct {
	gorm.Model
	ID        string `json:"ID"`
	Firstname string `json:"firstname" `
	Lastname  string `json:"lastname"`
}

//Get All Books
func getBooks(w http.ResponseWriter, r *http.Request) {

	var query []Book

	var authors []Author
	DB.Order("created_at").Find(&query)
	DB.Order("created_at").Find(&authors)

	template := template.Must(template.ParseFiles("template/books.html"))
	template.ExecuteTemplate(w, "books.html", query)

	//***********************api********************
	//w.Header().Set("Content-Type", "application/json")
	//for loop to get a book and its corresponding author
	//because of two tables linked by foreign key dependency.
	//for i, b := range query {
	//complete := Book{ID: b.ID, Isbn: b.Isbn, Title: b.Title, Author: authors[i]}
	//json.NewEncoder(w).Encode(complete)
	//fmt.Println(complete)
	//}

}

//Get Single Book

func getBook(w http.ResponseWriter, r *http.Request) {

	fmt.Println("**********Getting a Book**********")
	r.ParseForm()
	id := r.FormValue("id")
	var book Book
	var author Author
	DB.First(&book, id)
	DB.First(&author, id)
	complete := Book{ID: book.ID, Isbn: book.Isbn, Title: book.Title, Author: author}

	template := template.Must(template.ParseFiles("template/book.html"))
	template.ExecuteTemplate(w, "book.html", complete)

	//****api***
	//w.Header().Set("Content-Type", "application/json")
	//complete := Book{ID: book.ID, Isbn: book.Isbn, Title: book.Title, Author: author1}
	//json.NewEncoder(w).Encode(complete)
	//json.NewEncoder(w).Encode(author1)
	//fmt.Println(book, author1)
	//params := mux.Vars(r)

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

//*****Inserting a new Book**********************

func insert(w http.ResponseWriter, r *http.Request) {
	fmt.Println("***********inserting a book**********")
	if r.Method == "GET" {
		template := template.Must(template.ParseFiles("template/insert.html"))
		template.ExecuteTemplate(w, "insert.html", nil)
		return
	}
	var book Book
	r.ParseForm()
	book.Title = r.FormValue("titleTitle")
	book.Isbn = r.FormValue("isbnIsbn")
	book.ID = r.FormValue("idId")
	book.Author.Firstname = r.FormValue("authorFn")
	book.Author.Lastname = r.FormValue("authorLn")

	var err error
	if book.Title == "" || book.Isbn == "" || book.ID == "" || book.Author.Firstname == "" || book.Author.Lastname == "" {
		fmt.Println("Error inserting row:", err)
		template := template.Must(template.ParseFiles("template/insert.html"))
		template.ExecuteTemplate(w, "insert.html", "Error inserting data,Please check all fields.")
		return
	}
	DB.Create(&book)
	template := template.Must(template.ParseFiles("template/insert.html"))   //Parse the html file
	template.ExecuteTemplate(w, "insert.html", "Book Added Successfully!!!") //Execute the html file

}

//************Update a Book****************

func updateBook(w http.ResponseWriter, r *http.Request) {

	fmt.Println("**********Updating Book**********")
	r.ParseForm()
	id := r.FormValue("id")
	var book Book
	var author Author
	DB.First(&book, id)
	DB.First(&author, id)
	complete := Book{ID: book.ID, Isbn: book.Isbn, Title: book.Title, Author: author}
	template := template.Must(template.ParseFiles("template/update.html"))
	template.ExecuteTemplate(w, "update.html", complete)

	//******************************api******************************
	//w.Header().Set("Content-Type", "application/json")
	//params := mux.Vars(r)
	//json.NewDecoder(r.Body).Decode(&book)
	//json.NewDecoder(r.Body).Decode(&author)
	//DB.Save(&book)
	//json.NewEncoder(w).Encode(book)

}

func updateresult(w http.ResponseWriter, r *http.Request) {
	fmt.Println("**********updating result**********")
	var book Book
	var author Author
	r.ParseForm()
	book.ID = r.FormValue("id")
	book.Title = r.FormValue("titleTitle")
	book.Isbn = r.FormValue("isbnIsbn")
	author.Firstname = r.FormValue("authorFn")
	author.Lastname = r.FormValue("authorLn")
	DB.Save(&book)
	DB.Save(&author)
	template := template.Must(template.ParseFiles("template/result.html"))
	template.ExecuteTemplate(w, "result.html", "Book Updated Successfully!!!")

}

//Delete a Book

func deleteBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("**********Deleting Book***********")
	r.ParseForm()
	id := r.FormValue("id")

	var book Book
	var author Author

	DB.Delete(&book, id)
	DB.Delete(&author, id)
	template := template.Must(template.ParseFiles("template/result.html"))
	template.ExecuteTemplate(w, "result.html", "Book Deleted Successfully!!!")
	//************************api************************
	//w.Header().Set("Content-Type", "application/json")
	//params := mux.Vars(r)
	//DB.First(&book, params["id"])
	//DB.First(&author, params["id"])
	//json.NewEncoder(w).Encode("The Book was deleted successfully")
	//fmt.Println(book)

}

func InitialiseRouter() {

	//init router
	r := mux.NewRouter()
	r.HandleFunc("/", landingpage)
	//Route Handlers /Endpoints
	r.HandleFunc("/api/create", insert)
	r.HandleFunc("/api/books", getBooks)
	r.HandleFunc("/api/book/", getBook)
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/update/", updateBook)
	r.HandleFunc("/api/updateresult/", updateresult)
	r.HandleFunc("/api/delete/", deleteBook)
	log.Fatal(http.ListenAndServe(":8000", r))

}
func landingpage(w http.ResponseWriter, r *http.Request) {
	template := template.Must(template.ParseFiles("template/home.html"))
	template.ExecuteTemplate(w, "home.html", "")
}

func main() {

	InitialMigration()
	InitialiseRouter()

}
