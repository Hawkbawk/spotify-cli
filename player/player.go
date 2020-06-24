package player

import (
	"github.com/Hawkbawk/spotify-cli/authentication"
	"github.com/zmb3/spotify"
)

// Player defines a struct whose client property can be
// used to control a user's playback on any of their devices.
type Player struct {
	client *spotify.Client
}

// NewPlayer returns a pointer to a Player, which can be used
// to control a user's Spotify playback.
func NewPlayer() *Player {
	return &Player{authentication.NewClient()}
}

// NextTrack skips to the next track in a user's Spotify
// queue. Requires a valid player, obtained by calling NewPlayer.
func (p Player) NextTrack() error {
	return p.client.Next()
}

// PreviousTrack skips back to the previous track in a user's Spotify
// queue, if a previous track exists.
func (p Player) PreviousTrack() error {
	return p.client.Previous()
}

// Pause pauses playback on the user's current active device.
func (p Player) Pause() error {
	return p.client.Pause()
}

// Play resumes playback on the user's current active device.
func (p Player) Play() error {
	return p.client.Play()
}
