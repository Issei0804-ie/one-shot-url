package util

import (
	"github.com/joho/godotenv"
	"log"
	"path/filepath"
	"runtime"
)

func InitEnv() {
	_, pwd, _, _ := runtime.Caller(0)

	dir := filepath.Dir(pwd)
	log.Println(dir)
	err := godotenv.Load(dir + "/../.env")
	if err != nil {
		log.Fatal(err)
	}
}

func InitLog() {
	log.SetFlags(log.Ltime | log.Llongfile)
}
