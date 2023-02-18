package repository

import (
	"fmt"
	"log"
	"project31/internal/user"
)

func (r *Repository) GetUser(id int64) *user.User {

	fmt.Printf("select id, name, age from users where id=%d", id)

	row := r.database.QueryRow(fmt.Sprintf("select id, name, age from users where id=%d", id))
	u := user.User{}
	err := row.Scan(&u.Id, &u.Name, &u.Age)
	if err != nil {
		log.Println(err)
	}
	return &u
}
