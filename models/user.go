package models

import (
	"net/http"
	"appengine"
	"appengine/datastore"
	"time"
)

type User struct {
	Id int64
	Username string
	Password string
	Verify string 
	Email string
	Created time.Time
}

// User util methods

func UserByUsername(r *http.Request, username string) User {
	var user User

	c := appengine.NewContext(r)
	q:= datastore.NewQuery("User").Filter("Username =", username)
	for t := q.Run(c); ; {
		_, err := t.Next(&user)
		if err == datastore.Done {
				break
		}
	}
	return user
}

func UserByUsernameAndPassword(r *http.Request, username, password string) *User {
	var users []*User

	c := appengine.NewContext(r)
	// Fetch the user
	q:= datastore.NewQuery("User").Filter("Username =", username).Filter("Password =", password)
	_, err := q.GetAll(c, &users)
	if err != nil || users == nil {
		return nil
	}
	
	return users[0]
}