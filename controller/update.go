package controller

import (
	"encoding/json"
	"github.com/asyauqi15/go-flag/model"
	"github.com/go-chi/chi/v5"
	"html/template"
	"net/http"
)

func (c *Controller) Update(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "feature_name")

	val, err := c.rdb.Get(r.Context(), "flag:"+name).Result()
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

	err = tmpl.Execute(w, flag)
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
	err := c.rdb.Set(r.Context(), "flag:"+name, flagJSON, 0).Err()
	if err != nil {
		http.Error(w, "Error updating feature", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/flag", http.StatusSeeOther)
}
