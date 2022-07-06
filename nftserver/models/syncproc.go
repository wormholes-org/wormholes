package models

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/nftexchange/nftserver/ethhelper"
	"github.com/nftexchange/nftserver/ethhelper/database"
	"gorm.io/gorm"
	"log"
	"math/big"
	"strconv"
	"time"
)

//const InfuraWssPoint = "test"
const InfuraWssPoint = "wss://rinkeby.infura.io/ws/v3/97cb2119c79842b7818a7a37df749b2b"
const ReDialDelyTime = 5

func UpdateBlockNumber(sqldsn string) {
	var client *ethclient.Client
	var err error
	for {
		for  {
			client, err = ethclient.Dial(InfuraWssPoint)
			if err != nil {
				log.Println("UpdateBlockNumber() connect err=", err)
				fmt.Println("UpdateBlockNumber() connect err=", err)
				time.Sleep(ReDialDelyTime * time.Second)
			} else {
				log.Println("UpdateBlockNumber() connect OK!")
				fmt.Println("UpdateBlockNumber() connect OK!")
				break
			}
		}
		headers := make(chan *types.Header)
		sub, err := client.SubscribeNewHead(context.Background(), headers)
		if err != nil {
			log.Println("UpdateBlockNumber() connect OK!")
			fmt.Println("UpdateBlockNumber() connect OK!")
			client.Close()
			continue
		}
		nd, err := NewNftDb(sqldsn)
		if err != nil {
			fmt.Printf("UpdateBlockNumber() connect database err = %s\n", err)
			client.Close()
			continue
		}
	loop:
		for {
			select {
			case err := <-sub.Err():
				fmt.Println("UpdateBlockNumber() sub.err()=", err)
				sub.Unsubscribe()
				client.Close()
				log.Println("UpdateBlockNumber() restart.")
				break loop
			case header := <-headers:
				headNumber := header.Number.String()
				fmt.Println("headNumber=", headNumber)
				block, err := client.BlockByHash(context.Background(), header.Hash())
				if err != nil {
					fmt.Println("UpdateBlockNumber() get block err()=", err)
					continue
				}
				fmt.Println("block.number= ", block.Number().String())
				var params SysParams
				dbErr := nd.db.Last(&params)
				if dbErr.Error != nil {
					fmt.Println("UpdateBlockNumber() opendb err=", dbErr.Error)
					continue
				}
				dbErr = nd.db.Model(&SysParams{}).Where("id = ?", params.ID).Update("blocknumber", block.Number().Uint64())
				if dbErr.Error != nil {
					fmt.Println("UpdateBlockNumber() update db err=", dbErr.Error)
				}
			}
		}
		nd.Close()
		fmt.Println("UpdateBlockNumber() db close.")
	}
}

func GetCurrentBlockNumber() uint64{
	var client *ethclient.Client
	var err error
	for  {
		client, err = ethclient.Dial(InfuraWssPoint)
		if err != nil {
			log.Println("GetCurrentBlockNumber() connect err=", err)
			fmt.Println("GetCurrentBlockNumber() connect err=", err)
			time.Sleep(ReDialDelyTime * time.Second)
		} else {
			//log.Println("GetCurrentBlockNumber() connect OK!")
			//fmt.Println("GetCurrentBlockNumber() connect OK!")
			break
		}
	}
	for {
		header, err := client.HeaderByNumber(context.Background(), nil)
		if err != nil {
			fmt.Println("GetCurrentBlockNumber() get HeaderByNumber err=", err)
			continue
		} else {
			fmt.Println("GetCurrentBlockNumber() header.Number=", header.Number.String())
			return header.Number.Uint64()
		}
	}
}

func GetDbBlockNumber(sqldsn string) (uint64, error){
	nd, err := NewNftDb(sqldsn)
	if err != nil {
		fmt.Printf("GetDbBlockNumber() connect database err = %s\n", err)
		return 0, err
	}
	defer nd.Close()
	var params SysParams
	dbErr := nd.db.Last(&params)
	if dbErr.Error != nil {
		fmt.Println("GetDbBlockNumber() opendb err=", dbErr.Error)
		return 0, dbErr.Error
	}
	return params.Scannumber, nil
}

func MintProc(tx *gorm.DB, to, royalty, contractAddr, tokenId, txhash, ts string) error  {
	fmt.Println("MintProc() Start.")
	t, _ := strconv.ParseUint(ts, 10, 64)
	transTime := time.Unix(int64(t), 0)
	var nftRec Nfts
	err := tx.Select("id").Where("contract = ? AND tokenid = ?", contractAddr, tokenId).First(&nftRec)
	if err.Error != nil && err.Error == gorm.ErrRecordNotFound {
		fmt.Println("MintProc() err =", ErrNftNotExist)
		return nil
	} else if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
		return err.Error
	}
	trans := Trans{}
	err = tx.Where("contract = ? AND tokenid = ? AND selltype = ? AND txhash = ?",
		contractAddr, tokenId, SellTypeMintNft.String(), txhash).First(&trans)
	if err.Error == nil {
		fmt.Println("MintProc() err =", ErrTransExist)
		return nil
	} else if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
		return err.Error
	}
	trans = Trans{}
	trans.Contract = contractAddr
	trans.Fromaddr = ""
	trans.Toaddr = to
	trans.Tokenid = tokenId
	//trans.Transtime = time.Now().Unix()
	trans.Transtime = transTime.Unix()
	trans.Selltype = SellTypeMintNft.String()
	//trans.Name = nftRec.Name
	//trans.Meta = nftRec.Meta
	//trans.Desc = nftRec.Desc
	trans.Txhash = txhash
	return tx.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&trans).Create(&trans)
		if err.Error != nil {
			fmt.Println("MintProc() create trans err=", err.Error)
			return err.Error
		}
		nftrecord := Nfts{}
		//nftrecord.Royalty, _ = strconv.Atoi(royalty)
		//nftrecord.Royalty = nftrecord.Royalty / 100
		nftrecord.Mintstate = Minted.String()
		err = tx.Model(&Nfts{}).Where("contract = ? AND tokenid =?",
			contractAddr, tokenId).Updates(&nftrecord)
		if err.Error != nil {
			fmt.Println("MintProc() update nfts record err=", err.Error)
			return err.Error
		}
		fmt.Println("MintProc() Ok")
		return nil
	})
}

func TransProc(tx *gorm.DB, from, to, price, contractAddr, tokenId, txhash, ts string) error {
	t, _ := strconv.ParseUint(ts, 10, 64)
	transTime := time.Unix(int64(t), 0)
	trans := Trans{}
	err := tx.Select("id").Where("contract = ? AND tokenid = ? AND txhash = ? AND (selltype = ? or selltype = ? or selltype = ?)",
		contractAddr, tokenId, txhash, SellTypeFixPrice.String(), SellTypeBidPrice.String(), SellTypeHighestBid.String()).First(&trans)
	if err.Error == nil {
		fmt.Println("TransProc() err =", ErrTransExist)
		return nil
	}
	if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
		fmt.Println("TransProc() err =", err.Error)
		return err.Error
	}

	var nftRec Nfts
	err = tx.Select("id").Where("contract = ? AND tokenid = ?", contractAddr, tokenId).First(&nftRec)
	if err.Error != nil && err.Error == gorm.ErrRecordNotFound {
		fmt.Println("TransProc() err =", ErrNftNotExist)
		return nil
	} else if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
		return err.Error
	}
	err = tx.Where("contract = ? AND tokenid = ?", contractAddr, tokenId).First(&nftRec)
	if err.Error != nil && err.Error == gorm.ErrRecordNotFound {
		fmt.Println("TransProc() err =", ErrNftNotExist)
		return nil
	} else if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
		return err.Error
	}
	var auctionRec Auction
	err = tx.Where("contract = ? AND tokenid = ? AND ownaddr =?",
		contractAddr, tokenId, nftRec.Ownaddr).First(&auctionRec)
	if err.Error != nil {
		if err.Error != gorm.ErrRecordNotFound {
			fmt.Println("TransProc() dbase err=", err.Error)
			return err.Error
		}
		trans = Trans{}
		//trans.Auctionid = auctionRec.ID
		trans.Contract = contractAddr
		trans.Createaddr = nftRec.Createaddr
		trans.Fromaddr = from
		trans.Toaddr = to
		trans.Nftid = nftRec.ID
		trans.Tokenid = tokenId
		//trans.Paychan = auctionRec.Paychan
		//trans.Currency = auctionRec.Currency
		trans.Txhash = txhash
		trans.Name = nftRec.Name
		trans.Meta = nftRec.Meta
		trans.Desc = nftRec.Desc
		trans.Price, _ = strconv.ParseUint(price, 10, 64)
		trans.Transtime = transTime.Unix()
		//trans.Transtime = time.Now().Unix()
		trans.Selltype = SellTypeFixPrice.String()
		return tx.Transaction(func(tx *gorm.DB) error {
			err = tx.Model(&trans).Create(&trans)
			if err.Error != nil {
				fmt.Println("TransProc() create trans record err=", err.Error)
				return err.Error
			}
			var collectRec Collects
			err = tx.Where("name = ? AND createaddr =?",
				nftRec.Collections, nftRec.Collectcreator).First(&collectRec)
			if err.Error == nil {
				transCnt := collectRec.Transcnt + 1
				transAmt := collectRec.Transamt + trans.Price
				collectRec = Collects{}
				collectRec.Transcnt = transCnt
				collectRec.Transamt = transAmt
				err = tx.Model(&Collects{}).Where("name = ? AND createaddr =?",
					nftRec.Collections, nftRec.Collectcreator).Updates(&collectRec)
				if err.Error != nil {
					fmt.Println("TransProc() update collectRec err=", err.Error)
					return err.Error
				}
			}
			fmt.Println("TransProc() OK" )
			return nil
		})
	} else {
		trans = Trans{}
		trans.Auctionid = auctionRec.ID
		trans.Contract = auctionRec.Contract
		trans.Createaddr = nftRec.Createaddr
		trans.Fromaddr = from
		trans.Toaddr = to
		trans.Nftid = auctionRec.Nftid
		trans.Tokenid = auctionRec.Tokenid
		trans.Paychan = auctionRec.Paychan
		trans.Currency = auctionRec.Currency
		trans.Txhash = txhash
		trans.Name = nftRec.Name
		trans.Meta = nftRec.Meta
		trans.Desc = nftRec.Desc
		trans.Price, _ = strconv.ParseUint(price, 10, 64)
		//trans.Price = auctionRec.Price
		trans.Transtime = transTime.Unix()
		//trans.Transtime = time.Now().Unix()
		trans.Selltype = auctionRec.Selltype
		return tx.Transaction(func(tx *gorm.DB) error {
			err = tx.Model(&trans).Create(&trans)
			if err.Error != nil {
				fmt.Println("TransProc() create trans record err=", err.Error)
				return err.Error
			}
			var collectRec Collects
			err = tx.Where("name = ? AND createaddr =?",
				nftRec.Collections, nftRec.Collectcreator).First(&collectRec)
			if err.Error == nil {
				transCnt := collectRec.Transcnt + 1
				transAmt := collectRec.Transamt + trans.Price
				collectRec = Collects{}
				collectRec.Transcnt = transCnt
				collectRec.Transamt = transAmt
				err = tx.Model(&Collects{}).Where("name = ? AND createaddr =?",
					nftRec.Collections, nftRec.Collectcreator).Updates(&collectRec)
				if err.Error != nil {
					fmt.Println("TransProc() update collectRec err=", err.Error)
					return err.Error
				}
			}
			nftrecord := Nfts{}
			nftrecord.Ownaddr = to
			nftrecord.Selltype = SellTypeNotSale.String()
			nftrecord.Paychan = auctionRec.Paychan
			nftrecord.TransCur = auctionRec.Currency
			nftrecord.Transprice = trans.Price
			nftrecord.Transamt = nftRec.Transamt + trans.Price
			nftrecord.Transcnt = nftRec.Transcnt + 1
			nftrecord.Transtime = time.Now().Unix()
			err = tx.Model(&Nfts{}).Where("contract = ? AND tokenid =?",
				auctionRec.Contract, auctionRec.Tokenid).Updates(&nftrecord)
			if err.Error != nil {
				fmt.Println("TransProc() update record err=", err.Error)
				return err.Error
			}
			err = tx.Model(&Auction{}).Where("contract = ? AND tokenid = ?",
				auctionRec.Contract, auctionRec.Tokenid).Delete(&Auction{})
			if err.Error != nil {
				fmt.Println("TransProc() delete auction record err=", err.Error)
				return err.Error
			}
			err = tx.Model(&Bidding{}).Where("contract = ? AND tokenid = ?",
				auctionRec.Contract, auctionRec.Tokenid).Delete(&Bidding{})
			if err.Error != nil {
				fmt.Println("TransProc() delete bid record err=", err.Error)
				return err.Error
			}
			fmt.Println("TransProc() OK" )
			return nil
		})
	}
}

func SyncBlockTxs(sqldsn string, block uint64, blockTxs []*database.NftTx) error {
	nd, err := NewNftDb(sqldsn)
	if err != nil {
		fmt.Printf("SyncBlockTxs() connect database err = %s\n", err)
		return err
	}
	defer nd.Close()
	for _, nftTx := range blockTxs {
		if nftTx.From == "" {
			err = MintProc(nd.db, nftTx.To, "", nftTx.Contract, nftTx.TokenId, nftTx.TxHash, nftTx.Ts)
			if err != nil {
				break
			}
		}
		if nftTx.From != "" && nftTx.To != "" && nftTx.Value != "" &&
			nftTx.From != ZeroAddr && nftTx.To != ZeroAddr {
			fmt.Println("SyncBlockTxs() nftTx.Value=", nftTx.Value)
			var price string
			if len(nftTx.Value) >= 9 {
				price = nftTx.Value[:len(nftTx.Value)-9]
			} else {
				continue
				//price = "0"
			}
			fmt.Println("SyncBlockTxs() price=", price)
			err = TransProc(nd.db, nftTx.From, nftTx.To, price, nftTx.Contract, nftTx.TokenId, nftTx.TxHash, nftTx.Ts)
			if err != nil {
				break
			}
		}
	}
	if err == nil {
		var params SysParams
		dbErr := nd.db.Last(&params)
		if dbErr.Error != nil {
			fmt.Println("SyncBlockTxs() params err=", dbErr.Error)
			return dbErr.Error
		}
		dbErr = nd.db.Model(&SysParams{}).Where("id = ?", params.ID).Update("scannumber", block + 1)
		if dbErr.Error != nil {
			fmt.Println("SyncBlockTxs() update params err=", dbErr.Error)
			return dbErr.Error
		}
		fmt.Println("SyncBlockTxs() update block=", block)
	}
	return err
	/*return nd.db.Transaction(func(tx *gorm.DB) error {
		var err error
		for _, nftTx := range blockTxs {
			if nftTx.From == "" {
 				err = MintProc(tx, nftTx.To, "", nftTx.Contract, nftTx.TokenId, nftTx.TxHash, nftTx.Ts)
				if err != nil {
					break
				}
			}
			if nftTx.From != "" && nftTx.To != "" && nftTx.Value != "" &&
				nftTx.From != ZeroAddr && nftTx.To != ZeroAddr {
				fmt.Println("SyncBlockTxs() nftTx.Value=", nftTx.Value)
				var price string
				if len(nftTx.Value) >= 9 {
					price = nftTx.Value[:len(nftTx.Value)-9]
				} else {
					continue
					//price = "0"
				}
				fmt.Println("SyncBlockTxs() price=", price)
				err = TransProc(tx, nftTx.From, nftTx.To, price, nftTx.Contract, nftTx.TokenId, nftTx.TxHash, nftTx.Ts)
				if err != nil {
					break
				}
			}
		}
		if err == nil {
			var params SysParams
			dbErr := tx.Last(&params)
			if dbErr.Error != nil {
				fmt.Println("SyncBlockTxs() params err=", dbErr.Error)
				return dbErr.Error
			}
			dbErr = nd.db.Model(&SysParams{}).Where("id = ?", params.ID).Update("scannumber", block + 1)
			if dbErr.Error != nil {
				fmt.Println("SyncBlockTxs() update params err=", dbErr.Error)
				return dbErr.Error
			}
			fmt.Println("SyncBlockTxs() update block=", block)
		}
		return err
	})*/
}

func AmountValid(Price uint64, addr string) (bool, string, error) {
	auth, err := ethhelper.AllowanceOfWeth(addr)
	if err != nil {
		return false, "", err
	}
	balance, err := ethhelper.BalanceOfWeth(addr)
	if err != nil {
		return false, "", err
	}
	temp := strconv.FormatUint(balance, 10)
	if len(temp) < 9 {
		return false, ErrBalanceLess.Error(), nil
	} else {
		temp = temp[:len(temp)-9]
	}
	balance, _ = strconv.ParseUint(temp, 10, 64)
	if balance < Price {
		return false, ErrBalanceLess.Error(), nil
	}
	if auth == "0" || len(auth) <= 9 {
		return false, ErrAuthorizeLess.Error(), nil
	}
	auth = auth[:len(auth)-9]
	authAmount := new(big.Int)
	authAmount.SetString(auth, 10)
	fmt.Println("authAmount=", authAmount)
	pricet := new(big.Int)
	pricet = pricet.SetUint64(Price)
	fmt.Println("pricet=", pricet)
	if authAmount.Cmp(pricet) < 0 {
		return false, ErrAuthorizeLess.Error(), nil
	}
	return true, "",  nil
}

func ScanBiddings(sqldsn string, scanAddr map[string]bool) error {
	nd, err := NewNftDb(sqldsn)
	if err != nil {
		fmt.Printf("SyncBlockTxs() connect database err = %s\n", err)
		return err
	}
	defer nd.Close()
	var biddings []Bidding
	for s, _ := range scanAddr {
		dberr := nd.db.Model(Bidding{}).Where("Bidaddr = ?", s).Find(&biddings)
		if dberr.Error == nil && dberr.RowsAffected != 0 {
			for _, bidding := range biddings {
				valid, _, cerr := AmountValid(bidding.Price, bidding.Bidaddr)
				if cerr != nil {
					continue
				}
				if !valid {
					dberr = nd.db.Model(&Bidding{}).Where("id = ?", bidding.ID).Delete(&Bidding{})
				}
			}
		}
	}
	return err
}

func GetBlockTxs(block uint64) ([]*database.NftTx, []*ethhelper.WethTransfer, []*ethhelper.WethTransfer) {
	buyResultCh := make(chan []*database.NftTx)
	wethTransferCh := make(chan *ethhelper.WethTransfer)
	wethApproveCh := make(chan *ethhelper.WethTransfer)
	var BlockTxs []*database.NftTx
	var wethTransfers []*ethhelper.WethTransfer
	var wethApproves []*ethhelper.WethTransfer
	endCh := make(chan bool)
	go ethhelper.SyncNftFromChain(strconv.FormatUint((block), 10), true, buyResultCh, wethTransferCh, wethApproveCh, endCh)
loop:
	for {
		select {
			case buyResult := <-buyResultCh:
				BlockTxs = append(BlockTxs, buyResult...)
			case wethTransfer := <-wethTransferCh:
				wethTransfers = append(wethTransfers, wethTransfer)
			case wethApprove := <-wethApproveCh:
				wethApproves = append(wethApproves, wethApprove)
			case <-endCh:
				break loop
			default:
		}
	}
	fmt.Println("GetBlockTx() end blocktxs count=", len(BlockTxs))
	fmt.Println("GetBlockTx() end wethTransfers count=", len(wethTransfers))
	fmt.Println("GetBlockTx() end wethApproves count=", len(wethApproves))
	return BlockTxs, wethTransfers, wethApproves
}


func SyncProc(sqldsn string, syncCh chan uint64) chan struct{}{
	var procEnd = make(chan struct{}, 1)
	go func() {
		for {
			select {
				case blockE := <- syncCh:
					fmt.Println("SyncProc() start end blockE=", blockE)
					blockS, err := GetDbBlockNumber(sqldsn)
					if err != nil {
						fmt.Println("SyncProc() GetDbBlockNumber() err=", err)
						procEnd <- struct{}{}
						continue
					}
					mAddr := make(map[string]bool)
				loop:
					for blockS <= blockE {
						fmt.Println("SyncProc() call SyncNftFromChain() blockNum=", blockS)
						txs, wethts, _ := GetBlockTxs(blockS)
						fmt.Println("SyncProc() call SyncNftFromChain() OK")
						fmt.Println(time.Now().String()[:22],"SyncProc() call SyncBlockTxs() Start blockNum=", blockS)
						err := SyncBlockTxs(sqldsn, blockS, txs)
						if err != nil {
							fmt.Println("SyncProc() call SyncBlockTxs() err=", err)
							break loop
						}
						fmt.Println(time.Now().String()[:22],"SyncProc() call SyncBlockTxs() End blockNum=", blockS)
						for _, wetht := range wethts {
							mAddr[wetht.From] = true
						}
						blockS ++
					}
					if len(mAddr) != 0 {
						ScanBiddings(sqldsn, mAddr)
					}
					fmt.Println("SyncProc() sync blockE=", blockE)
					procEnd <- struct{}{}
			}
		}
	}()
	return procEnd
}

func InitSyncBlockTs(sqldsn string) error {
	blockS, err := GetDbBlockNumber(sqldsn)
	if err != nil {
		fmt.Println("SyncProc() get scan block num err=", err)
		return err
	}
	fmt.Println("SyncProc() end blockNum=", blockS)
	mAddr := make(map[string]bool)
	for blockS <= GetCurrentBlockNumber() {
		fmt.Println("SyncProc() call SyncNftFromChain() blockNum=", blockS)
		txs, wethts, _ := GetBlockTxs(blockS)
		fmt.Println("SyncProc() call SyncNftFromChain() OK")
		err := SyncBlockTxs(sqldsn, blockS, txs)
		if err != nil {
			return err
		}
		for _, wetht := range wethts {
			mAddr[wetht.From] = true
		}
		blockS ++
	}
	if len(mAddr) != 0 {
		ScanBiddings(sqldsn, mAddr)
	}
	return nil
}
