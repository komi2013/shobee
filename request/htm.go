package request

import (
	"html/template"
	"net/http"
	"strings"
)

// Htm can make many html pages
func Htm(w http.ResponseWriter, r *http.Request) {
	type View struct {
		CacheV string
	}
	var view View

	uri := strings.Split(r.URL.String(), "/")
	tpl := template.Must(template.ParseFiles("view/htm/" + uri[2] + ".html"))
	tpl.Execute(w, view)

}
