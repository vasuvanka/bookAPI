package route

import (
	"encoding/json"
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"

	"gopkg.in/mgo.v2"

	"../model"
)

const MONGO_URL = "localhost"

var sess *mgo.Session

type Message struct {
	Message string `json:"message"`
	Code    int    `json:"status"`
}

func Initialize() {
	session, err := mgo.Dial(MONGO_URL)
	if err != nil {
		log.Fatal("cannot dial mongo", err)
	}
	sess = session.Copy()
	defer session.Close()
}

func Welcome(w http.ResponseWriter, r *http.Request) {
	var m = &Message{"welcome to book api!", 200}
	WithError(w, *m)
}

func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	err, books := model.All(sess)
	if err != nil {
		WithError(w, Message{"invalid request", http.StatusBadRequest})
		return
	}
	Json(w, http.StatusOK, books)
}
func GetBookById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	err, book := model.ById(sess, params["id"])
	if err != nil {
		WithError(w, Message{"invalid request", http.StatusBadRequest})
		return
	}
	Json(w, http.StatusOK, book)
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var book model.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		WithError(w, Message{"invalid request", http.StatusBadRequest})
		return
	}
	err, b := model.Save(sess, book)
	if err != nil {
		WithError(w, Message{"invalid request", http.StatusBadRequest})
		return
	}
	Json(w, http.StatusOK, b)
}
func UpdateBookById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	defer r.Body.Close()
	var book model.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		WithError(w, Message{"invalid request", http.StatusBadRequest})
		return
	}
	book.ID = bson.ObjectIdHex(params["id"])
	err, b := model.Update(sess, book)
	if err != nil {
		WithError(w, Message{"invalid request", http.StatusBadRequest})
		return
	}
	Json(w, http.StatusOK, b)
}
func RemoveBookById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	err, book := model.Delete(sess, params["id"])
	if err != nil {
		WithError(w, Message{"invalid request", http.StatusBadRequest})
		return
	}
	Json(w, http.StatusOK, book)
}
func Json(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		WithError(w, Message{"invalid request", http.StatusBadRequest})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func WithError(w http.ResponseWriter, m Message) {
	response, _ := json.Marshal(m)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(m.Code)
	w.Write(response)
}
