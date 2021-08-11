package store

import (
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // for driver
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

//Store ...
type Store struct {
	DB *sqlx.DB
}

//DbStore ...
var DbStore *Store = NewStore()

//NewStore - create new db connection
func NewStore() *Store {
	loadEnvs()
	dbConnString := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ":3306)/" + os.Getenv("DB_NAME") + "?parseTime=true&loc=Europe%2FKiev"
	conn, err := sqlx.Connect("mysql", dbConnString)
	if err != nil {
		log.Fatal(err)
	}
	conn.SetMaxIdleConns(0)
	conn.SetMaxOpenConns(151)
	return &Store{
		DB: conn,
	}
}

//Close db connection
func (s *Store) Close() {
	s.DB.Close()
}

func loadEnvs() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}
