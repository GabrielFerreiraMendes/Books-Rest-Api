package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var books []Book
var client *mongo.Client

type Book struct {
	ID     primitive.ObjectID `json:"id,omitempty"      bson: "id, omitempty"`
	Isbn   string             `json: "isbn,omitempty"   bson: "isbn,omitempty"`
	Title  string             `json: "title,omitempty"  bson: "title,omitempty"`
	Author *Author            `json: "author,omitempty" bson: "author,omitempty"`
}

type Author struct {
	FirstName string `json: "firstname,omitempty" bson: "firstname,omitempty"`
	LastName  string `json: "lastname,omitempty"  bson: "astname,omitempty"`
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//params := mux.Vars(r)

	var books []Book

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:8080"))

	collection := client.Database("BookStore").Collection("Books")

	cursor, err := collection.Find(ctx, bson.M{})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `"}`))
		return
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var book Book
		cursor.Decode(&book)
		books = append(books, book)
	}

	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `"}`))
		return
	}
	//for _, item := range books {
	//	if item.ID == params["id"] {
	//		json.NewEncoder(w).Encode(item)
	//		return
	//	}
	//}

	json.NewEncoder(w).Encode(&Book{})
}

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:8080"))

	collection := client.Database("BookStore").Collection("Books")
	result, _ := collection.InsertOne(ctx, book)

	json.NewEncoder(w).Encode(result)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//params := mux.Vars(r)

	//for index, item := range books {
	//	if item.ID == params["id"] {
	//		var book Book
	//
	//		_ = json.NewDecoder(r.Body).Decode(&book)
	//		book.ID = item.ID
	//
	//		books[index] = book
	//		json.NewEncoder(w).Encode(book)
	//		return
	//	}
	//}
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//
	//params := mux.Vars(r)
	//
	//for index, item := range books {
	//	if item.ID == params["id"] {
	//		books = append(books[:index], books[index+1:]...)
	//		break
	//	}
	//}
	//
	//json.NewEncoder(w).Encode(books)
}

func main() {
	//books = append(books, Book{ID: 1, Isbn: "123", Title: "First book",
	//	Author: &Author{FirstName: "John", LastName: "Doe"}})
	//
	//books = append(books, Book{ID: 2, Isbn: "456", Title: "Second book",
	//	Author: &Author{FirstName: "Johnathan", LastName: "Dude"}})
	//
	//books = append(books, Book{ID: 3, Isbn: "789", Title: "Third book",
	//	Author: &Author{FirstName: "Joe", LastName: "Durelo"}})
	//
	//books = append(books, Book{ID: 4, Isbn: "147", Title: "Fourth book",
	//	Author: &Author{FirstName: "James", LastName: "Gun"}})

	//ctx, _ := context.WithTimeout(context.Background(), 10+time.Second)
	//client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:8080"))

	router := mux.NewRouter()
	router.HandleFunc("/api/books", getBooks).Methods("GET")
	//router.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/api/books", createBook).Methods("POST")
	//router.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	//router.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}
