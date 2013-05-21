package controllers

import (
	"html/template"
	"net/http"
	"regexp"
)

type SignupForm struct {
	Username string
	Password string
	Verify string 
	Email string
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
	form := new(SignupForm)
	
	if r.Method == "GET" {
		writeForm(w, form)	
	}
	if r.Method == "POST" {
		// Get form field values
		form.Username = r.FormValue("username")
		form.Password = r.FormValue("password")
		form.Verify = r.FormValue("verify")
		form.Email = r.FormValue("email")
		
		// Validate form fields
		if ! (validUsername(form.Username) && validPassword(form.Password) && (form.Password == form.Verify) && validEmail(form.Email)) {
			if !validUsername(form.Username) {
				form.ErrorUsername = "That's not a valid username"
			}
			if !validPassword(form.Password) {
				form.ErrorPassword = "That's not a valid password"
			}
			if(form.Password != form.Verify) {
				form.ErrorVerify = "Your passwords didn't match"
			}
			if !validEmail(form.Email) {
				form.ErrorEmail = "That's not a valid email"
			}
			
			form.Password = ""
			form.Verify = ""
			
			writeForm(w, form)	
		}else {
			http.Redirect(w, r, "/unit2/welcome?username="+form.Username, http.StatusFound)
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

func writeForm(w http.ResponseWriter, form *SignupForm) {
	t, _ := template.ParseFiles("templates/signup.html")
	err := t.Execute(w, form)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func validUsername(username string) bool {
	valid := regexp.MustCompile(`^[a-zA-Z0-9_-]{3,20}$`)

	return valid.MatchString(username)
}

func validPassword(password string) bool {
	valid := regexp.MustCompile(`^.{3,20}$`)

	return valid.MatchString(password)
}

func validEmail(email string) bool {
	valid := regexp.MustCompile(`^[\S]+@[\S]+\.[\S]+$`)

	return valid.MatchString(email)
}
