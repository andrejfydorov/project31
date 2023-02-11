package repository

import (
	"database/sql"
	"fmt"
	"log"
)

type Repository struct {
	database *sql.DB
}

func New() *Repository {
	r := Repository{}
	db, err := sql.Open("mysql", "root:12345678@/gousers")
	if err != nil {
		log.Println(err)
	}
	r.database = db
	return &r
}

func (r *Repository) Close() {
	r.database.Close()
}

func (r *Repository) QueryRow(requestMessage string) *sql.Row {
	row := r.database.QueryRow(requestMessage)
	return row
}

func (r *Repository) Exec(requestMessage string) error {
	fmt.Println(requestMessage)
	_, err := r.database.Exec(requestMessage)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (r *Repository) Query(requestMessage string) (*sql.Rows, error) {
	rows, err := r.database.Query(requestMessage)
	return rows, err
}
