const Web3 = require('web3');
const contracts = require("./NFT1155.json")
const web3 = new Web3(new Web3.providers.HttpProvider("https://rinkeby.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161"));
const privKey = '564ea566096d3de340fc5ddac98aef672f916624c8b0e4664a908cd2a2d156fe'
const from = '0x077d34394Ed01b3f31fBd9816cF35d4558146066'
const to = contracts.address
const networkId=4

async function payload() {
    try {
        const contract = new web3.eth.Contract(contracts.abi)
        const input = contract.methods.mint(from,147258,"https://api.coolcatsnft.com/cat/9929").encodeABI()
        let nonce = 0;
        await web3.eth.getTransactionCount(from, (err, count) => {
            nonce = count
        })
        let gasPrice = ''
        await web3.eth.getGasPrice((err, gas) => {
            gasPrice = '0x' + parseInt(gas).toString(16)
        })
        const ethTx = require('ethereumjs-tx').Transaction
        const privateKey = Buffer.from(
            privKey,
            'hex',
        )

        const Common = require('ethereumjs-common').default
        const customCommon = Common.forCustomChain(
            'rinkeby',
            {
                name: 'my-network',
                networkId: networkId,
                chainId: networkId,
            },
            'petersburg',
        )
        const txParams = {
            nonce: nonce,
            gasPrice: gasPrice,
            gasLimit: '0x2dc6c0',
            gas: '0x1e8480',
            from: from,
            to: to,
            value: '0x00',
            data: input,
            chainId: networkId
        }
        const tx = new ethTx(txParams, { common: customCommon })

        tx.sign(privateKey)
        const serializedTx = tx.serialize()
        let ret = await web3.eth.sendSignedTransaction('0x' + serializedTx.toString('hex'), (error, hash) => {
        })
        console.log(ret.logs[0].topics)
    }catch (e) {
        console.log(e)
    }

}
payload().then()