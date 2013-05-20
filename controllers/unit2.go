package controllers

import (
	"html/template"
	"net/http"
)

type User struct {
	Username string
	Password string
	Verify string 
	Email string
}

type Error struct {
	ErrorUsername string
	ErrorPassword string
	ErrorVerify string
	ErrorEmail string
}

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
	// Get form field values
	username := r.FormValue("username")
	password := r.FormValue("password")
	verify := r.FormValue("verify")
	email := r.FormValue("email")
	
	u := new(User)
	
	// Validate form fields
	if ! (validUsername(username) && validPassword(password) && validEmail(email)) {
		
		t, _ := template.ParseFiles("templates/signup.html")
		err := t.Execute(w, u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}	
	}else {
		t, _ := template.ParseFiles("templates/welcome.html")
		err := t.Execute(w, u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func validUsername(string) bool {
	return true
}

func validPassword(string) bool {
	return true
}

func validEmail(string) bool {
	return true
}