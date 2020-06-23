package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"syscall"

	"github.com/google/uuid"
	"github.com/zmb3/spotify"
)

const redirectURI = "http://localhost:8080/redirect"
const tokFilePath = "resources/tok.txt"

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
func Authenticate() {
	var client *spotify.Client
	var playerState *spotify.PlayerState

	http.HandleFunc("/redirect", CompleteAuth)

	go func() {
		// Opens the authorization URL in the user's browser of choice,
		// assuming OS X style command.
		url := auth.AuthURL(state)
		
		// TODO: Make this a generic open command so the program can
		// be run on any machine.
		command := exec.Command("open", url)
		if err := command.Start(); err != nil {
			log.Fatal(err)
		}

		client = <-ch

		user, err := client.CurrentUser()
		if err != nil {
			fmt.Println(err)
		}
		playerState, err = client.PlayerState()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("You are logged in as User ", user.ID)

		fmt.Printf("Found device of type %s (%s) for playback\n", playerState.Device.Type, playerState.Device.Name)
	}()

	http.ListenAndServe(":8080", nil)
}

// CompleteAuth is called after the user has been redirected to localhost:8080/redirect
// and finishes the authorization process by obtaining a OAuth2 token that is then
// used to create a client for interacting with the Spotify API. This token is then
// cached in a file on disk in plaintext (this is a Hackweek project, not a piece of bank
// code) so that the user doesn't have to reauthenticate every time they want to use
// the CLI.
func CompleteAuth(response http.ResponseWriter, request *http.Request) {
	println("About to finish authorization")
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
	if err := ioutil.WriteFile(tokFilePath, jsonTok, 0700); err != nil {
		log.Fatal("Couldn't stash OAuth2 token. Error: ", err)
	}

	// Display a success message on the redirect page
	// and send the finished client over the channel.
	client := auth.NewClient(tok)
	fmt.Println("Successfully logged in")
	data, err := ioutil.ReadFile("resources/home.html")
	io.WriteString(response, string(data))
	ch <- &client

}
