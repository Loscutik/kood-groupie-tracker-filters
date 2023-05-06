package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

const API = "https://groupietrackers.herokuapp.com/api"

type Api struct {
	Artists   string
	Locations string
	Dates     string
	Relation  string
}

type Artist struct {
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}
type Locations struct {
	Id        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type Dates struct {
	Id    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Relation struct {
	Id             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

// Gets data from URL usung the HTTP client "c" and parses it as JSON into the parameter "v". The "v" parameter must be a pointer.
func GetAndUnmarshalJSON(c *http.Client, url string, v any) error {
	if url == "" {
		return fmt.Errorf("no url is passed")
	}
	res, err := c.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v) // JSON parsing
}

// Creates an item of the struct Api and parses JSON from given URL to the item
func GetAPI(c *http.Client, url string) (Api, error) {
	api := Api{}
	err := GetAndUnmarshalJSON(c, url, &api)
	return api, err
}

// Creates a slice of the struct Artist and parses JSON from the URL kept in the "api" to the Artist item
func GetArtists(c *http.Client, api *Api) ([]Artist, error) {
	var artists []Artist
	err := GetAndUnmarshalJSON(c, api.Artists, &artists)
	return artists, err
}

// Creates an item of the struct Artist and gets JSON data from the URL kept in the "api" for the artist with given "num".
// Parses the data to the Artist item
func GetArtist(c *http.Client, api *Api, num int) (Artist, error) {
	var artist Artist
	err := GetAndUnmarshalJSON(c, api.Artists+"/"+strconv.Itoa(num), &artist)
	return artist, err
}

// Returns Relations for the given Artist
func GetArtistsRelation(c *http.Client, artist *Artist) (Relation, error) {
	var relations Relation
	err := GetAndUnmarshalJSON(c, artist.Relations, &relations)
	return relations, err
}
