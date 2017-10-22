package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"io/ioutil"
	"bytes"
	"os"
)

// request

const (
	RequestedViewDetailedMerged = "DetailedMerged"
	AttendeeTypeRequired = "Required"
	AttendeeTypeOptional = "Optional"
)

type TimeZoneThingy struct {
	Bias int
	Time string
	DayOrder int
	Month int
	DayOfWeek string
}

type TimeZone struct {
	XMLName xml.Name `xml:"http://schemas.microsoft.com/exchange/services/2006/types TimeZone"`
	Bias int
	StandardTime TimeZoneThingy
	DaylightTime TimeZoneThingy
}

type Email struct {
	Address string
}

type MailboxData struct {
	XMLName xml.Name `xml:"http://schemas.microsoft.com/exchange/services/2006/types MailboxData"`
	Email Email
	AttendeeType string
	ExcludeConflicts bool
}

type MailboxDataArray struct {
	MailboxData []MailboxData
}

type TimeWindow struct {
	StartTime string
	EndTime string
}

type FreeBusyViewOptions struct {
	XMLName xml.Name `xml:"http://schemas.microsoft.com/exchange/services/2006/types FreeBusyViewOptions"`
	TimeWindow TimeWindow
	MergedFreeBusyIntervalInMinutes int
	RequestedView string
}

type GetUserAvailabilityRequest struct {
	XMLName xml.Name `xml:"http://schemas.microsoft.com/exchange/services/2006/messages GetUserAvailabilityRequest"`
	TimeZone TimeZone
	MailboxDataArray MailboxDataArray
	FreeBusyViewOptions FreeBusyViewOptions
}

type Envelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Body Body
}

type Body struct {
	GetUserAvailabilityRequest GetUserAvailabilityRequest
	GetUserAvailabilityResponse GetUserAvailabilityResponse
}

type ResponseMessage struct {
	ResponseClass string `xml:"ResponseClass,attr"`
	ResponseCode string
}

type WorkingPeriod struct {
	DayOfWeek string
	StartTimeInMinutes int
	EndTimeInMinutes int
}

type WorkingPeriodArray struct {
	WorkingPeriod []WorkingPeriod
}

type WorkingHours struct {
	TimeZone TimeZone
	WorkingPeriodArray WorkingPeriodArray
}

const (
	BusyTypeFree = "Free"
	BusyTypeTentative = "Tentative"
	BusyTypeBusy = "Busy"
	BusyTypeOOF = "OOF"
	BusyTypeNoData = "NoData"
)

type CalendarEventDetails struct {
	ID string
	Subject string
	Location string
	IsMeeting bool
	IsRecurring bool
	IsException bool
	IsReminder bool
	IsPrivate bool
}

type CalendarEvent struct {
	StartTime string
	EndTime string
	BusyType string
	CalendarEventDetails CalendarEventDetails
}

type CalendarEventArray struct {
	CalendarEvent []CalendarEvent
}

type FreeBusyView struct {
	FreeBusyViewType string
	MergedFreeBusy string
	CalendarEventArray CalendarEventArray
	WorkingHours WorkingHours
}

type FreeBusyResponse struct {
	ResponseMessage ResponseMessage
	FreeBusyView FreeBusyView
}

type FreeBusyResponseArray struct {
	FreeBusyResponse []FreeBusyResponse
}

type GetUserAvailabilityResponse struct {
	FreeBusyResponseArray FreeBusyResponseArray
}

func SerializeRequest(request GetUserAvailabilityRequest) ([]byte, error) {
	r := Envelope{
		Body: Body{
			GetUserAvailabilityRequest: request,
			GetUserAvailabilityResponse:GetUserAvailabilityResponse{},
		},
	}

	marshalled, err := xml.MarshalIndent(r, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	return marshalled, nil
}

func DeserializeResponse(bytes []byte) (GetUserAvailabilityResponse, error) {
	r := Envelope{}
	err := xml.Unmarshal(bytes, &r)
	if err != nil {
		return GetUserAvailabilityResponse{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}
	response := r.Body.GetUserAvailabilityResponse
	return response, nil
}

func GetUserAvailability(client *http.Client, url string, username string, password string, request GetUserAvailabilityRequest) (*GetUserAvailabilityResponse, error) {

	//fmt.Sprintf("querying page: %s\n", url)

	requestBody, err := SerializeRequest(request)
	if err != nil {
		return nil, err
	}

	//fmt.Printf("sending body : %s\n", string(requestBody))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to init request: %v", err)
	}
	req.SetBasicAuth(username, password)
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %v", err)
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to retreive response: %v", err)
	}

	//fmt.Printf("received body : %s\n", string(responseBody))

	response, err := DeserializeResponse(responseBody)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return &response, nil
}

func main() {

	searchEmail := os.Getenv("SEARCH_EMAIL")

	request := GetUserAvailabilityRequest{
		TimeZone: TimeZone{
			Bias: 480,
			StandardTime: TimeZoneThingy{
				Bias: 0,
				Time: "02:00:00",
				DayOrder: 5,
				Month: 10,
				DayOfWeek: "Sunday",
			},
			DaylightTime: TimeZoneThingy{
				Bias: -60,
				Time: "02:00:00",
				DayOrder: 1,
				Month: 4,
				DayOfWeek: "Sunday",
			},
		},
		MailboxDataArray: MailboxDataArray{
			MailboxData: []MailboxData{
				{
					Email:            Email{Address: searchEmail},
					AttendeeType:     AttendeeTypeRequired,
					ExcludeConflicts: false,
				},
			},
		},
		FreeBusyViewOptions: FreeBusyViewOptions{
			TimeWindow: TimeWindow{
				StartTime: "2017-10-23T00:00:00",
				EndTime: "2017-10-23T23:59:59",
			},
			MergedFreeBusyIntervalInMinutes: 60,
			RequestedView: RequestedViewDetailedMerged,
		},
	}

	client := &http.Client{}
	url := os.Getenv("EWS_URL")
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")

	response, err := GetUserAvailability(client, url, username, password, request)
	if err != nil {
		fmt.Printf("failed to retrieve availability: %v\n", err)
	} else {
		for idx, response := range response.FreeBusyResponseArray.FreeBusyResponse {
			fmt.Printf("respondent: %s\n", request.MailboxDataArray.MailboxData[idx].Email.Address)
			for _, event := range response.FreeBusyView.CalendarEventArray.CalendarEvent {
				fmt.Printf("\t%s, %s, %s\n", event.StartTime, event.EndTime, event.CalendarEventDetails.Subject)
			}
		}

	}
}

