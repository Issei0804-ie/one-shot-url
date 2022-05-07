package main

import (
	"one-shot-url/api"
	"one-shot-url/database"
	"one-shot-url/short"
	"one-shot-url/util"
)

func main() {
	util.InitEnv()
	util.InitLog()

	db := database.NewDB(false)
	s := short.NewShort(db.IsExistShortUrl)
	api := api.NewAPI(s, db)
	api.Run(8080)
}
