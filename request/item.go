package request

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
	"encoding/json"

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
	FROM t_sku WHERE item_id = $1` // if multiple variation
	rows, err := db.Query(query, itemId)
	if err != nil {
		fmt.Println(query)
		fmt.Println(err)
	}
	type sku struct {
		SkuImgs     []string
		SkuPrice    float64
		ModelNumber string
		CategoryId  int
		GenreId     int
		SkuQuantity float64
	}
	var s sku
	for rows.Next() {
		imgs := ""
		if err := rows.Scan(&imgs, &s.SkuPrice, &s.ModelNumber, &s.CategoryId, &s.GenreId, &s.SkuQuantity); err != nil {
			log.Print(err)
		}
		// s.SkuImgs, _ := json.Marshal(arr)
		json.Unmarshal([]byte(imgs), &s.SkuImgs)
		
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
	topCategory := ""
	for rows.Next() {
		r := mCategoryTree{}
		if err := rows.Scan(&r.LeafID, &r.Level1, &r.Level2, &r.Level3, &r.Level4, &r.Level5, &r.Level6, &r.Level7, &r.Level8, &r.UpdatedAt); err != nil {
			log.Print(err)
		}
		whereIn = strconv.Itoa(r.Level1)
		topCategory = strconv.Itoa(r.Level1)
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

	rows, err = db.Query(`SELECT category_id, category_name_` + lang +
		` FROM m_category WHERE category_id in (` + whereIn + `)`)
	if err != nil {
		log.Print(err)
	}
	for rows.Next() {
		r := category{}
		if err := rows.Scan(&r.CategoryID, &r.CategoryName); err != nil {
			log.Print(err)
		}
		convert[r.CategoryID] = r.CategoryName
	}
	for _, v := range tree {
		y := category{}
		y.Level = v[0]
		y.CategoryID = v[1]
		y.CategoryName = convert[v[1]]
		breadCrumb = append(breadCrumb, y)
	}
	sort.Slice(breadCrumb, func(i, j int) bool { return breadCrumb[i].Level < breadCrumb[j].Level }) // DESC

	query = `SELECT variation_id, variation_type,variation_description_` + lang + // lang is not allow special characters
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
		list["variation_type"] = flg
		list["variation_value"] = description
		list["tier_type"] = tier
		vList[id] = list
	}
	fmt.Printf("%#v\n", vList)
	query = `SELECT variation_id, variation_name_` + lang + `,variation_value_` + lang +
		`, search_category_id, search_type FROM m_variation WHERE variation_id in (` + ids +
		`) OR search_category_id in (` + whereIn + `,0) `
	rows, err = db.Query(query)
	if err != nil {
		fmt.Println(query)
		fmt.Println(err)
	}
	searchList :=  map[string][][]string{}
	for rows.Next() {
		id := 0
		name := ""
		value := ""
		catId := -1
		sType := "0"
		if err := rows.Scan(&id, &name, &value, &catId, &sType); err != nil {
			log.Print(err)
		}

		if _, ok := vList[id]; ok {
			list := map[string]string{}
			list["variation_name"] = name
			if vList[id]["variation_type"] == "0" {
				list["variation_value"] = value
			} else {
				list["variation_value"] = vList[id]["variation_value"]
			}
			vList[id] = list
		}
		if catId > -1 { // search_category_id records < variation_id records
			list2 := []string{strconv.Itoa(id),value,sType}
			searchList[name] = append(searchList[name],list2)
		}
	}

	type search struct {
		Name    string
		SType   string
		Values  [][]string
	}
	var seaList []search
	for k, v := range searchList {
		var sea search
		sea.Name   = k
		var list2  [][]string
		for _, v2 := range v {
			sea.SType  = v2[2]
			list3 := []string{v2[0],v2[1]}
			list2 = append(list2,list3)
		}
		sea.Values = list2
		seaList = append(seaList,sea)
	}
	fmt.Printf("%#v\n", ids)
	
	query = `SELECT ROUND(AVG(review_score),1) FROM h_review 
		WHERE category_id = $1 AND genre_id = $2
		GROUP BY category_id, genre_id`
	rows, err = db.Query(query, s.CategoryId, s.GenreId)
	if err != nil {
		fmt.Println(query)
		fmt.Println(err)
	}
	var reviewScore float64
	for rows.Next() {
		if err := rows.Scan(&reviewScore); err != nil {
			log.Print(err)
		}
	}
	query = `SELECT genre_id, genre_name_` + lang +` FROM m_genre`
	rows, err = db.Query(query)
	if err != nil {
		fmt.Println(query)
		fmt.Println(err)
	}
	type genre struct {
		ID   string
		Name string
	}
	var genres []genre
	for rows.Next() {
		var g genre
		if err := rows.Scan(&g.ID,&g.Name); err != nil {
			log.Print(err)
		}
		genres = append(genres,g)
	}
	type View struct {
		Sku             sku
		TranslationItem translationItem
		VariationList   map[int]map[string]string
		BreadCrumb      []category
		ReviewScore     float64
		SearchList      []search
		TopCategory		string
		Genres          []genre
	}
	var view View
	view.Sku = s
	view.TranslationItem = t
	view.VariationList = vList
	view.BreadCrumb = breadCrumb
	view.ReviewScore = reviewScore
	view.SearchList = seaList
	view.TopCategory = topCategory
	view.Genres = genres
	tpl := template.Must(template.ParseFiles("view/item.tmpl"))
	tpl.Execute(w, view)
}
