package controllers

import (
	"regexp"
	"html/template"
	"net/http"
)

func writeForm(w http.ResponseWriter, data interface{}) {
	t, _ := template.ParseFiles("templates/signup.html")
	err := t.Execute(w, data)
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