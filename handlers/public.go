package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	spotifyOAuth "golang.org/x/oauth2/spotify"

	utils "github.com/jessicabuzzelli/running-playlist-generator/utils"
)

var (
	state string

	clientId     = os.Getenv("SPOTIFY_ID")
	clientSecret = os.Getenv("SPOTIFY_SECRET")

	OAuthScopes = []string{
		"user-read-private",
		"playlist-modify-private",
		"playlist-modify-public",
		"playlist-read-collaborative",
		"user-read-playback-position",
		"user-read-playback-state",
	}

	OAuthConf = &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Scopes:       OAuthScopes,
		Endpoint:     spotifyOAuth.Endpoint,
		RedirectURL:  redirectURI,
	}
)

const (
	baseURL     = "https://api.spotify.com/v1"
	redirectURI = "http://localhost:8080/callback"
)

func Login(c *gin.Context) {
	writer := c.Writer
	request := c.Request

	state = utils.GenerateRandomString(32)

	c.Set("OAuthRequestState", state)

	url := OAuthConf.AuthCodeURL(state, oauth2.AccessTypeOffline)

	http.Redirect(writer, request, url, http.StatusSeeOther)

	return
}

func OAuthCallback(c *gin.Context) {
	params := c.Request.URL.Query()
	code := params.Get("code")
	returnedState := params.Get("state")

	// todo: persist state between requests or omit entirely
	// sentState := c.GetString("OAuthRequestState")
	sentState := state

	if returnedState != sentState {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OAuth state mismatch!"})
		return
	}

	token, err := OAuthConf.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}

	client = OAuthConf.Client(ctx, token)

	c.Redirect(http.StatusSeeOther, "/app/loginSuccessful")

	return
}
