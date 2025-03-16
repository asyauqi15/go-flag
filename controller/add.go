package controller

import (
	"encoding/json"
	"github.com/asyauqi15/go-flag/model"
	"html/template"
	"net/http"
	"net/url"
)

func (c *Controller) Add(w http.ResponseWriter, r *http.Request) {
	var tmpl, err = template.ParseFS(c.templateFS, "views/add.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	queryParams := r.URL.Query()
	errorMsg := queryParams.Get("error")

	err = tmpl.Execute(w, templateData{
		RootPath: c.rootPath,
		ErrorMsg: errorMsg,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *Controller) AddProcess(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get form values
	name := r.FormValue("name")
	value := r.FormValue("value")
	active := r.FormValue("active") == "on"

	if name == "" {
		http.Error(w, "Feature name is required", http.StatusBadRequest)
		return
	}

	// Check if the key already exists in Redis
	exists, err := c.rdb.Exists(ctx, "flag:"+name).Result()
	if err != nil {
		http.Error(w, "Error checking feature existence", http.StatusInternalServerError)
		return
	}
	if exists > 0 {
		queryParams := url.Values{}
		queryParams.Set("error", "Feature name already exists")

		redirectURL := c.rootPath + "/add?" + queryParams.Encode()

		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	// Store in Redis
	flag := model.Flag{Name: name, Value: value, Active: active}
	flagJSON, _ := json.Marshal(flag)
	err = c.rdb.Set(ctx, "flag:"+name, flagJSON, 0).Err()
	if err != nil {
		http.Error(w, "Error saving feature", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, c.rootPath, http.StatusSeeOther)
}
