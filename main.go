package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"groupietracker/api"
)

const (
	TEMPLATES_PATH = "./templates/"
)

type application struct {
	errLog  *log.Logger
	infoLog *log.Logger
	client  *http.Client
	apies   *api.Api
	artists *[]api.Artist
}

var app application

func main() {
	// Creates logs of what happened
	errLog := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)               // Creates logs of errors
	infoLogFile, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND, 0666) // Puts log info into the specific file
	if err != nil {
		errLog.Printf("Cannot open a log file. Error is %s\nStdout will be used for the info log ", err)
		infoLogFile = os.Stdout
	}
	infoLog := log.New(infoLogFile, "INFO:  ", log.Ldate|log.Ltime|log.Lshortfile)
	// Specifies the time limits for server requests
	netClient := http.Client{
		Timeout: time.Second * 20,
	}
	// Gets the data from given API
	apies, err := api.GetAPI(&netClient, api.API)
	if err != nil {
		log.Fatalf("fail creating api: %s", err)
	}
	infoLog.Printf("apies: %#v\n", apies)
	// Gets the data of the specific Artists list
	artists, err := api.GetArtists(&netClient, &apies)
	if err != nil {
		log.Fatalf("fail fetcing artists' data: %s", err)
	}
	infoLog.Printf("artists: %#v\n", artists)
	app = application{
		errLog:  errLog,
		infoLog: infoLog,
		client:  &netClient,
		apies:   &apies,
		artists: &artists,
	}
	// Handlers to run the web pages
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.homePageHandler)
	mux.HandleFunc("/info", app.concertsInfo)
	fileServer := http.FileServer(http.Dir(TEMPLATES_PATH + "static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	// Starting the web server
	port, err := parseArgs()
	if err != nil {
		errLog.Fatal(err)
	}
	fmt.Printf("Starting server at port %s\n", *port)
	infoLog.Printf("Starting server at port %s\n", *port)
	if err := http.ListenAndServe(":"+*port, mux); err != nil {
		errLog.Fatal(err)
	}
}

// The handler of the main page
func (app *application) homePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.errLog.Printf("wrong path for home page: %s", r.URL.Path)
		api.NotFound(w, r)
		return
	}
	// Assembling the page from templates
	api.GetTemplate("home.page.tmpl", w, app.artists)
}

// The handler of the website data
func (app *application) concertsInfo(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	app.infoLog.Printf("concertInfo - requested id: %s", id)
	i, err := strconv.Atoi(id)
	if err != nil || i < 1 || i > len(*app.artists) {
		app.errLog.Printf("wrong artist's id (%s), %s", id, err)
		api.NotFound(w, r)
		return
	}
	rel, err := api.GetArtistsRelation(app.client, &(*app.artists)[i-1])
	if err != nil {
		app.errLog.Printf("error during getting relations %s", err)
		api.NotFound(w, r)
		return
	}
	app.infoLog.Printf("relations for id %s is: %#v\n", id, rel)
	if rel.Id != i {
		app.errLog.Printf("relations' Id is different from the artist's Id %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Parses locations and the dates of concerts in each location
	type datesLocations []struct {
		Location string
		Dates    []string
	}
	datesLocs := make(datesLocations, len(rel.DatesLocations))
	j := 0
	for l, d := range rel.DatesLocations {
		l = strings.ReplaceAll(l, "_", " ")
		datesLocs[j].Location = strings.ReplaceAll(l, "-", ", ")
		datesLocs[j].Dates = d
		j++
	}
	output := struct {
		Id             int
		Image          string
		Name           string
		DatesLocations *datesLocations
	}{(*app.artists)[i-1].Id, (*app.artists)[i-1].Image, (*app.artists)[i-1].Name, &datesLocs}
	api.GetTemplate("info.page.tmpl", w, output)
}

// Parses the program's arguments to obtain the server port. If no arguments found, it uses the 8080 port by default
// Usage: go run .  --port=PORT_NUMBER
func parseArgs() (*string, error) {
	port := flag.String("port", "8080", "server port")
	flag.Parse()
	if flag.NArg() > 0 {
		return nil, fmt.Errorf("wrong arguments\nUsage: go run .  --port=PORT_NUMBER")
	}
	_, err := strconv.ParseUint(*port, 10, 16)
	if err != nil {
		return nil, fmt.Errorf("error: port must be a 16-bit unsigned number ")
	}
	return port, nil
}
