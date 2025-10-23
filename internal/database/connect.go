package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type PSQL struct {
	User string
	Password string
	Host string
	Port string
	Database string
}

type DatabaseService struct {
	DB *sql.DB
}

func NewDatabaseService(db *sql.DB) *DatabaseService {
	return &DatabaseService{DB: db}
}

func NewPSQL() (*PSQL, error) {
    _ = godotenv.Load()

    return &PSQL{
        User:     os.Getenv("POSTGRES_USER"),
        Password: os.Getenv("POSTGRES_PASSWORD"),
        Host:     os.Getenv("POSTGRES_HOST"),
        Port:     os.Getenv("POSTGRES_PORT"),
        Database: os.Getenv("POSTGRES_DB"),
    }, nil
}

func TestPSQL() (*PSQL, error) {
    err := godotenv.Load("../.env.test")
	if err != nil {
		return nil, err
	}

    return &PSQL{
        User:     os.Getenv("POSTGRES_USER"),
        Password: os.Getenv("POSTGRES_PASSWORD"),
        Host:     os.Getenv("POSTGRES_HOST"),
        Port:     os.Getenv("POSTGRES_PORT"),
        Database: os.Getenv("POSTGRES_DB"),
    }, nil
}

func (p *PSQL) Connect() (*sql.DB, error) {
	connStr := fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s?sslmode=disable",
        p.User, p.Password, p.Host, p.Port, p.Database,
    )

	db, err := sql.Open("postgres",connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, err
}