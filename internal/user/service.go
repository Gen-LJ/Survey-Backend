package user

import (
	"errors"
	"fmt"
	"log"

	"survey-backend/internal/country"
	"survey-backend/internal/region"
	"survey-backend/internal/user/response"
	"survey-backend/pkg/config"
	"survey-backend/pkg/database"
	"survey-backend/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RegisterUser(u *User) error {
	// Check if user already exists. Anything other than "not found" is a real
	// database problem and must not be read as "the email is free".
	var existing User
	err := FindUserByEmail(u.Email, &existing)
	if err == nil {
		return fmt.Errorf("user with email %s already exists", u.Email)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to look up user: %w", err)
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

	// Hash the password only once the request is known to be valid.
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
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

	token, err := jwt.GenerateToken(email)
	if err != nil {
		return User{}, "", fmt.Errorf("failed to issue token: %w", err)
	}

	return user, token, nil
}

func GetRegisterForm() ([]response.CountryWithRegions, error) {
	countries, err := country.FindActive()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch active countries")
	}

	if len(countries) == 0 {
		return nil, fmt.Errorf("no active country found")
	}

	var result []response.CountryWithRegions

	for _, c := range countries {
		regions, err := region.FindActiveByCountryID(c.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch regions for country %s:%w", c.Name, err)
		}

		var regionList []response.RegionBrief

		for _, r := range regions {
			regionList = append(regionList, response.RegionBrief{
				ID:   r.ID,
				Name: r.Name,
				Code: r.Code,
			})
		}

		result = append(result, response.CountryWithRegions{
			ID:      c.ID,
			Name:    c.Name,
			Code:    c.Code,
			Regions: regionList,
		})

	}

	return result, nil
}

// MinAdminPasswordLength is the shortest password accepted for the seeded admin.
const MinAdminPasswordLength = 8

var ErrUserNotFound = errors.New("user not found")

// SeedAdmin creates the platform administrator from ADMIN_EMAIL and
// ADMIN_PASSWORD when that account does not exist yet.
//
// RegisterUser deliberately refuses the admin role, so without this there is no
// way to obtain an admin short of editing the database by hand - which also
// means the admin-only endpoints are unreachable on a fresh deploy.
//
// It never touches an account that already exists, so leaving the variables set
// across restarts is safe and re-running it will not reset a changed password.
func SeedAdmin() error {
	email := config.Get("ADMIN_EMAIL", "")
	password := config.Get("ADMIN_PASSWORD", "")

	if email == "" || password == "" {
		return nil
	}

	if len(password) < MinAdminPasswordLength {
		return fmt.Errorf("ADMIN_PASSWORD must be at least %d characters", MinAdminPasswordLength)
	}

	var existing User
	err := FindUserByEmail(email, &existing)
	if err == nil {
		if existing.Role != "admin" {
			log.Printf("warning: %s already exists with role %q; not changing it", email, existing.Role)
		}
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to look up the admin account: %w", err)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash the admin password: %w", err)
	}

	admin := User{
		Name:     config.Get("ADMIN_NAME", "Administrator"),
		Email:    email,
		Password: string(hashed),
		Role:     "admin",
	}

	if err := CreateUser(&admin); err != nil {
		return fmt.Errorf("failed to create the admin account: %w", err)
	}

	log.Println("created admin account:", email)
	return nil
}

// GrantPoints credits a user's balance. Points otherwise only enter the system
// when a survey is published, so without this an interviewer on a fresh deploy
// can never afford to publish anything.
func GrantPoints(email string, points uint) (User, error) {
	var updated User

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var target User
		if err := tx.Where("email = ?", email).First(&target).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrUserNotFound
			}
			return err
		}

		if err := AddPoints(tx, target.ID, points); err != nil {
			return err
		}

		return tx.First(&updated, target.ID).Error
	})
	if err != nil {
		return User{}, err
	}

	return updated, nil
}
