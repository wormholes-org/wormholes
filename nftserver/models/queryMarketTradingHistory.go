package models

import (
	"gorm.io/gorm"
	"strings"
)

var NftFields = "Ownaddr, Md5, Name, Desc, Meta, Url, Contract, Tokenid, Count, " +
	"Ownaddr, Md5, Name, Desc, Meta, Url, Contract, Tokenid, Count, " +
	"Approve, Categories, Collectcreator, Collections, Image, Hide, " +
	"Signdata, Createaddr, Verifyaddr, Currency, Price, Royalty, " +
	"Paychan, TransCur, Transprice, Transtime, Createdate, Favorited, " +
	"Transcnt, Transamt, Verified, Verifiedtime, Selltype, Mintstate, " +
	"Extend"

var TransFields = "auctioni, contract, createaddr, " +
	"fromaddr, toaddr, tradesig, signdata, txhash, tokenid, " +
	"count, transtime, paychan, currency, price, name, desc, " +
	"meta, selltype, error"

func (nft NftDb) QueryMarketTradingHistory(filter []StQueryField, sort []StSortField,
	start_index string, count string) ([]TradingHistory, int, error) {
	var tranRecs []Trans
	var recCount int64
	var queryWhere string
	var orderBy string

	sql := "SELECT trans.* FROM trans LEFT JOIN nfts ON trans.contract = nfts.contract AND trans.tokenid = nfts.tokenid"
	countSql := "SELECT count(*) FROM trans LEFT JOIN nfts ON trans.contract = nfts.contract AND trans.tokenid = nfts.tokenid"

	if len(filter) > 0 {
		for k, v := range filter {
			if strings.Contains(TransFields, v.Field) {
				filter[k].Field = "trans." +filter[k].Field
			} else if strings.Contains(NftFields, v.Field) {
				filter[k].Field = "nfts" + filter[k].Field
			}
		}
		queryWhere = nft.joinFilters(filter)
		if len(queryWhere) > 0 {
			sql = sql + " where trans.deleted_at is null and trans.price > 0 and trans.selltype != '" + SellTypeError.String() + "' AND trans.selltype != '" + SellTypeMintNft.String() + "' and" + queryWhere
			countSql = countSql + " where trans.deleted_at is null and trans.price > 0 and trans.selltype != '" + SellTypeError.String() + "' AND trans.selltype != '" + SellTypeMintNft.String() + "' and" + queryWhere
		} else {
			sql = sql + " where trans.deleted_at is null and trans.price > 0 and trans.selltype != '" + SellTypeError.String() + "' AND trans.selltype != '" + SellTypeMintNft.String() + "' "
			countSql = countSql + " where trans.deleted_at is null and trans.price > 0 and trans.selltype != '" + SellTypeError.String() + "' AND trans.selltype != '" + SellTypeMintNft.String() + "' "
		}
	} else {
		sql = sql + " where trans.deleted_at is null and trans.price > 0 and trans.selltype != '" + SellTypeError.String() + "' AND trans.selltype != '" + SellTypeMintNft.String() + "' "
		countSql = countSql + " where trans.deleted_at is null and trans.price > 0 and trans.selltype != '" + SellTypeError.String() + "' AND trans.selltype != '" + SellTypeMintNft.String() + "' "
	}
	if len(sort) > 0 {
		for k, v := range sort {
			if k >0 {
				orderBy = orderBy + ", "
			}
			if strings.Contains(TransFields, v.By) {
				orderBy += "trans." + v.By + " " + v.Order
			} else if strings.Contains(NftFields, v.By) {
				orderBy += "nfts." + v.By + " " + v.Order
			}
		}
	}
	if len(orderBy) > 0 {
		orderBy = orderBy + ", trans.id desc"
	} else {
		orderBy = "trans.id desc"
	}
	sql = sql + " order by " + orderBy
	countSql = countSql + " order by " + orderBy

	if len(start_index) > 0 {
		sql = sql + " limit " + start_index + ", " + count
	}

	err := nft.db.Raw(sql).Scan(&tranRecs)
	if err.Error != nil && err.Error != gorm.ErrRecordNotFound{
		return nil, 0, err.Error
	}
	err = nft.db.Raw(countSql).Scan(&recCount)
	if err.Error != nil {
		return nil, 0, err.Error
	}


	trans := make([]TradingHistory, 0, 20)
	for i := 0; i < len(tranRecs); i++ {
		var tran TradingHistory
		tran.NftContractAddr = tranRecs[i].Contract
		tran.NftTokenId = tranRecs[i].Tokenid
		tran.NftName = tranRecs[i].Name
		tran.Price = tranRecs[i].Price
		tran.Count = 1
		tran.From = tranRecs[i].Fromaddr
		tran.To = tranRecs[i].Toaddr
		tran.Date = tranRecs[i].Transtime
		tran.Selltype = tranRecs[i].Selltype
		tran.Txhash =  tranRecs[i].Txhash
		trans = append(trans, tran)
	}
	return trans, int(recCount), nil
}

