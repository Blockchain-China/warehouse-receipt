package util

import (
    "github.com/syndtr/goleveldb/leveldb"
    "github.com/syndtr/goleveldb/leveldb/util"
    "fmt"
    "strings"
    "strconv"
)

func Put(key string, value string){
    //var db *leveldb.DB
    //var err error
    db, err := leveldb.OpenFile("../db", nil)
    defer db.Close()
    if err != nil {
        fmt.Printf("LevelDB OpenFile Error: [%s]", err)
        return
    }

    err = db.Put([]byte(key), []byte(value), nil)
    if (err != nil) {
        fmt.Printf("LevelDB Put Error: [%s]", err)
        return
    }
}

func Put2(key string, value string){
    db, err := leveldb.OpenFile("../db", nil)
    defer db.Close()
    if err != nil {
        fmt.Printf("LevelDB OpenFile Error: [%s]", err)
        return
    }
    
    var data []byte
    data, err = db.Get([]byte(key), nil)
    if err != nil {
        if err.Error() == "leveldb: not found" {
            data = []byte("")
        } else {
            fmt.Printf("LevelDB Get Error: [%s]", err)
            return
        }
    }

    err = db.Put([]byte(key), []byte(value+","+string(data)), nil)
    if (err != nil) {
        fmt.Printf("LevelDB Put Error: [%s]", err)
        return
    }

}

func Get(key string)(string){
    db, err := leveldb.OpenFile("../db", nil)
    defer db.Close()
    if err != nil {
        fmt.Printf("LevelDB OpenFile Error: [%s]", err)
        return ""
    }

   var data []byte
    data, err = db.Get([]byte(key), nil)
    if err != nil {
        if err.Error() == "leveldb: not found" {
            data = []byte("")
        } else {
            fmt.Printf("LevelDB Get Error: [%s]", err)
            return ""
        }
    }
    return string(data)
}

func GetSum(key string)(string){
    db, err := leveldb.OpenFile("../db", nil)
    defer db.Close()
    if err != nil {
        fmt.Printf("LevelDB OpenFile Error: [%s]", err)
        return "0"
    }

    var data []byte
    data, err = db.Get([]byte(key), nil)
    if err != nil {
        if err.Error() == "leveldb: not found" {
            data = []byte("0")
        } else {
            fmt.Printf("LevelDB Get Error: [%s]", err)
            return "0"
        }
    }

    arr := strings.Split(string(data), ",")
    sum := len(arr)-1
    if sum < 0 {
        sum = 0
    }
    return strconv.Itoa(sum)
}

func GetByPrefix(key string)(string){
    db, err := leveldb.OpenFile("../db", nil)
    defer db.Close()
    if err != nil {
        fmt.Printf("LevelDB OpenFile Error: [%s]", err)
        return ""
    }

    var str string
    iter := db.NewIterator(util.BytesPrefix([]byte(key)), nil)
    for iter.Next() {
        str =  string(iter.Value())+","+str
    }
    iter.Release()
    if str != "" {
        str = "["+str[:len(str)-1]+"]"
    }
    return str
}

func DelValue(key string, value string){
    db, err := leveldb.OpenFile("../db", nil)
    defer db.Close()
    if err != nil {
        fmt.Printf("LevelDB OpenFile Error: [%s]", err)
        return 
    }

    var data []byte
    data, err = db.Get([]byte(key), nil)
    if err != nil {
        if err.Error() == "leveldb: not found" {
            data = []byte("")
        } else {
            fmt.Printf("LevelDB Put Error: [%s]", err)
            return 
        }
    }
    
    newValue := strings.Replace(string(data), value+",", "", -1)
    err = db.Put([]byte(key), []byte(newValue), nil)
    if (err != nil) {
        fmt.Printf("LevelDB Put Error: [%s]", err)
        return
    }
}

func PutMember(key string, value string){
    db, err := leveldb.OpenFile("../mdb", nil)
    defer db.Close()
    if err != nil {
        fmt.Printf("LevelDB OpenFile Error: [%s]", err)
        return
    }

    err = db.Put([]byte(key), []byte(value), nil)
    if (err != nil) {
        fmt.Printf("LevelDB Put Error: [%s]", err)
        return
    }
}

func GetMember(key string)(string){
    db, err := leveldb.OpenFile("../mdb", nil)
    defer db.Close()
    if err != nil {
        fmt.Printf("LevelDB OpenFile Error: [%s]", err)
        return ""
    }

   var data []byte
    data, err = db.Get([]byte(key), nil)
    if err != nil {
        if err.Error() == "leveldb: not found" {
            data = []byte("")
        } else {
            fmt.Printf("LevelDB Get Error: [%s]", err)
            return ""
        }
    }
    return string(data)
}