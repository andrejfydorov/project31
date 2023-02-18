package controllers

import (
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"strconv"
)

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
