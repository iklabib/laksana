package db

import (
	"database/sql"
	"log"

	_ "github.com/glebarez/go-sqlite"
)

type Record struct {
	db *sql.DB
}

func NewRecord() *Record {
	record := &Record{}
	record.Init()
	return record
}

func (rd *Record) Init() {
	db, err := sql.Open("sqlite", "local.db")
	if err != nil {
		log.Fatal("can't open record")
	}

	query := `CREATE TABLE IF NOT EXISTS containers(client_id TEXT PRIMARY KEY, container_id TEXT NOT NULL, last_interaction TEXT NOT NULL)`
	if _, err := db.Exec(query); err != nil {
		log.Fatal("failed to initiate containers table")
	}

	rd.db = db
}

func (rd *Record) Retrieve(clientId string) string {
	query := `SELECT container_id FROM containers WHERE client_id=?`
	rows, err := rd.db.Query(query, clientId)
	if (err != nil){
		log.Fatal(err)
	}

	if rows.Next() {
		var container_id string
		rows.Scan(&container_id)
		return container_id
	}

	return ""
}

func (rd *Record) Insert(clientId, containerId string) (sql.Result, error) {
	query := `INSERT INTO containers(client_id, container_id, last_interaction) VALUES(?, ?, datetime('now'))`
	return rd.db.Exec(query, clientId, containerId)
}

func (rd *Record) Update(clientId, containerId string) (sql.Result, error) {
	query := `UPDATE containers SET container_id = ?, last_interaction = datetime('now') WHERE client_id = ?`
	return rd.db.Exec(query, containerId, clientId)
}

func (rd *Record) Delete(clientId string) (sql.Result, error) {
	return rd.db.Exec(`DELETE FROM containers WHERE client_id = ?`, clientId)
}

func (rd *Record) CleanRecord() (sql.Result, error) {
	return rd.db.Exec("DELETE FROM containers")
}
