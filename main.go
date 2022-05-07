package main

import (
	"log"
	"one-shot-url/api"
	"one-shot-url/database"
	"one-shot-url/util"
)

func main() {
	util.InitEnv()
	util.InitLog()

	db := database.NewDB(false)
	err := db.Store("aaaaa", "aaaaaa")
	if err != nil {
		log.Printf(err.Error())
	}

	api := api.NewAPI()
	api.Run(8080)
}
