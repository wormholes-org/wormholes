package models

import (
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

func (nft NftDb) BuyResultInterface(admin_addr, from, to, contractAddr, tokenId, trade_sig, price, royalty, txhash, sig, admin_sig string) error {
	from = strings.ToLower(from)
	to = strings.ToLower(to)
	contractAddr = strings.ToLower(contractAddr)

	if IsUint64DataValid(price) != true {
		return ErrPrice
	}
	if admin_addr == "" {
		return ErrDataFormat
	}
	IsAdminAddr, err := IsAdminAddr(admin_addr)
	if err != nil || !IsAdminAddr {
		fmt.Println("BuyResultInterface() admin_addr is not admin.")
		return ErrNftUpAddrNotAdmin
	}

	fmt.Println(time.Now().String()[:25],"BuyResultInterface() Begin", "from=", from, "to=", to, "price=", price,
		"contractAddr=", contractAddr, "tokenId=", tokenId,
		"royalty=", royalty/*, "sig=", sig, "trade_sig=", trade_sig*/)
	fmt.Println("BuyResultInterface()++q++++++++++++++++++")
	if royalty != "" {
		fmt.Println("BuyResultInterface() royalty!=Null mint royalty=", royalty)
		var nftRec Nfts
		err := nft.db.Where("contract = ? AND tokenid = ?", contractAddr, tokenId).First(&nftRec)
		if err.Error != nil {
			fmt.Println("BuyResultInterface() royalty err =", ErrNftNotExist)
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
				fmt.Println("BuyResultInterface() royalty create trans err=", err.Error)
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
				fmt.Println("BuyResultInterface() royalty update nfts record err=", err.Error)
				return err.Error
			}
			fmt.Println("BuyResultInterface() royalty!=Null Ok")
			return nil
		})
	}
	fmt.Println("BuyResultInterface()-------------------")
	if from != "" && to != "" {
		fmt.Println("BuyResultInterface() from != Null && to != null" )
		var nftRec Nfts
		err := nft.db.Where("contract = ? AND tokenid = ?", contractAddr, tokenId).First(&nftRec)
		if err.Error != nil {
			fmt.Println("BuyResultInterface() auction not find err=", err.Error)
			return ErrNftNotExist
		}
		if price == "" {
			fmt.Println("BuyResultInterface() price == null" )
			return nft.db.Transaction(func(tx *gorm.DB) error {
				trans := Trans{}
				trans.Contract = contractAddr
				trans.Createaddr = nftRec.Createaddr
				trans.Fromaddr = from
				trans.Toaddr = to
				trans.Signdata = sig
				trans.Tradesig = trade_sig
				trans.Tokenid = tokenId
				trans.Nftid = nftRec.ID
				trans.Price = 0
				trans.Transtime = time.Now().Unix()
				trans.Selltype = SellTypeAsset.String()
				err := tx.Model(&trans).Create(&trans)
				if err.Error != nil {
					fmt.Println("BuyResultInterface() create trans record err=", err.Error)
					return err.Error
				}
				fmt.Println("BuyResultInterface() price == null OK" )
				return nil
			})
		}else{
			fmt.Println("BuyResultInterface() price != null" )
			return nft.db.Transaction(func(tx *gorm.DB) error {
				trans := Trans{}
				trans.Contract = contractAddr
				trans.Createaddr = nftRec.Createaddr
				trans.Fromaddr = from
				trans.Toaddr = to
				trans.Signdata = sig
				trans.Tradesig = trade_sig
				trans.Nftid = nftRec.ID
				trans.Tokenid = tokenId
				trans.Txhash = txhash
				trans.Name = nftRec.Name
				trans.Meta = nftRec.Meta
				trans.Desc = nftRec.Desc
				trans.Price, _ = strconv.ParseUint(price, 10, 64)
				trans.Transtime = time.Now().Unix()
				trans.Selltype = SellTypeForeignPrice.String()
				err := tx.Model(&trans).Create(&trans)
				if err.Error != nil {
					fmt.Println("BuyResultInterface() create trans record err=", err.Error)
					return err.Error
				}
				nftrecord := Nfts{}
				nftrecord.Ownaddr = to
				nftrecord.Selltype = SellTypeNotSale.String()
				nftrecord.Transprice = trans.Price
				nftrecord.Transamt += trans.Price
				nftrecord.Transcnt += 1
				nftrecord.Transtime = time.Now().Unix()
				err = tx.Model(&Nfts{}).Where("contract = ? AND tokenid =?",
					contractAddr, tokenId).Updates(&nftrecord)
				if err.Error != nil {
					fmt.Println("BuyResultInterface() update record err=", err.Error)
					return err.Error
				}
				fmt.Println("BuyResultInterface() from != Null && to != Null --> price != Null OK" )
				return nil
			})
		}
	}
	fmt.Println("BuyResultInterface() End.")
	return ErrFromToAddrZero
}

