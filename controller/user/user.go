package user

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/atomedgesoft/calendariq/config"
	"github.com/atomedgesoft/calendariq/model"
	"github.com/atomedgesoft/inputvalidator"
)

type Result struct {
	Result string `json:"result"`
}

// POST USER
func InsertUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	w.Header().Set("Content-Type", "application/json")
	inputvalidator.IsMethodValid(w, r, "POST")
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println("Error decoding user data:", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	// Validating Inputs
	response, errs := InputValidation(user)
	if response["output"] != "valid" {
		inputvalidator.WriteJson(response, w)
		return
	}
	if len(errs) > 0 {
		inputvalidator.WriteJson(errs, w)
		return
	}
	user.Id = inputvalidator.GenerateRandomKey(config.Charset)
	user.CreatedAt = config.CurrentDateTime(user.TimeZone)
	user.IsActive = true

	// Check if the email is already in use
	exist, _ := IsEmailExists(user.EmailAddress)
	if exist {
		http.Error(w, "Email-ID already exists!", http.StatusBadRequest)
		return
	}
	if !exist {
		resultID, err := model.InsertUser(user)
		if err != nil {
			log.Println("Error inserting user:", err)
			inputvalidator.ErrorHandler(err, http.StatusInternalServerError, w)
			return
		}
		if resultID != "" {
			out := Result{Result: "New user added: " + user.Id}
			inputvalidator.WriteJson(out, w)
			return
		} else {
			http.Error(w, "User could not be added.", http.StatusInternalServerError)
		}
	}
}

// GET USER
func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user model.User
	res, err := model.GetUser(user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	inputvalidator.WriteJson(res, w)
}

// Checking if the emailID is already exist or not
func IsEmailExists(email string) (bool, error) {
	db, err := config.ConnectDB()
	if err != nil {
		return false, err
	}
	defer db.Close()
	var emailaddress string
	sqlStmt := `select emailaddress from users where emailaddress=$1`
	err = db.QueryRow(sqlStmt, email).Scan(&emailaddress)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, err
		}
		return false, err
	}
	return true, nil
}
