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
)

const (
	TOKEN_PATH = "token.json"
)

func SaveToken(t *oauth2.Token) error {
	f, err := os.Create(TOKEN_PATH)
	if err != nil {
		return fmt.Errorf("unable to save oauth token: %v", err)
	}
	fmt.Printf("saved oauth token: %v\n", t)
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
		Scopes:       []string{"Calendars.Read", "User.Read"},
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
		fmt.Printf("%s\n", tok.AccessToken)
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
