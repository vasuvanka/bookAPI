package model

import (
	"time"

	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	COLLECTION = "books"
	DB         = "test"
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

func All(session *mgo.Session) (error, []Book) {
	var books []Book
	c := session.DB(DB).C(COLLECTION)
	if err := c.Find(bson.M{}).All(&books); err != nil {
		return err, nil
	}
	return nil, books
}

func ById(session *mgo.Session, id string) (error, Book) {
	var book Book
	c := session.DB(DB).C(COLLECTION)
	if !isValidObjectId(id) {
		return NewError("invalid objectID"), Book{}
	}
	if err := c.FindId(bson.ObjectIdHex(id)).One(&book); err != nil {
		return err, Book{}
	}
	return nil, book
}

func Save(session *mgo.Session, book Book) (error, Book) {
	book.ID = bson.NewObjectId()
	book.When = time.Now()
	c := session.DB(DB).C(COLLECTION)
	if err := c.Insert(book); err != nil {
		return err, Book{}
	}
	return nil, book
}

func Delete(session *mgo.Session, id string) (error, Book) {
	var book Book
	var err error
	c := session.DB(DB).C(COLLECTION)
	if !isValidObjectId(id) {
		return NewError("invalid objectID"), Book{}
	}

	if err, book = ById(session, id); err != nil {
		return err, Book{}
	}

	if err := c.RemoveId(bson.ObjectIdHex(id)); err != nil {
		return err, Book{}
	}

	return nil, book

}

func Update(session *mgo.Session, book Book) (error, Book) {
	c := session.DB(DB).C(COLLECTION)
	if !isValidObjectId(book.ID.Hex()) {
		return NewError("invalid objectID"), Book{}
	}
	book.When = time.Now()
	if err := c.UpdateId(book.ID, bson.M{"$set": book}); err != nil {
		return err, Book{}
	}
	return ById(session, book.ID.Hex())
}

func isValidObjectId(id string) bool {
	return bson.IsObjectIdHex(id)
}
