package models

import (
	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/pborman/uuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Account struct {
	AccountID string `json:"-" bson:"account_id"`
	Username  string `json:"username" bson:"username"`
	Password  string `json:"-" bson:"-"`
	Hashword  string `json:"-" bson:"password"`
	SessionID string `json:"-" bson:"session_id"`
}

var (
	AccountNotFound = errors.New("account not found")
)

func CreateAccount(db *mgo.Database, a Account) error {
	accounts := db.C("accounts")
	cost := 3
	hashword, err := bcrypt.GenerateFromPassword([]byte(a.Password), cost)
	if err != nil {
		return err
	}
	a.Hashword = string(hashword)

	a.AccountID = uuid.New()
	err = accounts.Insert(a)
	if err != nil {
		return err
	}
	return nil
}

func LoadAccount(db *mgo.Database, username string) (Account, error) {
	accounts := db.C("accounts")
	a := Account{}
	query := bson.M{"username": username}
	err := accounts.Find(query).One(&a)
	if err == mgo.ErrNotFound {
		return a, AccountNotFound
	} else if err != nil {
		return a, err
	}
	return a, nil
}

func CheckSession(db *mgo.Database, sessionID string) (Account, error) {
	accounts := db.C("accounts")
	a := Account{}
	query := bson.M{"session_id": sessionID}
	err := accounts.Find(query).One(&a)
	if err != nil {
		logrus.Error(err)
		return a, err
	}
	return a, nil
}

func (a Account) NewSession(db *mgo.Database) (string, error) {
	accounts := db.C("accounts")
	sessionID := uuid.New()
	query := bson.M{"account_id": a.AccountID}
	update := bson.M{"$set": bson.M{"session_id": sessionID}}
	err := accounts.Update(query, update)
	if err != nil {
		return "", err
	}
	return sessionID, nil
}

func (a Account) ClearSession(db *mgo.Database) error {
	accounts := db.C("accounts")
	query := bson.M{"account_id": a.AccountID}
	update := bson.M{"$unset": bson.M{"session_id": ""}}
	err := accounts.Update(query, update)
	if err != nil {
		return err
	}
	return nil
}
