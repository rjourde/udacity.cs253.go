package controllers

import (
	"net/http"
	"appengine"
	"appengine/datastore"
	"time"
	"fmt"
	"github.com/gorilla/securecookie"
	"models"
)

var wikiSecret []byte = securecookie.GenerateRandomKey(32)
var wikiUserIdCookie *securecookie.SecureCookie

func wikiFrontPage(w http.ResponseWriter, r *http.Request) {

}

func wikiSignup(w http.ResponseWriter, r *http.Request) {
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
				key, err := datastore.Put(c, datastore.NewIncompleteKey(c, "WikiUser", nil), &u)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
	
				wikiUserIdCookie = securecookie.New(wikiSecret, nil)
				
				stringID := fmt.Sprintf("%d", key.IntID())
				storeCookie(w, r, "user_id", stringID)
				
				// redirect to the wiki front page
				http.Redirect(w, r, "/wiki", http.StatusFound)
				return
			}
		}
	}
}

func wikiLogin(w http.ResponseWriter, r *http.Request) {
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
				if(wikiUserIdCookie == nil){
					wikiUserIdCookie = securecookie.New(wikiSecret, nil)
				}
				stringID := fmt.Sprintf("%d", userID)
				storeCookie(w, r, "user_id", stringID)
				
				// redirect to the wiki front page
				http.Redirect(w, r, "/wiki", http.StatusFound)
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

func wikiLogout(w http.ResponseWriter, r *http.Request) {
	clearCookie(w, "user_id")
	// redirect to the wiki front page
	http.Redirect(w, r, "/wiki", http.StatusFound)
	return
}

func wikiHistory(w http.ResponseWriter, r *http.Request) {

}

func wikiEdit(w http.ResponseWriter, r *http.Request) {

}

func wikiPage(w http.ResponseWriter, r *http.Request) {

}