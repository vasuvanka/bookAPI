package main

import (
	"log"
	"net/http"

	"./src/route"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", route.Welcome).Methods("GET")
	authRouter := router.PathPrefix("/auth").Subrouter().StrictSlash(true)
	authRouter.HandleFunc("/login", route.Login).Methods("POST")
	authRouter.HandleFunc("/register", route.Register).Methods("POST")
	apiRouter := router.PathPrefix("/api").Subrouter().StrictSlash(true)

	bookRouter := apiRouter.PathPrefix("/book").Subrouter()

	bookRouter.HandleFunc("/", route.ValidateMiddleware(route.GetAllBooks)).Methods("GET")
	bookRouter.HandleFunc("/", route.ValidateMiddleware(route.CreateBook)).Methods("POST")
	bookRouter.HandleFunc("/{id}", route.ValidateMiddleware(route.GetBookById)).Methods("GET")
	bookRouter.HandleFunc("/{id}", route.ValidateMiddleware(route.UpdateBookById)).Methods("PUT")
	bookRouter.HandleFunc("/{id}", route.ValidateMiddleware(route.RemoveBookById)).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":3000", router))
}
