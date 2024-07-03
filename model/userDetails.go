package model

import (
	"fmt"
	"time"

	"github.com/atomedgesoft/calendariq/config"
)

type User struct {
	Id            string    `json:id`
	FirstName     string    `json:firstname`
	LastName      string    `json:lastname`
	EmailAddress  string    `json:emailaddress`
	Signinthrough string    `json:signinthrough`
	CreatedAt     time.Time `json:createdate`
	TimeZone      string    `json:timezone`
	IsActive      bool      `json:isactive`
	Country       string    `json:country`
}

func InsertUser(user User) string {
	db, _ := config.ConnectDB()
	defer db.Close()

	sqlStmt := `insert into users(id, firstname, lastname, emailaddress, signinthrough, createdat, timezone, country, isactive)values($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`

	err := db.QueryRow(sqlStmt, user.Id, user.FirstName, user.LastName, user.EmailAddress, user.Signinthrough, user.CreatedAt, user.TimeZone, user.Country, user.IsActive).Scan(&user.Id)
	if err != nil {
		fmt.Println(err)
	}
	return user.Id
}

// getting user from the db
func getUser(user User) {
	fmt.Println(user)

}
