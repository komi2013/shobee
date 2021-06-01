package request

import (
	"bytes"
	"fmt"
	"database/sql"
	"encoding/base64"
	"log"
	"net/http"
	"net/smtp"
	"time"
	"html/template"

	"demo/common"
	"golang.org/x/crypto/bcrypt"
	_ "github.com/lib/pq"
)

// RegisterEmail
func RegisterEmail(w http.ResponseWriter, r *http.Request) {
	// common.SetUser(w, r, 2)
	lang := common.GetLang(w, r)

	connStr := common.DbConnect
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Print(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT buyer_id FROM t_buyer WHERE email = $1",r.FormValue("email"))
	if err != nil {
		log.Print(err)
	}
	var buyerID string
	for rows.Next() {
		if err := rows.Scan(&buyerID); err != nil {
			log.Print(err)
		}
		fmt.Fprint(w, `{"Status":"2","Msg":"already this email"}`)
		return
	}

	rows, err = db.Query("SELECT NEXTVAL('t_buyer_buyer_id_seq')")
	if err != nil {
		log.Print(err)
	}
	for rows.Next() {
		if err := rows.Scan(&buyerID); err != nil {
			log.Print(err)
		}
	}
	token := common.StringRand(20)
	password, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password1")), 14)
    
	_, err = db.Exec(`INSERT INTO t_buyer(
		buyer_id
		,auth_token   
		,email    
		,password    
		,family_name 
		,given_name
		,birthday
		,updated_at
		) VALUES($1,$2,$3,$4,$5,$6,$7,$8)`,
		buyerID,
		token,
		r.FormValue("email"),
		"waitingVerify" + string(password),
		r.FormValue("familyName"),
		r.FormValue("givenName"),
		r.FormValue("birthday"),
		time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		fmt.Println(err)
	}

	tpl := template.Must(template.ParseFiles("mailTmpl/registration.tmpl"))
	type View struct {
		Token     string
		BuyerID  string
	}
	var view View
	view.Token = token
	view.BuyerID = buyerID
	buf := new(bytes.Buffer)
	tpl.Execute(buf, view)
	// fmt.Println(buf.String())

	// Connect to the remote SMTP server.
	c, err := smtp.Dial("localhost:25")
	if err != nil {
		fmt.Println(err)
	}

	// Set the sender and recipient first
	if err := c.Mail("sender@quigen.info"); err != nil {
		fmt.Println(err)
	}
	if err := c.Rcpt("komatsu@beenos.com"); err != nil {
		fmt.Println(err)
	}

	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		fmt.Println(err)
	}
	msg := "To: komatsu@beenos.com \r\n" +
		"From: sender@quigen.info \r\n" +
		"Subject: you got comment \r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"Content-Transfer-Encoding: base64\r\n" +
		"\r\n" + base64.StdEncoding.EncodeToString([]byte(buf.String()))
	_, err = wc.Write([]byte(msg))
	if err != nil {
		fmt.Println(err)
	}
	err = wc.Close()
	if err != nil {
		fmt.Println(err)
	}
	err = c.Quit()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprint(w, `{"Status":"1"}`+ lang)
}
