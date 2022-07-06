const ethers = require("ethers");

(async () => {
    try{
        //contract,owner,tokenId,amount,royalty,metaUrl
        //owner 发行NFT签名 //nft合约地址 接收者（創建者） tokenid 数量  版税 meta
        let  args = process.argv.splice(2);
        const privateKey = "0x8c995fd78bddf528bd548cce025f62d4c3c0658362dbfd31b23414cf7ce2e8ed";
        const mintData = ethers.utils.solidityKeccak256(['address', 'address', 'uint256', 'uint256', 'uint16', 'string'],
            [args[0],args[1], parseInt(args[2]),  parseInt(args[3]),  parseInt(args[4]),  args[5]])

        const sig3 = await (new ethers.Wallet(privateKey)).signMessage(ethers.utils.arrayify(mintData));
        console.log(ethers.utils.joinSignature(sig3));
    }catch (e) {
        console.log(e)
    }

})()