package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var cliente *mongo.Client

// Author é uma classe de identificação do autor do livro.
type Author struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
}

// Book é uma classe que identifica o livro e o escritor.
type Book struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Isbn   string             `json:"isbn,omitempty"   bson:"isbn,omitempty"`
	Title  string             `json:"title,omitempty"  bson:"title,omitempty"`
	Author *Author            `json:"author,omitempty" bson:"author,omitempty"`
}

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	cliente, _ = mongo.Connect(context.TODO(), clientOptions)

	router := mux.NewRouter()
	router.HandleFunc("/api/books", getBooks).Methods("GET")
	router.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/api/books", createBook).Methods("POST")
	//router.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	//router.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	params := mux.Vars(r)

	id, _ := primitive.ObjectIDFromHex(params["id"])

	var book Book
	collection := cliente.Database("BookStore").Collection("Books")
	err := collection.FindOne(context.TODO(), Book{ID: id}).Decode(&book)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	json.NewEncoder(w).Encode(book)
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var books []Book
	collection := cliente.Database("BookStore").Collection("Books")
	ctx := context.TODO()
	cursor, err := collection.Find(ctx, bson.M{})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
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
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	json.NewEncoder(w).Encode(books)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)

	collection := cliente.Database("BookStore").Collection("Books")
	result, _ := collection.InsertOne(context.TODO(), book)

	json.NewEncoder(w).Encode(result)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")

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
	//w.Header().Set("Content-Type", "application/json")
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
