package request

import (
	"fmt"
	"net/http"
	"html/template"

	"demo/common"
	"demo/i18n"
	// _ "github.com/lib/pq" // this driver for postgres
)

// Sign page
func Sign(w http.ResponseWriter, r *http.Request) {

	// common.SetUser(w, r, 2)
	lang := common.GetLang(w, r)
	// r.FormValue("warn")

	// fmt.Println("%v", footer)
	fmt.Println("%v", i18n.Footer(lang))
	type View struct {
		Footer		map[string]string
		I18n		map[string]string
	}
	var view View
	view.Footer = i18n.Footer(lang)
	view.I18n = i18n.Sign(lang)
	tpl := template.Must(template.ParseFiles("view/sign.tmpl"))
	tpl.Execute(w, view)
}
