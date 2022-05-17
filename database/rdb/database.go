package rdb

import (
	"errors"
	sq "github.com/Masterminds/squirrel"
	"time"

	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
	"os"
)

type Interactor interface {
	Store(longURL string, shortURL string) error
	SearchLongURL(shortURL string) (longURL string, err error)
	GetDB() *sql.DB
	IsExistShortUrl(shortURL string) bool
}

func NewRDB(isTest bool) Interactor {
	var address, port, addr, user, passwd, dbName string
	if isTest {
		address = os.Getenv("TEST_RDB_ADDRESS")
		port = os.Getenv("TEST_RDB_PORT")
		addr = address + ":" + port
		user = os.Getenv("TEST_RDB_USER")
		passwd = os.Getenv("TEST_RDB_USER_PASSWORD")
		dbName = os.Getenv("TEST_RDB_NAME")
	} else {
		address = os.Getenv("RDB_ADDRESS")
		port = os.Getenv("RDB_PORT")
		addr = address + ":" + port
		user = os.Getenv("RDB_USER")
		passwd = os.Getenv("RDB_USER_PASSWORD")
		dbName = os.Getenv("RDB_NAME")
	}

	log.Println("RDB config is below.")
	log.Println("address is " + address)
	log.Println("port is " + port)
	log.Println("user is " + user)
	log.Println("db name is " + dbName)

	cfg := mysql.Config{
		User:   user,
		Passwd: passwd,
		Net:    "tcp",
		Addr:   addr,
		DBName: dbName,
	}

	log.Println("connect RDB now...")
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("connected RDB.")
	return DB{db: db}
}

type DB struct {
	db *sql.DB
}

func (d DB) Store(longURL string, shortURL string) error {
	tableName := d.makeTableName(shortURL)
	if !d.isExistTable(tableName) {
		err := d.createTable(tableName)
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

func (d DB) SearchLongURL(shortURL string) (longURL string, err error) {
	tableName := d.makeTableName(shortURL)
	row, err := sq.Select("long_url").From(tableName).Where("short_url = ?", shortURL).RunWith(d.db).Query()
	if err != nil {
		return "", err
	}
	if !row.Next() {
		return "", errors.New(shortURL + " is not found in database.")
	}
	err = row.Scan(&longURL)
	if err != nil {
		log.Println(err.Error())
		return "", errors.New("database error")
	}

	return longURL, nil
}

func (d DB) IsExistShortUrl(shortURL string) bool {
	tableName := d.makeTableName(shortURL)
	row, err := sq.Select("long_url").From(tableName).Where("short_url = ?", shortURL).Limit(1).RunWith(d.db).Query()
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return row.Next()
}

func (d DB) GetDB() *sql.DB {
	return d.db
}

func (d DB) createTable(tableName string) error {
	// SQL インジェクション が行える可能性があるため後ほど確認 or 書き換えましょう
	query := "CREATE TABLE " + tableName + " (id int NOT NULL PRIMARY KEY AUTO_INCREMENT, long_url varchar(1000), short_url VARCHAR(100), updated_at DATETIME, created_at DATETIME, deleted_at DATETIME)"
	_, err := d.db.Exec(query)

	if err != nil {
		log.Println("tableName is " + tableName + ", err is " + err.Error())
		return err
	}

	return nil
}

func (d DB) makeTableName(shortURL string) string {
	tableNameHead := shortURL[0:2]
	tableName := tableNameHead + "_urls"
	return tableName
}

func (d DB) isExistTable(tableName string) bool {
	query := "SELECT id FROM " + tableName + " LIMIT 1;"
	_, err := d.db.Exec(query)
	if err != nil {
		log.Println(tableName + " is not exist")
		return false
	}
	return true
}
