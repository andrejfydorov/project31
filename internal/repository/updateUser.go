package repository

import (
	"fmt"
	"log"
)

func (r *Repository) UpdateUser(id int, newAge int) (int64, error) {

	fmt.Printf("update users set age=%d where id=%d", newAge, id)

	res, err := r.database.Exec(fmt.Sprintf("update users set age=%d where id=%d", newAge, id))
	if err != nil {
		log.Println(err)
	}

	return res.RowsAffected()
}
