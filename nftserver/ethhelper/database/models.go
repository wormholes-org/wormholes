package database

import (
	"bytes"
	"encoding/json"
	"github.com/nftexchange/nftserver/ethhelper/common"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"time"
)

type Collection struct {
	Id           int64  `gorm:"column:id;not null;type:bigint primary key auto_increment"`
	UserAddr     string `gorm:"column:creator;type:varchar(50) " json:"user_addr"`
	CreateTs     string `gorm:"column:create_ts;type:varchar(20) " json:"ts"`
	CreateHash   string `gorm:"column:create_hash;type:varchar(80) " json:"hash"`
	Name         string `gorm:"column:name;type:varchar(50);charset utf8 collate utf8_general_ci NULL " json:"name"`
	Img          string `gorm:"column:img;type:varchar(50)" json:"img"`
	ContractAddr string `gorm:"column:contract_addr;null;type:varchar(50) " json:"contract_addr"`
	Desc         string `gorm:"column:desc;type:varchar(50);charset utf8 collate utf8_general_ci NULL" json:"desc"`
	Categories   string `gorm:"column:categories;type:varchar(50);charset utf8 collate utf8_general_ci NULL " json:"categories"`
	ContractType string `gorm:"column:contract_type;type:varchar(50);charset utf8 collate utf8_general_ci NULL " json:"contract_type"`
}
type Nft struct {
	Id       int64  `gorm:"column:id;not null;type:bigint primary key auto_increment"`
	Uri      string `gorm:"column:uri;type:text" json:"uri"`
	Desc     string `gorm:"column:desc;type:text;charset utf8 collate utf8_general_ci NULL " json:"description"`
	Name     string `gorm:"column:name;type:varchar(300);charset utf8 collate utf8_general_ci NULL " json:"name"`
	Img      string `gorm:"column:img;type:varchar(50)" json:"image"`
	Contract string `gorm:"column:contract;type:varchar(50) " json:"contract"`
	TokenId  string `gorm:"column:tokenId;type:varchar(200) " json:"tokenId"`
}
type NftTx struct {
	Id               int64  `gorm:"column:id;not null;type:bigint primary key auto_increment"`
	Value            string `gorm:"column:value;type:varchar(50) " json:"value"`
	From             string `gorm:"column:fromAddr;type:varchar(50) " json:"from"`
	To               string `gorm:"column:toAddr;type:varchar(50) " json:"to"`
	TxHash           string `gorm:"column:txHash;type:varchar(80) " json:"txHash"`
	Ts               string `gorm:"column:ts;type:varchar(20) " json:"ts"`
	Contract         string `gorm:"column:contract;type:varchar(50) " json:"contract"`
	TokenId          string `gorm:"column:tokenId;type:varchar(200) " json:"tokenId"`
	BlockNumber      string `gorm:"column:blockNumber;type:varchar(50) " json:"blockNumber"`
	TransactionIndex string `gorm:"column:transactionIndex;type:varchar(30) " json:"transactionIndex"`
}


func (c *Collection) Insert() {
	db.Insert(c)
}
func (c *Collection) GetByContract(addr string) *Collection {
	var col Collection
	if db.Db().Table("collections").Select("*").Where("contract_addr=?", addr).First(&col).Error == nil {
		return &col
	}
	return nil
}
func (c *Collection) Update() {
	db.Update(c)
}

// NftCollectionsExist 合集图片等同步完成后写程序来循环修改
func (n *Nft) NftCollectionsExist() bool {
	if db.Db().Table("collections").Select("Id").Where("contract_addr=?", n.Contract).Where("img=''").Error == nil {
		err := db.Db().Exec("update collections set img='" + n.Img + "' where contract='" + n.Contract + "'").Error
		if err != nil {
			log.Println("NftImgExist err:", err)
		}
	}
	return true
}
func (n *NftTx) Insert() {
	db.Insert(n)
}

func (n *Nft) Insert() {
	db.Insert(n)
}

type Col struct {
	User string ` json:"user_addr"`
	Name string ` json:"name"`
	Img  string ` json:"img"`
	Type string ` json:"contract_type"`
	Addr string ` json:"contract_addr"`
	Desc string ` json:"desc"`
	Cate string ` json:"categories"`
	Sig  string ` json:"sig"`
}

type NftModel struct {
	Uri               string ` json:"meta"`
	Desc              string ` json:"desc"`
	Name              string ` json:"name"`
	User_addr         string ` json:"user_addr"`
	Creator_addr      string ` json:"creator_addr"`
	Owner_addr        string ` json:"owner_addr"`
	Md5               string ` json:"md5"`
	Source_url        string ` json:"source_url"`
	Nft_contract_addr string ` json:"nft_contract_addr"`
	Nft_token_id      string ` json:"nft_token_id"`
	Categories        string ` json:"categories"`
	Collections       string ` json:"collections"`
	Asset_sample      string ` json:"asset_sample"`
	Hide              string ` json:"hide"`
	Royalty           string ` json:"royalty"`
	Count             string ` json:"count"`
	Sig               string ` json:"sig"`
}
type TxResult struct {
	Admin_addr        string ` json:"admin_addr"`
	From              string ` json:"from"`
	To                string ` json:"to"`
	Nft_contract_addr string ` json:"nft_contract_addr"`
	Nft_token_id      string ` json:"nft_token_id"`
	Trade_sig         string ` json:"trade_sig"`
	Price             string ` json:"price"`
	Royalty           string ` json:"royalty"`
	Txhash            string ` json:"txhash"`
	Sig               string ` json:"sig"`
	Admin_sig         string ` json:"admin_sig"`
}

const contract = "0x1b034413634915361d354fd4f9b96c35aa8e6094"

func GetMaxNfts() {
	var tt []Nft
	err := db.Db().Table("nfts").Select("*").
		Where("contract=?", contract).Scan(&tt).Error
	var colName string
	var colAddr string
	var ttt big.Int
	if err == nil {
		for i := 0; i < len(tt); i++ {
			//if i == 20 {
			//	break
			//}
			t := tt[i]
			var nftM NftModel
			base64Data, err := common.GetImgBase64(t.Img)
			if i == 0 {
				var col Collection
				var c Col
				err = db.Db().Table("collections").Select("*").
					Where("contract_addr=?", contract).First(&col).Error

				c.Addr = col.UserAddr
				colAddr = col.UserAddr
				c.Type = "ERC721"
				c.Name = col.Name + " T"
				c.Desc = col.Desc
				c.User = col.UserAddr
				c.Sig = "0x4a71940655b075316ae19b02457201ed0f719d14f2d20c986b8c16113233e047535d5d1cc4eb293609e79bc60daf622216b190d50a16519d6f826bee05e548051b"
				c.Img = base64Data

				Post(c, new_collection_api)
				colName = col.Name + " T"
			}
			nftM.Asset_sample = base64Data
			nftM.Uri = t.Uri
			nftM.Desc = t.Desc
			nftM.Name = t.Name
			nftM.User_addr = "0x7fbc8ad616177c6519228fca4a7d9ec7d1804900"
			nftM.Owner_addr = colAddr
			nftM.Creator_addr = colAddr
			nftM.Nft_contract_addr = t.Contract
			nftM.Nft_token_id = t.TokenId + "33"
			nftM.Source_url = t.Img
			nftM.Hide = "false"
			nftM.Royalty = "0"
			nftM.Count = "1"
			nftM.Collections = colName
			nftM.Sig = "0x4a71940655b075316ae19b02457201ed0f719d14f2d20c986b8c16113233e047535d5d1cc4eb293609e79bc60daf622216b190d50a16519d6f826bee05e548051b"
			_, err = Post(nftM, new_nft_api)
			if err == nil {
				var txs []NftTx
				err = db.Db().Table("nft_txes").Select("*").
					Where("contract=?", contract).Where("tokenId=?", t.TokenId).Order("id ASC").Scan(&txs).Error
				if err == nil {
					var tx TxResult
					for j := 0; j < len(txs); j++ {
						tx.From = txs[j].From
						tx.To = txs[j].To
						tx.Nft_contract_addr = nftM.Nft_contract_addr
						tx.Nft_token_id = nftM.Nft_token_id
						tx.Royalty = ""
						if tx.From == "" {
							tx.Royalty = "2"
						}
						ttt.SetString(txs[j].Value, 0)
						ttt.Div(&ttt, new(big.Int).SetUint64(1000000000))
						fmt.Println(ttt.String())
						tx.Price = ttt.String()
						tx.Txhash = txs[j].TxHash
						tx.Admin_addr = "0x7fbc8ad616177c6519228fca4a7d9ec7d1804900"
						tx.Admin_sig = "0x4a71940655b075316ae19b02457201ed0f719d14f2d20c986b8c16113233e047535d5d1cc4eb293609e79bc60daf622216b190d50a16519d6f826bee05e548051b"
						tx.Sig = "0x4a71940655b075316ae19b02457201ed0f719d14f2d20c986b8c16113233e047535d5d1cc4eb293609e79bc60daf622216b190d50a16519d6f826bee05e548051b"
						tx.Trade_sig = "0x4a71940655b075316ae19b02457201ed0f719d14f2d20c986b8c16113233e047535d5d1cc4eb293609e79bc60daf622216b190d50a16519d6f826bee05e548051b"
						_, err = Post(tx, new_tx_api)
					}
				}
			}
		}
	}
}

const postUrl = "http://192.168.1.238:8081/v2/"
const new_collection_api = "newCollections"
const new_nft_api = "upload"
const new_tx_api = "buyResultInterface"

func Post(data interface{}, api string) (string, error) {
	contentType := "application/json"
	client := &http.Client{Timeout: 10 * time.Second}
	jsonStr, _ := json.Marshal(data)
	resp, err := client.Post(postUrl+api, contentType, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Println("Post "+api+"  :", err)
		return "", err
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(result))
	return string(result), nil
}

func GetTopNft() {
	type Tmp struct {
		Tm string `gorm:"column:count;" json:"count"`
	}

	go func() {
		var t Tmp
		for i := 1; i < 5000; i++ {
			err := db.Db().Table("collections").Select("contract_addr as count").Where("id=?", fmt.Sprintf("%v", i)).First(&t).Error
			if err == nil {
				//var t1 Tmp
				//err = db.Db().Table("nft_txes").Select("count(*) as count").
				//	Where("contract=?", fmt.Sprintf("%v", t.Tm)).Where("value!='0'").
				//	Where("value is not NULL").Where(" LENGTH(trim(value))>0").First(&t1).Error
				//if v, _ := strconv.Atoi(t1.Tm); v > 20 {
				var t2 Tmp
				err = db.Db().Table("nfts").Select("count(*) as count").
					Where("contract=?", fmt.Sprintf("%v", t.Tm)).Where(" LENGTH(trim(img))>0").Where("img is not NULL ").First(&t2).Error
				if v2, _ := strconv.Atoi(t2.Tm); v2 < 100 && v2 > 50 {
					fmt.Println(t.Tm)
				}
			}
			//}
			if i%1000 == 0 {
				fmt.Println("_______________", i)
			}
		}
	}()
	go func() {
		var t Tmp
		for i := 5000; i < 10000; i++ {
			err := db.Db().Table("collections").Select("contract_addr as count").Where("id=?", fmt.Sprintf("%v", i)).First(&t).Error
			if err == nil {
				//var t1 Tmp
				//err = db.Db().Table("nft_txes").Select("count(*) as count").
				//	Where("contract=?", fmt.Sprintf("%v", t.Tm)).Where("value!='0'").Where(" LENGTH(trim(value))>0").Where("value is not NULL").First(&t1).Error
				//if v, _ := strconv.Atoi(t1.Tm); v > 20 {
				var t2 Tmp
				err = db.Db().Table("nfts").Select("count(*) as count").
					Where("contract=?", fmt.Sprintf("%v", t.Tm)).Where(" LENGTH(trim(img))>0").Where("img is not NULL ").First(&t2).Error
				if v2, _ := strconv.Atoi(t2.Tm); v2 < 100 && v2 > 50 {
					fmt.Println(t.Tm)
				}
				//}
			}
			if i%1000 == 0 {
				fmt.Println("_______________", i)
			}
		}
	}()
	go func() {
		var t Tmp
		for i := 10000; i < 16900; i++ {
			err := db.Db().Table("collections").Select("contract_addr as count").Where("id=?", fmt.Sprintf("%v", i)).First(&t).Error
			if err == nil {
				//var t1 Tmp
				//err = db.Db().Table("nft_txes").Select("count(*) as count").
				//	Where("contract=?", fmt.Sprintf("%v", t.Tm)).Where("value!='0'").Where(" LENGTH(trim(value))>0").Where("value is not NULL").First(&t1).Error
				//if v, _ := strconv.Atoi(t1.Tm); v > 20 {
				var t2 Tmp
				err = db.Db().Table("nfts").Select("count(*) as count").
					Where("contract=?", fmt.Sprintf("%v", t.Tm)).Where(" LENGTH(trim(img))>0").Where("img is not NULL ").First(&t2).Error
				if v2, _ := strconv.Atoi(t2.Tm); v2 < 100 && v2 > 50 {
					fmt.Println(t.Tm)
				}
				//}
			}
			if i%1000 == 0 {
				fmt.Println("_______________", i)
			}
		}
	}()

}
