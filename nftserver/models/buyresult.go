package models

import (
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

func (nft NftDb) BuyResult(from, to, contractAddr, tokenId, trade_sig, price, sig, royalty, txhash string) error {
	from = strings.ToLower(from)
	to = strings.ToLower(to)
	contractAddr = strings.ToLower(contractAddr)

	fmt.Println("BuyResult() price = ", price)
	if IsUint64DataValid(price) != true {
		fmt.Println("BuyResult() price err")
		return ErrPrice
	}
	fmt.Println(time.Now().String()[:25],"BuyResult() Begin", "from=", from, "to=", to, "price=", price,
		"contractAddr=", contractAddr, "tokenId=", tokenId, "txhash=", txhash,
		"royalty=", royalty/*, "sig=", sig, "trade_sig=", trade_sig*/)
	fmt.Println("BuyResult()++q++++++++++++++++++")
	if royalty != "" {
		fmt.Println("BuyResult() royalty!=Null mint royalty=", royalty)
		var nftRec Nfts
		err := nft.db.Where("contract = ? AND tokenid = ?", contractAddr, tokenId).First(&nftRec)
		if err.Error != nil {
			fmt.Println("BuyResult() royalty err =", ErrNftNotExist)
			return ErrNftNotExist
		}
		trans := Trans{}
		trans.Contract = contractAddr
		trans.Fromaddr = ""
		trans.Toaddr = to
		trans.Signdata = sig
		trans.Tradesig = trade_sig
		trans.Tokenid = tokenId
		trans.Price, _ = strconv.ParseUint(price, 10, 64)
		trans.Transtime = time.Now().Unix()
		trans.Selltype = SellTypeMintNft.String()
		trans.Name = nftRec.Name
		trans.Meta = nftRec.Meta
		trans.Desc = nftRec.Desc
		trans.Txhash = txhash
		return nft.db.Transaction(func(tx *gorm.DB) error {
			err := tx.Model(&trans).Create(&trans)
			if err.Error != nil {
				fmt.Println("BuyResult() royalty create trans err=", err.Error)
				return err.Error
			}
			nftrecord := Nfts{}
			nftrecord.Signdata = sig

			nftrecord.Royalty, _ = strconv.Atoi(royalty)
			//nftrecord.Royalty = nftrecord.Royalty / 100
			nftrecord.Mintstate = Minted.String()
			err = tx.Model(&Nfts{}).Where("contract = ? AND tokenid =?",
				contractAddr, tokenId).Updates(&nftrecord)
			if err.Error != nil {
				fmt.Println("BuyResult() royalty update nfts record err=", err.Error)
				return err.Error
			}
			fmt.Println("BuyResult() royalty!=Null Ok")
			return nil
		})
	}
	fmt.Println("BuyResult()-------------------")
	if from != "" && to != "" {
		fmt.Println("BuyResult() from != Null && to != Null" )
		var nftRec Nfts
		err := nft.db.Where("contract = ? AND tokenid = ?", contractAddr, tokenId).First(&nftRec)
		if err.Error != nil {
			fmt.Println("BuyResult() auction not find err=", err.Error)
			return ErrNftNotExist
		}
		if price == "" {
			fmt.Println("BuyResult() price == null" )
			return nft.db.Transaction(func(tx *gorm.DB) error {
				var auctionRec Auction
				err = tx.Set("gorm:query_option", "FOR UPDATE").Where("contract = ? AND tokenid = ? AND ownaddr =?",
					contractAddr, tokenId, nftRec.Ownaddr).First(&auctionRec)
				if err.Error != nil {
					fmt.Println("BuyResult() auction not find err=", err.Error)
					return err.Error
				}
				trans := Trans{}
				trans.Auctionid = auctionRec.ID
				trans.Contract = auctionRec.Contract
				trans.Createaddr = nftRec.Createaddr
				trans.Fromaddr = from
				trans.Toaddr = to
				trans.Signdata = sig
				trans.Tradesig = trade_sig
				trans.Tokenid = auctionRec.Tokenid
				trans.Nftid = auctionRec.Nftid
				trans.Paychan = auctionRec.Paychan
				trans.Currency = auctionRec.Currency
				trans.Price = 0
				trans.Transtime = time.Now().Unix()
				trans.Selltype = SellTypeAsset.String()
				err := tx.Model(&trans).Create(&trans)
				if err.Error != nil {
					fmt.Println("BuyResult() create trans record err=", err.Error)
					return err.Error
				}
				fmt.Println("BuyResult() price == null OK" )
				return nil
			})
		}else{
			fmt.Println("BuyResult() price != null" )
			return nft.db.Transaction(func(tx *gorm.DB) error {
				var auctionRec Auction
				err = tx.Where("contract = ? AND tokenid = ? AND ownaddr =?",
					contractAddr, tokenId, nftRec.Ownaddr).First(&auctionRec)
				if err.Error != nil {
					fmt.Println("BuyResult() auction not find err=", err.Error)
					return err.Error
				}
				trans := Trans{}
				trans.Auctionid = auctionRec.ID
				trans.Contract = auctionRec.Contract
				trans.Createaddr = nftRec.Createaddr
				trans.Fromaddr = from
				trans.Toaddr = to
				trans.Signdata = sig
				trans.Tradesig = trade_sig
				trans.Nftid = auctionRec.Nftid
				trans.Tokenid = auctionRec.Tokenid
				trans.Paychan = auctionRec.Paychan
				trans.Currency = auctionRec.Currency
				trans.Txhash = txhash
				trans.Name = nftRec.Name
				trans.Meta = nftRec.Meta
				trans.Desc = nftRec.Desc
				trans.Price, _ = strconv.ParseUint(price, 10, 64)
				trans.Transtime = time.Now().Unix()
				/*if auctionRec.Selltype == SellTypeWaitSale.String() {
					trans.Selltype = SellTypeHighestBid.String()
				}else {
					trans.Selltype = auctionRec.Selltype
				}*/
				trans.Selltype = auctionRec.Selltype
				err := tx.Model(&trans).Create(&trans)
				if err.Error != nil {
					fmt.Println("BuyResult() create trans record err=", err.Error)
					return err.Error
				}
				var collectRec Collects
				err = nft.db.Where("name = ? AND createaddr =?",
					nftRec.Collections, nftRec.Collectcreator).First(&collectRec)
				if err.Error == nil {
					transCnt := collectRec.Transcnt + 1
					transAmt := collectRec.Transamt + trans.Price
					collectRec = Collects{}
					collectRec.Transcnt = transCnt
					collectRec.Transamt = transAmt
					err = tx.Model(&Collects{}).Where("name = ? AND createaddr =?",
						nftRec.Collections, nftRec.Collectcreator).Updates(&collectRec)
					if err.Error != nil {
						fmt.Println("BuyResult() update collectRec err=", err.Error)
						return err.Error
					}
				}
				nftrecord := Nfts{}
				nftrecord.Ownaddr = to
				nftrecord.Selltype = SellTypeNotSale.String()
				nftrecord.Paychan = auctionRec.Paychan
				nftrecord.TransCur = auctionRec.Currency
				nftrecord.Transprice = trans.Price
				nftrecord.Transamt = nftRec.Transamt + trans.Price
				nftrecord.Transcnt = nftRec.Transcnt + 1
				nftrecord.Transtime = time.Now().Unix()
				err = tx.Model(&Nfts{}).Where("contract = ? AND tokenid =?",
					auctionRec.Contract, auctionRec.Tokenid).Updates(&nftrecord)
				if err.Error != nil {
					fmt.Println("BuyResult() update record err=", err.Error)
					return err.Error
				}
				err = tx.Model(&Auction{}).Where("contract = ? AND tokenid = ?",
					auctionRec.Contract, auctionRec.Tokenid).Delete(&Auction{})
				if err.Error != nil {
					fmt.Println("BuyResult() delete auction record err=", err.Error)
					return err.Error
				}
				err = nft.db.Model(&Bidding{}).Where("contract = ? AND tokenid = ?",
					auctionRec.Contract, auctionRec.Tokenid).Delete(&Bidding{})
				if err.Error != nil {
					fmt.Println("BuyResult() delete bid record err=", err.Error)
					return err.Error
				}
				fmt.Println("BuyResult() from != Null && to != Null --> price != Null OK" )
				return nil
			})
		}
	}
	fmt.Println("BuyResult() End.")
	return ErrFromToAddrZero
}
