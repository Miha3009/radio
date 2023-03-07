package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var conn *sql.DB

func OpenConnection(config Config) error {
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", config.User, config.Password, config.Host, config.Port, config.DBname)
	newConn, err := sql.Open("postgres", connString)
	if err != nil {
		return err
	}

	conn = newConn
	return nil
}

func GetConnection() *sql.DB {
	return conn
}

func CloseConnection() error {
	if conn != nil {
		return conn.Close()
	}
	return nil
}
