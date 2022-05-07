package database

import (
	"database/sql"
	"log"
	"one-shot-url/util"
	"testing"
)

func TestMain(m *testing.M) {
	util.InitEnv()
	util.InitLog()

	db := NewDB(true)
	log.Println("flash database")
	err := FlashTestData(db.db)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	log.Println("success")
	log.Println("start test")
	m.Run()
	log.Println("finish test")
	log.Println("flash database")
	err = FlashTestData(db.db)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	log.Println("success")
}

func TestDB_SearchLongURL(t *testing.T) {
	db := NewDB(true)
	err := FlashTestData(db.db)
	if err != nil {
		t.Fatalf(err.Error())
		return
	}
	type fields struct {
		db *sql.DB
	}
	type args struct {
		shortURL string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantLongURL string
		wantErr     bool
	}{
		// TODO: Add test cases.
		{
			name: "存在する shortURL を使用した正常系テスト",
			fields: fields{
				db: db.db,
			},
			args: args{
				shortURL: "00hu8jgt",
			},
			wantLongURL: "https://example.com",
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DB{
				db: tt.fields.db,
			}
			gotLongURL, err := d.SearchLongURL(tt.args.shortURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchLongURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotLongURL != tt.wantLongURL {
				t.Errorf("SearchLongURL() gotLongURL = %v, want %v", gotLongURL, tt.wantLongURL)
			}
		})
	}
}
