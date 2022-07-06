package models

import (
	"strings"
	"gorm.io/gorm"
)

type CollectionInfo struct{
	Createaddr		string	`json:"collection_creator_addr"`
	Contract		string	`json:"nft_contract_addr"`
	Contracttype	string	`json:"contracttype"`
	Name			string	`json:"name"`
	Desc			string	`json:"desc"`
	Categories		string	`json:"categories"`
	Totalcount		int		`json:"totalcount"`
	Tradeamount		uint64	`json:"trade_amount"`
	Tradeavgprice	uint64	`json:"trade_avg_price"`
	Tradefloorprice	uint64	`json:"trade_floor_price"`
	Extend			string	`json:"extend"`
}

func (nft * NftDb) QueryCollectionInfo(creatorAddr string,
	collectionName string) ([]CollectionInfo, error) {

	creatorAddr = strings.ToLower(creatorAddr)
	collection := Collects{}
	collectionInfo := CollectionInfo{}
	collectionInfos := []CollectionInfo{}

	result := nft.db.Model(&Collects{}).Where("createaddr = ? and name = ?",
		creatorAddr, collectionName).First(&collection)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}

	if result.Error == nil {
		collectionInfo.Createaddr = collection.Createaddr
		collectionInfo.Contract	= collection.Contract
		collectionInfo.Contracttype	= collection.Contracttype
		collectionInfo.Name = collection.Name
		collectionInfo.Desc	= collection.Desc
		collectionInfo.Categories = collection.Categories
		collectionInfo.Totalcount = collection.Totalcount
		collectionInfo.Extend = collection.Extend
	}


	type CollectionAmt struct {
		Tradeamount uint64
		Tradecnt uint64
	}
	collectionAmt := CollectionAmt{}
	nftSql := "SELECT SUM(transamt) AS Tradeamount, SUM(transcnt) AS Tradecnt FROM nfts " +
		"WHERE collectcreator = ? AND collections = ? AND deleted_at IS NULL"
	result = nft.db.Raw(nftSql, creatorAddr, collectionName).Scan(&collectionAmt)
	if result.Error == nil {
		collectionInfo.Tradeamount = collectionAmt.Tradeamount
		if collectionAmt.Tradecnt > 0 {
			collectionInfo.Tradeavgprice = collectionAmt.Tradeamount / collectionAmt.Tradecnt
		}
	}

	transSql := "SELECT MIN(trans.price) AS Tradefloorprice " +
		"FROM trans LEFT JOIN nfts ON trans.contract = nfts.contract AND trans.tokenid = nfts.tokenid " +
		"where nfts.collectcreator = ? AND nfts.collections = ? " +
		"AND nfts.deleted_at IS NULL AND trans.selltype != ? AND trans.selltype != ?"
	result = nft.db.Raw(transSql, creatorAddr, collectionName, SellTypeError.String(), SellTypeMintNft.String()).Scan(&collectionInfo.Tradefloorprice)

	collectionInfos = append(collectionInfos, collectionInfo)
	return collectionInfos, nil
}
