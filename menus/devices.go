package menus

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/zmb3/spotify"
)

var (
	deviceListTemplate = promptui.SelectTemplates{
		Label:    "{{ . | white | bold }}",
		Active:   "🎛️  {{ .Name | cyan | bold }}",
		Inactive: "🎛️  {{ .Name | green | faint }}",
		Details: `
{{ "- - - Device Details - - -" | white | bold }}
{{ "Device Name:" | white | bold }}  {{ .Name | cyan | bold }}
{{ "Device Active:" | white | bold }}  {{ .Active | yellow | bold }}
{{ "Device Type:" | white | bold}}  {{ .Type | green | bold }}`,
		Help: "Movement: ← ↑ → ↓  ||  h j k l\tSearch: \"/\"",
	}
)

// DisplayDeviceMenu displays a list of the user's currently available devices.
// The user is then given the option to select a device and then switch
// playback to that device.
func DisplayDeviceMenu(client *spotify.Client) {
	devices, err := client.PlayerDevices()

	if err != nil {
		log.Fatal("Couldn't list devices: ", err)
	}

	stdout := os.Stdout
	deviceList := promptui.Select{
		Label: "Current Devices",
		Items: devices,
		Searcher: func(input string, index int) bool {
			existsInName := strings.Contains(strings.ToLower(devices[index].Name), strings.ToLower(input))
			existsInType := strings.Contains(strings.ToLower(devices[index].Type), strings.ToLower(input))
			return existsInName || existsInType
		},
		HideSelected: true,
		Templates:    &deviceListTemplate,
		Stdout:       stdout,
	}

	index, _, err := deviceList.Run()

	if err != nil {
		stdout.Truncate(0)
		log.Fatal(err)
	}

	confirmSwitchPlayback(devices[index], client)
	DisplayHomeMenu(client)
}

func confirmSwitchPlayback(device spotify.PlayerDevice, client *spotify.Client) {
	confirmSwitch := promptui.Prompt{
		Label:     "Do you want to switch playback to " + device.Name + "?",
		IsConfirm: true,
		Default:   "n",
	}

	_, err := confirmSwitch.Run()

	if err == nil {
		err = client.TransferPlayback(device.ID, true)

		if err != nil {
			log.Fatal("Couldn't transfer playback: ", err)
		}

	} else {
		fmt.Println("Not switching playback.")
	}
}
