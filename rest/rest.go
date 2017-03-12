/*
Copyright DASE@ECNU. 2016 All Rights Reserved.
*/

package main

import (
	"net/http"
	"fmt"
	//"encoding/json"

	"github.com/ant0ine/go-json-rest/rest"

	"github.com/hyperledger/fabric/work/receipt/entity"
	"github.com/hyperledger/fabric/work/receipt/app1"
	"github.com/hyperledger/fabric/work/receipt/app2"
	"github.com/hyperledger/fabric/work/receipt/query"
	"github.com/hyperledger/fabric/work/receipt/member"
)

func main() {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	// api.Use(&rest.AuthBasicMiddleware{
	// 	Realm: "test zone",
	// 	Authenticator: func(userId string, password string) bool {
	// 		if userId == "admin" && password == "admin" {
	// 			return true
	// 		}
	// 		return false
	// 	},
	// })
	router, err := rest.MakeRouter(
		//rest.Get("/receipts", GetReceipts),
		rest.Get("/receipts/:username", GetReceipts),
		rest.Post("/receipts", PostReceipt),
		rest.Put("/receipts/:username", PutReceipt),

		rest.Post("/members", PostMember),
	)
	if err != nil {
		fmt.Println(err)
	}
	api.SetApp(router)
	fmt.Println(http.ListenAndServe(":8888", api.MakeHandler()))
}

func GetReceipt(w rest.ResponseWriter, r *rest.Request) {
	// code := r.PathParam("code")

	// var country *Country
	// if store[code] != nil {
	// 	country = store[code]
	// }

	// if country == nil {
	// 	rest.NotFound(w, r)
	// 	return
	// }
	// w.WriteJson(country)
}

func GetReceipts(w rest.ResponseWriter, r *rest.Request) {
	// var form map[string]interface{}
	// err := r.DecodeJsonPayload(&form)
	// if err != nil {
	//  	rest.Error(w, err.Error(), http.StatusInternalServerError)
	//  	return
	// }

	username := r.PathParam("username")
    msg := query.Query(username)
    w.(http.ResponseWriter).Write([]byte(msg))
}

func PostReceipt(w rest.ResponseWriter, r *rest.Request) {
    var receipt entity.WarehouseReceipt
    err := r.DecodeJsonPayload(&receipt)
    if err != nil{
	    rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
    }
	
    msg := app1.Assign(receipt)
	w.(http.ResponseWriter).Write([]byte(msg))
}

func PutReceipt(w rest.ResponseWriter, r *rest.Request) {
    var receipt entity.WarehouseReceipt
    err := r.DecodeJsonPayload(&receipt)
    if err != nil{
	    rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
    }
	
    msg := app2.Transfer(receipt)
    w.(http.ResponseWriter).Write([]byte(msg))
}

func PostMember(w rest.ResponseWriter, r *rest.Request) {
    var user User
    err := r.DecodeJsonPayload(&user)
    if err != nil{
	    rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
    }
	
    msg := member.Register(user.Username)
	w.(http.ResponseWriter).Write([]byte(msg))
}

type User struct {
    Username               string
}

// POST
// curl -i -H 'Content-Type: application/json' -d '{"ID":"","WARRANTY_NO":"W150514008272","CUSTOMER_CODE":"customer","CUSTOMER_NAME":"上海动产测试贸易商","WAREHOUSE_CODE":"warehouse","WAREHOUSE_NAME":"上海动产测试仓库","WARRANTY_STATUS":"regulator","WARRANTY_TYPE":"10","NUM":"1","UNIT":"件","WEIGHT":"8.68","WEIGHT_UNIT":"吨","PRICE":"34711.32","WARRANTY_APP_NO":"P150514008105","CREATE_TIME":"2015-05-14 14:54:26","CREATE_CODE":"U00122","CREATE_NAME":"张三","TYPE_CODE":"1001","TYPE_NAME":"电解铜","CGLOCATION":"A-1","CWHOPITION":"仓库已核","CINSSUEDPERSON":"李四","CINSSUEDPLACE":"上海","CINSSUEDDATE":"2015-09-14 14:54:26","CPLATFORMOPTION":"平台未核","CREGULEOPTION":"监管公司未核","COWNEROPTION":"存货人未验收"}' http://127.0.0.1:8888/receipts
// curl -i -H 'Content-Type: application/json' -d '{"ID":"","WARRANTY_NO":"W150514008273","CUSTOMER_CODE":"customer","CUSTOMER_NAME":"上海动产测试贸易商","WAREHOUSE_CODE":"warehouse","WAREHOUSE_NAME":"上海动产测试仓库","WARRANTY_STATUS":"regulator","WARRANTY_TYPE":"10","NUM":"1","UNIT":"件","WEIGHT":"8.68","WEIGHT_UNIT":"吨","PRICE":"34711.32","WARRANTY_APP_NO":"P150514008105","CREATE_TIME":"2015-05-14 14:54:26","CREATE_CODE":"U00122","CREATE_NAME":"张三","TYPE_CODE":"1001","TYPE_NAME":"电解铜","CGLOCATION":"A-1","CWHOPITION":"仓库已核","CINSSUEDPERSON":"李四","CINSSUEDPLACE":"上海","CINSSUEDDATE":"2015-09-14 14:54:26","CPLATFORMOPTION":"平台未核","CREGULEOPTION":"监管公司未核","COWNEROPTION":"存货人未验收"}' http://127.0.0.1:8888/receipts
// curl -i -H 'Content-Type: application/json' -d '{"ID":"","WARRANTY_NO":"W150514008274","CUSTOMER_CODE":"customer2","CUSTOMER_NAME":"上海动产测试贸易商2","WAREHOUSE_CODE":"warehouse","WAREHOUSE_NAME":"上海动产测试仓库","WARRANTY_STATUS":"regulator","WARRANTY_TYPE":"10","NUM":"1","UNIT":"件","WEIGHT":"8.68","WEIGHT_UNIT":"吨","PRICE":"34711.32","WARRANTY_APP_NO":"P150514008105","CREATE_TIME":"2015-05-14 14:54:26","CREATE_CODE":"U00122","CREATE_NAME":"张三","TYPE_CODE":"1001","TYPE_NAME":"电解铜","CGLOCATION":"A-1","CWHOPITION":"仓库已核","CINSSUEDPERSON":"李四","CINSSUEDPLACE":"上海","CINSSUEDDATE":"2015-09-14 14:54:26","CPLATFORMOPTION":"平台未核","CREGULEOPTION":"监管公司未核","COWNEROPTION":"存货人未验收"}' http://127.0.0.1:8888/receipts
// curl -i -H 'Content-Type: application/json' -d '{"ID":"","WARRANTY_NO":"W150514008275","CUSTOMER_CODE":"customer","CUSTOMER_NAME":"上海动产测试贸易商","WAREHOUSE_CODE":"warehouse2","WAREHOUSE_NAME":"上海动产测试仓库2","WARRANTY_STATUS":"regulator","WARRANTY_TYPE":"10","NUM":"1","UNIT":"件","WEIGHT":"8.68","WEIGHT_UNIT":"吨","PRICE":"34711.32","WARRANTY_APP_NO":"P150514008105","CREATE_TIME":"2015-05-14 14:54:26","CREATE_CODE":"U00122","CREATE_NAME":"张三","TYPE_CODE":"1001","TYPE_NAME":"电解铜","CGLOCATION":"A-1","CWHOPITION":"仓库已核","CINSSUEDPERSON":"李四","CINSSUEDPLACE":"上海","CINSSUEDDATE":"2015-09-14 14:54:26","CPLATFORMOPTION":"平台未核","CREGULEOPTION":"监管公司未核","COWNEROPTION":"存货人未验收"}' http://127.0.0.1:8888/receipts
// curl -i -H 'Content-Type: application/json' -d '{"ID":"","WARRANTY_NO":"W150514008276","CUSTOMER_CODE":"customer2","CUSTOMER_NAME":"上海动产测试贸易商2","WAREHOUSE_CODE":"warehouse2","WAREHOUSE_NAME":"上海动产测试仓库2","WARRANTY_STATUS":"regulator","WARRANTY_TYPE":"10","NUM":"1","UNIT":"件","WEIGHT":"8.68","WEIGHT_UNIT":"吨","PRICE":"34711.32","WARRANTY_APP_NO":"P150514008105","CREATE_TIME":"2015-05-14 14:54:26","CREATE_CODE":"U00122","CREATE_NAME":"张三","TYPE_CODE":"1001","TYPE_NAME":"电解铜","CGLOCATION":"A-1","CWHOPITION":"仓库已核","CINSSUEDPERSON":"李四","CINSSUEDPLACE":"上海","CINSSUEDDATE":"2015-09-14 14:54:26","CPLATFORMOPTION":"平台未核","CREGULEOPTION":"监管公司未核","COWNEROPTION":"存货人未验收"}' http://127.0.0.1:8888/receipts
// curl -i -H 'Content-Type: application/json' -d '' http://127.0.0.1:8888/receipts
// curl -i -H 'Content-Type: application/json' -d '{"Username":"testuser"}' http://127.0.0.1:8888/members


// GET all
// curl -i http://127.0.0.1:8888/receipts/warehouse
// curl -i http://127.0.0.1:8888/receipts/warehouse2
// curl -i http://127.0.0.1:8888/receipts/customer
// curl -i http://127.0.0.1:8888/receipts/customer2
// curl -i http://127.0.0.1:8888/receipts/ouyeel
// curl -i http://127.0.0.1:8888/receipts/regulator

// PUT
// curl -i -X PUT -H 'Content-Type: application/json' -d '{"ID":"warehouse-1","WARRANTY_NO":"W150514008272","CUSTOMER_CODE":"customer","CUSTOMER_NAME":"上海动产测试贸易商","WAREHOUSE_CODE":"warehouse","WAREHOUSE_NAME":"上海动产测试仓库","WARRANTY_STATUS":"regulator","WARRANTY_TYPE":"10","NUM":"1","UNIT":"件","WEIGHT":"8.68","WEIGHT_UNIT":"吨","PRICE":"34711.32","WARRANTY_APP_NO":"P150514008105","CREATE_TIME":"2015-05-14 14:54:26","CREATE_CODE":"U00122","CREATE_NAME":"张三","TYPE_CODE":"1001","TYPE_NAME":"电解铜","CGLOCATION":"A-1","CWHOPITION":"仓库已核","CINSSUEDPERSON":"李四","CINSSUEDPLACE":"上海","CINSSUEDDATE":"2015-09-14 14:54:26","CPLATFORMOPTION":"平台未核","CREGULEOPTION":"监管公司未核","COWNEROPTION":"存货人未验收"}' http://127.0.0.1:8888/receipts/warehouse
// curl -i -X PUT -H 'Content-Type: application/json' -d '' http://127.0.0.1:8888/receipts/warehouse

// GET one
// curl -i http://127.0.0.1:8888/receipts/FR
