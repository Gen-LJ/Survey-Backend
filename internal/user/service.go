package user

import (
	"errors"
	"fmt"

	"survey-backend/internal/region"
	"survey-backend/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(u *User) error {
	// Hash the password
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Check if user already exists
	var existing User
	err = FindUserByEmail(u.Email, &existing)
	if err == nil {
		// User already exists
		return fmt.Errorf("user with email %s already exists", u.Email)
	}

	regions, err := region.FindActiveByCountryID(u.CountryID)

	if err != nil {
		return fmt.Errorf("failed to find active regions: %w", err)
	}

	found := false
	for _, r := range regions {
		if r.ID == u.RegionID {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("you can't register as this country at the moment")
	}

	if u.Role == "admin" {
		return fmt.Errorf("bro what @_@. you are a genious but i am better. xD")
	}

	validRoles := map[string]bool{
		"interviewer": true,
		"respondent":  true,
	}

	if !validRoles[u.Role] {
		return fmt.Errorf("you can't register as this role: %s", u.Role)
	}

	u.Password = string(hashed)

	return CreateUser(u)
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
