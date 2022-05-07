package database

import (
	sq "github.com/Masterminds/squirrel"
	"time"

	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
	"os"
)

type Interactor interface {
	Store(longURL string, shortURL string) error
	SearchLongURL(shortURL string) (longURL string)
	SearchShortURL(longURL string) (shortURL string)
}

func NewDB() *DB {

	log.Println("connect RDB now...")
	address := os.Getenv("DB_ADDRESS")
	log.Println("RDB address is " + address)
	port := os.Getenv("DB_PORT")
	log.Println("RDB port is " + port)
	addr := address + ":" + port
	cfg := mysql.Config{
		User:   os.Getenv("DB_USER"),
		Passwd: os.Getenv("DB_USER_PASSWORD"),
		Net:    "tcp",
		Addr:   addr,
		DBName: os.Getenv("DB_NAME"),
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("connected RDB.")
	return &DB{db: db}
}

type DB struct {
	db *sql.DB
}

func (d DB) Store(longURL string, shortURL string) error {
	tableNameHead := shortURL[0:2]
	tableName := tableNameHead + "_urls"
	if !d.isExistTable(tableName) {
		err := d.makeTable(tableName)
		if err != nil {
			return err
		}
	}
	_, err := sq.Insert(tableName).Columns("long_url", "short_url", "updated_at", "created_at", "deleted_at").
		Values(longURL, shortURL, time.Now(), time.Now(), nil).RunWith(d.db).Exec()
	if err != nil {
		return err
	}
	return nil
}

func (d DB) SearchLongURL(shortURL string) (longURL string) {
	//TODO implement me
	panic("implement me")
}

func (d DB) SearchShortURL(longURL string) (shortURL string) {
	//TODO implement me
	panic("implement me")
}

func (d DB) makeTable(tableName string) error {
	// SQL インジェクション が行える可能性があるため後ほど確認 or 書き換えましょう
	query := "CREATE TABLE " + tableName + " (id int NOT NULL PRIMARY KEY AUTO_INCREMENT, long_url varchar(1000), short_url VARCHAR(100), updated_at DATETIME, created_at DATETIME, deleted_at DATETIME)"
	_, err := d.db.Exec(query)

	if err != nil {
		log.Println("tableName is " + tableName + ", err is " + err.Error())
		return err
	}

	return nil
}

func (d DB) isExistTable(tableName string) bool {
	query := "SELECT id FROM " + tableName + " LIMIT 1"
	_, err := d.db.Exec(query)
	if err != nil {
		log.Println(tableName + "is not exist")
		return false
	}
	return true
}
