package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

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
