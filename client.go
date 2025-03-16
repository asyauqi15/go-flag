package flag

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"github.com/asyauqi15/go-flag/controller"
	"github.com/asyauqi15/go-flag/model"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

//go:embed views/*.html
var templateFS embed.FS

type Client struct {
	rdb *redis.Client
}

func (c *Client) InitiateRoutes(r *chi.Mux) {
	cont := controller.New(c.rdb, templateFS)

	mux := chi.NewMux()
	mux.Get("/", cont.Index)
	mux.Get("/add", cont.Add)
	mux.Post("/add", cont.AddProcess)
	mux.Get("/feature/{feature_name}", cont.Update)
	mux.Post("/feature/{feature_name}", cont.UpdateProcess)
	mux.Post("/feature/{feature_name}/delete", cont.Delete)

	r.Mount("/flag", mux)
}

func (c *Client) IsActive(ctx context.Context, name string) (bool, error) {
	val, err := c.rdb.Get(ctx, "flag:"+name).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	var flag model.Flag
	err = json.Unmarshal([]byte(val), &flag)
	if err != nil {
		return false, err
	}

	return flag.Active, nil
}

func New(rdb *redis.Client) (*Client, error) {
	if rdb == nil {
		return nil, errors.New("redis connection is nil")
	}

	return &Client{
		rdb: rdb,
	}, nil
}
