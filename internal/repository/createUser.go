package repository

import (
	"fmt"
	"project31/internal/user"
)

func (r *Repository) CreateUser(u *user.User) (int64, error) {

	fmt.Printf("insert into users (name, age) values ('%s', %d)", u.Name, u.Age)

	res, err := r.database.Exec(fmt.Sprintf("insert into users (name, age) values ('%s', %d)", u.Name, u.Age))
	if err != nil {
		return -1, err
	} else {
		return res.LastInsertId()
	}
}
