package main

import (
	"github.com/joho/godotenv"
	"log"
	"one-shot-url/api"
	"one-shot-url/database"
	"path/filepath"
	"runtime"
)

func main() {
	_, pwd, _, _ := runtime.Caller(0)

	dir := filepath.Dir(pwd)
	log.Println(dir)
	err := godotenv.Load(dir + "/.env")
	if err != nil {
		log.Fatal(err)
	}
	log.SetFlags(log.Ltime | log.Llongfile)

	db := database.NewDB()
	err = db.Store("aaaaa", "aaaaaa")
	if err != nil {
		log.Printf(err.Error())
	}

	api := api.NewAPI()
	api.Run(8080)
}
