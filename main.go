package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var books []Book

type Book struct {
	ID     string  "json: id"
	Isbn   string  "json: isbn"
	Title  string  "json: title"
	Author *Author "json: author"
}

type Author struct {
	FirstName string "json: firstname"
	LastName  string "json: lastname"
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)

}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	json.NewEncoder(w).Encode(&Book{})
}

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(10))

	books = append(books, book)
	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	for index, item := range books {
		if item.ID == params["id"] {
			var book Book

			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = item.ID

			books[index] = book
			json.NewEncoder(w).Encode(book)
			return
		}
	}
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}

	json.NewEncoder(w).Encode(books)
}

func main() {
	books = append(books, Book{ID: "1", Isbn: "123", Title: "First book",
		Author: &Author{FirstName: "John", LastName: "Doe"}})

	books = append(books, Book{ID: "2", Isbn: "456", Title: "Second book",
		Author: &Author{FirstName: "Johnathan", LastName: "Dude"}})

	books = append(books, Book{ID: "3", Isbn: "789", Title: "Third book",
		Author: &Author{FirstName: "Joe", LastName: "Durelo"}})

	books = append(books, Book{ID: "4", Isbn: "147", Title: "Fourth book",
		Author: &Author{FirstName: "James", LastName: "Gun"}})

	router := mux.NewRouter()
	router.HandleFunc("/api/books", getBooks).Methods("GET")
	router.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/api/books", createBook).Methods("POST")
	router.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}
