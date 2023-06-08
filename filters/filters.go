package filters

import (
	"fmt"
	"strconv"
	"strings"

	"groupietracker/api"
)

type (
	filter   func(value string, data *api.Artist) bool
	filterOr func(values []int, data *api.Artist) bool
)

type SetGivenFilter struct {
	FilterFunc filter
	Value      string
}
type SetGivenFilterOr struct {
	FilterFunc filterOr
	Value      []int
}

func AddGivenFilter(givenFilters *[]SetGivenFilter, filterFunc filter, value string) {
	*givenFilters = append(*givenFilters, SetGivenFilter{
		FilterFunc: filterFunc,
		Value:      value,
	})
}

func AddGivenFilterOr(givenFilters *[]SetGivenFilterOr, filterFunc filterOr, value []int) {
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

func FilterFirstAlbumYearEq(value string, data *api.Artist) bool {
	year := strings.Split(data.FirstAlbum, "-")[2]
	return year == value
}

func FilterFirstAlbumYearLt(value string, data *api.Artist) bool {
	year := strings.Split(data.FirstAlbum, "-")[2]
	return year <= value
}

func FilterFirstAlbumYearGt(value string, data *api.Artist) bool {
	year := strings.Split(data.FirstAlbum, "-")[2]
	return year >= value
}

func FilterFirstAlbumMonthEq(value string, data *api.Artist) bool {
	date := strings.Split(data.FirstAlbum, "-")
	ym := date[2] + date[1]
	fmt.Printf("Filter album year=%s, year=%s", value, ym)
	return ym == value
}

func FilterFirstAlbumMonthLt(value string, data *api.Artist) bool {
	date := strings.Split(data.FirstAlbum, "-")
	ym := date[2] + date[1]
	return ym <= value
}

func FilterFirstAlbumMonthGt(value string, data *api.Artist) bool {
	date := strings.Split(data.FirstAlbum, "-")
	ym := date[2] + date[1]
	return ym >= value
}

func FilterFirstAlbumDateEq(value string, data *api.Artist) bool {
	date := strings.Split(data.FirstAlbum, "-")
	ymd := date[2] + date[1] + date[0]
	return ymd == value
}

func FilterFirstAlbumDateLt(value string, data *api.Artist) bool {
	date := strings.Split(data.FirstAlbum, "-")
	ymd := date[2] + date[1] + date[0]
	return ymd <= value
}

func FilterFirstAlbumDateGt(value string, data *api.Artist) bool {
	date := strings.Split(data.FirstAlbum, "-")
	ymd := date[2] + date[1] + date[0]
	return ymd >= value
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
