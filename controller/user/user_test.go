package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/atomedgesoft/calendariq/model"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	fmt.Println("Hello")
	os.Exit(m.Run())
}
func TestInsertUser(t *testing.T) {
	// Create a new mock database connection and the associated mock object

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	// Define the expected user input
	user := model.User{
		Id:            "12345",
		FirstName:     "Vikram",
		LastName:      "R",
		EmailAddress:  "Vikram2@gmail.com",
		Signinthrough: "Gmail",
		CreatedAt:     "2024-08-23T14:30:00Z",
		TimeZone:      "UTC",
		Country:       "INDIA",
		IsActive:      true,
	}
	mock.ExpectQuery(`insert into "user"(id, firstname, lastname, emailaddress, signinthrough, createdat,  timezone, country,isactive)`).
		WithArgs(user.Id, user.FirstName, user.LastName, user.EmailAddress, user.Signinthrough, user.CreatedAt, user.TimeZone, user.Country, user.IsActive)
	reqBody, _ := json.Marshal(user)
	req, err := http.NewRequest("POST", "/user", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		InsertUser(w, r)
	})
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code, "Expected status code 200")
	Output, _ := json.Marshal("New user added")

	expectedResponse := Output
	assert.EqualValues(t, expectedResponse, rr.Body.String(), "Response body differs")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
