package request

import (
	"database/sql"
	// "encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	// "sort"
	// "strconv"
	// "strings"
	"time"
	// "encoding/json"

	"demo/common"
	"demo/i18n"

	_ "github.com/lib/pq"
)

// Confirm page
func Confirm(w http.ResponseWriter, r *http.Request) {
	
	uid := common.GetUser(w, r)

	if uid == 0 {
		http.Redirect(w, r, "/sign/?redirect="+url.QueryEscape(r.URL.Path), 302)
	}
	lang := common.GetLang(w, r)
	connStr := common.DbConnect
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Print(err)
	}
	defer db.Close()
	_, err = strconv.Atoi(r.FormValue("id"))
	if err != nil {
		log.Print("sku does not exist")
		http.Redirect(w, r, "/", 302)
		return
	}
	skuID := r.FormValue("id")
	
	query := `SELECT sku_quantity FROM t_sku WHERE sku_id = $1`
	rows, err := db.Query(query,skuID)
	if err != nil {
		log.Print(query,"\n",skuID,"\n",err)
	}
	is := false
	for rows.Next() {
		quantity := ""
		if err := rows.Scan(&quantity); err != nil {
			log.Print(err)
		}
		is = true
		if quantity < r.FormValue("quantity") {
			log.Print("quantity fraud")
			http.Redirect(w, r, "/", 302)
			return
		}
	}
	if !is {
		log.Print("sku does not exist")
		http.Redirect(w, r, "/", 302)
		return
	} 

	query = `SELECT sku_id, sku_quantity FROM t_cart WHERE buyer_id = $1`
	rows, err = db.Query(query, uid)
	if err != nil {
		log.Print(query,"\n",uid,"\n",err)
	}

	carts := map[string]map[string]string{}
	is = false
	diff := false
	skuIDs := r.FormValue("id")
	tmpSkuID := ""
	for rows.Next() {
		q := ""
		if err := rows.Scan(&tmpSkuID,&q); err != nil {
			log.Print(err)
		}
		skuIDs = skuIDs + "," + tmpSkuID
		if tmpSkuID == skuID {
			is = true
			if q != r.FormValue("quantity") {
				diff = true
				q = r.FormValue("quantity")
			}
		}
		cart := map[string]string{}
		cart["SkuQuantity"] = q
		carts[tmpSkuID] = cart
	}

	if is {
		if diff {
			query = `UPDATE t_cart SET sku_quantity = $1, updated_at = $2 
					WHERE sku_id = $3`
			_, err = db.Exec(query,r.FormValue("quantity"),
				time.Now().Format("2006-01-02 15:04:05"),skuID)
			if err != nil {
				log.Print(query,"\n",r.FormValue("quantity"),
					time.Now().Format("2006-01-02 15:04:05"),skuID,"\n",err)
			}			
		}

	} else {
		query = `INSERT INTO t_cart(
			sku_id   
			,buyer_id    
			,sku_quantity    
			,updated_at
			) VALUES($1,$2,$3,$4)`
		_, err = db.Exec(query,
		skuID,
		uid,
		r.FormValue("quantity"),
		time.Now().Format("2006-01-02 15:04:05") )
		if err != nil {
			log.Print(query,"\n",r.FormValue("id"),uid,r.FormValue("quantity"),
				time.Now().Format("2006-01-02 15:04:05"),"\n",err)
		}
	}
	query = `SELECT sku_id, sku_img, item_id, sku_price FROM t_sku WHERE sku_id in (`+skuIDs+`)`
	rows, err = db.Query(query)
	if err != nil {
		log.Print(query,"\n",skuIDs,"\n",err)
	}
	log.Print(skuIDs)
	itemIDs := "0"
	tmpItemID := ""
	for rows.Next() {
		img := ""
		p := ""
		if err := rows.Scan(&tmpSkuID,&img,&tmpItemID,&p); err != nil {
			log.Print(err)
		}
		itemIDs = itemIDs + "," + tmpItemID
		cart := map[string]string{}
		cart["SkuImg"] = img
		cart["SkuPrice"] = p
		cart["ItemID"] = tmpItemID
		cart["SkuQuantity"] = carts[tmpSkuID]["SkuQuantity"]
		carts[tmpSkuID] = cart
	}
	query = `SELECT item_id, item_name FROM t_translation_item 
		WHERE item_id in (`+itemIDs+`) AND language = $1`
	rows, err = db.Query(query,lang)
	if err != nil {
		log.Print(query,"\n",itemIDs,"\n",lang,"\n",err)
	}
	for rows.Next() {
		itemName := ""
		if err := rows.Scan(&tmpItemID,&itemName); err != nil {
			log.Print(err)
		}
		for k, v := range carts {
			if tmpItemID == v["ItemID"] {
				carts[k]["ItemName"] = itemName
			}

		}
	}
	fmt.Println("%v", carts)
	type View struct {
		Footer		map[string]string
		I18n		map[string]string
		Carts		map[string]map[string]string
	}
	var view View
	view.Footer = i18n.Footer(lang)
	view.I18n = i18n.Confirm(lang)
	view.Carts = carts
	tpl := template.Must(template.ParseFiles("view/confirm.tmpl"))
	tpl.Execute(w, view)
}
