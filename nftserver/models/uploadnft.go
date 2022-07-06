package models

import (
	"fmt"
	"github.com/nftexchange/nftserver/ethhelper"

	//"github.com/nftexchange/ethhelper"
	"gorm.io/gorm"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func (nft NftDb) UploadNft(
	user_addr string,
	creator_addr string,
	owner_addr string,
	md5 string,
	name string,
	desc string,
	meta string,
	source_url string,
	nft_contract_addr string,
	nft_token_id string,
	categories string,
	collections string,
	asset_sample string,
	hide string,
	royalty string,
	count string,
	sig string) error {
	user_addr = strings.ToLower(user_addr)
	creator_addr = strings.ToLower(creator_addr)
	owner_addr = strings.ToLower(owner_addr)
	nft_contract_addr = strings.ToLower(nft_contract_addr)

	fmt.Println("UploadNft() begin ->> time = ", time.Now().String()[:22])
	fmt.Println("UploadNft() user_addr = ", user_addr)
	fmt.Println("UploadNft() creator_addr = ", creator_addr)
	fmt.Println("UploadNft() owner_addr = ", owner_addr)
	fmt.Println("UploadNft() md5 = ", md5)
	fmt.Println("UploadNft() name = ", name)
	fmt.Println("UploadNft() desc = ", desc)
	fmt.Println("UploadNft() meta = ", meta)
	fmt.Println("UploadNft() source_url = ", source_url)
	fmt.Println("UploadNft() nft_contract_addr = ", nft_contract_addr)
	fmt.Println("UploadNft() nft_token_id = ", nft_token_id)
	fmt.Println("UploadNft() categories = ", categories)
	fmt.Println("UploadNft() collections = ", collections)
	//fmt.Println("UploadNft() asset_sample = ", asset_sample)
	fmt.Println("UploadNft() hide = ", hide)
	fmt.Println("UploadNft() royalty = ", royalty)
	//fmt.Println("UploadNft() sig = ", sig)

	if IsIntDataValid(count) != true {
		return ErrDataFormat
	}
	if IsIntDataValid(royalty) != true {
		return ErrDataFormat
	}
	r, _ := strconv.Atoi(royalty)
	fmt.Println("UploadNft() royalty=", r, "SysRoyaltylimit=", SysRoyaltylimit, "RoyaltyLimit", RoyaltyLimit )
	if r > SysRoyaltylimit || r > RoyaltyLimit {
		return ErrRoyalty
	}
	if count == "" {
		count = "1"
	}
	if c, _ := strconv.Atoi(count); c < 1 {
		fmt.Println("UploadNft() contract count < 1.")
		return ErrContractCountLtZero
	}
	if nft.IsValidCategory(categories) {
		return ErrNoCategory
	}

	var collectRec Collects
	if collections != "" {
		err := nft.db.Where("name = ? AND createaddr =?",
			collections, creator_addr).First(&collectRec)
		if err.Error != nil {
			fmt.Println("UploadNft() err=Collection not exist.")
			return ErrCollectionNotExist
		}
	} else {
		return ErrCollectionNotExist
	}
	if nft_contract_addr == "" && nft_token_id == "" {
		var NewTokenid string
		rand.Seed(time.Now().UnixNano())
		var i int
		for i = 0; i < genTokenIdRetry ; i++ {
			//NewTokenid := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
			s := fmt.Sprintf("%d", rand.Int63())
			if len(s) < 15 {
				continue
			}
			s = s[len(s)-13:]
			NewTokenid = s
			if s[0] == '0' {
				continue
			}
			fmt.Println("UploadNft() NewTokenid=", NewTokenid)
			nfttab :=  Nfts{}
			err := nft.db.Where("contract = ? AND tokenid = ? ", NFT1155Addr, NewTokenid).First(&nfttab)
			if err.Error == gorm.ErrRecordNotFound {
				break
			}
		}
		if i >= 20 {
			fmt.Println("UploadNft() generate tokenId error.")
			return ErrGenerateTokenId
		}
		nfttab :=  Nfts{}
		nfttab.Tokenid = NewTokenid
		//nfttab.Contract = strings.ToLower(ExchangAddr) //nft_contract_addr
		nfttab.Contract = collectRec.Contract
		nfttab.Createaddr = creator_addr
		nfttab.Ownaddr = owner_addr
		nfttab.Name = name
		nfttab.Desc = desc
		nfttab.Meta = meta
		nfttab.Categories = categories
		nfttab.Collectcreator = collectRec.Createaddr
		nfttab.Collections = collections
		nfttab.Signdata = sig
		nfttab.Url = source_url
		nfttab.Image = asset_sample
		nfttab.Md5 = md5
		nfttab.Selltype = SellTypeNotSale.String()
		nfttab.Verified = NoVerify.String()
		if collectRec.Contract == strings.ToLower(NFT1155Addr) {
			nfttab.Mintstate = NoMinted.String()
		} else {
			//调用查询函数
			nfttab.Mintstate = Minted.String()
		}
		nfttab.Createdate = time.Now().Unix()
		nfttab.Royalty, _ = strconv.Atoi(royalty)
		//nfttab.Royalty /= 100
		nfttab.Count, _ = strconv.Atoi(count)
		nfttab.Hide = hide
		err0, approve := ethhelper.GenCreateNftSign(NFT1155Addr, nfttab.Ownaddr, nfttab.Meta,
			nfttab.Tokenid, count, royalty)
		if err0 != nil {
			fmt.Println("UploadNft() GenCreateNftSign() err=", err0)
			return err0
		}
		fmt.Println("UploadNft() GenCreateNftSign() approve=", approve)
		nfttab.Approve = approve
		return nft.db.Transaction(func(tx *gorm.DB) error {
			err := tx.Model(&Nfts{}).Create(&nfttab)
			if err.Error != nil {
				fmt.Println("UploadNft() err=", err.Error)
				return err.Error
			}
			if collections != "" {
				var collectListRec CollectLists
				collectListRec.Collectsid = collectRec.ID
				collectListRec.Nftid = nfttab.ID
				err = tx.Model(&CollectLists{}).Create(&collectListRec)
				if err.Error != nil {
					fmt.Println("UploadNft() create CollectLists err=", err.Error)
					return err.Error
				}
				err = tx.Model(&Collects{}).Where("name = ? AND createaddr =?",
					collections, creator_addr).Update("totalcount",collectRec.Totalcount+1)
				if err.Error != nil {
					fmt.Println("UploadNft() add collectins totalcount err= ", err.Error )
					return err.Error
				}
			}
			return nil
		})
	} else {
		var nfttab Nfts
		dberr := nft.db.Where("contract = ? AND tokenid = ? ", nft_contract_addr, nft_token_id).First(&nfttab)
		if dberr.Error == nil {
			fmt.Println("UploadNft() err=nft already exist.")
			return ErrNftAlreadyExist
		}
		/*ownAddr, royalty, err := func(contract, tokenId string) (string, string, error) {
			return "ownAddr", "200", nil
		}(nft_contract_addr, nft_token_id)
		if ownAddr == user_addr {
			var nfttab Nfts
			nfttab.Tokenid = nft_token_id
			nfttab.Contract = nft_contract_addr //nft_contract_addr
			nfttab.Createaddr = creator_addr
			nfttab.Ownaddr = ownAddr
			nfttab.Name = name
			nfttab.Desc = desc
			nfttab.Meta = meta
			nfttab.Categories = categories
			nfttab.Collections = collections
			nfttab.Signdata = sig
			nfttab.Url = source_url
			nfttab.Image = asset_sample
			nfttab.Md5 = md5
			nfttab.Selltype = SellTypeNotSale.String()
			nfttab.Verified = NoVerify.String()
			nfttab.Mintstate = Minted.String()
			nfttab.Royalty, _ = strconv.Atoi(royalty)
			nfttab.Royalty = nfttab.Royalty / 100
			nfttab.Createdate = time.Now().Unix()
			nfttab.Hide = hide
			return nft.db.Transaction(func(tx *gorm.DB) error {
				err := tx.Model(&Nfts{}).Create(&nfttab)
				if err.Error != nil {
					fmt.Println("UploadNft() create exist nft err=", err.Error)
					return err.Error
				}
				if collections != "" {
					var collectListRec CollectLists
					collectListRec.Collectsid = collectRec.ID
					collectListRec.Nftid = nfttab.ID
					err = tx.Model(&CollectLists{}).Create(&collectListRec)
					if err.Error != nil {
						fmt.Println("UploadNft() create CollectLists err=", err.Error)
						return err.Error
					}
				}
				return nil
			})
		}*/
		IsAdminAddr, err := IsAdminAddr(user_addr)
		if err != nil {
			fmt.Println("UploadNft() upload address is not admin.")
			return ErrNftUpAddrNotAdmin
		}
		if IsAdminAddr {
			var nfttab Nfts
			nfttab.Tokenid = nft_token_id
			nfttab.Contract = nft_contract_addr //nft_contract_addr
			nfttab.Createaddr = creator_addr
			nfttab.Ownaddr = owner_addr
			nfttab.Name = name
			nfttab.Desc = desc
			nfttab.Meta = meta
			nfttab.Categories = categories
			nfttab.Collectcreator = creator_addr
			nfttab.Collections = collections
			nfttab.Signdata = sig
			nfttab.Url = source_url
			nfttab.Image = asset_sample
			nfttab.Md5 = md5
			nfttab.Selltype = SellTypeNotSale.String()
			nfttab.Verified = Passed.String()
			nfttab.Mintstate = Minted.String()
			/*nfttab.Royalty, _ = strconv.Atoi(royalty)
			nfttab.Royalty = nfttab.Royalty / 100*/
			nfttab.Createdate = time.Now().Unix()
			nfttab.Royalty, _ = strconv.Atoi(royalty)
			//nfttab.Royalty /= 100
			nfttab.Count, _ = strconv.Atoi(count)
			nfttab.Hide = hide
			return nft.db.Transaction(func(tx *gorm.DB) error {
				err := tx.Model(&Nfts{}).Create(&nfttab)
				if err.Error != nil {
					fmt.Println("UploadNft() admin create nft err=", err.Error)
					return err.Error
				}
				if collections != "" {
					var collectListRec CollectLists
					collectListRec.Collectsid = collectRec.ID
					collectListRec.Nftid = nfttab.ID
					err = tx.Model(&CollectLists{}).Create(&collectListRec)
					if err.Error != nil {
						fmt.Println("UploadNft() create CollectLists err=", err.Error)
						return err.Error
					}
					err = tx.Model(&Collects{}).Where("name = ? AND createaddr =?",
						collections, creator_addr).Update("totalCount",collectRec.Totalcount+1)
					if err.Error != nil {
						fmt.Println("UploadNft() add collectins totalcount err= ", err.Error )
						return err.Error
					}
				}
				return nil
			})
		} else {
			fmt.Println("UploadNft() upload address is not admin.")
			return ErrNftUpAddrNotAdmin
		}
	}
	return nil
}
