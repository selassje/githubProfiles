package main

import (
	"flag"

	"github.com/selassje/githubProfiles/view"
)

func getUserName() (name string) {
	flag.StringVar(&name, "n", "selassje", "User name to provide info for")
	flag.Parse()
	return
}

func main() {
	view.RunGui()
}
