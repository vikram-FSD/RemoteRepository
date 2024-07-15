package user

import (
	"github.com/atomedgesoft/calendariq/config"
	"github.com/atomedgesoft/calendariq/inputvalidator"
	"github.com/atomedgesoft/calendariq/model"
)

func InputValidation(user model.User) (map[string]string, map[string]string) {
	response := make(map[string]string)
	errors := make(map[string]error)
	messageArray := make(map[string]string)
	errorArray := make(map[string]string)

	response["firstname"], errors["firstname"] = inputvalidator.IsStringValid(config.Lang, user.FirstName, 50, true, "firstname")
	response["lastname"], errors["lastname"] = inputvalidator.IsStringValid(config.Lang, user.LastName, 50, true, "lastname")
	response["email"] = inputvalidator.IsEmailAddressValid(config.Lang, user.EmailAddress, 50)
	response["signinthrough"], errors["signinthrough"] = inputvalidator.IsStringValid(config.Lang, user.Signinthrough, 50, false, "signinthrough")
	response["timezone"], errors["timezone"] = inputvalidator.IsTimeZoneValid(config.Lang, user.TimeZone, 50, true, "timezone")
	response["country"], errors["country"] = inputvalidator.IsStringWitSpaceValid(config.Lang, user.Country, 50, true, "country")

	// Loop for response
	for key, value := range response {
		if value != "valid" {
			messageArray[key] = value
		}
	}

	// Map key="output" value="valid" if the length of the map is 0
	if len(messageArray) == 0 {
		messageArray["output"] = "valid"
	}

	// Process errors for errorArray
	for key, err := range errors {
		if err != nil {
			errorArray[key] = err.Error()
		}
	}

	return messageArray, errorArray
}
