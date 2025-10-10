package user

import (
	"errors"
	"fmt"

	"survey-backend/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(name, email, password, role string) error {
	// Hash the password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Check if user already exists
	var existing User
	err = FindUserByEmail(email, &existing)
	if err == nil {
		// User already exists
		return fmt.Errorf("user with email %s already exists", email)
	}

	if role == "admin" {
		return fmt.Errorf("bro what @_@. you are a genious but i am better. xD")
	}

	validRoles := map[string]bool{
		"interviewer": true,
		"respondent":  true,
	}

	if !validRoles[role] {
		return fmt.Errorf("you can't register as this role: %s", role)
	}

	// Create new user
	newUser := User{
		Name:     name,
		Email:    email,
		Password: string(hashed),
		Role:     role,
	}

	return CreateUser(&newUser)
}

func LoginUser(email, password string) (User, string, error) {
	var user User
	err := FindUserByEmail(email, &user)
	if err != nil {
		return User{}, "", errors.New("invalid credentials")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return User{}, "", errors.New("invalid credentials")
	}

	token, _ := jwt.GenerateToken(email)

	return user, token, nil
}
