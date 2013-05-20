package controllers

import (
	"fmt"
	"net/http"
)

func unit1(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, Udacity!")
}