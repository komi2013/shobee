package app

import (
	// "html/template"
	"fmt"
	// "log"
	"net/http"
	"strings"
)

// List
func List(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("%#v\n", r.FormValue("word"))
	fmt.Printf("%#v\n", r.FormValue("column1"))
	fmt.Printf("%#v\n", r.FormValue("column2"))
	fmt.Printf("%#v\n", r.FormValue("price"))

	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	prices := r.FormValue("price")
	fmt.Printf("%#v\n", prices)
	// log.Println(prices[0])

}
