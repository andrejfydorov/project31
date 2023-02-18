package repository

import (
	"database/sql"
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

//func (r *Repository) GetId() int {
//	row := r.database.QueryRow("select LAST_INSERT_ID()")
//	var i int
//	err := row.Scan(&i)
//	if err != nil {
//		log.Println(err)
//	}
//	return i
//}
