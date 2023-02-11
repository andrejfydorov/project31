package main

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"project31/pkg/controllers"
	"project31/pkg/repository"
	"project31/pkg/router"
)

func main() {
	repo := repository.New()
	defer repo.Close()

	c := controllers.New(repo)

	r := router.New()
	r.InitRoutes(c)

	err := http.ListenAndServe("localhost:8080", r.Router())
	if err != nil {
		log.Fatalln(err)
	}

}
