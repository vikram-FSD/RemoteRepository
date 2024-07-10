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
		user    model.User
		Message = make(map[string]string)
	)
	w.Header().Set("Content-Type", "application/json")
	inputvalidator.IsMethodValid(w, r, "POST")
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return
	}
	//validating Inputs using IsStringValid Function
	isValid, _ := inputvalidator.IsStringValid("en", user.FirstName, 50, true, "firstname")
	if isValid != "valid" {
		http.Error(w, isValid, http.StatusBadRequest)
		return
	}
	isValid, _ = inputvalidator.IsStringValid("en", user.LastName, 50, true, "lastname")
	if isValid != "valid" {
		http.Error(w, isValid, http.StatusBadRequest)
		return
	}

	isValid, _ = inputvalidator.IsStringValid("en", user.Signinthrough, 50, true, "signinthrough")
	if isValid != "valid" {
		http.Error(w, isValid, http.StatusBadRequest)
		return
	}
	isValid, _ = inputvalidator.IsStringValid("en", user.TimeZone, 50, true, "Timezone")
	if isValid != "valid" {
		http.Error(w, isValid, http.StatusBadRequest)
		return
	}
	isEmailValid := inputvalidator.IsEmailAddressValid("English", user.EmailAddress, 100)
	if isEmailValid != "valid" {
		http.Error(w, isEmailValid, http.StatusBadRequest)
		return
	}
	//Autogenerating Id with length of 12 character & only contains letters & numbers and send to model to store with DB.
	user.Id = inputvalidator.GenerateRandomKey("generateID12345")
	//validating country string with space
	res, _ := inputvalidator.IsStringWitSpaceValid("EN", user.Country, 50, false, "country")
	if res != "valid" {
		http.Error(w, res, http.StatusBadRequest)
		return
	}
	createdTimeAndDate := config.CurrentDateTime(user.TimeZone)
	user.CreatedAt = createdTimeAndDate
	user.IsActive = true
	//To check if the data is already exist in DB with emailAddress, Id, and firstname
	email, id, fname, err := model.IsUserExists(user)
	if err != nil {
		log.Println(err)
	}
	if user.EmailAddress == email || user.Id == id || user.FirstName == fname {
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
