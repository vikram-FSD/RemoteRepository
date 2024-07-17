package model

import (
	"log"

	"github.com/atomedgesoft/calendariq/config"
)

type User struct {
	Id            string `json:"id"`
	FirstName     string `json:"firstname"`
	LastName      string `json:"lastname"`
	EmailAddress  string `json:"emailaddress"`
	Signinthrough string `json:"signinthrough"`
	CreatedAt     string `json:"createdat"`
	TimeZone      string `json:"timezone"`
	IsActive      bool   `json:"isactive"`
	Country       string `json:"country"`
}

// POST API for USER
func InsertUser(user User) (string, error) {
	db, _ := config.ConnectDB()
	defer db.Close()

	sqlStmt := `insert into users(id, firstname, lastname, emailaddress, signinthrough, createdat,  timezone, country,isactive)values($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`

	err := db.QueryRow(sqlStmt, user.Id, user.FirstName, user.LastName, user.EmailAddress, user.Signinthrough, user.CreatedAt, user.TimeZone, user.Country, user.IsActive).Scan(&user.Id)
	if err != nil {
		return "", err
	}
	return user.Id, nil
}

// GET API for USER
func GetUser(user User) ([]User, error) {
	db, err := config.ConnectDB()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	//query for retrieving the users from the Database
	sqlStmt := `select * from users`
	rows, err := db.Query(sqlStmt)
	if err != nil {
		return nil, err
	}
	var Usercontainer []User
	for rows.Next() {
		var get User
		err := rows.Scan(&get.Id, &get.FirstName, &get.LastName, &get.EmailAddress, &get.Signinthrough, &get.TimeZone, &get.Country, &get.IsActive, &get.CreatedAt)
		if err != nil {
			return Usercontainer, err
		}
		Usercontainer = append(Usercontainer, get)
	}
	if err = rows.Err(); err != nil {
		return Usercontainer, err
	}
	defer rows.Close()
	return Usercontainer, nil
}
