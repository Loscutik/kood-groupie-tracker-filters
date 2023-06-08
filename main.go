package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"groupietracker/api"
	"groupietracker/filters"
)

// constrains for filters (output on the page)
type filtersConstrains struct {
	MinYearCreation,
	MaxYearCreation int
	Members   []int
	Locations map[string]bool
}
type application struct {
	errLog            *log.Logger
	infoLog           *log.Logger
	client            *http.Client
	apies             api.Api
	artists           []*api.Artist
	filtersConstrains *filtersConstrains
}

var app application

func main() {
	// Creates logs of what happened
	errLog := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)                          // Creates logs of errors
	infoLogFile, err := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o664) // Puts log info into the specific file
	if err != nil {
		errLog.Printf("Cannot open a log file. Error is %s\nStdout will be used for the info log ", err)
		infoLogFile = os.Stdout
	}
	infoLog := log.New(infoLogFile, "INFO:  ", log.Ldate|log.Ltime|log.Lshortfile)

	app = application{
		errLog:  errLog,
		infoLog: infoLog,
	}

	// Specifies the time limits for server requests
	app.client = &http.Client{
		Timeout: time.Second * 20,
	}

	// Gets data from the given API
	app.apies, err = api.GetAPI(app.client, api.API)
	if err != nil {
		log.Fatalf("fail creating api: %s", err)
	}
	infoLog.Printf("apies: %#v\n", app.apies)

	// Get data of the specific Artists list
	app.artists, err = api.GetArtists(app.client, app.apies)
	if err != nil {
		log.Fatalf("fail fetcing artists' data: %s", err)
	}
	infoLog.Printf("artists: %#v\n", app.artists)

	// preparing data for filters
	app.filtersConstrains = &filtersConstrains{
		MinYearCreation: 3000,
		MaxYearCreation: 0,
		Locations:       make(map[string]bool),
	}

	// Get a list of locations for the filter
	for _, a := range app.artists {
		for _, l := range a.Locations.Locations {
			if _, ok := app.filtersConstrains.Locations[l]; !ok {
				app.filtersConstrains.Locations[l] = false
			}
		}
	}
	infoLog.Printf("locations' list: %#v\n", app.filtersConstrains.Locations)

	maxMembers := 0
	for _, a := range app.artists {
		if a.CreationDate < app.filtersConstrains.MinYearCreation {
			app.filtersConstrains.MinYearCreation = a.CreationDate
		} else {
			if a.CreationDate > app.filtersConstrains.MaxYearCreation {
				app.filtersConstrains.MaxYearCreation = a.CreationDate
			}
		}
		if len(a.Members) > maxMembers {
			maxMembers = len(a.Members)
		}
	}
	app.filtersConstrains.Members = make([]int, maxMembers)
	for i := 0; i < maxMembers; i++ {
		app.filtersConstrains.Members[i] = i + 1
	}
	infoLog.Printf("filters: %#v\n", app.filtersConstrains)

	// Handlers to run the web pages
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.homePageHandler)
	mux.HandleFunc("/info", app.concertsInfo)
	fileServer := http.FileServer(neuteredFileSystem{http.Dir(TEMPLATES_PATH + "static/")})
	// fileServer := http.FileServer( http.Dir(TEMPLATES_PATH + "static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	// Starting the web server
	port, err := parseArgs()
	if err != nil {
		errLog.Fatal(err)
	}
	fmt.Printf("Starting server at http://localhost:%s\n", *port)
	infoLog.Printf("Starting server at http://localhost:%s\n", *port)
	if err := http.ListenAndServe(":"+*port, mux); err != nil {
		errLog.Fatal(err)
	}
}

// this type doesn't allow FileServer to open directions in the static direction
type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}
	s, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if s.IsDir() {
		if _, err := nfs.fs.Open(filepath.Join(path, "index.html")); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}

	return f, nil
}

// The handler of the main page
func (app *application) homePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.errLog.Printf("wrong path for home page: %s", r.URL.Path)
		NotFound(w, r)
		return
	}
	// get filters from the http request
	type filtersValues struct {
		Name,
		CreationDateFrom,
		CreationDateTo,
		FirstAlbumYear,
		FirstAlbumMonth,
		FirstAlbumDay,
		Location string
		Members []int
	}
	var filtersValue filtersValues
	filtersValue.Members = make([]int, len(app.filtersConstrains.Members))
	var filtersSet []filters.SetGivenFilter
	var filtersOrSet []filters.SetGivenFilterOr

	// getting filters
	// filter by the name of groupe
	filtersValue.Name = r.URL.Query().Get("Name")
	if filtersValue.Name != "" {
		filters.AddGivenFilter(&filtersSet, filters.FilterNameContain, filtersValue.Name)
	}

	// filter by the number of members
	if r.URL.Query().Get("Members") != "" {
		membersFilters := r.URL.Query()["Members"]
		var valuesMembersFilter []int
		for _, mf := range membersFilters {
			i, err := strconv.Atoi(mf)
			if err == nil {
				valuesMembersFilter = append(valuesMembersFilter, i)
			}
		}
		for _, v := range valuesMembersFilter {
			filtersValue.Members[v-1] = 1
		}

		filters.AddGivenFilterOr(&filtersOrSet, filters.FilterNumberOfMembers, valuesMembersFilter)
	}

	// filter by the date of creation
	filtersValue.CreationDateFrom = r.URL.Query().Get("CreationDateFrom")
	filtersValue.CreationDateTo = r.URL.Query().Get("CreationDateTo")
	if filtersValue.CreationDateFrom == "" {
		if filtersValue.CreationDateTo != "" {
			filters.AddGivenFilter(&filtersSet, filters.FilterYearCreatingLt, filtersValue.CreationDateTo)
		}
	} else /*filtersValue.CreationDateFrom !=""*/ if filtersValue.CreationDateTo == "" {
		filters.AddGivenFilter(&filtersSet, filters.FilterYearCreatingGt, filtersValue.CreationDateFrom)
	} else { /*filtersValue.CreationDateFrom !="" && filtersValue.CreationDateTo !="" */
		switch {
		case filtersValue.CreationDateFrom == filtersValue.CreationDateTo:
			filters.AddGivenFilter(&filtersSet, filters.FilterYearCreatingEq, filtersValue.CreationDateFrom)
		case filtersValue.CreationDateFrom < filtersValue.CreationDateTo:
			filters.AddGivenFilter(&filtersSet, filters.FilterYearCreatingGt, filtersValue.CreationDateFrom)
			filters.AddGivenFilter(&filtersSet, filters.FilterYearCreatingLt, filtersValue.CreationDateTo)
		}
	}

	// filter by the first album date
	filtersValue.FirstAlbumYear = r.URL.Query().Get("FirstAlbumYear")
	filtersValue.FirstAlbumMonth = r.URL.Query().Get("FirstAlbumMonth")
	filtersValue.FirstAlbumDay = r.URL.Query().Get("FirstAlbumDay")
	if filtersValue.FirstAlbumYear != "" {
		if year, err := strconv.Atoi(filtersValue.FirstAlbumYear); err != nil || year < 1900 || year > time.Now().Year() {
			app.errLog.Printf("wrong year for album creation %s", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		value := filtersValue.FirstAlbumYear
		if filtersValue.FirstAlbumMonth != "" && filtersValue.FirstAlbumMonth != "0" {
			if month, err := strconv.Atoi(filtersValue.FirstAlbumMonth); err != nil || month < 1 || month > 12 {
				app.errLog.Printf("wrong month for album creation %s", err)
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}
			value += filtersValue.FirstAlbumMonth
			if filtersValue.FirstAlbumDay != "" && filtersValue.FirstAlbumDay != "0" {
				if day, err := strconv.Atoi(filtersValue.FirstAlbumMonth); err != nil || day < 1 || day > 31 {
					app.errLog.Printf("wrong month for album creation %s", err)
					http.Error(w, "Bad Request", http.StatusBadRequest)
					return
				}
				value += filtersValue.FirstAlbumDay
				filters.AddGivenFilter(&filtersSet, filters.FilterFirstAlbumDateEq, value)
			} else {
				filters.AddGivenFilter(&filtersSet, filters.FilterFirstAlbumMonthEq, value)
			}
		} else {
			filters.AddGivenFilter(&filtersSet, filters.FilterFirstAlbumYearEq, value)
		}
	}

	// filter by the location
	filtersValue.Location = r.URL.Query().Get("Location")
	if filtersValue.Location != "" {
		filters.AddGivenFilter(&filtersSet, filters.FilterLocationContain, filtersValue.Location)
	}

	// apply all filters to set of artists
	output := &struct {
		Artists           []*api.Artist
		FiltersConstrains *filtersConstrains
		FiltersValues     *filtersValues
	}{
		Artists:           filters.ApplyFilters(app.artists, filtersSet, filtersOrSet),
		FiltersConstrains: app.filtersConstrains,
		FiltersValues:     &filtersValue,
	}

	// Assembling the page from templates
	GetTemplate("home.page.tmpl", w, output)
}

// The handler of the website data
func (app *application) concertsInfo(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	app.infoLog.Printf("concertInfo - requested id: %s", id)
	i, err := strconv.Atoi(id)
	if err != nil || i < 1 || i > len(app.artists) {
		app.errLog.Printf("wrong artist's id (%s), %s", id, err)
		NotFound(w, r)
		return
	}
	rel, err := api.GetArtistsRelation(app.client, (app.artists)[i-1])
	if err != nil {
		app.errLog.Printf("error during getting relations %s", err)
		NotFound(w, r)
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
	}{(app.artists)[i-1].Id, (app.artists)[i-1].Image, (app.artists)[i-1].Name, &datesLocs}
	GetTemplate("info.page.tmpl", w, output)
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
