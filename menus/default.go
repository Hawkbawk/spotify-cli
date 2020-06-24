package menus

import (
	"log"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/zmb3/spotify"
)

const play string = "Play"
const stop string = "Stop"
const next string = "Next"
const previous string = "Previous"
const search string = "Search"
const quit string = "Quit"

// DisplayDefaultMenu displays the default menu for the CLI that the user
// interacts with. It lets the user choose what action they'd like to take,
// offering playback control, search functionality, and the ability to add
// tracks/episodes to the user's Spotify queue.
func DisplayDefaultMenu(client *spotify.Client) {
	actions := []string {play, stop, next, previous, search, quit}
	slct := promptui.Select{
		Label: "What do you want to do?",
		Items: actions,
		Searcher: func(input string, index int) bool {
			return strings.Contains(strings.ToLower(actions[index]), strings.ToLower(input))
		},
		HideSelected: true,
		CursorPos: 2,
	}

	_, action, err := slct.Run()

	if err != nil {
		log.Fatal(err)
	}

	switch action {
	case play:
		client.Play()
	case stop:
		client.Pause()
	case next:
		client.Next()
	case previous:
		client.Previous()
	case search:
		DisplaySearchMenu(client)
	case quit:
		os.Stdout.Close()
	}
}