package repository

import (
	"fmt"
	"project31/internal/user"
)

func (r *Repository) CreateFriends(u1 *user.User, u2 *user.User) (int64, error) {

	fmt.Printf("insert into friends (userid, friendid) values (%d, %d)", u1.Id, u2.Id)

	res, err := r.database.Exec(fmt.Sprintf("insert into friends (userid, friendid) values (%d, %d)", u1.Id, u2.Id))
	if err != nil {
		return -1, err
	} else {
		return res.RowsAffected()
	}
}
