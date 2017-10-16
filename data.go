/*
	Everything that's not alexa-related goes here for now.
 */
package jerkins

import (
	"golang.org/x/oauth2"
	"fmt"
	"os"
	"encoding/json"
	"net/http"
	"golang.org/x/net/context"
	"log"
	"io/ioutil"
	"net/url"
)

const (
	TOKEN_PATH = "token.json"
)

func SaveToken(t *oauth2.Token) error {
	f, err := os.Create(TOKEN_PATH)
	if err != nil {
		return fmt.Errorf("unable to save oauth token: %v", err)
	}
	//fmt.Printf("saved oauth token: %v\n", t)
	defer f.Close()

	// Encode the token and write to disk
	if err := json.NewEncoder(f).Encode(t); err != nil {
		return fmt.Errorf("could not encode oauth token: %v", err)
	}

	return nil
}

func LoadToken() (*oauth2.Token, error) {

	f, err := os.Open(TOKEN_PATH)
	if err != nil {
		return nil, fmt.Errorf("could not open cache file at %s: %v", TOKEN_PATH, err)
	}
	defer f.Close()

	// Decode the JSON token cache
	token := new(oauth2.Token)
	if err := json.NewDecoder(f).Decode(token); err != nil {
		return nil, fmt.Errorf("could not decode token in cache file at %s: %v", TOKEN_PATH, err)
	}
	return token, nil
}

func OAuthConfig() *oauth2.Config {

	clientId := os.Getenv("CLIENT_ID")
	secret := os.Getenv("CLIENT_SECRET")

	conf := &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: secret,
		Scopes:       []string{"Calendars.Read", "User.Read", "offline_access"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://login.microsoftonline.com/common/oauth2/v2.0/authorize",
			TokenURL: "https://login.microsoftonline.com/common/oauth2/v2.0/token",
		},
		RedirectURL: "http://localhost:8081/example",
	}

	return conf;
}

func Client() (*http.Client, error) {

	conf := OAuthConfig()
	ctx := context.Background()

	tok, err := LoadToken()
	if err != nil {
		return nil, fmt.Errorf("no token read: %v", err)
	} else {
		//fmt.Printf("%s\n", tok.AccessToken)
	}

	ts := conf.TokenSource(ctx, tok)
	newToken, err := ts.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to get token from source: %v", err)
	}

	err = SaveToken(newToken)
	if err != nil {
		return nil, fmt.Errorf("failed to save new token: %v", err)
	}

	return conf.Client(ctx, tok), nil
}

type event struct {
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

type events struct {
	Odata_context  string `json:"@odata.context"`
	Odata_nextLink string `json:"@odata.nextLink"`
	Value          []event
}

type Email struct {
	Address string
	Name string
}

type Status struct {
	Response string
	Time string
}

type Attendee struct {
	Email Email
	Status Status
}

type EventSummary struct {
	Start string
	End string
	Attendees []Attendee
	Subject string
}

func summarizeEvent(a event) EventSummary {

	b := EventSummary{}
	b.Start = a.Start.DateTime
	b.End = a.End.DateTime
	b.Subject = a.Subject
	b.Attendees = []Attendee{}

	for _, attendeeA := range a.Attendees {
		attendeeB := Attendee{
			Email: Email{Address: attendeeA.EmailAddress.Address, Name: attendeeA.EmailAddress.Name},
			Status: Status{Response: attendeeA.Status.Response, Time: attendeeA.Status.Time},
		}
		b.Attendees = append(b.Attendees, attendeeB)
	}

	return b
}

const (
	GRAPH_BASE_URL = "https://graph.microsoft.com/"
	EVENTS_API = GRAPH_BASE_URL + "v1.0/me/events/?"
)

func QueryEventsPage(client *http.Client, url string) *events {

	fmt.Sprintf("querying page: %s\n", url)

	req, err := http.NewRequest("GET", url, nil)

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	events := events{}
	json.Unmarshal(body, &events)

	return &events
}

func QueryEvents(client *http.Client, query string) []EventSummary {
	summaries := []EventSummary{}

	nextUrl := EVENTS_API + url.PathEscape(query)
	for len(nextUrl) > 0 {
		events := QueryEventsPage(client, nextUrl)
		for _, event := range events.Value {
			summaries = append(summaries, summarizeEvent(event))
		}
		nextUrl = events.Odata_nextLink
	}

	return summaries
}
