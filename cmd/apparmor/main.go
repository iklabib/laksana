package main

import (
	"errors"
	"fmt"

	"codeberg.org/iklabib/markisa/profiles/apparmor"
)

func main() {
	err := apparmor.GenerateProfile("markisa")
	if err != nil {
		err = errors.New("failed to generate apparmor profile: " + err.Error())
		panic(err)
	}
	fmt.Println("apparmor profile generated")
}
