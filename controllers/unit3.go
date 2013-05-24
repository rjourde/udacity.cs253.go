package controllers

import (
	"html/template"
	"net/http"
	"appengine"
    "appengine/datastore"
	"time"
)

type Art struct {
	Title string
	Art string
	Created time.Time
	Error string `datastore:"-"`
}

func unit3AsciiChan(w http.ResponseWriter, r *http.Request) {
	
	if r.Method == "GET" {
		renderForm(w, r, Art{ "", "", time.Now(), ""})
	}
	if r.Method == "POST" {
		title := r.FormValue("title")
		art := r.FormValue("art")
		
		if len(title) <= 0 || len(art) <= 0 {
			error := "we need both a title and some artwork!"
			
			renderForm(w, r, Art{ Title: title, Art: art, Created: time.Now(), Error: error } )
		} else {
			// create new art
			c := appengine.NewContext(r)
			a := Art{ Title: title, Art: art, Created: time.Now(), Error: "" }
			
			_, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Art", nil), &a)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			
			http.Redirect(w, r, "/unit3/asciichan", http.StatusFound)
			return
		}
	}
}

func renderForm(w http.ResponseWriter, r *http.Request, a Art) {
	t, _ := template.ParseFiles("templates/front.html", "templates/arts.html")
	
	c := appengine.NewContext(r)
    q := datastore.NewQuery("Art").Order("-Created")
	
    var arts []*Art
    if _, err := q.GetAll(c, &arts); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
	
	if err := t.ExecuteTemplate(w, "tmpl_front", a); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := t.ExecuteTemplate(w, "tmpl_arts", arts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}