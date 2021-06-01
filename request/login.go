package request

import (
	"fmt"
	"database/sql"
	"log"
	"net/http"
	// "time"

	"demo/common"
	"golang.org/x/crypto/bcrypt"
	_ "github.com/lib/pq"
)

// Login
func Login(w http.ResponseWriter, r *http.Request) {
	connStr := common.DbConnect
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Print(err)
	}
	defer db.Close()
	
	query := `SELECT buyer_id, password FROM t_buyer WHERE email = $1`
	rows, err := db.Query(query, r.FormValue("email"))
	if err != nil {
		log.Print(err)
	}
	buyerID := 0
	hash := ""
	for rows.Next() {
		if err := rows.Scan(&buyerID,&hash); err != nil {
			log.Print(err)
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(r.FormValue("password")))
	if err != nil {
		fmt.Fprint(w, `{"Status":"2"}`)
	} else {
		common.SetUser(w, r, buyerID)
		cookie, err := r.Cookie("redirect")
		redirect := "/"
		if err == nil {
			redirect = cookie.Value
		}
		fmt.Fprint(w, `{"Status":"1","Redirect":"`+redirect+`"}`)
	}
}
