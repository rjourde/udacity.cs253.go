package models

import (
	"time"
	"errors"
	"net/http"
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"appengine/memcache"
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
	err = cache(r, page);
	
	return err;
}

func GetPage(r *http.Request, name string) (*Page, error) {
	var page *Page
	
	c := appengine.NewContext(r)
	if item, err := memcache.Get(c, "page" + name); err == memcache.ErrCacheMiss {
		// Fetch the page
		q:= datastore.NewQuery("Page").Filter("Name =", name)
		
		var pages []*Page
		keys, err := q.GetAll(c, &pages)
		if err != nil {
			return nil, err
		}
		if(keys == nil || pages == nil) {
			return nil, errors.New(name + " page doesn't exit")
		}
		
		page = pages[0]
	} else{
		if err := json.Unmarshal(item.Value, &page); err != nil {
			return nil, err
		}
	}
	
	return page, nil
}

func UpdatePage(r *http.Request, page Page, name string, content string) error {
	page.Name = name
	page.Content = content
	
	// update the datastore
	c := appengine.NewContext(r)
	key := datastore.NewKey(c, "Page", "", page.Id, nil)
	
	_, err := datastore.Put(c, key, &page)
	if err != nil {
		return err;
	}
	
	// update memcache
	err = cache(r, page);

	return err
}

func cache(r *http.Request, page Page) error {
	// Create an Item
	encodedPage, _ := json.Marshal(page)
	
	item := &memcache.Item{
		Key:   "page" + page.Name,
		Value: encodedPage,
	}
	// Set the item to the memcache
	c := appengine.NewContext(r)
	err := memcache.Set(c, item)
	
	return err
}