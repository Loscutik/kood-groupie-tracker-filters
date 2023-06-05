package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const API = "https://groupietrackers.herokuapp.com/api"

type Api struct {
	Artists   string
	Locations string
	Dates     string
	Relation  string
}

type Artist struct {
	Id           int        `json:"id"`
	Image        string     `json:"image"`
	Name         string     `json:"name"`
	Members      []string   `json:"members"`
	CreationDate int        `json:"creationDate"`
	FirstAlbum   string     `json:"firstAlbum"`
	Locations    *Locations `json:"locations"`
	ConcertDates string     `json:"concertDates"`
	Relations    string     `json:"relations"`
}
type Locations struct {
	Id        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     *Dates   `json:"dates"`
}

type Dates struct {
	Id    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Relation struct {
	Id             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

/*
implements  the Unmarshaler interface for the Location type
*/
func (l *Locations) UnmarshalJSON(data []byte) error {
	netClient := http.Client{
		Timeout: time.Second * 20,
	}

	data = data[1 : len(data)-1]
	type tmpLoc Locations
	var tmp tmpLoc
	err := GetAndUnmarshalJSON(&netClient, string(data), &tmp)
	//*l = Locations(tmp)
	l.Id = tmp.Id
	l.Dates = tmp.Dates
	l.Locations = make([]string, len(tmp.Locations))
	for i := 0; i < len(tmp.Locations); i++ {
		l.Locations[i] = strings.ReplaceAll(strings.ReplaceAll(tmp.Locations[i], "-", ", "), "_", " ")
	}
	return err
}

/*
implements  the Unmarshaler interface for the Dates type
*/
func (d *Dates) UnmarshalJSON(data []byte) error {
	netClient := http.Client{
		Timeout: time.Second * 20,
	}

	data = data[1 : len(data)-1]
	type tmpDate Dates
	var tmp tmpDate
	err := GetAndUnmarshalJSON(&netClient, string(data), &tmp)
	*d = Dates(tmp)
	return err
}

/*
get data from `url` usung HTTP client `c` and parse them as JSON into parametr `v`. Parametr `v` must be a pointer
*/
func GetAndUnmarshalJSON(c *http.Client, url string, v any) error {
	if url == "" {
		return fmt.Errorf("no url is passed")
	}

	res, err := c.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// parse JSON
	return json.Unmarshal(data, v)
}

/*
creates an item of the struct Api and parse JSON from given url to the item
*/
func GetAPI(c *http.Client, url string) (Api, error) {
	api := Api{}
	err := GetAndUnmarshalJSON(c, url, &api)
	return api, err
}

/*
creates a slice of the struct Artist and parse JSON from the url kept in the `api` to the Artist item
*/
func GetArtists(c *http.Client, api Api) ([]*Artist, error) {
	var artists []*Artist
	err := GetAndUnmarshalJSON(c, api.Artists, &artists)
	return artists, err
}

/*
get data from `url` usung HTTP client `c` and parse them as JSON into parametr v. Parametr v must be a pointer
*/
func getSet(c *http.Client, url string, v any) error {
	// get JSON array which is a value of object "index"
	var tmp struct{ Index json.RawMessage }
	err := GetAndUnmarshalJSON(c, url, &tmp)
	if err != nil {
		return err
	}
	// parse set of locations
	return json.Unmarshal(tmp.Index, v)
}

/*
creates a slice of the struct Locations and parse JSON from the url kept in the `api` to the Locations item
*/
func GetLocationsSet(c *http.Client, api Api) ([]*Locations, error) {
	var locations []*Locations

	type tmpLoc Locations
	var tmp []*tmpLoc
	err := getSet(c, api.Locations, &tmp)
	locations = make([]*Locations, len(tmp))
	for i, l := range tmp {
		locations[i] = (*Locations)(l)
	}
	return locations, err
}

/*
creates a slice of the struct Dates and parse JSON from the url kept in the `api` to the Dates item
*/
func GetDataSet(c *http.Client, api Api) ([]*Dates, error) {
	var dates []*Dates
	err := getSet(c, api.Dates, &dates)
	return dates, err
}

/*
creates a slice of the struct Relation and parse JSON from the url kept in the `api` to the Relation item
*/
func GetRelationSet(c *http.Client, api Api) ([]*Relation, error) {
	var relations []*Relation
	err := getSet(c, api.Relation, &relations)
	return relations, err
}

/*
creates an item of the struct Artist
and get JSON data from the url kept in the `api` for the artist with given `num`.
Parse the data to the Artist item
*/
func GetArtist(c *http.Client, api Api, num int) (Artist, error) {
	var artist Artist
	err := GetAndUnmarshalJSON(c, api.Artists+"/"+strconv.Itoa(num), &artist)
	return artist, err
}

/*
creates an item of the struct Locations
and get JSON data from the url kept in the `api` for the locations with given `num`.
Parse the data to the Locations item
*/
func GetLocationsNum(c *http.Client, api Api, num int) (Locations, error) {
	var locations Locations
	err := GetAndUnmarshalJSON(c, api.Locations+"/"+strconv.Itoa(num), &locations)
	return locations, err
}

/*
creates an item of the struct Dates
and get JSON data from the url kept in the `api` for the date with given `num`.
Parse the data to the Dates item
*/
func GetDatasNum(c *http.Client, api Api, num int) (Dates, error) {
	var dates Dates
	err := GetAndUnmarshalJSON(c, api.Dates+"/"+strconv.Itoa(num), &dates)
	return dates, err
}

/*
creates an item of the struct Relation
and get JSON data from the url kept in the `api` for the relation with given `num`.
Parse the data to the Relation item
*/
func GetRelationNum(c *http.Client, api Api, num int) (Relation, error) {
	var relations Relation
	err := GetAndUnmarshalJSON(c, api.Relation+"/"+strconv.Itoa(num), &relations)
	return relations, err
}

/*
	returns Locations for the given Artist
	func GetArtistsLocations(c *http.Client, artist *Artist) (Locations, error) {
		var locations Locations
		err := GetAndUnmarshalJSON(c, artist.Locations, &locations)
		return locations, err
	}

*/
/*
	returns Data for the given Artist
*/
func GetArtistsDatas(c *http.Client, artist *Artist) (Dates, error) {
	var dates Dates
	err := GetAndUnmarshalJSON(c, artist.ConcertDates, &dates)
	return dates, err
}

/*
returns Relations for the given Artist
*/
func GetArtistsRelation(c *http.Client, artist *Artist) (Relation, error) {
	var relations Relation
	err := GetAndUnmarshalJSON(c, artist.Relations, &relations)
	return relations, err
}
