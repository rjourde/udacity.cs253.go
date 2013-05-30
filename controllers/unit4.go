package controllers

import (
	"html/template"
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
				//store the hashed user id in a cookie
				storeCookie(r, "user_id", key.StringID())
				
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

func storeCookie(r *http.Request, cookieName, cookieValue string) {
	// Hash the value
	
	// Set the value to a cookie
	r.Header().Add("Set-Cookie", cookieName + "=" + cookieValue + "Path=/") 
}

func unit4Welcome(w http.ResponseWriter, r *http.Request) {
	// get user id from the cookie
	// get the user from the id
	// get the username
	username := "Myself"
	
	t, _ := template.ParseFiles("templates/welcome.html")
	err := t.Execute(w, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}


