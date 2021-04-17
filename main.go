package main

import (
	"fmt"
	"log"
	"net/http"

	"demo/app"
	"demo/common"
)

func main() {

	http.HandleFunc("/test/", app.Test)
	http.HandleFunc("/item/", app.Item)
	http.HandleFunc("/htm/", app.Htm)

	// http.HandleFunc("/", app.Top)
	// fmt.Println("starting.." + common.CacheV)
	fmt.Println(common.DbConnect)

	log.Fatal(http.ListenAndServe(common.GoPort, nil))

}
