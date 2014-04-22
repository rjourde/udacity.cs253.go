package controllers

import (
	"html/template"
	"net/http"
	"appengine"
  "appengine/datastore"
  "appengine/user"
	"time"
	"github.com/gorilla/securecookie"
	"fmt"
	"strconv"
	"models"
	"tools"
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
				
				userID, _, _ := datastore.AllocateIDs(c, "User", nil, 1)
				key := datastore.NewKey(c, "WikiUser", "", userID, nil)
				u := models.User{ userID, username, password, verify, email, time.Now() }
				_, err := datastore.Put(c, key, &u)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				userIdCookie = securecookie.New(secret, nil)
				
				stringID := fmt.Sprintf("%d", key.IntID())
				tools.StoreCookie(w, r, userIdCookie, "user_id", stringID)
				
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
		user := models.UserByUsernameAndPassword(r, username, password)
		if(len(user.Username) > 0) {
			if(username == user.Username && password == user.Password) {
				if(userIdCookie == nil){
					userIdCookie = securecookie.New(secret, nil)
				}
				stringID := fmt.Sprintf("%d", user.Id)
				tools.StoreCookie(w, r, userIdCookie, "user_id", stringID)
				
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

func unit4GoogleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
    c := appengine.NewContext(r)
    u := user.Current(c)
    if u == nil {
        url, _ := user.LoginURL(c, "/")
        fmt.Fprintf(w, `<a href="%s">Sign in or register</a>`, url)
        return
    }
    url, _ := user.LogoutURL(c, "/")
    fmt.Fprintf(w, `Welcome, %s! (<a href="%s">sign out</a>)`, u, url)
}

func unit4Logout(w http.ResponseWriter, r *http.Request) {
	tools.ClearCookie(w, "user_id")
	// redirect to the signup
	http.Redirect(w, r, "/unit4/signup", http.StatusFound)
	return
}

func unit4Welcome(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	// read the secure cookie
	if userId := tools.FetchCookie(r, userIdCookie, "user_id"); len(userId) > 0 {
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



