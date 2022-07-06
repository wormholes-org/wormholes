package contracts

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/nftexchange/nftserver/common/contracts/admin"
	"github.com/nftexchange/nftserver/models"
	"log"
	"math/big"
)

func init() {
	//TimeProc()
}

func AdminList() ([]string, error) {
	client, err := ethclient.Dial(models.EthersNode)
	if err != nil {
		log.Println(err)
	}
	privateKey, err := crypto.HexToECDSA(AdminListPrv)
	if err != nil {
		log.Println(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Println("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	address := common.HexToAddress(models.AdminAddr)
	instance, err := admin.NewAdmin(address, client)
	if err != nil {
		log.Println(err)
	}
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Println(err)
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println(err)
	}
	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice
	result, err := instance.List(&bind.CallOpts{From: fromAddress, Context: context.Background()})
	if err != nil {
		log.Println(err)
	}
	var addrs []string
	for _, addr := range result {
		addrs = append(addrs, addr.Hex())
	}
	return addrs, err
}

func AdminListWithClient(client *ethclient.Client) ([]string, error) {
	/*client, err := ethclient.Dial(InfuraPoint)
	if err != nil {
		log.Println(err)
	}*/
	privateKey, err := crypto.HexToECDSA(AdminListPrv)
	if err != nil {
		log.Println(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Println("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	address := common.HexToAddress(models.AdminAddr)
	instance, err := admin.NewAdmin(address, client)
	if err != nil {
		log.Println(err)
	}
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Println(err)
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println(err)
	}
	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice
	result, err := instance.List(&bind.CallOpts{From: fromAddress, Context: context.Background()})
	if err != nil {
		log.Println(err)
	}
	var addrs []string
	for _, addr := range result {
		addrs = append(addrs, addr.Hex())
	}
	return addrs, err
}
