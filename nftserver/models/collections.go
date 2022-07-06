package models

import (
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

func (nft NftDb) NewCollections(useraddr, name, img, contract_type, contract_addr,
	desc, categories, sig string) error {
	useraddr = strings.ToLower(useraddr)
	contract_addr = strings.ToLower(contract_addr)

	var collectRec Collects
	err := nft.db.Where("Createaddr = ? AND name = ? ", useraddr, name).First(&collectRec)
	if err.Error == nil {
		fmt.Println("NewCollections() err=Collection already exist." )
		return ErrCollectionExist
	} else if err.Error == gorm.ErrRecordNotFound {
		fmt.Println("NewCollections() err=Collection already exist.")
		collectRec = Collects{}
		collectRec.Createaddr = useraddr
		collectRec.Name = name
		collectRec.Desc = desc
		collectRec.Img = img
		if contract_addr != "" {
			collectRec.Contract = contract_addr
		} else {
			collectRec.Contract = strings.ToLower(NFT1155Addr)
		}
		collectRec.Contracttype = contract_type
		collectRec.Categories = categories
		collectRec.SigData = sig
		return nft.db.Transaction(func(tx *gorm.DB) error {
			err := tx.Model(&Collects{}).Create(&collectRec)
			if err.Error != nil {
				fmt.Println("NewCollections() err=", err.Error)
				return err.Error
			}
			return nil
		})
	}
	fmt.Println("NewCollections() dbase err=.", err)
	return err.Error
}

func (nft NftDb) ModifyCollections(useraddr, name, img, contract_type, contract_addr,
	desc, categories, sig string) error {
	useraddr = strings.ToLower(useraddr)
	contract_addr = strings.ToLower(contract_addr)
	var collectRec Collects
	err := nft.db.Where("Createaddr = ? AND name = ? ", useraddr, name).First(&collectRec)
	if err.Error != nil {
		fmt.Println("NewCollections() err=Collection not exist." )
		return ErrCollectionNotExist
	}
	collectRec = Collects{}
	if img != "" {
		collectRec.Img = img
	}
	if contract_type != "" {
		collectRec.Contracttype = contract_type
	}
	if contract_addr != "" {
		collectRec.Contract = contract_addr
	}
	if desc != "" {
		collectRec.Desc = desc
	}
	if categories != "" {
		collectRec.Categories = categories
	}
	if sig != "" {
		collectRec.SigData = sig
	}
	return nft.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&Collects{}).Where("Createaddr = ? AND name = ? ", useraddr, name).Updates(&collectRec)
		if err.Error != nil {
			fmt.Println("NewCollections() err=", err.Error)
			return err.Error
		}
		return nil
	})
}

func (nft NftDb) QueryNFTCollectionList(start_index, count string) ([]UserCollection, int, error) {
	var collectRecs []Collects
	var recCount int64
	if IsIntDataValid(start_index) != true {
		return nil, 0, ErrDataFormat
	}
	if IsIntDataValid(count) != true {
		return nil, 0, ErrDataFormat
	}
	err := nft.db.Model(Collects{}).Where("totalcount > 0").Count(&recCount)
	if err.Error != nil {
		fmt.Println("QueryNFTCollectionList() recCount err=", err)
		return nil, 0, ErrNftNotExist
	}
	startIndex, _ := strconv.Atoi(start_index)
	nftCount, _ := strconv.Atoi(count)
	if int64(startIndex) >= recCount || recCount == 0{
		return nil, 0, ErrNftNotExist
	} else {
		temp := recCount - int64(startIndex)
		if int64(nftCount) > temp {
			nftCount = int(temp)
		}
		err = nft.db.Model(Collects{}).Where("totalcount > 0").Order("transamt desc, id desc").Limit(nftCount).Offset(startIndex).Find(&collectRecs)
		if err.Error != nil {
			fmt.Println("QueryNFTCollectionList() find record err=", err)
			return nil, 0, ErrNftNotExist
		}
		userCollects := make([]UserCollection, 0, 20)
		for i := 0; i < len(collectRecs); i++ {
			var userCollect UserCollection
			userCollect.CreatorAddr = collectRecs[i].Createaddr
			userCollect.Name = collectRecs[i].Name
			//userCollect.Img = collectRecs[i].Img
			userCollect.ContractAddr = collectRecs[i].Contract
			userCollect.Desc = collectRecs[i].Desc
			//userCollect.Royalty = collectRecs[i].Royalty
			userCollect.Categories = collectRecs[i].Categories
			userCollects = append(userCollects, userCollect)
		}
		return userCollects, int(recCount), nil
	}
}