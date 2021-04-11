package server

import (
	"fmt"
	"os"
)

type DB struct {
	in       chan []byte
	out      chan string
	filename string
	file     *os.File
}

func newDB(filename string) *DB {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Can't open file ", filename, " for writing")
		os.Exit(1)
	}
	fmt.Println("Open DB file ", filename)
	return &DB{
		in:       make(chan []byte),
		out:      make(chan string),
		filename: filename,
		file:     file,
	}
}

func (db *DB) run() {
	fmt.Println("start DB server")
	for {
		d := <-db.in
		fmt.Printf("DB recive %d bytes of data: %s\n", len(d), string(d[0:15]))
		_, err := db.file.Write(d)
		if err != nil {
			fmt.Println("Can't write to file ", db.filename, " err: ", err)
			db.out <- "err"
		}
		_, err = db.file.Write([]byte("\n"))
		if err != nil {
			fmt.Println("Can't write to file ", db.filename, " err: ", err)
			db.out <- "err"
		}
		err = db.file.Sync()
		if err != nil {
			fmt.Println("Can't sync file ", db.filename, " err: ", err)
			db.out <- "err"
		}
		db.out <- "ok"
	}
}

func (db *DB) save(data []byte) string {
	db.in <- data
	resp := <-db.out
	return resp
}
