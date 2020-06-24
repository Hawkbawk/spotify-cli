package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/Hawkbawk/spotify-cli/authentication"
	"github.com/Hawkbawk/spotify-cli/menus"
	"github.com/manifoldco/promptui"
)

func main() {

	obtainPermission()

	client := authentication.NewClient()
	menus.DisplayHomeMenu(client)

}

func obtainPermission() {
	if _, err := ioutil.ReadFile(authentication.TokFilePath); err != nil {
		prompt := promptui.Prompt{
			Label:     "In order for this Spotify CLI to work, it needs to be connected to your Spotify account. Is that okay",
			Default:   "no",
			IsConfirm: true,
		}

		_, err := prompt.Run()
		if err != nil {
			log.Fatal("Permission denied. Exiting...")
		} else {
			fmt.Println("Permission granted, going to authenticate...")
		}
	}
}
