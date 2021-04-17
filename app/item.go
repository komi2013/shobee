package app

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	// "strconv"
	"strings"

	"demo/common"

	_ "github.com/lib/pq"
)

// Item page
func Item(w http.ResponseWriter, r *http.Request) {

	// common.SetUser(w, r, 2)
	// uid := common.GetUser(w, r)
	// fmt.Println("uid", uid)

	lang := common.GetLang(w, r)
	u := strings.Split(r.URL.Path, "/")
	itemId := u[2]
	connStr := common.DbConnect
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	query := `SELECT article_number,shipping_fee,sku_status,sku_imgs,sku_price,model_number,
	cost,service_charge,weight,category_id,genre_id,sku_quantity
	FROM t_sku WHERE item_id = $1 limit 1` // if multiple variation
	rows, err := db.Query(query, itemId)
	if err != nil {
		fmt.Println(query)
		fmt.Println(err)
	}
	type Sku struct {
		ArticleNumber string
		ShippingFee   float64
		SkuStatus     string
		SkuImgs       string
		SkuPrice      float64
		ModelNumber   string
		Cost          float64
		ServiceCharge float64
		Weight        float64
		CategoryId    int
		GenreId       int
		SkuQuantity   float64
	}
	var s Sku
	for rows.Next() {
		if err := rows.Scan(&s.ArticleNumber, &s.ShippingFee, &s.SkuStatus, &s.SkuImgs, &s.SkuPrice, &s.ModelNumber,
			&s.Cost, &s.ServiceCharge, &s.Weight, &s.CategoryId, &s.GenreId, &s.SkuQuantity); err != nil {
			log.Print(err)
		}
	}

	query = `SELECT item_name,item_description,item_exact_name
	FROM t_translation_item WHERE item_id = $1 AND language = $2`
	rows, err = db.Query(query, itemId, lang)
	if err != nil {
		fmt.Println(query)
		fmt.Println(err)
	}
	type TranslationItem struct {
		ItemName        string
		ItemDescription string
		ItemExactName   string
	}
	var t TranslationItem
	for rows.Next() {
		if err := rows.Scan(&t.ItemName, &t.ItemDescription, &t.ItemExactName); err != nil {
			log.Print(err)
		}
	}

	query = `SELECT variation_id, description_flg,variation_description_` + lang + // lang is not allow special characters
		`, tier_type FROM t_variation WHERE item_id = $1`
	rows, err = db.Query(query, itemId)
	if err != nil {
		fmt.Println(query)
		fmt.Println(err)
	}
	vList := map[int]map[string]string{}
	ids := "-1"
	for rows.Next() {
		list := map[string]string{}
		id := 0
		flg := ""
		description := ""
		tier := ""
		if err := rows.Scan(&id, &flg, &description, &tier); err != nil {
			log.Print(err)
		}
		ids = ids + "," + strconv.Itoa(id)
		list["description_flg"] = flg
		list["variation_description"] = description
		list["tier_type"] = tier
		vList[id] = list
	}

	query = `SELECT variation_id, variation_name_` + lang + `,variation_value_` + lang +
		` FROM m_variation WHERE variation_id in (` + ids + `)`
	rows, err = db.Query(query)
	if err != nil {
		fmt.Println(query)
		fmt.Println(err)
	}
	for rows.Next() {
		list := map[string]string{}
		id := 0
		name := ""
		value := ""
		if err := rows.Scan(&id, &name, &value); err != nil {
			log.Print(err)
		}
		list["variation_name"] = name
		list["variation_value"] = value
		vList[id] = list
	}
	fmt.Printf("%#v\n", vList)
	// fmt.Printf("%#v\n", ids)

	type View struct {
		// CacheV       string
		// Q            common.TQuestion
		// Available    int
		// BreadCrumb   []BreadCrumb
		// Title        string
		// Qtxt         template.HTML
		Sku             Sku
		TranslationItem TranslationItem
		VariationList   map[int]map[string]string
	}
	var view View
	view.Sku = s
	view.TranslationItem = t
	view.VariationList = vList
	//m_category_tree, m_category
	tpl := template.Must(template.ParseFiles("view/item.tmpl"))
	tpl.Execute(w, view)
	// fmt.Fprint(w, `{"Status":"1"}`)
}
