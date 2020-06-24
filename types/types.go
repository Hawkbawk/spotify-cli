package types

import (
	"github.com/zmb3/spotify"
	"strings"
)

type Action struct {
	Action string
	Icon   string
}

type Track struct {
	Title  string
	Artist string
}

func NewTrack(track *spotify.SimpleTrack) Track {
	artists := strings.Builder{}
	for _, artist := range track.Artists {
		artists.Write([]byte(artist.Name))
		artists.WriteByte(' ')
	}

	return Track{
		Title:  track.Name,
		Artist: artists.String(),
	}
}
