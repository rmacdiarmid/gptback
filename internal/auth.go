package internal

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/rmacdiarmid/gptback/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(input map[string]interface{}) (int64, error) {
	// Input validation and user creation
	// Hash the password using bcrypt
	password := input["password"].(string)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("error hashing password: %v", err)
	}

	// Create a user struct with the provided input and hashed password
	user := database.User{
		Email:        input["email"].(string),
		PasswordHash: string(hashedPassword),
		// Add any additional fields if necessary
	}

	if firstName, ok := input["firstName"]; ok {
		user.FirstName = firstName.(string)
	}

	if lastName, ok := input["lastName"]; ok {
		user.LastName = lastName.(string)
	}

	if gender, ok := input["gender"]; ok {
		user.Gender = gender.(string)
	}

	if dateOfBirth, ok := input["dateOfBirth"]; ok {
		user.DateOfBirth = dateOfBirth.(string)
	}

	// Insert the user data into the database
	// Assuming you have a function `database.CreateUser` that takes the user struct and returns the user ID
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
