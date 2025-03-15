package controller

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path"
)

func (c *Controller) Add(w http.ResponseWriter, r *http.Request) {
	var filepath = path.Join("views", "add.html")
	var tmpl, err = template.ParseFiles(filepath)
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
	active := r.FormValue("active") == "on" // Checkbox sends "on" if checked

	if name == "" {
		http.Error(w, "Feature name is required", http.StatusBadRequest)
		return
	}

	// Store in Redis
	ctx := r.Context()
	flag := Flag{Name: name, Active: active}
	flagJSON, _ := json.Marshal(flag)
	err := c.rdb.Set(ctx, "flag:"+name, flagJSON, 0).Err()
	if err != nil {
		http.Error(w, "Error saving feature", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/flag", http.StatusSeeOther)
}
