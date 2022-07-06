package nftexchangev2

func (nft *NftExchangeControllerV2) GetVersion() {
	var version string
	version= `
			<h1>0.5.17</h1>
			<p>login 添加approve_addr 字段, 增加答题机制，详情见说明。					[Y]</p>
			<h1>0.5.16.1</h1>
			<p>fix bugs for querynftlist()											[Y]</p>
			<p>remove validation for search()										[Y]</p>
			<h1>0.5.16</h1>
			<p>[new] buyResultInterface												[Y]</p>
			<h1>0.5.15</h1>
			<p>[update] queryHomePage 添加字段total									[Y]</p>
			<p>[update] modifyCollectionsImage 输入接口更新							[Y]</p>
			<h1>0.5.14</h1>
			<p>[update] queryCollectionInfo 添加3个交易字段 
			<p>- trade_amount - trade_avg_price - trade_floor_price					[Y]</p>
			<h1>0.5.13</h1>
			<p>[update] queryUserBidlist 添加name字段 NFT名称							[Y]</p>
			<p>[update] queryUserOfferlist 添加name字段 NFT名称						[Y]</p>
			<h1>0.5.12</h1>
			<p>[update] 所有地址全转为小写写法(用户地址，合约地址)							[Y]</p>
			<h1>0.5.11</h1>
			<p>[update] queryNFT 添加字段collection_desc								[Y]</p>
			<h1>0.5.10</h1>
			<p>[new] version 添加version 接口										[Y]</p>
			<p>[new] set_sys_para 后台接口, 前端无视									[Y]</p>
			<p>[new] get_sys_para 后台接口, 前端无视									[Y]</p>
			<h1>0.5.9.1</h1>
			<p>[update] queryUserNFTList 增加字段collection_creator_addr				[Y]</p>
			<p>[update] queryUserCollectionList 增加字段collection_creator_addr		[Y]</p>
			<p>[update] queryUserFavoriteList 增加字段collection_creator_addr		[Y]</p>
			<p>[update] queryNFTCollectionList 增加字段collection_creator_addr		[Y]</p>
			<p>[update] queryCollectionInfo 增加字段collection_creator_addr			[Y]</p>
			<p>[update] queryHomePage 增加字段collection_creator_addr				[Y]</p>
			<p>[update] queryNFT 增加字段trade_hash									[Y]</p>
			<p>[update] queryUserTradingHistroy 增加字段trade_hash					[Y]</p>
			<p>[update] queryMarketTradingHistroy 增加字段trade_hash					[Y]</p>
			`

	nft.Ctx.ResponseWriter.Write([]byte(version))
}