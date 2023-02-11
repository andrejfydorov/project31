package repository

import (
	"database/sql"
	"fmt"
	"log"
)

type Repository struct {
	database *sql.DB
	MaxId    int
}

func New() *Repository {
	r := Repository{}
	db, err := sql.Open("mysql", "root:12345678@/gousers")
	if err != nil {
		log.Println(err)
	}
	r.database = db
	r.GetId()
	return &r
}

func (r *Repository) Close() {
	r.database.Close()
}

func (r *Repository) GetId() {
	row := r.database.QueryRow("SELECT LAST_INSERT_ID()")
	var i int
	err := row.Scan(&i)
	if err != nil {
		log.Println(err)
	}
	r.MaxId = i
}

func (r *Repository) QueryRow(requestMessage string) *sql.Row {
	fmt.Println(requestMessage)
	row := r.database.QueryRow(requestMessage)
	r.GetId()
	return row
}

func (r *Repository) Exec(requestMessage string) error {
	fmt.Println(requestMessage)
	_, err := r.database.Exec(requestMessage)
	if err != nil {
		return err
	} else {
		r.GetId()
		return nil
	}
}

func (r *Repository) Query(requestMessage string) (*sql.Rows, error) {
	fmt.Println(requestMessage)
	rows, err := r.database.Query(requestMessage)
	r.GetId()
	return rows, err
}
