package rdb

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"log"
	"sync"
	"time"
)

type bulkInserter struct {
	urlSchemas []urlSchema
	urlMu      sync.Mutex
}

type urlSchema struct {
	ShortURL  string
	LongURL   string
	UpdatedAt *time.Time
	CreatedAt *time.Time
	DeletedAt *time.Time
}

func newBulkInserter() *bulkInserter {
	b := bulkInserter{}
	return &b
}

func (b *bulkInserter) BulkInsert(db *sql.DB) {
	log.Println("in BulkInsert function")
	b.urlMu.Lock()
	if len(b.urlSchemas) == 0 {
		b.urlMu.Unlock()
		return
	}

	insertQuery := sq.Insert("urls").Columns("long_url", "short_url", "updated_at", "created_at", "deleted_at")

	for _, urlSchema := range b.urlSchemas {
		insertQuery = insertQuery.Values(urlSchema.LongURL, urlSchema.ShortURL, time.Now(), time.Now(), nil)
	}
	_, err := insertQuery.RunWith(db).Exec()
	if err != nil {
		log.Println(err.Error())
	}
	b.urlSchemas = []urlSchema{}
	b.urlMu.Unlock()
}

func (b *bulkInserter) AppendUrlSchema(schema urlSchema) {
	b.urlMu.Lock()
	log.Println("append url schema")
	b.urlSchemas = append(b.urlSchemas, schema)
	b.urlMu.Unlock()
}
