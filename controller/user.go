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
		log.Print(err)
	}
	//validating Empty Input
	if user.FirstName == "" {
		http.Error(w, "FirstName is required", http.StatusBadRequest)
		return
	}
	if user.LastName == "" {
		http.Error(w, "LastName is required", http.StatusBadRequest)
		return
	}
	if user.EmailAddress == "" {
		http.Error(w, "EmailId is required", http.StatusBadRequest)
		return
	}
	if user.Signinthrough == "" {
		http.Error(w, "SignInThrough is required", http.StatusBadRequest)
		return
	}
	if user.TimeZone == "" {
		http.Error(w, "Timezone is required", http.StatusBadRequest)
		return
	}
	//Validating emailAddress and send to model to store with DB
	email := user.EmailAddress
	isEmailValid := inputvalidator.IsEmailAddressValid("English", email, 100)
	if isEmailValid == "valid" {
		user.EmailAddress = email
	} else {
		http.Error(w, "Email ID is invalid", http.StatusBadRequest)
		log.Fatal("Email ID is Invalid Please enter the valid one!")
	}
	//Autogenerating Id with length of 12 character & only contains letters & numbers and send to model to store with DB.
	user.Id = inputvalidator.GenerateRandomKey("generateID12345")
	//validating country string with space
	Country := user.Country
	res, err := inputvalidator.IsStringWitSpaceValid("EN", Country, 50, false, "country")
	if res == "valid" {
		user.Country = Country
	} else {
		log.Fatal(err)
	}
	//Validating the firstname
	fName := user.FirstName
	str, _ := inputvalidator.IsFirstNameValid("en", fName, 60, false, "firstname")
	if str == "valid" {

		user.FirstName = fName

	} else {
		http.Error(w, "FirstName invalid", http.StatusBadRequest)
		return
	}
	createdTimeAndDate := config.CurrentDateTime(user.TimeZone)
	user.CreatedAt = createdTimeAndDate
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(user.CreatedAt)
	}
	user.IsActive = true
	response, err := model.InsertUser(user)
	if err != nil {
		log.Fatal(err)
	} else {
		Message["message"] = "New User Inserted: " + response
		fmt.Println(Message)
		output, err := json.Marshal(Message)
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Write(output)
		}
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	w.Header().Set("content-type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Print(err)
	}
	res, err := model.ReturnUser(user)
	if err != nil {
		fmt.Println(err)
	} else {
		jData, _ := json.Marshal(res)
		w.Write(jData)
	}
	fmt.Fprintf(w, `Data Retrieved Successfully.`)
}
