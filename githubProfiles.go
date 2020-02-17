package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/selassje/githubProfiles/controller"
)

func getUserName() (name string) {
	flag.StringVar(&name, "n", "selassje", "User name to provide info for")
	flag.Parse()
	return
}

func main() {
	userName := getUserName()
	user, err := controller.GetUserInfo(userName)
	if err != nil {
		log.Fatal("Could not perform query for username `" + userName + "`." + " with error `" + err.Error() + "`")
	}
	fmt.Print(user)
}
