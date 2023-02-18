package repository

import (
	"fmt"
	"log"
)

func (r *Repository) DeleteUser(id int64) (int64, error) {

	fmt.Printf("delete from friends where userid=%d", id)

	res, err := r.database.Exec(fmt.Sprintf("delete from friends where userid=%d", id))
	if err != nil {
		log.Println(err)
	}

	fmt.Printf("delete from users where id=%d", id)

	res, err = r.database.Exec(fmt.Sprintf("delete from users where id=%d", id))
	if err != nil {
		return -1, err
	}

	return res.RowsAffected()
}
