package controllers 

import "net/http"

func init() {
	http.HandleFunc("/", index)
	http.HandleFunc("/unit1", unit1)
	http.HandleFunc("/unit2/rot13", unit2Rot13)
	http.HandleFunc("/unit2/signup", unit2Signup)
	http.HandleFunc("/unit2/welcome", unit2Welcome)
}