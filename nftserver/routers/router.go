package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/nftexchange/nftserver/controllers/nftexchangev1"
	"github.com/nftexchange/nftserver/controllers/nftexchangev2"
	"github.com/nftexchange/nftserver/models"
)

func init() {
	registFilters()
	registRouterV1()
	registRouterV2()
	nftdb := new(models.NftDb)
	nftdb.InitDb(models.SqlSvr, "nftdb")
}

func registRouterV1() {
	//用户登录
	beego.Router("/v1/login", &nftexchangev1.NftExchangeControllerV1{}, "post:UserLogin")
	//上传nft作品
	beego.Router("/v1/upload", &nftexchangev1.NftExchangeControllerV1{}, "post:UploadNft")
	//购买nft作品
	beego.Router("/v1/buy", &nftexchangev1.NftExchangeControllerV1{}, "post:BuyNft")
	//查询nft作品
	beego.Router("/v1/queryNFTList", &nftexchangev1.NftExchangeControllerV1{}, "get:QueryAllNftProducts")
	//查询用户
	beego.Router("/v1/queryUser", &nftexchangev1.NftExchangeControllerV1{}, "post:QueryUserInfo")
	//查询用户
	//beego.Router("/v1/getimage", &controllers.NftExchangeController{}, "post:GetImageFromIPFS")
	beego.Router("/v1/ipfsaddtest", &nftexchangev1.NftExchangeControllerV1{}, "get:IpfsTest")
}

func registRouterV2() {
	//用户登录
	beego.Router("/v2/login", &nftexchangev2.NftExchangeControllerV2{}, "post:UserLogin")
	//上传nft作品
	beego.Router("/v2/upload", &nftexchangev2.NftExchangeControllerV2{}, "post:UploadNft")
	//购买nft作品
	beego.Router("/v2/buy", &nftexchangev2.NftExchangeControllerV2{}, "post:BuyNft")
	//取消购买nft作品
	beego.Router("/v2/cancelBuy", &nftexchangev2.NftExchangeControllerV2{}, "post:CancelBuyNft")

	//购买nft作品完成通知
	beego.Router("/v2/buyResult", &nftexchangev2.NftExchangeControllerV2{}, "post:BuyResult")
	//查询nft作品
	beego.Router("/v2/queryNFTList", &nftexchangev2.NftExchangeControllerV2{}, "post:QueryNftList")
	//查询用户
	beego.Router("/v2/queryUserInfo", &nftexchangev2.NftExchangeControllerV2{}, "post:QueryUserInfo")
	//修改用户信息
	beego.Router("/v2/modifyUserInfo", &nftexchangev2.NftExchangeControllerV2{}, "post:ModifyUserInfo")
	//查询单个NFT信息
	beego.Router("/v2/queryNFT", &nftexchangev2.NftExchangeControllerV2{}, "post:QueryNft")
	//售卖(上架)
	beego.Router("/v2/sell", &nftexchangev2.NftExchangeControllerV2{}, "post:Sell")
	//取消售卖(下架)
	beego.Router("/v2/cancelSell", &nftexchangev2.NftExchangeControllerV2{}, "post:CancelSell")
	//owner修改价格
	//beego.Router("/v2/modifyNFT", &controllers.NftExchangeControllerV2{}, "post:ModifyNFT")
	//查询NFT待审核列表
	beego.Router("/v2/queryPendingVrfList", &nftexchangev2.NftExchangeControllerV2{}, "get:QueryPendingVerificationList")
	//审核NFT
	beego.Router("/v2/vrfNFT", &nftexchangev2.NftExchangeControllerV2{}, "post:VerifyNft")
	//获取市场数据
	beego.Router("/v2/queryMarketInfo", &nftexchangev2.NftExchangeControllerV2{}, "post:QueryMarketInfo")
	//用户关注的NFT
	beego.Router("/v2/like", &nftexchangev2.NftExchangeControllerV2{}, "post:SetFavoriteNft")
	//新建集合
	beego.Router("/v2/newCollections", &nftexchangev2.NftExchangeControllerV2{}, "post:CreateCollection")
	//修改集合信息
	beego.Router("/v2/modifyCollections", &nftexchangev2.NftExchangeControllerV2{}, "post:ModifyCollection")
	//修改集合封面
	beego.Router("/v2/modifyCollectionsImage", &nftexchangev2.NftExchangeControllerV2{}, "post:ModifyCollectionsImage")
	//获取用户NFT列表
	beego.Router("/v2/queryUserNFTList", &nftexchangev2.NftExchangeControllerV2{}, "post:QueryUserNftList")
	//查询用户集合列表
	beego.Router("/v2/queryUserCollectionList", &nftexchangev2.NftExchangeControllerV2{}, "post:QueryUserCollectionList")
	//查询用户交易历史
	beego.Router("/v2/queryUserTradingHistory", &nftexchangev2.NftExchangeControllerV2{}, "post:QueryUserTradingHistory")
	//获取用户NFT出价列表
	beego.Router("/v2/queryUserOfferList", &nftexchangev2.NftExchangeControllerV2{}, "post:QueryUserOfferList")
	//获取用户NFT被出价列表
	beego.Router("/v2/queryUserBidList", &nftexchangev2.NftExchangeControllerV2{}, "post:QueryUserBidList")
	//获取用户NFT关注列表
	beego.Router("/v2/queryUserFavoriteList", &nftexchangev2.NftExchangeControllerV2{}, "post:QueryUserFavoriteList")
	//获取NFT合集列表
	beego.Router("/v2/queryNFTCollectionList", &nftexchangev2.NftExchangeControllerV2{}, "post:QueryNftCollectionList")
	//获取合集详情
	beego.Router("/v2/queryCollectionInfo", &nftexchangev2.NftExchangeControllerV2{}, "post:QueryCollectionInfo")

	//获取NFT交易所的交易历史
	beego.Router("/v2/queryMarketTradingHistory", &nftexchangev2.NftExchangeControllerV2{}, "post:QueryMarketTradingHistory")
	//用户申请KYC
	beego.Router("/v2/userRequireKYC", &nftexchangev2.NftExchangeControllerV2{}, "post:UserRequireKYC")
	//审核KYC
	beego.Router("/v2/userKYC", &nftexchangev2.NftExchangeControllerV2{}, "post:UserKYC")
	//获取用户KYC列表
	beego.Router("/v2/queryPendingKYCList", &nftexchangev2.NftExchangeControllerV2{}, "get:QueryPendingKYCList")
	//模糊查询nft,集合,账户
	beego.Router("/v2/search", &nftexchangev2.NftExchangeControllerV2{}, "post:Search")
	//获取交易所授权
	beego.Router("/v2/askForApprove", &nftexchangev2.NftExchangeControllerV2{}, "post:AskForApprove")
	//获取首页信息
	beego.Router("/v2/queryHomePage", &nftexchangev2.NftExchangeControllerV2{}, "get:QueryHomePage")
	//获得sysparams数据
	beego.Router("/v2/get_sys_para", &nftexchangev2.NftExchangeControllerV2{}, "get:GetSysParams")
	//设置sysparams数据
	beego.Router("/v2/set_sys_para", &nftexchangev2.NftExchangeControllerV2{}, "post:SetSysParams")
	//购买nft作品完成通知
	beego.Router("/v2/buyResultInterface", &nftexchangev2.NftExchangeControllerV2{}, "post:BuyResultInterface")

	//查询
	//beego.Router("/v1/getimage", &controllers.NftExchangeController{}, "post:GetImageFromIPFS")
	beego.Router("/v2/ipfsaddtest", &nftexchangev2.NftExchangeControllerV2{}, "get:IpfsTest")

	//查询版本号
	beego.Router("v2/version", &nftexchangev2.NftExchangeControllerV2{}, "get:GetVersion")

	//批量导入
	//beego.Router("/v2/batchNewCollections", &nftexchangev2.NftExchangeControllerV2{}, "post:BatchCreateCollection")
	//beego.Router("/v2/batchUpload", &nftexchangev2.NftExchangeControllerV2{}, "post:BatchUploadNft")
	//beego.Router("/v2/batchBuyResultInterface", &nftexchangev2.NftExchangeControllerV2{}, "post:BatchBuyResultInterface")

}

func registFilters() {
	//路由前进行授权检查
	beego.InsertFilter("/v2/upload", beego.BeforeRouter, nftexchangev2.CheckToken)
	beego.InsertFilter("/v2/buy", beego.BeforeRouter, nftexchangev2.CheckToken)
	beego.InsertFilter("/v2/cancelBuy", beego.BeforeRouter, nftexchangev2.CheckToken)
	beego.InsertFilter("/v2/buyResult", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/queryNFTList", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/queryUserInfo", beego.BeforeRouter, nftexchangev2.CheckToken)
	beego.InsertFilter("/v2/modifyUserInfo", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/queryNFT", beego.BeforeRouter, nftexchangev2.CheckToken)
	beego.InsertFilter("/v2/sell", beego.BeforeRouter, nftexchangev2.CheckToken)
	beego.InsertFilter("/v2/cancelSell", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/queryPendingVrfList", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/vrfNFT", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/queryMarketInfo", beego.BeforeRouter, nftexchangev2.CheckToken)
	beego.InsertFilter("/v2/like", beego.BeforeRouter, nftexchangev2.CheckToken)
	beego.InsertFilter("/v2/newCollections", beego.BeforeRouter, nftexchangev2.CheckToken)
	beego.InsertFilter("/v2/modifyCollections", beego.BeforeRouter, nftexchangev2.CheckToken)
	beego.InsertFilter("/v2/modifyCollectionsImage", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/queryUserNFTList", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/queryUserCollectionList", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/queryUserTradingHistory", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/queryUserOfferList", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/queryUserBidList", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/queryUserFavoriteList", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/queryNFTCollectionList", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/queryCollectionInfo", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/queryMarketTradingHistory", beego.BeforeRouter, nftexchangev2.CheckToken)
	beego.InsertFilter("/v2/userRequireKYC", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/userKYC", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/queryPendingKYCList", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/search", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/askForApprove", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/queryHomePage", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/get_sys_para", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/set_sys_para", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/buyResultInterface", beego.BeforeRouter, nftexchangev2.CheckToken)
	//beego.InsertFilter("/v2/version", beego.BeforeRouter, nftexchangev2.CheckToken)

}