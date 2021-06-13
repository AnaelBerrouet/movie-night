package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	//Application Packages
	"github.com/AnaelBerrouet/movie-night/handlers"

	//Third Party Packages
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Specification struct {
	Debug              bool   `default:"true"`
	Port               string `default:"4006"`
	DBConnectionString string `default:"user:user-pw@/db?multiStatements=true&parseTime=true"`
	Hostname           string `default:"localhost"`
	// HTTPSCertFilePath     string `default:"./certs/selfsigned.pem"`
	// HTTPSKeyFilePath      string `default:"./certs/selfsigned_key.pem"`
}

const APP_VERSION string = "0.0.1"

// Routing - Gorilla Mux https://www.gorillatoolkit.org/pkg/mux
// Tokenization - JWT https://github.com/dgrijalva/jwt-go
// Database migrations - https://github.com/pressly/goose
// ORM - https://gorm.io/
// Validations - https://github.com/go-playground/validator
// Env variables management - https://github.com/joho/godotenv

func main() {
	log.Println("Application Booting")

	s := fetchENV()

	db := createDBClient(s.DBConnectionString)
	// Close database connection once main function completes
	sqldb, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqldb.Close()
	// Check DB connectivity

	pingDB(sqldb)
	// Create Application Router
	r := mux.NewRouter()
	r.Handle("/", handlers.HomeHandler(db)).Methods("POST")
	// r.HandleFunc("/articles/{category}/", ArticlesCategoryHandler)
	// r.HandleFunc("/articles/{category}/{id:[0-9]+}", ArticleHandler)

	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Start server with tls on:", s.Port)
	log.Fatal(srv.ListenAndServe())
	// log.Fatal(srv.ListenAndServeTLS(s.HTTPSCertFilePath, s.HTTPSKeyFilePath))
}

// fetchENV - Fetch the ENV vars to boot the application
func fetchENV() *Specification {

	var s Specification
	err := envconfig.Process("collections", &s)
	if err != nil {
		log.Println("Error: Building Spec")
		log.Fatal(err)
	}

	return &s
}

// createClient - Create a Database Client
func createDBClient(connStr string) (db *gorm.DB) {

	db, err := gorm.Open(mysql.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	return
}

// pingDB - Ping DB to check connection on Boot
func pingDB(db *sql.DB) {
	err := db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully Pinged the DB")
}
