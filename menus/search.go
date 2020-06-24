package menus

import (
	"fmt"
	"log"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/zmb3/spotify"
)
// Maximum number of search results that a query can get.
var maxNumResults int = 20
// DisplaySearchMenu displays a menu that allows the user
// to search for a specific track on Spotify, then view the results.
func DisplaySearchMenu(client *spotify.Client) {
	prompt := promptui.Select{
		Label: "What do you want to search for?",
		Items: []string{"Track", "Artist"},
	}

	_, result, _ := prompt.Run()
	
	switch result {
	case "Track":
		prompt := promptui.Prompt {
			Label: "What track do you want to search for?",
		}
		result, err := prompt.Run()
		if err != nil {
			log.Fatal("Exiting...")
		}
		searchResults := searchForTrack(result, client)
		desiredTrack := displayTrackResults(searchResults, client)
		confirmAddTrackToQueue(desiredTrack, client)
	case "Artist":
		fmt.Println("Not yet implemented.")
	}
	
}

// DisplayTrackResults displays the results of searching for a track
// using promptui
func displayTrackResults(results []spotify.SimpleTrack, client *spotify.Client) *spotify.SimpleTrack {
	tracks := make([]string, maxNumResults)
	for i, track := range(results) {
		result := strings.Builder{}
		result.Write([]byte(track.Name))
		result.Write([]byte(" by: "))
		for _, artist := range(track.Artists) {
			result.WriteByte(' ')
			result.Write([]byte(artist.Name))
		}
		tracks[i] = result.String()
	}
	slct := promptui.Select{
		Label: "Track Search Results",
		Items: tracks,
	}

	index, _, err := slct.Run()

	if err != nil {
		log.Fatal("Exiting...")
	}

	return &results[index]
	
}

// confirmAddTrackToQueue confirms if a user wants to add the passed in
// track to their Spotify queue.
func confirmAddTrackToQueue(track *spotify.SimpleTrack, client *spotify.Client) {
	prompt := promptui.Prompt {
		Label: "Do you want to add " + track.Name + " to your queue?",
		IsConfirm: true,
		Default: "y",
	}

	_, err := prompt.Run()

	if err == nil {
		client.QueueSong(track.ID)
		fmt.Println("Added song " + track.Name + " to your queue!")
	} else {
		log.Fatal("Exiting...")
	}

	DisplayDefaultMenu(client)
}

// searchForTrack searches Spotify for a track who matches the keyword
// specified by query and returns a slice of tracks.
func searchForTrack(query string, client *spotify.Client) []spotify.SimpleTrack {
	options := &spotify.Options{Limit: &maxNumResults}
	searchResult, err := client.SearchOpt(query, spotify.SearchTypeTrack, options)
	if err != nil {
		log.Fatal(err)
	}
	results := make([]spotify.SimpleTrack, maxNumResults)
	for i, item := range(searchResult.Tracks.Tracks) {
		results[i] = item.SimpleTrack
	}
	return results
}