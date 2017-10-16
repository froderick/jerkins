/*
	Here I'm experimenting with fetching calendar data.
 */
package main

import (
	"log"
	"fmt"
	"github.com/froderick/jerkins"
)

func main() {

	client, err := jerkins.Client()
	if err != nil {
		log.Fatal(err)
	}

	events := jerkins.QueryEvents(client, "$select=subject,start,end,attendees&$orderby=start/dateTime&$filter=start/dateTime gt '2017-10-15'")

	for _, event := range events {
		fmt.Printf("%+v : \"%+v\"\n", event.Start, event.Subject)
	}
}
