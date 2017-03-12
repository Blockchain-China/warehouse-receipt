/*
Copyright DASE@ECNU. 2016 All Rights Reserved.
*/

package main

import (
	"net/http"
	"fmt"

	"github.com/ant0ine/go-json-rest/rest"
)

func main() {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/transactions/:id", GetTransactions),
	)
	if err != nil {
		fmt.Println(err)
	}
	api.SetApp(router)
	fmt.Println(http.ListenAndServe(":8100", api.MakeHandler()))
}

func GetTransactions(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
    txArr := getTransactions2("120.132.12.128", "7050", id)
    if len(txArr)>0 {
    	w.WriteJson(txArr)
    }else{
        msg := "{\"msg\":\"Not Found\"}"
        w.(http.ResponseWriter).Write([]byte(msg))
    } 
    

 //    txid := r.PathParam("txid")
 //    blockNumber := strconv.Itoa(getBlockNumber("120.132.12.128", "7050", txid))
 //    msg := "{\"blockNumber\":\""+ blockNumber + "\"}"
 //    w.(http.ResponseWriter).Write([]byte(msg))
}

// GET
// curl -i http://127.0.0.1:8100/transactions/warehouse-1
