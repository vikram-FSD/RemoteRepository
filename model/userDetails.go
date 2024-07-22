package model

import (
	"fmt"
	"log"

	"github.com/atomedgesoft/calendariq/config"
)

type User struct {
	Id            *string `json:"id"`
	FirstName     *string `json:"firstname"`
	LastName      *string `json:"lastname"`
	EmailAddress  *string `json:"emailaddress"`
	Signinthrough *string `json:"signinthrough"`
	CreatedAt     *string `json:"createdat"`
	TimeZone      *string `json:"timezone"`
	IsActive      *bool   `json:"isactive"`
	Country       *string `json:"country"`
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
	return *user.Id, nil
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

func GetUserDetailsById(userId string) (User, bool, error) {
	var users User
	db, _ := config.ConnectDB()
	defer db.Close()
	sql := "select * from users where id=$1"
	data, err := db.Query(sql, userId)

	for data.Next() {
		err := data.Scan(&users.Id, &users.FirstName, &users.LastName, &users.EmailAddress, &users.Signinthrough, &users.TimeZone, &users.Country, &users.IsActive, &users.CreatedAt)
		if err != nil {
			return User{}, false, err
		}
		return users, true, err
	}
	return User{}, false, err
}

// UPDATE USER
func UpdateUser(get User) (string, error) {
	var id string
	db, _ := config.ConnectDB()
	defer db.Close()
	sql := "insert into users(id, firstname, lastname, emailaddress, signinthrough, createdat, timezone, country, isactive) values($1, $2, $3, $4, $5, $6, $7, $8, $9) on conflict(id) do update set firstname=Excluded.firstname,lastname=Excluded.lastname,emailaddress=Excluded.emailaddress,signinthrough=Excluded.signinthrough,createdat=Excluded.createdat,timezone=Excluded.timezone,country=Excluded.country,isactive=Excluded.isactive RETURNING id"
	err := db.QueryRow(sql, &get.Id, &get.FirstName, &get.LastName, &get.EmailAddress, &get.Signinthrough, &get.CreatedAt, &get.TimeZone, &get.Country, &get.IsActive).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func CreateUserDetails(user *User, Id string) (string, error) {

	db, _ := config.ConnectDB()
	defer db.Close()
	sqlStmt, sqldata := SqlStmtAndDataForUser(*user, Id)
	var id string
	err := db.QueryRow(sqlStmt, sqldata...).Scan(&id)
	return id, err
}
func SqlStmtAndDataForUser(user User, Id string) (string, []interface{}) {
	sql := "insert into users (id, firstname, lastname, emailaddress, signinthrough, timezone, country, isactive, createdat)"
	var args []interface{}
	args = append(args, Id, user.Id)
	if user.FirstName != nil {
		sql += ", firstname"
		args = append(args, user.FirstName)
	}
	if user.LastName != nil {
		sql += ", lastname"
		args = append(args, user.LastName)
	}
	if user.EmailAddress != nil {
		sql += ", emailaddress"
		args = append(args, user.EmailAddress)
	}
	if user.Signinthrough != nil {
		sql += ", signinthrough"
		args = append(args, user.Signinthrough)
	}
	if user.TimeZone != nil {
		sql += ", timezone"
		args = append(args, user.TimeZone)
	}
	if user.Country != nil {
		sql += ", country"
		args = append(args, user.Country)
	}
	if user.IsActive != nil {
		sql += ", firstname"
		args = append(args, user.IsActive)
	}
	if user.CreatedAt != nil {
		sql += ", createdat"
		args = append(args, user.CreatedAt)
	}
	sql += ")VALUES ($1, $2 "
	for i := range args[2:] {
		sql += ", $" + fmt.Sprint(i+3)
	}
	sql += ") RETURNING id"
	return sql, args
}
