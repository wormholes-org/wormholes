package models

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
)

func (nft NftDb) MakeOffer(userAddr,
	contractAddr,
	tokenId string,
	PayChannel string,
	CurrencyType string,
	price uint64,
	TradeSig string,
	dead_time int64,
	sigdata string) error {
	userAddr = strings.ToLower(userAddr)
	contractAddr = strings.ToLower(contractAddr)

	var nftrecord Nfts
	err := nft.db.Where("contract = ? AND tokenid =?", contractAddr, tokenId).First(&nftrecord)
	if err.Error != nil {
		fmt.Println("MakeOffer() bidprice not find nft err= ", err.Error )
		return ErrNftNotExist
	}
	if nftrecord.Ownaddr == userAddr {
		fmt.Println("MakeOffer() don't buy your own nft.")
		return ErrBuyOwn
	}

	var auctionRec Auction
	err = nft.db.Where("contract = ? AND tokenid = ?", contractAddr, tokenId).First(&auctionRec)
	if err.Error != nil {
		if err.Error == gorm.ErrRecordNotFound {
			fmt.Println("MakeOffer() RecordNotFound")
			auctionRec = Auction{}
			auctionRec.Selltype = SellTypeBidPrice.String()
			auctionRec.Paychan = PayChannel
			auctionRec.Ownaddr = nftrecord.Ownaddr
			auctionRec.Nftid = nftrecord.ID
			auctionRec.Contract = contractAddr
			auctionRec.Tokenid = tokenId
			auctionRec.Currency = CurrencyType
			//auctionRec.Startprice = price
			//auctionRec.Endprice = price
			auctionRec.Startdate = time.Now().Unix()
			auctionRec.Enddate = time.Now().Unix()
			auctionRec.Signdata = sigdata
			auctionRec.Tradesig = TradeSig
			auctionHistory := AuctionHistory{}
			auctionHistory.AuctionRecord = auctionRec.AuctionRecord
			return nft.db.Transaction(func(tx *gorm.DB) error {
				err = tx.Model(&auctionRec).Create(&auctionRec)
				if err.Error != nil {
					fmt.Println("MakeOffer() create auctionRec record err=", err.Error)
					return err.Error
				}
				err = tx.Model(&AuctionHistory{}).Create(&auctionHistory)
				if err.Error != nil {
					fmt.Println("MakeOffer() create auctionHistory record err=", err.Error)
					return err.Error
				}
				nftrecord = Nfts{}
				nftrecord.Selltype = auctionRec.Selltype
				err = tx.Model(&Nfts{}).Where("contract = ? AND tokenid =?",
					auctionRec.Contract, auctionRec.Tokenid).Updates(&nftrecord)
				if err.Error != nil {
					fmt.Println("MakeOffer() update record err=", err.Error)
					return err.Error
				}
				bidRec := Bidding{}
				bidRec.Bidaddr = userAddr
				bidRec.Auctionid = auctionRec.ID
				bidRec.Contract = contractAddr
				bidRec.Tokenid = tokenId
				bidRec.Price = price
				bidRec.Currency = CurrencyType
				bidRec.Paychan = PayChannel
				bidRec.Tradesig = TradeSig
				bidRec.Bidtime = time.Now().Unix()
				bidRec.Signdata = sigdata
				bidRec.Deadtime = dead_time
				bidRec.Nftid = auctionRec.Nftid
				bidRecHistory := BiddingHistory{}
				bidRecHistory.BidRecord = bidRec.BidRecord
				err := tx.Model(&bidRec).Create(&bidRec)
				if err.Error != nil {
					fmt.Println("MakeOffer() create bidRec record err=", err.Error)
					return err.Error
				}
				err = tx.Model(&BiddingHistory{}).Create(&bidRecHistory)
				if err.Error != nil {
					fmt.Println("MakeOffer() create bidRecHistory record err=", err.Error)
					return err.Error
				}
				fmt.Println("MakeOffer() RecordNotFound OK")
				return nil
			})
		}
		return ErrNftNotSell
	}
	//if time.Now().Unix() < auctionRec.Startdate {
	//	return ErrAuctionNotBegan
	//}

	if auctionRec.Selltype == SellTypeHighestBid.String() {
		//addrs, err := ethhelper.BalanceOfWeth()
		fmt.Println("MakeOffer() Selltype == SellTypeHighestBid")
		if time.Now().Unix() >= auctionRec.Enddate {
			fmt.Println("MakeOffer() time.Now().Unix() >= auctionRec.Enddate")
			return ErrAuctionEnd
		}
		if auctionRec.Startprice > price {
			fmt.Println("MakeOffer() auctionRec.Startprice > price")
			return ErrBidOutRange
		}
		var bidRec Bidding
		err = nft.db.Where("contract = ? AND tokenid = ? AND bidAddr = ?", contractAddr, tokenId, userAddr).First(&bidRec)
		if err.Error == nil {
			fmt.Println("MakeOffer() first bidding.")
			bidRec = Bidding{}
			bidRec.Price = price
			bidRec.Currency = CurrencyType
			bidRec.Paychan = PayChannel
			bidRec.Tradesig = TradeSig
			bidRec.Bidtime = time.Now().Unix()
			bidRec.Signdata = sigdata
			return nft.db.Transaction(func(tx *gorm.DB) error {
				err := tx.Model(&bidRec).Where("contract = ? AND tokenid = ? AND bidAddr = ?", contractAddr, tokenId, userAddr).Updates(&bidRec)
				if err.Error != nil {
					fmt.Println("MakeOffer() update Bidding record err=", err.Error)
					return err.Error
				}
				bidRecHistory := BiddingHistory(bidRec)
				err = tx.Model(&BiddingHistory{}).Where("contract = ? AND tokenid = ? AND bidAddr = ?", contractAddr, tokenId, userAddr).Updates(&bidRecHistory)
				if err.Error != nil {
					fmt.Println("MakeOffer() update bidRecHistory record err=", err.Error)
					return err.Error
				}
				fmt.Println("MakeOffer() first bidding OK.")
				return nil
			})
		} else{
			bidRec = Bidding{}
			bidRec.Bidaddr = userAddr
			bidRec.Auctionid = auctionRec.ID
			bidRec.Nftid = auctionRec.Nftid
			bidRec.Contract = contractAddr
			bidRec.Tokenid = tokenId
			bidRec.Price = price
			bidRec.Currency = CurrencyType
			bidRec.Paychan = PayChannel
			bidRec.Tradesig = TradeSig
			bidRec.Bidtime = time.Now().Unix()
			bidRec.Signdata = sigdata
			return nft.db.Transaction(func(tx *gorm.DB) error {
				err := tx.Model(&bidRec).Create(&bidRec)
				if err.Error != nil {
					fmt.Println("MakeOffer() create record err=", err.Error)
					return err.Error
				}
				bidRecHistory := BiddingHistory{}
				bidRecHistory.BidRecord = bidRec.BidRecord
				err = tx.Model(&BiddingHistory{}).Create(&bidRecHistory)
				if err.Error != nil {
					fmt.Println("MakeOffer() create bidRecHistory record err=", err.Error)
					return err.Error
				}
				fmt.Println("MakeOffer() change bidding OK.")
				return nil
			})
		}
	}
	if auctionRec.Selltype == SellTypeBidPrice.String() {
		fmt.Println("MakeOffer() Selltype == SellTypeBidPrice")
		var bidRec Bidding
		err = nft.db.Where("contract = ? AND tokenid = ? AND bidAddr = ?", contractAddr, tokenId, userAddr).First(&bidRec)
		if err.Error == nil {
			bidRec = Bidding{}
			bidRec.Price = price
			bidRec.Currency = CurrencyType
			bidRec.Paychan = PayChannel
			bidRec.Tradesig = TradeSig
			bidRec.Bidtime = time.Now().Unix()
			bidRec.Signdata = sigdata
			return nft.db.Transaction(func(tx *gorm.DB) error {
				err := tx.Model(&bidRec).Where("contract = ? AND tokenid = ? AND bidAddr = ?", contractAddr, tokenId, userAddr).Updates(&bidRec)
				if err.Error != nil {
					fmt.Println("MakeOffer() update Bidding record err=", err.Error)
					return err.Error
				}
				bidRecHistory := BiddingHistory(bidRec)
				err = tx.Model(&BiddingHistory{}).Where("contract = ? AND tokenid = ? AND bidAddr = ?", contractAddr, tokenId, userAddr).Updates(&bidRecHistory)
				if err.Error != nil {
					fmt.Println("MakeOffer() update bidRecHistory record err=", err.Error)
					return err.Error
				}
				fmt.Println("MakeOffer() change bidding OK.")
				return nil
			})
		} else {
			return nft.db.Transaction(func(tx *gorm.DB) error {
				bidRec := Bidding{}
				bidRec.Bidaddr = userAddr
				bidRec.Auctionid = auctionRec.ID
				bidRec.Nftid = auctionRec.Nftid
				bidRec.Contract = contractAddr
				bidRec.Tokenid = tokenId
				bidRec.Price = price
				bidRec.Currency = CurrencyType
				bidRec.Paychan = PayChannel
				bidRec.Tradesig = TradeSig
				bidRec.Bidtime = time.Now().Unix()
				bidRec.Deadtime = dead_time
				bidRec.Signdata = sigdata
				bidRecHistory := BiddingHistory{}
				bidRecHistory.BidRecord = bidRec.BidRecord
				err := tx.Model(&bidRec).Create(&bidRec)
				if err.Error != nil {
					fmt.Println("MakeOffer() create bidRec record err=", err.Error)
					return err.Error
				}
				err = tx.Model(&BiddingHistory{}).Create(&bidRecHistory)
				if err.Error != nil {
					fmt.Println("MakeOffer() create bidRecHistory record err=", err.Error)
					return err.Error
				}
				fmt.Println("MakeOffer() first bidding OK.")
				return nil
			})
		}
	}
	if auctionRec.Selltype == SellTypeFixPrice.String() {
		fmt.Println("MakeOffer() Selltype == SellTypeFixPrice")
		var bidRec Bidding
		err = nft.db.Where("contract = ? AND tokenid = ? AND bidAddr = ?", contractAddr, tokenId, userAddr).First(&bidRec)
		if err.Error == nil {
			bidRec = Bidding{}
			bidRec.Price = price
			bidRec.Currency = CurrencyType
			bidRec.Paychan = PayChannel
			bidRec.Tradesig = TradeSig
			bidRec.Bidtime = time.Now().Unix()
			bidRec.Signdata = sigdata
			return nft.db.Transaction(func(tx *gorm.DB) error {
				err := tx.Model(&bidRec).Where("contract = ? AND tokenid = ? AND bidAddr = ?", contractAddr, tokenId, userAddr).Updates(&bidRec)
				if err.Error != nil {
					fmt.Println("MakeOffer() update Bidding record err=", err.Error)
					return err.Error
				}
				bidRecHistory := BiddingHistory(bidRec)
				err = tx.Model(&BiddingHistory{}).Where("contract = ? AND tokenid = ? AND bidAddr = ?", contractAddr, tokenId, userAddr).Updates(&bidRecHistory)
				if err.Error != nil {
					fmt.Println("MakeOffer() update bidRecHistory record err=", err.Error)
					return err.Error
				}
				fmt.Println("MakeOffer() change bidding OK.")
				return nil
			})
		} else {
			return nft.db.Transaction(func(tx *gorm.DB) error {
				bidRec := Bidding{}
				bidRec.Bidaddr = userAddr
				bidRec.Auctionid = auctionRec.ID
				bidRec.Nftid = auctionRec.Nftid
				bidRec.Contract = contractAddr
				bidRec.Tokenid = tokenId
				bidRec.Price = price
				bidRec.Currency = CurrencyType
				bidRec.Paychan = PayChannel
				bidRec.Tradesig = TradeSig
				bidRec.Bidtime = time.Now().Unix()
				bidRec.Deadtime = dead_time
				bidRec.Signdata = sigdata
				bidRecHistory := BiddingHistory{}
				bidRecHistory.BidRecord = bidRec.BidRecord
				err := tx.Model(&bidRec).Create(&bidRec)
				if err.Error != nil {
					fmt.Println("MakeOffer() create bidRec record err=", err.Error)
					return err.Error
				}
				err = tx.Model(&BiddingHistory{}).Create(&bidRecHistory)
				if err.Error != nil {
					fmt.Println("MakeOffer() create bidRecHistory record err=", err.Error)
					return err.Error
				}
				fmt.Println("MakeOffer() first bidding OK.")
				return nil
			})
		}
	}
	return ErrNftNotSell
}

func (nft NftDb) Sell(ownAddr,
	PrivAddr string,
	contractAddr,
	tokenId string,
	sellType string,
	payChan string,
	days int,
	startPrice,
	endPrice uint64,
	royalty string,
	currency string,
	hide string,
	sigData string,
	tradeSig string) error {

	ownAddr = strings.ToLower(ownAddr)
	PrivAddr = strings.ToLower(PrivAddr)
	contractAddr = strings.ToLower(contractAddr)

	fmt.Println(time.Now().String()[:22], "Sell() Start.",
		"tokenId=", tokenId,
		"SellType=", sellType,
		"startPrice=", startPrice,
		"endPrice=", endPrice)
	defer fmt.Println(time.Now().String()[:22], "Sell() end.")
	var nftrecord Nfts
	err := nft.db.Where("contract = ? AND tokenid =? AND ownaddr = ?", contractAddr, tokenId, ownAddr).First(&nftrecord)
	if err.Error != nil {
		fmt.Println("Sell() err= ", err.Error )
		return err.Error
	}
	if nftrecord.Verified != Passed.String() {
		return ErrNotVerify
	}
	/*if nftrecord.Mintstate != Minted.String() {
		return ErrNftNotMinted
	}*/
	//if startDate.After(endDate) {
	//	return ErrAuctionStartAfterEnd
	//}
	//if startDate.Before(time.Now()) {
	//	startDate = time.Now()
	//	//return ErrAuctionStartBeforeNow
	//}
	var auctionRec Auction
	err = nft.db.Where("contract = ? AND nftid = ? AND ownaddr = ?",
		nftrecord.Contract, nftrecord.ID, ownAddr).First(&auctionRec)
	if err.Error == nil {
		if auctionRec.Selltype != SellTypeBidPrice.String() {
			fmt.Println("Sell() err=", err.Error, ErrNftSelling)
			return ErrNftSelling
		} else {
			err := nft.db.Transaction(func(tx *gorm.DB) error {
				err = tx.Model(&Bidding{}).Where("contract = ? AND tokenid = ?",
					auctionRec.Contract, auctionRec.Tokenid).Delete(&Bidding{})
				if err.Error != nil {
					fmt.Println("Sell() delete bid record err=", err.Error)
					return err.Error
				}
				err = tx.Model(&Auction{}).Where("contract = ? AND tokenid = ?",
					auctionRec.Contract, auctionRec.Tokenid).Delete(&Auction{})
				if err.Error != nil {
					fmt.Println("Sell() delete bidprice auction record err=", err.Error)
					return err.Error
				}
				return nil
			})
			if err != nil {
				fmt.Println("Sell() delete bidprice err=", err)
				return err
			}
		}
	}
	auctionRec = Auction{}
	auctionRec.Selltype = sellType
	auctionRec.Paychan = payChan
	auctionRec.Ownaddr = ownAddr
	auctionRec.Nftid = nftrecord.ID
	auctionRec.Contract = contractAddr
	auctionRec.Tokenid = tokenId
	auctionRec.Currency = currency
	auctionRec.Startprice = startPrice
	auctionRec.Endprice = endPrice
	auctionRec.Privaddr = PrivAddr
	auctionRec.Startdate = time.Now().Unix()
	auctionRec.Enddate = time.Now().AddDate(0, 0, days).Unix()
	//auctionRec.Enddate = time.Now().Add(1 * time.Minute).Unix()
	auctionRec.Signdata = sigData
	auctionRec.Tradesig = tradeSig
	auctionRec.SellState = SellStateStart.String()

	if sellType == SellTypeFixPrice.String() {
		auctionRec.Startprice = startPrice
		auctionRec.Endprice = startPrice
	}
	auctionHistory := AuctionHistory{}
	auctionHistory.AuctionRecord = auctionRec.AuctionRecord
	return nft.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&auctionRec).Create(&auctionRec)
		if err.Error != nil {
			fmt.Println("Sell() create auctionRec record err=", err.Error)
			return err.Error
		}
		err = tx.Model(&AuctionHistory{}).Create(&auctionHistory)
		if err.Error != nil {
			fmt.Println("Sell() create auctionHistory record err=", err.Error)
			return err.Error
		}
		nftrecord = Nfts{}
		nftrecord.Hide = hide
		nftrecord.Selltype = sellType
		err = tx.Model(&Nfts{}).Where("contract = ? AND tokenid =?",
			auctionRec.Contract, auctionRec.Tokenid).Updates(&nftrecord)
		if err.Error != nil {
			fmt.Println("Sell() update record err=", err.Error)
			return err.Error
		}
		/*nftrecord = Nfts{}
		nftrecord.Royalty, _ = strconv.Atoi(royalty)
		nftrecord.Royalty = nftrecord.Royalty / 100
		err = tx.Model(&Nfts{}).Where("contract = ? AND tokenid =? AND royalty = ?",
			auctionRec.Contract, auctionRec.Tokenid, 0).Updates(&nftrecord)
		if err.Error != nil {
			fmt.Println("Sell() update record err=", err.Error)
			return err.Error
		}*/
		return nil
	})
}

