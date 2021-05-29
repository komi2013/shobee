package request

import (
	"fmt"
	// "database/sql"
	// "log"
	"net/http"
	// "strconv"

	"demo/common"
	// _ "github.com/lib/pq" // this driver for postgres
)

// Pay page
func Pay(w http.ResponseWriter, r *http.Request) {

	// common.SetUser(w, r, 2)
	lang := common.GetLang(w, r)
	fmt.Println("lang", lang)
	fmt.Fprint(w, `{"Status":"1"}`)
}
