package controller

import "github.com/redis/go-redis/v9"

type Flag struct {
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

type Controller struct {
	rdb *redis.Client
}

func New(rdb *redis.Client) *Controller {
	return &Controller{
		rdb: rdb,
	}
}
