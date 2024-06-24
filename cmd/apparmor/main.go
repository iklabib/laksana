package main

import (
	"errors"
	"fmt"

	"codeberg.org/iklabib/laksana/profiles/apparmor"
)

func main() {
	err := apparmor.GenerateProfile("laksana")
	if err != nil {
		err = errors.New("failed to generate apparmor profile: " + err.Error())
		panic(err)
	}
	fmt.Println("apparmor profile generated")
}
