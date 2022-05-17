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
	ShortURL string
	LongURL  string
}

func newBulkInserter() *bulkInserter {
	b := bulkInserter{}
	return &b
}

func (b *bulkInserter) BulkInsert(db *sql.DB) {
	b.urlMu.Lock()

	insertQuery := sq.Insert("urls").Columns("long_url", "short_url", "updated_at", "created_at", "deleted_at")

	for _, urlSchema := range b.urlSchemas {
		insertQuery = insertQuery.Values(urlSchema.LongURL, urlSchema.ShortURL, time.Now(), time.Now(), nil)
	}
	_, err := insertQuery.RunWith(db).Exec()
	if err != nil {
		log.Println(err.Error())
	}
	b.urlSchemas = nil
	b.urlMu.Unlock()
}

func (b *bulkInserter) AppendUrlSchema(schema urlSchema) {
	b.urlMu.Lock()
	b.urlSchemas = append(b.urlSchemas, schema)
	b.urlMu.Unlock()
}
