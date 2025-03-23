package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ParsaAminpour/robix/backend/handler"
	"github.com/ParsaAminpour/robix/backend/models"
	_ "github.com/ParsaAminpour/robix/backend/models"
	"github.com/labstack/echo/v4"

	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openTestDatabaseOnMEM() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Problem in Opening Database on MEM")
	}
	return db, nil
}

func closeTestDatabaseOnMEM(db *gorm.DB) error {
	sqkDB, _ := db.DB()
	err := sqkDB.Close()
	if err != nil {
		log.Fatal("Problem in Closing Database on MEM")
	}
	return nil
}

type TmpUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func WriteMockDataToDatabase(db *gorm.DB, t *testing.T) error {
	file, err := os.Open("MOCK_DATA.json")
	if err != nil {
		return err
	}
	defer file.Close()

	var users []models.User
	decoder := json.NewDecoder(file)
	decoder.Decode(&users)
	for _, user := range users {
		database := models.Database{DB: db}
		database.CreateUser(&user)
	}
	return nil
}

func sendRequest(method, endpoint string, _endpoint func(c echo.Context, db *gorm.DB) error, req_body map[string]interface{}, db *gorm.DB, t *testing.T) *httptest.ResponseRecorder {
	e := echo.New()
	var req *http.Request
	var rec *httptest.ResponseRecorder

	switch method {
	case http.MethodPost:
		jsonData, err := json.Marshal(req_body)
		assert.NoError(t, err)
		req = httptest.NewRequest(http.MethodPost, endpoint, bytes.NewReader(jsonData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	case http.MethodGet:
		req = httptest.NewRequest(http.MethodGet, endpoint, nil)

	case http.MethodPut:
		jsonData, err := json.Marshal(req_body)
		assert.NoError(t, err)
		req = httptest.NewRequest(http.MethodPut, endpoint, bytes.NewReader(jsonData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	case http.MethodDelete:
		jsonData, err := json.Marshal(req_body)
		assert.NoError(t, err)
		req = httptest.NewRequest(http.MethodDelete, endpoint, bytes.NewReader(jsonData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}
	rec = httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	err := _endpoint(ctx, db)
	assert.NoError(t, err)

	return rec
}

func TestUserOperationsWithPreAddedUsers(t *testing.T) {
	db, err := openTestDatabaseOnMEM()
	assert.NoError(t, err)
	defer closeTestDatabaseOnMEM(db)

	err = WriteMockDataToDatabase(db, t)
	assert.NoError(t, err)

	t.Run("CheckTmpDataAdded", func(t *testing.T) {
		var allUsers []models.User
		database := models.Database{DB: db}
		database.GetAllUsers(&allUsers)

		for _, user := range allUsers {
			assert.NotNil(t, user)
		}
	})

	t.Run("CheckNewUserAddedViaSignup", func(t *testing.T) {
		new_user := map[string]interface{}{
			"username": "json",
			"email":    "json@gmail.com",
			"password": "testPassword",
		}

		databaes := models.Database{DB: db}
		var users_len_before int64
		var users_len_after int64
		databaes.GetUsersLength(&users_len_before)

		rec := sendRequest(http.MethodPost, "/users/signup", handler.Signup, new_user, db, t)
		var bodyResponse models.User
		err = json.Unmarshal(rec.Body.Bytes(), &bodyResponse)
		assert.NoError(t, err)

		assert.Equal(t, bodyResponse.Username, new_user["username"])
		assert.Equal(t, bodyResponse.Email, new_user["email"])
		assert.NotEqual(t, bodyResponse.Password, new_user["password"])
		databaes.GetUsersLength(&users_len_after)
		assert.Equal(t, users_len_after-1, users_len_before)
	})

	t.Run("CheckUpdateUser", func(t *testing.T) {
		plain_data := map[string]interface{}{
			"username":    "Martschke",
			"newUsername": "NewMartschke",
			"newEmail":    "gmartschke0@rakuten.co.jp",
		}
		rec := sendRequest(http.MethodPut, "/users/update", handler.UpdateUser, plain_data, db, t)

		var body_response models.User
		err = json.Unmarshal(rec.Body.Bytes(), &body_response)
		assert.NoError(t, err)

		assert.Equal(t, body_response.Username, plain_data["newUsername"])
		assert.Equal(t, body_response.Email, plain_data["newEmail"])
	})
}
