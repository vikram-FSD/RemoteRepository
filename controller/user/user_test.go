package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestInsertUser(t *testing.T) {
	user := map[string]interface{}{
		"firstname":     "John",
		"lastname":      "Doe",
		"emailaddress":  "john.doe@example.com",
		"signinthrough": "email",
		"createdat":     "2024-08-23T07:57:28Z",
		"timezone":      "UTC",
		"country":       "USA",
		"isactive":      true,
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Failed to marshal the user data %v", err)
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error Initializing the Mock DB: %v", err)
	}
	defer db.Close()
	//Expected query and response for mock database
	mock.ExpectQuery(`insert into "user"`).WithArgs(sqlmock.AnyArg(), "John", "Doe", "john.doe@example.com", "email", sqlmock.AnyArg(), "UTC", "USA", true).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("12345"))
	req, err := http.NewRequest("POST", "/user", bytes.NewBuffer(userJSON))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		InsertUser(db, w, r)
	})
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %v , but get %v", http.StatusOK, status)
	}
	expectedResponse := `"New user added"`
	actualResponse := rr.Body.String()
	assert.Equal(t, expectedResponse, actualResponse, "Expected Response %v, but got %v", expectedResponse, actualResponse)

}

func TestGetUser(t *testing.T) {

	req, err := http.NewRequest("GET", "/user", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetUser)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	expected := `
	[
    {
        "id": "123",
        "firstname": "John",
        "lastname": "Doe",
        "emailaddress": "john.doe@example.com",
        "signinthrough": "email",
        "createdat": "2024-08-23T00:00:00Z",
        "timezone": "UTC",
        "isactive": true,
        "country": "USA"
    },
    {
        "id": "1001",
        "firstname": "John",
        "lastname": "Doe",
        "emailaddress": "john.doe@example.com",
        "signinthrough": "email",
        "createdat": "2024-08-23T00:00:00Z",
        "timezone": "UTC",
        "isactive": true,
        "country": "USA"
    }
]`
	assert.JSONEq(t, expected, rr.Body.String())
}
func TestGetUserById(t *testing.T) {
	req, err := http.NewRequest("GET", "/user/{id}", nil)
	if err != nil {
		t.Error(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetUser)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}
