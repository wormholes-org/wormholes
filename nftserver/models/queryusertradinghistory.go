package models

import (
	"fmt"
	"strconv"
	"strings"
)

type TradingHistory struct {
	NftContractAddr string `json:"nft_contract_addr"`
	NftTokenId      string `json:"nft_token_id"`
	NftName         string `json:"nft_name"`
	Price           uint64 `json:"price"`
	Count           uint64 `json:"count"`
	From            string `json:"from"`
	To              string `json:"to"`
	Txhash 			string `json:"trade_hash"`
	Selltype        string `json:"selltype"`
	Date	        int64  `json:"date"`
}

func (nft NftDb) QueryUserTradingHistory(user_addr , start_index, count string) ([]TradingHistory, int, error) {
	user_addr = strings.ToLower(user_addr)
	if IsIntDataValid(start_index) != true {
		return nil, 0, ErrDataFormat
	}
	if IsIntDataValid(count) != true {
		return nil, 0, ErrDataFormat
	}
	var tranRecs []Trans
	var recCount int64
	db := nft.db.Model(Trans{}).Where("(toaddr = ? OR fromaddr = ?) AND (selltype != ? AND selltype != ? AND price > 0)",
		user_addr, user_addr, SellTypeError.String(), SellTypeMintNft.String()).Count(&recCount)
	if db.Error != nil {
		fmt.Println("QueryUserTradingHistory() recCount err=", db)
		return nil, 0, ErrNoTrans
	}
	if recCount == 0 {
		fmt.Println("QueryUserTradingHistory() recCount == 0")
		return nil, 0, ErrNoTrans
	}

	startIndex, _ := strconv.Atoi(start_index)
	nftCount, _ := strconv.Atoi(count)
	if int64(startIndex) > recCount || recCount == 0{
		return nil, 0, ErrNftNoMore
	} else {
		temp := recCount - int64(startIndex)
		if int64(nftCount) > temp {
			nftCount = int(temp)
		}
		err := db.Model(Trans{}).Order("transtime desc").Limit(nftCount).Offset(startIndex).Find(&tranRecs)
		if err.Error != nil {
			fmt.Println("QueryUserTradingHistory() find record err=", err)
			return nil, 0, ErrNftNotExist
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
			tran.Txhash = tranRecs[i].Txhash
			trans = append(trans, tran)
		}
		return trans, int(recCount), nil
	}
}

