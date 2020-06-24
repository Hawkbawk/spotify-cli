package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Hawkbawk/spotify-cli/authentication"
	"github.com/Hawkbawk/spotify-cli/player"
	"github.com/manifoldco/promptui"
)

const play string = "Play"
const stop string = "Stop"
const next string = "Next"
const previous string = "Previous"
const logout string = "Logout"

func main() {
	validate := func(input string) error {
		return nil
	}
	
	prompt := promptui.Prompt{
		Label:    "In order for this Spotify CLI to work, it needs to be connected to your Spotify account. Is that okay",
		Validate: validate,
		Default: "no",
		IsConfirm: true,
	}
	
	result, err := prompt.Run()
	if err != nil {
		log.Fatal("Permission denied. Exiting...")
		} else {
		fmt.Printf("You choose %q\n", result)
		fmt.Println("Permission granted, going to authenticate...")
	}
	player := player.NewPlayer()
	actions := []string {play, stop, next, previous, logout}
	slct := promptui.Select{
		Label: "What do you want to do?",
		Items: actions,
		Searcher: func(input string, index int) bool {
			return strings.Contains(strings.ToLower(actions[index]), strings.ToLower(input))
		},
	}

	_, action, err := slct.Run()

	if err != nil {
		log.Fatal(err)
	}

	switch action {
	case "Play":
		defer player.Play()
	case "Stop":
		defer player.Pause()
	case "Next":
		defer player.NextTrack()
	case "Previous":
		defer player.PreviousTrack()
	case "Logout":
		os.Remove(authentication.TokFilePath)
	}

}
