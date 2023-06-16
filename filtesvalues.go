package main

import (
	"fmt"
	"strconv"
	"time"

	"groupietracker/filters"
)

type date struct {
	year, month, day string
}

type filtersValues struct {
	Name,
	CreationDateFrom,
	CreationDateTo,
	Location string
	FirstAlbumFrom,
	FirstAlbumTo date
	Members []int
}

func checkStringDate(date string, from, to int) bool {
	dateInt, err := strconv.Atoi(date)
	return err == nil && dateInt >= from && dateInt <= to
}

func (d *date) notEmptyYear() bool {
	return d.year != ""
}

func (d *date) notEmptyMonth() bool {
	return d.month != "" && d.month != "0"
}

func (d *date) notEmptyDay() bool {
	return d.day != "" && d.day != "0"
}

func (d *date) checkYear(from, to int) bool {
	return checkStringDate(d.year, from, to)
}

func (d *date) isValidMonth() bool {
	return checkStringDate(d.month, 1, 12)
}

func (d *date) isValidDay() bool {
	return checkStringDate(d.day, 1, 31)
}

func (d *date) mergeMonthYear() string {
	return d.year + d.month
}

func (d *date) mergeDayMonthYear() string {
	return d.year + d.month + d.day
}

/*
runs a filter for the given date.
It checks if the date consists of the year, month, day, or of the year and month, or of only a year.
Then it runs the filter for the given operation. The operation must be a one of the follow strings: "Eq", "Lt", "Gt".
*/
func runDateFilter(currentFilter date, operation string, filtersSet *[]filters.SetGivenFilter) error {
	operationFuncs := map[string](map[string]filters.Filter){
		"Date": {
			"Eq": filters.FilterFirstAlbumDateEq,
			"Lt": filters.FilterFirstAlbumDateLt,
			"Gt": filters.FilterFirstAlbumDateGt,
		},
		"Month": {
			"Eq": filters.FilterFirstAlbumMonthEq,
			"Lt": filters.FilterFirstAlbumMonthLt,
			"Gt": filters.FilterFirstAlbumMonthGt,
		},
		"Year": {
			"Eq": filters.FilterFirstAlbumYearEq,
			"Lt": filters.FilterFirstAlbumYearLt,
			"Gt": filters.FilterFirstAlbumYearGt,
		},
	}


	if currentFilter.notEmptyYear() {
		if !currentFilter.checkYear(1900, time.Now().Year()) {
			return fmt.Errorf("wrong years for album creation %s", currentFilter.year)
		}

		if currentFilter.notEmptyMonth() {
			if !currentFilter.isValidMonth() {
				return fmt.Errorf("wrong months for album creation %s.%s", currentFilter.month, currentFilter.year)
			}

			if currentFilter.notEmptyDay() {
				if !currentFilter.isValidDay() {
					return fmt.Errorf("wrong day for album creation %s.%s.%s", currentFilter.day, currentFilter.month, currentFilter.year)
				}


				filters.AddGivenFilter(filtersSet, operationFuncs["Date"][operation], currentFilter.mergeDayMonthYear())
			} else {
				filters.AddGivenFilter(filtersSet, operationFuncs["Month"][operation], currentFilter.mergeMonthYear())
			}
		} else {
			filters.AddGivenFilter(filtersSet, operationFuncs["Year"][operation], currentFilter.year)
		}
	}
	return nil
}
