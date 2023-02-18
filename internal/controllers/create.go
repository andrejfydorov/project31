package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"project31/internal/user"
)

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
