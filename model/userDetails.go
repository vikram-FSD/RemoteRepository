package model

import (
	"fmt"

	"github.com/atomedgesoft/calendariq/config"
)

type User struct {
	Id            string `json:id`
	FirstName     string `json:firstname`
	LastName      string `json:lastname`
	EmailAddress  string `json:emailaddress`
	Signinthrough string `json:signinthrough`
	CreatedAt     string `json:createdate`
	TimeZone      string `json:timezone`
	IsActive      bool   `json:isactive`
	Country       string `json:country`
}

func InsertUser(user User) (string, error) {
	db, _ := config.ConnectDB()
	defer db.Close()

	sqlStmt := `insert into users(id, firstname, lastname, emailaddress, signinthrough,  timezone, country,isactive, createdat)values($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`

	err := db.QueryRow(sqlStmt, user.Id, user.FirstName, user.LastName, user.EmailAddress, user.Signinthrough, user.TimeZone, user.Country, user.IsActive, user.CreatedAt).Scan(&user.Id)
	if err != nil {
		return "", err
	}
	return user.FirstName, nil
}

// getting user from the db
func ReturnUser(user User) ([]User, error) {
	db, err := config.ConnectDB()
	if err != nil {
		fmt.Println(err)
	}
	//query for retrieving the users from the Database
	sqlStmt := `select * from users`
	rows, err := db.Query(sqlStmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
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
	return Usercontainer, nil
}

// Checking if the emailID is already exist or not
func IsEmailExists(user User) (email string, error error) {
	db, err := config.ConnectDB()
	if err != nil {
		fmt.Println(err)
	}
	sqlStmt := `select * from users`
	rows, err := db.Query(sqlStmt)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	var Usercontainer []User
	var get User
	for rows.Next() {

		err := rows.Scan(&get.Id, &get.FirstName, &get.LastName, &get.EmailAddress, &get.Signinthrough, &get.TimeZone, &get.Country, &get.IsActive, &get.CreatedAt)
		if err != nil {
			return get.EmailAddress, err
		}
		Usercontainer = append(Usercontainer, get)
		fmt.Println(Usercontainer)
	}
	if err = rows.Err(); err != nil {
		return "", err
	}
	return get.EmailAddress, nil
}
