package models

import (
	"encoding/json"
	"gorm.io/gorm"
)

type HomePageReq struct {
	Announcement []string `json:"announcement"`
	NftLoop []NftLoopKey `json:"nft_loop"`
	Collections []CollectionKey `json:"collections"`
	NftList []NftKey `json:"nfts"`
}
type NftLoopKey struct {
	Contract string `json:"contract"`
	Tokenid string `json:"tokenid"`
}
type CollectionKey struct {
	Creator string `json:"creator"`
	Name string `json:"name"`
}
type NftKey struct {
	Contract string `json:"contract"`
	Tokenid string `json:"tokenid"`
}

type HomePageResp struct {
	Announcement []Announcements `json:"announcement"`
	NftLoop []Nfts `json:"nft_loop"`
	Collections []Collects `json:"collections"`
	NftList []Nfts `json:"nfts"`
	Total	int64 `json:"total"`
}

func (nft *NftDb) QueryHomePage() ([]HomePageResp, error) {
	sysParams := SysParams{}

	result := nft.db.Model(&SysParams{}).Last(&sysParams)
	if result.Error != nil {
		return nil, result.Error
	}

	var homePageReq HomePageReq
	err := json.Unmarshal([]byte(sysParams.Homepage), &homePageReq)
	if err != nil {
		return nil, err
	}

	homePageResp := HomePageResp{}

	announcementList, err := nft.QueryAnnouncement()
	if err != nil {
		return nil, err
	}
	homePageResp.Announcement = append(homePageResp.Announcement, announcementList...)

	for _, v := range homePageReq.NftLoop {
		nftData := Nfts{}
		result := nft.db.Model(&Nfts{}).Where("contract = ? and tokenid = ?", v.Contract, v.Tokenid).
			First(&nftData)
		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			return nil, result.Error
		}
		nftData.Image = ""
		homePageResp.NftLoop = append(homePageResp.NftLoop, nftData)
	}

	for _, v := range homePageReq.Collections {
		collectData := Collects{}
		result := nft.db.Model(&Collects{}).Where("createaddr = ? and name = ?", v.Creator, v.Name).
			First(&collectData)
		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			return nil, result.Error
		}
		collectData.Img = ""
		homePageResp.Collections = append(homePageResp.Collections, collectData)
	}

	for _, v := range homePageReq.NftList {
		nftData := Nfts{}
		result := nft.db.Model(&Nfts{}).Where("contract = ? and tokenid = ?", v.Contract, v.Tokenid).
			First(&nftData)
		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			return nil, result.Error
		}
		nftData.Image = ""
		homePageResp.NftList = append(homePageResp.NftList, nftData)
	}

	result = nft.db.Model(&Nfts{}).Count(&homePageResp.Total)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}

	return []HomePageResp{homePageResp}, nil
}
