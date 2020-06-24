package authentication

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"syscall"

	"github.com/google/uuid"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

const redirectURI = "http://localhost:8080/redirect"

// TokFilePath is the path that the OAuth2 token is stored at.
var TokFilePath string = os.TempDir() + "tok.txt"

var (
	auth = spotify.NewAuthenticator(redirectURI,
		spotify.ScopeUserReadCurrentlyPlaying,
		spotify.ScopeUserReadPlaybackState,
		spotify.ScopeUserModifyPlaybackState)
	ch    = make(chan *spotify.Client)
	state = uuid.New().String()
)

// Authenticate authenticates the application with the user's account
// so that the application can interact with Spotify and their API on
// the user's behalf.
func authenticate() *spotify.Client {

	go http.HandleFunc("/redirect", completeAuth)

	go func() {
		// Opens the authorization URL in the user's browser of choice,
		// assuming OS X style command.
		url := auth.AuthURL(state)

		// TODO: Make this a generic open command so the program can be run on any machine.
		command := exec.Command("open", url)
		if err := command.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	go http.ListenAndServe(":8080", nil)
	client := <-ch
	return client
}

// CompleteAuth is called after the user has been redirected to localhost:8080/redirect
// and finishes the authorization process by obtaining a OAuth2 token that is then
// used to create a client for interacting with the Spotify API. This token is then
// cached in a file on disk in plaintext (this is a Hackweek project, not a piece of bank
// code) so that the user doesn't have to reauthenticate every time they want to use
// the CLI.
func completeAuth(response http.ResponseWriter, request *http.Request) {
	tok, err := auth.Token(state, request)
	if err != nil {
		http.Error(response, "Couldn't get valid token from Spotify", http.StatusForbidden)
		log.Fatal(err)
	}

	if responseState := request.FormValue("state"); responseState != state {
		http.NotFound(response, request)
		log.Fatalf("State mismatch: %s != %s", responseState, state)
	}

	// For caching the OAuth2 token on disk
	// TODO: Maybe there's a better way to safely store this token on the machine?
	jsonTok, err := json.Marshal(tok)
	syscall.Umask(0177)
	if err := ioutil.WriteFile(TokFilePath, jsonTok, 0700); err != nil {
		log.Fatal("Couldn't stash OAuth2 token. Error: ", err)
	}

	// Display a success message on the redirect page
	// and send the finished client over the channel.
	client := auth.NewClient(tok)
	data, err := ioutil.ReadFile("resources/home.html")
	io.WriteString(response, string(data))
	ch <- &client

}

// NewClient constructs a new Spotify client out of the cached OAuth2 token.
// The client is used to actually interact with Spotify.
func NewClient() *spotify.Client {
	if _, present := os.LookupEnv("SPOTIFY_ID"); !present {
		log.Fatal("You haven't set your environment variables correctly! See the README for more details.")
	}
	var client *spotify.Client
	var tok []byte
	var err error
	if tok, err = ioutil.ReadFile(TokFilePath); err != nil {
		client = authenticate()
	} else {
		var token *oauth2.Token = &oauth2.Token{}
		if err := json.Unmarshal(tok, token); err != nil {
			log.Fatal(err)
		}
		temp := auth.NewClient(token)
		client = &temp
	}
	return client
}
