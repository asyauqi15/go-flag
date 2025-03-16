package controller

import (
	"embed"
	"github.com/redis/go-redis/v9"
)

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
