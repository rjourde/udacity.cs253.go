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
	"models"
)

var secret []byte = securecookie.GenerateRandomKey(32)
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
			user := models.UserByUsername(r, username)
			
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
				
				u := models.User{ username, password, verify, email, time.Now() }
				key, err := datastore.Put(c, datastore.NewIncompleteKey(c, "User", nil), &u)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				userIdCookie = securecookie.New(secret, nil)
				
				stringID := fmt.Sprintf("%d", key.IntID())
				storeCookie(w, r, "user_id", stringID)
				
				// redirect to the page of the newly registered user
				http.Redirect(w, r, "/unit4/welcome", http.StatusFound)
				return
			}
		}
	}
}

func unit4Login(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		form := struct {
			Username string
			Password string
			ErrorLogin string
		}{
			"", "", "",
		}
		writeLoginForm(w, form)	
	}
	if r.Method == "POST" {
		// Get form field values
		username := r.FormValue("username")
		password := r.FormValue("password")
		
		// Validate form fields
		userID, user := models.UserByUsernameAndPassword(r, username, password)
		if(userID != 0 && len(user.Username) > 0) {
			if(username == user.Username && password == user.Password) {
				if(userIdCookie == nil){
					userIdCookie = securecookie.New(secret, nil)
				}
				stringID := fmt.Sprintf("%d", userID)
				storeCookie(w, r, "user_id", stringID)
				
				// redirect to the welcome page
				http.Redirect(w, r, "/unit4/welcome", http.StatusFound)
				return
			}
		}
		
		form := struct {
			Username string
			Password string
			ErrorLogin string
		}{
			username,
			password,
			"Invalid Login",
		}
		
		writeLoginForm(w, form)
	}
}

func unit4Logout(w http.ResponseWriter, r *http.Request) {
	clearCookie(w, "user_id")
	// redirect to the signup
	http.Redirect(w, r, "/unit4/signup", http.StatusFound)
	return
}

func unit4Welcome(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	// read the secure cookie
	if userId := fetchCookie(r, "user_id"); len(userId) > 0 {
		var user models.User
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
	} else {
		// redirect to the signup
		http.Redirect(w, r, "/unit4/signup", http.StatusFound)
		return
	}
}

// Cookie util methods

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

func fetchCookie(r *http.Request, cookieName string) string {
	if cookie, err := r.Cookie(cookieName); err == nil {
		value := make(map[string]string)
		if(userIdCookie != nil) {
			err = userIdCookie.Decode(cookieName, cookie.Value, &value)
			if (len(value[cookieName]) > 0 && err == nil) {
				return value[cookieName]
			}
		}
	}

	return ""
}

func clearCookie(w http.ResponseWriter, cookieName string) {
	cookie := &http.Cookie{
		Name:  cookieName,
		Value: "",
		Path:  "/",
	}
	http.SetCookie(w, cookie)
}



