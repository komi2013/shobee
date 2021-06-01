package request

import (
	// "fmt"
	"database/sql"
	"log"
	"net/http"
	"time"
	"strings"

	
	"demo/common"	
	_ "github.com/lib/pq" // this driver for postgres
)

// MailVerify
func MailVerify(w http.ResponseWriter, r *http.Request) {
	connStr := common.DbConnect
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	u := strings.Split(r.URL.Path, "/")
	if u[2] == "" {
		http.Redirect(w, r, "/sign/?warn=3", 302)
	}
	rows, err := db.Query("SELECT auth_token, updated_at FROM t_buyer WHERE buyer_id = $1",u[2])
	if err != nil {
		log.Print(err)
	}
	expire := time.Now().Add(time.Duration(-24) * time.Hour)
	var token string
	var updatedAt string
	for rows.Next() {
		if err := rows.Scan(&token,&updatedAt); err != nil {
			log.Print(err)
		}
	}
	if token != u[3] {
		http.Redirect(w, r, "/sign/?warn=1", 302)
	}
	tokenExpired, err := time.Parse("2006-01-02T15:04:05Z", updatedAt)
	if err != nil {
		log.Print(err)
	}
	log.Print(expire.Format("2006-01-02 15:04:05"))
	log.Print(tokenExpired.Format("2006-01-02 15:04:05"))
	if tokenExpired.Before(expire) {
		http.Redirect(w, r, "/sign/?warn=2", 302)
	}
	_, err = db.Exec(`UPDATE t_buyer SET password = REPLACE(password, 'waitingVerify', '')
		WHERE buyer_id = $1`, u[2])
	if err != nil {
		log.Print(err)
	}
	http.Redirect(w, r, "/sign/", 302)

}
