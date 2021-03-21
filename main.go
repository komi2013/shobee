package main

import (
	"fmt"
	"log"
	"net/http"

	"./common"
)

func main() {

	// http.HandleFunc("/category/", app.Category)
	// http.HandleFunc("/item/", app.Item)
	// http.HandleFunc("/profile/", app.Profile)
	// http.HandleFunc("/bought/", app.Bought)
	// http.HandleFunc("/cart/", app.Cart)

	// http.HandleFunc("/", app.Top)
	fmt.Println("starting.." + common.CacheV)
	fmt.Println(common.DbConnect)

	log.Fatal(http.ListenAndServe(common.GoPort, nil))
}
