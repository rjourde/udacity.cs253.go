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
	r.HandleFunc("/blog", blogFrontPage)
	r.HandleFunc("/blog/newpost", blogNewPost)
	r.HandleFunc("/blog/{id:[0-9]+}", blogViewPost)
	r.HandleFunc("/signup", unit4Signup)
	r.HandleFunc("/login", unit4Login)
	r.HandleFunc("/logout", unit4Logout)
	r.HandleFunc("/welcome", unit4Welcome)
	r.HandleFunc("/blog.json", jsonBlogFrontPage)
	r.HandleFunc("/blog/{id:[0-9]+}.json", jsonBlogViewPost)
	http.Handle("/", r)
}