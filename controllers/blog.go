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
)

type Post struct {
	Subject string
	Content string
	Created time.Time
}

func blogFrontPage(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		// Display all blog entries
		c := appengine.NewContext(r)
		q := datastore.NewQuery("Post").Order("-Created")
		
		var posts []*Post
		if _, err := q.GetAll(c, &posts); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		renderFrontPage(w, posts)
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
			p := Post{ subject, content, time.Now() }
			
			key, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Post", nil), &p)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// redirect to the page of the newly created post
			stringID := fmt.Sprintf("%d", key.IntID())
			http.Redirect(w, r, "/unit3/blog/" + stringID, http.StatusFound)
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

func renderFrontPage(w http.ResponseWriter, posts []*Post) {
	t, _ := template.ParseFiles("templates/blog.html")
	err := t.Execute(w, posts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func renderNewPostForm(w http.ResponseWriter, data interface{}) {
	t, _ := template.ParseFiles("templates/newpost.html")
	err := t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func renderPostView(w http.ResponseWriter, post Post) {
	t, _ := template.ParseFiles("templates/post.html")
	err := t.Execute(w, post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}


