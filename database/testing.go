package database

import (
	"database/sql"
	"github.com/tanimutomo/sqlfile"
	"log"
)

// FlashTestData Remove and Add test data/***
func FlashTestData(d *sql.DB) error {
	err := removeTestData(d)
	if err != nil {
		log.Println("failed removeTestData")
		return err
	}

	err = addTestData(d)
	if err != nil {
		log.Println("failed addTestData")
		return err
	}
	return nil
}

func addTestData(d *sql.DB) error {
	s := sqlfile.New()
	err := s.File("./../DB/test-data.sql")
	if err != nil {
		return err
	}
	_, err = s.Exec(d)
	if err != nil {
		return err
	}
	return nil
}

func removeTestData(d *sql.DB) error {
	_, err := d.Exec("DROP DATABASE IF EXISTS one_shot_url_test")
	if err != nil {
		return err
	}
	_, err = d.Exec("CREATE DATABASE IF NOT EXISTS one_shot_url_test")
	return nil
}
