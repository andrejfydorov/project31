package repository

import (
	"fmt"
	"log"
	"project31/internal/user"
)

func (r *Repository) GetUsers() ([]*user.User, error) {

	fmt.Printf("select * from users")

	var users []*user.User

	rows, err := r.database.Query(fmt.Sprintf("select * from users"))
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		u := user.User{}
		err := rows.Scan(&u.Id, &u.Name, &u.Age)
		if err != nil {
			log.Println(err)
			continue
		}

		users = append(users, &u)
	}

	for _, _user := range users {
		rows, err = r.database.Query(fmt.Sprintf(
			"select u2.id, u2.name, u2.age"+
				" from users u"+
				" join friends f on f.userid=u.id"+
				" join users u2 on u2.id=f.friendid"+
				" where u.id=%d", _user.Id))
		if err != nil {
			log.Println(err)
		}
		defer rows.Close()

		for rows.Next() {
			u := user.User{}
			err := rows.Scan(&u.Id, &u.Name, &u.Age)
			if err != nil {
				fmt.Println(err)
				continue
			}
			_user.Friends = append(_user.Friends, &u)
		}
	}

	return users, nil
}
