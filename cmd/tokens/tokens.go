/*
	Customized example copied from the oauth2 lib docs, useful for getting a new access token.
*/

package main

import (
	"context"
	"golang.org/x/oauth2"
	"log"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {

	clientId := os.Getenv("CLIENT_ID")
	secret := os.Getenv("CLIENT_SECRET")

	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: secret,
		Scopes:       []string{"Calendars.Read", "User.Read"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://login.microsoftonline.com/common/oauth2/v2.0/authorize",
			TokenURL: "https://login.microsoftonline.com/common/oauth2/v2.0/token",
		},
		RedirectURL: "http://localhost:8081/example",
	}

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth dialog: %v\n", url)

	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}

	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s\n", tok.AccessToken)
	}

	client := conf.Client(ctx, tok)
	resp, resperr := client.Get("https://graph.microsoft.com/v1.0/me/events/")
	if resperr != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("%s", body)
}
