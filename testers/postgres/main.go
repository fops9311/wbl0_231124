package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "demo1"
)

func main() {
	fmt.Println(insertOrder("hello", "world"))
}
func insertOrder(key string, data string) error {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return err
	}

	defer db.Close()

	insertDynStmt := `insert into "orders1"("id", "data") values($1, $2)`
	_, err = db.Exec(insertDynStmt, key, data)
	if err != nil {
		return err
	}
	return nil
}
