package controllers 

import (
	"net/http"
	"github.com/gorilla/mux"
)

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/", index)
	r.HandleFunc("/unit1", unit1)
	r.HandleFunc("/unit2/rot13", unit2Rot13)
	r.HandleFunc("/unit2/signup", unit2Signup)
	r.HandleFunc("/unit2/welcome", unit2Welcome)
	r.HandleFunc("/unit3/asciichan", unit3AsciiChan)
	r.HandleFunc("/unit3/blog", blogFrontPage)
	r.HandleFunc("/unit3/blog/newpost", blogNewPost)
	r.HandleFunc("/unit3/blog/{id:[0-9]+}", blogViewPost)
	r.HandleFunc("/unit4/signup", unit4Signup)
	r.HandleFunc("/unit4/login", unit4Login)
	r.HandleFunc("/unit4/welcome", unit4Welcome)
	http.Handle("/", r)
}