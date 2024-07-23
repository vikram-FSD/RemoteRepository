package model

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/atomedgesoft/calendariq/config"
	"github.com/atomedgesoft/inputvalidator"
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

func GetUserDetailsById(userId string) (User, error) {
	var users User
	db, _ := config.ConnectDB()
	defer db.Close()
	sql := "select * from users where id=$1"
	data, err := db.Query(sql, userId)
	if err != nil {
		log.Println(err)
		return User{}, err
	}

	for data.Next() {
		err := data.Scan(&users.Id, &users.FirstName, &users.LastName, &users.EmailAddress, &users.Signinthrough, &users.TimeZone, &users.Country, &users.IsActive, &users.CreatedAt)
		if err != nil {
			return User{}, err
		}
		return users, nil
	}
	return User{}, err
}

// UPDATE USER
func UpdateUser(get User) (string, error) {
	var w http.ResponseWriter
	var id string
	db, _ := config.ConnectDB()
	defer db.Close()
	IsIdExist, err := IsIdExist(get.Id)
	if err != nil {
		inputvalidator.WriteJson("Id not Matching with the table", w)
	}
	if IsIdExist {
		sql := "insert into users(id, firstname, lastname, emailaddress, signinthrough, timezone, country, isactive, createdat) values($1, $2, $3, $4, $5, $6, $7, $8, $9) on conflict(id) do update set firstname=Excluded.firstname,lastname=Excluded.lastname,emailaddress=Excluded.emailaddress,signinthrough=Excluded.signinthrough,createdat=Excluded.createdat,timezone=Excluded.timezone,country=Excluded.country,isactive=Excluded.isactive RETURNING id"
		err := db.QueryRow(sql, &get.Id, &get.FirstName, &get.LastName, &get.EmailAddress, &get.Signinthrough, &get.TimeZone, &get.Country, &get.IsActive, &get.CreatedAt).Scan(&id)
		if err != nil {
			return "", err
		}
	}
	return id, nil
}

// To check if the Id is existing in the table or not
func IsIdExist(id string) (bool, error) {
	db, err := config.ConnectDB()
	if err != nil {
		return false, err
	}
	defer db.Close()
	var userId string
	sqlStmt := `select id from users where id=$1`
	err = db.QueryRow(sqlStmt, id).Scan(&userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, err
		}
		return false, err
	}
	return true, nil
}
func DeleteUserDetailsById(id string) (string, error) {
	var w http.ResponseWriter
	db, err := config.ConnectDB()
	if err != nil {
		inputvalidator.ErrorHandler(err, 500, w)
	}
	defer db.Close()
	sql := "delete from users where id=$1"
	db.QueryRow(sql, id)
	return id, nil
}
