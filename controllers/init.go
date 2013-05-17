package controllers 

import "net/http"

func init() {
	http.HandleFunc("/", index)
}