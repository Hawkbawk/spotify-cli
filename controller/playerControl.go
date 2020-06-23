package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/dixonwille/wmenu"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

const maxRetries = 5

type Player struct {
	client *spotify.Client
}

func NewPlayer() *Player {
	var tok []byte
	var err error
	retryCount := 0
	for tok, err = ioutil.ReadFile(tokFilePath); err != nil && retryCount < maxRetries; retryCount++ {
		authenticate()
	}
	if retryCount >= maxRetries {
		log.Fatal("Unable to authenticate. Please try again.")
	}
	var token *oauth2.Token

	if err := json.Unmarshal(tok, token); err != nil {
		log.Fatal("Couldn't parse the token. This is likely a bug. Please let me know on the project's Github page.")
	}
	client := auth.NewClient(token)
	return &Player{&client}
}

func (p Player) NextTrack(opt wmenu.Opt) error {
	return p.client.Next()
}

func (p Player) PreviousTrack(opt wmenu.Opt) error {
	return p.client.Previous()
}

func (p Player) Pause(opt wmenu.Opt) error {
	return p.client.Pause()
}

func (p Player) Play(opt wmenu.Opt) error {
	return p.client.Play()
}
