package model

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"

	"../shared"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

const (
	COLLECTION_USER = "gousers"
	COLLECTION_AUTH = "goauthorization"
	INVALID_USER    = "invalid user info"
	JWT_SECRET      = "m;}YW-JCq5:h^.uu"
)

type User struct {
	ID       bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string        `json:"name,omitempty" bson:"name,omitempty"`
	When     time.Time     `json:"when,omitempty" bson:"when,omitempty"`
	Username string        `json:"username,omitempty" bson:"username,omitempty"`
	Password string        `json:"password,omitempty" bson:"password,omitempty"`
	LastSeen time.Time     `json:"last,omitempty" bson:"last,omitempty"`
	Token    string        `json:"token"`
}

type Authorization struct {
	ID     bson.ObjectId `json:"id" bson:"_id"`
	token  string        `json:"token" bson:"token"`
	UserId bson.ObjectId `json:"userId" bson:"userId"`
}

type LoginInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func Authorize(jwtToken string) (error, Authorization) {
	var user User
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, NewError("invalid token")
		}
		return []byte(JWT_SECRET), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		mapstructure.Decode(claims, &user)
	} else {
		return NewError("not a valid token"), Authorization{}
	}
	if err != nil {
		return err, Authorization{}
	}

	session := shared.GetSession()
	var authObj Authorization
	c := session.DB(DB).C(COLLECTION_AUTH)
	if err := c.Find(bson.M{"token": jwtToken, "userId": user.ID}).One(&authObj); err != nil {
		session.Close()
		return err, Authorization{}
	}
	session.Close()
	return nil, authObj
}

func HashPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	return string(hash), err
}

func CheckPasswordHash(pwd, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	return err
}

func Login(username, pwd string) (error, User) {
	var user User
	session := shared.GetSession()
	c := session.DB(DB).C(COLLECTION_USER)
	if err := c.Find(bson.M{"username": username}).One(&user); err != nil {
		session.Close()
		return err, User{}
	}

	if err := CheckPasswordHash(pwd, user.Password); err != nil {
		session.Close()
		return err, User{}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"id":       user.ID,
		"name":     user.Name,
		"last":     user.When,
	})
	tokenStr, err := token.SignedString([]byte(JWT_SECRET))
	if err != nil {
		session.Close()
		return err, User{}
	}
	var authorization Authorization
	authorization.ID = bson.NewObjectId()
	authorization.token = tokenStr
	authorization.UserId = user.ID
	c = session.DB(DB).C(COLLECTION_AUTH)
	if err := c.Insert(authorization); err != nil {
		session.Close()
		return err, User{}
	}
	user.Token = tokenStr
	session.Close()
	return nil, toJson(user)
}

func toJson(user User) User {
	user.Password = ""
	return user
}

func Register(loginInfo LoginInfo) (error, User) {
	var user User
	session := shared.GetSession()
	user.ID = bson.NewObjectId()
	user.When = time.Now()
	password, err := HashPassword(loginInfo.Password)
	if err != nil {
		return err, User{}
	}
	user.Password = password
	user.Name = loginInfo.Name
	user.Username = loginInfo.Username
	c := session.DB(DB).C(COLLECTION_USER)
	if err := c.Insert(user); err != nil {
		session.Close()
		return err, User{}
	}
	session.Close()
	return nil, user
}
