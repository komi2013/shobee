package request

import (
	"fmt"
	"net/http"
	"html/template"

	"demo/common"
	// _ "github.com/lib/pq" // this driver for postgres
)

// Sign page
func Sign(w http.ResponseWriter, r *http.Request) {

	// common.SetUser(w, r, 2)
	lang := common.GetLang(w, r)
	fmt.Println("lang", lang)
	type View struct {
		ReviewScore     float64
	}
	var view View
	tpl := template.Must(template.ParseFiles("view/sign.tmpl"))
	tpl.Execute(w, view)
}
