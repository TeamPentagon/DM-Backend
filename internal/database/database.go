package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path"

	_ "github.com/glebarez/go-sqlite"
	"github.com/syndtr/goleveldb/leveldb"
)

// Create a database connection to a file path
func CreateSQLiteDatabase(params ...string) (*sql.DB, error) {
	dir := buildPath(params...)
	fmt.Println(dir)
	db, err := sql.Open("sqlite", dir)
	if err != nil {
		log.Fatal(err)
	}
	return db, err
}

// Create a LevelDB database connection to a file path
func CreateLevelDBDatabase(params ...string) (*leveldb.DB, error) {
	dir := buildPath(params...)
	fmt.Println(dir)
	db, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		return nil, err
	}
	return db, err
}

// Helper function to build the directory path
func buildPath(params ...string) string {
	dir := path.Join(params[0])
	dir = path.Join(dir, params[1])
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatal(err)
		}
	}
	if len(params) > 2 {
		dir = path.Join(dir, params[2])
	}
	return dir
}
