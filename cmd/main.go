package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"project31/pkg/user"
	"strconv"
)

var database *sql.DB

type Service struct {
	Store map[int]*user.User
}

func GetMaxId() int {
	row := database.QueryRow("select max(id) from users")
	var i int
	err := row.Scan(&i)
	if err != nil {
		fmt.Println(err)
	}
	return i
}

func main() {
	db, err := sql.Open("mysql", "root:12345678@/gousers")
	if err != nil {
		log.Println(err)
	}
	database = db
	defer db.Close()

	r := chi.NewRouter()
	srv := Service{map[int]*user.User{}}
	r.Post("/create", srv.Create)
	r.Post("/make_friends", srv.MakeFriends)
	r.Post("/user", srv.Delete)
	r.Get("/friends/{id:[0-9]+}", srv.GetFriends)
	r.Post("/{id:[0-9]+}", srv.Update)
	r.Get("/get", srv.GetAll)

	err = http.ListenAndServe("localhost:8080", r)
	if err != nil {
		log.Fatalln(err)
	}

}

func (s *Service) Add(u *user.User) {
	s.Store[GetMaxId()+1] = u
}

func (s *Service) Create(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
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

		_, err = database.Exec("insert into users (id, name, age) values (?, ?, ?)", GetMaxId()+1, u.Name, u.Age)
		if err != nil {
			log.Println(err)
		}
		currentUserID := GetMaxId()
		for _, friend := range u.Friends {
			i := GetMaxId() + 1
			_, err = database.Exec("insert into users (id, name, age) values (?, ?, ?)", i, friend.Name, friend.Age)
			if err != nil {
				log.Println(err)
			}
			_, err = database.Exec("insert into friends (userid, friendid) values (?, ?)", currentUserID, i)
			if err != nil {
				log.Println(err)
			}
		}

		writer.WriteHeader(http.StatusCreated)
		//fmt.Printf("User was created %s and id is %d\n", u.Name, currentUserID)
		writer.Write([]byte(fmt.Sprintf("User was created %s and id is %d\n", u.Name, currentUserID)))
		return
	}
	writer.WriteHeader(http.StatusBadRequest)
}

func (s *Service) MakeFriends(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
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

		_, err = database.Exec("insert into friends (userid, friendid) values (?, ?)", Friends.SourceId, Friends.TargetId)
		if err != nil {
			log.Println(err)
		}

		row := database.QueryRow(fmt.Sprintf("select * from users where id=%d", Friends.SourceId))

		u1 := user.User{}
		err = row.Scan(&u1.Id, &u1.Name, &u1.Age)
		if err != nil {
			fmt.Println(err)
		}

		row = database.QueryRow(fmt.Sprintf("select * from users where id=%d", Friends.TargetId))
		u2 := user.User{}
		err = row.Scan(&u2.Id, &u2.Name, &u2.Age)
		if err != nil {
			fmt.Println(err)
		}

		writer.WriteHeader(http.StatusOK)
		fmt.Printf("%s и %s теперь друзья\n", u1.Name, u2.Name)
		writer.Write([]byte(fmt.Sprintf("%s и %s теперь друзья\n", u1.Name, u2.Name)))
		return

		//writer.WriteHeader(http.StatusBadRequest)
		//writer.Write([]byte(fmt.Sprintf("User %d или %d не найден.\n", Friends.SourceId, Friends.TargetId)))
		//fmt.Printf("User %d или %d не найден.\n", Friends.SourceId, Friends.TargetId)
		//return

	}
	writer.WriteHeader(http.StatusBadRequest)
}

func (s *Service) GetFriends(writer http.ResponseWriter, request *http.Request) {
	responce := ""

	id := chi.URLParam(request, "id")
	fmt.Println(id + "\n")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Fatalln(err)
	}

	rows, err := database.Query(fmt.Sprintf(
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

func (s *Service) Delete(writer http.ResponseWriter, request *http.Request) {

	if request.Method == "POST" {
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

		row := database.QueryRow(fmt.Sprintf(fmt.Sprintf("select id, name, age from users where id=%d", UserId.TargetId)))

		u := user.User{}
		err = row.Scan(&u.Id, &u.Name, &u.Age)
		if err != nil {
			fmt.Println(err)
		}

		_, err = database.Exec(fmt.Sprintf("delete from friends where userid=%d", UserId.TargetId))
		if err != nil {
			log.Println(err)
		}

		_, err = database.Exec(fmt.Sprintf("delete from users where id=%d", UserId.TargetId))
		if err != nil {
			log.Println(err)
		}

		writer.WriteHeader(http.StatusOK)
		//fmt.Printf("User was deleted %d == %d\n", ID, userId.Target_id)
		writer.Write([]byte(fmt.Sprintf("User was deleted %s\n", u.Name)))
		return

		//writer.WriteHeader(http.StatusBadRequest)
		//writer.Write([]byte(fmt.Sprintf("User %d не найден.\n", UserId.TargetId)))
		//return

	}
	writer.WriteHeader(http.StatusBadRequest)
}

func (s *Service) Update(writer http.ResponseWriter, request *http.Request) {
	content, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		return
	}
	defer request.Body.Close()

	id := chi.URLParam(request, "id")
	fmt.Println(id + "\n")
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
	_, err = database.Exec(fmt.Sprintf("update users set age=%d where id=%d", Age.Age, idInt))
	if err != nil {
		log.Println(err)
	}

	row := database.QueryRow(fmt.Sprintf(fmt.Sprintf("select id, name, age from users where id=%d", idInt)))

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

func (s *Service) GetAll(writer http.ResponseWriter, request *http.Request) {
	s.Store = map[int]*user.User{}
	responce := ""

	rows, err := database.Query("select * from users")
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

		s.Store[u.Id] = &u
	}

	for _, _user := range s.Store {
		rows, err = database.Query(fmt.Sprintf(
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

	for _, user := range s.Store {
		responce += user.ToString() + "\n"
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(responce))
	return
}
