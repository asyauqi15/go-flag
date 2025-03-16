package controller

import (
	"embed"
	"github.com/redis/go-redis/v9"
)

type Controller struct {
	rdb        *redis.Client
	keyPrefix  string
	templateFS embed.FS
	rootPath   string
}

type templateData struct {
	RootPath string
	Data     any
	ErrorMsg string
}

func New(rdb *redis.Client, keyPrefix string, templateFS embed.FS, rootPath string) *Controller {
	return &Controller{
		rdb:        rdb,
		keyPrefix:  keyPrefix,
		templateFS: templateFS,
		rootPath:   rootPath,
	}
}
