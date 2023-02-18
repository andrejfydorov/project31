package controllers

import (
	"project31/internal/repository"
)

type Controllers struct {
	database *repository.Repository
}

func New(_database *repository.Repository) *Controllers {
	c := Controllers{}
	c.database = _database
	return &c
}
