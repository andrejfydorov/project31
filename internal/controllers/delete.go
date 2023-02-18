package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

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
