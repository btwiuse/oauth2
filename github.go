package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/btwiuse/pretty"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	CLIENT_ID     = "e448c84e25b2ed5aa242"
	CLIENT_SECRET = "88c18d789caade09222bb61cede5d13276cbfd53"
	REDIRECT_URL  = "https://hub.k0s.io/oauth2/github/callback"
	AUTHORIZE_URL = "https://github.com/login/oauth/authorize"
	TOKEN_URL     = "https://github.com/login/oauth/access_token"
)

func main() {
	local()
	os.Exit(0)
	http.HandleFunc("/auth", HandleCallback)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func local() {
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

	code := os.Args[1]

	tkn, err := oauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		// Handle Error
		log.Fatalln(err)
	}

	if !tkn.Valid() {
		// Handle error
		log.Fatalln("token invalid")
	}

	log.Println("code:", code)
	log.Println("token:", pretty.JSONString(tkn))

	client := github.NewClient(oauthConfig.Client(oauth2.NoContext, tkn))
	user, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		// Handle Error
		log.Fatalln("token invalid")
	}

	// fmt.Printf("Name: %s\n", *user.Name) // Name: Your Name
	pretty.JSON(user)
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
