package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

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
