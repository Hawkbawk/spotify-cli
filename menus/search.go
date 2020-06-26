package menus

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Hawkbawk/spotify-cli/types"
	"github.com/manifoldco/promptui"
	"github.com/zmb3/spotify"
)

// searchTemplate is the template that is used by when
// rendering the default search menu.
var searchTemplate = promptui.SelectTemplates{
	Label:    "{{ . | white | bold }}",
	Active:   "{{ . | cyan | bold }}",
	Inactive: "{{ . | green | faint }}",
	Help:     "Movement: ← ↑ → ↓  ||  h j k l\tSearch: \"/\"",
}

// trackResultsTemplate is the template that is used when
// rendering the results of searching for a track.
var trackResultsTemplate = promptui.SelectTemplates{
	Label:    "{{ . | white | bold }}",
	Active:   "{{ .Title | cyan | bold }} {{ \"by:\" | red | bold }} {{ .Artist | cyan | bold }}",
	Inactive: "{{ .Title | green | faint }} {{ \"by:\" | red | faint }} {{ .Artist | green | faint }}",
	Help:     "Movement: ← ↑ → ↓  ||  h j k l",
}

var artistResultsTemplate = promptui.SelectTemplates{
	Label:    "{{ . | white | bold }}",
	Active:   "{{ .Name | cyan | bold }}",
	Inactive: "{{ .Name | green | faint }}",
	Help:     "Movement: ← ↑ → ↓  ||  h j k l",
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
		Label:        "What do you want to search for?",
		Items:        []string{"Track", "Artist"},
		Templates:    &searchTemplate,
		HideSelected: true,
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
		prompt := promptui.Prompt{
			Label:     "What artist do you want to search for?",
			Templates: &queryTemplate,
		}

		query, err := prompt.Run()
		if err != nil {
			os.Exit(0)
		}
		searchResults := searchForArtist(query, client)
		desiredArtist := displayArtistsResults(searchResults)
		confirmListenToArtist(desiredArtist, client)
	}

	DisplayHomeMenu(client)

}

// DisplayTrackResults displays the results of searching for a track
// using promptui
func displayTrackResults(results []spotify.SimpleTrack) spotify.SimpleTrack {
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
		HideSelected: true,
	}

	index, _, err := prompt.Run()

	if err != nil {
		log.Fatal(err)
	}

	return results[index]

}

func displayArtistsResults(artists []spotify.SimpleArtist) spotify.SimpleArtist {
	prompt := promptui.Select{
		Label:     "Artist Search Results",
		Items:     artists,
		Templates: &artistResultsTemplate,
		Searcher: func(input string, index int) bool {
			return strings.Contains(strings.ToLower(artists[index].Name), strings.ToLower(input))
		},
		HideSelected: true,
	}

	index, _, err := prompt.Run()

	if err != nil {
		log.Fatal(err)
	}

	return artists[index]
}

// confirmAddTrackToQueue confirms if a user wants to add the passed in
// track to their Spotify queue.
func confirmAddTrackToQueue(track spotify.SimpleTrack, client *spotify.Client) {
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

func confirmListenToArtist(artist spotify.SimpleArtist, client *spotify.Client) {
	prompt := promptui.Prompt{
		Label:     "Do you want to start listening to " + artist.Name + "?",
		IsConfirm: true,
		Default:   "n",
	}

	_, err := prompt.Run()

	if err == nil {
		options := &spotify.PlayOptions{
			PlaybackContext: &artist.URI,
		}
		if err = client.PlayOpt(options); err != nil {
			log.Fatal("Couldn't start listening to the specified artist. Error: ", err)
		}
	} else {
		log.Fatal(err)
	}

	DisplayHomeMenu(client)
}

// searchForTrack searches Spotify for a track who matches the keyword
// specified by query and returns a slice of tracks.
func searchForTrack(query string, client *spotify.Client) []spotify.SimpleTrack {
	searchResult, err := client.Search(query, spotify.SearchTypeTrack)
	if err != nil {
		log.Fatal("Failed with error: ", err)
	}
	results := make([]spotify.SimpleTrack, len(searchResult.Tracks.Tracks))
	for i, item := range searchResult.Tracks.Tracks {
		results[i] = item.SimpleTrack
	}
	return results
}

func searchForArtist(query string, client *spotify.Client) []spotify.SimpleArtist {
	searchResult, err := client.Search(query, spotify.SearchTypeArtist)
	if err != nil {
		log.Fatal("Failed with error: ", err)
	}

	results := make([]spotify.SimpleArtist, len(searchResult.Artists.Artists))

	for i, item := range searchResult.Artists.Artists {
		results[i] = item.SimpleArtist
	}

	return results

}
