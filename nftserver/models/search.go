package models

import "gorm.io/gorm"

type SearchData struct {
	NftsRecords []Nfts			`json:"nfts"`
	CollectsRecords []Collects	`json:"collections"`
	UserAddrs []string			`json:"user_addrs"`
}
func (nft *NftDb) Search(cond string) ([]SearchData, error) {
	var searchData SearchData
	nfts := []Nfts{}
	findNftsResult := nft.db.Model(&Nfts{}).Where("name like ?", "%"+cond+"%").
		Order("name asc").Offset(0).Limit(5).Find(&nfts)
	if findNftsResult.Error != nil && findNftsResult.Error != gorm.ErrRecordNotFound {
		return nil, findNftsResult.Error
	}
	for k, _ := range nfts {
		nfts[k].Image = ""
	}
	searchData.NftsRecords = append(searchData.NftsRecords, nfts...)

	collects := []Collects{}
	findCollectsResult := nft.db.Model(&Collects{}).Where("createaddr like ? or name like ?", "%"+cond+"%", "%"+cond+"%").
		Order("name asc").Offset(0).Limit(5).Find(&collects)
	if findCollectsResult.Error != nil && findCollectsResult.Error != gorm.ErrRecordNotFound {
		return nil, findCollectsResult.Error
	}
	for k, _ := range collects {
		collects[k].Img = ""
	}
	searchData.CollectsRecords = append(searchData.CollectsRecords, collects...)

	users := []Users{}
	findUsersResult := nft.db.Model(&Users{}).Where("useraddr like ? or username like ?", "%"+cond+"%", "%"+cond+"%").
		Order("username asc").Offset(0).Limit(5).Find(&users)
	if findUsersResult.Error != nil && findUsersResult.Error != gorm.ErrRecordNotFound {
		return nil, findUsersResult.Error
	}

	searchData.UserAddrs = make([]string, 0)
	for _, user := range users {
		user.Portrait = ""
		user.Kycpic = ""
		searchData.UserAddrs = append(searchData.UserAddrs, user.Useraddr)
	}
	return []SearchData{searchData}, nil
}

