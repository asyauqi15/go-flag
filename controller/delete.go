package controller

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (c *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "feature_name")

	c.rdb.Del(r.Context(), "flag:"+name)
	http.Redirect(w, r, "/flag", http.StatusSeeOther)
}
