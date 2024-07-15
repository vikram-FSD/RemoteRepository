
# Detail

This Project is about to validate the user Data and perform insertion and retrieval of data from the Database

**POST** `call`:

Parameter | Type | Description
--- | --- | --- |
`firstname` | `string`| **Required.** string only contains letters & numbers.
`lastname`  | `string` |**Required.** string only contains letters & numbers.
`Email Address`|`string`|**Required.** string should be valid access
`Signin_through`|`string`|**Required.** string ( Ex: Gmail, Outlook etc., )
`CreatedAt`|`string`|**Required.** string Date Time based on the timezone.
`TimeZone `|`string`|**Required.** Only must contain String & "/"
`IsActive`|`boolean`|When you receive post request, ensure you are making IsActive is true.
`Country`| `string`|**Required.** should be in string and may have space

**Method**



InsertUser

```func InsertUser(w http.ResponseWriter, r *http.Request) {
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
	exist, err := IsEmailExists(user)
	if err != nil {
		log.Println("Error checking email existence:", err)
		inputvalidator.ErrorHandler(err, http.StatusInternalServerError, w)
		return
	}
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
		} else {
			http.Error(w, "User could not be added.", http.StatusInternalServerError)
		}
	}
}
```
***Method***

Retrieve User
```
func GetUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	w.Header().Set("content-type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		inputvalidator.ErrorHandler(err, 500, w)
		return
	}
	res, err := model.GetUser(user)
	if err != nil {
		inputvalidator.ErrorHandler(err, http.StatusInternalServerError, w)
		return
	}
	jData, _ := json.Marshal(res)
	w.Write(jData)
}
```


**Input :**
```
{
  "firstname"    :    "Ganesh",
  "lastname"     :    "P",
  "emailaddress" :    "Ganes@gmail.com",
  "signinthrough":    "Gmail",
  "timezone"     :    "Asia/Calcutta",
  "country"      :    "India"
}
```
***Sequence Diagram***


