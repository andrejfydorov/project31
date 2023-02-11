package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"io/ioutil"
	"log"
	"net/http"
	"project31/pkg/repository"
	"project31/pkg/user"
	"strconv"
)

type Controllers struct {
	database *repository.Repository
}

func New(_database *repository.Repository) *Controllers {
	c := Controllers{}
	c.database = _database
	return &c
}

func (c *Controllers) getMaxId() int {
	row := c.database.QueryRow("select max(id) from users")
	var i int
	err := row.Scan(&i)
	if err != nil {
		fmt.Println(err)
	}
	return i
}

func (c *Controllers) Create(writer http.ResponseWriter, request *http.Request) {
	content, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		return
	}
	defer request.Body.Close()

	var u user.User
	if err := json.Unmarshal(content, &u); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		fmt.Println(err)
		return
	}

	err = c.database.Exec(fmt.Sprintf("insert into users (id, name, age) values (%d, '%s', %d)", c.getMaxId()+1, u.Name, u.Age))
	if err != nil {
		log.Println(err)
	}
	currentUserID := c.getMaxId()
	for _, friend := range u.Friends {
		i := c.getMaxId() + 1
		err = c.database.Exec(fmt.Sprintf("insert into users (id, name, age) values (%d, '%s', %d)", i, friend.Name, friend.Age))
		if err != nil {
			log.Println(err)
		}
		err = c.database.Exec(fmt.Sprintf("insert into friends (userid, friendid) values (%d, %d)", currentUserID, i))
		if err != nil {
			log.Println(err)
		}
	}

	writer.WriteHeader(http.StatusCreated)
	writer.Write([]byte(fmt.Sprintf("User was created %s and id is %d\n", u.Name, currentUserID)))
	return

}

func (c *Controllers) MakeFriends(writer http.ResponseWriter, request *http.Request) {
	content, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		return
	}
	defer request.Body.Close()

	var Friends struct {
		SourceId int `json:"source_id"`
		TargetId int `json:"target_id"`
	}
	if err := json.Unmarshal(content, &Friends); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		return
	}

	err = c.database.Exec(fmt.Sprintf("insert into friends (userid, friendid) values (%d, %d)", Friends.SourceId, Friends.TargetId))
	if err != nil {
		log.Println(err)
	}

	row := c.database.QueryRow(fmt.Sprintf("select * from users where id=%d", Friends.SourceId))

	u1 := user.User{}
	err = row.Scan(&u1.Id, &u1.Name, &u1.Age)
	if err != nil {
		fmt.Println(err)
	}

	row = c.database.QueryRow(fmt.Sprintf("select * from users where id=%d", Friends.TargetId))
	u2 := user.User{}
	err = row.Scan(&u2.Id, &u2.Name, &u2.Age)
	if err != nil {
		fmt.Println(err)
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(fmt.Sprintf("%s и %s теперь друзья\n", u1.Name, u2.Name)))
	return

	//writer.WriteHeader(http.StatusBadRequest)
	//writer.Write([]byte(fmt.Sprintf("User %d или %d не найден.\n", Friends.SourceId, Friends.TargetId)))
	//fmt.Printf("User %d или %d не найден.\n", Friends.SourceId, Friends.TargetId)
	//return
}

func (c *Controllers) GetFriends(writer http.ResponseWriter, request *http.Request) {
	responce := ""

	id := chi.URLParam(request, "id")
	fmt.Println(id + "\n")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Fatalln(err)
	}

	rows, err := c.database.Query(fmt.Sprintf(
		"select u2.id, u2.name, u2.age"+
			" from users u"+
			" join friends f on f.userid=u.id"+
			" join users u2 on u2.id=f.friendid"+
			" where u.id=%d", idInt))
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
		responce += u.ToStringShort()
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(responce))
	return

	//writer.WriteHeader(http.StatusBadRequest)
	//writer.Write([]byte(fmt.Sprintf("User %d не найден.\n", idInt)))
	//return

}

func (c *Controllers) Delete(writer http.ResponseWriter, request *http.Request) {
	content, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		return
	}
	defer request.Body.Close()

	var UserId struct {
		TargetId int `json:"target_id"`
	}

	if err := json.Unmarshal(content, &UserId); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		//fmt.Println(err)
		return
	}

	row := c.database.QueryRow(fmt.Sprintf(fmt.Sprintf("select id, name, age from users where id=%d", UserId.TargetId)))

	u := user.User{}
	err = row.Scan(&u.Id, &u.Name, &u.Age)
	if err != nil {
		fmt.Println(err)
	}

	err = c.database.Exec(fmt.Sprintf("delete from friends where userid=%d", UserId.TargetId))
	if err != nil {
		log.Println(err)
	}

	err = c.database.Exec(fmt.Sprintf("delete from users where id=%d", UserId.TargetId))
	if err != nil {
		log.Println(err)
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(fmt.Sprintf("User was deleted %s\n", u.Name)))
	return

	//writer.WriteHeader(http.StatusBadRequest)
	//writer.Write([]byte(fmt.Sprintf("User %d не найден.\n", UserId.TargetId)))
	//return
}

func (c *Controllers) Update(writer http.ResponseWriter, request *http.Request) {
	content, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		return
	}
	defer request.Body.Close()

	id := chi.URLParam(request, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Fatalln(err)
	}

	var Age struct {
		Age int `json:"new age"`
	}
	if err := json.Unmarshal(content, &Age); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		fmt.Println(err)
		return
	}
	err = c.database.Exec(fmt.Sprintf("update users set age=%d where id=%d", Age.Age, idInt))
	if err != nil {
		log.Println(err)
	}

	row := c.database.QueryRow(fmt.Sprintf(fmt.Sprintf("select id, name, age from users where id=%d", idInt)))

	u := user.User{}
	err = row.Scan(&u.Id, &u.Name, &u.Age)
	if err != nil {
		fmt.Println(err)
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(fmt.Sprintf("Возраст пользователя %s успешно обновлен.\n", u.Name)))
	return

	//writer.WriteHeader(http.StatusBadRequest)
	//writer.Write([]byte(fmt.Sprintf("User %d не найден.\n", idInt)))
	//return

}

func (c *Controllers) GetAll(writer http.ResponseWriter, request *http.Request) {
	var store []*user.User
	responce := ""

	rows, err := c.database.Query("select * from users")
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

		store = append(store, &u)
	}

	for _, _user := range store {
		rows, err = c.database.Query(fmt.Sprintf(
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

	for _, user := range store {
		responce += user.ToString() + "\n"
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(responce))
	return
}
