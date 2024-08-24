package user

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/atomedgesoft/calendariq/config"
	"github.com/atomedgesoft/calendariq/model"
	"github.com/atomedgesoft/inputvalidator"
	"github.com/gorilla/mux"
)

// POST USER
func InsertUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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
		resultID, err := model.InsertUser(db, user)
		if err != nil {
			log.Println("Error inserting user:", err)
			inputvalidator.ErrorHandler(err, http.StatusInternalServerError, w)
			return
		}
		if resultID != "" {
			out := "New user added"
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
	return
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

// * Create / Update User Details on user_details table
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var patch model.User
	w.Header().Set("Content-Type", "application/json")
	inputvalidator.IsMethodValid(w, r, "PATCH")
	err := json.NewDecoder(r.Body).Decode(&patch)
	if err != nil {
		inputvalidator.ErrorHandler(err, 500, w)
		return
	}
	validData, _ := InputValidation(patch)
	if validData["output"] == "valid" {
		updatedId, _ := model.UpdateUser(patch)
		if updatedId == "" {
			inputvalidator.WriteJson("No records found ", w)
			return
		}
		result := make(map[string]string)
		result["User Updated Successfulluy: "] = updatedId
		inputvalidator.WriteJson(result, w)
		return
	}
	inputvalidator.WriteJson(validData, w)
}

// Get User by ID function
func GetUserById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var Result model.User
	inputs := mux.Vars(r) //mux.Vars(r)    which get the request Input id sent by the url in post man
	Result, err := model.GetUserDetailsById(inputs["id"])
	if err != nil {
		inputvalidator.ErrorHandler(err, 500, w)
		return
	}
	if len(Result.Id) == 0 {
		message := make(map[string]string)
		message["output"] = "No record Found with this ID !"
		inputvalidator.WriteJson(message, w)
		return
	}
	inputvalidator.WriteJson(Result, w)
}
func DeleteUserById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	inputs := mux.Vars(r)
	idFound, err := model.GetUserDetailsById(inputs["id"])
	if idFound.Id != "" {
		_, err := model.DeleteUserDetailsById(inputs["id"])
		if err != nil {
			inputvalidator.ErrorHandler(err, 500, w)
		}
		out := "Deleted Successfully"
		inputvalidator.WriteJson(out, w)
	}
	inputvalidator.ErrorHandler(err, 500, w)

}
