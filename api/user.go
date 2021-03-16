package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type User struct {
	ID       int
	Username string
	Password string
	Email    string
}

func StartApi(string) {
	router := mux.NewRouter()
	router.HandleFunc("/api/login", login).Methods("GET")
	fmt.Println("App is working on port :8081")
	log.Fatal(http.ListenAndServe(":8081", router))

}
func UserHandleFunc(w http.ResponseWriter, r *http.Request) {
	login := r.URL.Path[len("/api/login"):]

	switch method := r.Method; method {
	case http.MethodGet:
		StartApi(login)

		users := [2]User{
			{Username: "Sabina", Email: "sab@sa.com"},
			{Username: "Talga", Email: "talg@tal.com"},
		}
		writeJSON(w, users)
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unsupported request method."))
	}
}
