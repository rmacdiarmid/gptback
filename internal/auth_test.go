package internal

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"testing"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/stretchr/testify/assert"

// 	"github.com/rmacdiarmid/gptback/logger"
// 	"github.com/rmacdiarmid/gptback/pkg/database"
// )

// func init() {
// 	// Initialize logger with a dummy logger
// 	logger.DualLog = log.New(ioutil.Discard, "", 0)
// }

// func TestRegisterUser(t *testing.T) {
// 	// 1. Setup mock database
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()
// 	database.DB = db

// 	// Expect transaction Begin, Commit, and Rollback
// 	mock.ExpectBegin()
// 	mock.ExpectCommit()

// 	// configure mock to expect a call to CreateUser with the provided user
// 	user := database.User{
// 		Email:        "test@example.com",
// 		PasswordHash: "hashedpassword",
// 		FirstName:    "John",
// 		LastName:     "Doe",
// 		Gender:       "Male",
// 		DateOfBirth:  "1990-01-01",
// 	}
// 	// Add the expectation for the additional query
// 	mock.ExpectExec("INSERT INTO user_account_6007").
// 		WithArgs(user.FirstName, user.LastName, user.Gender, user.DateOfBirth, 0). // Adjust the arguments as needed
// 		WillReturnResult(sqlmock.NewResult(1, 1))

// 	// Debugging: print out all expectations before calling RegisterUser
// 	for _, e := range mock.Expectations() {
// 		fmt.Printf("Expectation: %s\n", e.String())
// 	}

// 	// call the function being tested
// 	t.Run("RegisterUser", func(t *testing.T) {
// 		_, err = RegisterUser(map[string]interface{}{
// 			"email":       user.Email,
// 			"password":    "password",
// 			"firstName":   user.FirstName,
// 			"lastName":    user.LastName,
// 			"gender":      user.Gender,
// 			"dateOfBirth": user.DateOfBirth,
// 		})
// 		assert.NoError(t, err)
// 	})

// 	// assert that all expectations were met
// 	assert.NoError(t, mock.ExpectationsWereMet())
// }

// func TestLoginUser(t *testing.T) {
// 	// create a new mock database connection
// 	db, mock, err := sqlmock.New()
// 	assert.NoError(t, err)
// 	defer db.Close()

// 	// configure mock to expect a call to GetUserByEmail with the provided email
// 	email := "test@example.com"
// 	user := database.User{
// 		Email:        email,
// 		PasswordHash: "hashedpassword",
// 	}
// 	mock.ExpectQuery("SELECT.*FROM users.*WHERE.*email.*=.*").
// 		WithArgs(email).
// 		WillReturnRows(sqlmock.NewRows([]string{"email", "password_hash"}).AddRow(user.Email, user.PasswordHash))

// 		// Debugging: print out all expectations before calling RegisterUser
// 	for _, e := range mock.Expectations() {
// 		fmt.Printf("Expectation: %s\n", e.String())
// 	}

// 	// call the function being tested
// 	_, err = LoginUser(map[string]interface{}{
// 		"email":    email,
// 		"password": "password",
// 	})
// 	assert.NoError(t, err)

// 	// assert that all expectations were met
// 	assert.NoError(t, mock.ExpectationsWereMet())
// }
