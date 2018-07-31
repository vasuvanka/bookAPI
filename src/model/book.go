package model

import (
	"log"
	"time"

	"../shared"
	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	COLLECTION_BOOK = "books"
	DB              = "test"
)

type Config struct {
	router  *mux.Router
	session *mgo.Session
}

type Book struct {
	ID          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Author      string        `json:"author" bson:"author,omitempty"`
	Name        string        `json:"name" bson:"name,omitempty"`
	When        time.Time     `json:"when" bson:"when,omitempty"`
	Publication string        `json:"publication" bson:"pub,omitempty"`
}

type errorString struct {
	s string
}

func (e errorString) Error() string {
	return e.s
}

func NewError(text string) error {
	return errorString{text}
}

func All() (error, []Book) {
	session := shared.GetSession()
	if session == nil {
		log.Fatal("session null")
	}
	var books []Book
	c := session.DB(DB).C(COLLECTION_BOOK)
	if err := c.Find(bson.M{}).All(&books); err != nil {
		session.Close()
		return err, nil
	}
	session.Close()
	return nil, books
}

func ById(id string) (error, Book) {
	var book Book
	session := shared.GetSession()
	c := session.DB(DB).C(COLLECTION_BOOK)
	if !isValidObjectId(id) {
		session.Close()
		return NewError("invalid objectID"), Book{}
	}
	if err := c.FindId(bson.ObjectIdHex(id)).One(&book); err != nil {
		session.Close()
		return err, Book{}
	}
	session.Close()
	return nil, book
}

func Save(book Book) (error, Book) {
	book.ID = bson.NewObjectId()
	book.When = time.Now()
	session := shared.GetSession()
	c := session.DB(DB).C(COLLECTION_BOOK)
	if err := c.Insert(book); err != nil {
		session.Close()
		return err, Book{}
	}
	session.Close()
	return nil, book
}

func Delete(id string) (error, Book) {
	var book Book
	var err error
	session := shared.GetSession()
	c := session.DB(DB).C(COLLECTION_BOOK)
	if !isValidObjectId(id) {
		session.Close()
		return NewError("invalid objectID"), Book{}
	}

	if err, book = ById(id); err != nil {
		session.Close()
		return err, Book{}
	}

	if err := c.RemoveId(bson.ObjectIdHex(id)); err != nil {
		session.Close()
		return err, Book{}
	}
	session.Close()
	return nil, book

}

func Update(book Book) (error, Book) {
	session := shared.GetSession()
	c := session.DB(DB).C(COLLECTION_BOOK)
	if !isValidObjectId(book.ID.Hex()) {
		session.Close()
		return NewError("invalid objectID"), Book{}
	}
	book.When = time.Now()
	if err := c.UpdateId(book.ID, bson.M{"$set": book}); err != nil {
		session.Close()
		return err, Book{}
	}
	return ById(book.ID.Hex())
}

func isValidObjectId(id string) bool {
	return bson.IsObjectIdHex(id)
}
