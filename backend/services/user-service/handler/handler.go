package handler

import (
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/ParsaAminpour/robix/backend/models"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	mu = &sync.Mutex{}
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func PasswordHashValidation(password, password_hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password_hash), []byte(password))
	return err == nil
}

func HomePage(c echo.Context, db *gorm.DB) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "Welcome to Home page!"})
}

func Signup(c echo.Context, db *gorm.DB) error {
	mu.Lock()
	defer mu.Unlock()

	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
	}
	user.Password = hashedPassword

	database := &models.Database{DB: db}
	if err := database.CreateUser(user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	log.Printf("Created: Username %s | Email %s\n", user.Username, user.Email)
	return c.JSON(http.StatusCreated, user)
}

func GetUser(c echo.Context, db *gorm.DB) error {
	mu.Lock()
	defer mu.Unlock()

	username := c.Param("username")
	database := &models.Database{DB: db}
	user := new(models.User)
	err := database.FetchUser(user, username)

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	if err != nil {
		return c.JSON(http.StatusNoContent, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}

func DeleteUser(c echo.Context, db *gorm.DB) error {
	mu.Lock()
	defer mu.Unlock()

	username := c.Param("username")
	user_to_delete := new(models.User)
	database := &models.Database{DB: db}
	err := database.DeleteUser(user_to_delete, username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "user deleted"})
}

func GetAllUsers(c echo.Context, db *gorm.DB) error {
	mu.Lock()
	defer mu.Unlock()

	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Record not found"})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"record": users})
}

type NewUserInformation struct {
	Username    string `json:"username"`
	NewUsername string `json:"newUsername"`
	NewEmail    string `json:"newEmail"`
}

func UpdateUser(c echo.Context, db *gorm.DB) error {
	mu.Lock()
	defer mu.Unlock()

	user := new(models.User)
	new_user := new(NewUserInformation)
	if err := c.Bind(new_user); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	database := &models.Database{DB: db}
	if err := database.DB.Where("username = ?", new_user.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Record not found!"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	user.Username = new_user.NewUsername
	if new_user.NewEmail != "" {
		user.Email = new_user.NewEmail
	}
	db.Save(&user)
	return c.JSON(http.StatusOK, user)
}
