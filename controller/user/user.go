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
	Message string `json:message`
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

	//To check if the data is already exist in DB with emailAddress
	email, err := model.IsEmailExists(user)
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
			out := Result{inputvalidator.Message(config.Lang, "insert-success")}
			Output, _ = json.Marshal(out)
			fmt.Println(out)
		}
		w.Write(Output)

		////////////////////////////////
		// 	response, err := model.InsertUser(user)
		// 	if err != nil {
		// 		log.Print(err)
		// 		return
		// 	}
		// 	Message := make(map[string]string)
		// 	Message["message"] = "New User Inserted: " + response
		// 	output, err := json.Marshal(Message)
		// 	if err != nil {
		// 		w.Write([]byte(err.Error()))
		// 	}
		// 	w.Write(output)
		// }
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
