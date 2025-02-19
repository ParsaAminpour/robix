package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/ParsaAminpour/robix/backend/handler"
	_ "github.com/ParsaAminpour/robix/backend/handler"
	"github.com/ParsaAminpour/robix/backend/models"
	_ "github.com/ParsaAminpour/robix/backend/models"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	mu   = &sync.Mutex{}
	db   *gorm.DB
	once sync.Once
)

type Config struct {
	Host     string
	Port     string
	Password string
	User     string
	DBName   string
	SSLMode  string
}

func (conf *Config) getDBConfig() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		return Config{}, fmt.Errorf("Error in loading dotenv")
	}
	db_config := Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		User:     os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}
	return db_config, nil
}

func getDB() *gorm.DB {
	var err error
	once.Do(func() {
		var dbConf Config
		dbConf, err = dbConf.getDBConfig()
		if err != nil {
			log.Fatalf("Failed to get DB config: %v", err)
		}
		db, err = initializeDB(&dbConf)
		if err != nil {
			log.Fatalf("Failed to initialize DB: %v", err)
		}
	})
	return db
}

func initializeDB(conf *Config) (*gorm.DB, error) {
	if db != nil {
		return nil, fmt.Errorf("database already initialized")
	}
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		conf.Host, conf.Port, conf.User, conf.Password, conf.DBName, conf.SSLMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error opening DSN: %w", err)
	}
	return db, nil
}

func endpointHandler(_handler func(c echo.Context, db *gorm.DB) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		return _handler(c, db)
	}
}

func main() {
	getDB()
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("error while migraing: ", err.Error())
	}
	fmt.Println("The DB initialized: ")

	e := echo.New()

	e.GET("/", endpointHandler(handler.HomePage))
	e.Group("/users")
	e.GET("users/all", endpointHandler(handler.GetAllUsers))
	e.GET("/users/:username", endpointHandler(handler.GetUser))
	e.POST("/users/signup", endpointHandler(handler.Signup))
	e.PUT("/users/update", endpointHandler(handler.UpdateUser))
	e.DELETE("/users/delete/:username", endpointHandler(handler.DeleteUser))

	e.Logger.Fatal(e.Start(":8080"))
}
