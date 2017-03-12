/*
    author:DASE@ECNU
*/
package entity

import (
        "io/ioutil"
        "strconv"
        "fmt"
)

type User struct{
    Username string
    Password string
    Realname string
}

type WarehouseReceipt struct{
    //主健, 后台生成
    ID string
    //仓单号
    WARRANTY_NO	string
    //持有人代码, 值为:customr, customer2
    CUSTOMER_CODE string	
    //持有人名称
    CUSTOMER_NAME string
    //保管人代码, 值为:warehouse, warehouse2
    WAREHOUSE_CODE string	
    //保管人名称
    WAREHOUSE_NAME string
    //保管人地址
    WAREHOUSE_ADDRESS string	
    //仓单状态: 值为:regulator, ouyeel, customer
    WARRANTY_STATUS	string
    //仓单类型
    WARRANTY_TYPE string
    //总数量
    NUM	string
    //单位
    UNIT string
    //总重量
    WEIGHT string
    //重量单位
    WEIGHT_UNIT string
    //总价值
    PRICE string	
    //有效期
    //VALIDITY_TIME string	
    //版本号
    //HAND_NO string
    //制单申请号
    WARRANTY_APP_NO string	
    //第三方监管标志(0--不需要;1--需要监管)
    //SUPERVISON_FLAG	string
    //参保标记(0-无保险，1-参保)
    //INSURED_FLAG string
    //是否违禁仓单（1-是、0-否）    	
    //ILLICIT_FLAG string
    //创建时间	
    CREATE_TIME string
    //创建人代码
    CREATE_CODE string
    //创建人姓名
    CREATE_NAME string
    //最后修改人代码
    //UPDATE_CODE	string
    //最后修改人姓名
    //UPDATE_NAME string
    //品种代码
    TYPE_CODE string
    //品种名称
    TYPE_NAME string
    //规格
    CSPEC string
    //产地
    CORIGIN string
    //材质
    CMATERIAL string
    //货位
    CGLOCATION string
    //仓库复核, 值为:仓库未核, 仓库已核
    CWHOPITION string
    //填发人
    CINSSUEDPERSON string
    //填发地
    CINSSUEDPLACE string
    //填发日期
    CINSSUEDDATE string
    //平台复核, 值为:平台未核, 平台已核
    CPLATFORMOPTION string
    //监管公司复核, 值为:监管公司未核, 监管公司已核
    CREGULEOPTION string
    //存货人验收仓单
    COWNEROPTION string
}

func NewWarehouseReceipt(filename string) WarehouseReceipt {
    line, _ := ioutil.ReadFile("../config/"+filename)
    numStr := string(line)
    fmt.Println("numStr="+numStr)
    num, _ := strconv.Atoi(numStr)
    ioutil.WriteFile("../config/"+filename, []byte(strconv.Itoa(num+1)), 0644)
    
    return WarehouseReceipt{ID:filename+"-"+numStr,WARRANTY_NO:"W150514008272",CUSTOMER_CODE:"customer",CUSTOMER_NAME:"上海动产测试贸易商",WAREHOUSE_CODE:filename,WAREHOUSE_NAME:"上海动产测试仓库",WAREHOUSE_ADDRESS:"上海市宝山区xx路xx号",
    WARRANTY_STATUS:"regulator",WARRANTY_TYPE:"10",NUM:"1",UNIT:"件",WEIGHT:"8.68",WEIGHT_UNIT:"吨",PRICE:"34711.32",WARRANTY_APP_NO:"P150514008105",CREATE_TIME:"2015-05-14 14:54:26",CREATE_CODE:"U00122",
    CREATE_NAME:"张三",TYPE_CODE:"1001",TYPE_NAME:"电解铜",CSPEC:"Φ100*900",CORIGIN:"xx钢",CMATERIAL:"H234",CGLOCATION:"A-1",CWHOPITION:"仓库已核",CINSSUEDPERSON:"李四",CINSSUEDPLACE:"上海",CINSSUEDDATE:"2015-09-14 14:54:26",
    CPLATFORMOPTION:"平台未核",CREGULEOPTION:"监管公司未核",COWNEROPTION:"存货人未验收"}
}