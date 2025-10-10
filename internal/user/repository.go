package user

import (
	"survey-backend/pkg/database"
)

func CreateUser(u *User) error {
	return  database.DB.Create(u).Error
}

func FindUserByEmail(email string,u *User) error {
	return database.DB.Where("email = ?",email).First(u).Error
}

func FindUserById(id uint,u *User) error {
	return  database.DB.First(u,id).Error
}

func UpdateUser(u *User) error {
	return  database.DB.Save(u).Error
}

func UpdateUserField(id uint,fields map[string]any) error {
	return  database.DB.Model(&User{}).Where("id = ?").Updates(fields).Error
}

func DeleteUser(id uint) error {
	return  database.DB.Delete(&User{},id).Error
}

func GetAllUsers(users *[]User) error {
	return  database.DB.Find(users).Error
}

func SearchUsers(keyword string,users *[]User) error {
	query := "%" + keyword + "%"
	return database.DB.Where("name LIKE ? OR email LIKE ?",query,query).Find(users).Error
}
