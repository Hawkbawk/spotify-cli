package menus

import (
	"fmt"
	"github.com/Hawkbawk/spotify-cli/types"
	"log"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/zmb3/spotify"
)

// searchTemplate is the template that is used by when
// rendering the default search menu.
var searchTemplate = promptui.SelectTemplates{
	Label:    "{{ . | green | bold }}",
	Active:   "{{ . | cyan | bold }}",
	Inactive: "{{ . | green | faint }}",
	Help:     "Movement: ← ↑ → ↓  ||  h j k l\tSearch: \"/\"",
	FuncMap:  promptui.FuncMap,
}

// trackResultsTemplate is the template that is used when
// rendering the results of searching for a track.
var trackResultsTemplate = promptui.SelectTemplates{
	Label:    "{{ . | white | bold }}",
	Active:   "{{ .Title | cyan | bold }} {{ \"by:\" | red | bold }} {{ .Artist | cyan | bold }}",
	Inactive: "{{ .Title | green | faint }} {{ \"by:\" | red | faint }} {{ .Artist | green | faint }}",
	Help:     "Movement: ← ↑ → ↓  ||  h j k l",
	Selected: " ✅ {{ .Title | cyan | bold }} {{ \"by:\" | red | bold}} {{ .Artist | cyan | bold }}",
	FuncMap:  promptui.FuncMap,
}

var queryTemplate = promptui.PromptTemplates{
	Prompt:  "{{ . | white | bold }}",
	Success: "",
	Valid:   "{{ . | white | bold }}\t",
}

// DisplaySearchMenu displays a menu that allows the user
// to search for a specific track on Spotify, then view the results.
func DisplaySearchMenu(client *spotify.Client) {
	prompt := promptui.Select{
		Label:     "What do you want to search for?",
		Items:     []string{"Track", "Artist"},
		Templates: &searchTemplate,
	}

	_, action, err := prompt.Run()

	if err != nil {
		os.Exit(0)
	}

	switch action {
	case "Track":
		prompt := promptui.Prompt{
			Label:     "What track do you want to search for?",
			Templates: &queryTemplate,
		}
		query, err := prompt.Run()
		if err != nil {
			os.Exit(0)
		}
		searchResults := searchForTrack(query, client)
		desiredTrack := displayTrackResults(searchResults)
		confirmAddTrackToQueue(desiredTrack, client)
	case "Artist":
		fmt.Println("Not yet implemented.")
	}

}

// DisplayTrackResults displays the results of searching for a track
// using promptui
func displayTrackResults(results []spotify.SimpleTrack) *spotify.SimpleTrack {
	tracks := make([]types.Track, len(results))
	for i, track := range results {
		tracks[i] = types.NewTrack(&track)
	}

	prompt := promptui.Select{
		Label:     "Track Search Results",
		Items:     tracks,
		Templates: &trackResultsTemplate,
		Searcher: func(input string, index int) bool {
			existsInTitle := strings.Contains(strings.ToLower(tracks[index].Title), strings.ToLower(input))
			existsInArtist := strings.Contains(strings.ToLower(tracks[index].Artist), strings.ToLower(input))
			return existsInArtist || existsInTitle
		},
	}

	index, _, err := prompt.Run()

	if err != nil {
		os.Exit(0)
	}

	return &results[index]

}

// confirmAddTrackToQueue confirms if a user wants to add the passed in
// track to their Spotify queue.
func confirmAddTrackToQueue(track *spotify.SimpleTrack, client *spotify.Client) {
	prompt := promptui.Prompt{
		Label:     "Do you want to add " + track.Name + " to your queue?",
		IsConfirm: true,
		Default:   "y",
	}

	_, err := prompt.Run()

	if err == nil {
		if err = client.QueueSong(track.ID); err != nil {
			log.Fatal("Couldn't queue song.")
		}
		fmt.Println("Added song " + track.Name + " to your queue!")
	} else {
		os.Exit(0)
	}

	DisplayHomeMenu(client)
}

// searchForTrack searches Spotify for a track who matches the keyword
// specified by query and returns a slice of tracks.
func searchForTrack(query string, client *spotify.Client) []spotify.SimpleTrack {
	searchResult, err := client.Search(query, spotify.SearchTypeTrack)
	if err != nil {
		os.Exit(0)
	}
	results := make([]spotify.SimpleTrack, len(searchResult.Tracks.Tracks))
	for i, item := range searchResult.Tracks.Tracks {
		results[i] = item.SimpleTrack
	}
	return results
}
