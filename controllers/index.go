package controllers

import (
    "html/template"
    "net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
    t, _ := template.ParseFiles("templates/index.html")
    t.Execute(w, nil)
}