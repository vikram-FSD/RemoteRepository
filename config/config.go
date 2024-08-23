package config

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

var (
	Active     = true
	Lang       = "en"
	Location   = "Asia/Calcutta"
	Timenow, _ = time.LoadLocation(Location)
)

const Charset = "abcdefghijklmnopqrstuvwxyz0123456789"

func ConnectDB() (*sql.DB, error) {
	const (
		host     = "localhost"
		port     = "5432"
		username = "postgres"
		password = "password"
		dbname   = "CalendarIQ"
	)
	SqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, username, password, dbname)
	db, err := sql.Open("postgres", SqlInfo)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
func CurrentDateTime(Locate string) string {
	Timenow, _ := time.LoadLocation(Locate)
	DefaultDate := time.Now().In(Timenow)
	return DefaultDate.Format("2006-01-02 15:04:05 MST")
}
