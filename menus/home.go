package menus

import (
	"log"
	"os"
	"strings"

	"github.com/Hawkbawk/spotify-cli/types"
	"github.com/manifoldco/promptui"
	"github.com/zmb3/spotify"
)

// A list of all the possible actions a user can take
// on the default menu.
var (
	togglePlayback = types.Action{
		Action: "Toggle Playback",
		Icon:   "⏯️ ",
	}
	next = types.Action{
		Action: "Next Song",
		Icon:   "⏭️ ",
	}
	previous = types.Action{
		Action: "Previous Song",
		Icon:   "⏮️ ",
	}
	playlists = types.Action{
		Action: "Select Playlist",
		Icon:   "💿",
	}
	search = types.Action{
		Action: "Search Spotify",
		Icon:   "🔎",
	}
	quit = types.Action{
		Action: "Quit Program",
		Icon:   "🚫",
	}
)

// The template for the default menu
var homeMenuTemplate = promptui.SelectTemplates{
	Label:    "🎼   {{ .Title | cyan }} {{ \"by:\" | red }} {{ .Artist | green }}  🎼",
	Active:   "{{ .Icon }}  {{ .Action | cyan | bold }}",
	Inactive: "{{ .Icon }}  {{ .Action | green | faint }}",
	Help: "Movement: ← ↑ → ↓  ||  h j k l	Search: \"/\"",
	FuncMap: promptui.FuncMap,
}

// DisplayHomeMenu displays the default menu for the CLI that the user
// interacts with. It lets the user choose what Action they'd like to take,
// offering playback control, search functionality, and the ability to add
// tracks/episodes to the user's Spotify queue.
func DisplayHomeMenu(client *spotify.Client) {
	actions := []types.Action{togglePlayback, next, previous, playlists, search, quit}
	currentlyPlaying := getCurrentPlayingTrack(client)
	prompt := promptui.Select{
		Label: currentlyPlaying,
		Items: actions,
		Searcher: func(input string, index int) bool {
			return strings.Contains(strings.ToLower(actions[index].Action), strings.ToLower(input))
		},
		HideSelected: true,
		CursorPos:    1, // default the cursor to next song, cause that's what most people will want
		Templates:    &homeMenuTemplate,
	}

	i, _, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	switch i {
	case 0:
		if state, _ := client.PlayerState(); state.Playing {
			err = client.Pause()
		} else {
			err = client.Play()
		}
		if err != nil {
			log.Fatal("Couldn't adjust playback. Error: " + err.Error())
		}
	case 1:
		if err = client.Next(); err != nil {
			log.Fatal("Couldn't skip forward. Error: " + err.Error())
		}
	case 2:
		if err = client.Previous(); err != nil {
			log.Fatal("Couldn't skip back. Error: " + err.Error())
		}
	case 3:
		DisplayPlaylistsMenu(client)
	case 4:
		DisplaySearchMenu(client)
	case 5:
		os.Exit(0)
	}

	DisplayHomeMenu(client)
}

func getCurrentPlayingTrack(client *spotify.Client) *types.Track {

	playing, _ := client.PlayerCurrentlyPlaying()
	var currentlyPlaying types.Track
	if playing.Item == nil {
		currentlyPlaying = types.Track{
			Title:  "nothing",
			Artist: "nobody",
		}
	} else {
		currentlyPlaying = types.NewTrack(&playing.Item.SimpleTrack)
	}
	return &currentlyPlaying
}
