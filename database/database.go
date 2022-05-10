package database

import (
	"errors"
	sq "github.com/Masterminds/squirrel"
	"sync"
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

func NewDB(isTest bool) Interactor {
	var address, port, addr, user, passwd, dbName string
	if isTest {
		address = os.Getenv("TEST_DB_ADDRESS")
		port = os.Getenv("TEST_DB_PORT")
		addr = address + ":" + port
		user = os.Getenv("TEST_DB_USER")
		passwd = os.Getenv("TEST_DB_USER_PASSWORD")
		dbName = os.Getenv("TEST_DB_NAME")
	} else {
		address = os.Getenv("DB_ADDRESS")
		port = os.Getenv("DB_PORT")
		addr = address + ":" + port
		user = os.Getenv("DB_USER")
		passwd = os.Getenv("DB_USER_PASSWORD")
		dbName = os.Getenv("DB_NAME")
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
	now := time.Now()
	var t []table
	mu := sync.Mutex{}
	return DB{
		db:             db,
		lastInsertTime: &now,
		tables:         &t,
		mu:             &mu,
	}
}

type DB struct {
	db             *sql.DB
	lastInsertTime *time.Time
	tables         *[]table
	mu             *sync.Mutex
}

type table struct {
	tableName string
	longURL   string
	shortURL  string
	updatedAt time.Time
	createdAt time.Time
	deletedAt *time.Time
}

func (d DB) Store(longURL string, shortURL string) error {
	tableName := d.makeTableName(shortURL)
	if !d.isExistTable(tableName) {
		err := d.createTable(tableName)
		if err != nil {
			return err
		}
	}

	t := table{
		tableName: tableName,
		longURL:   longURL,
		shortURL:  shortURL,
		updatedAt: time.Now(),
		createdAt: time.Now(),
		deletedAt: nil,
	}
	return d.bulkInsert(t)

}

func (d DB) bulkInsert(tt table) error {
	d.mu.Lock()
	*d.tables = append(*d.tables, tt)
	if len(*d.tables) >= 10 {
		log.Println("start bulk insert")
		tx, err := d.db.Begin()
		if err != nil {
			log.Println(err.Error())
			d.mu.Unlock()
			return errors.New("can not start transaction. this is database error")
		}
		for _, t := range *d.tables {
			_, err = sq.Insert(t.tableName).Columns("long_url", "short_url", "updated_at", "created_at", "deleted_at").
				Values(t.longURL, t.shortURL, time.Now(), time.Now(), nil).RunWith(tx).Exec()
			if err != nil {
				log.Println(err.Error())
				tx.Rollback()
				d.mu.Unlock()
				return errors.New("bulk insert error")
			}
		}
		err = tx.Commit()
		if err != nil {
			log.Println(err.Error())
			d.mu.Unlock()
			return errors.New("transaction error")
		}
		log.Println("bulk insert !!")
		*d.tables = []table{}
		d.mu.Unlock()
	} else {
		d.mu.Unlock()
		time.Sleep(time.Second * 1)
		d.mu.Lock()
		defer d.mu.Unlock()
		if d.isInTables(tt.shortURL) {
			_, err := sq.Insert(tt.tableName).Columns("long_url", "short_url", "updated_at", "created_at", "deleted_at").Values(tt.longURL, tt.shortURL, time.Now(), time.Now(), nil).RunWith(d.db).Exec()
			if err != nil {
				log.Println(err)
				return errors.New("insert error")
			}
			*d.tables = removeTable(*d.tables, tt)
		}
	}
	return nil
}

func (d DB) isInTables(shortURL string) bool {
	for _, t := range *d.tables {
		if shortURL == t.shortURL {
			return true
		}
	}
	return false
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

// tables の中から table.tableName を unique なリストで返します．
func uniqueTableNames(tables []table) []string {
	var tableNames []string
	for _, t := range tables {
		if tableNames == nil {
			tableNames = append(tableNames, t.tableName)
			continue
		}
		for _, tableName := range tableNames {
			if tableName == t.tableName {
				continue
			}
		}
		tableNames = append(tableNames, t.tableName)
	}
	return tableNames
}

func removeTable(source []table, target table) []table {
	var newTable []table
	for _, t := range source {
		if t.shortURL != target.shortURL {
			newTable = append(newTable, t)
		}
	}
	return newTable
}
