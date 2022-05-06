package main

import (
	"one-shot-url/api"
)

func main() {
	api := api.NewAPI()
	api.Run(80)
}
