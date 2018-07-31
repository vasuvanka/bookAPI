package shared

import (
	"context"
	"log"
	"net/http"

	mgo "gopkg.in/mgo.v2"
)

const MONGO_URL = "localhost"
const SESSION_KEY = "session"

func GetSession() *mgo.Session {
	session, err := mgo.Dial(MONGO_URL)
	if err != nil {
		log.Fatal("cannot dial mongo", err)
	}
	sess := session.Copy()
	defer session.Close()
	return sess
}

func AttachSession(w http.ResponseWriter, r *http.Request, next http.Handler) {
	session := GetSession()
	ctx := context.WithValue(r.Context(), SESSION_KEY, session)
	go func() {
		select {
		case <-ctx.Done():
			session.Close()
		}
	}()
}
