package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	CLIENT_ID     = "YOUR_CLIENT_ID"
	CLIENT_SECRET = "YOUR_CLIENT_SECRET"
	REDIRECT_URL  = "https://SOME_REDIRECT_URL.com"
	AUTHORIZE_URL = "https://github.com/login/oauth/authorize"
	TOKEN_URL     = "https://github.com/login/oauth/access_token"
)

func main() {
	http.HandleFunc("/auth", HandleCallback)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func HandleCallback(w http.ResponseWriter, r *http.Request) {
	oauthConfig := &oauth2.Config{
		ClientID:     CLIENT_ID,
		ClientSecret: CLIENT_SECRET,
		Endpoint: oauth2.Endpoint{
			AuthURL:  AUTHORIZE_URL,
			TokenURL: TOKEN_URL,
		},
		RedirectURL: REDIRECT_URL,
		Scopes:      getScopes(),
	}

	tkn, err := oauthConfig.Exchange(oauth2.NoContext, r.URL.Query().Get("code"))
	if err != nil {
		// Handle Error
	}

	if !tkn.Valid() {
		// Handle error
	}

	client := github.NewClient(oauthConfig.Client(oauth2.NoContext, tkn))
	user, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		// Handle Error
	}

	fmt.Printf("Name: %s\n", *user.Name) // Name: Your Name
	http.Redirect(w, r, REDIRECT_URL, http.StatusPermanentRedirect)
}

func getScopes() []string {
	return []string{"user:email"}
}
