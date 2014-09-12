package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

// A struct matching the format of the JSON returned from YQL
type attendeeYQLResult struct {
	Query struct {
		Results struct {
			A []struct {
				Content string
			}
		}
	}
}

// Takes a lanyrd event id and prints a list of all the names of the attendees
func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: lanyrd_attendees <event_id>\nWhere event_id is the fragment of the url after lanyrd.com. eg To retrieve a list of attendees for the event whose lanyard page is http://lanyrd.com/2015/my-event use lanyrd_attendees \"2015/my-event\"")
		return
	}

	names, err := get_event_attendees(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := range names {
		fmt.Println(names[i])
	}
}

// Retrieves the attendees from the attendees pages for the event in lanyrd via YQL. The
// attendees list is paginated using the page= URL parameter, so we increment this by one
// each time until we're returned no results
func get_event_attendees(event_code string) ([]string, error) {

	var names []string
	page_number := 1

	for {

		fmt.Println("Retrieving attendees for " + event_code + " page " + strconv.Itoa(page_number))

		// Retrieve the list of attendees for the current page in json form from YQL
		resp, err := http.Get("https://query.yahooapis.com/v1/public/yql?q=select%20*%20from%20html%20where%20url%3D%22lanyrd.com%2F" + event_code + "%2Fattendees%2F%3Fpage%3D" + strconv.Itoa(page_number) + "%22%20and%20xpath%3D'%2F%2Fdiv%5B%40class%3D%22mini-profile%22%5D%2Fspan%2Fa'&format=json&callback=")
		if err != nil {
			return []string{}, err
		}
		defer resp.Body.Close()

		// Convert the json to our attendeeYQLStruct
		jsondata, _ := ioutil.ReadAll(resp.Body)
		var yql_results attendeeYQLResult
		err = json.Unmarshal(jsondata, &yql_results)
		if err != nil {
			return []string{}, err
		}

		// If there are no results for this page then we've reached the
		// last of the pages
		if len(yql_results.Query.Results.A) == 0 {
			return names, nil
		}

		// Append the names of the attendees onto our array
		for i := range yql_results.Query.Results.A {
			names = append(names, yql_results.Query.Results.A[i].Content)
		}

		// Increment the page number
		page_number++

	}

}
