package main

import (
	"log"
	"one-shot-url/api"
	"one-shot-url/database"
	"one-shot-url/short"
	"one-shot-url/util"
	"os"
	"strconv"
)

func main() {
	util.InitEnv()
	util.InitLog()
	rawPort := os.Getenv("PORT")
	port, err := strconv.Atoi(rawPort)
	if err != nil {
		log.Fatalln("can not get a port number that environment value. did you modify .env?")
	}

	db := database.NewDB(false)
	s := short.NewShort(db.IsExistShortUrl)
	api := api.NewAPI(s, db, port)
	api.Run(port)
}
