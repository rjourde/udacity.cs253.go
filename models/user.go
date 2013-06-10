package models

import (
	"net/http"
	"appengine"
	"appengine/datastore"
	"time"
)

type User struct {
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


func UserByUsernameAndPassword(r *http.Request, username, password string) (int64, *User) {
	var users []*User

	c := appengine.NewContext(r)
	// Fetch the user
	q:= datastore.NewQuery("User").Filter("Username =", username).Filter("Password =", password)
	keys, err := q.GetAll(c, &users)
	if err != nil {
		return 0, nil
	}
	if(keys == nil || users == nil) {
		return 0, nil
	}
	return keys[0].IntID(), users[0]
}