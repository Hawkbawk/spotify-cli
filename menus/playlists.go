package menus

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/zmb3/spotify"
	"log"
	"os"
	"strings"
)

var (
	playlistsTemplate = promptui.SelectTemplates{
		Label:    "{{ . | green | bold }}",
		Active:   "{{ .Name | cyan | bold }}",
		Inactive: "{{ .Name | green | faint }}",
		Help:     "Movement: ← ↑ → ↓  ||  h j k l\tSearch: \"/\"",
		Selected: " ✅ {{ .Name | cyan | bold }}",
	}

	confirmAddTracksTemplate = promptui.PromptTemplates{
		Prompt: "{{ . | white | bold }}",
		Valid:  "{{ . | white | bold }}",
	}
)

// DisplayPlaylistsMenu displays a list of the current users playlist
// (up to their first ten playlists) and let's them choose what playlist
// they would like to listen to, if they so choose.
func DisplayPlaylistsMenu(client *spotify.Client) {
	playlists, err := client.CurrentUsersPlaylists()

	if err != nil {
		log.Fatal(err)
	}

	list := promptui.Select{
		Label:     "What playlist do you want to listen to?",
		Templates: &playlistsTemplate,
		Items:     playlists.Playlists,
		Searcher: func(input string, index int) bool {
			return strings.Contains(strings.ToLower(playlists.Playlists[index].Name), strings.ToLower(input))
		},
	}

	index, _, err := list.Run()

	if err != nil {
		os.Exit(0)
	}
	desiredPlaylist := playlists.Playlists[index]
	prompt := promptui.Prompt{
		Label:     "Do you want to listen to " + desiredPlaylist.Name + "?",
		IsConfirm: true,
		Default:   "y",
		Templates: &confirmAddTracksTemplate,
	}

	_, err = prompt.Run()

	if err == nil {
		fmt.Println("Now listening to " + desiredPlaylist.Name)

		playbackOptions := &spotify.PlayOptions{
			PlaybackContext: &desiredPlaylist.URI,
		}
		err := client.PlayOpt(playbackOptions)
		if err != nil {
			log.Fatal(err)
		}

	} else {
		os.Exit(0)
	}
	DisplayHomeMenu(client)
}
