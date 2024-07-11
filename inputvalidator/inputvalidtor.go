package inputvalidator

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	emailverifier "github.com/AfterShip/email-verifier"
)

/*
* Error Handling
 */
func ErrorHandler(err error, statuscode int, w http.ResponseWriter) {
	if err != nil {
		log.Default().Print(err.Error())

		w.WriteHeader(statuscode)
		w.Write([]byte(err.Error()))
		return
	}
	return

}

/*
* Validate Given string with the max length
 */
func IsStringValid(Lang string, str string, maxChar int, Isnullable bool, fieldname string) (string, error) {

	var Output string

	if !Isnullable && len(str) == 0 {
		Output = fieldname + " " + message[Lang]["must-be-valid"]
		return Output, nil
	}

	pattern := `^[A-Za-z0-9]+$`
	matched, err := regexp.MatchString(pattern, str)

	CheckErrWithReturn(err)

	if !matched {
		Output = message[Lang]["invalid-char-found"]
	} else if len(str) > maxChar {
		Output = message[Lang]["len-is-greater-than-expected"]
	} else {
		Output = "valid"
	}
	return Output, nil
}

/*
* Validate Given string with space & the max length -
 */
func IsStringWitSpaceValid(Lang string, str string, maxChar int, Isnullable bool, fieldname string) (string, error) {

	var Output string

	if !Isnullable && len(str) == 0 {
		Output = fieldname + " " + message[Lang]["must-be-valid"]
	}

	pattern := `^[A-Za-z0-9 ]+$`
	matched, err := regexp.MatchString(pattern, str)

	CheckErrWithReturn(err)

	if !matched {
		Output = message[Lang]["invalid-char-found"]
	} else if len(str) > maxChar {
		Output = message[Lang]["len-is-greater-than-expected"]
	} else {
		Output = "valid"
	}
	return Output, nil
}

/*
* Validate Email Address & Input Length
 */
func IsEmailAddressValid(Lang string, emailaddress string, maxChar int) string {

	req, err := emailverifier.NewVerifier().Verify(emailaddress)

	if strings.ContainsAny(emailaddress, "` | ' | \" | * | % | ; | :  | ,  | / | \\ | + | ") {
		return message[Lang]["invalid-char-found"]
	} else if len(emailaddress) > maxChar {
		return message[Lang]["len-is-greater-than-20-char"]
	} else if err != nil {
		log.Default().Print(err.Error())
		return message[Lang]["invalid-email-address"]
	} else if !req.Syntax.Valid {
		return message[Lang]["invalid-email-address"]
	} else if strings.ContainsAny(emailaddress, " ' ' | ''") {
		return message[Lang]["invalid-email-address"]
	}

	return "valid"

}

/*
*	Generate Key
 */

func GenerateRandomKey(Charset string) string {
	// Seed the random number generator
	rand.NewSource(time.Now().UnixNano())

	var key string
	for i := 0; i < 12; i++ {
		index := rand.Intn(len(Charset))
		key += string(Charset[index])
	}
	return key

}

/*
*	Validate Contact Number & it the length should be 10 Numbers
 */
func IsContactValid(Lang string, contact string) string {

	var Output string
	_, err := strconv.Atoi(contact)

	if err != nil {
		Output = message[Lang]["invalid-contact"]
	} else if len(contact) > 10 || len(contact) < 5 {
		Output = message[Lang]["invalid-contact"]
	} else if strings.ContainsAny(contact, "` | ' | \" | * | % | ; | :  | ,  | / | \\ | ' '") {
		Output = message[Lang]["invalid-contact"]
	} else {
		Output = "valid"
	}

	return Output
}

/*
* Validate Given Company Name with the max length
 */
func IsCompanyNameValid(Lang string, CompanyName string, maxChar int) (string, error) {

	var Output string
	pattern := `^[A-Za-z0-9 ]+$`
	matched, err := regexp.MatchString(pattern, CompanyName)

	CheckErrWithReturn(err)

	switch {
	case CompanyName == "":
		Output = message[Lang]["company-name-invalid"]
	case !matched:
		Output = message[Lang]["invalid-char-found"]
	case len(CompanyName) > maxChar || len(CompanyName) < 6:
		Output = message[Lang]["len-is-greater-than-expected"]
	default:
		Output = "valid"
	}

	return Output, nil
}

/*
* Validate Given string with the max length
 */
func IsPasswordValid(Lang string, str string, minChar int) (string, error) {

	var Output string
	pattern := `^[A-Za-z0-9]+$`
	matched, err := regexp.MatchString(pattern, str)

	CheckErrWithReturn(err)

	if len(str) == 0 {
		Output = message[Lang]["enter-valid-password"]
	} else if strings.Contains(str, " ") {
		Output = message[Lang]["should-not-contain-space"]
	} else if strings.IndexFunc(str, unicode.IsLower) < 0 {
		Output = message[Lang]["must-contain-lower-char"]
	} else if matched == true {
		Output = message[Lang]["must-contain-special-char"]
	} else if len(str) < minChar {
		Output = message[Lang]["len-is-less-than-expected"]
	} else {
		Output = "valid"
	}
	return Output, nil
}

/*
*	Validate URL Input
 */
func IsURLValid(Lang string, site string) string {
	var err error = nil

	u, err := url.Parse(site)

	if err != nil || u.Scheme == "" || u.Host == "" {
		return message[Lang]["invalid-url"]
	}
	return "valid"
}

/*
* Note: Key will be generated internally & the length should 5
* Reason: Ensure they shouldnt pass any special characters to harm our system.
 */
func IsKeyValid(Lang string, Key string) string {

	if len(Key) != 14 {
		return message[Lang]["invalid-length-key"]
	} else if strings.ContainsAny(Key, "` | ' | \" | * | % | ; | :  | , | . | / | \\ | ") {
		return message[Lang]["invalid-char-found"]
	} else {
		return "valid"
	}
}

func Message(Lang string, key string) string {
	return message[Lang][key]
}

// Return current Date & time
func Timenow(Timenow *time.Location) time.Time {
	return time.Now().In(Timenow)
}

// func GetDate(date string)  {
// 	layout := "2006-01-02" // Adjust layout if needed
// 	dob, err := time.Parse(layout, date)
// 	ErrorHandler(err)
// 	parsedDate := dob.Format("2006-01-02") // YYYY-MM-DD format
// 	return parsedDate
// }

/*
*	Calculate Age
 */
func YearsBetween(startDate, endDate time.Time) int {
	years := endDate.Year() - startDate.Year()

	// Check if the end date is before the anniversary of the start date in the end year
	if endDate.Before(startDate.AddDate(years, 0, 0)) {
		years--
	}
	return years
}

/*
* Is Age Valid
 */

func IsAgeValid(Lang string, Age int) string {
	if Age > 90 || Age < 18 {
		return message[Lang]["age-is-under-limit"]
	}
	return "valid"
}

func IsValidDOB(Lang string, DOB string) string {
	layout := "2006-01-02" // Adjust layout if needed
	_, err := time.Parse(layout, DOB)

	if err != nil {
		return message[Lang]["dob-invalid"]
	}

	return "valid"
}

/*
*	Validate Extension Number & it the length should be less than 10
 */
func IsExtensionValid(Lang string, contact int) string {

	var str string
	str = fmt.Sprintf("%d", contact)

	if len(str) > 10 || len(str) < 1 {
		return message[Lang]["invalid-contact"]
	} else {
		return "valid"
	}

}

/*
*	Validate Extension Number & it the length should be less than 10
 */
func IsContactNoValid(Lang string, contact *int) string {

	var str string

	if contact != nil {

		fmt.Println("Contact:", *contact)

		str = fmt.Sprintf("%d", *contact)

		fmt.Println("Contact after string:", str)
		fmt.Println("Length of Phone Number:", len(str))
		if len(str) > 10 || len(str) <= 4 {
			return message[Lang]["invalid-contact"]
		}
	}

	return "valid"
}

func IsAddressValid(Lang string, address string) string {

	if len(address) > 100 {
		return message[Lang]["address-length-greater-10-char"]
	}

	return "valid"
}

// It validates the Request Method and throw an error
func IsMethodValid(w http.ResponseWriter, r *http.Request, method string) {
	if r.Method != method {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}
}

// Return Json Object
func WriteJson(message any, w http.ResponseWriter) {

	Output, _ := json.Marshal(message)

	w.Write(Output)

}

func Layout() string {
	return "2006-01-02"
}

// Validate Date
func IsDateValid(Lang string, locate *time.Location, date *string, fieldname string) (string, error) {
	layout := Layout()

	if date == nil {
		return "valid", nil
	}

	parsedDate, err := time.Parse(layout, *date) // Adjust layout if needed

	if err != nil {
		return "valid", err
	}

	if parsedDate.After(Timenow(locate)) {
		return fieldname + " " + message[Lang]["future-date"], nil
	}
	return "valid", nil
}

func IsToDateIsGreaterThanFromDate(Lang string, fromdate string, todate *string) (string, error) {

	layout := Layout()

	if todate == nil {
		return "valid", nil
	}

	parsedFromDate, err := time.Parse(layout, fromdate)

	if err != nil {
		return "valid", err
	}
	parsedToDate, err := time.Parse(layout, *todate)

	if parsedFromDate.After(parsedToDate) {
		return message[Lang]["from-date-greater-than-to-date"], nil
	}

	return "valid", nil

}

func CheckErrWithReturn(err error) (string, error) {

	if err != nil {
		return "", err
	}

	return "valid", nil

}

// Check whether the input is int or return false
func IsIntValid(Lang string, input string) (bool, string) {
	var num int
	_, err := fmt.Sscanln(input, &num) // Scan the string into an int

	if err != nil {
		return false, message[Lang]["invalid-data-type"]
	}

	return true, ""

}
func IsRequired(input, data string) string {
	if data == "" {
		return fmt.Sprint(input, " is missing")
	}
	return data
}
func IsFirstNameValid(Lang string, str string, maxChar int, Isnullable bool, fieldname string) (string, error) {
	var output string
	if !Isnullable && len(str) == 0 {
		output = fieldname + " " + message[Lang]["must-be-valid"]
		return output, nil
	}
	pattern := `^[A-Za-z0-9]+$`
	matched, err := regexp.MatchString(pattern, str)
	if err != nil {
		return "", err
	}
	if !matched {
		output = message[Lang]["invalid-firstName"]
	} else if len(str) > maxChar {
		output = message[Lang]["Len-is-greater-than-expected"]

	} else {
		output = "valid"
	}
	return output, nil
}

func IsJSONValid(permission map[string]int) {

}
