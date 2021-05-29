package main

import (
	"fmt"
	"log"
	"net/http"

	"demo/request"

	"demo/common"
)

func main() {

	http.HandleFunc("/htm/", request.Htm)
	http.HandleFunc("/item/", request.Item)
	http.HandleFunc("/list/", request.List)
	http.HandleFunc("/pay/", request.Pay)
	http.HandleFunc("/registerEmail/", request.RegisterEmail)
	http.HandleFunc("/sign/", request.Sign)
	http.HandleFunc("/test/", request.Test)

	
	

	// http.HandleFunc("/", app.Top)
	// fmt.Println("starting.." + common.CacheV)
	fmt.Println(common.DbConnect)

	log.Fatal(http.ListenAndServe(common.GoPort, nil))

}
