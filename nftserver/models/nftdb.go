package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/nftexchange/nftserver/ethhelper"
	"golang.org/x/crypto/sha3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
	//"github.com/nftexchange/ethhelper"
)

var (
	ErrNftAlreadyExist = errors.New("nft already exist.")
	ErrNftNotExist = errors.New("nft Not exist.")
	ErrNftDelete = errors.New("nft delete error.")
	ErrNftIllegalMint = errors.New("nft mint illegal.")
	ErrAlreadyNftFavorited = errors.New("already NftFavorited.")
	ErrAlreadyUserFavorited = errors.New("already UserFavorited.")
	ErrNotNftFavorited = errors.New("not NftFavorited.")
	ErrNftSelling = errors.New("nft be selling.")
	ErrNftNotSell = errors.New("nft not on sale.")
	ErrNftNotMinted = errors.New("nft not Minted.")
	ErrAlreadyBid = errors.New("Already bid.")
	ErrAuctionEnd = errors.New("The auction ended.")
	ErrAuctionNotBegan = errors.New("The auction not began.")
	ErrBidOutRange = errors.New("Bid is out-of-range.")
	ErrNotVerify = errors.New("Not verify.")
	ErrAuctionStartAfterEnd = errors.New("start time > end time.")
	ErrAuctionStartBeforeNow = errors.New("start time > end time.")
	ErrNoRightSell = errors.New("have no right to sell.")
	ErrFromToAddrZero = errors.New("from or to addr = 0.")
	ErrNoAuthorize = errors.New("No authorize.")
	ErrAuthorizeLess = errors.New("Less authorize amount.")
	ErrBalanceLess = errors.New("Less balance amount.")
	ErrCollectionExist = errors.New("Collection already exist.")
	ErrCollectionNotExist = errors.New("Collection not exist.")
	ErrNftUpAddrNotOwn = errors.New("Nft upload address not own.")
	ErrNftUpAddrNotAdmin = errors.New("Nft upload address not Admin.")
	ErrNftNoMore = errors.New("Nft no more.")
	ErrGenerateTokenId = errors.New("generate token id error.")
	ErrContractCountLtZero = errors.New("contract count < 0.")
	ErrNoTrans = errors.New("no trade.")
	ErrNoCategory = errors.New("category err.")
	ErrPrice = errors.New("Price error.")
	ErrDataFormat = errors.New("data format error.")
	ErrRoyalty = errors.New("royalty too big error.")
	ErrBuyOwn = errors.New("buy your own nft.")
	ErrTransExist = errors.New("transaction exist.")
)

//const (
//	/*//"Admin" at 0x56c971ebBC0cD7Ba1f977340140297C0B48b7955
//	//"NFT1155" at 0x53d76f1988B50674089e489B5ad1217AaC08CC85
//	//"NFT721" at 0x5402AcE68556CC74aBB8861ceddc8F49401ac5D5
//	//"TradeCore" at 0x3dE836C28a578da26D846f27353640582761909f
//	initExchangAddr = "0x53d76f1988B50674089e489B5ad1217AaC08CC85"
//	initNftAddr = "0x56c971ebBC0cD7Ba1f977340140297C0B48b7955"*/
//
//	//"Admin" at 0x56c971ebBC0cD7Ba1f977340140297C0B48b7955
//	//"NFT1155" at 0xA1e67a33e090Afe696D7317e05c506d7687Bb2E5
//	//"TradeCore" at 0xD8D5D49182d7Abf3cFc1694F8Ed17742886dDE82
//
//	initExchangAddr = "0xA1e67a33e090Afe696D7317e05c506d7687Bb2E5"
//	initNftAddr = "0x56c971ebBC0cD7Ba1f977340140297C0B48b7955"
//
//	initNFT1155 = "0x53d76f1988B50674089e489B5ad1217AaC08CC85"
//	initTrade = "0x3dE836C28a578da26D846f27353640582761909f"
//	initLowprice = 1000000000
//	initRoyaltylimit = 50 * 100
//	SysRoyaltylimit = 50 * 100
//	ZeroAddr = "0x0000000000000000000000000000000000000000"
//	genTokenIdRetry = 20
//	initCategories = "art,music,domain_names,virtual_worlds,trading_cards,collectibles,sports,utility"
//)

//var (
//	ExchangAddr string
//	NftAddr string
//	Lowprice uint64
//	RoyaltyLimit int
//)

type SellState int
const (
	SellStateStart SellState = iota
	SellStateWait
)

func (this SellState) String() string {
	switch this {
	case SellStateStart:
		return "SellStart"
	case SellStateWait:
		return "SellWait"
	default:
		return "Unknow"
	}
}
type SellType int
const (
	SellTypeNotSale SellType = iota
	SellTypeSetPrice
	SellTypeFixPrice
	SellTypeForeignPrice
	SellTypeHighestBid
	SellTypeBidPrice
	SellTypeMintNft
	SellTypeForeignMint
	SellTypeDelNft
	SellTypeWaitSale
	SellTypeAsset
	SellTypeError
)

func (this SellType) String() string {
	switch this {
	case SellTypeNotSale:
		return "NotSale"
	case SellTypeSetPrice:
		return "SetPrice"
	case SellTypeFixPrice:
		return "FixPrice"
	case SellTypeHighestBid:
		return "HighestBid"
	case SellTypeBidPrice:
		return "BidPrice"
	case SellTypeForeignPrice:
		return "ForeignPrice"
	case SellTypeMintNft:
		return "MintNft"
	case SellTypeForeignMint:
		return "ForeignMint"
	case SellTypeDelNft:
		return "DelNft"
	case SellTypeWaitSale:
		return "WaitSale"
	case SellTypeAsset:
		return "AssetTransfer"
	case SellTypeError:
		return "Error"
	default:
		return "Unknow"
	}
}

type Userrec struct {
	Useraddr	string		`json:"useraddr" gorm:"type:char(42) NOT NULL;comment:'用户地址'"`
	Signdata	string		`json:"sig" gorm:"type:longtext NOT NULL;comment:'签名数据'"`
	Username 	string 		`json:"user_name" gorm:"type:char(200) CHARACTER SET utf8mb4 NOT NULL;comment:'用户名称'"`
	Country 	string 		`json:"country" gorm:"type:char(200) CHARACTER SET utf8mb4 NOT NULL;comment:'国籍'"`
	Bio			string		`json:"user_info" gorm:"type:longtext CHARACTER SET utf8mb4 NOT NULL;comment:'用户简介'"`
	Portrait	string		`json:"portrait" gorm:"type:longtext NOT NULL;comment:'用户头像'"`
	Kycpic		string		`json:"kycpic" gorm:"type:longtext NOT NULL;comment:'kyc审核照片'"`
	Email		string		`json:"user_mail" gorm:"type:longtext NOT NULL;comment:'用户邮箱'"`
	Userregd	int64		`json:"userregd" gorm:"type:bigint DEFAULT NULL;comment:'用户注册时间'"`
	Userlogin	int64		`json:"userlogin" gorm:"type:bigint DEFAULT NULL;comment:'用户登录时间'"`
	Userlogout	int64		`json:"userlogout" gorm:"type:bigint DEFAULT NULL;comment:'用户登出时间'"`
	Verified	string 		`json:"verified" gorm:"type:char(20)  DEFAULT NULL;comment:'是否通过审核'"`
	Verifyaddr	string		`json:"vrf_addr" gorm:"type:char(42) NOT NULL;comment:'验证人地址'"`
	Desc		string		`json:"desc" gorm:"type:longtext CHARACTER SET utf8mb4 NOT NULL;comment:'审核描述：未通过审核描述'"`
	Favorited	int			`json:"favorited" gorm:"type:int unsigned zerofill DEFAULT NULL;comment:'被关注计数'"`
}

type Users struct {
	gorm.Model
	Userrec
}

func (v Users) TableName() string {
	return "users"
}

type Announce struct {
	Title		string		`json:"title" gorm:"type:longtext CHARACTER SET utf8mb4 NOT NULL;comment:'主题'"`
	Content		string		`json:"content" gorm:"type:longtext CHARACTER SET utf8mb4 NOT NULL;comment:'内容'"`
}

type Announcements struct {
	gorm.Model
	Announce
}

func (v Announcements) TableName() string {
	return "announcements"
}

//type SysParamsRec struct {
//	Exchangaddr		string		`json:"exchangaddr" gorm:"type:char(42) NOT NULL;comment:'交易所合约地址'"`
//	Nftaddr			string		`json:"exchangaddr" gorm:"type:char(42) NOT NULL;comment:'nft合约地址'"`
//	Lowprice		uint64		`json:"lowprice" gorm:"type:bigint unsigned DEFAULT NULL;comment:'底价'"`
//	Royaltylimit    int 		`json:"Royaltylimit" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'版税'"`
//	Signdata		string		`json:"sig" gorm:"type:longtext NOT NULL;comment:'签名数据'"`
//	Homepage		string		`json:"homepage" gorm:"type:longtext CHARACTER SET utf8mb4 NOT NULL;comment:'homepage数据'"`
//	Categories		string 		`json:"categories" gorm:"type:longtext CHARACTER SET utf8mb4 NOT NULL;comment:'nft分类'"`
//	Extend			string		`json:"extend" gorm:"type:longtext NOT NULL;comment:'扩展'"`
//}
//
//type SysParams struct {
//	gorm.Model
//	SysParamsRec
//}
//
//func (v SysParams) TableName() string {
//	return "sysparams"
//}

type Sigmsgrec struct {
	Signdata	string		`json:"sig" gorm:"type:longtext NOT NULL;comment:'签名'"`
	Signmsg		string		`json:"sigmsg" gorm:"type:longtext NOT NULL;comment:'原始数据'"`
}

type Sigmsgs struct {
	gorm.Model
	Sigmsgrec
}

func (v Sigmsgs) TableName() string {
	return "sigmsgs"
}

type Verified int
const (
	NoVerify	Verified = iota
	Passed
	NoPass
)

func (this Verified) String() string {
	switch this {
	case NoVerify:
		return "NoVerify"
	case NoPass:
		return "NoPass"
	case Passed:
		return "Passed"
	default:
		return "Unknow"
	}
}

type NftRecord struct {
	Ownaddr			string		`json:"ownaddr" gorm:"type:char(42) ;comment:'nft拥有者地址'"`
	Md5				string		`json:"md5" gorm:"type:longtext ;comment:'图片md5值'"`
	Name			string 		`json:"name" gorm:"type:varchar(200) CHARACTER SET utf8mb4 ;comment:'nft名称'"`
	Desc			string		`json:"desc" gorm:"type:longtext CHARACTER SET utf8mb4  ;comment:'nft描述'"`
	Meta			string		`json:"meta" gorm:"type:longtext CHARACTER SET utf8mb4  ;comment:'元信息'"`
	Url				string		`json:"source_url" gorm:"type:varchar(200) DEFAULT NULL;comment:'nfc原始数据保持地址'"`
	Contract		string		`json:"nft_contract_addr" gorm:"type:char(42) ;comment:'合约地址'"`
	Tokenid			string 		`json:"nft_token_id" gorm:"type:char(42) ;comment:'唯一标识nft标志'"`
	Count	     	int 		`json:"count" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'nft可卖数量'"`
	Approve			string		`json:"approve" gorm:"type:longtext ;comment:'授权'"`
	Categories		string 		`json:"categories" gorm:"type:varchar(200) CHARACTER SET utf8mb4 ;comment:'nft分类'"`
	Collectcreator	string		`json:"collection_creator_addr" gorm:"type:char(42) ;comment:'合集创建者地址'"`
	Collections 	string  	`json:"collections" gorm:"type:varchar(200) CHARACTER SET utf8mb4 ;comment:'NFT合集名'"`
	Image			string		`json:"asset_sample" gorm:"type:longtext ;comment:'缩略图二进制数据'"`
	Hide			string		`json:"hide" gorm:"type:char(20) ;comment:'是否让其他人看到'"`
	Signdata		string		`json:"sig" gorm:"type:longtext ;comment:'签名数据，创建时产生'"`
	Createaddr		string		`json:"user_addr" gorm:"type:char(42) ;comment:'创建nft地址'"`
	Verifyaddr		string		`json:"vrf_addr" gorm:"type:char(42) ;comment:'验证人地址'"`
	Currency    	string 		`json:"currency" gorm:"type:varchar(10) CHARACTER SET utf8mb4 DEFAULT NULL;COMMENT:'交易币种'"`
	Price			uint64		`json:"price" gorm:"type:bigint unsigned DEFAULT NULL;comment:'创建时定的价格'"`
	Royalty     	int 		`json:"royalty" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'版税'"`
	Paychan    		string 		`json:"paychan" gorm:"type:char(20) DEFAULT NULL;COMMENT:'交易通道'"`
	TransCur    	string 		`json:"trans_cur" gorm:"type:char(20) CHARACTER SET utf8mb4 DEFAULT NULL;COMMENT:'交易币种'"`
	Transprice		uint64		`json:"transprice" gorm:"type:bigint unsigned DEFAULT NULL;comment:'交易成交价格'"`
	Transtime		int64		`json:"last_trans_time" gorm:"type:bigint DEFAULT NULL;comment:'最后交易时间'"`
	Createdate		int64		`json:"createdate" gorm:"type:bigint DEFAULT NULL;comment:'nft创建时间'"`
	Favorited		int			`json:"favorited" gorm:"type:int unsigned zerofill DEFAULT 0;comment:'被关注计数'"`
	Transcnt		int			`json:"transcnt" gorm:"type:int unsigned zerofill DEFAULT NULL;comment:'交易次数，每交易一次加一'"`
	Transamt		uint64		`json:"transamt" gorm:"type:bigint DEFAULT NULL;comment:'交易总金额'"`
	Verified		string 		`json:"verified" gorm:"type:char(20) DEFAULT NULL;comment:'nft作品是否通过审核'"`
	Verifieddesc	string		`json:"Verifieddesc" gorm:"type:longtext CHARACTER SET utf8mb4  ;comment:'审核描述：未通过审核描述'"`
	Verifiedtime	int64		`json:"vrf_time" gorm:"type:bigint DEFAULT NULL;comment:'审核时间'"`
	Selltype    	string      `json:"selltype" gorm:"type:char(20) DEFAULT NULL;COMMENT:'nft交易类型'"`
	//Sellprice		uint64		`json:"sellingprice" gorm:"type:bigint unsigned DEFAULT NULL;comment:'正在销售价格'"`
	Mintstate   	string      `json:"mintstate" gorm:"type:char(20) DEFAULT NULL;COMMENT:'铸币状态'"`
	Extend			string		`json:"extend" gorm:"type:longtext ;comment:'扩展字段'"`
}

type Nfts struct {
	gorm.Model
	NftRecord
}

func (v Nfts) TableName() string {
	return "nfts"
}

type ContractType int
const (
	ERC1155	ContractType = iota
	ERC721
)

func (this ContractType) String() string {
	switch this {
	case ERC1155:
		return "ERC1155"
	case ERC721:
		return "ERC721"
	default:
		return "Unknow"
	}
}

type CollectRec struct {
	Createaddr		string		`json:"collection_creator_addr" gorm:"type:char(42) NOT NULL;comment:'创建者者地址'"`
	Contract    	string		`json:"nft_contract_addr" gorm:"type:char(42) NOT NULL;comment:'合约地址'"`
	Contracttype	string		`json:"contracttype" gorm:"type:char(20) CHARACTER SET utf8mb4 NOT NULL;comment:'合约类型'"`
	Name			string		`json:"name" gorm:"type:varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL;comment:'集合名称'"`
	Desc			string		`json:"desc" gorm:"type:longtext CHARACTER SET utf8mb4 NOT NULL;comment:'集合描述'"`
	Categories		string 		`json:"categories" gorm:"type:char(200) CHARACTER SET utf8mb4 NOT NULL;comment:'集合分类'"`
	Totalcount		int			`json:"total_count" gorm:"type:int unsigned zerofill DEFAULT 0;comment:'集合中nft总数'"`
	Transcnt		int			`json:"transcnt" gorm:"type:int unsigned zerofill DEFAULT NULL;comment:'交易次数，每交易一次加一'"`
	Transamt		uint64		`json:"transamt" gorm:"type:bigint DEFAULT NULL;comment:'交易总金额'"`
	SigData			string		`json:"sig" gorm:"type:longtext NOT NULL;comment:'签名'"`
	Img				string		`json:"img" gorm:"type:longtext NOT NULL;comment:'logo图片'"`
	Extend			string		`json:"extend" gorm:"type:longtext NOT NULL;comment:'扩展字段'"`
}

type Collects struct {
	gorm.Model
	CollectRec
}

func (v Collects) TableName() string {
	return "collects"
}

type CollectListRec struct {
	Collectsid	uint		`json:"collectid" gorm:"type:bigint unsigned DEFAULT NULL;comment:'集合索引'"`
	Nftid		uint		`json:"nftid" gorm:"type:bigint unsigned DEFAULT NULL;comment:'nft索引'"`

}

type CollectLists struct {
	gorm.Model
	CollectListRec
}

func (v CollectLists) TableName() string {
	return "collectlists"
}

type TranRecord struct {
	Auctionid	uint		`json:"auctionid" gorm:"type:bigint DEFAULT NULL;COMMENT:'出价索引'"`
	Contract	string		`json:"contract" gorm:"type:char(42) NOT NULL;comment:'合约地址'"`
	Createaddr	string		`json:"user_addr" gorm:"type:char(42) NOT NULL;comment:'创建nft地址'"`
	Fromaddr	string		`json:"fromaddr" gorm:"type:char(42) NOT NULL;comment:'卖家地址'"`
	Toaddr		string		`json:"toaddr" gorm:"type:char(42) NOT NULL;comment:'买家地址'"`
	Tradesig	string		`json:"tradesig" gorm:"type:longtext NOT NULL;comment:'交易签名'"`
	Signdata	string		`json:"signdata" gorm:"type:longtext NOT NULL;comment:'签名数据，创建时产生'"`
	Txhash		string		`json:"txhash" gorm:"type:longtext NOT NULL;comment:'交易签名'"`
	Tokenid		string 		`json:"nft_token_id" gorm:"type:char(42) NOT NULL;comment:'唯一标识nft标志'"`
	Count     	int 		`json:"count" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'nft可卖数量'"`
	Transtime	int64		`json:"transtime" gorm:"type:bigint DEFAULT NULL;comment:'nft创建时间'"`
	Nftid		uint		`json:"nftid" gorm:"type:int DEFAULT NULL;COMMENT:'nft索引'"`
	Paychan     string 		`json:"paychan" gorm:"type:char(20) CHARACTER SET utf8mb4 DEFAULT NULL;COMMENT:'交易通道'"`
	Currency    string 		`json:"currency" gorm:"type:char(20) CHARACTER SET utf8mb4 DEFAULT NULL;COMMENT:'交易币种'"`
	Price		uint64		`json:"price" gorm:"bigint unsigned DEFAULT NULL;comment:'成交价格'"`
	Name 		string		`json:"name" gorm:"type:varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL;COMMENT:'nft名称'"`
	Desc		string		`json:"desc" gorm:"type:longtext CHARACTER SET utf8mb4 NOT NULL;comment:'审核描述：未通过审核描述'"`
	Meta		string		`json:"meta" gorm:"type:longtext NOT NULL;comment:'元信息'"`
	Selltype    string      `json:"selltype" gorm:"type:char(20) DEFAULT NULL;COMMENT:'nft交易类型'"`
	Error	    string      `json:"error" gorm:"type:char(200) DEFAULT NULL;COMMENT:'nft交易出错原因'"`
}

type Trans struct {
	gorm.Model
	TranRecord
}

func (v Trans) TableName() string {
	return "trans"
}

type MintState int
const (
	NoMinted MintState = iota
	Minted
	Minting
)

func (this MintState) String() string {
	switch this {
		case NoMinted:
			return "NoMinted"
		case Minted:
			return "Minted"
		case Minting:
			return "Minting"
		default:
			return "Unknow"
	}
}

type AuctionRecord struct {
	Selltype 	string 		`json:"selltype" gorm:"type:char(20) DEFAULT NULL;COMMENT:'拍卖类型'"`
	Ownaddr		string		`json:"ownaddr" gorm:"type:char(42) NOT NULL;comment:'nft拥有者地址'"`
	Privaddr	string		`json:"privaddr" gorm:"type:char(42) NOT NULL;comment:'nft拥有者地址'"`
	//Privaddrid	uint		`gorm:"type:int(11) DEFAULT NULL;COMMENT:'指定买家地址索引'"`
	Nftid		uint		`json:"nftid" gorm:"type:int DEFAULT NULL;COMMENT:'拍卖nft索引'"`
	Tokenid		string 		`json:"nft_token_id" gorm:"type:char(42) NOT NULL;comment:'唯一标识nft标志'"`
	Count	   	int 		`json:"count" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'nft可卖数量'"`
	Contract	string		`json:"nft_contract_addr" gorm:"type:char(42) NOT NULL;comment:'合约地址'"`
	Paychan     string 		`json:"paychan" gorm:"type:char(20) CHARACTER SET utf8mb4 DEFAULT NULL;COMMENT:'交易通道'"`
	Currency    string 		`json:"currency" gorm:"type:char(20) CHARACTER SET utf8mb4 DEFAULT NULL;COMMENT:'交易币种'"`
	Startprice	uint64		`json:"startprice" gorm:"type:bigint unsigned DEFAULT NULL;COMMENT:'起拍价'"`
	Endprice	uint64		`json:"endprice" gorm:"type:bigint unsigned DEFAULT NULL;COMMENT:'结束价'"`
	Startdate	int64		`json:"startdate" gorm:"type:bigint DEFAULT NULL;comment:'拍卖开始时间'"`
	Enddate		int64		`json:"enddate" gorm:"type:bigint DEFAULT NULL;comment:'拍卖结束时间'"`
	Tradesig	string		`json:"tradesig" gorm:"type:longtext NOT NULL;comment:'交易签名'"`
	Signdata	string		`json:"sig" gorm:"type:longtext NOT NULL;comment:'签名数据'"`
	Toaddr		string		`json:"toaddr" gorm:"type:char(42) NOT NULL;comment:'nft拥有者地址'"`
	Price		uint64		`json:"price" gorm:"bigint unsigned DEFAULT NULL;comment:'成交价格'"`
	Blocknumber	int64		`json:"blocknumber" gorm:"type:bigint DEFAULT NULL;comment:'售卖时的区块高度'"`
	Txhash		string		`json:"txhash" gorm:"type:longtext NOT NULL;comment:'交易签名'"`
	SellState 	string 		`json:"sellstate" gorm:"type:char(20) DEFAULT NULL;COMMENT:'售卖状态'"`
}

type Auction struct {
	gorm.Model
	AuctionRecord
}

func (v Auction) TableName() string {
	return "auctions"
}

type AuctionHistory struct {
	gorm.Model
	AuctionRecord
}

func (v AuctionHistory) TableName() string {
	return "auctionhistorys"
}

type BidRecord struct {
	Bidaddr		string		`json:"bidaddr" gorm:"varchar(42) DEFAULT NULL;COMMENT:'出价客户地址'"`
	Auctionid	uint		`json:"auctionid" gorm:"type:int DEFAULT NULL;COMMENT:'拍卖索引'"`
	Tokenid		string 		`json:"nft_token_id" gorm:"type:char(42) NOT NULL;comment:'唯一标识nft标志'"`
	Count     	int 		`json:"count" gorm:"type:int unsigned zerofill DEFAULT 0;COMMENT:'nft可卖数量'"`
	Contract	string		`json:"nft_contract_addr" gorm:"type:char(42) NOT NULL;comment:'合约地址'"`
	Paychan     string 		`json:"paychan" gorm:"type:char(20) CHARACTER SET utf8mb4 DEFAULT NULL;COMMENT:'交易通道'"`
	Currency    string 		`json:"currency" gorm:"type:char(20) CHARACTER SET utf8mb4 DEFAULT NULL;COMMENT:'交易币种'"`
	Price		uint64		`json:"price" gorm:"bigint unsigned DEFAULT NULL;COMMENT:'出价'"`
	Nftid		uint		`json:"nftid" gorm:"type:int DEFAULT NULL;COMMENT:'nft索引'"`
	Bidtime		int64		`json:"bidtime" gorm:"bigint DEFAULT NULL;COMMENT:'出价时间'"`
	Deadtime	int64		`json:"dead_time" gorm:"bigint DEFAULT NULL;COMMENT:'出价失效时间'"`
	Tradesig	string		`json:"tradesig" gorm:"type:longtext NOT NULL;comment:'交易签名'"`
	Signdata	string		`json:"sig" gorm:"type:longtext NOT NULL;comment:'签名数据，创建时产生'"`
}

type Bidding struct {
	gorm.Model
	BidRecord
}

func (v Bidding) TableName() string {
	return "biddings"
}

type BiddingHistory struct {
	gorm.Model
	BidRecord
}

func (v BiddingHistory) TableName() string {
	return "biddinghistorys"
}

type Collected struct {
	gorm.Model
	Nftid		string		`gorm:"type:int DEFAULT NULL;COMMENT:'nft索引'"`
}

type Created struct {
	gorm.Model
	//Nftid		string		`gorm:"type:int DEFAULT NULL;COMMENT:'nft索引'"`
	NftRecord
}

type NftFavoriteRec struct {
	Useraddr		string		`json:"user_addr" gorm:"type:char(42) NOT NULL;comment:'关注人地址'"`
	Tokenid			string 		`json:"nft_token_id" gorm:"type:char(42) NOT NULL;comment:'唯一标识nft标志'"`
	Contract    	string		`json:"nft_contract_addr" gorm:"type:char(42) NOT NULL;comment:'合约地址'"`
	Name 			string		`json:"name" gorm:"type:varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL;COMMENT:'nft名称'"`
	Image			string		`json:"asset_sample" gorm:"type:longtext NOT NULL;comment:'缩略图二进制数据'"`
	Img				string		`json:"img" gorm:"type:longtext NOT NULL;comment:'logo图片'"`
	Collectcreator	string		`json:"collection_creator_addr" gorm:"type:char(42) NOT NULL;comment:'合集创建者地址'"`
	Collections 	string  	`json:"collections" gorm:"type:varchar(200) CHARACTER SET utf8mb4 NOT NULL;comment:'NFT合集名'"`
	Signdata		string		`json:"sig" gorm:"type:longtext NOT NULL;comment:'签名数据，创建时产生'"`
	Nftid			uint		`json:"nftid" gorm:"type:bigint DEFAULT NULL;COMMENT:'nft索引'"`
}

type NftFavorited struct {
	gorm.Model
	NftFavoriteRec
}

func (v NftFavorited) TableName() string {
	return "nftfavoriteds"
}

type UserFavorited struct {
	gorm.Model
	Useraddr 		string		`gorm:"type:char(42) NOT NULL;comment:'关注者地址'"`
	Favoritedaddr 	string		`gorm:"type:char(42) NOT NULL;comment:'被关注者地址'"`
	//Nftid			uint		`gorm:"type:int DEFAULT 0;COMMENT:'nft索引'"`
}

func (v UserFavorited) TableName() string {
	return "userfavoriteds"
}

type Exchange struct {
	gorm.Model
	addr 		string		`json:"addr" gorm:"type:char(42) NOT NULL;comment:'地址'"`
}

func (v Exchange) TableName() string {
	return "exchangs"
}

type NftDb struct {
	db *gorm.DB
}

type Portrait struct {
	Useraddr	string		`json:"useraddr" gorm:"type:char(42) NOT NULL;comment:'用户地址'"`
	Portrait	string		`json:"portrait" gorm:"type:longtext NOT NULL;comment:'用户头像'"`
}

type NftTransInfo struct {
	Nft     NftRecord     `json:"nft"`
	Auction AuctionRecord `json:"auction"`
	Trans   []TranRecord  `json:"trans"`
	Bids    []BidRecord   `json:"bids"`
	Sigs    []Sigmsgrec   `json:"sigs"`
	Portraits []Portrait  `json:"portraits"`
}

type StQueryField struct {
	Field string `json:"field"`
	Operation string `json:"operation"`
	Value string `json:"value"`
}
type StSortField struct {
	By string `json:"by"`
	Order string `json:"order"`
}

//const sqlsvrLocal = "demo:123456@tcp(192.168.56.128:3306)/"
//const vpnsvr = "demo:123456@tcp(192.168.1.238:3306)/"
var SqlSvr string
//const dbName = "nftdb"
var dbName string
const localtime = "?parseTime=true&loc=Local"

func (nft *NftDb) ConnectDB(sqldsn string) error {
	var err error
	//nft.db, err = gorm.Open("mysql", sqldsn)
	nft.db, err = gorm.Open(mysql.Open(sqldsn), &gorm.Config{})
	if err != nil {
		fmt.Println("failed to connect database", err.Error())
	}
	return err
}

func NewNftDb(sqldsn string) (*NftDb, error) {
	nft := new(NftDb)
	var err error
	nft.db, err = gorm.Open(mysql.Open(sqldsn), &gorm.Config{})
	if err != nil {
		fmt.Println("failed to connect database", err.Error())
	}
	return nft, err
}

func (nft NftDb) Close()  {
	//nft.db.Close()
	sqlDB, _ := nft.db.DB()
	sqlDB.Close()
}

func (nft *NftDb) GetDB() *gorm.DB{
	return nft.db
}

func (nft NftDb) createDb(dbName string) error {
	strOrder := "create database if not exists " + dbName +";"
	db := nft.db.Exec(strOrder)
	if db.Error != nil {
		fmt.Printf("CreateDataBase err=%s\n", db.Error)
		return db.Error
	}
	strOrder = "use " + dbName
	db = nft.db.Exec(strOrder)
	if db.Error != nil {
		fmt.Printf("use database err=%s\n", db.Error)
	}
	return db.Error
}

func (nft NftDb) InitDb(sqlsvr string, dbName string) error  {
	err := nft.ConnectDB(sqlsvr)
	if err != nil {
		fmt.Printf("InitDb()->connectDb() err=%s\n", err)
		return err
	}
	err = nft.createDb(dbName)
	if err != nil {
		fmt.Printf("Create Db err=%s\n", err)
		return err
	}
	err = nft.db.AutoMigrate(&Users{})
	if err != nil {
		fmt.Println("create table Users{} err=", err)
	}
	err = nft.db.AutoMigrate(&Nfts{})
	if err != nil {
		fmt.Println("create table Nfts{} err=", err)
	}
	err = nft.db.AutoMigrate(&Trans{})
	if err != nil {
		fmt.Println("create table Trans{} err=", err)
	}
	err = nft.db.AutoMigrate(&Auction{})
	if err != nil {
		fmt.Println("create table Auction{} err=", err)
	}
	err = nft.db.AutoMigrate(&AuctionHistory{})
	if err != nil {
		fmt.Println("create table AuctionHistory{} err=", err)
	}
	err = nft.db.AutoMigrate(&Bidding{})
	if err != nil {
		fmt.Println("create table Bidding{} err=", err)
	}
	err = nft.db.AutoMigrate(&BiddingHistory{})
	if err != nil {
		fmt.Println("create table BiddingHistory{} err=", err)
	}
	err = nft.db.AutoMigrate(&NftFavorited{})
	if err != nil {
		fmt.Println("create table NftFavorited{} err=", err)
	}
	err = nft.db.AutoMigrate(&UserFavorited{})
	if err != nil {
		fmt.Println("create table UserFavorited{} err=", err)
	}
	err = nft.db.AutoMigrate(&Sigmsgs{})
	if err != nil {
		fmt.Println("create table Sigmsg{} err=", err)
	}
	err = nft.db.AutoMigrate(&SysParams{})
	if err != nil {
		fmt.Println("create table SysParams{} err=", err)
	}
	err = nft.db.AutoMigrate(&Collects{})
	if err != nil {
		fmt.Println("create table Collects{} err=", err)
	}
	err = nft.db.AutoMigrate(&CollectLists{})
	if err != nil {
		fmt.Println("create table CollectLists{} err=", err)
	}
	err = nft.db.AutoMigrate(&Announcements{})
	if err != nil {
		fmt.Println("create table Announcements{} err=", err)
	}
	nft.Close()
	return err
}

//func (nft NftDb) Login(userAddr, sigData string) error {
//	userAddr = strings.ToLower(userAddr)
//	user := Users{}
//	db := nft.db.Model(&user).Where("useraddr = ?", userAddr).First(&user)
//	if db.Error != nil {
//		if db.Error == gorm.ErrRecordNotFound {
//			user.Useraddr = userAddr
//			user.Signdata = sigData
//			user.Userlogin = time.Now().Unix()
//			user.Userlogout = time.Now().Unix()
//			user.Username = ""
//			user.Userregd = time.Now().Unix()
//			db = nft.db.Model(&user).Create(&user)
//			if db.Error != nil {
//				fmt.Println("loging()->create() err=", db.Error)
//				return db.Error
//			}
//		}
//	} else {
//		db = nft.db.Model(&Users{}).Where("useraddr = ?", userAddr).Update("userlogin", time.Now().Unix())
//		if db.Error != nil {
//			fmt.Printf("login()->UPdate() users err=%s\n", db.Error)
//		}
//	}
//	return db.Error
//}

func IsAdminAddr(userAddr string) (bool, error) {
	adminAddrs, err := ethhelper.AdminList()
	if err != nil {
		fmt.Println("IsAdminAddr() get admin addr err=", err)
		return false, err
	}
	userAddr = userAddr[2:]
	var IsAdminAddr bool
	for _, addr := range adminAddrs {
		if addr == userAddr {
			IsAdminAddr = true
			break
		}
	}
	return IsAdminAddr, nil
}

//func (nft NftDb) UploadNft(
//	user_addr string,
//	creator_addr string,
//	owner_addr string,
//	md5 string,
//	name string,
//	desc string,
//	meta string,
//	source_url string,
//	nft_contract_addr string,
//	nft_token_id string,
//	categories string,
//	collections string,
//	asset_sample string,
//	hide string,
//	royalty string,
//	count string,
//	sig string) error {
//	user_addr = strings.ToLower(user_addr)
//	creator_addr = strings.ToLower(creator_addr)
//	owner_addr = strings.ToLower(owner_addr)
//	nft_contract_addr = strings.ToLower(nft_contract_addr)
//
//	if IsIntDataValid(count) != true {
//		return ErrDataFormat
//	}
//	if IsIntDataValid(royalty) != true {
//		return ErrDataFormat
//	}
//	r, _ := strconv.Atoi(royalty)
//	r = r / 100
//	fmt.Println("UploadNft() royalty=", r, "SysRoyaltylimit=", SysRoyaltylimit, "RoyaltyLimit", RoyaltyLimit )
//	if r > SysRoyaltylimit || r > RoyaltyLimit {
//		return ErrRoyalty
//	}
//	if count == "" {
//		count = "1"
//	}
//	if c, _ := strconv.Atoi(count); c < 1 {
//		fmt.Println("UploadNft() contract count < 1.")
//		return ErrContractCountLtZero
//	}
//	if nft.IsValidCategory(categories) {
//		return ErrNoCategory
//	}
//
//	var collectRec Collects
//	if collections != "" {
//		err := nft.db.Where("name = ? AND createaddr =?",
//			collections, creator_addr).First(&collectRec)
//		if err.Error != nil {
//			fmt.Println("UploadNft() err=Collection not exist.")
//			return ErrCollectionNotExist
//		}
//	} else {
//		return ErrCollectionNotExist
//	}
//	if nft_contract_addr == "" && nft_token_id == "" {
//		var NewTokenid string
//		rand.Seed(time.Now().UnixNano())
//		var i int
//		for i = 0; i < genTokenIdRetry ; i++ {
//			//NewTokenid := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
//			s := fmt.Sprintf("%d", rand.Int63())
//			if len(s) < 15 {
//				continue
//			}
//			s = s[len(s)-13:]
//			NewTokenid = s
//			nfttab :=  Nfts{}
//			err := nft.db.Where("contract = ? AND tokenid = ? ", ExchangAddr, NewTokenid).First(&nfttab)
//			if err.Error == gorm.ErrRecordNotFound {
//				break
//			}
//		}
//		if i >= 20 {
//			fmt.Println("UploadNft() generate tokenId error.")
//			return ErrGenerateTokenId
//		}
//		nfttab :=  Nfts{}
//		nfttab.Tokenid = NewTokenid
//		nfttab.Contract = strings.ToLower(ExchangAddr) //nft_contract_addr
//		nfttab.Createaddr = creator_addr
//		nfttab.Ownaddr = owner_addr
//		nfttab.Name = name
//		nfttab.Desc = desc
//		nfttab.Meta = meta
//		nfttab.Categories = categories
//		nfttab.Collectcreator = collectRec.Createaddr
//		nfttab.Collections = collections
//		nfttab.Signdata = sig
//		nfttab.Url = source_url
//		nfttab.Image = asset_sample
//		nfttab.Md5 = md5
//		nfttab.Selltype = SellTypeNotSale.String()
//		nfttab.Verified = NoVerify.String()
//		nfttab.Mintstate = NoMinted.String()
//		nfttab.Createdate = time.Now().Unix()
//		nfttab.Royalty, _ = strconv.Atoi(royalty)
//		nfttab.Royalty /= 100
//		nfttab.Count, _ = strconv.Atoi(count)
//		nfttab.Hide = hide
//		err0, approve := ethhelper.GenCreateNftSign(initExchangAddr, nfttab.Ownaddr, nfttab.Meta,
//			nfttab.Tokenid, count, royalty)
//		if err0 != nil {
//			fmt.Println("UploadNft() GenCreateNftSign() err=", err0)
//			return err0
//		}
//		fmt.Println("UploadNft() GenCreateNftSign() approve=", approve)
//		nfttab.Approve = approve
//		return nft.db.Transaction(func(tx *gorm.DB) error {
//			err := tx.Model(&Nfts{}).Create(&nfttab)
//			if err.Error != nil {
//				fmt.Println("UploadNft() err=", err.Error)
//				return err.Error
//			}
//			if collections != "" {
//				var collectListRec CollectLists
//				collectListRec.Collectsid = collectRec.ID
//				collectListRec.Nftid = nfttab.ID
//				err = tx.Model(&CollectLists{}).Create(&collectListRec)
//				if err.Error != nil {
//					fmt.Println("UploadNft() create CollectLists err=", err.Error)
//					return err.Error
//				}
//				err = tx.Model(&Collects{}).Where("name = ? AND createaddr =?",
//					collections, creator_addr).Update("totalcount",collectRec.Totalcount+1)
//				if err.Error != nil {
//					fmt.Println("UploadNft() add collectins totalcount err= ", err.Error )
//					return err.Error
//				}
//			}
//			return nil
//		})
//	} else {
//		var nfttab Nfts
//		dberr := nft.db.Where("contract = ? AND tokenid = ? ", nft_contract_addr, nft_token_id).First(&nfttab)
//		if dberr.Error == nil {
//			fmt.Println("UploadNft() err=nft already exist.")
//			return ErrNftAlreadyExist
//		}
//		/*ownAddr, royalty, err := func(contract, tokenId string) (string, string, error) {
//			return "ownAddr", "200", nil
//		}(nft_contract_addr, nft_token_id)
//		if ownAddr == user_addr {
//			var nfttab Nfts
//			nfttab.Tokenid = nft_token_id
//			nfttab.Contract = nft_contract_addr //nft_contract_addr
//			nfttab.Createaddr = creator_addr
//			nfttab.Ownaddr = ownAddr
//			nfttab.Name = name
//			nfttab.Desc = desc
//			nfttab.Meta = meta
//			nfttab.Categories = categories
//			nfttab.Collections = collections
//			nfttab.Signdata = sig
//			nfttab.Url = source_url
//			nfttab.Image = asset_sample
//			nfttab.Md5 = md5
//			nfttab.Selltype = SellTypeNotSale.String()
//			nfttab.Verified = NoVerify.String()
//			nfttab.Mintstate = Minted.String()
//			nfttab.Royalty, _ = strconv.Atoi(royalty)
//			nfttab.Royalty = nfttab.Royalty / 100
//			nfttab.Createdate = time.Now().Unix()
//			nfttab.Hide = hide
//			return nft.db.Transaction(func(tx *gorm.DB) error {
//				err := tx.Model(&Nfts{}).Create(&nfttab)
//				if err.Error != nil {
//					fmt.Println("UploadNft() create exist nft err=", err.Error)
//					return err.Error
//				}
//				if collections != "" {
//					var collectListRec CollectLists
//					collectListRec.Collectsid = collectRec.ID
//					collectListRec.Nftid = nfttab.ID
//					err = tx.Model(&CollectLists{}).Create(&collectListRec)
//					if err.Error != nil {
//						fmt.Println("UploadNft() create CollectLists err=", err.Error)
//						return err.Error
//					}
//				}
//				return nil
//			})
//		}*/
//		IsAdminAddr, err := IsAdminAddr(user_addr)
//		if err != nil {
//			fmt.Println("UploadNft() upload address is not admin.")
//			return ErrNftUpAddrNotAdmin
//		}
//		if IsAdminAddr {
//			var nfttab Nfts
//			nfttab.Tokenid = nft_token_id
//			nfttab.Contract = nft_contract_addr //nft_contract_addr
//			nfttab.Createaddr = creator_addr
//			nfttab.Ownaddr = owner_addr
//			nfttab.Name = name
//			nfttab.Desc = desc
//			nfttab.Meta = meta
//			nfttab.Categories = categories
//			nfttab.Collectcreator = creator_addr
//			nfttab.Collections = collections
//			nfttab.Signdata = sig
//			nfttab.Url = source_url
//			nfttab.Image = asset_sample
//			nfttab.Md5 = md5
//			nfttab.Selltype = SellTypeNotSale.String()
//			nfttab.Verified = Passed.String()
//			nfttab.Mintstate = Minted.String()
//			/*nfttab.Royalty, _ = strconv.Atoi(royalty)
//			nfttab.Royalty = nfttab.Royalty / 100*/
//			nfttab.Createdate = time.Now().Unix()
//			nfttab.Royalty, _ = strconv.Atoi(royalty)
//			nfttab.Royalty /= 100
//			nfttab.Count, _ = strconv.Atoi(count)
//			nfttab.Hide = hide
//			return nft.db.Transaction(func(tx *gorm.DB) error {
//				err := tx.Model(&Nfts{}).Create(&nfttab)
//				if err.Error != nil {
//					fmt.Println("UploadNft() admin create nft err=", err.Error)
//					return err.Error
//				}
//				if collections != "" {
//					var collectListRec CollectLists
//					collectListRec.Collectsid = collectRec.ID
//					collectListRec.Nftid = nfttab.ID
//					err = tx.Model(&CollectLists{}).Create(&collectListRec)
//					if err.Error != nil {
//						fmt.Println("UploadNft() create CollectLists err=", err.Error)
//						return err.Error
//					}
//					err = tx.Model(&Collects{}).Where("name = ? AND createaddr =?",
//						collections, creator_addr).Update("totalCount",collectRec.Totalcount+1)
//					if err.Error != nil {
//						fmt.Println("UploadNft() add collectins totalcount err= ", err.Error )
//						return err.Error
//					}
//				}
//				return nil
//			})
//		} else {
//			fmt.Println("UploadNft() upload address is not admin.")
//			return ErrNftUpAddrNotAdmin
//		}
//	}
//	return nil
//}

//function buy_nft(user_addr,sig,nft_contract_addr,nft_token_id)
func (nft NftDb) BuyNft(userAddr, tradeSig, sigdata, contract, nftTokenId string) error {
	userAddr = strings.ToLower(userAddr)
	contract = strings.ToLower(contract)

	var ownaddr string

	trans := Trans{}
	nfts := Nfts{}
	ntfstab := nft.db.Model(&nfts).Where("contract = ? AND tokenid =? ", contract, nftTokenId).First(&nfts)
	if ntfstab.Error != nil {
		return ErrNftNotExist
	}
	ownaddr = nfts.Ownaddr
	//trans.Transid = 0
	trans.Contract = contract
	trans.Fromaddr = ownaddr
	trans.Toaddr = userAddr
	trans.Signdata = sigdata
	trans.Tradesig = tradeSig
	trans.Tokenid = nftTokenId
	trans.Price = nfts.Price
	trans.Transtime = time.Now().Unix()
	return nft.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&trans).Create(&trans)
		if err.Error != nil {
			fmt.Println("buyNft() insert failed, ", err)
			return err.Error
		}
		err = tx.Model(&Nfts{}).Where("contract = ? AND tokenid =? ", contract, nftTokenId).Update("ownaddr", userAddr)
		if err.Error != nil {
			fmt.Println("buyNft() update err=", err.Error)
			return err.Error
		}
		return nil
	})
}

func (nft NftDb) QueryNft() ([]Nfts, error){
	nfts := []Nfts{}
	err := nft.db.Model(&Nfts{}).Find(&nfts)
	if err.Error != nil {
		fmt.Println("queryNft, err=\n ", err)
		return nil, err.Error
	}
	marshal, _ := json.Marshal(nfts)
	fmt.Printf("%s\n",string(marshal))
	//return string(marshal), nil
	//return  marshal, nil
	return nfts, err.Error
}


func (nft *NftDb) joinFilters(filter []StQueryField) string {
	var joinString string
	joinString = ""

	for k1, v1 := range filter {
		if strings.Contains(joinString, v1.Field) {
			// 如果已经处理过该字段, 则继续处理下一个
			continue
		}
		// 如果没有处理过该字段, 则加入查询条件字符串
		if k1 == 0 {
			if !strings.Contains(v1.Field, "price") &&
				!strings.Contains(v1.Field, "date") &&
				!strings.Contains(v1.Field, "time") {
				joinString = joinString + "(" + v1.Field + v1.Operation + "'" + v1.Value + "'"
			} else {
				joinString = joinString + "(" + v1.Field + v1.Operation + v1.Value
			}

		} else {
			if !strings.Contains(v1.Field, "price") &&
				!strings.Contains(v1.Field, "date") &&
				!strings.Contains(v1.Field, "time") {
				joinString = joinString + " and (" + v1.Field + v1.Operation + "'" + v1.Value + "'"
			} else {
				joinString = joinString + " and (" + v1.Field + v1.Operation + v1.Value
			}

		}

		for k2, v2 := range filter {
			// 处理与v1相同的字段
			// k1之前的的数据都已经处理过了, 直接跳过, 只处理k1后的,
			// 且与v1相同的值
			if k2 <= k1 || v2.Field != v1.Field {
				continue
			}
			if !strings.Contains(v2.Field, "price") &&
				!strings.Contains(v2.Field, "date") &&
				!strings.Contains(v2.Field, "time") {
				joinString = joinString + " or " + v2.Field + v2.Operation + "'" + v2.Value + "'"
			} else {
				joinString = joinString + " and " + v2.Field + v2.Operation + v2.Value
			}

		}
		joinString = joinString + ")"
	}

	return joinString
}

func (nft NftDb) QueryNftbyUser(userAddr string)  ([]Nfts, error){
	userAddr = strings.ToLower(userAddr)

	nfts := []Nfts{}
	err := nft.db.Model(&Nfts{}).Where("ownaddr = ?", userAddr).Find(&nfts)
	if err.Error != nil && err.Error != gorm.ErrRecordNotFound{
		fmt.Println("queryNft, err=\n ", err.Error)
		return nil, err.Error
	}
	marshal, _ := json.Marshal(nfts)
	fmt.Printf("%s\n",string(marshal))
	//return string(marshal), nil
	//return marshal, nil
	return nfts, err.Error
}

//type UserInfo struct {
//	Name         	string 	`json:"user_name"`          //用户名
//	//Portrait     	string 	`json:"portrait"`           //用户头像
//	Email		 	string 	`json:"user_mail"`          //邮箱
//	Bio	         	string	`json:"user_info"`          //自我描述
//	Verified	    string	`json:"verified"`           //是否通过验证
//	NftCount	 	int    	`json:"nft_count"`          //用户持有的NFT总数
//	CreateCount	 	int		`json:"create_count"`       //用户创作的NFT总数
//	OwnerCount	 	int    	`json:"owner_count"`        //用户创作的NFT的拥有者数量
//	TradeAmount	 	uint64	`json:"trade_amount"`       //用户创作的NFT的成交额,
//	TradeAvgPrice	uint64 	`json:"trade_avg_price"`    //用户创作的NFT均价,
//	TradeFloorPrice	uint64	`json:"trade_floor_price"`  //用户创作的NFT最低价
//}
//
//func (nft NftDb) QueryUserInfo(userAddr string) (UserInfo, error){
//	userAddr = strings.ToLower(userAddr)
//
//	var uinfo UserInfo
//	user := Users{}
//	err := nft.db.Model(&user).Where("useraddr = ?", userAddr).First(&user)
//	if err.Error != nil {
//		if err.Error == gorm.ErrRecordNotFound {
//			return UserInfo{}, nil
//		}else {
//			fmt.Println("QueryUserInfo() query users err=", err)
//			return UserInfo{}, err.Error
//		}
//	}
//
//	uinfo.Name = user.Username
//	//uinfo.Portrait = user.Portrait
//	uinfo.Email = user.Email
//	uinfo.Bio = user.Bio
//	uinfo.Verified = user.Verified
//	var recCount int64
//	err = nft.db.Model(Nfts{}).Where("ownaddr = ?", userAddr).Count(&recCount)
//	if err.Error == nil {
//		uinfo.NftCount = int(recCount)
//	}
//	err = nft.db.Model(Nfts{}).Where("createaddr = ?", userAddr).Count(&recCount)
//	if err.Error == nil {
//		uinfo.CreateCount = int(recCount)
//	}
//	err = nft.db.Model(Nfts{}).Where("createaddr = ? AND ownaddr != ?",
//			userAddr, userAddr).Count(&recCount)
//	if err.Error == nil {
//		uinfo.OwnerCount = int(recCount)
//	}
//
//	/*type SumInfo struct {
//		SumCount int
//		SumPrice uint64
//	}
//	sum := SumInfo{}
//	err = nft.db.Raw("SELECT SUM(Transcnt) as SumCount, SUM(Transamt) as SumPrice FROM nfts WHERE createaddr = ?", userAddr).Scan(&sum)
//	if err.Error != nil {
//		fmt.Println("QueryUserInfo() query Sum err=", err)
//		return UserInfo{}, err.Error
//	}
//	uinfo.TradeAmount = sum.SumPrice
//	if sum.SumCount != 0 {
//		uinfo.TradeAvgPrice = sum.SumPrice / uint64(sum.SumCount)
//	}
//
//	var nftRec Nfts
//	err = nft.db.Order("transprice desc").Where("createaddr = ?", userAddr).Last(&nftRec)
//	if err.Error != nil {
//		if err.Error != gorm.ErrRecordNotFound {
//			fmt.Println("QueryUserInfo() query statistics err=", err)
//			return UserInfo{}, err.Error
//		}
//	}
//	uinfo.TradeFloorPrice = nftRec.Transprice*/
//
//	type TransInfo struct {
//		TradeAmount	 	uint64
//		TradeAvgPrice	float64
//		TradeFloorPrice	uint64
//		TradeMaxPrice	uint64
//		TradeCount		uint64
//	}
//	tInfo := TransInfo{}
//	sql := "SELECT sum(trans.price) as TradeAmount, avg(trans.price) as TradeAvgPrice, " +
//		"min(trans.price) as TradeFloorPrice, max(trans.price) as TradeMaxPrice, " +
//		"COUNT(trans.price) AS TradeCount " +
//		//"FROM trans" +" WHERE createaddr = ? AND selltype != ? AND selltype != ?"
//		"FROM trans" +" WHERE ( trans.fromaddr = ? OR trans.toaddr = ?) AND selltype != ? AND selltype != ?"
//	err = nft.db.Raw(sql, userAddr, userAddr, SellTypeMintNft.String(), SellTypeError.String()).Scan(&tInfo)
//	if err.Error == nil {
//		uinfo.TradeAmount = tInfo.TradeAmount
//		uinfo.TradeAvgPrice = uint64(tInfo.TradeAvgPrice)
//		uinfo.TradeFloorPrice = tInfo.TradeFloorPrice
//	}
//	return uinfo, err.Error
//}
//
//func (nft NftDb) ModifyUserInfo(user_addr, user_name, portrait, user_mail, user_info, sig string) error{
//	user_addr = strings.ToLower(user_addr)
//
//	fmt.Println("ModifyUserInfo() start.")
//	user := Users{}
//	err := nft.db.Model(&user).Where("useraddr = ?", user_addr).First(&user)
//	if err.Error != nil {
//		fmt.Println("ModifyUserInfo() err= not find user.")
//		return err.Error
//	}
//	user.Username = user_name
//	user.Bio = user_info
//	user.Email = user_mail
//	user.Portrait = portrait
//	user.Signdata = sig
//	err = nft.db.Model(&Users{}).Where("useraddr = ?", user_addr).Updates(user)
//	if err.Error != nil {
//		fmt.Println("ModifyUserInfo() update err= ", err.Error )
//		return err.Error
//	}
//	fmt.Println("ModifyUserInfo() Ok.")
//	return err.Error
//}

func (nft NftDb) Like(userAddr, contractAddr, tokenId, sig string) error {
	userAddr = strings.ToLower(userAddr)
	contractAddr = strings.ToLower(contractAddr)
	var nftrecord Nfts
	err := nft.db.Where("contract = ? AND tokenid =? ", contractAddr, tokenId).First(&nftrecord)
	if err.Error != nil {
		fmt.Println("AddFavor() err= ", err.Error )
		return ErrNftNotExist
	}
	var favorrecord NftFavorited
	err = nft.db.Where("Nftid = ? AND useraddr = ?",
		nftrecord.ID, userAddr).First(&favorrecord)
	if err.Error == nil {
		return nft.db.Transaction(func(tx *gorm.DB) error {
			err = tx.Model(&NftFavorited{}).Where("Nftid = ? AND useraddr = ?",
					nftrecord.ID, userAddr).Delete(&NftFavorited{})
			if err.Error != nil {
				fmt.Println("AddFavor() create record err=", err.Error)
				return err.Error
			}
			if nftrecord.Favorited > 0 {
				favorited := nftrecord.Favorited - 1
				nftrecord = Nfts{}
				nftrecord.Favorited = favorited
				err = tx.Model(&nftrecord).Where("contract = ? AND tokenid =? ", contractAddr, tokenId).Update("favorited",nftrecord.Favorited)
				if err.Error != nil {
					fmt.Println("AddFavor() update NftFavorited err= ", err.Error )
					return err.Error
				}
			}
			return nil
		})
	}
	favorrecord = NftFavorited{}
	favorrecord.Useraddr = userAddr
	favorrecord.Contract = contractAddr
	favorrecord.Tokenid = tokenId
	favorrecord.Nftid = nftrecord.ID
	favorrecord.Signdata = sig
	favorrecord.Name = nftrecord.Name
	favorrecord.Image = nftrecord.Image
	favorrecord.Collectcreator = nftrecord.Collectcreator
	favorrecord.Collections = nftrecord.Collections
	var collectRec Collects
	err = nft.db.Where("createaddr = ? AND name =? ",
		nftrecord.Createaddr, nftrecord.Collections).First(&collectRec)
	if err.Error == nil {
		favorrecord.Img = collectRec.Img
	}
	return nft.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&NftFavorited{}).Create(&favorrecord)
		if err.Error != nil {
			fmt.Println("AddFavor() create record err=", err.Error)
			return err.Error
		}
		favorited := nftrecord.Favorited + 1
		nftrecord = Nfts{}
		nftrecord.Favorited = favorited
		err = tx.Model(&nftrecord).Where("contract = ? AND tokenid =? ", contractAddr, tokenId).Update("favorited",nftrecord.Favorited)
		if err.Error != nil {
			fmt.Println("AddFavor() update NftFavorited err= ", err.Error )
			return err.Error
		}
		return nil
	})
}

func (nft NftDb) DelNftFavor(userAddr, contractAddr, tokenId string) error {
	userAddr = strings.ToLower(userAddr)
	contractAddr = strings.ToLower(contractAddr)

	var nftrecord Nfts
	err := nft.db.Where("contract = ? AND tokenid =? ", contractAddr, tokenId).First(&nftrecord)
	if err.Error != nil {
		fmt.Println("DelFavor() err= ", err.Error )
		return err.Error
	}
	return nft.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&NftFavorited{}).Where("nftid = ? AND useraddr = ?", nftrecord.ID, userAddr).Delete(&NftFavorited{})
		if err.Error != nil {
			if err.Error == gorm.ErrRecordNotFound {
				return ErrNotNftFavorited
			}
			fmt.Println("DelFavor() err=", err.Error )
			return err.Error
		}
		err = tx.Model(&nftrecord).Where("contract = ? AND tokenid =? ", contractAddr, tokenId).Update("Favorited",nftrecord.Favorited-1)
		if err.Error != nil {
			fmt.Println("AddFavor() update NftFavorited err= ", err.Error )
			return err.Error
		}
		return nil
	})
}

func (nft NftDb) QueryNftFavorited(userAddr string)  ([]Nfts, error){
	userAddr = strings.ToLower(userAddr)

	favors := []NftFavorited{}
	err := nft.db.Where("useraddr = ?", userAddr).Find(&favors)
	if err.Error != nil {
		fmt.Println("queryNft, err=\n ", err.Error)
		return nil, err.Error
	}
	nfts := []Nfts{}
	for _, favor := range favors {
		var nftrecord Nfts
		err = nft.db.Where("ID = ?", favor.Nftid).First(&nftrecord)
		if err.Error != nil {
			fmt.Println("AddFavor() err= ", err.Error )
			break
		}
		nftrecord.Image = ""
		nfts = append(nfts, nftrecord)
	}
	marshal, _ := json.Marshal(nfts)
	fmt.Printf("%s\n",string(marshal))
	//return string(marshal), nil
	//return marshal, nil
	return nfts, err.Error
}

//获取NFT待审核列表
func (nft NftDb) QueryUnverifiedNfts(verified string) ([]Nfts, error){
	nfts := []Nfts{}
	//nft.db.Table("nfts").Model()
	queryResult := nft.db.Where("verified =?", verified).Find(&nfts)
	if queryResult.Error != nil {
		fmt.Println("QueryUnverifiedNfts, err =", queryResult.Error)
		return nil, queryResult.Error
	}
	for k, _ := range nfts {
		nfts[k].Image = ""
	}
	fmt.Println("rowsnum=", queryResult.RowsAffected)

	return nfts, nil
}

//审核NFT*
func (nft NftDb) VerifyNft(vrfaddr string, owner string, contractaddr string,
	tokenid string, desc string, verified string, sig string) error {

	vrfaddr = strings.ToLower(vrfaddr)
	owner = strings.ToLower(owner)
	contractaddr = strings.ToLower(contractaddr)

	//modify the database value of verified field if the valification address is valid.
	nftData := Nfts{}
	takeResult := nft.db.Where("contract = ? and tokenid = ?",
		contractaddr, tokenid).Take(&nftData)
	if takeResult.Error != nil {
		return takeResult.Error
	}
	updateValue := make(map[string]interface{})
	updateValue["verified"] = verified
	updateValue["verifieddesc"] = desc
	updateValue["signdata"] = sig
	updateValue["verifiedtime"] = time.Now().Unix()
	updateResult := nft.db.Model(&nftData).Updates(updateValue)
	if updateResult.Error != nil {
		return updateResult.Error
	}

	return nil
}

//func (nft NftDb) Sell(ownAddr,
//	PrivAddr string,
//	contractAddr,
//	tokenId string,
//	sellType string,
//	payChan string,
//	days int,
//	startPrice,
//	endPrice uint64,
//	royalty string,
//	currency string,
//	hide string,
//	sigData string,
//	tradeSig string) error {
//
//	ownAddr = strings.ToLower(ownAddr)
//	PrivAddr = strings.ToLower(PrivAddr)
//	contractAddr = strings.ToLower(contractAddr)
//
//	fmt.Println(time.Now().String()[:22], "Sell() Start.",
//		"tokenId=", tokenId,
//		"SellType=", sellType,
//		"startPrice=", startPrice,
//		"endPrice=", endPrice)
//	defer fmt.Println(time.Now().String()[:22], "Sell() end.")
//	var nftrecord Nfts
//	err := nft.db.Where("contract = ? AND tokenid =? AND ownaddr = ?", contractAddr, tokenId, ownAddr).First(&nftrecord)
//	if err.Error != nil {
//		fmt.Println("Sell() err= ", err.Error )
//		return err.Error
//	}
//	if nftrecord.Verified != Passed.String() {
//		return ErrNotVerify
//	}
//	/*if nftrecord.Mintstate != Minted.String() {
//		return ErrNftNotMinted
//	}*/
//	//if startDate.After(endDate) {
//	//	return ErrAuctionStartAfterEnd
//	//}
//	//if startDate.Before(time.Now()) {
//	//	startDate = time.Now()
//	//	//return ErrAuctionStartBeforeNow
//	//}
//	var auctionRec Auction
//	err = nft.db.Where("contract = ? AND nftid = ? AND ownaddr = ?",
//		nftrecord.Contract, nftrecord.ID, ownAddr).First(&auctionRec)
//	if err.Error == nil {
//		if auctionRec.Selltype != SellTypeBidPrice.String() {
//			fmt.Println("Sell() err=", err.Error, ErrNftSelling)
//			return ErrNftSelling
//		} else {
//			err := nft.db.Transaction(func(tx *gorm.DB) error {
//				err = tx.Model(&Bidding{}).Where("contract = ? AND tokenid = ?",
//					auctionRec.Contract, auctionRec.Tokenid).Delete(&Bidding{})
//				if err.Error != nil {
//					fmt.Println("Sell() delete bid record err=", err.Error)
//					return err.Error
//				}
//				err = tx.Model(&Auction{}).Where("contract = ? AND tokenid = ?",
//					auctionRec.Contract, auctionRec.Tokenid).Delete(&Auction{})
//				if err.Error != nil {
//					fmt.Println("Sell() delete bidprice auction record err=", err.Error)
//					return err.Error
//				}
//				return nil
//			})
//			if err != nil {
//				fmt.Println("Sell() delete bidprice err=", err)
//				return err
//			}
//		}
//	}
//	auctionRec = Auction{}
//	auctionRec.Selltype = sellType
//	auctionRec.Paychan = payChan
//	auctionRec.Ownaddr = ownAddr
//	auctionRec.Nftid = nftrecord.ID
//	auctionRec.Contract = contractAddr
//	auctionRec.Tokenid = tokenId
//	auctionRec.Currency = currency
//	auctionRec.Startprice = startPrice
//	auctionRec.Endprice = endPrice
//	auctionRec.Privaddr = PrivAddr
//	auctionRec.Startdate = time.Now().Unix()
//	//auctionRec.Enddate = time.Now().AddDate(0, 0, days).Unix()
//	auctionRec.Enddate = time.Now().Add(1 * time.Minute).Unix()
//	auctionRec.Signdata = sigData
//	auctionRec.Tradesig = tradeSig
//	auctionRec.SaleState = "SaleStart"
//
//	if sellType == SellTypeFixPrice.String() {
//		auctionRec.Startprice = startPrice
//		auctionRec.Endprice = startPrice
//	}
//	auctionHistory := AuctionHistory{}
//	auctionHistory.AuctionRecord = auctionRec.AuctionRecord
//	return nft.db.Transaction(func(tx *gorm.DB) error {
//		err = tx.Model(&auctionRec).Create(&auctionRec)
//		if err.Error != nil {
//			fmt.Println("Sell() create auctionRec record err=", err.Error)
//			return err.Error
//		}
//		err = tx.Model(&AuctionHistory{}).Create(&auctionHistory)
//		if err.Error != nil {
//			fmt.Println("Sell() create auctionHistory record err=", err.Error)
//			return err.Error
//		}
//		nftrecord = Nfts{}
//		nftrecord.Hide = hide
//		nftrecord.Selltype = sellType
//		err = tx.Model(&Nfts{}).Where("contract = ? AND tokenid =?",
//			auctionRec.Contract, auctionRec.Tokenid).Updates(&nftrecord)
//		if err.Error != nil {
//			fmt.Println("Sell() update record err=", err.Error)
//			return err.Error
//		}
//		/*nftrecord = Nfts{}
//		nftrecord.Royalty, _ = strconv.Atoi(royalty)
//		nftrecord.Royalty = nftrecord.Royalty / 100
//		err = tx.Model(&Nfts{}).Where("contract = ? AND tokenid =? AND royalty = ?",
//			auctionRec.Contract, auctionRec.Tokenid, 0).Updates(&nftrecord)
//		if err.Error != nil {
//			fmt.Println("Sell() update record err=", err.Error)
//			return err.Error
//		}*/
//		return nil
//	})
//}

//func (nft NftDb) ModifyPrice(ownAddr,
//							 contractAddr,
//							 tokenId string,
//							 sellType string,
//							 price,
//							 fixedPrice uint64,
//							 currency string,
//							 hide string,
//							 tradeSig, sigData string,
//							 time1, time2 time.Time) error {
//	var auctionRec Auction
//	err := nft.db.Where("contract = ? AND tokenid = ? AND ownAddr = ?", contractAddr, tokenId, ownAddr).First(&auctionRec)
//	if err.Error != nil {
//		return ErrNftNotSell
//	}
//	auctionRec = Auction{}
//	auctionRec.Selltype = sellType
//	auctionRec.Contract = contractAddr
//	auctionRec.Tokenid = tokenId
//	auctionRec.Currency = currency
//	auctionRec.Startprice = price
//	auctionRec.Endprice = fixedPrice
//	auctionRec.Startdate = time1
//	auctionRec.Enddate = time2
//	auctionRec.Tradesig = tradeSig
//	auctionRec.Signdata = sigData
//	return nft.db.Transaction(func(tx *gorm.DB) error {
//		err = tx.Model(&auctionRec).Updates(&auctionRec)
//		if err.Error != nil {
//			fmt.Println("ModifyPrice() updates record err=", err.Error)
//			return err.Error
//		}
//		nftrecord := Nfts{}
//		nftrecord.Hide = hide
//		nftrecord.Selltype = sellType
//		err = tx.Model(&Nfts{}).Where("contract = ? AND tokenid =?",
//			auctionRec.Contract, auctionRec.Tokenid).Updates(&nftrecord)
//		if err.Error != nil {
//			fmt.Println("ModifyPrice() update record err=", err.Error)
//			return err.Error
//		}
//		return nil
//	})
//}

func (nft NftDb) CancelBuy(UserAddr, NftContractAddr, NftTokenId, TradeSig, Sig string) error {
	UserAddr = strings.ToLower(UserAddr)
	NftContractAddr = strings.ToLower(NftContractAddr)

	err := nft.db.Model(&Bidding{}).Where("Bidaddr = ? AND Contract = ? AND Tokenid =?",
					UserAddr, NftContractAddr, NftTokenId).Delete(&Bidding{})
	if err.Error != nil {
		fmt.Println("cancelBuy() update record err=", err.Error)
	}
	return err.Error
}

func (nft NftDb) CancellSell(ownAddr, contractAddr, tokenId, sigData string) error {
	ownAddr = strings.ToLower(ownAddr)
	contractAddr = strings.ToLower(contractAddr)

	var nftrecord Nfts
	err := nft.db.Where("contract = ? AND tokenid = ? AND ownaddr = ?", contractAddr, tokenId, ownAddr).First(&nftrecord)
	if err.Error != nil {
		fmt.Println("CancellSell() err= ", err.Error )
		return err.Error
	}
	var auctionRec Auction
	err = nft.db.Where("nftid = ? AND ownaddr = ?", nftrecord.ID, ownAddr).First(&auctionRec)
	if err.Error != nil {
		return ErrNftNotSell
	}
	return nft.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Where("nftid = ? AND ownaddr = ?", nftrecord.ID, ownAddr).Delete(&auctionRec)
		if err.Error != nil {
			return err.Error
		}
		err = tx.Model(&Bidding{}).Where("auctionid = ?", auctionRec.ID).Delete(&Bidding{})
		if err.Error != nil {
			fmt.Println("CancellSell() delete bid record err=", err.Error)
			return err.Error
		}
		err = tx.Model(&Nfts{}).Where("contract = ? AND tokenid =?",
			auctionRec.Contract, auctionRec.Tokenid).Update("selltype", SellTypeNotSale.String())
		if err.Error != nil {
			fmt.Println("CancellSell() update record err=", err.Error)
			return err.Error
		}
		return nil
	})
}

//func (nft NftDb) MakeOffer(userAddr,
//	                       contractAddr,
//	                       tokenId string,
//	                       PayChannel string,
//	                       CurrencyType string,
//	                       price uint64,
//	                       TradeSig string,
//	                       dead_time int64,
//	                       sigdata string) error {
//	userAddr = strings.ToLower(userAddr)
//	contractAddr = strings.ToLower(contractAddr)
//	var auctionRec Auction
//	err := nft.db.Where("contract = ? AND tokenid = ?", contractAddr, tokenId).First(&auctionRec)
//	if err.Error != nil {
//		if err.Error == gorm.ErrRecordNotFound {
//			fmt.Println("MakeOffer() RecordNotFound")
//			var nftrecord Nfts
//			err := nft.db.Where("contract = ? AND tokenid =?", contractAddr, tokenId).First(&nftrecord)
//			if err.Error != nil {
//				fmt.Println("MakeOffer() bidprice not find nft err= ", err.Error )
//				return ErrNftNotExist
//			}
//			auctionRec = Auction{}
//			auctionRec.Selltype = SellTypeBidPrice.String()
//			auctionRec.Paychan = PayChannel
//			auctionRec.Ownaddr = nftrecord.Ownaddr
//			auctionRec.Nftid = nftrecord.ID
//			auctionRec.Contract = contractAddr
//			auctionRec.Tokenid = tokenId
//			auctionRec.Currency = CurrencyType
//			auctionRec.Startprice = price
//			auctionRec.Endprice = price
//			auctionRec.Startdate = time.Now().Unix()
//			auctionRec.Enddate = time.Now().Unix()
//			auctionRec.Signdata = sigdata
//			auctionRec.Tradesig = TradeSig
//			auctionHistory := AuctionHistory{}
//			auctionHistory.AuctionRecord = auctionRec.AuctionRecord
//			return nft.db.Transaction(func(tx *gorm.DB) error {
//				err = tx.Model(&auctionRec).Create(&auctionRec)
//				if err.Error != nil {
//					fmt.Println("MakeOffer() create auctionRec record err=", err.Error)
//					return err.Error
//				}
//				err = tx.Model(&AuctionHistory{}).Create(&auctionHistory)
//				if err.Error != nil {
//					fmt.Println("MakeOffer() create auctionHistory record err=", err.Error)
//					return err.Error
//				}
//				nftrecord = Nfts{}
//				nftrecord.Selltype = auctionRec.Selltype
//				err = tx.Model(&Nfts{}).Where("contract = ? AND tokenid =?",
//					auctionRec.Contract, auctionRec.Tokenid).Updates(&nftrecord)
//				if err.Error != nil {
//					fmt.Println("MakeOffer() update record err=", err.Error)
//					return err.Error
//				}
//				bidRec := Bidding{}
//				bidRec.Bidaddr = userAddr
//				bidRec.Auctionid = auctionRec.ID
//				bidRec.Contract = contractAddr
//				bidRec.Tokenid = tokenId
//				bidRec.Price = price
//				bidRec.Currency = CurrencyType
//				bidRec.Paychan = PayChannel
//				bidRec.Tradesig = TradeSig
//				bidRec.Bidtime = time.Now().Unix()
//				bidRec.Signdata = sigdata
//				bidRec.Deadtime = dead_time
//				bidRec.Nftid = auctionRec.Nftid
//				bidRecHistory := BiddingHistory{}
//				bidRecHistory.BidRecord = bidRec.BidRecord
//				err := tx.Model(&bidRec).Create(&bidRec)
//				if err.Error != nil {
//					fmt.Println("MakeOffer() create bidRec record err=", err.Error)
//					return err.Error
//				}
//				err = tx.Model(&BiddingHistory{}).Create(&bidRecHistory)
//				if err.Error != nil {
//					fmt.Println("MakeOffer() create bidRecHistory record err=", err.Error)
//					return err.Error
//				}
//				fmt.Println("MakeOffer() RecordNotFound OK")
//				return nil
//			})
//		}
//		return ErrNftNotSell
//	}
//	//if time.Now().Unix() < auctionRec.Startdate {
//	//	return ErrAuctionNotBegan
//	//}
//	if auctionRec.Selltype == SellTypeHighestBid.String() {
//		//addrs, err := ethhelper.BalanceOfWeth()
//		fmt.Println("MakeOffer() Selltype == SellTypeHighestBid")
//		if time.Now().Unix() >= auctionRec.Enddate {
//			fmt.Println("MakeOffer() time.Now().Unix() >= auctionRec.Enddate")
//			return ErrAuctionEnd
//		}
//		if auctionRec.Startprice > price {
//			fmt.Println("MakeOffer() auctionRec.Startprice > price")
//			return ErrBidOutRange
//		}
//		var bidRec Bidding
//		err = nft.db.Where("contract = ? AND tokenid = ? AND bidAddr = ?", contractAddr, tokenId, userAddr).First(&bidRec)
//		if err.Error == nil {
//			fmt.Println("MakeOffer() first bidding.")
//			bidRec = Bidding{}
//			bidRec.Price = price
//			bidRec.Currency = CurrencyType
//			bidRec.Paychan = PayChannel
//			bidRec.Tradesig = TradeSig
//			bidRec.Bidtime = time.Now().Unix()
//			bidRec.Signdata = sigdata
//			return nft.db.Transaction(func(tx *gorm.DB) error {
//				err := tx.Model(&bidRec).Where("contract = ? AND tokenid = ? AND bidAddr = ?", contractAddr, tokenId, userAddr).Updates(&bidRec)
//				if err.Error != nil {
//					fmt.Println("MakeOffer() update Bidding record err=", err.Error)
//					return err.Error
//				}
//				bidRecHistory := BiddingHistory(bidRec)
//				err = tx.Model(&BiddingHistory{}).Where("contract = ? AND tokenid = ? AND bidAddr = ?", contractAddr, tokenId, userAddr).Updates(&bidRecHistory)
//				if err.Error != nil {
//					fmt.Println("MakeOffer() update bidRecHistory record err=", err.Error)
//					return err.Error
//				}
//				fmt.Println("MakeOffer() first bidding OK.")
//				return nil
//			})
//		} else{
//			bidRec = Bidding{}
//			bidRec.Bidaddr = userAddr
//			bidRec.Auctionid = auctionRec.ID
//			bidRec.Nftid = auctionRec.Nftid
//			bidRec.Contract = contractAddr
//			bidRec.Tokenid = tokenId
//			bidRec.Price = price
//			bidRec.Currency = CurrencyType
//			bidRec.Paychan = PayChannel
//			bidRec.Tradesig = TradeSig
//			bidRec.Bidtime = time.Now().Unix()
//			bidRec.Signdata = sigdata
//			return nft.db.Transaction(func(tx *gorm.DB) error {
//				err := tx.Model(&bidRec).Create(&bidRec)
//				if err.Error != nil {
//					fmt.Println("MakeOffer() create record err=", err.Error)
//					return err.Error
//				}
//				bidRecHistory := BiddingHistory{}
//				bidRecHistory.BidRecord = bidRec.BidRecord
//				err = tx.Model(&BiddingHistory{}).Create(&bidRecHistory)
//				if err.Error != nil {
//					fmt.Println("MakeOffer() create bidRecHistory record err=", err.Error)
//					return err.Error
//				}
//				fmt.Println("MakeOffer() change bidding OK.")
//				return nil
//			})
//		}
//	}
//	if auctionRec.Selltype == SellTypeBidPrice.String() {
//		fmt.Println("MakeOffer() Selltype == SellTypeBidPrice")
//		var bidRec Bidding
//		err = nft.db.Where("contract = ? AND tokenid = ? AND bidAddr = ?", contractAddr, tokenId, userAddr).First(&bidRec)
//		if err.Error == nil {
//			bidRec = Bidding{}
//			bidRec.Price = price
//			bidRec.Currency = CurrencyType
//			bidRec.Paychan = PayChannel
//			bidRec.Tradesig = TradeSig
//			bidRec.Bidtime = time.Now().Unix()
//			bidRec.Signdata = sigdata
//			return nft.db.Transaction(func(tx *gorm.DB) error {
//				err := tx.Model(&bidRec).Where("contract = ? AND tokenid = ? AND bidAddr = ?", contractAddr, tokenId, userAddr).Updates(&bidRec)
//				if err.Error != nil {
//					fmt.Println("MakeOffer() update Bidding record err=", err.Error)
//					return err.Error
//				}
//				bidRecHistory := BiddingHistory(bidRec)
//				err = tx.Model(&BiddingHistory{}).Where("contract = ? AND tokenid = ? AND bidAddr = ?", contractAddr, tokenId, userAddr).Updates(&bidRecHistory)
//				if err.Error != nil {
//					fmt.Println("MakeOffer() update bidRecHistory record err=", err.Error)
//					return err.Error
//				}
//				fmt.Println("MakeOffer() change bidding OK.")
//				return nil
//			})
//		} else {
//			return nft.db.Transaction(func(tx *gorm.DB) error {
//				bidRec := Bidding{}
//				bidRec.Bidaddr = userAddr
//				bidRec.Auctionid = auctionRec.ID
//				bidRec.Nftid = auctionRec.Nftid
//				bidRec.Contract = contractAddr
//				bidRec.Tokenid = tokenId
//				bidRec.Price = price
//				bidRec.Currency = CurrencyType
//				bidRec.Paychan = PayChannel
//				bidRec.Tradesig = TradeSig
//				bidRec.Bidtime = time.Now().Unix()
//				bidRec.Deadtime = dead_time
//				bidRec.Signdata = sigdata
//				bidRecHistory := BiddingHistory{}
//				bidRecHistory.BidRecord = bidRec.BidRecord
//				err := tx.Model(&bidRec).Create(&bidRec)
//				if err.Error != nil {
//					fmt.Println("MakeOffer() create bidRec record err=", err.Error)
//					return err.Error
//				}
//				err = tx.Model(&BiddingHistory{}).Create(&bidRecHistory)
//				if err.Error != nil {
//					fmt.Println("MakeOffer() create bidRecHistory record err=", err.Error)
//					return err.Error
//				}
//				fmt.Println("MakeOffer() first bidding OK.")
//				return nil
//			})
//		}
//	}
//	if auctionRec.Selltype == SellTypeFixPrice.String() {
//		fmt.Println("MakeOffer() Selltype == SellTypeFixPrice")
//		var bidRec Bidding
//		err = nft.db.Where("contract = ? AND tokenid = ? AND bidAddr = ?", contractAddr, tokenId, userAddr).First(&bidRec)
//		if err.Error == nil {
//			bidRec = Bidding{}
//			bidRec.Price = price
//			bidRec.Currency = CurrencyType
//			bidRec.Paychan = PayChannel
//			bidRec.Tradesig = TradeSig
//			bidRec.Bidtime = time.Now().Unix()
//			bidRec.Signdata = sigdata
//			return nft.db.Transaction(func(tx *gorm.DB) error {
//				err := tx.Model(&bidRec).Where("contract = ? AND tokenid = ? AND bidAddr = ?", contractAddr, tokenId, userAddr).Updates(&bidRec)
//				if err.Error != nil {
//					fmt.Println("MakeOffer() update Bidding record err=", err.Error)
//					return err.Error
//				}
//				bidRecHistory := BiddingHistory(bidRec)
//				err = tx.Model(&BiddingHistory{}).Where("contract = ? AND tokenid = ? AND bidAddr = ?", contractAddr, tokenId, userAddr).Updates(&bidRecHistory)
//				if err.Error != nil {
//					fmt.Println("MakeOffer() update bidRecHistory record err=", err.Error)
//					return err.Error
//				}
//				fmt.Println("MakeOffer() change bidding OK.")
//				return nil
//			})
//		} else {
//			return nft.db.Transaction(func(tx *gorm.DB) error {
//				bidRec := Bidding{}
//				bidRec.Bidaddr = userAddr
//				bidRec.Auctionid = auctionRec.ID
//				bidRec.Nftid = auctionRec.Nftid
//				bidRec.Contract = contractAddr
//				bidRec.Tokenid = tokenId
//				bidRec.Price = price
//				bidRec.Currency = CurrencyType
//				bidRec.Paychan = PayChannel
//				bidRec.Tradesig = TradeSig
//				bidRec.Bidtime = time.Now().Unix()
//				bidRec.Deadtime = dead_time
//				bidRec.Signdata = sigdata
//				bidRecHistory := BiddingHistory{}
//				bidRecHistory.BidRecord = bidRec.BidRecord
//				err := tx.Model(&bidRec).Create(&bidRec)
//				if err.Error != nil {
//					fmt.Println("MakeOffer() create bidRec record err=", err.Error)
//					return err.Error
//				}
//				err = tx.Model(&BiddingHistory{}).Create(&bidRecHistory)
//				if err.Error != nil {
//					fmt.Println("MakeOffer() create bidRecHistory record err=", err.Error)
//					return err.Error
//				}
//				fmt.Println("MakeOffer() first bidding OK.")
//				return nil
//			})
//		}
//	}
//	return ErrNftNotSell
//}
//
//func (nft NftDb) BuyResult(from, to, contractAddr, tokenId, trade_sig, price, sig, royalty, txhash string) error {
//	from = strings.ToLower(from)
//	to = strings.ToLower(to)
//	contractAddr = strings.ToLower(contractAddr)
//
//	if IsUint64DataValid(price) != true {
//		return ErrPrice
//	}
//	fmt.Println(time.Now().String()[:25],"BuyResult() Begin", "from=", from, "to=", to, "price=", price,
//		"contractAddr=", contractAddr, "tokenId=", tokenId,
//		"royalty=", royalty/*, "sig=", sig, "trade_sig=", trade_sig*/)
//	fmt.Println("BuyResult()++q++++++++++++++++++")
//	if royalty != "" {
//		fmt.Println("BuyResult() royalty!=\"\" mint royalty=", royalty)
//		var nftRec Nfts
//		err := nft.db.Where("contract = ? AND tokenid = ?", contractAddr, tokenId).First(&nftRec)
//		if err.Error != nil {
//			fmt.Println("BuyResult() royalty err =", ErrNftNotExist)
//			return ErrNftNotExist
//		}
//		trans := Trans{}
//		trans.Contract = contractAddr
//		trans.Fromaddr = ""
//		trans.Toaddr = to
//		trans.Signdata = sig
//		trans.Tokenid = tokenId
//		trans.Price, _ = strconv.ParseUint(price, 10, 64)
//		trans.Transtime = time.Now().Unix()
//		trans.Selltype = SellTypeMintNft.String()
//		trans.Name = nftRec.Name
//		trans.Meta = nftRec.Meta
//		trans.Desc = nftRec.Desc
//		trans.Txhash = txhash
//		return nft.db.Transaction(func(tx *gorm.DB) error {
//			err := tx.Model(&trans).Create(&trans)
//			if err.Error != nil {
//				fmt.Println("BuyResult() royalty create trans err=", err.Error)
//				return err.Error
//			}
//			nftrecord := Nfts{}
//			//nftrecord.Ownaddr = to
//			nftrecord.Signdata = sig
//
//			nftrecord.Royalty, _ = strconv.Atoi(royalty)
//			nftrecord.Royalty = nftrecord.Royalty / 100
//			//nftrecord.Selltype = SellTypeNotSale.String()
//			nftrecord.Mintstate = Minted.String()
//			err = tx.Model(&Nfts{}).Where("contract = ? AND tokenid =?",
//				contractAddr, tokenId).Updates(&nftrecord)
//			if err.Error != nil {
//				fmt.Println("BuyResult() royalty update nfts record err=", err.Error)
//				return err.Error
//			}
//			fmt.Println("BuyResult() royalty!=\"\" Ok")
//			return nil
//		})
//	}
//	fmt.Println("BuyResult()-------------------")
//	if from != "" && to != "" {
//		fmt.Println("BuyResult() 1 from != \"\" && to != \"\"" )
//		var nftRec Nfts
//		err := nft.db.Where("contract = ? AND tokenid = ?", contractAddr, tokenId).First(&nftRec)
//		if err.Error != nil {
//			fmt.Println("BuyResult() auction not find err=", err.Error)
//			return ErrNftNotExist
//		}
//		if price == "" {
//			fmt.Println("BuyResult() price == null" )
//			return nft.db.Transaction(func(tx *gorm.DB) error {
//				var auctionRec Auction
//				err = tx.Set("gorm:query_option", "FOR UPDATE").Where("contract = ? AND tokenid = ? AND ownaddr =?",
//					contractAddr, tokenId, nftRec.Ownaddr).First(&auctionRec)
//				if err.Error != nil {
//					fmt.Println("BuyResult() auction not find err=", err.Error)
//					return err.Error
//				}
//				fmt.Println("BuyResult() 1  price = null SaleState=", SaleWait.String())
//				trans := Trans{}
//				trans.Auctionid = auctionRec.ID
//				trans.Contract = auctionRec.Contract
//				trans.Createaddr = nftRec.Createaddr
//				trans.Fromaddr = from
//				trans.Toaddr = to
//				trans.Signdata = sig
//				trans.Tokenid = auctionRec.Tokenid
//				trans.Nftid = auctionRec.Nftid
//				trans.Paychan = auctionRec.Paychan
//				trans.Currency = auctionRec.Currency
//				trans.Price = 0
//				trans.Transtime = time.Now().Unix()
//				trans.Selltype = SellTypeAsset.String()
//				err := tx.Model(&trans).Create(&trans)
//				if err.Error != nil {
//					fmt.Println("BuyResult() create trans record err=", err.Error)
//					return err.Error
//				}
//				fmt.Println("BuyResult() 2  price == null OK" )
//				return nil
//			})
//		}else{
//			fmt.Println("BuyResult() price != null" )
//			return nft.db.Transaction(func(tx *gorm.DB) error {
//				var auctionRec Auction
//				err = tx.Where("contract = ? AND tokenid = ? AND ownaddr =?",
//					contractAddr, tokenId, nftRec.Ownaddr).First(&auctionRec)
//				if err.Error != nil {
//					fmt.Println("BuyResult() auction not find err=", err.Error)
//					return err.Error
//				}
//				fmt.Println("BuyResult() 1 price !=0 SaleState=", SaleWait.String())
//				trans := Trans{}
//				trans.Auctionid = auctionRec.ID
//				trans.Contract = auctionRec.Contract
//				trans.Createaddr = nftRec.Createaddr
//				trans.Fromaddr = from
//				trans.Toaddr = to
//				trans.Signdata = sig
//				trans.Nftid = auctionRec.Nftid
//				trans.Tokenid = auctionRec.Tokenid
//				trans.Paychan = auctionRec.Paychan
//				trans.Currency = auctionRec.Currency
//				trans.Txhash = txhash
//				trans.Name = nftRec.Name
//				trans.Meta = nftRec.Meta
//				trans.Desc = nftRec.Desc
//				trans.Price, _ = strconv.ParseUint(price, 10, 64)
//				trans.Transtime = time.Now().Unix()
//				if auctionRec.Selltype == SellTypeWaitSale.String() {
//					trans.Selltype = SellTypeHighestBid.String()
//				}else {
//					trans.Selltype = auctionRec.Selltype
//				}
//				err := tx.Model(&trans).Create(&trans)
//				if err.Error != nil {
//					fmt.Println("BuyResult() create trans record err=", err.Error)
//					return err.Error
//				}
//				nftrecord := Nfts{}
//				nftrecord.Ownaddr = to
//				nftrecord.Selltype = SellTypeNotSale.String()
//				nftrecord.Paychan = auctionRec.Paychan
//				nftrecord.TransCur = auctionRec.Currency
//				nftrecord.Transprice = trans.Price
//				nftrecord.Transamt += trans.Price
//				nftrecord.Transcnt += 1
//				nftrecord.Transtime = time.Now().Unix()
//				err = tx.Model(&Nfts{}).Where("contract = ? AND tokenid =?",
//					auctionRec.Contract, auctionRec.Tokenid).Updates(&nftrecord)
//				if err.Error != nil {
//					fmt.Println("BuyResult() update record err=", err.Error)
//					return err.Error
//				}
//				err = tx.Model(&Auction{}).Where("contract = ? AND tokenid = ?",
//					auctionRec.Contract, auctionRec.Tokenid).Delete(&Auction{})
//				if err.Error != nil {
//					fmt.Println("BuyResult() delete auction record err=", err.Error)
//					return err.Error
//				}
//				err = nft.db.Model(&Bidding{}).Where("contract = ? AND tokenid = ?",
//					auctionRec.Contract, auctionRec.Tokenid).Delete(&Bidding{})
//				if err.Error != nil {
//					fmt.Println("BuyResult() delete bid record err=", err.Error)
//					return err.Error
//				}
//				fmt.Println("BuyResult() from != \"\" && to != \"\" --> price != \"\" OK" )
//				return nil
//			})
//		}
//	}
//	fmt.Println("BuyResult() End.")
//	return ErrFromToAddrZero
//}

func (nft NftDb) QuerySellNfts() ([]Auction, error) {
	var auctionRecs []Auction
	err := nft.db.Find(&auctionRecs)
	if err.Error != nil {
		return nil, ErrNftNotSell
	}
	return auctionRecs, err.Error
}

func (nft NftDb) QuerySingleSellNft(contract, tokenId string) (*Auction, error)  {
	contract = strings.ToLower(contract)

	var auctionRec Auction
	err := nft.db.Where("contract = ? AND tokenid = ?", contract, tokenId).First(&auctionRec)
	if err.Error != nil {
		return nil, ErrNftNotSell
	}
	return &auctionRec, err.Error
}

func (nft NftDb) QuerySigInfo(signData  string) (Sigmsgrec, error)  {
	var sig Sigmsgs
	err := nft.db.Where("signdata = ?", signData).First(&sig)
	if err.Error != nil {
		return Sigmsgrec{}, err.Error
	}
	return sig.Sigmsgrec, err.Error
}

//type NftAuction struct {
//	Selltype        string `json:"selltype"`
//	Ownaddr         string `json:"ownaddr"`
//	NftTokenId      string `json:"nft_token_id"`
//	NftContractAddr string `json:"nft_contract_addr"`
//	Paychan         string `json:"paychan"`
//	Currency        string `json:"currency"`
//	Startprice      uint64 `json:"startprice"`
//	Endprice        uint64 `json:"endprice"`
//	Startdate       int64  `json:"startdate"`
//	Enddate         int64  `json:"enddate"`
//	Tradesig       	string `json:"tradesig"`
//}
//
//type NftTran struct {
//	NftContractAddr string `json:"nft_contract_addr"`
//	Fromaddr        string `json:"fromaddr"`
//	Toaddr          string `json:"toaddr"`
//	NftTokenId      string `json:"nft_token_id"`
//	Transtime       int64  `json:"transtime"`
//	Paychan         string `json:"paychan"`
//	Currency        string `json:"currency"`
//	Price           uint64 `json:"price"`
//	Txhash			string `json:"trade_hash"`
//	Selltype        string `json:"selltype"`
//}
//
//type NftBid struct {
//	Bidaddr         string `json:"bidaddr"`
//	NftTokenId      string `json:"nft_token_id"`
//	NftContractAddr string `json:"nft_contract_addr"`
//	Paychan         string `json:"paychan"`
//	Currency        string `json:"currency"`
//	Price           uint64 `json:"price"`
//	Bidtime         int64  `json:"bidtime"`
//	Tradesig       	string `json:"tradesig"`
//}
//
//type NftSingleInfo struct {
//	Name            string 			`json:"name"`
//	CreatorAddr     string 			`json:"creator_addr"`
//	//CreatorPortrait string 			`json:"creator_portrait"`
//	OwnerAddr       string 			`json:"owner_addr"`
//	//OwnerPortrait   string 			`json:"owner_portrait"`
//	Md5             string 			`json:"md5"`
//	//AssetSample     string 			`json:"asset_sample"`
//	Desc            string 			`json:"desc"`
//	Collectiondesc  string 			`json:"collection_desc"`
//	Meta            string 			`json:"meta"`
//	SourceUrl       string 			`json:"source_url"`
//	NftContractAddr string 			`json:"nft_contract_addr"`
//	NftTokenId      string 			`json:"nft_token_id"`
//	Categories      string 			`json:"categories"`
//	CollectionCreatorAddr string    `json:"collection_creator_addr"`
//	Collections     string 			`json:"collections"`
//	//Img             string 			`json:"img"`
//	Approve         string 			`json:"approve"`
//	Royalty         int 			`json:"royalty"`
//	Verified        string 			`json:"verified"`
//	Selltype        string 			`json:"selltype"`
//	Mintstate       string	 		`json:"mintstate"`
//	Likes	        int 			`json:"likes"`
//
//	Auction 		NftAuction		`json:"auction"`
//	Trans   		[]NftTran		`json:"trans"`
//	Bids    		[]NftBid	 	`json:"bids"`
//}
//
//func (nft NftDb) QuerySingleNft(contract, tokenId string) (NftSingleInfo, error) {
//	contract = strings.ToLower(contract)
//
//	var nftInfo NftSingleInfo
//
//	var nftRecord Nfts
//	err := nft.db.Where("contract = ? AND tokenid = ?", contract, tokenId).First(&nftRecord)
//	if err.Error != nil {
//		return NftSingleInfo{}, ErrNftNotExist
//	}
//	nftInfo.Name = nftRecord.Name
//	nftInfo.CreatorAddr = nftRecord.Createaddr
//	nftInfo.OwnerAddr = nftRecord.Ownaddr
//	nftInfo.Md5 = nftRecord.Md5
//	//nftInfo.AssetSample = nftRecord.Image
//	nftInfo.Desc = nftRecord.Desc
//	nftInfo.Meta =  nftRecord.Meta
//	nftInfo.SourceUrl = nftRecord.Url
//	nftInfo.NftContractAddr = nftRecord.Contract
//	nftInfo.NftTokenId = nftRecord.Tokenid
//	nftInfo.Categories = nftRecord.Categories
//	nftInfo.Collections = nftRecord.Collections
//	nftInfo.Approve = nftRecord.Approve
//	nftInfo.Royalty = nftRecord.Royalty
//	nftInfo.Verified = nftRecord.Verified
//	nftInfo.Selltype = nftRecord.Selltype
//	nftInfo.Mintstate = nftRecord.Mintstate
//	nftInfo.Likes = nftRecord.Favorited
//
//	user := Users{}
//	err = nft.db.Where("useraddr = ?", nftRecord.Createaddr).First(&user)
//	if err.Error == nil {
//		//nftInfo.CreatorPortrait = user.Portrait
//	}
//	user = Users{}
//	err = nft.db.Where("useraddr = ?", nftRecord.Ownaddr).First(&user)
//	if err.Error == nil {
//		//nftInfo.OwnerPortrait = user.Portrait
//	}
//	var collectRec Collects
//	err = nft.db.Where("Createaddr = ? AND name = ? ", nftRecord.Createaddr, nftRecord.Collections).First(&collectRec)
//	if err.Error == nil {
//		//nftInfo.Img = collectRec.Img
//		nftInfo.CollectionCreatorAddr = collectRec.Createaddr
//		nftInfo.Collectiondesc = collectRec.Desc
//	}
//
//	var auctionRec Auction
//	err = nft.db.Where("contract = ? AND tokenid = ?", contract, tokenId).First(&auctionRec)
//	if err.Error == nil {
//		nftInfo.Auction.Selltype = auctionRec.Selltype
//		nftInfo.Auction.Ownaddr = auctionRec.Ownaddr
//		nftInfo.Auction.NftTokenId = auctionRec.Tokenid
//		nftInfo.Auction.NftContractAddr = auctionRec.Contract
//		nftInfo.Auction.Paychan = auctionRec.Paychan
//		nftInfo.Auction.Currency = auctionRec.Currency
//		nftInfo.Auction.Startprice = auctionRec.Startprice
//		nftInfo.Auction.Endprice = auctionRec.Endprice
//		nftInfo.Auction.Startdate = auctionRec.Startdate
//		nftInfo.Auction.Enddate = auctionRec.Enddate
//		nftInfo.Auction.Tradesig = auctionRec.Tradesig
//	}
//
//	trans := make([]Trans, 0, 20)
//	err = nft.db.Where("contract = ? AND tokenid = ? AND selltype != ? AND selltype != ? AND price != ? ",
//		contract, tokenId, SellTypeMintNft.String(), SellTypeError.String(), 0).Find(&trans)
//	/*err = nft.db.Raw("SELECT * FROM trans\n WHERE id IN (SELECT MAX(id) AS o FROM trans GROUP BY contract, tokenId, Auctionid) " +
//	"and contract = ? and tokenid = ?  and \n  Selltype !=\"MintNft\"",
//	contract, tokenId).Find(&trans)*/
//	if err.Error == nil {
//		if err.RowsAffected != 0 {
//			for _, tran := range trans {
//				var nfttran NftTran
//				nfttran.NftContractAddr = tran.Contract
//				nfttran.Fromaddr = tran.Fromaddr
//				nfttran.Toaddr = tran.Toaddr
//				nfttran.NftTokenId = tran.Tokenid
//				nfttran.Transtime = tran.Transtime
//				nfttran.Paychan = tran.Paychan
//				nfttran.Currency = tran.Currency
//				nfttran.Price = tran.Price
//				nfttran.Selltype = tran.Selltype
//				nfttran.Txhash = tran.Txhash
//				nftInfo.Trans = append(nftInfo.Trans, nfttran)
//			}
//		}
//	}
//	bids := make([]Bidding, 0, 20)
//	err = nft.db.Where("contract = ? AND tokenid = ?", contract, tokenId).Find(&bids)
//	if err.Error == nil {
//		if err.RowsAffected != 0 {
//			for _, bid := range bids {
//				var nftbid NftBid
//				nftbid.Bidaddr = bid.Bidaddr
//				nftbid.NftTokenId = bid.Tokenid
//				nftbid.NftContractAddr = bid.Contract
//				nftbid.Paychan = bid.Paychan
//				nftbid.Currency = bid.Currency
//				nftbid.Price = bid.Price
//				nftbid.Bidtime = bid.Bidtime
//				nftbid.Tradesig = bid.Tradesig
//				nftInfo.Bids = append(nftInfo.Bids, nftbid)
//			}
//		}
//	}
//	return nftInfo, nil
//}

type UserNft struct {
	UserAddr        string `json:"user_addr"`
	CreatorAddr     string `json:"creator_addr"`
	OwnerAddr       string `json:"owner_addr"`
	Md5             string `json:"md5"`
	Name            string `json:"name"`
	Desc            string `json:"desc"`
	Meta            string `json:"meta"`
	SourceUrl       string `json:"source_url"`
	NftContractAddr string `json:"nft_contract_addr"`
	NftTokenId      string `json:"nft_token_id"`
	Categories      string `json:"categories"`
	Collections     string `json:"collections"`
	//AssetSample     string `json:"asset_sample"`
	Hide            string `json:"hide"`
}

func (nft NftDb) QueryUserNFTList(user_addr, start_index, count string) ([]UserNft, int, error) {
	user_addr = strings.ToLower(user_addr)
	if IsIntDataValid(start_index) != true {
		return nil, 0, ErrDataFormat
	}
	if IsIntDataValid(count) != true {
		return nil, 0, ErrDataFormat
	}
	var nftRecords []Nfts
	var recCount int64
	err := nft.db.Model(Nfts{}).Where("ownaddr = ?", user_addr).Count(&recCount)
	if err.Error != nil {
		fmt.Println("QueryUserNFTList() recCount err=", err)
		return nil, 0, ErrNftNotExist
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
		err = nft.db.Model(Nfts{}).Where("ownaddr = ?", user_addr).Limit(nftCount).Offset(startIndex).Find(&nftRecords)
		if err.Error != nil {
			fmt.Println("QueryUserNFTList() find record err=", err)
			return nil, 0, gorm.ErrRecordNotFound
		}
		userNfts := make([]UserNft, 0, 20)
		for i := 0; i < len(nftRecords); i++ {
			var userNft UserNft
			//userNft.UserAddr =
			userNft.CreatorAddr = nftRecords[i].Createaddr
			userNft.OwnerAddr = nftRecords[i].Ownaddr
			userNft.Md5 = nftRecords[i].Md5
			userNft.Name = nftRecords[i].Name
			userNft.Desc = nftRecords[i].Desc
			userNft.Meta = nftRecords[i].Meta
			userNft.SourceUrl = nftRecords[i].Url
			userNft.NftContractAddr = nftRecords[i].Contract
			userNft.NftTokenId = nftRecords[i].Tokenid
			userNft.Categories = nftRecords[i].Categories
			userNft.Collections = nftRecords[i].Collections
			//userNft.AssetSample = nftRecords[i].Image
			userNft.Hide = nftRecords[i].Hide
			userNfts = append(userNfts, userNft)
		}
		return userNfts, int(recCount), nil
	}
}

type UserCollection struct {
	CreatorAddr  string `json:"collection_creator_addr"`
	Name         string `json:"name"`
	//Img          string `json:"img"`
	ContractAddr string `json:"contract_addr"`
	Desc         string `json:"desc"`
	//Royalty      int    `json:"royalty"`
	Categories   string `json:"categories"`
	Totalcount   int    `json:"total_count"`
}

func (nft NftDb) QueryUserCollectionList(user_addr , start_index, count string) ([]UserCollection, int, error) {
	user_addr = strings.ToLower(user_addr)

	if IsIntDataValid(start_index) != true {
		return nil, 0, ErrDataFormat
	}
	if IsIntDataValid(count) != true {
		return nil, 0, ErrDataFormat
	}

	var collectRecs []Collects
	var recCount int64
	err := nft.db.Model(Collects{}).Where("createaddr = ?", user_addr).Count(&recCount)
	if err.Error != nil {
		fmt.Println("QueryUserCollectionList() recCount err=", err)
		return nil, 0, ErrNftNotExist
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
		err = nft.db.Model(Collects{}).Where("createaddr = ?", user_addr).Limit(nftCount).Offset(startIndex).Find(&collectRecs)
		if err.Error != nil {
			fmt.Println("QueryUserCollectionList() find record err=", err)
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
			userCollect.Totalcount = collectRecs[i].Totalcount
			userCollects = append(userCollects, userCollect)
		}
		return userCollects, int(recCount), nil
	}
}

//func (nft NftDb) QueryNFTCollectionList(start_index, count string) ([]UserCollection, int, error) {
//	var collectRecs []Collects
//	var recCount int64
//	if IsIntDataValid(start_index) != true {
//		return nil, 0, ErrDataFormat
//	}
//	if IsIntDataValid(count) != true {
//		return nil, 0, ErrDataFormat
//	}
//	err := nft.db.Model(Collects{}).Count(&recCount)
//	if err.Error != nil {
//		fmt.Println("QueryNFTCollectionList() recCount err=", err)
//		return nil, 0, ErrNftNotExist
//	}
//	startIndex, _ := strconv.Atoi(start_index)
//	nftCount, _ := strconv.Atoi(count)
//	if int64(startIndex) > recCount || recCount == 0{
//		return nil, 0, ErrNftNoMore
//	} else {
//		temp := recCount - int64(startIndex)
//		if int64(nftCount) > temp {
//			nftCount = int(temp)
//		}
//		err = nft.db.Model(Collects{}).Limit(nftCount).Offset(startIndex).Find(&collectRecs)
//		if err.Error != nil {
//			fmt.Println("QueryNFTCollectionList() find record err=", err)
//			return nil, 0, ErrNftNotExist
//		}
//		userCollects := make([]UserCollection, 0, 20)
//		for i := 0; i < len(collectRecs); i++ {
//			var userCollect UserCollection
//			userCollect.CreatorAddr = collectRecs[i].Createaddr
//			userCollect.Name = collectRecs[i].Name
//			//userCollect.Img = collectRecs[i].Img
//			userCollect.ContractAddr = collectRecs[i].Contract
//			userCollect.Desc = collectRecs[i].Desc
//			//userCollect.Royalty = collectRecs[i].Royalty
//			userCollect.Categories = collectRecs[i].Categories
//			userCollects = append(userCollects, userCollect)
//		}
//		return userCollects, int(recCount), nil
//	}
//}

//type TradingHistory struct {
//	NftContractAddr string `json:"nft_contract_addr"`
//	NftTokenId      string `json:"nft_token_id"`
//	NftName         string `json:"nft_name"`
//	Price           uint64 `json:"price"`
//	Count           uint64 `json:"count"`
//	From            string `json:"from"`
//	To              string `json:"to"`
//	Txhash 			string `json:"trade_hash"`
//	Selltype        string `json:"selltype"`
//	Date	        int64  `json:"date"`
//}
//
//func (nft NftDb) QueryUserTradingHistory(user_addr , start_index, count string) ([]TradingHistory, int, error) {
//	user_addr = strings.ToLower(user_addr)
//	if IsIntDataValid(start_index) != true {
//		return nil, 0, ErrDataFormat
//	}
//	if IsIntDataValid(count) != true {
//		return nil, 0, ErrDataFormat
//	}
//	var tranRecs []Trans
//	var recCount int64
//	db := nft.db.Model(Trans{}).Where("(toaddr = ? OR fromaddr = ?) AND (selltype != ? AND selltype != ?)",
//				user_addr, user_addr, SellTypeError.String(), SellTypeMintNft.String()).Count(&recCount)
//	if db.Error != nil {
//		fmt.Println("QueryUserTradingHistory() recCount err=", db)
//		return nil, 0, ErrNoTrans
//	}
//	if recCount == 0 {
//		fmt.Println("QueryUserTradingHistory() recCount == 0")
//		return nil, 0, ErrNoTrans
//	}
//
//	startIndex, _ := strconv.Atoi(start_index)
//	nftCount, _ := strconv.Atoi(count)
//	if int64(startIndex) > recCount || recCount == 0{
//		return nil, 0, ErrNftNoMore
//	} else {
//		temp := recCount - int64(startIndex)
//		if int64(nftCount) > temp {
//			nftCount = int(temp)
//		}
//		err := db.Model(Trans{}).Limit(nftCount).Offset(startIndex).Find(&tranRecs)
//		if err.Error != nil {
//			fmt.Println("QueryUserTradingHistory() find record err=", err)
//			return nil, 0, ErrNftNotExist
//		}
//		trans := make([]TradingHistory, 0, 20)
//		for i := 0; i < len(tranRecs); i++ {
//			var tran TradingHistory
//			tran.NftContractAddr = tranRecs[i].Contract
//			tran.NftTokenId = tranRecs[i].Tokenid
//			tran.NftName = tranRecs[i].Name
//			tran.Price = tranRecs[i].Price
//			tran.Count = 1
//			tran.From = tranRecs[i].Fromaddr
//			tran.To = tranRecs[i].Toaddr
//			tran.Date = tranRecs[i].Transtime
//			tran.Selltype = tranRecs[i].Selltype
//			tran.Txhash = tranRecs[i].Txhash
//			trans = append(trans, tran)
//		}
//		return trans, int(recCount), nil
//	}
//}

//func (nft NftDb) QueryMarketTradingHistory(filter []StQueryField, sort []StSortField,
//	start_index string, count string) ([]TradingHistory, int, error) {
//	var tranRecs []Trans
//	var recCount int64
//	var queryWhere string
//	var orderBy string
//
//	if len(filter) > 0 {
//		queryWhere = nft.joinFilters(filter)
//	}
//	if len(sort) > 0 {
//		for k, v := range sort {
//			if k >0 {
//				orderBy = orderBy + ", "
//			}
//			orderBy = v.By + " " + v.Order
//		}
//	} else {
//		orderBy = "transtime desc"
//	}
//
//	tx := nft.db.Model(Trans{})
//	if len(queryWhere) > 0 {
//		tx = tx.Where(queryWhere)
//	}
//	tx = tx.Where("selltype != ? AND selltype != ?",
//		SellTypeError.String(), SellTypeMintNft.String())
//	if len(orderBy) > 0 {
//		tx = tx.Order(orderBy)
//	}
//	tx = tx.Count(&recCount)
//	if tx.Error != nil {
//		fmt.Println("QueryMarketTradingHistory() recCount err=", tx.Error)
//		return nil, 0, ErrNftNotExist
//	}
//	startIndex, _ := strconv.Atoi(start_index)
//	nftCount, _ := strconv.Atoi(count)
//	if int64(startIndex) > recCount || recCount == 0{
//		return nil, 0, ErrNftNoMore
//	} else {
//		temp := recCount - int64(startIndex)
//		if int64(nftCount) > temp {
//			nftCount = int(temp)
//		}
//		tx = tx.Limit(nftCount).Offset(startIndex).Find(&tranRecs)
//		if tx.Error != nil {
//			fmt.Println("QueryMarketTradingHistory() find record err=", tx.Error)
//			return nil, 0, ErrNftNotExist
//		}
//		//var trans []TradingHistory
//		trans := make([]TradingHistory, 0, 20)
//		for i := 0; i < len(tranRecs); i++ {
//			var tran TradingHistory
//			tran.NftContractAddr = tranRecs[i].Contract
//			tran.NftTokenId = tranRecs[i].Tokenid
//			tran.NftName = tranRecs[i].Name
//			tran.Price = tranRecs[i].Price
//			tran.Count = 1
//			tran.From = tranRecs[i].Fromaddr
//			tran.To = tranRecs[i].Toaddr
//			tran.Date = tranRecs[i].Transtime
//			tran.Selltype = tranRecs[i].Selltype
//			tran.Txhash =  tranRecs[i].Txhash
//			trans = append(trans, tran)
//		}
//		return trans, int(recCount), nil
//	}
//}

type UserOffer struct {
	Contract	 string `json:"nft_contract_addr"`
	Tokenid      string `json:"nft_token_id"`
	Name      	 string `json:"name"`
	Price        uint64 `json:"price"`
	Count        uint64 `json:"count"`
	Bidtime      int64  `json:"date"`
}

func (nft NftDb) QueryUserOfferList(user_addr , start_index, count string) ([]UserOffer, int, error) {
	user_addr = strings.ToLower(user_addr)
	if IsIntDataValid(start_index) != true {
		return nil, 0, ErrDataFormat
	}
	if IsIntDataValid(count) != true {
		return nil, 0, ErrDataFormat
	}
	var Recount int64
	sql := "SELECT biddings.contract as Contract, biddings.tokenid as Tokenid, nfts.name as Name, biddings.price as Price, " +
		"biddings.count as Count, biddings.bidtime as Bidtime FROM biddings LEFT JOIN nfts ON biddings.contract = nfts.contract AND biddings.tokenid = nfts.tokenid " +
		"WHERE ownaddr = ? AND biddings.deleted_at is null"
	sqlCount := "SELECT count(*) as Reccnt FROM biddings LEFT JOIN nfts ON biddings.contract = nfts.contract AND biddings.tokenid = nfts.tokenid " +
	"WHERE ownaddr = ? AND biddings.deleted_at is null"
	err := nft.db.Raw(sqlCount, user_addr).Scan(&Recount)
	if err.Error != nil {
		fmt.Println("QueryUserInfo() query Sum err=", err)
		return nil, 0, err.Error
	}
	sql = sql + " limit" +" " + start_index + "," + count
	var useroffer []UserOffer
	err = nft.db.Raw(sql, user_addr).Scan(&useroffer)
	if err.Error != nil {
		fmt.Println("QueryUserInfo() query Sum err=", err)
		return nil, 0, err.Error
	}
	return useroffer, int(Recount), nil
}

type UserBid struct {
	NftContractAddr string `json:"nft_contract_addr"`
	NftTokenId      string `json:"nft_token_id"`
	Name      	 	string `json:"name"`
	Price           uint64 `json:"price"`
	Count           uint64 `json:"count"`
	Date            int64  `json:"date"`
}

func (nft NftDb) QueryUserBidList(user_addr , start_index, count string) ([]UserBid, int, error) {
	user_addr = strings.ToLower(user_addr)

	if IsIntDataValid(start_index) != true {
		return nil, 0, ErrDataFormat
	}
	if IsIntDataValid(count) != true {
		return nil, 0, ErrDataFormat
	}
	var offerRecs []Bidding
	var recCount int64
	err := nft.db.Model(Bidding{}).Where("Bidaddr = ?", user_addr).Count(&recCount)
	if err.Error != nil {
		fmt.Println("QueryUserBidList() recCount err=", err)
		return nil, 0, ErrNftNotExist
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
		err = nft.db.Model(Bidding{}).Where("Bidaddr = ?", user_addr).Limit(nftCount).Offset(startIndex).Find(&offerRecs)
		if err.Error != nil {
			fmt.Println("QueryUserBidList() find record err=", err)
			return nil, 0, ErrNftNotExist
		}
		userBids := make([]UserBid, 0, 20)
		for i := startIndex; i < len(offerRecs); i++ {
			var userBid UserBid
			userBid.NftContractAddr = offerRecs[i].Contract
			userBid.NftTokenId = offerRecs[i].Tokenid
			userBid.Price = offerRecs[i].Price
			userBid.Count = 1
			userBid.Date = offerRecs[i].Bidtime
			nftrec := Nfts{}
			err := nft.db.Model(&Nfts{}).Where("contract = ? AND tokenid = ?",
				userBid.NftContractAddr, userBid.NftTokenId).First(&nftrec)
			if err.Error == nil {
				userBid.Name = nftrec.Name
			}
			userBids = append(userBids, userBid)
		}
		return userBids, int(recCount), nil
	}
}

type UserFavorite struct {
	CreatorAddr  	string `json:"collection_creator_addr"`
	NftContractAddr string `json:"nft_contract_addr"`
	NftTokenId      string `json:"nft_token_id"`
	Name            string `json:"name"`
	//AssetSample     string `json:"asset_sample"`
	Collections     string `json:"collections"`
	//Img             string `json:"img"`
}

func (nft NftDb) QueryUserFavoriteList(user_addr , start_index, count string) ([]UserFavorite, int, error) {
	user_addr = strings.ToLower(user_addr)
	if IsIntDataValid(start_index) != true {
		return nil, 0, ErrDataFormat
	}
	if IsIntDataValid(count) != true {
		return nil, 0, ErrDataFormat
	}
	var favoritedRecs []NftFavorited
	var recCount int64
	err := nft.db.Model(NftFavorited{}).Where("useraddr = ?", user_addr).Count(&recCount)
	if err.Error != nil {
		fmt.Println("QueryUserFavoriteList() recCount err=", err)
		return nil, 0, ErrNftNotExist
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
		err = nft.db.Model(NftFavorited{}).Where("useraddr = ?", user_addr).Limit(nftCount).Offset(startIndex).Find(&favoritedRecs)
		if err.Error != nil {
			fmt.Println("QueryUserCollectionList() find record err=", err)
			return nil, 0, ErrNftNotExist
		}
		userFavorites := make([]UserFavorite, 0, 20)
		for i := 0; i < len(favoritedRecs); i++ {
			var favorite UserFavorite
			favorite.NftContractAddr = favoritedRecs[i].Contract
			favorite.NftTokenId = favoritedRecs[i].Tokenid
			favorite.Name = favoritedRecs[i].Name
			//favorite.Img = favoritedRecs[i].Img
			//favorite.AssetSample = favoritedRecs[i].Image
			favorite.CreatorAddr = favoritedRecs[i].Collectcreator
			favorite.Collections = favoritedRecs[i].Collections
			userFavorites = append(userFavorites, favorite)
		}
		return userFavorites, int(recCount), nil
	}
}

func (nft NftDb) AddUserFavor(userAddr, favoritedaddr string) error {
	userAddr = strings.ToLower(userAddr)
	favoritedaddr = strings.ToLower(favoritedaddr)
	var favorrecord UserFavorited
	err := nft.db.Where("favoritedaddr = ? AND useraddr = ?", favoritedaddr, userAddr).First(&favorrecord)
	if err.Error == nil {
		fmt.Println("AddUserFavor() UserFavorited already exist." )
		return ErrAlreadyUserFavorited
	}
	favorrecord = UserFavorited{}
	favorrecord.Useraddr = userAddr
	favorrecord.Favoritedaddr = favoritedaddr
	return nft.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&favorrecord).Create(&favorrecord)
		if err.Error != nil {
			fmt.Println("AddUserFavor() create record err=", err.Error)
			return err.Error
		}
		user := Users{}
		err = tx.Where("useraddr = ?", favoritedaddr).First(&user)
		if err.Error != nil {
			fmt.Println("AddUserFavor() find err= ", err.Error )
			return err.Error
		}
		err = tx.Model(&user).Where("useraddr = ?", favoritedaddr).Update("Favorited",user.Favorited+1)
		if err.Error != nil {
			fmt.Println("AddUserFavor() update NftFavorited err= ", err.Error )
			return err.Error
		}
		return nil
	})
}

func (nft NftDb) DelUserFavor(userAddr, favoritedaddr  string) error {
	userAddr = strings.ToLower(userAddr)
	favoritedaddr = strings.ToLower(favoritedaddr)

	var favorrecord UserFavorited
	err := nft.db.Where("favoritedaddr = ? AND useraddr = ?", favoritedaddr, userAddr).First(&favorrecord)
	if err.Error != nil {
		fmt.Println("DelUserFavor() err= ", err.Error )
		return err.Error
	}
	return nft.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&UserFavorited{}).Where("favoritedaddr = ? AND useraddr = ?", favoritedaddr, userAddr).Delete(&UserFavorited{})
		if err.Error != nil {
			if err.Error == gorm.ErrRecordNotFound {
				return ErrNotNftFavorited
			}
			fmt.Println("DelUserFavor() err=", err.Error )
			return err.Error
		}
		user := Users{}
		err = tx.Model(&user).Where("useraddr = ?", favoritedaddr).First(&user)
		if err.Error != nil {
			fmt.Println("DelUserFavor() find err= ", err.Error )
			return err.Error
		}
		err = tx.Model(&user).Where("useraddr = ?", favoritedaddr).Update("Favorited",user.Favorited-1)
		if err.Error != nil {
			fmt.Println("DelUserFavor() update err= ", err.Error )
			return err.Error
		}
		return nil
	})
}

func (nft NftDb) QueryUserFavorited(userAddr string)  ([]UserFavorited, error){
	userAddr = strings.ToLower(userAddr)

	favors := []UserFavorited{}
	err := nft.db.Where("favoruseraddr = ?", userAddr).Find(&favors)
	if err.Error != nil {
		fmt.Println("queryNft, err=\n ", err.Error)
		return nil, err.Error
	}
	marshal, _ := json.Marshal(favors)
	fmt.Printf("%s\n",string(marshal))
	//return string(marshal), nil
	//return marshal, nil
	return favors, err.Error
}

//获取市场数据
func (nft *NftDb) QueryMarketInfo() (uint64, error){
	transData := []Trans{}
	var totalAmount7 uint64
	//统计七天内交易金额
	before7daysTime := time.Now().AddDate(0, 0, -7)
	before7Date := time.Date(before7daysTime.Year(), before7daysTime.Month(), before7daysTime.Day(),
		0, 0, 0, 0, time.Local)
	currentDate := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(),
		0, 0, 0, 0, time.Local)
	fmt.Println(before7Date, currentDate)
	findResult := nft.db.Where("transtime >= ? and transtime <= ?", before7Date, currentDate).Find(&transData)
	if findResult.Error != nil {
		return 0, findResult.Error
	}
	for _, row := range transData {
		totalAmount7 = totalAmount7 + row.Price
	}

	return totalAmount7, nil
}

func (nft NftDb) TextAndHash(data []byte) ([]byte, string) {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), string(data))
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write([]byte(msg))
	return hasher.Sum(nil), msg
}

func (nft NftDb) GetEthAddr(msg string, sigStr string) (common.Address, error){
	sigData := hexutil.MustDecode(sigStr)
	if len(sigData) != 65 {
		return common.Address{}, fmt.Errorf("signature must be 65 bytes long")
	}
	if sigData[64] != 27 && sigData[64] != 28 {
		return common.Address{}, fmt.Errorf("invalid Ethereum signature (V is not 27 or 28)")
	}
	sigData[64] -= 27
	hash, _ := NftDb{}.TextAndHash([]byte(msg))
	rpk, err := crypto.SigToPub(hash, sigData)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(*rpk), nil
}

func (nft *NftDb) isValidVerifyAddr(rawData string, sig string) (bool, error) {
	addrList, err := ethhelper.AdminList()
	if err != nil {
		return false, err
	}

	verificationAddr, err := nft.GetEthAddr(rawData, sig)
	if err != nil {
		 return false, err
	}
	verificationAddrS := verificationAddr.String()

	for _, addr := range addrList {
		if verificationAddrS == addr {
			return true, nil
		}
	}

	return false, errors.New("verification address is invalid")
}

func (nft NftDb) InsertSigData(SigData, msg string) error  {
	sigmsg := Sigmsgs{}
	sigmsg.Signdata = SigData
	sigmsg.Signmsg = msg
	db := nft.db.Model(&sigmsg).Create(&sigmsg)
	if db.Error != nil {
		fmt.Println("InsertSigData()->create() err=", db.Error)
		return db.Error
	}
	return nil
}

//func (nft NftDb) QuerySysParams() (SysParamsRec, error) {
//	var params SysParams
//	err := nft.db.Last(&params)
//	if err.Error != nil {
//		if err.Error == gorm.ErrRecordNotFound {
//			params = SysParams{}
//			params.Exchangaddr = strings.ToLower(initExchangAddr)
//			params.Nftaddr = strings.ToLower(initNftAddr)
//			params.Lowprice = initLowprice
//			params.Royaltylimit = initRoyaltylimit
//			params.Categories = initCategories
//			err = nft.db.Model(&SysParams{}).Create(&params)
//			if err.Error != nil {
//				fmt.Println("SetSysParams() create SysParams err= ", err.Error )
//				return SysParamsRec{}, err.Error
//			}
//		} else {
//			fmt.Println("QuerySysParams() not find err=", err.Error)
//			return SysParamsRec{}, err.Error
//		}
//	}
//	return params.SysParamsRec, err.Error
//}
//
//func (nft NftDb) SetSysParams(param SysParamsRec) error {
//	var paramRec, updateP SysParams
//	err := nft.db.Last(&paramRec)
//	if err.Error != nil {
//		if nft.db.Error == gorm.ErrRecordNotFound {
//			updateP.Exchangaddr = initExchangAddr
//			updateP.Nftaddr = initNftAddr
//			updateP.Lowprice = initLowprice
//			updateP.Royaltylimit = initRoyaltylimit
//			updateP.Categories = initCategories
//		} else {
//			fmt.Println("QuerySysParams() not find err=", err.Error)
//			return err.Error
//		}
//	} else {
//		if param.Exchangaddr != "" {
//			updateP.Exchangaddr = param.Exchangaddr
//		} else{
//			updateP.Exchangaddr = paramRec.Exchangaddr
//		}
//		if param.Nftaddr != "" {
//			updateP.Nftaddr = param.Nftaddr
//		} else {
//			updateP.Nftaddr = paramRec.Nftaddr
//		}
//		if param.Lowprice != 0 {
//			updateP.Lowprice = param.Lowprice
//		} else {
//			updateP.Lowprice = paramRec.Lowprice
//		}
//	}
//	updateP.Signdata = param.Signdata
//	err = nft.db.Model(&SysParams{}).Create(&updateP)
//	if err.Error != nil {
//		fmt.Println("SetSysParams() create SysParams err= ", err.Error )
//		return err.Error
//	}
//	return nil
//}
//
//func InitSysParams(Sqldsndb string) {
//	nd, err := NewNftDb(Sqldsndb)
//	if err != nil {
//		fmt.Printf("connect database err = %s\n", err)
//	}
//	params, err := nd.QuerySysParams()
//	if err != nil {
//		ExchangAddr = initExchangAddr
//		NftAddr = initNftAddr
//		Lowprice = initLowprice
//		RoyaltyLimit = initRoyaltylimit
//	} else {
//		ExchangAddr = params.Exchangaddr
//		NftAddr = params.Nftaddr
//		Lowprice = params.Lowprice
//		RoyaltyLimit = params.Royaltylimit
//	}
//	nd.Close()
//}

//func (nft NftDb) NewCollections(useraddr, name, img, contract_type, contract_addr,
//	desc, categories, sig string) error {
//	useraddr = strings.ToLower(useraddr)
//	contract_addr = strings.ToLower(contract_addr)
//
//	var collectRec Collects
//	err := nft.db.Where("Createaddr = ? AND name = ? ", useraddr, name).First(&collectRec)
//	if err.Error == nil {
//		fmt.Println("NewCollections() err=Collection already exist." )
//		return ErrCollectionExist
//	} else if err.Error == gorm.ErrRecordNotFound {
//		fmt.Println("NewCollections() err=Collection already exist.")
//		collectRec = Collects{}
//		collectRec.Createaddr = useraddr
//		collectRec.Name = name
//		collectRec.Desc = desc
//		collectRec.Img = img
//		collectRec.Contract = contract_addr
//		collectRec.Contracttype = contract_type
//		collectRec.Categories = categories
//		collectRec.SigData = sig
//		return nft.db.Transaction(func(tx *gorm.DB) error {
//			err := tx.Model(&Collects{}).Create(&collectRec)
//			if err.Error != nil {
//				fmt.Println("NewCollections() err=", err.Error)
//				return err.Error
//			}
//			return nil
//		})
//	}
//	fmt.Println("NewCollections() dbase err=.", err)
//	return err.Error
//}
//
//func (nft NftDb) ModifyCollections(useraddr, name, img, contract_type, contract_addr,
//	desc, categories, sig string) error {
//	useraddr = strings.ToLower(useraddr)
//	contract_addr = strings.ToLower(contract_addr)
//	var collectRec Collects
//	err := nft.db.Where("Createaddr = ? AND name = ? ", useraddr, name).First(&collectRec)
//	if err.Error != nil {
//		fmt.Println("NewCollections() err=Collection not exist." )
//		return ErrCollectionNotExist
//	}
//	collectRec = Collects{}
//	if img != "" {
//		collectRec.Img = img
//	}
//	if contract_type != "" {
//		collectRec.Contracttype = contract_type
//	}
//	if contract_addr != "" {
//		collectRec.Contract = contract_addr
//	}
//	if desc != "" {
//		collectRec.Desc = desc
//	}
//	if categories != "" {
//		collectRec.Categories = categories
//	}
//	collectRec.SigData = sig
//	return nft.db.Transaction(func(tx *gorm.DB) error {
//		err := tx.Model(&Collects{}).Where("Createaddr = ? AND name = ? ", useraddr, name).Updates(&collectRec)
//		if err.Error != nil {
//			fmt.Println("NewCollections() err=", err.Error)
//			return err.Error
//		}
//		return nil
//	})
//}

func (nft NftDb) ModifyCollectionsImage(name, collection_creator_addr, img, sig string) error {
	collection_creator_addr = strings.ToLower(collection_creator_addr)
	var collectRec Collects
	err := nft.db.Where("createaddr = ? AND name = ?", collection_creator_addr, name).First(&collectRec)
	if err.Error != nil {
		fmt.Println("modifyCollectionsImage() err=Collection not exist." )
		return ErrCollectionNotExist
	}
	collectRec = Collects{}
	collectRec.Img = img
	collectRec.SigData = sig
	return nft.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&Collects{}).Where("createaddr = ? AND name = ?", collection_creator_addr, name).Updates(&collectRec)
		if err.Error != nil {
			fmt.Println("NewCollections() err=", err.Error)
			return err.Error
		}
		return nil
	})
}

func (nft NftDb) SaveHistoryTrans(NftContractAddr, NftTokenId, Price, Count, From, To, Date string) error  {
	NftContractAddr = strings.ToLower(NftContractAddr)
	From = strings.ToLower(From)
	To = strings.ToLower(To)
	if IsUint64DataValid(Price) != true {
		return ErrPrice
	}
	if IsIntDataValid(Count) != true {
		return ErrDataFormat
	}
	return nft.db.Transaction(func(tx *gorm.DB) error {
		trans := Trans{}
		trans.Contract = NftContractAddr
		trans.Fromaddr = From
		trans.Toaddr = To
		trans.Tokenid = NftTokenId
		trans.Price, _ = strconv.ParseUint(Price, 10, 64)
		trans.Transtime, _ = strconv.ParseInt(Date, 10, 64)
			trans.Selltype = SellTypeForeignPrice.String()
		err := tx.Model(&trans).Create(&trans)
		if err.Error != nil {
			fmt.Println("SaveHistoryTrans() create trans record err=", err.Error)
			return err.Error
		}
		nftrecord := Nfts{}
		nftrecord.Ownaddr = To
		nftrecord.Selltype = SellTypeNotSale.String()
		nftrecord.Transprice = trans.Price
		err = tx.Model(&Nfts{}).Where("contract = ? AND tokenid =?",
			NftContractAddr, NftTokenId).Updates(&nftrecord)
		if err.Error != nil {
			fmt.Println("SaveHistoryTrans() update record err=", err.Error)
			return err.Error
		}
		fmt.Println("SaveHistoryTrans() from != \"\" && to != \"\" --> price != \"\" OK" )
		return nil
	})
}

func (nft NftDb) HasCollectionsImage(contract_addr string) (bool, error) {
	contract_addr = strings.ToLower(contract_addr)
	var collectRec Collects
	err := nft.db.Where("Contract = ?", contract_addr).First(&collectRec)
	if err.Error != nil {
		fmt.Println("HasCollectionsImage() dbase err=", err )
		return false, err.Error
	}
	if collectRec.Img != "" {
		return true, nil
	} else {
		return false, nil
	}
}

func (nft NftDb) HasContractAddr(contract_addr string) (bool, error) {
	contract_addr = strings.ToLower(contract_addr)
	var nfttab Nfts
	err := nft.db.Where("contract = ?", contract_addr).First(&nfttab)
	if err.Error != nil {
		if err.Error == gorm.ErrRecordNotFound {
			fmt.Println("HasContractAddr() contract not exist.")
			return false, nil
		}
		fmt.Println("HasContractAddr() dbase err=", err)
		return true, err.Error
	} else {
		return true, nil
	}
}

func (nft *NftDb) QueryPendingKYCList(unverified string) ([]Users, error) {
	users := []Users{}
	queryResult := nft.db.Where("verified = ?", unverified).Find(&users)
	if queryResult.Error != nil {
		return nil, queryResult.Error
	}

	return users, nil
}

//审核用户KYC
func (nft NftDb) UserKYC(vrfaddr string, useraddr string, desc string,
	verified string, sig string) error {
	vrfaddr = strings.ToLower(vrfaddr)
	useraddr = strings.ToLower(useraddr)

	user := Users{}

	takeResult := nft.db.Where("useraddr = ?", useraddr).Take(&user)
	if takeResult.Error != nil {
		return takeResult.Error
	}
	updateValue := make(map[string]interface{})
	updateValue["verifyaddr"] = vrfaddr
	updateValue["desc"] = desc
	updateValue["verified"] = verified
	updateValue["signdata"] = sig
	updateResult := nft.db.Model(&user).Updates(updateValue)
	if updateResult.Error != nil {
		return updateResult.Error
	}
	return nil
}

//申请用户KYC*
func (nft NftDb) UserRequireKYC(useraddr string, country string, pic string, sig string) error {
	useraddr = strings.ToLower(useraddr)

	user := Users{}

	takeResult := nft.db.Where("useraddr = ?", useraddr).Take(&user)
	if takeResult.Error != nil {
		return takeResult.Error
	}
	updateValue := make(map[string]interface{})
	updateValue["kycpic"] = pic
	updateValue["signdata"] = sig
	updateValue["verified"] = NoVerify.String()
	updateValue["country"] = country
	updateResult := nft.db.Model(&user).Updates(updateValue)
	if updateResult.Error != nil {
		return updateResult.Error
	}
	return nil
}


func (nft NftDb) AskForApprove(nft_contract_addr, nft_token_id string) (UserNft, error) {
	nft_contract_addr = strings.ToLower(nft_contract_addr)
	nftRecords := Nfts{}
	err := nft.db.Where("contract = ? AND tokenid = ? ", nft_contract_addr, nft_token_id).First(&nftRecords)
	if err.Error == gorm.ErrRecordNotFound {
	}
	var userNft UserNft
	userNft.CreatorAddr = nftRecords.Createaddr
	userNft.OwnerAddr = nftRecords.Ownaddr
	userNft.Md5 = nftRecords.Md5
	userNft.Name = nftRecords.Name
	userNft.Desc = nftRecords.Desc
	userNft.Meta = nftRecords.Meta
	userNft.SourceUrl = nftRecords.Url
	userNft.NftContractAddr = nftRecords.Contract
	userNft.NftTokenId = nftRecords.Tokenid
	userNft.Categories = nftRecords.Categories
	userNft.Collections = nftRecords.Collections
	//userNft.AssetSample = nftRecords.Image
	userNft.Hide = nftRecords.Hide
	return userNft, nil
}

func (nft NftDb) QueryAnnouncement() ([]Announcements, error) {
	var announce []Announcements
	err := nft.db.Order("id desc").Limit(5).Find(&announce)
	if err.Error != nil {
		fmt.Println("QueryAnnouncement() not find err=", err.Error)
		return nil, err.Error
	}
	return announce, err.Error
}

func (nft NftDb) QuerySingleAnnouncement(title string) (Announcements, error) {
	var announce Announcements
	err := nft.db.Where("title = ?", title).Last(&announce)
	if err.Error != nil {
		fmt.Println("QuerySingleAnnouncement() not find err=", err.Error)
		return Announcements{}, err.Error
	}
	return announce, err.Error
}

func (nft NftDb) SetAnnouncement(title, content string) error {
	var announce Announcements
	announce.Title = title
	announce.Content = content
	err := nft.db.Model(&Announcements{}).Create(&announce)
	if err.Error != nil {
		fmt.Println("SetAnnouncement() create  err= ", err.Error )
		return err.Error
	}
	return nil
}



func (nft *NftDb) IsValidCategory(category string) bool {
	sysParams := SysParams{}

	result := nft.db.Model(&SysParams{}).Last(&sysParams)
	if result.Error != nil {
		return false
	}

	categories := strings.Split(sysParams.Categories, ",")
	for _, v := range categories {
		if v == category {
			return true
		}
	}
	return false
}


func IsIntDataValid(dataStr string) bool {
	if dataStr == "" {
		return true
	}
	data, err := strconv.Atoi(dataStr)
	if err != nil {
		return false
	}
	if data < 0 {
		return false
	}
	return true
}

func IsUint64DataValid(dataStr string) bool  {
	if dataStr == "" {
		return true
	}
	data, err := strconv.ParseUint(dataStr, 10, 64)
	if err != nil {
		return false
	}
	if data < 0 {
		return false
	}
	return true
}