/*
	Useful for getting a token.json if you don't already have one.
*/

package main

import (
	"golang.org/x/oauth2"
	"log"
	"fmt"
	"github.com/froderick/jerkins"
	"golang.org/x/net/context"
)

func main() {

	ctx := context.Background()
	conf := jerkins.OAuthConfig()

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
		fmt.Printf("access-token: %s\n", tok.AccessToken)
		fmt.Printf("refresh-token: %s\n", tok.RefreshToken)
	}

	jerkins.SaveToken(tok)
}
