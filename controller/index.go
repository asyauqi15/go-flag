package controller

import (
	"encoding/json"
	"github.com/asyauqi15/go-flag/model"
	"html/template"
	"log"
	"net/http"
	"sort"
)

func (c *Controller) Index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	keys, err := c.rdb.Keys(ctx, "flag:*").Result()
	if err != nil {
		http.Error(w, "Error fetching keys", http.StatusInternalServerError)
		return
	}
	sort.Strings(keys)

	var flags []model.Flag
	for _, key := range keys {
		data, err := c.rdb.Get(ctx, key).Result()
		if err != nil {
			log.Println("Error retrieving", key, err)
			continue
		}

		var flag model.Flag
		if err := json.Unmarshal([]byte(data), &flag); err != nil {
			log.Println("Error unmarshaling JSON", key, err)
			continue
		}

		flags = append(flags, flag)
	}

	tmpl, err := template.ParseFS(c.templateFS, "views/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, flags)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
