package models

import (
	"time"
	"errors"
	"net/http"
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"appengine/memcache"
	"fmt"
)

type Page struct {
	Id int64
	Name string
	Content string
	Created time.Time
}

func AddPage(r *http.Request, name string, content string) error {
	// create new page
	c := appengine.NewContext(r)
	pageID, _, _ := datastore.AllocateIDs(c, "Page", nil, 1)
	key := datastore.NewKey(c, "Page", "", pageID, nil)
	
	page := Page{ pageID, name, content, time.Now() }
	
	_, err := datastore.Put(c, key, &page)
	if err != nil {
		return err;
	}
	// Add the item to the memcache
	err := cache(r, page);
	
	return err;
}

func GetPage(r *http.Request, name string) (*Page, error) {
	var pages []*Page
	
	c := appengine.NewContext(r)
	if item, err := memcache.Get(c, "page"); err == memcache.ErrCacheMiss {
		// Fetch the page
		q:= datastore.NewQuery("Page").Filter("Name =", name)
		keys, err := q.GetAll(c, &pages)
		if err != nil {
			return nil, err
		}
		if(keys == nil || pages == nil) {
			return nil, errors.New(name + " page doesn't exit")
		}
	} else{
		if err := json.Unmarshal(item.Value, &pages); err != nil {
			return nil, err
		}
	}
	
	return pages[0], nil
}

func UpdatePage(r *http.Request, name string, content string) error {
	// update the datastore
	
	page := Page{ 0, name, content, time.Now() }
	
	// update memcache
	err := cache(r, page);

	return err
}

func cache(r *http.Request, page Page) error {
	// Create an Item
	item := &memcache.Item{
		Key:   fmt.Sprintf("page%d", page.Name),
		Value: json.Marshal(page),
	}
	// Set the item to the memcache
	c := appengine.NewContext(r)
	err := memcache.Set(c, item)
	
	return err
}