package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/btwiuse/pretty"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var (
	CLIENT_ID     = os.Getenv("GITHUB_CLIENT_ID")
	CLIENT_SECRET = os.Getenv("GITHUB_CLIENT_SECRET")
	REDIRECT_URL  = os.Getenv("OAUTH2_REDIRECT_URL")
	AUTHORIZE_URL = "https://github.com/login/oauth/authorize"
	TOKEN_URL     = "https://github.com/login/oauth/access_token"
)

func main() {
	http.HandleFunc("/", HandleIndex)
	http.HandleFunc("/oauth2/github/callback", HandleCallback)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, fmt.Sprintf(
		`<a target="_blank" href="https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s">Login<a>`,
		CLIENT_ID, REDIRECT_URL))
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
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
		// Handle Error
	}

	if !tkn.Valid() {
		log.Println("token invalid")
		return
		// Handle error
	}
	pretty.JSONLine(tkn)

	client := github.NewClient(oauthConfig.Client(oauth2.NoContext, tkn))
	user, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		log.Println("token invalid")
		return
		// Handle Error
	}

	// fmt.Printf("Name: %s\n", *user.Name) // Name: Your Name
	pretty.JSONLine(user)
	io.WriteString(w, pretty.JSONString(user))
	// http.Redirect(w, r, REDIRECT_URL, http.StatusPermanentRedirect)
}

func getScopes() []string {
	return []string{"user:email"}
}
