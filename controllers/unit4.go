package controllers

import (
	"html/template"
	"net/http"
	"appengine"
    "appengine/datastore"
	"time"
	"github.com/gorilla/securecookie"
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
			
			if(user != nil) {
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

				userIdCookie = securecookie.New(makeSalt(), makeSalt())
				
				storeCookie(w, r, "user_id", key.StringID())
				
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				// redirect to the page of the newly created post
				http.Redirect(w, r, "/unit4/welcome", http.StatusFound)
				return
			}
		}
	}
}

func userByUsername(r *http.Request,username string) *User {
	var user *User
	
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
    if encoded, err := userIdCookie.Encode(cookieValue, value); err == nil {
        cookie := &http.Cookie{
            Name:  cookieName,
            Value: encoded,
            Path:  "/",
        }
        http.SetCookie(w, cookie)
    }
}

func fetchCookie(c appengine.Context, r *http.Request, cookieName string) string{
	// read the secure cookie
	cookie, err := r.Cookie(cookieName)
	c.Debugf("cookieName : %s", cookieName)
	c.Debugf("cookie : %v", cookie)
	c.Debugf("cookie.Value : %v", cookie.Value)
	
	if err == nil {
		value := make(map[string]string)
		err = userIdCookie.Decode(cookieName, cookie.Value, &value)
        if err == nil {
			c.Debugf("cookie.Value : %s", cookie.Value)
			c.Debugf("value[cookieName] : %s", value[cookieName])
            return value[cookieName]
        } else {
			c.Debugf("ERROR : %s", err.Error())
		}
    }
	
	return ""
}

func unit4Welcome(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	// read the secure cookie
	if userId := fetchCookie(c, r, "user_id"); len(userId) > 0 {
		c.Debugf("UserID: %v", userId)
		var user *User
		
		key := datastore.NewKey(c, "User", userId, 0, nil)
		if err := datastore.Get(c, key, &user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		var username string = "Nobody"
		
		if(user != nil) {
			username = user.Username
		}
		
		t, _ := template.ParseFiles("templates/welcome.html")
		err := t.Execute(w, username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		c.Debugf("Nobody")
	}
}

func makeSalt() []byte {
	return securecookie.GenerateRandomKey(32)
}


