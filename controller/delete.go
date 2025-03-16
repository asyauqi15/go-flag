package controller

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (c *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "feature_name")

	key := fmt.Sprintf("%s:%s", c.keyPrefix, name)

	c.rdb.Del(r.Context(), key)
	http.Redirect(w, r, c.rootPath, http.StatusSeeOther)
}
