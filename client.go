package flag

import (
	"github.com/asyauqi15/go-flag/controller"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	rdb *redis.Client
}

func (c *Client) InitiateRoutes(r *chi.Mux) {
	cont := controller.New(c.rdb)

	mux := chi.NewMux()
	mux.Get("/", cont.Index)
	mux.Get("/add", cont.Add)
	mux.Post("/add", cont.AddProcess)
	mux.Get("/feature/{feature_name}", cont.Update)
	mux.Post("/feature/{feature_name}", cont.UpdateProcess)
	mux.Post("/feature/{feature_name}/delete", cont.Delete)

	r.Mount("/flag", mux)
}

func New(rdb *redis.Client) *Client {
	return &Client{
		rdb: rdb,
	}
}
