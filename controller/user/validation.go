package user

import (
	"github.com/atomedgesoft/calendariq/config"
	"github.com/atomedgesoft/calendariq/inputvalidator"
	"github.com/atomedgesoft/calendariq/model"
)

var (
	response      = make(map[string]string)
	errors        = make(map[string]error)
	message_array = make(map[string]string)
)

func InputValidation(user model.User) map[string]string {
	user.Id = inputvalidator.GenerateRandomKey(config.Charset) //Autogenerating Id
	response["firstname"], errors["firstname"] = inputvalidator.IsStringValid(config.Lang, user.FirstName, 50, true, "firstname")
	response["lastname"], errors["lastname"] = inputvalidator.IsStringValid(config.Lang, user.LastName, 50, true, "lastname")
	response["email"] = inputvalidator.IsEmailAddressValid(config.Lang, user.EmailAddress, 50)
	response["signinthrough"], errors["signinthrough"] = inputvalidator.IsStringValid(config.Lang, user.Signinthrough, 50, false, "signinthrough")
	response["timezone"], errors["timezone"] = inputvalidator.IsStringValid(config.Lang, user.TimeZone, 50, true, "Timezone")
	response["country"], errors["country"] = inputvalidator.IsStringWitSpaceValid(config.Lang, user.Country, 50, true, "country")

	message_array = make(map[string]string)

	for key, value := range response {
		if value != "valid" {
			message_array[key] = value
		}
	}
	if len(message_array) == 0 {
		message_array["output"] = "valid"
	}
	return message_array

}
