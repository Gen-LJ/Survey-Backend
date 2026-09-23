package user

import (
	"survey-backend/pkg/database"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreateUser(u *User) error {
	return database.DB.Create(u).Error
}

func FindUserByEmail(email string, u *User) error {
	return database.DB.Where("email = ?", email).First(u).Error
}

func FindUserById(id uint, u *User) error {
	return database.DB.First(u, id).Error
}

// FindUserByIdForUpdate loads a user with a row lock so a point balance can be
// read and written without another request interleaving.
func FindUserByIdForUpdate(tx *gorm.DB, id uint, u *User) error {
	return tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(u, id).Error
}

func UpdateUser(u *User) error {
	return database.DB.Save(u).Error
}

func UpdateUserField(id uint, fields map[string]any) error {
	return database.DB.Model(&User{}).Where("id = ?", id).Updates(fields).Error
}

// UpdatePoints writes a user's point balance inside the caller's transaction.
func UpdatePoints(tx *gorm.DB, id uint, points uint) error {
	return tx.Model(&User{}).Where("id = ?", id).Update("points", points).Error
}

// AddPoints credits a user inside the caller's transaction.
func AddPoints(tx *gorm.DB, id uint, points uint) error {
	return tx.Model(&User{}).
		Where("id = ?", id).
		Update("points", gorm.Expr("points + ?", points)).Error
}

func DeleteUser(id uint) error {
	return database.DB.Delete(&User{}, id).Error
}

func GetAllUsers(users *[]User) error {
	return database.DB.Find(users).Error
}

func SearchUsers(keyword string, users *[]User) error {
	query := "%" + keyword + "%"
	return database.DB.Where("name LIKE ? OR email LIKE ?", query, query).Find(users).Error
}
