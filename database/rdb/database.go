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
	b := newBulkInserter()
	go func() {
		for {
			b.BulkInsert(db)
			time.Sleep(time.Second * 2)
		}
	}()
	return DB{
		db:   db,
		bulk: b,
	}
}

type DB struct {
	db   *sql.DB
	bulk *bulkInserter
}

func (d DB) Store(longURL string, shortURL string) error {
	now := time.Now()
	schema := urlSchema{
		ShortURL:  shortURL,
		LongURL:   longURL,
		UpdatedAt: &now,
		CreatedAt: &now,
		DeletedAt: nil,
	}
	d.bulk.AppendUrlSchema(schema)
	return nil
}

func (d DB) SearchLongURL(shortURL string) (longURL string, err error) {
	row, err := sq.Select("long_url").From("urls").Where("short_url = ?", shortURL).RunWith(d.db).Query()
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

func (d DB) GetDB() *sql.DB {
	return d.db
}
