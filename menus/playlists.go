package menus

import (
	"fmt"
	"log"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/zmb3/spotify"
)

var (
	playlistsTemplate = promptui.SelectTemplates{
		Label:    "{{ . | white | bold }}",
		Active:   "{{ .Name | cyan | bold }}",
		Inactive: "{{ .Name | green | faint }}",
		Help:     "Movement: ← ↑ → ↓  ||  h j k l\tSearch: \"/\"",
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
		HideSelected: true,
	}

	index, _, err := list.Run()

	if err != nil {
		log.Fatal(err)
	}
	desiredPlaylist := playlists.Playlists[index]
	confirmListenToPlaylist(desiredPlaylist, client)
	DisplayHomeMenu(client)
}

func confirmListenToPlaylist(playlist spotify.SimplePlaylist, client *spotify.Client) {
	prompt := promptui.Prompt{
		Label:     "Do you want to listen to " + playlist.Name + "?",
		IsConfirm: true,
		Default:   "y",
		Templates: &confirmAddTracksTemplate,
	}

	_, err := prompt.Run()

	if err == nil {
		fmt.Println("Now listening to " + playlist.Name)

		playbackOptions := &spotify.PlayOptions{
			PlaybackContext: &playlist.URI,
		}
		err := client.PlayOpt(playbackOptions)
		if err != nil {
			log.Fatal(err)
		}

	} else {
		log.Fatal(err)
	}
}
