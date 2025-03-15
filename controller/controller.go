package controller

import (
	"embed"
	"github.com/redis/go-redis/v9"
)

type Flag struct {
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

type Controller struct {
	rdb        *redis.Client
	templateFS embed.FS
}

func New(rdb *redis.Client, templateFS embed.FS) *Controller {
	return &Controller{
		rdb:        rdb,
		templateFS: templateFS,
	}
}
