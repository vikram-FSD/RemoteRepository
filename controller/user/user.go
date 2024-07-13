package user

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/atomedgesoft/calendariq/config"
	"github.com/atomedgesoft/calendariq/inputvalidator"
	"github.com/atomedgesoft/calendariq/model"
)

type Result struct {
	Result string `json:result`
}

var Output []byte

// POST USER
func InsertUser(w http.ResponseWriter, r *http.Request) {
	var (
		user model.User
	)
	w.Header().Set("Content-Type", "application/json")
	inputvalidator.IsMethodValid(w, r, "POST")
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		return
	}
	// validating Inputs
	response, errs := InputValidation(user)
	if response["output"] != "valid" {
		inputvalidator.WriteJson(response, w)
		return
	}
	if errs["output"] != "valid" {
		inputvalidator.WriteJson(errs, w)
		return
	}
	user.Id = inputvalidator.GenerateRandomKey(config.Charset)
	user.CreatedAt = config.CurrentDateTime(user.TimeZone)
	user.IsActive = true
	//To check if the data is already exist in DB with emailAddress
	email, err := IsEmailExists(user)
	if err != nil {
		log.Println(err)
		return
	}
	if user.EmailAddress == email {
		http.Error(w, "Email-ID already Exist !", http.StatusBadRequest)
		return
	} else {
		results, err := model.InsertUser(user)
		if err != nil {
			inputvalidator.ErrorHandler(err, 500, w)
		}
		if len(results) > 0 {
			out := Result{"New user added: " + user.Id}
			Output, _ = json.Marshal(out)
		}
		w.Write(Output)
	}
}

// GET USER
func GetUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	w.Header().Set("content-type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Print(err)
	}
	res, err := model.GetUser(user)
	if err != nil {
		return
	}
	jData, _ := json.Marshal(res)
	fmt.Fprintf(w, `Data Retrieved Successfully.`)
	w.Write(jData)
}

// Checking if the emailID is already exist or not
func IsEmailExists(user model.User) (email string, error error) {
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
	db.Close()
	var get model.User
	for rows.Next() {
		err := rows.Scan(&get.Id, &get.FirstName, &get.LastName, &get.EmailAddress, &get.Signinthrough, &get.TimeZone, &get.Country, &get.IsActive, &get.CreatedAt)
		if err != nil {
			return "", err
		}
	}
	if err = rows.Err(); err != nil {
		return "", err
	}
	return get.EmailAddress, nil
}
