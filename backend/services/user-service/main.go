package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/ParsaAminpour/robix/backend/handler"
	_ "github.com/ParsaAminpour/robix/backend/handler"
	"github.com/ParsaAminpour/robix/backend/models"
	_ "github.com/ParsaAminpour/robix/backend/models"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus"
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

	skipper := func(c echo.Context) bool {
		return c.Request().URL.Path == "/metrics"
	}
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		Skipper:   skipper,
		BeforeNextFunc: func(c echo.Context) {
			c.Set("customValueFromContext", 42)
		},
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			value, _ := c.Get("customValueFromContext").(int)
			color.Green("REQUEST: uri: %v, status: %v, custom-value: %v\n", v.URI, v.Status, value)
			return nil
		},
	}))

	CustomCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "user_service",
			Name:      "requests_total",
			Help:      "Total number of requests processed by the MyApp web server.",
		},
		[]string{"path", "status"},
	)

	ErrorCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "user_service",
			Name:      "errors_total",
			Help:      "Total number of errors processed by the MyApp web server.",
		},
		[]string{"path", "status"},
	)

	if err := prometheus.Register(CustomCounter); err != nil {
		color.Red("error while registering custom counter: ", err.Error())
	}
	if err := prometheus.Register(ErrorCounter); err != nil {
		color.Red("error while registering error counter: ", err.Error())
	}

	e.Use(echoprometheus.NewMiddlewareWithConfig(echoprometheus.MiddlewareConfig{
		AfterNext: func(c echo.Context, err error) {
			CustomCounter.WithLabelValues(c.Request().URL.Path, strconv.Itoa(c.Response().Status)).Inc()
			if err != nil {
				ErrorCounter.WithLabelValues(c.Request().URL.Path, strconv.Itoa(c.Response().Status)).Inc()
			}
		},
	}))
	e.GET("/metrics", echoprometheus.NewHandler())

	e.GET("/", endpointHandler(handler.HomePage))
	e.Group("/users")
	e.GET("users/all", endpointHandler(handler.GetAllUsers))
	e.GET("/users/:username", endpointHandler(handler.GetUser))
	e.POST("/users/signup", endpointHandler(handler.Signup))
	e.PUT("/users/update", endpointHandler(handler.UpdateUser))
	e.DELETE("/users/delete/:username", endpointHandler(handler.DeleteUser))

	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		e.Logger.Fatal(err)
	}
}
