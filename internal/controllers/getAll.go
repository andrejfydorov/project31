package controllers

import (
	"log"
	"net/http"
)

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
