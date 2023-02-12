package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"io/ioutil"
	"log"
	"net/http"
	"project31/internal/repository"
	"project31/internal/user"
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

	id, err := c.database.CreateUser(&u)
	if err != nil {
		log.Println(err)
	}

	if id != -1 {
		currentUser := c.database.GetUser(id)
		for _, friend := range u.Friends {
			id, err = c.database.CreateUser(friend)
			if err != nil {
				log.Println(err)
			}
			_, err := c.database.CreateFriends(currentUser, friend)
			if err != nil {
				log.Println(err)
			}
		}

		writer.WriteHeader(http.StatusCreated)
		writer.Write([]byte(fmt.Sprintf("User was created %s and id is %d\n", currentUser.Name, currentUser.Id)))
		return
	}

	writer.WriteHeader(http.StatusInternalServerError)
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
		SourceId int64 `json:"source_id"`
		TargetId int64 `json:"target_id"`
	}
	if err := json.Unmarshal(content, &Friends); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		return
	}

	u1 := c.database.GetUser(Friends.SourceId)
	u2 := c.database.GetUser(Friends.TargetId)

	i, err := c.database.CreateFriends(u1, u2)
	if err != nil {
		log.Println(err)
	}
	if i > 0 {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(fmt.Sprintf("%s и %s теперь друзья\n", u1.Name, u2.Name)))
		return
	}

	writer.WriteHeader(http.StatusBadRequest)
	return
}

func (c *Controllers) GetFriends(writer http.ResponseWriter, request *http.Request) {
	responce := ""

	id := chi.URLParam(request, "id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Fatalln(err)
	}

	users, err := c.database.GetFriends(int64(idInt))
	if err != nil {
		log.Println(err)
	}
	if len(users) > 0 {
		for _, u := range users {
			responce += u.ToStringShort()
		}

		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(responce))
		return
	}

	writer.WriteHeader(http.StatusBadRequest)
	return
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
		TargetId int64 `json:"target_id"`
	}

	if err := json.Unmarshal(content, &UserId); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		return
	}

	user := c.database.GetUser(UserId.TargetId)

	i, err := c.database.DeleteUser(UserId.TargetId)
	if err != nil {
		log.Println(err)
	}

	if i > 0 {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(fmt.Sprintf("User was deleted %s\n", user.Name)))
		return
	}

	writer.WriteHeader(http.StatusBadRequest)
	return
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
	i, err := c.database.UpdateUser(idInt, Age.Age)
	if err != nil {
		log.Println(err)
	}

	if i > 0 {
		user := c.database.GetUser(int64(idInt))

		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(fmt.Sprintf("Возраст пользователя %s успешно обновлен.\n", user.Name)))
		return
	}

	writer.WriteHeader(http.StatusBadRequest)
	return
}

func (c *Controllers) GetAll(writer http.ResponseWriter, request *http.Request) {

	responce := ""

	users, err := c.database.GetUsers()
	if err != nil {
		log.Println(err)
	}

	if len(users) > 0 {
		for _, user := range users {
			responce += user.ToString() + "\n"
		}

		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(responce))
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("No users"))
	return
}
