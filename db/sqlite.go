package db

import (
	"database/sql"
	"fmt"
	"log"
)

func initDB() {
	DB, err := sql.Open("sqlite3", "./notes.db")
	if err != nil {
		log.Fatal(err)
	}

	tableCreationQuery := `
	CREATE TABLE IF NOT EXISTS notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		description TEXT NOT NULL,
		completed INTEGER NOT NULL CHECK (completed IN (0, 1) )
	);`

	_, err = DB.Exec(tableCreationQuery)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Table created successfully (if not present)")
}
