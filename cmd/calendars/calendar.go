/*
	Here I'm experimenting with fetching calendar data.
 */
package main

import (
	"log"
	"io/ioutil"
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/froderick/jerkins"
)


func main() {

    client, err := jerkins.Client()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "https://graph.microsoft.com/v1.0/me/events", nil)
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	fmt.Printf("%s\n", body)

	events := Events{}
	json.Unmarshal(body, &events)

	fmt.Println("-------------------")
	fmt.Printf("%+v\n", events)
}

type Event struct {
	Odata_etag string `json:"@odata.etag"`
	Attendees  []struct {
		EmailAddress struct {
			Address string `json:"address"`
			Name    string `json:"name"`
		} `json:"emailAddress"`
		Status struct {
			Response string `json:"response"`
			Time     string `json:"time"`
		} `json:"status"`
		Type string `json:"type"`
	} `json:"attendees"`
	Body struct {
		Content     string `json:"content"`
		ContentType string `json:"contentType"`
	} `json:"body"`
	BodyPreview     string        `json:"bodyPreview"`
	Categories      []interface{} `json:"categories"`
	ChangeKey       string        `json:"changeKey"`
	CreatedDateTime string        `json:"createdDateTime"`
	End             struct {
		DateTime string `json:"dateTime"`
		TimeZone string `json:"timeZone"`
	} `json:"end"`
	HasAttachments       bool   `json:"hasAttachments"`
	ICalUID              string `json:"iCalUId"`
	ID                   string `json:"id"`
	Importance           string `json:"importance"`
	IsAllDay             bool   `json:"isAllDay"`
	IsCancelled          bool   `json:"isCancelled"`
	IsOrganizer          bool   `json:"isOrganizer"`
	IsReminderOn         bool   `json:"isReminderOn"`
	LastModifiedDateTime string `json:"lastModifiedDateTime"`
	Location             struct {
		DisplayName string `json:"displayName"`
	} `json:"location"`
	OnlineMeetingURL string `json:"onlineMeetingUrl"`
	Organizer        struct {
		EmailAddress struct {
			Address string `json:"address"`
			Name    string `json:"name"`
		} `json:"emailAddress"`
	} `json:"organizer"`
	OriginalEndTimeZone   string `json:"originalEndTimeZone"`
	OriginalStartTimeZone string `json:"originalStartTimeZone"`
	Recurrence            struct {
		Pattern struct {
			DayOfMonth     int64    `json:"dayOfMonth"`
			DaysOfWeek     []string `json:"daysOfWeek"`
			FirstDayOfWeek string   `json:"firstDayOfWeek"`
			Index          string   `json:"index"`
			Interval       int64    `json:"interval"`
			Month          int64    `json:"month"`
			Type           string   `json:"type"`
		} `json:"pattern"`
		Range struct {
			EndDate             string `json:"endDate"`
			NumberOfOccurrences int64  `json:"numberOfOccurrences"`
			RecurrenceTimeZone  string `json:"recurrenceTimeZone"`
			StartDate           string `json:"startDate"`
			Type                string `json:"type"`
		} `json:"range"`
	} `json:"recurrence"`
	ReminderMinutesBeforeStart int64 `json:"reminderMinutesBeforeStart"`
	ResponseRequested          bool  `json:"responseRequested"`
	ResponseStatus             struct {
		Response string `json:"response"`
		Time     string `json:"time"`
	} `json:"responseStatus"`
	Sensitivity    string      `json:"sensitivity"`
	SeriesMasterID interface{} `json:"seriesMasterId"`
	ShowAs         string      `json:"showAs"`
	Start          struct {
		DateTime string `json:"dateTime"`
		TimeZone string `json:"timeZone"`
	} `json:"start"`
	Subject string `json:"subject"`
	Type    string `json:"type"`
	WebLink string `json:"webLink"`
}

type Events struct {
	Odata_context  string `json:"@odata.context"`
	Odata_nextLink string `json:"@odata.nextLink"`
	Value          []Event
}
