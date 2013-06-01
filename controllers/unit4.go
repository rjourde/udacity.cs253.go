package controllers

import (
	"html/template"
	"net/http"
	"appengine"
    "appengine/datastore"
	"time"
	"github.com/gorilla/securecookie"
	"fmt"
	"strconv"
)

type User struct {
	Username string
	Password string
	Verify string 
	Email string
	Created time.Time
}

var userIdCookie *securecookie.SecureCookie

func unit4Signup(w http.ResponseWriter, r *http.Request) {
	
	if r.Method == "GET" {
		form := struct {
			Username string
			Password string
			Verify string
			Email string
			ErrorUsername string
			ErrorPassword string
			ErrorVerify string
			ErrorEmail string
		}{
			"", "", "", "", "", "", "", "",
		}
		writeForm(w, form)	
	}
	if r.Method == "POST" {
		errorUsername := ""
		errorPassword := ""
		errorVerify := ""
		errorEmail := ""
		// Get form field values
		username := r.FormValue("username")
		password := r.FormValue("password")
		verify := r.FormValue("verify")
		email := r.FormValue("email")
		// Validate form fields
		if ! (validUsername(username) && validPassword(password) && (password == verify) && validEmail(email)) {
			if !validUsername(username) {
				errorUsername = "That's not a valid username"
			}
			if !validPassword(password) {
				errorPassword = "That's not a valid password"
			}
			if(password != verify) {
				errorVerify = "Your passwords didn't match"
			}
			if !validEmail(email) {
				errorEmail = "That's not a valid email"
			}
			
			password = ""
			verify = ""
			
			form := struct {
				Username string
				Password string
				Verify string
				Email string
				ErrorUsername string
				ErrorPassword string
				ErrorVerify string
				ErrorEmail string
			}{
				username,
				password,
				verify,
				email,
				errorUsername,
				errorPassword,
				errorVerify,
				errorEmail,
			}
			
			writeForm(w, form)	
		} else {
			user := userByUsername(r, username)
			
			if(len(user.Username) > 0) {
				errorUsername = "That user already exists"
				
				form := struct {
					Username string
					Password string
					Verify string
					Email string
					ErrorUsername string
					ErrorPassword string
					ErrorVerify string
					ErrorEmail string
				}{
					username,
					password,
					verify,
					email,
					errorUsername,
					errorPassword,
					errorVerify,
					errorEmail,
				}
				
				writeForm(w, form)	
			} else {
				c := appengine.NewContext(r)
				
				u := User{ username, password, verify, email, time.Now() }
				key, err := datastore.Put(c, datastore.NewIncompleteKey(c, "User", nil), &u)

				userIdCookie = securecookie.New(securecookie.GenerateRandomKey(32), nil)
				
				stringID := fmt.Sprintf("%d", key.IntID())
				storeCookie(w, r, "user_id", stringID)
				
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				// redirect to the page of the newly registered user
				http.Redirect(w, r, "/unit4/welcome", http.StatusFound)
				return
			}
		}
	}
}

func userByUsername(r *http.Request,username string) User {
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

func storeCookie(w http.ResponseWriter, r *http.Request, cookieName, cookieValue string) {
	value := map[string]string{
        cookieName : cookieValue,
    }
    if encoded, err := userIdCookie.Encode(cookieName, value); err == nil {
        cookie := &http.Cookie{
            Name:  cookieName,
            Value: encoded,
            Path:  "/",
        }
        http.SetCookie(w, cookie)
    }
}

func fetchCookie(c appengine.Context, r *http.Request, cookieName string) string{
	if cookie, err := r.Cookie(cookieName); err == nil {
		value := make(map[string]string)
		err = userIdCookie.Decode(cookieName, cookie.Value, &value)
        if err == nil {
            return value[cookieName]
        }
    }
	
	return ""
}

func unit4Welcome(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	// read the secure cookie
	if userId := fetchCookie(c, r, "user_id"); len(userId) > 0 {
		var user User
		intID, _ := strconv.ParseInt(userId, 10, 64)
		key := datastore.NewKey(c, "User", "", intID, nil)
		
		if err := datastore.Get(c, key, &user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		t, _ := template.ParseFiles("templates/welcome.html")
		err := t.Execute(w, user.Username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}



