package database

import (
	"database/sql"
	"io/ioutil"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rmacdiarmid/gptback/logger"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByEmail(t *testing.T) {
	// 1. Setup mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	DB = db

	// 2. Create test user
	testUser := User{
		UserId:       1,
		FirstName:    "John",
		LastName:     "Doe",
		Gender:       "M",
		DateOfBirth:  "2000-01-01",
		RoleId:       1,
		LoginName:    "johndoe",
		PasswordHash: "hashedpassword",
		PasswordSalt: "salt",
		Email:        "john.doe@example.com",
	}

	// 3. Test GetUserByEmail function
	email := "john.doe@example.com"
	rows := sqlmock.NewRows([]string{"UserId", "FirstName", "LastName", "Gender", "DateOfBirth", "RoleId", "LoginName", "PasswordHash", "PasswordSalt", "Email"}).AddRow(testUser.UserId, testUser.FirstName, testUser.LastName, testUser.Gender, testUser.DateOfBirth, testUser.RoleId, testUser.LoginName, testUser.PasswordHash, testUser.PasswordSalt, testUser.Email)
	mock.ExpectQuery("SELECT ua.UserId, ua.FirstName, ua.LastName, ua.Gender, ua.DateOfBirth, ua.RoleId, uld.LoginName, uld.PasswordHash, uld.PasswordSalt, uld.EmailAddress FROM user_account_6007 AS ua JOIN user_login_data_4231 AS uld ON ua.UserId = uld.UserId WHERE uld.EmailAddress = ?").WithArgs(email).WillReturnRows(rows)

	user, err := GetUserByEmail(email)
	assert.Nil(t, err)
	assert.Equal(t, testUser, user)

	// 4. Test for the case when the email is not found in the database
	email = "notfound@example.com"
	mock.ExpectQuery("SELECT ua.UserId, ua.FirstName, ua.LastName, ua.Gender, ua.DateOfBirth, ua.RoleId, uld.LoginName, uld.PasswordHash, uld.PasswordSalt, uld.EmailAddress FROM user_account_6007 AS ua JOIN user_login_data_4231 AS uld ON ua.UserId = uld.UserId WHERE uld.EmailAddress = ?").WithArgs(email).WillReturnError(sql.ErrNoRows)

	_, err = GetUserByEmail(email)
	assert.NotNil(t, err)
}

func TestCreateUser(t *testing.T) {
	// Initialize logger with a dummy logger
	logger.DualLog = log.New(ioutil.Discard, "", 0)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	DB = db

	mock.ExpectBegin()

	mock.ExpectExec("INSERT INTO user_account_6007").
		WithArgs("John", "Doe", "M", "2000-01-01", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("INSERT INTO user_login_data_4231").
		WithArgs(1, "johndoe", "hashedpassword", "salt", "john.doe@example.com").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	testUser := User{
		FirstName:    "John",
		LastName:     "Doe",
		Gender:       "M",
		DateOfBirth:  "2000-01-01",
		RoleId:       1,
		LoginName:    "johndoe",
		PasswordHash: "hashedpassword",
		PasswordSalt: "salt",
		Email:        "john.doe@example.com",
	}

	var userId int64
	var testErr error

	userId, testErr = CreateUser(testUser)
	if testErr != nil {
		t.Errorf("Expected nil, but got: %v", testErr)
	}

	if userId != 1 {
		t.Errorf("Expected user ID to be 1, but got: %d", userId)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}
