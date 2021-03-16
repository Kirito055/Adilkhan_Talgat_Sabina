package api

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v4/pgxpool"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Book type with Name, Author and ISBN
type Book struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	ISBN        string `json:"isbn"`
	Description string `json:"description,omitempty"`
}

var books = map[string]Book{
	"0345391802": Book{Title: "The Hitchhiker's Guide to the Galaxy", Author: "Douglas Adams", ISBN: "0345391802"},
	"0000000000": Book{Title: "Cloud Native Go", Author: "M.-Leander Reimer", ISBN: "0000000000"},
}

// ToJSON to be used for marshalling of Book type
func (b Book) ToJSON() []byte {
	ToJSON, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}
	return ToJSON
}

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	books    *BookModel
}

func FromDB() []*Book {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	pool, err := pgxpool.Connect(context.Background(), "user=postgres password=123 host=localhost port=5432 dbname=books sslmode=disable pool_max_conns=10")
	if err != nil {
		log.Fatalf("Unable to connection to database: %v\n", err)
	}
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		books:    &BookModel{DB: pool},
	}

	defer pool.Close()

	s, err := app.books.GetALL()
	if err != nil {
		panic(err)

	}

	return s
}

// FromJSON to be used for unmarshalling of Book type
func FromJSON(data []byte) Book {

	book := Book{}
	err := json.Unmarshal(data, &book)
	if err != nil {
		panic(err)
	}
	return book
}

// AllBooks returns a slice of all books
func AllBooks() []*Book {
	s := FromDB()
	return s
}

// AllBooks returns a slice of all books
//func AllBooks() []Book {
//	values := make([]Book, len(books))
//	idx := 0
//	for _, book := range books {
//		values[idx] = book
//		idx++
//	}
//	return values
//}

// BooksHandleFunc to be used as http.HandleFunc for Book API
func BooksHandleFunc(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodGet:
		books := FromDB()
		writeJSON(w, books)
	case http.MethodPost:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		book := FromJSON(body)
		isbn, created := CreateBook(book)
		if created {
			w.Header().Add("Location", "/api/books/"+isbn)
			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusConflict)
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unsupported request method."))
	}
}

// BookHandleFunc to be used as http.HandleFunc for Book API
func BookHandleFunc(w http.ResponseWriter, r *http.Request) {
	isbn := r.URL.Path[len("/api/books/"):]

	switch method := r.Method; method {
	case http.MethodGet:
		book, found := GetBook(isbn)
		if found {
			writeJSON(w, book)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	case http.MethodPut:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		book := FromJSON(body)
		exists := UpdateBook(isbn, book)
		if exists {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	case http.MethodDelete:
		DeleteBook(isbn)
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unsupported request method."))
	}
}

func writeJSON(w http.ResponseWriter, i interface{}) {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Write(b)
}

// GetBook returns the book for a given ISBN
func GetBook(isbn string) (*Book, bool) {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	pool, err := pgxpool.Connect(context.Background(), "user=postgres password=123 host=localhost port=5432 dbname=books sslmode=disable pool_max_conns=10")
	if err != nil {
		log.Fatalf("Unable to connection to database: %v\n", err)
	}
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		books:    &BookModel{DB: pool},
	}

	defer pool.Close()

	s, err := app.books.GetByIsbn(isbn)
	if err != nil {
		panic(err)

	}

	return s, false
}

// CreateBook creates a new Book if it does not exist
func CreateBook(book Book) (string, bool) {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	pool, err := pgxpool.Connect(context.Background(), "user=postgres password=123 host=localhost port=5432 dbname=books sslmode=disable pool_max_conns=10")
	if err != nil {
		log.Fatalf("Unable to connection to database: %v\n", err)
	}
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		books:    &BookModel{DB: pool},
	}

	defer pool.Close()

	s, err := app.books.Insert(book.Title, book.Author, book.Description, book.ISBN)
	if s != "" {
		return s, true
	}
	if err != nil {
		panic(err)
		return "", false

	}

	return s, false
}

// UpdateBook updates an existing book
func UpdateBook(isbn string, book Book) bool {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	pool, err := pgxpool.Connect(context.Background(), "user=postgres password=123 host=localhost port=5432 dbname=books sslmode=disable pool_max_conns=10")
	if err != nil {
		log.Fatalf("Unable to connection to database: %v\n", err)
	}
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		books:    &BookModel{DB: pool},
	}

	defer pool.Close()

	s, err := app.books.Update(book.Title, book.Author, book.Description, isbn)
	if err != nil {
		panic(err)

	}

	return s
}

// DeleteBook removes a book from the map by ISBN key
func DeleteBook(isbn string) {
	delete(books, isbn)
}
