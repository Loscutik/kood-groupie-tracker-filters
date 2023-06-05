package api

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestGetAndUnmarshalJSON(t *testing.T) {
	want := Api{
		Artists:   "https://groupietrackers.herokuapp.com/api/artists",
		Locations: "https://groupietrackers.herokuapp.com/api/locations",
		Dates:     "https://groupietrackers.herokuapp.com/api/dates",
		Relation:  "https://groupietrackers.herokuapp.com/api/relation",
	}

	netClient := http.Client{
		Timeout: time.Second * 10,
	}

	// res := response{}
	resa := Api{}
	err := GetAndUnmarshalJSON(&netClient, "https://groupietrackers.herokuapp.com/api", &resa)
	if err != nil {
		t.Fatal(err)
	}

	if reflect.ValueOf(want).Interface() != reflect.ValueOf(resa).Interface() {
		t.Fatalf("res is\n %#v\nwanted\n %#v\n", resa, want)
	}
}

func TestGetAtist(t *testing.T) {
	/*wLocations:= locations{
		Id:        11,
		Locations: []string{"doha-qatar", "minnesota-usa", "illinois-usa", "california-usa", "mumbai-india"},
	}

	wConcertDates:= dates{
		Id:    11,
		Dates: []string{"*15-12-2019", "*09-12-2019", "*07-12-2019", "*06-12-2019", "*16-11-2019"},
	}

	wRelations:= relations{
		Id: 11,
		DatesLocations: map[string][]string{
			"california-usa": {"06-12-2019"},
			"doha-qatar":     {"15-12-2019"},
			"illinois-usa":   {"07-12-2019"},
			"minnesota-usa":  {"09-12-2019"},
			"mumbai-india":   {"16-11-2019"},
		},
	}*/
	want := Artist{
		Id:           11,
		Image:        "https://groupietrackers.herokuapp.com/api/images/katyperry.jpeg",
		Name:         "Katy Perry",
		Members:      []string{"Katheryn Elizabeth Hudson"},
		CreationDate: 2001,
		FirstAlbum:   "04-10-2008",
		Locations: &Locations{
			Id:        11,
			Locations: []string{"doha-qatar", "minnesota-usa", "illinois-usa", "california-usa", "mumbai-india"},
			Dates:     &Dates{},
		},
		ConcertDates: "https://groupietrackers.herokuapp.com/api/dates/11",
		Relations:    "https://groupietrackers.herokuapp.com/api/relation/11",
	}
	netClient := http.Client{
		Timeout: time.Second * 10,
	}
	ap, err := GetAPI(&netClient, API)
	if err != nil {
		t.Fatal(err)
	}

	artists, err := GetArtists(&netClient, ap)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(artists[10].Locations.Dates)
	if want.Id != artists[10].Id ||
		want.Image != artists[10].Image ||
		want.Name != artists[10].Name ||
		want.Members[0] != artists[10].Members[0] ||
		want.CreationDate != artists[10].CreationDate ||
		want.FirstAlbum != artists[10].FirstAlbum ||
		want.Locations.Id != artists[10].Locations.Id ||
		want.Locations.Locations[0] != artists[10].Locations.Locations[0] ||
		want.ConcertDates != artists[10].ConcertDates ||
		want.Relations != artists[10].Relations {
		t.Fatalf("res is\n %v\nwanted\n %v\n", artists[10], want)
	}
}

func TestLocation(t *testing.T) {
	netClient := http.Client{
		Timeout: time.Second * 10,
	}
	ap, err := GetAPI(&netClient, API)
	if err != nil {
		t.Fatal(err)
	}

	byNum, err := GetLocationsNum(&netClient, ap, 2)
	if err != nil {
		t.Fatal("GetLocations - ", err)
	}

	bySet, err := GetLocationsSet(&netClient, ap)
	if err != nil {
		t.Fatal("GetLocationsSet - ", err)
	}

	byAddr := Locations{}
	err = GetAndUnmarshalJSON(&netClient, ap.Locations+"/2", &byAddr)
	if err != nil {
		t.Fatal("GetAndUnmarshalJSON - ", err)
	}

	if byAddr.Id == byNum.Id {
		t.Fatalf("err= %q\nlocatin obtained from GetLocations(&netClient, ap, 2) is \n%#v\nlocatin obtained from GetLocationsSet(&netClient,ap) is \n%#v\nlocatin obtained from GetAndUnmarshalJSON(&netClient, ap.Locations+/2, &byAddr) is \n%#v\n", err, byNum, bySet[1], byAddr)
	}
}

func TestData(t *testing.T) {
	netClient := http.Client{
		Timeout: time.Second * 10,
	}
	ap, err := GetAPI(&netClient, API)
	if err != nil {
		t.Fatal(err)
	}

	byNum, err := GetDatasNum(&netClient, ap, 2)
	if err != nil {
		t.Fatal("GetDatas - ", err)
	}

	bySet, err := GetDataSet(&netClient, ap)
	if err != nil {
		t.Fatal("GetDataSet - ", err)
	}

	byAddr := Dates{}
	err = GetAndUnmarshalJSON(&netClient, ap.Dates+"/2", &byAddr)
	if err != nil {
		t.Fatal("GetAndUnmarshalJSON - ", err)
	}

	if byAddr.Id == byNum.Id {
		t.Fatalf("err= %q\nlocatin obtained from GetDatas(&netClient, ap, 2) is \n%#v\nlocatin obtained from GetDataSet(&netClient,ap) is \n%#v\nlocatin obtained from GetAndUnmarshalJSON(&netClient, ap.Datas+/2, &byAddr) is \n%#v\n", err, byNum, bySet[1], byAddr)
	}
}

func TestRelation(t *testing.T) {
	netClient := http.Client{
		Timeout: time.Second * 10,
	}
	ap, err := GetAPI(&netClient, API)
	if err != nil {
		t.Fatal(err)
	}

	byNum, err := GetRelationNum(&netClient, ap, 2)
	if err != nil {
		t.Fatal("GetRelation - ", err)
	}

	bySet, err := GetRelationSet(&netClient, ap)
	if err != nil {
		t.Fatal("GetRelationSet - ", err)
	}

	byAddr := Relation{}
	err = GetAndUnmarshalJSON(&netClient, ap.Relation+"/2", &byAddr)
	if err != nil {
		t.Fatal("GetAndUnmarshalJSON - ", err)
	}

	if byAddr.Id == byNum.Id {
		t.Fatalf("err= %q\nlocatin obtained from GetRelation(&netClient, ap, 2) is \n%#v\nlocatin obtained from GetRelationSet(&netClient,ap) is \n%#v\nlocatin obtained from GetAndUnmarshalJSON(&netClient, ap.Relation+/2, &byAddr) is \n%#v\n", err, byNum, bySet[1], byAddr)
	}
}
