package controller

import (
	"encoding/json"
	"fmt"
	"github.com/asyauqi15/go-flag/model"
	"github.com/go-chi/chi/v5"
	"html/template"
	"net/http"
)

func (c *Controller) Update(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "feature_name")

	key := fmt.Sprintf("%s:%s", c.keyPrefix, name)

	val, err := c.rdb.Get(r.Context(), key).Result()
	if err != nil {
		http.Error(w, "Feature not found", http.StatusNotFound)
		return
	}

	var flag model.Flag
	err = json.Unmarshal([]byte(val), &flag)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFS(c.templateFS, "views/update.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, templateData{
		RootPath: c.rootPath,
		Data:     flag,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *Controller) UpdateProcess(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "feature_name")
	value := r.FormValue("value")
	active := r.FormValue("active") == "on"

	flag := model.Flag{Name: name, Value: value, Active: active}
	flagJSON, _ := json.Marshal(flag)

	key := fmt.Sprintf("%s:%s", c.keyPrefix, name)

	err := c.rdb.Set(r.Context(), key, flagJSON, 0).Err()
	if err != nil {
		http.Error(w, "Error updating feature", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, c.rootPath, http.StatusSeeOther)
}
