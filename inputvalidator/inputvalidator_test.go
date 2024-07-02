package inputvalidator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var StringInput = []struct {
	name string
	want string
}{
	{name: "Validfirstname", want: "valid"}, // Replace with a valid name
	{name: "Short name", want: "valid"},
	{name: "Name &8", want: "Invalid Character Found"},
	{name: "Name 8", want: "valid"},
	{name: "", want: "Invalid Character Found"}, // Empty string should be invalid
}

var EmailInput = []struct {
	name string
	want string
}{
	{name: "vijs@atomedgesoft.com", want: "valid"},
	{name: "vijs+j@atomedgesoft.com", want: message["en"]["invalid-char-found"]},
	{name: "*8988@gmail.com", want: message["en"]["invalid-char-found"]},
	{name: "askjdhaskdjasdasasdasdasdasdasdasdasdasdasdasddsadhaskdhasasdasdasdd@gmail.com", want: message["en"]["len-is-greater-than-20-char"]},
	{name: "asksa@gmailleee.coi", want: message["en"]["invalid-email-address"]},
	{name: "", want: message["en"]["invalid-email-address"]},
}

var IsUserIdExists = []struct {
	name   string
	userID string
	want   bool
}{
	{name: "Existing user", userID: "user123", want: true},
	{name: "Non-existing user", userID: "user456", want: false},
}

func TestIsStringValid(t *testing.T) {

	for _, tt := range StringInput {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := IsStringValid("en", tt.name, 100, true, "Name"); got != tt.want {
				t.Errorf("IsStringValid(%q) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestIsEmailAddressValid(t *testing.T) {

	for _, tt := range EmailInput {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmailAddressValid("en", tt.name, 50); got != tt.want {
				t.Errorf("IsEmailAddressValid(%q) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

// func TestIsEmailExists(t *testing.T) {
// 	tests := []struct {
// 		name         string
// 		emailaddress string
// 		table        string
// 		mockRows     *sqlmock.Rows
// 		expected     string
// 		expectError  error
// 	}{
// 		{
// 			name:         "Email exists",
// 			emailaddress: "vijdsssss@atomedgesoft.com",
// 			table:        "user_primary_info",
// 			mockRows:     sqlmock.NewRows([]string{"Id"}).AddRow(1),
// 			expected:     message[config.Lang]["email-exists"],
// 			expectError:  nil,
// 		},
// 		{
// 			name:         "Email does not exist",
// 			emailaddress: "test@example.com",
// 			table:        "user_primary_info",
// 			mockRows:     sqlmock.NewRows([]string{"Id"}),
// 			expected:     "valid",
// 			expectError:  nil,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {

// 			result, err := IsEmailExists(tt.emailaddress, tt.table)
// 			assert.Equal(t, tt.expected, result)
// 			assert.Equal(t, tt.expectError, err)
// 		})
// 	}
// }

func TestIsContactValid(t *testing.T) {

	var ContactInput = []struct {
		name string
		want string
	}{
		{name: "8971133050", want: "valid"},
		{name: "sdfsdfsdf", want: message["en"]["invalid-contact"]},
		{name: "*8988@gmail.com", want: message["en"]["invalid-contact"]},
		{name: "", want: message["en"]["invalid-contact"]},
	}

	for _, tt := range ContactInput {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsContactValid("en", tt.name); got != tt.want {
				t.Errorf("IsContactValid(%q) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestIsCompanyNameValid(t *testing.T) {

	var CompanyNameInput = []struct {
		name        string
		expectError error
		want        string
	}{
		{name: "32323", expectError: nil, want: message["en"]["len-is-greater-than-expected"]},
		{name: "sdfsdfsdf", expectError: nil, want: "valid"},
		{name: "*8988@gmail.com", expectError: nil, want: message["en"]["invalid-char-found"]},
		{name: "", expectError: nil, want: message["en"]["company-name-invalid"]},
	}

	for _, tt := range CompanyNameInput {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := IsCompanyNameValid("en", tt.name, 100); got != tt.want {
				t.Errorf("IsCompanyNameValid(%q) = %v, want %v", tt.name, got, tt.want)
				assert.Equal(t, tt.expectError, err)
			}
		})
	}
}

func TestIsPasswordValid(t *testing.T) {

	var PasswordInput = []struct {
		name        string
		expectError error
		want        string
	}{
		{name: "32323", expectError: nil, want: message["en"]["must-contain-special-char"]},
		{name: "32df3f8", expectError: nil, want: message["en"]["must-contain-special-char"]},
		{name: "sdfsdfsdf^&asd", expectError: nil, want: "valid"},
		{name: "8988gmailcom", expectError: nil, want: message["en"]["must-contain-special-char"]},
		{name: " ", expectError: nil, want: message["en"]["should-not-contain-space"]},
	}

	for _, tt := range PasswordInput {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := IsPasswordValid("en", tt.name, 6); got != tt.want {
				t.Errorf("IsPasswordValid(%q) = %v, want %v", tt.name, got, tt.want)
				assert.Equal(t, tt.expectError, err)
			}
		})
	}
}

func TestIsURLValid(t *testing.T) {

	var URLInput = []struct {
		name string
		want string
	}{
		{name: "https://atomedgesoft.com", want: "valid"},
		{name: " ", want: message["en"]["invalid-url"]},
		{name: "atomedgesoft", want: message["en"]["invalid-url"]},
	}

	for _, tt := range URLInput {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsURLValid("en", tt.name); got != tt.want {
				t.Errorf("IsURLValid(%q) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
