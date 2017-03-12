package main

import (
	"fmt"
	"os"
	"bufio"
	"io"
	"strings"
	"time"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"strconv"
)

func getTransactions(id string)([]map[string]string){
	var txArr []map[string]string
    file, err := os.Open("../config/blockchain")
    if err!= nil {
	    fmt.Printf("Error openFile [%s]", err)
	    return nil
    }
    defer file.Close()
    reader := bufio.NewReader(file)
    for {
	    line, err := reader.ReadString('\n') //每次读取一行
	    if err!= nil {
		    if err == io.EOF {
                break
            }
            fmt.Printf("Error readFile [%s]", err)
            return nil
	    }
	    line = strings.Replace(line, "\n", "", -1)
	    arr := strings.Split(line, ",")
	    if arr[0] == id {
            txMap := map[string]string{
	            "user" : arr[1],
	            "txid" : arr[2],
            }
            txArr = append(txArr, txMap)
	    }
    }
    return txArr
}

func getTimestamp(ip string, port string, txid string)(string) {
	resp, err := http.Get("http://"+ ip +":"+ port + "/transactions/"+txid)
	if err != nil {
	    fmt.Printf("%s", err)
	    return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body) 
	if err != nil { 
	    fmt.Printf("%s", err)
	    return ""
	} 

	var tx map[string]interface{}
	err = json.Unmarshal(body, &tx)
	if err != nil{
		fmt.Println("Error unmarshalling body [%s]", err)
		return ""
	}
    
    if tx["timestamp"]==nil {
    	return ""
    }
	ts := tx["timestamp"].(map[string]interface{})
    return time.Unix(int64(ts["seconds"].(float64)), int64(ts["nanos"].(float64))).Format("2006-01-02 15:04:05")
} 

func getTransactions2(ip string, port string, id string)([]map[string]string) {
    var arr2 []map[string]string
    arr := getTransactions(id)
    for _, m := range arr {
    	num := getBlockNumber(ip, port,  m["txid"])
    	if num>0 {
    		tx := map[string]string{"user":m["user"], "txid":m["txid"], "blockNumber":strconv.Itoa(num), "time":getTimestamp(ip, port, m["txid"])}
    		arr2 = append(arr2, tx)
    	}
	}
	return arr2
}

//func main() {
	//getTransactions("warehouse-1")
	//fmt.Println(getTimestamp("120.132.12.128", "7050", "1dff7d2a-7638-4d09-80fa-a7c233426aa2"))
//}