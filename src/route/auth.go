package route

import (
	"encoding/json"
	"log"
	"net/http"

	"../model"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var login model.LoginInfo
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		log.Fatal("error ", err)
		WithError(w, Message{err.Error(), 400})
		return
	}

	if err, userLogin := model.Login(login.Username, login.Password); err == nil {
		Json(w, http.StatusOK, userLogin)
	} else {
		WithError(w, Message{err.Error(), 400})
		return
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	var loginInfo model.LoginInfo
	if err := json.NewDecoder(r.Body).Decode(&loginInfo); err != nil {
		log.Fatal(err)
		WithError(w, Message{err.Error(), 400})
		return
	}

	if err, newUser := model.Register(loginInfo); err == nil {
		Json(w, http.StatusOK, newUser)
	} else {
		log.Fatal(err)
		WithError(w, Message{"invalid user data", 400})
	}
}
