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

// Top
func Top(w http.ResponseWriter, r *http.Request) {

	// common.SetUser(w, r, 2)
	lang := common.GetLang(w, r)
	buyerID := common.GetUser(w, r)
	fmt.Println("lang", lang)
	fmt.Println(buyerID)
	fmt.Fprint(w, `{"Status":"1"}`)
}
