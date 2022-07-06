package database

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
	"time"
)

type MysqlDb struct {
	db *gorm.DB
}
var db *MysqlDb
type Config struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
	MaxIdle  int    `json:"maxIdle"`
	MaxOpen  int    `json:"maxOpen"`
}

func init() {
	//db = NewMysqlDb()
	//CreateTable()
	//dbExecInitSql()
}

func CreateTable() {
	db.Db().AutoMigrate(&Collection{})
	db.Db().AutoMigrate(&Nft{})
	db.Db().AutoMigrate(&NftTx{})
}
func dbExecInitSql() {
	db.Db().Exec("ALTER  table nfts add  unique(contract,tokenId)")
	db.Db().Exec("ALTER  table collections add  unique(contract_addr)")
	db.Db().Exec("ALTER  table nft_txes add  unique(fromAddr,toAddr,txHash,contract,tokenId)")
}
func NewMysqlDb() *MysqlDb {
	type DatabaseConfig struct {
		Conf Config `json:"mysqlConfig"`
	}
	var config = &DatabaseConfig{}
	filePtr, err := os.Open("./config.json")
	defer filePtr.Close()
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(config)
	db, err := gorm.Open("mysql", config.Conf.User+":"+config.Conf.Password+"@tcp("+config.Conf.Host+":"+config.Conf.Port+")/"+config.Conf.Database+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
		fmt.Println(err.Error())
	} else {
		db.DB().SetMaxIdleConns(config.Conf.MaxIdle)  //最大空闲连接数
		db.DB().SetMaxOpenConns(config.Conf.MaxOpen)  //最大连接数
		db.DB().SetConnMaxLifetime(time.Second * 300) //设置连接空闲超时
		return &MysqlDb{db: db}
	}
	return nil
}

func (mysql *MysqlDb) Close() error { return mysql.db.Close() }
func (mysql *MysqlDb) Db() *gorm.DB { return mysql.db }
func (mysql *MysqlDb) Status() {
	data, _ := json.Marshal(mysql.Db().DB().Stats()) //获得当前的SQL配置情况
	fmt.Println(string(data))
}

func (mysql MysqlDb) CreateTable(value interface{}) {
	if !mysql.db.HasTable(value) {
		mysql.db.CreateTable(value)
	}
}

func (mysql MysqlDb) Insert(value interface{}) {
	mysql.db.Create(value)
}
func (mysql MysqlDb) Delete(value interface{}) {
	mysql.db.Delete(value)
}
func (mysql MysqlDb) Update(value interface{}) {
	mysql.db.Update(value)
}
