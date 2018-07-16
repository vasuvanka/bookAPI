package main

import (
	"log"
	"net/http"

	"./route"

	"github.com/gorilla/mux"
)

func main() {
	route.Initialize()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", route.Welcome).Methods("GET")
	bookRouter := router.PathPrefix("/api/book").Subrouter()
	bookRouter.HandleFunc("/", route.GetAllBooks).Methods("GET")
	bookRouter.HandleFunc("/", route.CreateBook).Methods("POST")
	bookRouter.HandleFunc("/{id}", route.GetBookById).Methods("GET")
	bookRouter.HandleFunc("/{id}", route.UpdateBookById).Methods("PUT")
	bookRouter.HandleFunc("/{id}", route.RemoveBookById).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":3000", router))
}
