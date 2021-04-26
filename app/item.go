package app

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"

	// "strconv"
	"strings"
	"time"

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

	query := `SELECT sku_imgs,sku_price,model_number,category_id,genre_id,sku_quantity
	FROM t_sku WHERE item_id = $1 limit 1` // if multiple variation
	rows, err := db.Query(query, itemId)
	if err != nil {
		fmt.Println(query)
		fmt.Println(err)
	}
	type sku struct {
		SkuImgs     string
		SkuPrice    float64
		ModelNumber string
		CategoryId  int
		GenreId     int
		SkuQuantity float64
	}
	var s sku
	for rows.Next() {
		if err := rows.Scan(&s.SkuImgs, &s.SkuPrice, &s.ModelNumber, &s.CategoryId, &s.GenreId, &s.SkuQuantity); err != nil {
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
	type translationItem struct {
		ItemName        string
		ItemDescription string
		ItemExactName   string
	}
	var t translationItem
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

	type category struct {
		Level        int
		CategoryID   int
		CategoryName string
	}

	type mCategoryTree struct {
		LeafID    int // leaf_id
		Level1    int // level_1
		Level2    int // level_2
		Level3    int // level_3
		Level4    int // level_4
		Level5    int // level_5
		Level6    int // level_6
		Level7    int // level_7
		Level8    int // level_8
		UpdatedAt time.Time
	}

	convert := make(map[int]string)
	var tree [][2]int
	var x [2]int
	rows, err = db.Query("SELECT * FROM m_category_tree WHERE leaf_id = $1", s.CategoryId)
	if err != nil {
		log.Print(err)
	}
	whereIn := ""
	for rows.Next() {
		r := mCategoryTree{}
		if err := rows.Scan(&r.LeafID, &r.Level1, &r.Level2, &r.Level3, &r.Level4, &r.Level5, &r.Level6, &r.Level7, &r.Level8, &r.UpdatedAt); err != nil {
			log.Print(err)
		}
		whereIn = strconv.Itoa(r.Level1)
		x[0] = 1
		x[1] = r.Level1
		tree = append(tree, x)
		convert[r.Level1] = ""
		if r.Level2 > 0 {
			whereIn = whereIn + "," + strconv.Itoa(r.Level2)
			x[0] = 2
			x[1] = r.Level2
			tree = append(tree, x)
			convert[r.Level2] = ""
		}
		if r.Level3 > 0 {
			whereIn = whereIn + "," + strconv.Itoa(r.Level3)
			x[0] = 3
			x[1] = r.Level3
			tree = append(tree, x)
			convert[r.Level3] = ""
		}
		if r.Level4 > 0 {
			whereIn = whereIn + "," + strconv.Itoa(r.Level4)
			x[0] = 4
			x[1] = r.Level4
			tree = append(tree, x)
			convert[r.Level4] = ""
		}
		if r.Level5 > 0 {
			whereIn = whereIn + "," + strconv.Itoa(r.Level5)
			x[0] = 5
			x[1] = r.Level5
			tree = append(tree, x)
			convert[r.Level5] = ""
		}
		if r.Level6 > 0 {
			whereIn = whereIn + "," + strconv.Itoa(r.Level6)
			x[0] = 6
			x[1] = r.Level6
			tree = append(tree, x)
			convert[r.Level6] = ""
		}
		if r.Level7 > 0 {
			whereIn = whereIn + "," + strconv.Itoa(r.Level7)
			x[0] = 7
			x[1] = r.Level7
			tree = append(tree, x)
			convert[r.Level7] = ""
		}
		if r.Level8 > 0 {
			whereIn = whereIn + "," + strconv.Itoa(r.Level8)
			x[0] = 8
			x[1] = r.Level8
			tree = append(tree, x)
			convert[r.Level8] = ""
		}
	}
	var breadCrumb []category
	if whereIn != "" {
		rows, err = db.Query(`SELECT category_id, category_name_` + lang + ` FROM m_category WHERE category_id in (` + whereIn + `)`)
		if err != nil {
			log.Print(err)
		}
		for rows.Next() {
			r := category{}
			if err := rows.Scan(&r.CategoryID, &r.CategoryName); err != nil {
				log.Print(err)
			}
			convert[r.CategoryID] = r.CategoryName
			// categoryName = append(categoryName, r)
			// breadCrumb[r.CategoryID]["name"] = r.CategoryName
		}
		for _, v := range tree {
			y := category{}
			y.Level = v[0]
			y.CategoryID = v[1]
			y.CategoryName = convert[v[1]]
			breadCrumb = append(breadCrumb, y)
		}
		sort.Slice(breadCrumb, func(i, j int) bool { return breadCrumb[i].Level < breadCrumb[j].Level }) // DESC
	}

	type View struct {
		Sku             sku
		TranslationItem translationItem
		VariationList   map[int]map[string]string
		BreadCrumb      []category
	}
	var view View
	view.Sku = s
	view.TranslationItem = t
	view.VariationList = vList
	view.BreadCrumb = breadCrumb
	tpl := template.Must(template.ParseFiles("view/item.tmpl"))
	tpl.Execute(w, view)
}
