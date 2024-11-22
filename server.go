package main

import (
	"timedev/api"
)

func main() {

	// setup and run app
	err := api.SetupAndRunApp()
	if err != nil {
		panic(err)
	}
}
