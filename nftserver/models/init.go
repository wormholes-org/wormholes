package models

import (
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
)

var sqldsn string
var sqldsndb string
var Sqldsndb string
var sqllocaldsndb string

func init() {
	dbName, _ = beego.AppConfig.String("dbname")
	dbUserName, _ := beego.AppConfig.String("dbusername")
	dbUserPassword, _ := beego.AppConfig.String("dbuserpassword")
	dbServerIP, _ := beego.AppConfig.String("dbserverip")
	dbServerPort, _ := beego.AppConfig.String("dbserverport")
	//const SqlSvr = "admin:user123456@tcp(192.168.1.238:3306)/"
	SqlSvr = dbUserName + ":" + dbUserPassword + "@tcp(" + dbServerIP + ":" + dbServerPort + ")/"
	fmt.Println("SqlSvr=", SqlSvr)
	sqldsn = SqlSvr + localtime
	sqldsndb = SqlSvr + dbName + localtime
	sqllocaldsndb = SqlSvr + dbName + localtime
	Sqldsndb = sqldsndb
	TradeAddr, _ = beego.AppConfig.String("TradeAddr")
	NFT1155Addr, _ = beego.AppConfig.String("NFT1155Addr")
	AdminAddr, _ = beego.AppConfig.String("AdminAddr")
	EthersNode, _ = beego.AppConfig.String("EthersNode")
	EthersWsNode, _ = beego.AppConfig.String("EthersWsNode")
}