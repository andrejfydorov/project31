package repository

import (
	"database/sql"
	"fmt"
	"log"
	"project31/internal/user"
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

func (r *Repository) CreateUser(u *user.User) (int64, error) {

	fmt.Printf("insert into users (name, age) values ('%s', %d)", u.Name, u.Age)

	res, err := r.database.Exec(fmt.Sprintf("insert into users (name, age) values ('%s', %d)", u.Name, u.Age))
	if err != nil {
		return -1, err
	} else {
		return res.LastInsertId()
	}
}

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

func (r *Repository) UpdateUser(id int, newAge int) (int64, error) {

	fmt.Printf("update users set age=%d where id=%d", newAge, id)

	res, err := r.database.Exec(fmt.Sprintf("update users set age=%d where id=%d", newAge, id))
	if err != nil {
		log.Println(err)
	}

	return res.RowsAffected()
}

func (r *Repository) CreateFriends(u1 *user.User, u2 *user.User) (int64, error) {

	fmt.Printf("insert into friends (userid, friendid) values (%d, %d)", u1.Id, u2.Id)

	res, err := r.database.Exec(fmt.Sprintf("insert into friends (userid, friendid) values (%d, %d)", u1.Id, u2.Id))
	if err != nil {
		return -1, err
	} else {
		return res.RowsAffected()
	}
}

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
