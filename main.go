package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Hawkbawk/spotify-cli/playerControl"
	"github.com/dixonwille/wmenu"
)

func main() {
	fmt.Println("Going to authorize user")
	player := playerControl.NewPlayer()
	menu := wmenu.NewMenu("What would you like to do?")
	menu.ChangeReaderWriter(os.Stdin, os.Stdout, os.Stderr)
	menu.Option("Play", "p", true, player.playerControl.Play)
	menu.Option("Pause", "s", false, player.playerControl.Pause)
	menu.Option("Next Track", "n", false, player.playerControl.NextTrack)
	menu.Option("Previous Track", "p", false, player.playerControl.PreviousTrack)
	if err := menu.Run(); err != nil {
		log.Fatal(err)
	}
}

