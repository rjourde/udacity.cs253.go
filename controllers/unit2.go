package controllers

import (
	"html/template"
	"net/http"
)

func unit2Rot13(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/rot13.html")
	err := t.Execute(w, encodeROT13(r.FormValue("text")))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func encodeROT13(text string) string {
	b := []byte(text)
	
	for i := 0; i < len(text); i++ {
		b[i] = rot13(b[i])
	}
	
	return string(b)
}

func rot13(b byte) byte {
	var a, z byte
	switch {
	case 'a' <= b && b <= 'z':
		a, z = 'a', 'z'
	case 'A' <= b && b <= 'Z':
		a, z = 'A', 'Z'
	default:
		return b
	}
	return (b-a+13)%(z-a+1) + a
}

func unit2Signup(w http.ResponseWriter, r *http.Request) {
	
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
		}else {
			http.Redirect(w, r, "/unit2/welcome?username="+username, http.StatusFound)
			return
		}
	}
}

func unit2Welcome(w http.ResponseWriter, r *http.Request) {
	// get 'username' parameter
	parameter := r.FormValue("username")
	
	t, _ := template.ParseFiles("templates/welcome.html")
	err := t.Execute(w, parameter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
