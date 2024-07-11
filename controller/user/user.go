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

func InsertUser(w http.ResponseWriter, r *http.Request) {
	var (
		user model.User
	)
	w.Header().Set("Content-Type", "application/json")
	inputvalidator.IsMethodValid(w, r, "POST")
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return
	}

	// validating Inputs
	response := InputValidation(user)
	if response["output"] != "valid" {
		inputvalidator.WriteJson(response, w)
		return
	}
	user.Id = inputvalidator.GenerateRandomKey(config.Charset)
	user.CreatedAt = config.CurrentDateTime(user.TimeZone)
	user.IsActive = true
	results, err := model.InsertUser(user)
	isValid, _ = inputvalidator.IsStringValid(config.Lang, user.LastName, 50, true, "lastname")
	if isValid != "valid" {
		http.Error(w, isValid, http.StatusBadRequest)
		return
	}

	isValid, _ = inputvalidator.IsStringValid(config.Lang, user.Signinthrough, 50, true, "signinthrough")
	if isValid != "valid" {
		http.Error(w, isValid, http.StatusBadRequest)
		return
	}
	if isValid != "valid" {
		http.Error(w, isValid, http.StatusBadRequest)
		return
	}
	isEmailValid := inputvalidator.IsEmailAddressValid(config.Lang, user.EmailAddress, 100)
	if isEmailValid != "valid" {
		http.Error(w, isEmailValid, http.StatusBadRequest)
		return
	}

	const Charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	user.Id = inputvalidator.GenerateRandomKey(Charset)
	//validating country string with space
	res, _ := inputvalidator.IsStringWitSpaceValid(config.Lang, user.Country, 50, false, "country")
	if res != "valid" {
		http.Error(w, res, http.StatusBadRequest)
		return
	}

	//To check if the data is already exist in DB with emailAddress
	email, err := model.IsUserExists(user)
	if err != nil {
		log.Println(err)
	}
	if user.EmailAddress == email {
		http.Error(w, "User Data already Exist !", http.StatusBadRequest)
		return
	} else {
		response, err := model.InsertUser(user)
		if err != nil {
			log.Print(err)
			return
		}
		Message["message"] = "New User Inserted: " + response
		output, err := json.Marshal(Message)
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		w.Write(output)
	}
}

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
