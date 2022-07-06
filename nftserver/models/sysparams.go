package models

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
)

const (
	/*//"Admin" at 0x56c971ebBC0cD7Ba1f977340140297C0B48b7955
	//"NFT1155" at 0x53d76f1988B50674089e489B5ad1217AaC08CC85
	//"NFT721" at 0x5402AcE68556CC74aBB8861ceddc8F49401ac5D5
	//"TradeCore" at 0x3dE836C28a578da26D846f27353640582761909f
	initExchangAddr = "0x53d76f1988B50674089e489B5ad1217AaC08CC85"
	initNftAddr = "0x56c971ebBC0cD7Ba1f977340140297C0B48b7955"*/

	//"Admin" at 0x56c971ebBC0cD7Ba1f977340140297C0B48b7955
	//"NFT1155" at 0xA1e67a33e090Afe696D7317e05c506d7687Bb2E5
	//"TradeCore" at 0xD8D5D49182d7Abf3cFc1694F8Ed17742886dDE82

	initNFT1155Addr = "0xA1e67a33e090Afe696D7317e05c506d7687Bb2E5"
	initAdminAddr = "0x56c971ebBC0cD7Ba1f977340140297C0B48b7955"

	initNFT1155 = "0x53d76f1988B50674089e489B5ad1217AaC08CC85"
	initTrade = "0x3dE836C28a578da26D846f27353640582761909f"

	initLowprice = 1000000000
	initRoyaltylimit = 50 * 100
	SysRoyaltylimit = 50 * 100
	ZeroAddr = "0x0000000000000000000000000000000000000000"
	genTokenIdRetry = 20
	initCategories = "art,music,domain_names,virtual_worlds,trading_cards,collectibles,sports,utility"
)

var (
	TradeAddr string
	NFT1155Addr string
	AdminAddr string
	EthersNode string
	EthersWsNode string
	Lowprice uint64
	RoyaltyLimit int
)

type SysParamsRec struct {
	NFT1155addr		string		`json:"nft1155addr" gorm:"type:char(42) NOT NULL;comment:'nft1155合约地址'"`
	Adminaddr		string		`json:"adminaddr" gorm:"type:char(42) NOT NULL;comment:'管理员合约地址'"`
	Lowprice		uint64		`json:"lowprice" gorm:"type:bigint unsigned DEFAULT NULL;comment:'底价'"`
	Blocknumber		uint64		`json:"blocknumber" gorm:"type:bigint unsigned DEFAULT NULL;comment:'区块高度'"`
	Scannumber		uint64		`json:"scannumber" gorm:"type:bigint unsigned DEFAULT NULL;comment:'已扫描区块高度'"`
	Royaltylimit    int 		`json:"Royaltylimit" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'版税'"`
	Signdata		string		`json:"sig" gorm:"type:longtext NOT NULL;comment:'签名数据'"`
	Homepage		string		`json:"homepage" gorm:"type:longtext CHARACTER SET utf8mb4 NOT NULL;comment:'homepage数据'"`
	Categories		string 		`json:"categories" gorm:"type:longtext CHARACTER SET utf8mb4 NOT NULL;comment:'nft分类'"`
	Extend			string		`json:"extend" gorm:"type:longtext NOT NULL;comment:'扩展'"`
}

type SysParams struct {
	gorm.Model
	SysParamsRec
}

func (v SysParams) TableName() string {
	return "sysparams"
}

func (nft NftDb) QuerySysParams() (*SysParamsRec, error) {
	var params SysParams
	err := nft.db.Last(&params)
	if err.Error != nil {
		if err.Error == gorm.ErrRecordNotFound {
			params = SysParams{}
			params.NFT1155addr = strings.ToLower(NFT1155Addr)
			params.Adminaddr = strings.ToLower(AdminAddr)
			params.Lowprice = initLowprice
			params.Royaltylimit = initRoyaltylimit
			params.Categories = initCategories
			params.Blocknumber = GetCurrentBlockNumber()
			params.Scannumber = params.Blocknumber
			err = nft.db.Model(&SysParams{}).Create(&params)
			if err.Error != nil {
				fmt.Println("SetSysParams() create SysParams err= ", err.Error )
				return nil, err.Error
			}
		} else {
			fmt.Println("QuerySysParams() not find err=", err.Error)
			return nil, err.Error
		}
	}
	return &params.SysParamsRec, err.Error
}

func (nft NftDb) SetSysParams(param SysParamsRec) error {
	var paramRec, updateP SysParams
	err := nft.db.Last(&paramRec)
	if err.Error != nil {
		if nft.db.Error == gorm.ErrRecordNotFound {
			updateP.NFT1155addr = NFT1155Addr
			updateP.Adminaddr = AdminAddr
			updateP.Lowprice = initLowprice
			updateP.Royaltylimit = initRoyaltylimit
			updateP.Categories = initCategories
		} else {
			fmt.Println("QuerySysParams() not find err=", err.Error)
			return err.Error
		}
	} else {
		if param.NFT1155addr != "" {
			updateP.NFT1155addr = param.NFT1155addr
		} else{
			updateP.NFT1155addr = paramRec.NFT1155addr
		}
		if param.Adminaddr != "" {
			updateP.Adminaddr = param.Adminaddr
		} else {
			updateP.Adminaddr = paramRec.Adminaddr
		}
		if param.Lowprice != 0 {
			updateP.Lowprice = param.Lowprice
		} else {
			updateP.Lowprice = paramRec.Lowprice
		}
	}
	updateP.Signdata = param.Signdata
	err = nft.db.Model(&SysParams{}).Create(&updateP)
	if err.Error != nil {
		fmt.Println("SetSysParams() create SysParams err= ", err.Error )
		return err.Error
	}
	return nil
}

func InitSysParams(Sqldsndb string) error {
	nd, err := NewNftDb(Sqldsndb)
	defer nd.Close()
	if err != nil {
		fmt.Printf("connect database err = %s\n", err)
		return err
	}
	params, err := nd.QuerySysParams()
	if err != nil {
		return err
	} else {
		Lowprice = params.Lowprice
		RoyaltyLimit = params.Royaltylimit
	}
	return nil
}

