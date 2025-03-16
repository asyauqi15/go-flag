package controller

import (
	"encoding/json"
	"github.com/asyauqi15/go-flag/model"
	"html/template"
	"net/http"
)

func (c *Controller) Add(w http.ResponseWriter, r *http.Request) {
	var tmpl, err = template.ParseFS(c.templateFS, "views/add.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *Controller) AddProcess(w http.ResponseWriter, r *http.Request) {
	// Get form values
	name := r.FormValue("name")
	value := r.FormValue("value")
	active := r.FormValue("active") == "on"

	if name == "" {
		http.Error(w, "Feature name is required", http.StatusBadRequest)
		return
	}

	// Store in Redis
	ctx := r.Context()
	flag := model.Flag{Name: name, Value: value, Active: active}
	flagJSON, _ := json.Marshal(flag)
	err := c.rdb.Set(ctx, "flag:"+name, flagJSON, 0).Err()
	if err != nil {
		http.Error(w, "Error saving feature", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/flag", http.StatusSeeOther)
}
