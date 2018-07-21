package route

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"gopkg.in/mgo.v2/bson"

	"../model"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

const AUTH_KEY = "auth_key"
const JWT_SECRET = "m;}YW-JCq5:h^.uu"

type Message struct {
	Message string `json:"message"`
	Code    int    `json:"status"`
}

func Welcome(w http.ResponseWriter, r *http.Request) {
	var m = &Message{"welcome to book api!", 200}
	WithError(w, *m)
}

func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	err, books := model.All()
	if err != nil {
		WithError(w, Message{"invalid request", http.StatusBadRequest})
		return
	}
	Json(w, http.StatusOK, books)
}
func GetBookById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	err, book := model.ById(params["id"])
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
	err, b := model.Save(book)
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
	err, b := model.Update(book)
	if err != nil {
		WithError(w, Message{"invalid request", http.StatusBadRequest})
		return
	}
	Json(w, http.StatusOK, b)
}
func RemoveBookById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	err, book := model.Delete(params["id"])
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

func ValidateMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					return []byte(model.JWT_SECRET), nil
				})
				if error != nil {
					WithError(w, Message{error.Error(), 401})
					return
				}
				if token.Valid {
					context.Set(req, AUTH_KEY, token.Claims)
					next(w, req)
				} else {
					WithError(w, Message{"invalid token", 401})
				}
			}
		} else {
			WithError(w, Message{"invalid token", 401})
		}
	})
}
