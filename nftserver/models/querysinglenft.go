package models

import "strings"

type NftAuction struct {
	Selltype        string `json:"selltype"`
	Ownaddr         string `json:"ownaddr"`
	NftTokenId      string `json:"nft_token_id"`
	NftContractAddr string `json:"nft_contract_addr"`
	Paychan         string `json:"paychan"`
	Currency        string `json:"currency"`
	Startprice      uint64 `json:"startprice"`
	Endprice        uint64 `json:"endprice"`
	Startdate       int64  `json:"startdate"`
	Enddate         int64  `json:"enddate"`
	Tradesig       	string `json:"tradesig"`
}

type NftTran struct {
	NftContractAddr string `json:"nft_contract_addr"`
	Fromaddr        string `json:"fromaddr"`
	Toaddr          string `json:"toaddr"`
	NftTokenId      string `json:"nft_token_id"`
	Transtime       int64  `json:"transtime"`
	Paychan         string `json:"paychan"`
	Currency        string `json:"currency"`
	Price           uint64 `json:"price"`
	Txhash			string `json:"trade_hash"`
	Selltype        string `json:"selltype"`
}

type NftBid struct {
	Bidaddr         string `json:"bidaddr"`
	NftTokenId      string `json:"nft_token_id"`
	NftContractAddr string `json:"nft_contract_addr"`
	Paychan         string `json:"paychan"`
	Currency        string `json:"currency"`
	Price           uint64 `json:"price"`
	Bidtime         int64  `json:"bidtime"`
	Tradesig       	string `json:"tradesig"`
}

type NftSingleInfo struct {
	Name            string 			`json:"name"`
	CreatorAddr     string 			`json:"creator_addr"`
	//CreatorPortrait string 			`json:"creator_portrait"`
	OwnerAddr       string 			`json:"owner_addr"`
	//OwnerPortrait   string 			`json:"owner_portrait"`
	Md5             string 			`json:"md5"`
	//AssetSample     string 			`json:"asset_sample"`
	Desc            string 			`json:"desc"`
	Collectiondesc  string 			`json:"collection_desc"`
	Meta            string 			`json:"meta"`
	SourceUrl       string 			`json:"source_url"`
	NftContractAddr string 			`json:"nft_contract_addr"`
	NftTokenId      string 			`json:"nft_token_id"`
	Categories      string 			`json:"categories"`
	CollectionCreatorAddr string    `json:"collection_creator_addr"`
	Collections     string 			`json:"collections"`
	//Img             string 			`json:"img"`
	Approve         string 			`json:"approve"`
	Royalty         int 			`json:"royalty"`
	Verified        string 			`json:"verified"`
	Selltype        string 			`json:"selltype"`
	Mintstate       string	 		`json:"mintstate"`
	Likes	        int 			`json:"likes"`

	Auction 		NftAuction		`json:"auction"`
	Trans   		[]NftTran		`json:"trans"`
	Bids    		[]NftBid	 	`json:"bids"`
}

func (nft NftDb) QuerySingleNft(contract, tokenId string) (NftSingleInfo, error) {
	contract = strings.ToLower(contract)

	var nftInfo NftSingleInfo

	var nftRecord Nfts
	err := nft.db.Where("contract = ? AND tokenid = ?", contract, tokenId).First(&nftRecord)
	if err.Error != nil {
		return NftSingleInfo{}, ErrNftNotExist
	}
	nftInfo.Name = nftRecord.Name
	nftInfo.CreatorAddr = nftRecord.Createaddr
	nftInfo.OwnerAddr = nftRecord.Ownaddr
	nftInfo.Md5 = nftRecord.Md5
	//nftInfo.AssetSample = nftRecord.Image
	nftInfo.Desc = nftRecord.Desc
	nftInfo.Meta =  nftRecord.Meta
	nftInfo.SourceUrl = nftRecord.Url
	nftInfo.NftContractAddr = nftRecord.Contract
	nftInfo.NftTokenId = nftRecord.Tokenid
	nftInfo.Categories = nftRecord.Categories
	nftInfo.Collections = nftRecord.Collections
	nftInfo.Approve = nftRecord.Approve
	nftInfo.Royalty = nftRecord.Royalty
	nftInfo.Verified = nftRecord.Verified
	nftInfo.Selltype = nftRecord.Selltype
	nftInfo.Mintstate = nftRecord.Mintstate
	nftInfo.Likes = nftRecord.Favorited

	user := Users{}
	err = nft.db.Where("useraddr = ?", nftRecord.Createaddr).First(&user)
	if err.Error == nil {
		//nftInfo.CreatorPortrait = user.Portrait
	}
	user = Users{}
	err = nft.db.Where("useraddr = ?", nftRecord.Ownaddr).First(&user)
	if err.Error == nil {
		//nftInfo.OwnerPortrait = user.Portrait
	}
	var collectRec Collects
	err = nft.db.Where("Createaddr = ? AND name = ? ", nftRecord.Collectcreator, nftRecord.Collections).First(&collectRec)
	if err.Error == nil {
		//nftInfo.Img = collectRec.Img
		nftInfo.CollectionCreatorAddr = collectRec.Createaddr
		nftInfo.Collectiondesc = collectRec.Desc
	}

	var auctionRec Auction
	err = nft.db.Where("contract = ? AND tokenid = ?", contract, tokenId).First(&auctionRec)
	if err.Error == nil {
		nftInfo.Auction.Selltype = auctionRec.Selltype
		nftInfo.Auction.Ownaddr = auctionRec.Ownaddr
		nftInfo.Auction.NftTokenId = auctionRec.Tokenid
		nftInfo.Auction.NftContractAddr = auctionRec.Contract
		nftInfo.Auction.Paychan = auctionRec.Paychan
		nftInfo.Auction.Currency = auctionRec.Currency
		nftInfo.Auction.Startprice = auctionRec.Startprice
		nftInfo.Auction.Endprice = auctionRec.Endprice
		nftInfo.Auction.Startdate = auctionRec.Startdate
		nftInfo.Auction.Enddate = auctionRec.Enddate
		nftInfo.Auction.Tradesig = auctionRec.Tradesig
	}

	trans := make([]Trans, 0, 20)
	err = nft.db.Where("contract = ? AND tokenid = ? AND selltype != ? AND selltype != ? AND price > 0",
		contract, tokenId, SellTypeMintNft.String(), SellTypeError.String()).Find(&trans)
	/*err = nft.db.Raw("SELECT * FROM trans\n WHERE id IN (SELECT MAX(id) AS o FROM trans GROUP BY contract, tokenId, Auctionid) " +
	"and contract = ? and tokenid = ?  and \n  Selltype !=\"MintNft\"",
	contract, tokenId).Find(&trans)*/
	if err.Error == nil {
		if err.RowsAffected != 0 {
			for _, tran := range trans {
				var nfttran NftTran
				nfttran.NftContractAddr = tran.Contract
				nfttran.Fromaddr = tran.Fromaddr
				nfttran.Toaddr = tran.Toaddr
				nfttran.NftTokenId = tran.Tokenid
				nfttran.Transtime = tran.Transtime
				nfttran.Paychan = tran.Paychan
				nfttran.Currency = tran.Currency
				nfttran.Price = tran.Price
				nfttran.Selltype = tran.Selltype
				nfttran.Txhash = tran.Txhash
				nftInfo.Trans = append(nftInfo.Trans, nfttran)
			}
		}
	}
	bids := make([]Bidding, 0, 20)
	err = nft.db.Where("contract = ? AND tokenid = ?", contract, tokenId).Find(&bids)
	if err.Error == nil {
		if err.RowsAffected != 0 {
			for _, bid := range bids {
				var nftbid NftBid
				nftbid.Bidaddr = bid.Bidaddr
				nftbid.NftTokenId = bid.Tokenid
				nftbid.NftContractAddr = bid.Contract
				nftbid.Paychan = bid.Paychan
				nftbid.Currency = bid.Currency
				nftbid.Price = bid.Price
				nftbid.Bidtime = bid.Bidtime
				nftbid.Tradesig = bid.Tradesig
				nftInfo.Bids = append(nftInfo.Bids, nftbid)
			}
		}
	}
	return nftInfo, nil
}

