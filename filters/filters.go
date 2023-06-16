package filters

import (
	"strconv"
	"strings"

	"groupietracker/api"
)

type (
	Filter   func(value string, data *api.Artist) bool
	FilterOr func(values []int, data *api.Artist) bool
)

type SetGivenFilter struct {
	FilterFunc Filter
	Value      string
}
type SetGivenFilterOr struct {
	FilterFunc FilterOr
	Value      []int
}

func AddGivenFilter(givenFilters *[]SetGivenFilter, filterFunc Filter, value string) {
	*givenFilters = append(*givenFilters, SetGivenFilter{
		FilterFunc: filterFunc,
		Value:      value,
	})
}

func AddGivenFilterOr(givenFilters *[]SetGivenFilterOr, filterFunc FilterOr, value []int) {
	*givenFilters = append(*givenFilters, SetGivenFilterOr{
		FilterFunc: filterFunc,
		Value:      value,
	})
}

func ApplyFilters(records []*api.Artist, filters []SetGivenFilter, filtersOr []SetGivenFilterOr) []*api.Artist {
	if len(filters) == 0 {
		return records
	}
	filteredRecords := make([]*api.Artist, 0, len(records))
	for _, r := range records {
		passedFilter := true
		for _, f := range filters {
			if !f.FilterFunc(f.Value, r) {
				passedFilter = false
				break
			}
		}

		for _, f := range filtersOr {
			if !f.FilterFunc(f.Value, r) {
				passedFilter = false
				break
			}
		}

		if passedFilter {
			filteredRecords = append(filteredRecords, r)
		}
	}
	return filteredRecords
}

func FilterYearCreatingEq(value string, data *api.Artist) bool {
	year, err := strconv.Atoi(value)
	// if value isn't correct, filter will not apply as though there is no filters at all
	if err != nil {
		return true
	}
	return data.CreationDate == year
}

func FilterYearCreatingLt(value string, data *api.Artist) bool {
	year, err := strconv.Atoi(value)
	// if value isn't correct, filter will not apply as though there is no filters at all
	if err != nil {
		return true
	}
	return data.CreationDate <= year
}

func FilterYearCreatingGt(value string, data *api.Artist) bool {
	year, err := strconv.Atoi(value)
	// if value isn't correct, filter will not apply as though there is no filters at all
	if err != nil {
		return true
	}
	return data.CreationDate >= year
}

func FilterFirstAlbumYearEq(givenYear string, data *api.Artist) bool {
	year := strings.Split(data.FirstAlbum, "-")[2]
	return year == givenYear
}

func FilterFirstAlbumYearLt(givenYear string, data *api.Artist) bool {
	year := strings.Split(data.FirstAlbum, "-")[2]
	return year <= givenYear
}

func FilterFirstAlbumYearGt(givenYear string, data *api.Artist) bool {
	year := strings.Split(data.FirstAlbum, "-")[2]
	return year >= givenYear
}

func FilterFirstAlbumMonthEq(givenYM string, data *api.Artist) bool {
	date := strings.Split(data.FirstAlbum, "-")
	ym := date[2] + date[1]
	return ym == givenYM
}

func FilterFirstAlbumMonthLt(givenYM string, data *api.Artist) bool {
	date := strings.Split(data.FirstAlbum, "-")
	ym := date[2] + date[1]
	return ym <= givenYM
}

func FilterFirstAlbumMonthGt(givenYM string, data *api.Artist) bool {
	date := strings.Split(data.FirstAlbum, "-")
	ym := date[2] + date[1]
	return ym >= givenYM
}

func FilterFirstAlbumDateEq(givenYMD string, data *api.Artist) bool {
	date := strings.Split(data.FirstAlbum, "-")
	ymd := date[2] + date[1] + date[0]
	return ymd == givenYMD
}

func FilterFirstAlbumDateLt(givenYMD string, data *api.Artist) bool {
	date := strings.Split(data.FirstAlbum, "-")
	ymd := date[2] + date[1] + date[0]
	return ymd <= givenYMD
}

func FilterFirstAlbumDateGt(givenYMD string, data *api.Artist) bool {
	date := strings.Split(data.FirstAlbum, "-")
	ymd := date[2] + date[1] + date[0]
	return ymd >= givenYMD
}

func FilterLocationContain(value string, data *api.Artist) bool {
	for _, l := range data.Locations.Locations {
		if strings.Contains(strings.ToLower(l), strings.ToLower(value)) {
			return true
		}
	}
	return false
}

func FilterNameContain(value string, data *api.Artist) bool {
	return strings.Contains(strings.ToLower(data.Name), strings.ToLower(value))
}

func FilterNumberOfMembers(values []int, data *api.Artist) bool {
	if len(values) == 0 {
		return true
	}
	for _, v := range values {
		if len(data.Members) == v {
			return true
		}
	}
	return false
}
