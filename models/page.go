package models

import (
	"time"
	"errors"
	"net/http"
	"appengine"
	"appengine/datastore"
)

type Page struct {
	Id int64
	Name string
	Content string
	Created time.Time
}

func GetPage(r *http.Request, name string) (*Page, error) {
	var pages []*Page
	
	c := appengine.NewContext(r)
	// Fetch the page
	q:= datastore.NewQuery("Page").Filter("Name =", name)
	keys, err := q.GetAll(c, &pages)
	if err != nil {
		return nil, err
	}
	if(keys == nil || pages == nil) {
		return nil, errors.New(name + " page doesn't exit")
	}
	return pages[0], nil
}