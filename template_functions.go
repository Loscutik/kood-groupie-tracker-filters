package main

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

const (
	TEMPLATES_PATH = "./templates/"
)

func GetTemplate(a string, w http.ResponseWriter, s any) {
	site := []string{
		TEMPLATES_PATH + "base.layout.tmpl",
		TEMPLATES_PATH + a,
	}

	funcMap := template.FuncMap{
		// The name "title" is what the function will be called in the template text.
		"nowYear": time.Now().Year,
	}
	tm, err := template.New("base.layout.tmpl").Funcs(funcMap).ParseFiles(site...)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatalln(err)
		return
	}
	err = tm.Execute(w, s)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatalln(err)
		return
	}
}

// Opens a beautiful HTML 404 web page instead of the status 404 "Page not found"
func NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)                             // Sets status code at 404
	tm, _ := template.ParseFiles(TEMPLATES_PATH + "error404.html") // Opens the HTML web page
	err := tm.Execute(w, nil)
	if err != nil {
		http.NotFound(w, r)
	}
}
