package main

import (
	"fmt"

	"github.com/Hawkbawk/spotify-cli/controller"
)

func main() {
	fmt.Println("Going to authorize user")
	controller.Authenticate()
	fmt.Println("Authorized user")
}
