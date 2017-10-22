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

	events, err := jerkins.QueryAllFutureEvents(client)

	if err != nil {
		fmt.Printf("failed to fetch data: %v", err)
		return
	}

	for _, event := range events {
		if jerkins.ContainsAttendee(event, "jack") {
			fmt.Printf("%+v : \"%+v\"\n", event.Start, event.Subject)
			for _, at := range event.Attendees {
				fmt.Printf("\t%s\n", at.Email.Address)
			}
		}
	}
}
