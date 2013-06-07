package controllers

import (
	"html/template"
	"net/http"
	"appengine"
	"appengine/datastore"
	"time"
	"github.com/gorilla/mux"
	"fmt"
	"strconv"
	"encoding/json"
	"appengine/memcache"
)

type Post struct {
	Id int64
	Subject string
	Content string
	Created time.Time
}

func blogFrontPage(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		// Display all blog entries
		var posts []*Post
		// Get the item from the memcache
		c := appengine.NewContext(r)
		if err := memcache.Gob.Get(c, "posts", &post); err == memcache.ErrCacheMiss {
			c := appengine.NewContext(r)
			q := datastore.NewQuery("Post").Order("-Created")
			
			keys, err := q.GetAll(c, &posts)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// Create an Item
			item := &memcache.Item{
				Key:   "posts",
				Value: posts,
			}
			// Add the item to the memcache
			if err := memcache.Gob.Set(c, item); err == memcache.ErrNotStored {
				c.Infof("item with key %q already exists", item.Key)
			} else if err != nil {
				c.Errorf("error adding item: %v", err)
			}
		}
		
		renderFrontPage(w, posts)
	}
}

func jsonBlogFrontPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Display all blog entries
		c := appengine.NewContext(r)
		q := datastore.NewQuery("Post").Order("-Created")

		var posts []*Post
		_, err := q.GetAll(c, &posts)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		renderJsonFrontPage(w, posts)
	}
}

func blogNewPost(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		// Display empty form
		data := struct {
			Subject string
			Content string
			Error string
		}{ "", "", "", }
		
		renderNewPostForm(w, data)
	}
	if r.Method == "POST" {
		subject := r.FormValue("subject")
		content := r.FormValue("content")
		
		if len(subject) <= 0 || len(content) <= 0 {
			error := "subject and content, please!"
			
			data := struct {
				Subject string
				Content string
				Error string
			}{
				subject,
				content,
				error,
			}
			
			renderNewPostForm(w, data)
		} else {
			// create new post
			c := appengine.NewContext(r)
			p := Post{ 0, subject, content, time.Now() }
			
			key := datastore.NewIncompleteKey(c, "Post", nil)
			p.Id = key.IntID()
			
			_, err := datastore.Put(c, key, &p)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// redirect to the page of the newly created post
			stringID := fmt.Sprintf("%d", p.Id)
			http.Redirect(w, r, "/blog/" + stringID, http.StatusFound)
			return
		}

	}
}

func blogViewPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	intID, _ := strconv.ParseInt(id, 10, 64)
	// fetch the post from its ID
	var post Post
	c := appengine.NewContext(r)
	key := datastore.NewKey(c, "Post", "", intID, nil)
	if err := datastore.Get(c, key, &post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	renderPostView(w, post)
}

func jsonBlogViewPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	intID, _ := strconv.ParseInt(id, 10, 64)
	// fetch the post from its ID
	var post Post
	c := appengine.NewContext(r)
	key := datastore.NewKey(c, "Post", "", intID, nil)
	if err := datastore.Get(c, key, &post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	renderJsonPostView(w, post)
}

func renderFrontPage(w http.ResponseWriter, posts []*Post) {
	t := template.ParseFiles("templates/blog.html")
	if err := t.Execute(w, post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func renderJsonFrontPage(w http.ResponseWriter, posts []*Post) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	
	data := make([]JsonResponse, len(posts))
	
	for index, post := range posts {
		response := JsonResponse {
			"subject": post.Subject,
			"content": post.Content,
			"created": post.Created.String(),
		}
		data[index] = response
	}
	 
	fmt.Fprint(w, data)
}

func renderNewPostForm(w http.ResponseWriter, data interface{}) {
	t, _ := template.ParseFiles("templates/newpost.html")
	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func renderPostView(w http.ResponseWriter, post Post) {
	t, _ := template.ParseFiles("templates/post.html")
	if err := t.Execute(w, post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func renderJsonPostView(w http.ResponseWriter, post Post) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	
	data := JsonResponse {
		"subject": post.Subject,
		"content": post.Content,
		"created": post.Created.String(),
	} 
	
	fmt.Fprint(w, data)
}

type JsonResponse map[string]interface{}

func (r JsonResponse) String() (s string) {
	b, err := json.Marshal(r)
	if err != nil {
			s = ""
			return
	}
	s = string(b)
	return
}
