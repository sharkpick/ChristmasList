package main

import (
	"database/sql"
	"log"
)

func createRecipientTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE ` + RecipientTable + ` (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"name" TEXT NOT NULL UNIQUE,
		"finished" integer
		);`
	log.Println("creating", RecipientTable)
	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Println(err)
		return
	}
	statement.Exec()
	log.Println(RecipientTable, "created")
}

func createGiftsTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE ` + GiftTable + ` (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"name" TEXT NOT NULL,
		"price" integer,
		"url" TEXT,
		"purchased" integer,
		"userid" integer,
		FOREIGN KEY(userid) REFERENCES ` + RecipientTable + `(id) 
		);`
	log.Println("creating", GiftTable)
	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Println(err)
		return
	}
	statement.Exec()
	log.Println(GiftTable, "created")
}
