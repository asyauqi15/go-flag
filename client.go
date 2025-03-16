package flag

import (
	"embed"
	"errors"
	"github.com/asyauqi15/go-flag/controller"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

//go:embed views/*.html
var templateFS embed.FS

type Client struct {
	rdb      *redis.Client
	rootPath string
}

func (c *Client) InitiateRoutes(r *chi.Mux, path string) {
	cont := controller.New(c.rdb, templateFS, path)

	mux := chi.NewMux()
	mux.Get("/", cont.Index)
	mux.Get("/add", cont.Add)
	mux.Post("/add", cont.AddProcess)
	mux.Get("/feature/{feature_name}", cont.Update)
	mux.Post("/feature/{feature_name}", cont.UpdateProcess)
	mux.Post("/feature/{feature_name}/delete", cont.Delete)

	r.Mount(path, mux)
}

func New(rdb *redis.Client) (*Client, error) {
	if rdb == nil {
		return nil, errors.New("redis connection is nil")
	}

	return &Client{
		rdb: rdb,
	}, nil
}
