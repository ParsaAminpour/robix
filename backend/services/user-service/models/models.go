package models

import (
	"errors"
	"fmt"

	"github.com/ParsaAminpour/robix/backend/utils"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique" json:"username"`
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"password"`
}

type Database struct {
	DB *gorm.DB
}

func (db *Database) FetchUser(user_ref *User, byUsername string) error {
	res := db.DB.Where("username = ?", byUsername).First(user_ref)
	if res.RowsAffected == 0 {
		return fmt.Errorf("error in fetching %s", byUsername)
	}
	return nil
}

func (db *Database) CreateUser(user_ref *User) error {
	if err := db.DB.Create(&user_ref).Error; err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == "23505" {
			return fmt.Errorf("username or email already exists")
		}
		return fmt.Errorf("could not create user")
	}
	return nil
}

func (db *Database) DeleteUser(user_ref *User, byUsername string) error {
	if err := db.DB.Where("username = ?", byUsername).First(user_ref).Error; err != nil {
		return errors.New("user is not found to delete")
	}
	if err := db.DB.Delete(user_ref).Error; err != nil {
		return err
	}
	return nil
}

func (db *Database) GetAllUsers(users_ref *[]User) error {
	if err := db.DB.Find(&users_ref).Error; err != nil {
		return err
	}
	return nil
}

func (db *Database) GetUsersLength(len *int64) error {
	if err := db.DB.Model(&User{}).Count(len).Error; err != nil {
		return err
	}
	return nil
}

func (db *Database) UpdateUser(user_ref *User, newUsername, newEmail string) error {
	user_ref.Username = newUsername
	if newEmail != "" {
		user_ref.Email = newEmail
	}
	db.DB.Save(&user_ref)
	return nil
}

func (db *Database) UpdateUserPassword(user_ref *User, username, old_password, new_password string) error {
	if res := db.DB.Model(&User{}).Where("username = ?", username).First(user_ref); res.RowsAffected == 0 {
		return fmt.Errorf("user is not found to update")
	}
	hashed_password, _ := utils.HashPassword(old_password)
	if user_ref.Password != hashed_password {
		return fmt.Errorf("old password is incorrect")
	}
	hashed_new_password, _ := utils.HashPassword(new_password)
	user_ref.Password = hashed_new_password
	db.DB.Save(&user_ref)

	return nil
}
