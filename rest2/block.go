package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strconv"
)

func getHeight(ip string, port string)(int, error) {
	resp, err := http.Get("http://"+ ip +":"+ port + "/chain")
	if err != nil {
	    fmt.Printf("%s", err)
	    return 0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body) 
	if err != nil { 
	    fmt.Printf("%s", err)
	    return 0, err
	} 

	var chain map[string]interface{}
	err = json.Unmarshal(body, &chain)
	if err != nil{
		fmt.Println("Error unmarshalling body [%s]", err)
		return 0, err
	}
	return int(chain["height"].(float64)), nil
} 

func containsTxid(ip string, port string, height int, txid string)(bool, error){
	resp, err := http.Get("http://"+ ip +":"+ port + "/chain/blocks/"+ strconv.Itoa(height))
	if err != nil {
	    fmt.Printf("%s", err)
	    return false, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body) 
	if err != nil { 
	    fmt.Printf("%s", err)
	    return false, err
	} 

	var block map[string]interface{}
	err = json.Unmarshal(body, &block)
	if err != nil{
		fmt.Println("Error unmarshalling body [%s]", err)
		return false, err
	}
	
	if block["transactions"] ==nil {
		return false, nil
    }
	transactions := block["transactions"].([]interface{})
	for _, tx := range transactions {
		txMap := tx.(map[string]interface{})
		if(txMap["txid"]==txid){
		    return true, nil	
		}
	}
	return false, nil
} 

func getBlockNumber(ip string, port string, txid string)(int){
    height, _ := getHeight(ip, port)
    for i:=height-1;i>0;i-- {
    	contains, _ := containsTxid(ip, port, i, txid)
    	if contains {
    		return i
    	} 
    }
    return -1  
}

// func main() {
// 	fmt.Println(strconv.Itoa(getBlockNumber("120.132.12.128", "7050", "1dff7d2a-7638-4d09-80fa-a7c233426aa2")))
// }