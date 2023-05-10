package internal

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/rmacdiarmid/gptback/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

func IsEmailUnique(email string) (bool, error) {
	_, err := database.GetUserByEmail(email)
	if err == nil {
		return false, fmt.Errorf("email already used")
	}
	return true, nil
}

func RegisterUser(input map[string]interface{}) (int64, error) {
	// Check if the email is unique
	email := input["email"].(string)
	unique, err := IsEmailUnique(email)
	if !unique {
		return 0, err
	}

	// Compare the provided passwords
	password := input["password"].(string)
	passwordConfirmation := input["passwordConfirmation"].(string)
	if password != passwordConfirmation {
		return 0, fmt.Errorf("passwords do not match")
	}

	// Hash the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("error hashing password: %v", err)
	}

	// Create a user struct with the provided input and hashed password
	user := database.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	// Insert the user data into the database
	userID, err := database.CreateUser(user)
	if err != nil {
		return 0, fmt.Errorf("error creating user: %v", err)
	}

	return userID, nil
}

func LoginUser(input map[string]interface{}) (string, error) {
	email := input["email"].(string)
	password := input["password"].(string)

	// Retrieve user data from the database
	// Assuming you have a function `database.GetUserByEmail` that takes the email and returns the user data
	user, err := database.GetUserByEmail(email)
	if err != nil {
		return "", fmt.Errorf("error retrieving user: %v", err)
	}

	// Compare the provided password with the stored password hash
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid password")
	}

	// Create a JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.UserId,
		"email":  user.Email,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign and return the token
	// Assuming you have an environment variable `JWT_SECRET` with your JWT secret
	jwtSecret := os.Getenv("JWT_SECRET")
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("error signing token: %v", err)
	}

	return signedToken, nil
}
