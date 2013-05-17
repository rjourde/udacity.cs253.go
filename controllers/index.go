package controllers

import (
    "fmt"
    "net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hello, Udacity!")
}