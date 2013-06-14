package controllers

import (
	"html/template"
	"net/http"
	"appengine"
	"appengine/datastore"
	"time"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"models"
	"tools"
)

var wikiSecret []byte = securecookie.GenerateRandomKey(32)
var wikiUserIdCookie *securecookie.SecureCookie

var currentUser *models.User

type NavItem struct {
	URL string
	Name string
}

func wikiFrontPage(w http.ResponseWriter, r *http.Request) {
	renderWikiFrontPage(w)
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
				
				userID, _, _ := datastore.AllocateIDs(c, "User", nil, 1)
				key := datastore.NewKey(c, "User", "", userID, nil)
				u := models.User{ userID, username, password, verify, email, time.Now() }
				_, err := datastore.Put(c, key, &u)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
	
				wikiUserIdCookie = securecookie.New(wikiSecret, nil)
				
				stringID := fmt.Sprintf("%d", u.Id)
				tools.StoreCookie(w, r, wikiUserIdCookie, "user_id", stringID)
				// set the current user
				currentUser = &u
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
		user := models.UserByUsernameAndPassword(r, username, password)
		if(len(user.Username) > 0) {
			if(username == user.Username && password == user.Password) {
				if(wikiUserIdCookie == nil){
					wikiUserIdCookie = securecookie.New(wikiSecret, nil)
				}
				stringID := fmt.Sprintf("%d", user.Id)
				tools.StoreCookie(w, r, wikiUserIdCookie, "user_id", stringID)
				
				// set the current user
				currentUser = user
				
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
	tools.ClearCookie(w, "user_id")
	// clear the current user
	currentUser = nil
	// redirect to the wiki front page
	http.Redirect(w, r, "/wiki", http.StatusFound)
	return
}

func wikiEdit(w http.ResponseWriter, r *http.Request) {
	// get the page name in the URL
	vars := mux.Vars(r)
	pageName := vars["page"]
	
	if r.Method == "GET" {
		// fetch the page only if you are already looged in
		if(currentUser != nil) {
			// if the page does not exist redirect to new page form
			if page, err := models.GetPage(r, pageName); err != nil {
				// redirect to the wiki page
				http.Redirect(w, r, "/wiki/" + pageName, http.StatusFound)
				return
			} else {
				renderNewPageForm(w, page.Content)
			}
		} else {
			// redirect to the login page
			http.Redirect(w, r, "/wiki/login", http.StatusFound)
			return
		}
	}
	if r.Method == "POST" {
		content := r.FormValue("content")
		// if the page does not exist redirect to the new page form
		if page, err := models.GetPage(r, pageName); err != nil {
			renderNewPageForm(w, nil)
		} else {
			// update page
			models.UpdatePage(r, *page, pageName, content)
			
			// redirect to the wiki page
			http.Redirect(w, r, "/wiki/" + pageName, http.StatusFound)
			return
		}
	}
}

func wikiPage(w http.ResponseWriter, r *http.Request) {
	// get the page name in the URL
	vars := mux.Vars(r)
	pageName := vars["page"]
	
	if r.Method == "GET" {
		// fetch the page
		// if the page does not exist redirect to the new page form
		if page, err := models.GetPage(r, pageName); err != nil {
			renderNewPageForm(w, nil)
		} else {
			renderPageView(w, *page)
		}
	}
	if r.Method == "POST" {
		content := r.FormValue("content")
		
		err := models.AddPage(r, pageName, content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		// redirect to the wiki front page
		http.Redirect(w, r, "/wiki", http.StatusFound)
		return
	}
}

func authenticationItems() []NavItem {
	if currentUser != nil {
		return []NavItem{ {URL: "/wiki/logout", Name: "logout(" + currentUser.Username + ")" } }
	} 
	
	return []NavItem{ {URL: "/wiki/signup", Name: "signup" },
					  {URL: "/wiki/login", Name: "login"} }
}

func navigationItems(pageURL string) []NavItem {
	if currentUser != nil {
		return []NavItem{ {URL: "/wiki/logout", Name: "logout(" + currentUser.Username + ")" },
						  {URL: "/wiki/_edit/" + pageURL, Name: "edit"} }
	} 

	return []NavItem{ {URL: "/wiki/login", Name: "login"},
					  {URL: "/wiki/signup", Name: "signup" } }
}

func renderWikiFrontPage(w http.ResponseWriter) {
	t, _ := template.ParseFiles("templates/navigation.html", "templates/wiki.html")

	
	if err := t.ExecuteTemplate(w, "tmpl_wiki", authenticationItems()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func renderNewPageForm(w http.ResponseWriter, data interface{}) {
	t, _ := template.ParseFiles("templates/newpage.html")
	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func renderPageView(w http.ResponseWriter, page models.Page) {
	t, _ := template.ParseFiles("templates/navigation.html", "templates/page.html")
	
	article := struct {
		Navigation []NavItem
		Page models.Page
	}{
		navigationItems(page.Name),
		page,
	}
	
	if err := t.ExecuteTemplate(w, "tmpl_page", article); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}