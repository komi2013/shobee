package main

import (
	"fmt"
	"log"
	"net/http"

	"demo/app"
	"demo/common"
)

func main() {

	http.HandleFunc("/htm/", app.Htm)
	http.HandleFunc("/item/", app.Item)
	http.HandleFunc("/list/", app.List)
	http.HandleFunc("/test/", app.Test)

	// http.HandleFunc("/", app.Top)
	// fmt.Println("starting.." + common.CacheV)
	fmt.Println(common.DbConnect)

	log.Fatal(http.ListenAndServe(common.GoPort, nil))

}
