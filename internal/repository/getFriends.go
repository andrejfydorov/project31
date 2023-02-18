package repository

import (
	"fmt"
	"log"
	"project31/internal/user"
)

func (r *Repository) GetFriends(id int64) ([]*user.User, error) {

	fmt.Printf(
		"select u2.id, u2.name, u2.age"+
			" from users u"+
			" join friends f on f.userid=u.id"+
			" join users u2 on u2.id=f.friendid"+
			" where u.id=%d", id)

	rows, err := r.database.Query(fmt.Sprintf(
		"select u2.id, u2.name, u2.age"+
			" from users u"+
			" join friends f on f.userid=u.id"+
			" join users u2 on u2.id=f.friendid"+
			" where u.id=%d", id))
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	var users []*user.User

	for rows.Next() {
		u := user.User{}
		err := rows.Scan(&u.Id, &u.Name, &u.Age)
		if err != nil {
			log.Println(err)
			continue
		}
		users = append(users, &u)
	}
	return users, nil
}
