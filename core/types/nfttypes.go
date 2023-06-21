package types

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
	"regexp"
	"strings"
)

type MintDeep struct {
	UserMint     *big.Int
	OfficialMint *big.Int
	//ExchangeList SNFTExchangeList
}

//type SNFTExchange struct {
//	InjectedInfo
//	NFTAddress         common.Address
//	MergeLevel         uint8
//	CurrentMintAddress common.Address
//	BlockNumber        *big.Int
//}
//type InjectedInfo struct {
//	MetalUrl string
//	Royalty  uint32
//	Creator  string
//}
//
//type SNFTExchangeList struct {
//	SNFTExchanges []*SNFTExchange
//}
//
//func (ex *SNFTExchange) MinNFTAddress() common.Address {
//	return ex.NFTAddress
//}
//func (ex *SNFTExchange) MaxNFTAddress() common.Address {
//	if ex.MergeLevel == 0 {
//		return ex.NFTAddress
//	}
//	minAddrInt := big.NewInt(0)
//	minAddrInt.SetBytes(ex.NFTAddress.Bytes())
//	nftNumber := math.BigPow(256, int64(ex.MergeLevel))
//	maxAddrInt := big.NewInt(0)
//	maxAddrInt.Add(minAddrInt, new(big.Int).Sub(nftNumber, big.NewInt(1)))
//	maxAddr := common.BytesToAddress(maxAddrInt.Bytes())
//	return maxAddr
//}
//func (ex *SNFTExchange) MaxNFTAddress16() common.Address {
//	if ex.MergeLevel == 0 {
//		return ex.NFTAddress
//	}
//	minAddrInt := big.NewInt(0)
//	minAddrInt.SetBytes(ex.NFTAddress.Bytes())
//	nftNumber := math.BigPow(16, int64(ex.MergeLevel))
//	maxAddrInt := big.NewInt(0)
//	maxAddrInt.Add(minAddrInt, new(big.Int).Sub(nftNumber, big.NewInt(1)))
//	maxAddr := common.BytesToAddress(maxAddrInt.Bytes())
//	return maxAddr
//}
//
//func (list *SNFTExchangeList) PopAddress(blocknumber *big.Int) (common.Address, *InjectedInfo, bool) {
//	if len(list.SNFTExchanges) > 0 {
//		log.Info("PopAddress()", "SNFTExchanges[0].BlockNumber=", list.SNFTExchanges[0].BlockNumber.Uint64())
//		log.Info("PopAddress()", "-----------------blocknumber=", blocknumber.Uint64())
//		if list.SNFTExchanges[0].BlockNumber.Cmp(blocknumber) >= 0 {
//			return common.Address{}, nil, false
//		}
//		addr := list.SNFTExchanges[0].CurrentMintAddress
//		InjectedInfo := &InjectedInfo{
//			MetalUrl: list.SNFTExchanges[0].MetalUrl,
//			Royalty:  list.SNFTExchanges[0].Royalty,
//			Creator:  list.SNFTExchanges[0].Creator,
//		}
//		//if list.SNFTExchanges[0].CurrentMintAddress == list.SNFTExchanges[0].MaxNFTAddress() {
//		if list.SNFTExchanges[0].CurrentMintAddress == list.SNFTExchanges[0].MaxNFTAddress16() {
//			if len(list.SNFTExchanges) > 1 {
//				list.SNFTExchanges = list.SNFTExchanges[1:]
//			} else {
//				list.SNFTExchanges = list.SNFTExchanges[:0]
//			}
//		} else {
//			currentMintInt := new(big.Int).SetBytes(list.SNFTExchanges[0].CurrentMintAddress.Bytes())
//			currentMintInt.Add(currentMintInt, big.NewInt(1))
//			list.SNFTExchanges[0].CurrentMintAddress = common.BytesToAddress(currentMintInt.Bytes())
//		}
//		return addr, InjectedInfo, true
//	}
//	return common.Address{}, nil, false
//}

type PledgedToken struct {
	Address      common.Address
	Amount       *big.Int
	Flag         bool
	ProxyAddress common.Address
}

type InjectedOfficialNFT struct {
	Dir        string         `json:"dir"`
	StartIndex *big.Int       `json:"start_index"`
	Number     uint64         `json:"number"`
	Royalty    uint16         `json:"royalty"`
	Creator    string         `json:"creator"`
	Address    common.Address `json:"address"`
	VoteWeight *big.Int       `json:"vote_weight"`
}

type InjectedOfficialNFTList struct {
	InjectedOfficialNFTs []*InjectedOfficialNFT
}

func (list *InjectedOfficialNFTList) GetInjectedInfo(addr common.Address) *InjectedOfficialNFT {
	maskB, _ := big.NewInt(0).SetString("8000000000000000000000000000000000000000", 16)
	addrInt := new(big.Int).SetBytes(addr.Bytes())
	addrInt.Sub(addrInt, maskB)
	tempInt := new(big.Int)
	for _, injectOfficialNFT := range list.InjectedOfficialNFTs {
		if injectOfficialNFT.StartIndex.Cmp(addrInt) == 0 {
			return injectOfficialNFT
		}
		if injectOfficialNFT.StartIndex.Cmp(addrInt) < 0 {
			tempInt.SetInt64(0)
			tempInt.Add(injectOfficialNFT.StartIndex, new(big.Int).SetUint64(injectOfficialNFT.Number))
			if tempInt.Cmp(addrInt) > 0 {
				return injectOfficialNFT
			}
		}
	}

	return nil
}

func (list *InjectedOfficialNFTList) DeleteExpireElem(num *big.Int) {
	var index int
	maskB, _ := big.NewInt(0).SetString("8000000000000000000000000000000000000000", 16)
	for k, injectOfficialNFT := range list.InjectedOfficialNFTs {
		sum := new(big.Int).Add(injectOfficialNFT.StartIndex, new(big.Int).SetUint64(injectOfficialNFT.Number))
		sum.Add(sum, maskB)
		if sum.Cmp(num) > 0 {
			index = k
			break
		}
	}

	list.InjectedOfficialNFTs = list.InjectedOfficialNFTs[index:]
}

func (list *InjectedOfficialNFTList) RemainderNum(addrInt *big.Int) uint64 {
	var sum uint64
	maskB, _ := big.NewInt(0).SetString("8000000000000000000000000000000000000000", 16)
	tempInt := new(big.Int)
	for _, injectOfficialNFT := range list.InjectedOfficialNFTs {
		if injectOfficialNFT.StartIndex.Cmp(addrInt) >= 0 {
			sum = sum + injectOfficialNFT.Number
		}
		if injectOfficialNFT.StartIndex.Cmp(addrInt) < 0 {
			tempInt.SetInt64(0)
			tempInt.Add(injectOfficialNFT.StartIndex, new(big.Int).SetUint64(injectOfficialNFT.Number))
			tempInt.Add(tempInt, maskB)
			if tempInt.Cmp(addrInt) >= 0 {
				sum = sum + new(big.Int).Sub(tempInt, addrInt).Uint64()
			}
		}
	}

	return sum
}

func (list *InjectedOfficialNFTList) MaxIndex() *big.Int {
	max := big.NewInt(0)
	for _, injectOfficialNFT := range list.InjectedOfficialNFTs {
		index := new(big.Int).Add(injectOfficialNFT.StartIndex, new(big.Int).SetUint64(injectOfficialNFT.Number))
		if index.Cmp(max) > 0 {
			max.Set(index)
		}
	}

	return max
}

func (list *InjectedOfficialNFTList) DeepCopy() *InjectedOfficialNFTList {
	tempList := &InjectedOfficialNFTList{
		InjectedOfficialNFTs: make([]*InjectedOfficialNFT, 0, len(list.InjectedOfficialNFTs)),
	}

	for _, v := range list.InjectedOfficialNFTs {
		tempInjected := &InjectedOfficialNFT{
			Dir:        v.Dir,
			StartIndex: new(big.Int).Set(v.StartIndex),
			Number:     v.Number,
			Royalty:    v.Royalty,
			Creator:    v.Creator,
			Address:    v.Address,
			VoteWeight: new(big.Int).Set(v.VoteWeight),
		}

		tempList.InjectedOfficialNFTs = append(tempList.InjectedOfficialNFTs, tempInjected)
	}

	return tempList
}

// Wormholes struct for handling NFT transactions
type Wormholes struct {
	Type         uint8  `json:"type"`
	NFTAddress   string `json:"nft_address,omitempty"`
	ProxyAddress string `json:"proxy_address,omitempty"`
	ProxySign    string `json:"proxy_sign,omitempty"`
	Exchanger    string `json:"exchanger,omitempty"`
	Royalty      uint16 `json:"royalty,omitempty"`
	MetaURL      string `json:"meta_url,omitempty"`
	//ApproveAddress string		`json:"approve_address"`
	FeeRate       uint16           `json:"fee_rate,omitempty"`
	Name          string           `json:"name,omitempty"`
	Url           string           `json:"url,omitempty"`
	Dir           string           `json:"dir,omitempty"`
	StartIndex    string           `json:"start_index,omitempty"`
	Number        uint64           `json:"number,omitempty"`
	Buyer         Payload          `json:"buyer,omitempty"`
	Seller1       Payload          `json:"seller1,omitempty"`
	Seller2       MintSellPayload  `json:"seller2,omitempty"`
	ExchangerAuth ExchangerPayload `json:"exchanger_auth,omitempty"`
	Creator       string           `json:"creator,omitempty"`
	Version       string           `json:"version,omitempty"`
	RewardFlag    uint8            `json:"reward_flag,omitempty"`
	BuyerAuth     TraderPayload    `json:"buyer_auth,omitempty"`
	SellerAuth    TraderPayload    `json:"seller_auth,omitempty"`
	NoAutoMerge   bool             `json:"no_automerge,omitempty"`
}

const WormholesVersion = "v0.0.1"
const PattenAddr = "^0x[0-9a-fA-F]{40}$"

// var PattenAddr = "^0[xX][0-9a-fA-F]{40}$"
// var PattenHex = "^[0-9a-fA-F]+$"
func (w *Wormholes) CheckFormat() error {
	//regHex, _ := regexp.Compile(PattenHex)
	//regAddr, _ := regexp.Compile(PattenAddr)

	switch w.Type {
	case 0:
		if len(w.MetaURL) > 256 {
			return errors.New("metaurl too long")
		}
		//exchangerMatch := regAddr.MatchString(w.Exchanger)
		//if !exchangerMatch {
		//	return errors.New("exchanger format error")
		//}

	case 1:

	case 2:
	case 3:
	case 4:
	case 5:
	case 6:
	//case 7:
	//case 8:
	case 9:
	case 10:
	case 11:

		if len(w.Name) > 64 {
			w.Name = string([]byte(w.Name)[:64])
		}
		if len(w.Url) > 128 {
			w.Url = string([]byte(w.Url)[:128])
		}

	case 12:
	case 13:
	case 14:
	case 15:
	case 16:
		if len(w.Seller2.MetaURL) > 256 {
			return errors.New("metaurl too long")
		}

	case 17:
		if len(w.Seller2.MetaURL) > 256 {
			return errors.New("metaurl too long")
		}

	case 18:
	case 19:
		if len(w.Seller2.MetaURL) > 256 {
			return errors.New("metaurl too long")
		}

	case 20:
	case 21:
	case 22:
	case 23:
		if len(w.Dir) > 256 {
			return errors.New("dir too long")
		}

		if len(w.Creator) > 0 {
			regAddr, err := regexp.Compile(PattenAddr)
			if err != nil {
				return err
			}
			match := regAddr.MatchString(w.Creator)
			if !match {
				return errors.New("invalid creator")
			}
		}

	case 24:
		if len(w.Dir) > 256 {
			return errors.New("dir too long")
		}

		if len(w.Creator) > 0 {
			regAddr, err := regexp.Compile(PattenAddr)
			if err != nil {
				return err
			}
			match := regAddr.MatchString(w.Creator)
			if !match {
				return errors.New("invalid creator")
			}
		}

	case 25:
		recipient := strings.ToLower(w.ProxyAddress)
		regAddr, err := regexp.Compile(PattenAddr)
		if err != nil {
			return err
		}
		match := regAddr.MatchString(recipient)
		if !match {
			return errors.New("invalid proxy address")
		}
	case 26:
	case 27:
	case 28:
	case 29:
	case 30:
	case 31:
	default:
		return errors.New("not exist nft type")
	}

	return nil
}

func (w *Wormholes) TxGas() (uint64, error) {

	switch w.Type {
	case 0:
		return params.WormholesTx0, nil
	case 1:
		return params.WormholesTx1, nil
	case 2:
		return params.WormholesTx2, nil
	case 3:
		return params.WormholesTx3, nil
	case 4:
		return params.WormholesTx4, nil
	case 5:
		return params.WormholesTx5, nil
	case 6:
		return params.WormholesTx6, nil
	//case 7:
	//	return params.WormholesTx7, nil
	//case 8:
	//	return params.WormholesTx8, nil
	case 9:
		return params.WormholesTx9, nil
	case 10:
		return params.WormholesTx10, nil
	case 11:
		return params.WormholesTx11, nil

	case 12:
		return params.WormholesTx12, nil
	//case 13:
	//	return params.WormholesTx13, nil
	case 14:
		return params.WormholesTx14, nil
	case 15:
		return params.WormholesTx15, nil
	case 16:
		return params.WormholesTx16, nil
	case 17:
		return params.WormholesTx17, nil
	case 18:
		return params.WormholesTx18, nil
	case 19:
		return params.WormholesTx19, nil
	case 20:
		return params.WormholesTx20, nil
	case 21:
		return params.WormholesTx21, nil
	case 22:
		return params.WormholesTx22, nil
	case 23:
		return params.WormholesTx23, nil
	case 24:
		return params.WormholesTx24, nil
	case 25:
		return params.WormholesTx25, nil
	case 26:
		return params.WormholesTx26, nil
	case 27:
		return params.WormholesTx27, nil
	case 28:
		return params.WormholesTx28, nil
	case 29:
		return params.WormholesTx29, nil
	case 30:
		return params.WormholesTx30, nil
	case 31:
		return params.WormholesTx31, nil
	default:
		return 0, errors.New("not exist nft type")
	}
}

type Payload struct {
	Amount      string `json:"price"`
	NFTAddress  string `json:"nft_address"`
	Exchanger   string `json:"exchanger"`
	BlockNumber string `json:"block_number"`
	Seller      string `json:"seller"`
	Sig         string `json:"sig"`
}

type MintSellPayload struct {
	Amount        string `json:"price"`
	Royalty       string `json:"royalty"`
	MetaURL       string `json:"meta_url"`
	ExclusiveFlag string `json:"exclusive_flag"`
	Exchanger     string `json:"exchanger"`
	BlockNumber   string `json:"block_number"`
	Sig           string `json:"sig"`
}

type ExchangerPayload struct {
	ExchangerOwner string `json:"exchanger_owner"`
	To             string `json:"to"`
	BlockNumber    string `json:"block_number"`
	Sig            string `json:"sig"`
}

type TraderPayload struct {
	Exchanger   string `json:"exchanger"`
	BlockNumber string `json:"block_number"`
	Sig         string `json:"sig"`
}

// *** modify to support nft transaction 20211215 end ***

type NominatedOfficialNFT struct {
	InjectedOfficialNFT
}
