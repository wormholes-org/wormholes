package types

import (
	"bytes"
	"encoding/hex"

	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

func TestCalculateAddressRange(t *testing.T) {
	rand.Seed(time.Now().Unix())
	var validators []*Validator

	for i := 0; i < 10; i++ {
		validators = append(validators, NewValidator(RandomAddr(), big.NewInt(rand.Int63()), common.Address{}))
	}
	validatorList := NewValidatorList(validators)
	for _, vl := range validatorList.Validators {
		validatorList.CalculateAddressRange(vl.Addr, validatorList.StakeBalance(vl.Addr))
	}

	//---------------------------------------------//
	//for _, v := range validatorList.Validators {
	//	fmt.Println("address---", v.Addr, "|---weight", v.Weight)
	//}
}

func TestCalculateAddressRange2(t *testing.T) {
	var validators []*Validator

	stakeAmt := []*big.Int{
		big.NewInt(10),
		big.NewInt(10),
		big.NewInt(10),
		big.NewInt(10),
	}

	for i := 0; i < len(stakeAmt); i++ {
		validators = append(validators, NewValidator(RandomAddr(), stakeAmt[i], common.Address{}))
	}
	validatorList := NewValidatorList(validators)
	for _, vl := range validatorList.Validators {
		validatorList.CalculateAddressRange(vl.Addr, validatorList.StakeBalance(vl.Addr))
	}

	hash := randomHash()
	err, res := validatorList.CollectValidators(hash, 11)
	if err != nil {
		fmt.Println("error collect validators", err)
	}
	for i, v := range res {
		fmt.Println("i===", i, "v====", v, "balance=====", validatorList.StakeBalance(v))
	}
}

func TestCalculateAddressRangeV2(t *testing.T) {
	rand.Seed(time.Now().Unix())
	var validators []*Validator

	addrs := []common.Address{
		common.HexToAddress("0xe2fa892cc5cc268a0cc1d924ec907c796351c645"),
		common.HexToAddress("0xc067825f4b7a53bb9f2daf72ff22c8ee39736aff"),
		common.HexToAddress("0x9e4d5c72569465270232ed7af71981ee82d08dbf"),
		common.HexToAddress("0x612dfa56dca1f581ed34b9c60da86f1268ab6349"),
		common.HexToAddress("0x4110e56ed25e21267fbeef79244f47ada4e2e963"),
		common.HexToAddress("0x31534d5c7b1eabb73425c2361661b878f4322f9d"),
		common.HexToAddress("0x20cb28ae861c322a9a86b4f9e36ad6977930fa05"),
		common.HexToAddress("0x091dbba95b26793515cc9acb9beb5124c479f27f"),
		common.HexToAddress("0xffac4cd934f026dcaf0f9d9eeddcd9af85d8943e"),
		common.HexToAddress("0xdb33217fe3f74bd41c550b06b624e23ab7f55d05"),

		common.HexToAddress("0xbbae84e9879f908700c6ef5d15e928abfb556a21"),
		common.HexToAddress("0xa270bbdff450ebbc2d0413026de5545864a1b6d6"),
		common.HexToAddress("0x84d84e6073a06b6e784241a9b13aa824ab455326"),
		common.HexToAddress("0x7bf72621dd7c4fe4af77632e3177c08f53fdaf09"),

		common.HexToAddress("0x52eae6d396e82358d703bedec2ab11e723127230"),
		common.HexToAddress("0x107837ea83f8f06533ddd3fc39451cd0aa8da8bd"),
	}
	var amts []*big.Int
	for i := 0; i < 8; i++ {
		amt, _ := new(big.Int).SetString("750000000000000000000000", 10)
		amts = append(amts, amt)
	}

	for i := 0; i < 8; i++ {
		amt, _ := new(big.Int).SetString("70000000000000000000000", 10)
		amts = append(amts, amt)
	}
	validators = append(validators, NewValidator(addrs[0], amts[0], common.HexToAddress("0x4d0a8127d3120684cc70ec12e6e8f44ee990b5ac")))
	validators = append(validators, NewValidator(addrs[1], amts[1], common.HexToAddress("0xbad3f0edd751b3b8def4aaddbcf5533ec93452c2")))
	validators = append(validators, NewValidator(addrs[2], amts[2], common.HexToAddress("0x96f2a9f08c92c174700a0bdb452ea737633382a0")))
	validators = append(validators, NewValidator(addrs[3], amts[3], common.HexToAddress("0x66f9e46b49eddc40f0da18d67c07ae755b3643ce")))
	validators = append(validators, NewValidator(addrs[4], amts[4], common.HexToAddress("0x3e6a45b12e2a4e25fb0176c7aa1855459e8e862b")))
	validators = append(validators, NewValidator(addrs[5], amts[5], common.HexToAddress("0x2dbdacc91fd967e2a5c3f04d321752d99a7741c8")))
	validators = append(validators, NewValidator(addrs[6], amts[6], common.HexToAddress("0x36c1550f16c43b5dd85f1379e708d89da9789d9b")))
	validators = append(validators, NewValidator(addrs[7], amts[7], common.HexToAddress("0x8520dc57a2800e417696bdf93553e63bcf31e597")))
	validators = append(validators, NewValidator(addrs[8], amts[8], common.HexToAddress("0x0000000000000000000000000000000000000000")))
	validators = append(validators, NewValidator(addrs[9], amts[9], common.HexToAddress("0x0000000000000000000000000000000000000000")))
	validators = append(validators, NewValidator(addrs[10], amts[10], common.HexToAddress("0x0000000000000000000000000000000000000000")))
	validators = append(validators, NewValidator(addrs[11], amts[11], common.HexToAddress("0x0000000000000000000000000000000000000000")))
	validators = append(validators, NewValidator(addrs[12], amts[12], common.HexToAddress("0x0000000000000000000000000000000000000000")))
	validators = append(validators, NewValidator(addrs[13], amts[13], common.HexToAddress("0x0000000000000000000000000000000000000000")))
	validators = append(validators, NewValidator(addrs[14], amts[14], common.HexToAddress("0x0000000000000000000000000000000000000000")))
	validators = append(validators, NewValidator(addrs[15], amts[15], common.HexToAddress("0x0000000000000000000000000000000000000000")))

	validatorList := NewValidatorList(validators)
	for _, vl := range validatorList.Validators {
		validatorList.CalculateAddressRangeV2(vl.Addr, validatorList.StakeBalance(vl.Addr), big.NewInt(70))
	}

	//---------------------------------------------//
	for _, v := range validatorList.Validators {
		fmt.Println("address---", v.Addr, "|---weight", v.Weight)
	}
}

type ValidatorInfo struct {
	Addr    common.Address
	Balance *big.Int
	Proxy   common.Address
	Weight  []*big.Int
}
type OnlineValidators struct {
	Validators []ValidatorInfo
}

func TestCalculateAddressRangeV2ByOnlineData(t *testing.T) {
	file, err := os.Open("/Users/carver/Desktop/online.json")
	defer file.Close()
	if err != nil {
		fmt.Println("err open file")
	}

	var ovs OnlineValidators
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&ovs)
	if err != nil {
		fmt.Println("err decode")
	}

	var validators []*Validator
	for _, v := range ovs.Validators {
		validators = append(validators, NewValidator(v.Addr, v.Balance, v.Proxy))
	}

	// calculate address range
	validatorList := NewValidatorList(validators)
	for _, vl := range validatorList.Validators {
		validatorList.CalculateAddressRangeV2(vl.Addr, validatorList.StakeBalance(vl.Addr), big.NewInt(70))
	}

	//---------------------------------------------//
	for _, v := range validatorList.Validators {
		fmt.Println("address---", v.Addr, "|---weight", v.Weight)
	}
	for i := 0; i < len(validatorList.Validators); i++ {
		localWeight := validatorList.Validators[i].Weight
		remoteWeight := ovs.Validators[i].Weight
		if remoteWeight[0].Cmp(localWeight[0]) != 0 || remoteWeight[1].Cmp(localWeight[1]) != 0 {
			fmt.Println("failed compare")
		}
	}
	fmt.Println("success compare")
}

func TestCollectValidators(t *testing.T) {
	var count int
	for i := 0; i < 1000; i++ {
		validatorList := prepareFixedValidator()
		for _, vl := range validatorList.Validators {
			validatorList.CalculateAddressRange(vl.Addr, validatorList.StakeBalance(vl.Addr))
		}
		hash := randomHash()
		err, res := validatorList.CollectValidators(hash, 11)
		if err != nil {
			fmt.Println("error collect validators", err)
		}
		if len(res) == 11 {
			count++
		}
	}
	fmt.Println("------", count)
}

func TestStability(t *testing.T) {
	//cannot have random addresses
	validatorList := prepareFixedValidator()
	for _, vl := range validatorList.Validators {
		validatorList.CalculateAddressRange(vl.Addr, validatorList.StakeBalance(vl.Addr))
	}
	hash := common.HexToHash("0xba84ce252cfa4ac640ddbbec2e634d544cb4aacf1880b24126d9300b9a38534d")
	err, res := validatorList.CollectValidators(hash, 11)
	if err != nil {
		fmt.Println("error collect validators", err)
	}
	random11Amt := big.NewInt(0)
	for _, v := range res {
		random11Amt.Add(random11Amt, validatorList.StakeBalance(v))
	}
	fmt.Println("11StakeAmt====", random11Amt, "totalAmt", validatorList.TotalStakeBalance(), "hash", hash.Hex())
}

func prepareFixedValidator() *ValidatorList {
	var validators []*Validator
	stakeAmt := []*big.Int{
		big.NewInt(14),
		big.NewInt(14),
		big.NewInt(13),
		big.NewInt(14),
		big.NewInt(13),
		big.NewInt(12),
		big.NewInt(13),
		big.NewInt(13),
		big.NewInt(14),
		big.NewInt(12),
		big.NewInt(13),
		big.NewInt(13),
		big.NewInt(14),
		big.NewInt(15),
		big.NewInt(16),
		big.NewInt(17),

		//
		big.NewInt(17),
		big.NewInt(1),
		big.NewInt(1),
		big.NewInt(13),
		big.NewInt(13),
		big.NewInt(14),
	}

	addrs := []common.Address{
		common.HexToAddress("0x091DBBa95B26793515cc9aCB9bEb5124c479f27F"),
		common.HexToAddress("0x107837Ea83f8f06533DDd3fC39451Cd0AA8DA8BD"),
		common.HexToAddress("0x612DFa56DcA1F581Ed34b9c60Da86f1268Ab6349"),
		common.HexToAddress("0x84d84e6073A06B6e784241a9B13aA824AB455326"),
		common.HexToAddress("0x9e4d5C72569465270232ed7Af71981Ee82d08dBF"),
		common.HexToAddress("0xa270bBDFf450EbbC2d0413026De5545864a1b6d6"),
		common.HexToAddress("0x4110E56ED25e21267FBeEf79244f47ada4e2E963"),
		common.HexToAddress("0xdb33217fE3F74bD41c550B06B624E23ab7f55d05"),
		common.HexToAddress("0xE2FA892CC5CC268a0cC1d924EC907C796351C645"),
		common.HexToAddress("0x52EAE6D396E82358D703BEDeC2ab11E723127230"),
		common.HexToAddress("0x31534d5C7b1eabb73425c2361661b878F4322f9D"),
		common.HexToAddress("0xbbaE84E9879F908700c6ef5D15e928Abfb556a21"),
		common.HexToAddress("0x20cb28AE861c322A9A86b4F9e36Ad6977930fA05"),
		common.HexToAddress("0xFfAc4cd934f026dcAF0f9d9EEDDcD9af85D8943e"),
		common.HexToAddress("0xc067825f4B7a53Bb9f2Daf72fF22C8EE39736afF"),
		common.HexToAddress("0x7bf72621Dd7C4Fe4AF77632e3177c08F53fdAF09"),

		//
		common.HexToAddress("0xB6FD5851a8c1d9B1C22a210664Fbe7187C137582"),
		common.HexToAddress("0xa4E91908d98aC1b0F232B6873F0989cDE07c7C71"),
		common.HexToAddress("0x1778B78658dDb31a8F0b8ba80E8471225050c62d"),
		common.HexToAddress("0xC12703f9708eE5A5A704696Ea3Dcb0f1c784273a"),
		common.HexToAddress("0xb67Db7D7e97486b3f23369160363430c0e98dFC9"),
		common.HexToAddress("0x7d5412AeA1e796fC58b3A8Be12a2d853528007dC"),
	}

	for i := 0; i < len(addrs); i++ {
		validators = append(validators, NewValidator(addrs[i], stakeAmt[i], common.Address{}))
	}
	validatorList := NewValidatorList(validators)

	return validatorList
}

func RandomAddr() common.Address {
	priKey, _ := crypto.GenerateKey()
	return crypto.PubkeyToAddress(priKey.PublicKey)
}

func TestRandomValidatorV2(t *testing.T) {
	validatorList := prepareFixedValidator()
	for _, vl := range validatorList.Validators {
		validatorList.CalculateAddressRange(vl.Addr, validatorList.StakeBalance(vl.Addr))
	}
	// This hash will calculate less than 15 validators
	// Stable calculation of a validator
	hash := common.HexToHash("0xd95b0361af0635770474311d57c543c90c609b41d18fa0af95583872c3e2ad6f")
	vals := validatorList.RandomValidatorV2(15, hash)
	for i, val := range vals {
		fmt.Println("====i====", i, "====val====", val)
	}
}

func TestProbability(t *testing.T) {
	var count int
	vl := MockValidatorList()
	for i := 0; i < 1000; i++ {
		hash := randomHash()
		selectedValidators := vl.RandomValidatorV2(11, hash)

		// 1.The probability of an address appearing=========================start//
		//for _, v := range selectedValidators {
		//	if v == common.HexToAddress("0x53182359904a07925b3338260b3ee8CD7bAf1935") {
		//		count++
		//	}
		//}
		// The probability of an address appearing=========================end//

		//2. All the validators with 70,000 pledged amount are recorded in===============start//
		var specifiedValidator int
		for _, v := range selectedValidators {
			if vl.StakeBalance(v).Cmp(big.NewInt(800000)) < 0 {
				specifiedValidator++
			}
		}
		if specifiedValidator == 11 {
			count++
		}
		// All the validators with 70,000 pledged amount are recorded in===============end//
	}
	fmt.Println("====count===", count)
}

func MockValidatorList() *ValidatorList {
	var validators []*Validator
	stakeAmt := []*big.Int{
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700001),
		big.NewInt(700001),
		big.NewInt(700001),
		big.NewInt(700001),
		big.NewInt(700001),
		big.NewInt(700001),
		big.NewInt(700001),

		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
	}

	addrs := []common.Address{
		common.HexToAddress("0x091DBBa95B26793515cc9aCB9bEb5124c479f27F"),
		common.HexToAddress("0x107837Ea83f8f06533DDd3fC39451Cd0AA8DA8BD"),
		common.HexToAddress("0x612DFa56DcA1F581Ed34b9c60Da86f1268Ab6349"),
		common.HexToAddress("0x84d84e6073A06B6e784241a9B13aA824AB455326"),
		common.HexToAddress("0x9e4d5C72569465270232ed7Af71981Ee82d08dBF"),
		common.HexToAddress("0xa270bBDFf450EbbC2d0413026De5545864a1b6d6"),
		common.HexToAddress("0x4110E56ED25e21267FBeEf79244f47ada4e2E963"),
		common.HexToAddress("0xdb33217fE3F74bD41c550B06B624E23ab7f55d05"),
		common.HexToAddress("0xE2FA892CC5CC268a0cC1d924EC907C796351C645"),
		common.HexToAddress("0x52EAE6D396E82358D703BEDeC2ab11E723127230"),
		common.HexToAddress("0x31534d5C7b1eabb73425c2361661b878F4322f9D"),
		common.HexToAddress("0xbbaE84E9879F908700c6ef5D15e928Abfb556a21"),
		common.HexToAddress("0x20cb28AE861c322A9A86b4F9e36Ad6977930fA05"),
		common.HexToAddress("0xFfAc4cd934f026dcAF0f9d9EEDDcD9af85D8943e"),
		common.HexToAddress("0xc067825f4B7a53Bb9f2Daf72fF22C8EE39736afF"),
		common.HexToAddress("0x7bf72621Dd7C4Fe4AF77632e3177c08F53fdAF09"),

		//
		common.HexToAddress("0x53182359904a07925b3338260b3ee8CD7bAf1935"),
		common.HexToAddress("0xa01c49F5206DEB26cA7e8b6336E8D61F73b77ba8"),
		common.HexToAddress("0x007a585aC6607F550b6d60eCaEB88b88Ac479daD"),
		common.HexToAddress("0xEe7696c15fc6F3e72Db61867db46d95cC8C8A54F"),
		common.HexToAddress("0x4862713256F2D306029cA4e034292338aF1ef52f"),
		common.HexToAddress("0xc3a690d02FB7132378E7875a03E08Aff7b7Abb27"),
		common.HexToAddress("0xb107cCFd67fc281D9F00D1C41d168A948F0cC469"),
		common.HexToAddress("0x36B651f4CF1E7622783D6E534C667fEd02c559B9"),
		common.HexToAddress("0x7cb6bb45c1C379ADb1665826bcEDAB1c10FC520E"),
		common.HexToAddress("0xa8e9127D4F708aF1c8120dBdaEbb424A9dBb4245"),
		common.HexToAddress("0xA727FB70b907582Db1Ece0098b0784c5d6C128b3"),
		common.HexToAddress("0x8aaa8d0A8cB518b1979e4cF1b9E86607FbCFB9aF"),
		common.HexToAddress("0x84C2fE50aca4BeCF37ca8B240C6D71bA63FF3D93"),
		common.HexToAddress("0x5ff38AeFa2377e93190e05258302d9ffB18ADa0b"),
		common.HexToAddress("0xDc058854d613F2791064CE9Fa9E9Fb46CAE1D312"),
		common.HexToAddress("0x3a78bbE2Cd8863b9aC7A7e6c4906fCF64Ef2c484"),
	}

	for i := 0; i < len(addrs); i++ {
		validators = append(validators, NewValidator(addrs[i], stakeAmt[i], common.Address{}))
	}
	validatorList := NewValidatorList(validators)

	for _, vl := range validatorList.Validators {
		validatorList.CalculateAddressRange(vl.Addr, validatorList.StakeBalance(vl.Addr))
	}
	return validatorList
}

func TestRlpToHash(t *testing.T) {
	validator := NewValidator(RandomAddr(), big.NewInt(1000000), common.Address{})
	encValidator, _ := rlp.EncodeToBytes(validator)
	hash := common.BytesToHash(encValidator)
	t.Log("===hash===", hash.Hex())
}

func AddrList() []common.Address {
	addrs1 := []string{
		// self
		"0x091DBBa95B26793515cc9aCB9bEb5124c479f27F",
		"0x107837Ea83f8f06533DDd3fC39451Cd0AA8DA8BD",
		"0x612DFa56DcA1F581Ed34b9c60Da86f1268Ab6349",
		"0x84d84e6073A06B6e784241a9B13aA824AB455326",
		"0x9e4d5C72569465270232ed7Af71981Ee82d08dBF",
		"0xa270bBDFf450EbbC2d0413026De5545864a1b6d6",
		"0x4110E56ED25e21267FBeEf79244f47ada4e2E963",
		"0xdb33217fE3F74bD41c550B06B624E23ab7f55d05",
		"0xE2FA892CC5CC268a0cC1d924EC907C796351C645",
		"0x52EAE6D396E82358D703BEDeC2ab11E723127230",
		"0x31534d5C7b1eabb73425c2361661b878F4322f9D",
		"0xbbaE84E9879F908700c6ef5D15e928Abfb556a21",
		"0x20cb28AE861c322A9A86b4F9e36Ad6977930fA05",
		"0xFfAc4cd934f026dcAF0f9d9EEDDcD9af85D8943e",
		"0xc067825f4B7a53Bb9f2Daf72fF22C8EE39736afF",
		"0x7bf72621Dd7C4Fe4AF77632e3177c08F53fdAF09",

		// other
		"0xa9D7C42f60879c8Bf5002857D3f943D492A3a4eE",
		"0xf484c55531e0DE69a0dD14AcEE55A18363c0bB18",
		"0x7243552A90505D1e67D84106Ab13eB72DB8337E2",
		"0x412f52Ba4350139b7bf0469781Ac2AB0b5Aa8034",
		"0xEE34dEa224839aCd7d16AB0F160203a2c8DB6e9B",
		"0x37fF4076b8cA98f5ce00EcDB5841033A7D231142",
		"0xde31fb555169028FcD34CBf99927875E10b552f1",
		"0x785067Bf5Da2d72d0feE45b51e04f82F81527174",
		"0xF93A2E5c94a315272AAD31c6fB02052E121121e1",
		"0xF6a45b30DF36105048A22f550b63f84AB52fd6AA",
		"0x79B9E94d490151fdc797fCD7B174dc0561ec5740",
		"0x2f21d7D75ECD9ac488B44E6d4295A9d7BFCB44Ad",
		"0x1f63bDC4dF28799689119829334B5b584Cae3Fd6",
		"0xD50230828117B4801013065B890B597c8F563428",
		"0xEAa83d8bd0c0362430F7A09774512FE2D37e1Afb",
		"0xE56f92DE3789B53fc5198153f9b53fdc2Ee0a778",
		"0xdD1F8D80d758766A7D959D1c3348d808ff7C2102",
		"0xb249F79f1A752Caf3f5378a93671D2996A7fBFF7",
		"0x83BDd29e1Bddacd37295BC2033160BEc66F47c23",
		"0xb0584C9b370497788261dE542aCb33d4aAe6952f",
		"0x4448C92dfB560c6F4D2Df371CF12A4fF441a2fe4",
		"0xB53df3E6d295ADCC61C40820E9CaEf5653d4D044",
		"0x27009F63C4d01Bb5deEFecdf70B6D08bC8edA720",
		"0x41D7FDaF014A850D1AE8D14b76bF6A91445647b7",
		"0xd4989D676F893e0c9E585b288dc12c527a3F0f99",
		"0xfa0db344597DE0bDa9a3dfb65eD93C89ff2fb883",
		"0xa2d0E0fc679A8b86B9214bBb3F25FCd541f6a3BE",
		"0x97AC70B023Eb8D921E8F2CE58cdB321F857Cf8aF",
		"0xB476A63842f0479A68C2c276B32BA741f0AB4347",
		"0x887DB52dfB96C742Ca475EF8eA33969DAE5ea7Be",
		"0xC2BAc3E82f5f47156e0bb4b53cd5667EC5eD3488",
		"0x9329f2370b56DC82cBDe5f927e15Fe29a6b4CFa1",
		"0x8563c57aF7d7B38b8D1859030e23cf2eF7e8134A",
		"0x5A9e4B0D5Ed9358017f5789314FbBb47Cf74d6C6",
		"0xA002Fb8E7eBD8633124CE0ffcd63b2D435FF4429",
		"0xd0735120EE48Ef2c86C1200A48F9096d57A48f97",
		"0x1662b48A65c4883F2a6C1a758041929a81B5528a",
		"0xC3e2a3aB58fF8aF53761Ef99EE1fb69244dcC018",
	}

	var addrs2 []common.Address
	for _, addr := range addrs1 {
		addrs2 = append(addrs2, common.HexToAddress(addr))
	}
	return addrs2
}

func GetOtherMinerAddr() []common.Address {
	addr1 := []string{
		// other
		"0xa9D7C42f60879c8Bf5002857D3f943D492A3a4eE",
		"0xf484c55531e0DE69a0dD14AcEE55A18363c0bB18",
		"0x7243552A90505D1e67D84106Ab13eB72DB8337E2",
		"0x412f52Ba4350139b7bf0469781Ac2AB0b5Aa8034",
		"0xEE34dEa224839aCd7d16AB0F160203a2c8DB6e9B",
		"0x37fF4076b8cA98f5ce00EcDB5841033A7D231142",
		"0xde31fb555169028FcD34CBf99927875E10b552f1",
		"0x785067Bf5Da2d72d0feE45b51e04f82F81527174",
		"0xF93A2E5c94a315272AAD31c6fB02052E121121e1",
		"0xF6a45b30DF36105048A22f550b63f84AB52fd6AA",
		"0x79B9E94d490151fdc797fCD7B174dc0561ec5740",
		"0x2f21d7D75ECD9ac488B44E6d4295A9d7BFCB44Ad",
		"0x1f63bDC4dF28799689119829334B5b584Cae3Fd6",
		"0xD50230828117B4801013065B890B597c8F563428",
		"0xEAa83d8bd0c0362430F7A09774512FE2D37e1Afb",
		"0xE56f92DE3789B53fc5198153f9b53fdc2Ee0a778",
		"0xdD1F8D80d758766A7D959D1c3348d808ff7C2102",
		"0xb249F79f1A752Caf3f5378a93671D2996A7fBFF7",
		"0x83BDd29e1Bddacd37295BC2033160BEc66F47c23",
		"0xb0584C9b370497788261dE542aCb33d4aAe6952f",
		"0x4448C92dfB560c6F4D2Df371CF12A4fF441a2fe4",
		"0xB53df3E6d295ADCC61C40820E9CaEf5653d4D044",
		"0x27009F63C4d01Bb5deEFecdf70B6D08bC8edA720",
		"0x41D7FDaF014A850D1AE8D14b76bF6A91445647b7",
		"0xd4989D676F893e0c9E585b288dc12c527a3F0f99",
		"0xfa0db344597DE0bDa9a3dfb65eD93C89ff2fb883",
		"0xa2d0E0fc679A8b86B9214bBb3F25FCd541f6a3BE",
		"0x97AC70B023Eb8D921E8F2CE58cdB321F857Cf8aF",
		"0xB476A63842f0479A68C2c276B32BA741f0AB4347",
		"0x887DB52dfB96C742Ca475EF8eA33969DAE5ea7Be",
		"0xC2BAc3E82f5f47156e0bb4b53cd5667EC5eD3488",
		"0x9329f2370b56DC82cBDe5f927e15Fe29a6b4CFa1",
		"0x8563c57aF7d7B38b8D1859030e23cf2eF7e8134A",
		"0x5A9e4B0D5Ed9358017f5789314FbBb47Cf74d6C6",
		"0xA002Fb8E7eBD8633124CE0ffcd63b2D435FF4429",
		"0xd0735120EE48Ef2c86C1200A48F9096d57A48f97",
		"0x1662b48A65c4883F2a6C1a758041929a81B5528a",
		"0xC3e2a3aB58fF8aF53761Ef99EE1fb69244dcC018",
	}
	var addrs2 []common.Address
	for _, addr := range addr1 {
		addrs2 = append(addrs2, common.HexToAddress(addr))
	}
	return addrs2
}

func GetSelfAddr() []common.Address {
	addrs1 := []string{
		// self
		"0x8520dc57A2800e417696bdF93553E63bCF31e597",
		"0x107837Ea83f8f06533DDd3fC39451Cd0AA8DA8BD",
		"0x66f9e46b49EDDc40F0dA18D67C07ae755b3643CE",
		"0x84d84e6073A06B6e784241a9B13aA824AB455326",
		"0x96f2A9f08c92c174700A0bdb452EA737633382A0",
		"0xa270bBDFf450EbbC2d0413026De5545864a1b6d6",
		"0x3E6a45b12E2A4E25fb0176c7Aa1855459E8e862b",
		"0xdb33217fE3F74bD41c550B06B624E23ab7f55d05",
		"0x4d0A8127D3120684CC70eC12e6E8F44eE990b5aC",
		"0x52EAE6D396E82358D703BEDeC2ab11E723127230",
		"0x2DbdaCc91fd967E2A5C3F04D321752d99a7741C8",
		"0xbbaE84E9879F908700c6ef5D15e928Abfb556a21",
		"0x36c1550F16c43B5Dd85f1379E708d89DA9789d9b",
		"0xFfAc4cd934f026dcAF0f9d9EEDDcD9af85D8943e",
		"0xbad3F0edd751B3b8DeF4AaDDbcF5533eC93452C2",
		"0x7bf72621Dd7C4Fe4AF77632e3177c08F53fdAF09",
	}

	var addrs2 []common.Address
	for _, addr := range addrs1 {
		addrs2 = append(addrs2, common.HexToAddress(addr))
	}
	return addrs2
}

func TestValidatorAccordToDistance(t *testing.T) {
	addrs := AddrList()

	c := 140000
	c2 := 70000
	stakeAmt := []*big.Int{
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),

		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
	}

	var validators []*Validator
	for i := 0; i < len(addrs); i++ {
		validators = append(validators, NewValidator(addrs[i], stakeAmt[i], common.Address{}))
	}
	validatorList := NewValidatorList(validators)

	for _, vl := range validatorList.Validators {
		validatorList.CalculateAddressRange(vl.Addr, validatorList.StakeBalance(vl.Addr))
	}

	otherMiners := GetOtherMinerAddr()
	var count int
	for i := 0; i < 2; i++ {
		randomHash := randomHash()
		consensusValidator := validatorList.RandomValidatorV2(11, randomHash)
		rewardAddr := validatorList.ValidatorByDistance(ConvertToBigInt(consensusValidator), 6, randomHash)
		flg := false
		for i := 0; i < len(otherMiners); i++ {
			for _, v := range rewardAddr {
				if v.Hex() == otherMiners[i].Hex() {
					flg = true
					fmt.Println("===vvvvv===", v.Hex())
					break
				} else {
					flg = false
					continue
				}
			}
			if flg {
				count++
				break
			}
		}
	}
	fmt.Println("===reward  other addr count===", count, "time", time.Now().Unix())
}

func TestValidatorByDistance(t *testing.T) {
	randomHash := randomHash()

	addrs := AddrList()

	stakeAmt := []*big.Int{
		big.NewInt(80000),
		big.NewInt(80000),
		big.NewInt(80000),
		big.NewInt(80000),
		big.NewInt(80000),
		big.NewInt(80000),
		big.NewInt(80000),
		big.NewInt(80000),
		big.NewInt(80000),
		big.NewInt(80000),
		big.NewInt(80000),
		big.NewInt(80000),
		big.NewInt(80000),
		big.NewInt(80000),
		big.NewInt(80000),
		big.NewInt(80000),

		big.NewInt(72000),
		big.NewInt(72000),
		big.NewInt(72000),
		big.NewInt(72000),
		big.NewInt(72000),
		big.NewInt(72000),
		big.NewInt(72000),
		big.NewInt(72000),
		big.NewInt(72000),
		big.NewInt(72000),
		big.NewInt(72000),
		big.NewInt(72000),
		big.NewInt(72000),
		big.NewInt(72000),
		big.NewInt(72000),
		big.NewInt(72000),
		big.NewInt(72000),
		big.NewInt(72000),
		big.NewInt(72000),
		big.NewInt(72000),
		big.NewInt(72000),
		big.NewInt(72000),
		big.NewInt(72000),
		big.NewInt(72000),
		big.NewInt(72000),
		big.NewInt(72000),
		big.NewInt(72000),
		big.NewInt(72000),
		big.NewInt(72000),
		big.NewInt(72000),
		big.NewInt(72000),
		big.NewInt(72000),
		big.NewInt(72000),
		big.NewInt(72000),
		big.NewInt(72000),
		big.NewInt(72000),
		big.NewInt(72000),
		big.NewInt(72000),
	}

	var validators []*Validator
	for i := 0; i < len(addrs); i++ {
		validators = append(validators, NewValidator(addrs[i], stakeAmt[i], common.Address{}))
	}
	validatorList := NewValidatorList(validators)

	rewardAddr := validatorList.ValidatorByDistance(ConvertToBigInt(addrs), 6, randomHash)

	for _, v := range rewardAddr {
		fmt.Println("====rewardAddr====", v)
	}
}

func ConvertToBigInt(addrs []common.Address) (bigIntSlice []*big.Int) {
	for _, m := range addrs {
		bigIntSlice = append(bigIntSlice, m.Hash().Big())
	}
	return
}

func TestExsistSevenValidatorProbability(t *testing.T) {
	addrs := AddrList()

	c := 400000
	c2 := 70000
	stakeAmt := []*big.Int{
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),

		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
	}

	var validators []*Validator
	for i := 0; i < len(addrs); i++ {
		validators = append(validators, NewValidator(addrs[i], stakeAmt[i], common.Address{}))
	}
	validatorList := NewValidatorList(validators)

	for _, vl := range validatorList.Validators {
		validatorList.CalculateAddressRange(vl.Addr, validatorList.StakeBalance(vl.Addr))
	}

	selfAddrs := GetSelfAddr()
	var count int
	for i := 0; i < 20000; i++ {
		randomHash := randomHash()
		consensusValidator := validatorList.RandomValidatorV2(11, randomHash)

		var occurCount int

		for j := 0; j < len(selfAddrs); j++ {
			for _, v := range consensusValidator {
				if v.Hex() == selfAddrs[j].Hex() {
					occurCount++
					break
				}
			}
		}
		if occurCount >= 7 {
			count++
		}
	}
	fmt.Println("===Probability of occurrence===", count, "time", time.Now().Unix())
}

func TestOtherValidatorProbability(t *testing.T) {
	addrs := AddrList()

	c := 400000
	c2 := 70000
	stakeAmt := []*big.Int{
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),

		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
	}

	var validators []*Validator
	for i := 0; i < len(addrs); i++ {
		validators = append(validators, NewValidator(addrs[i], stakeAmt[i], common.Address{}))
	}
	validatorList := NewValidatorList(validators)

	for _, vl := range validatorList.Validators {
		validatorList.CalculateAddressRange(vl.Addr, validatorList.StakeBalance(vl.Addr))
	}

	otherAddrs := GetOtherMinerAddr()
	var count int
	for i := 0; i < 200; i++ {
		randomHash := randomHash()
		consensusValidator := validatorList.RandomValidatorV2(11, randomHash)

		var occurCount int

		for j := 0; j < len(otherAddrs); j++ {
			for _, v := range consensusValidator {
				if v.Hex() == otherAddrs[j].Hex() {
					occurCount++
					break
				}
			}
		}
		count += occurCount
	}
	fmt.Println("===Probability of occurrence===", float32(count)/20000, "time", time.Now().Unix())
}

func TestConsensus(t *testing.T) {
	c := 400000
	c2 := 70000
	stakeAmt := []*big.Int{
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),

		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
	}

	addrs := []common.Address{
		common.HexToAddress("0xFfAc4cd934f026dcAF0f9d9EEDDcD9af85D8943e"),
		common.HexToAddress("0xE2FA892CC5CC268a0cC1d924EC907C796351C645"),
		common.HexToAddress("0xdb33217fE3F74bD41c550B06B624E23ab7f55d05"),
		common.HexToAddress("0xc067825f4B7a53Bb9f2Daf72fF22C8EE39736afF"),
		common.HexToAddress("0xbbaE84E9879F908700c6ef5D15e928Abfb556a21"),
		common.HexToAddress("0xa270bBDFf450EbbC2d0413026De5545864a1b6d6"),
		common.HexToAddress("0x9e4d5C72569465270232ed7Af71981Ee82d08dBF"),
		common.HexToAddress("0x84d84e6073A06B6e784241a9B13aA824AB455326"),
		common.HexToAddress("0x7bf72621Dd7C4Fe4AF77632e3177c08F53fdAF09"),
		common.HexToAddress("0x612DFa56DcA1F581Ed34b9c60Da86f1268Ab6349"),
		common.HexToAddress("0x52EAE6D396E82358D703BEDeC2ab11E723127230"),
		common.HexToAddress("0x4110E56ED25e21267FBeEf79244f47ada4e2E963"),
		common.HexToAddress("0x31534d5C7b1eabb73425c2361661b878F4322f9D"),
		common.HexToAddress("0x20cb28AE861c322A9A86b4F9e36Ad6977930fA05"),
		common.HexToAddress("0x107837Ea83f8f06533DDd3fC39451Cd0AA8DA8BD"),
		common.HexToAddress("0x091DBBa95B26793515cc9aCB9bEb5124c479f27F"),

		//
		common.HexToAddress("0xAc17A48d782DE6985D471F6EEcB1023C08f0CB05"),
		common.HexToAddress("0xa337175D23bB0ee4f879e38B73973CDe829b9d29"),
		common.HexToAddress("0x9baeC4D975B0Dd5baa904200Ea11727eDD593be6"),
		common.HexToAddress("0x77c33e2951d2851D9D648161BEfAC8f4C60D1181"),
		common.HexToAddress("0x5a7f652BC51Fb99747fe2641A8CDd6CbFE51b201"),
		common.HexToAddress("0x2c2909db351764D92ecc313a0D8baF72735C5165"),
	}

	var validators []*Validator
	for i := 0; i < len(addrs); i++ {
		validators = append(validators, NewValidator(addrs[i], stakeAmt[i], common.Address{}))
	}
	validatorList := NewValidatorList(validators)

	for _, vl := range validatorList.Validators {
		validatorList.CalculateAddressRange(vl.Addr, validatorList.StakeBalance(vl.Addr))
	}

	var count1 int
	var count2 int
	var count3 int
	var count4 int
	var count5 int
	var count6 int
	for i := 0; i < 200; i++ {
		randomHash := randomHash()
		fmt.Println(randomHash)
		consensusValidator := validatorList.RandomValidatorV2(11, randomHash)
		for _, v := range consensusValidator {
			if v.Hex() == "0xAc17A48d782DE6985D471F6EEcB1023C08f0CB05" {
				count1++
			}
			if v.Hex() == "0xa337175D23bB0ee4f879e38B73973CDe829b9d29" {
				count2++
			}
			if v.Hex() == "0x9baeC4D975B0Dd5baa904200Ea11727eDD593be6" {
				count3++
			}
			if v.Hex() == "0x77c33e2951d2851D9D648161BEfAC8f4C60D1181" {
				count4++
			}
			if v.Hex() == "0x5a7f652BC51Fb99747fe2641A8CDd6CbFE51b201" {
				count5++
			}
			if v.Hex() == "0x2c2909db351764D92ecc313a0D8baF72735C5165" {
				count6++
			}
		}
	}
	fmt.Print("=====v====", count1, count2, count3, count4, count5, count6)
}

func TestRandomValidatorsV3By16Addr(t *testing.T) {
	c, _ := new(big.Int).SetString("750000000000000000000000000", 10)
	c2, _ := new(big.Int).SetString("70000000000000000000000000", 10)

	stakeAmt := []*big.Int{
		c,
		c2,
		c,
		c2,
		c,
		c2,
		c,
		c2,
		c,
		c2,
		c,
		c2,
		c,
		c2,
		c,
		c2,
	}

	addrs := GetSelfAddr()

	// initial count
	countMap := make(map[common.Address]int, 0)
	for _, v := range addrs {
		countMap[v] = 0
	}

	var validators []*Validator
	for i := 0; i < len(addrs); i++ {
		validators = append(validators, NewValidator(addrs[i], stakeAmt[i], common.Address{}))
	}
	validatorList := NewValidatorList(validators)

	for _, vl := range validatorList.Validators {
		validatorList.CalculateAddressRange(vl.Addr, validatorList.StakeBalance(vl.Addr))
	}

	for i := 0; i < 200; i++ {
		randomHash := randomHash()
		fmt.Println("====randomhash===", randomHash)
		consensusValidator := validatorList.RandomValidatorV3(11, randomHash)
		for _, v := range consensusValidator {
			countMap[v]++
		}
	}

	for k, v := range countMap {
		fmt.Println("addr====", k, "====count====", v)
	}
}

func TestRandomValidatorsV3By400Addr(t *testing.T) {
	c, _ := new(big.Int).SetString("750000000000000000000000000", 10)
	c2, _ := new(big.Int).SetString("70000000000000000000000000", 10)

	stakeAmt := []*big.Int{
		c,
		c2,
		c,
		c2,
		c,
		c2,
		c,
		c2,
		c,
		c2,
		c,
		c2,
		c,
		c2,
		c,
		c2,
	}

	for i := 0; i < 384; i++ {
		stakeAmt = append(stakeAmt, c2.Add(c2, big.NewInt(int64(i))))
	}

	addrs := GetSelfAddr()

	for i := 0; i < 384; i++ {
		addrs = append(addrs, RandomAddr())
	}

	// initial count
	countMap := make(map[common.Address]int, 0)
	for _, v := range addrs {
		countMap[v] = 0
	}

	var validators []*Validator
	for i := 0; i < len(addrs); i++ {
		validators = append(validators, NewValidator(addrs[i], stakeAmt[i], common.Address{}))
	}
	validatorList := NewValidatorList(validators)

	for _, vl := range validatorList.Validators {
		validatorList.CalculateAddressRange(vl.Addr, validatorList.StakeBalance(vl.Addr))
	}

	for i := 0; i < 1000; i++ {
		randomHash := randomHash()
		fmt.Println("====randomhash===", randomHash)
		consensusValidator := validatorList.RandomValidatorV3(11, randomHash)
		for _, v := range consensusValidator {
			countMap[v]++
		}
	}

	for k, v := range countMap {
		fmt.Println("addr====", k, "====count====", v)
	}
}

const allocData = "" +
	"0x091DBBa95B26793515cc9aCB9bEb5124c479f27F:0x9ed194db19b238c00000," +
	"0x107837Ea83f8f06533DDd3fC39451Cd0AA8DA8BD:0xed2b525841adfc00000," +
	"0x612DFa56DcA1F581Ed34b9c60Da86f1268Ab6349:0x9ed194db19b238c00000," +
	"0x84d84e6073A06B6e784241a9B13aA824AB455326:0xed2b525841adfc00000," +
	"0x9e4d5C72569465270232ed7Af71981Ee82d08dBF:0x9ed194db19b238c00000," +
	"0xa270bBDFf450EbbC2d0413026De5545864a1b6d6:0xed2b525841adfc00000," +
	"0x4110E56ED25e21267FBeEf79244f47ada4e2E963:0x9ed194db19b238c00000," +
	"0xdb33217fE3F74bD41c550B06B624E23ab7f55d05:0xed2b525841adfc00000," +
	"0xE2FA892CC5CC268a0cC1d924EC907C796351C645:0x9ed194db19b238c00000," +
	"0x52EAE6D396E82358D703BEDeC2ab11E723127230:0xed2b525841adfc00000," +
	"0x31534d5C7b1eabb73425c2361661b878F4322f9D:0x9ed194db19b238c00000," +
	"0xbbaE84E9879F908700c6ef5D15e928Abfb556a21:0xed2b525841adfc00000," +
	"0x20cb28AE861c322A9A86b4F9e36Ad6977930fA05:0x9ed194db19b238c00000," +
	"0xFfAc4cd934f026dcAF0f9d9EEDDcD9af85D8943e:0xed2b525841adfc00000," +
	"0xc067825f4B7a53Bb9f2Daf72fF22C8EE39736afF:0x9ed194db19b238c00000," +
	"0x7bf72621Dd7C4Fe4AF77632e3177c08F53fdAF09:0xed2b525841adfc00000"

func GetBalanceInfos(data string) []*BalanceInfo {

	var balanceInfos []*BalanceInfo

	accountInfos := strings.Split(data, ",")
	for _, accountInfo := range accountInfos {
		splitAccountInfo := strings.Split(accountInfo, ":")
		address := common.HexToAddress(splitAccountInfo[0])
		strBalance := string([]byte(splitAccountInfo[1])[2:])
		balance, _ := new(big.Int).SetString(strBalance, 16)
		balanceInfo := BalanceInfo{
			Address: address,
			Balance: balance,
		}

		balanceInfos = append(balanceInfos, &balanceInfo)
	}

	return balanceInfos
}

func TestRandomValidatorsV3(t *testing.T) {
	addressNum := make(map[common.Address]int, 0)
	//hash := md5.New()
	randomAddress := make(map[string]int, 0)

	//hash := sha256.New()
	infos := GetBalanceInfos(allocData)

	tempValidatorList := &ValidatorList{
		Validators: make([]*Validator, 0, len(infos)),
	}

	for _, info := range infos {
		tempValidator := Validator{
			Addr:    info.Address,
			Balance: new(big.Int).Set(info.Balance),
		}
		tempValidatorList.Validators = append(tempValidatorList.Validators, &tempValidator)
	}

	//for _, v := range tempValidatorList.Validators {
	//	t.Log("before", v.Weight)
	//}

	for _, vl := range tempValidatorList.Validators {
		tempValidatorList.CalculateAddressRange(vl.Addr, tempValidatorList.StakeBalance(vl.Addr))
	}

	//for _, v := range tempValidatorList.Validators {
	//	t.Log("after", v.Weight)
	//}

	for i := 0; i < 20000; i++ {
		sum := crypto.Keccak256([]byte(strconv.Itoa(i)))
		//sum = crypto.Keccak256Hash(sum.Bytes())
		//sum := hash.Sum([]byte(strconv.Itoa(i)))
		//sum = hash.Sum(sum)

		t.Log("sum: ", hex.EncodeToString(sum))
		rr := common.BytesToHash(sum).Hex()
		pri, _ := crypto.HexToECDSA(rr[2:])
		addr := crypto.PubkeyToAddress(pri.PublicKey)
		t.Log(addr.Hex())
		randomAddress[strings.ToLower(string([]byte(addr.String())[:3]))]++
		addressArr := tempValidatorList.RandomValidatorV3(11, common.BytesToHash(sum))
		//arr := tempValidatorList.InitAddressArr(infos)
		//addressArr := tempValidatorList.SelectRandom11Address(11, arr, sum)
		for _, address := range addressArr {
			//t.Log(address.Hex())
			addressNum[address]++
		}
	}

	//infos := GetBalanceInfos(devnetAllocData)
	base, _ := new(big.Int).SetString("1000000000000000000000", 10)
	for _, info := range infos {
		t.Log(info.Address.Hex(), new(big.Int).Div(info.Balance, base), addressNum[info.Address])
	}
	t.Log(tempValidatorList.Len())

	for k, v := range randomAddress {
		t.Log(k, v)
	}
}

func TestMockData(t *testing.T) {
	c1 := 5000000
	c2 := 140000
	c3 := 70000
	var stakeAmt []*big.Int
	for i := 0; i < 8; i++ {
		stakeAmt = append(stakeAmt, big.NewInt(int64(c1)))
	}

	for i := 0; i < 8; i++ {
		stakeAmt = append(stakeAmt, big.NewInt(int64(c2)))
	}

	for i := 0; i < 131; i++ {
		stakeAmt = append(stakeAmt, big.NewInt(int64(c3)))
	}

	addrs := []common.Address{
		common.HexToAddress("0x8520dc57A2800e417696bdF93553E63bCF31e597"),
		common.HexToAddress("0x66f9e46b49EDDc40F0dA18D67C07ae755b3643CE"),
		common.HexToAddress("0x96f2A9f08c92c174700A0bdb452EA737633382A0"),
		common.HexToAddress("0x3E6a45b12E2A4E25fb0176c7Aa1855459E8e862b"),
		common.HexToAddress("0x4d0A8127D3120684CC70eC12e6E8F44eE990b5aC"),
		common.HexToAddress("0x2DbdaCc91fd967E2A5C3F04D321752d99a7741C8"),
		common.HexToAddress("0x36c1550F16c43B5Dd85f1379E708d89DA9789d9b"),
		common.HexToAddress("0xbad3F0edd751B3b8DeF4AaDDbcF5533eC93452C2"),

		common.HexToAddress("0x107837Ea83f8f06533DDd3fC39451Cd0AA8DA8BD"),
		common.HexToAddress("0x84d84e6073A06B6e784241a9B13aA824AB455326"),
		common.HexToAddress("0xa270bBDFf450EbbC2d0413026De5545864a1b6d6"),
		common.HexToAddress("0xdb33217fE3F74bD41c550B06B624E23ab7f55d05"),
		common.HexToAddress("0x52EAE6D396E82358D703BEDeC2ab11E723127230"),
		common.HexToAddress("0xbbaE84E9879F908700c6ef5D15e928Abfb556a21"),
		common.HexToAddress("0xFfAc4cd934f026dcAF0f9d9EEDDcD9af85D8943e"),
		common.HexToAddress("0x7bf72621Dd7C4Fe4AF77632e3177c08F53fdAF09"),

		//

		common.HexToAddress("0xef2B17d61867e9504E8E7e6f18DE1f4c7EA62c9a"),
		common.HexToAddress("0x4854F8324009AFDC20C5f651D70fFA5eF6c036B8"),
		common.HexToAddress("0x93c70c1f04E0932630F197E7FecDD0BB51eC507c"),
		common.HexToAddress("0xfb3e3f2B6A57D1A01e9dD76b77d3579A2AD4A431"),
		common.HexToAddress("0xfa0db344597DE0bDa9a3dfb65eD93C89ff2fb883"),
		common.HexToAddress("0xF6a45b30DF36105048A22f550b63f84AB52fd6AA"),
		common.HexToAddress("0xf51525083022a1037CF7a18Fa580433d7d36f1a8"),
		common.HexToAddress("0xf408671688cfbDCCF74796Bba1DA03Ba0b19C8E5"),
		common.HexToAddress("0xeb998e3095Fd10999ff93856C6FA4243596b2A1C"),
		common.HexToAddress("0xeAFbC420F4c4E6a811d77BC5dBCF413fd48C4A22"),
		common.HexToAddress("0xEAa83d8bd0c0362430F7A09774512FE2D37e1Afb"),
		common.HexToAddress("0xEA0359b30CfD247fA1039336618A150c4A960205"),
		common.HexToAddress("0xE7509C40581aEEcC4e2A763776149A736DccF257"),
		common.HexToAddress("0xE56f92DE3789B53fc5198153f9b53fdc2Ee0a778"),
		common.HexToAddress("0xE375006eF60470cB2583444033Be7302b2A16f16"),
		common.HexToAddress("0xde31fb555169028FcD34CBf99927875E10b552f1"),
		common.HexToAddress("0xdD1F8D80d758766A7D959D1c3348d808ff7C2102"),
		common.HexToAddress("0xd88D58a309Fa42f3552A78EC17a1390B56AbB675"),
		common.HexToAddress("0xD50230828117B4801013065B890B597c8F563428"),
		common.HexToAddress("0xd4989D676F893e0c9E585b288dc12c527a3F0f99"),
		common.HexToAddress("0xd0735120EE48Ef2c86C1200A48F9096d57A48f97"),
		common.HexToAddress("0xcf76A03eC42bBd719D0f3d24bb39346b55BA4883"),
		common.HexToAddress("0xcf615F48877448292B48D209803b0E531d70CAAa"),
		common.HexToAddress("0xCeEA027DC415F859D4cAF6d876c58c84800f3b2b"),
		common.HexToAddress("0xc5c7Efe59bd13D9Ee1583C799CfFb5F4d0183c5c"),
		common.HexToAddress("0xC3e2a3aB58fF8aF53761Ef99EE1fb69244dcC018"),
		common.HexToAddress("0xC2BAc3E82f5f47156e0bb4b53cd5667EC5eD3488"),
		common.HexToAddress("0xC1FCB96582203B5E57C552b07C8142bde16aAB10"),
		common.HexToAddress("0xC0Fe83dd51372AdB41e86caCF6C2daEdF8642d0c"),
		common.HexToAddress("0xBa655cf6890E4FaD0Ea17CaF150D4d8c3518383b"),
		common.HexToAddress("0xB9f6da86b3aF1e75a168158cDF033bAAc21F227c"),
		common.HexToAddress("0xb760d060dd94604693250Ed957CDD30329D11BDc"),
		common.HexToAddress("0xB53df3E6d295ADCC61C40820E9CaEf5653d4D044"),
		common.HexToAddress("0xB476A63842f0479A68C2c276B32BA741f0AB4347"),
		common.HexToAddress("0xb4289dbd91D061900394937AA3EABAbF7bAf443b"),
		common.HexToAddress("0xb249F79f1A752Caf3f5378a93671D2996A7fBFF7"),
		common.HexToAddress("0xb13cDCaC9A005Cc7B8360C981e9d7d645b73e8be"),
		common.HexToAddress("0xb0584C9b370497788261dE542aCb33d4aAe6952f"),
		common.HexToAddress("0xaD65E213e361EA2d0591F762565334f300199035"),
		common.HexToAddress("0xAAb0620588c411d4d5E17238b912b8A34C220D72"),
		common.HexToAddress("0xa9D7C42f60879c8Bf5002857D3f943D492A3a4eE"),
		common.HexToAddress("0xA9199f935f922CDbe759101c8f1509d1E95bf7Eb"),
		common.HexToAddress("0xa8270d82827440A13DB2086d96eCaBACC1f2e047"),
		common.HexToAddress("0xA6ea760Ad0C57954A4E93cf8d3873e4786141EeD"),
		common.HexToAddress("0xa5DaC53e972b35911Eb46ce0e0573aE7d911c206"),
		common.HexToAddress("0xA36a61046DD92aF5544348A48bF9Ab0cF29979D7"),
		common.HexToAddress("0xa2d0E0fc679A8b86B9214bBb3F25FCd541f6a3BE"),
		common.HexToAddress("0xA19c58Cd85125B8D659a7047DA709957E4063B84"),
		common.HexToAddress("0xA002Fb8E7eBD8633124CE0ffcd63b2D435FF4429"),
		common.HexToAddress("0x9eA566069E5984807dd69F8Fb1aE848f7a41c283"),
		common.HexToAddress("0x9E9D8EF8cb9141ba9d33aa9eB83F355db27A3b75"),
		common.HexToAddress("0x9E17B9A2a577d41f6D21dc79cD5396C10A66Eef1"),
		common.HexToAddress("0x984800d6469b4F72A7b344f4E1A71EaB79d22362"),
		common.HexToAddress("0x97AC70B023Eb8D921E8F2CE58cdB321F857Cf8aF"),
		common.HexToAddress("0x9329f2370b56DC82cBDe5f927e15Fe29a6b4CFa1"),
		common.HexToAddress("0x91E7A9740AE36d01ab04f6aCEb9De570C2f67C89"),
		common.HexToAddress("0x91559e92bEb6637864Cf44AB24D658814f0b75D5"),
		common.HexToAddress("0x8cE2Ea28977ce112b27eD3954Ea9C90b75623159"),
		common.HexToAddress("0x8b13Ba3418B1E692eFDFdf6e83BAe2E91dEa9755"),
		common.HexToAddress("0x887DB52dfB96C742Ca475EF8eA33969DAE5ea7Be"),
		common.HexToAddress("0x87eACE22c89b12D4B05fdE2c61142673efD931A0"),
		common.HexToAddress("0x869ccE356aD7421F20961D215a713041e03F886d"),
		common.HexToAddress("0x8563c57aF7d7B38b8D1859030e23cf2eF7e8134A"),
		common.HexToAddress("0x8436Eba72B44bc4EB1e8BDf634b357b2EDc6086c"),
		common.HexToAddress("0x83c80534F316148B726646d1c1cfD81fCb209645"),
		common.HexToAddress("0x83BDd29e1Bddacd37295BC2033160BEc66F47c23"),
		common.HexToAddress("0x8358A78eBb16e92F44726d6fBB863760C9131196"),
		common.HexToAddress("0x8103ebe9f220505BEa89eA087d334557a95df5e5"),
		common.HexToAddress("0x7fe4359ce70a1222acAc3a9b3e286e567042Bb19"),
		common.HexToAddress("0x7dE68FCEe94edD20b5e197b2D47B662b10276302"),
		common.HexToAddress("0x79B9E94d490151fdc797fCD7B174dc0561ec5740"),
		common.HexToAddress("0x785067Bf5Da2d72d0feE45b51e04f82F81527174"),
		common.HexToAddress("0x77DD8E405e6CFAEfc40Baa6355c8b772EcB615D7"),
		common.HexToAddress("0x7499268010Ce2D2D9D01ed66E33FbC67ECA6A627"),
		common.HexToAddress("0x7316759F684a589b1e03A00d7B98ad5D5E6aD303"),
		common.HexToAddress("0x71a7d3FBC09b41B9328Ba1cEae77AB1423015cF9"),
		common.HexToAddress("0x701984B285Ec098f84eF1D6B3381398D9bd2f6dE"),
		common.HexToAddress("0x6EAC39375777d8ce44d97Bc2fEf2AcAe1cbc5750"),
		common.HexToAddress("0x6bD0441aF0F4c2d972E3498456c5feC4127abaFF"),
		common.HexToAddress("0x67CBfA0fFCc28C25E750A08a8ca162F62742459D"),
		common.HexToAddress("0x6751B9E0938d817Df9Bf28e260DfD4D85444BCa2"),
		common.HexToAddress("0x630400eb2846a16A75E20530DB2549b3FEdD9309"),
		common.HexToAddress("0x5e9FCA1351587F4df9c8ffb72094683eb0E0e7C7"),
		common.HexToAddress("0x5D725C416a3Dde0620D6eE1B3f1192B8Df993cda"),
		common.HexToAddress("0x5bD7498DA0FDE7f38ff82b65f09Fa10e669AC65d"),
		common.HexToAddress("0x5A9e4B0D5Ed9358017f5789314FbBb47Cf74d6C6"),
		common.HexToAddress("0x597E6E36ee36182D8AEDeC34812206CD6533727e"),
		common.HexToAddress("0x51bd1A35e1dC6b9daa82f9890ccEA1aeBcaa16D1"),
		common.HexToAddress("0x4f0B31840Bce10feaC50ac9E6363cc1f455c36Cb"),
		common.HexToAddress("0x4Ec691BEB23cEEf0721693Ea201A83c56E8B2B46"),
		common.HexToAddress("0x4e326ad1C286757dE7cf4b95b3accf05671Ea76E"),
		common.HexToAddress("0x4b2f2Ab47fE076d7763B939DdF5aA0e25AC269bF"),
		common.HexToAddress("0x48C76039D09E556011CE8aE3E0c2e2e247019a78"),
		common.HexToAddress("0x4689145EDB408168F4C10231254B4af994A1B249"),
		common.HexToAddress("0x4641563492052984CDCb046Cc7098719D3f9Db97"),
		common.HexToAddress("0x4448C92dfB560c6F4D2Df371CF12A4fF441a2fe4"),
		common.HexToAddress("0x424d5BC35e7B7953D93C4f2B4b359684A87E647F"),
		common.HexToAddress("0x41D7FDaF014A850D1AE8D14b76bF6A91445647b7"),
		common.HexToAddress("0x412f52Ba4350139b7bf0469781Ac2AB0b5Aa8034"),
		common.HexToAddress("0x40a7f453120eF7AaDC47a0203B75f643e378159F"),
		common.HexToAddress("0x3BE2358a41Bf7D67B2bfA3A21fc4a90ae4b97d9d"),
		common.HexToAddress("0x393a7Bd22eF8939d20F12a16390dfc37Ee06AFB4"),
		common.HexToAddress("0x37fF4076b8cA98f5ce00EcDB5841033A7D231142"),
		common.HexToAddress("0x371F4992D2691f821374656b5de3605De4520109"),
		common.HexToAddress("0x3711754Ae85Cb5BB2fd2106d00aD6A943283b9B0"),
		common.HexToAddress("0x364efBa8EE57C5911069d5c1Cf28947207C441a6"),
		common.HexToAddress("0x3412A3e8E31325a0E173D17D12572656E192Ad92"),
		common.HexToAddress("0x33DFb24d65dfCB5F3b95D5cD2088a870952Eb7B6"),
		common.HexToAddress("0x32329e0422A6222B847d22bF3a0EFEFbDaBd6885"),
		common.HexToAddress("0x2f21d7D75ECD9ac488B44E6d4295A9d7BFCB44Ad"),
		common.HexToAddress("0x2de5674d11e4957Af4Fb4E9Acd64aCc028F1Ec80"),
		common.HexToAddress("0x2A7D0620255fd25502f8c6e9ae7af56393Dc24cd"),
		common.HexToAddress("0x27c439e33E18aeFE4bf85f7fE81b55eb0E43a5D3"),
		common.HexToAddress("0x27009F63C4d01Bb5deEFecdf70B6D08bC8edA720"),
		common.HexToAddress("0x262432667BedA2fAB43376D15A569B87102A8ef5"),
		common.HexToAddress("0x2556d6Cfc3a514f1A215211A5d2594743C744210"),
		common.HexToAddress("0x23A820a708454F94c72F06DEB0A29a7Fa774A4B2"),
		common.HexToAddress("0x2136C5c260B34bd188d01d63A9dc24aE18d22A61"),
		common.HexToAddress("0x1f63bDC4dF28799689119829334B5b584Cae3Fd6"),
		common.HexToAddress("0x1e3dB3B37D75a3F1DaEbA39F433C073472eE784D"),
		common.HexToAddress("0x1D1A53098fea123633f02AA17b4BB7DaA843d3b0"),
		common.HexToAddress("0x1694de7F63667f58f051F375BCD821339B8Cd5A5"),
		common.HexToAddress("0x1662b48A65c4883F2a6C1a758041929a81B5528a"),
		common.HexToAddress("0x15CBd86a925979a69e5eCa71F1baA427bEE7E040"),
		common.HexToAddress("0x14cdC05099BA71a25442939f04F0863ff9f04584"),
		common.HexToAddress("0x13FBd71bFa711BF9Db4B14Df2ca4Dd37f3aeE40c"),
		common.HexToAddress("0x1244013858388Faf893cCf892c98dF4d9F72325c"),
		common.HexToAddress("0x113F701525431A0a4dF11a1f69372Da3e56442AD"),
		common.HexToAddress("0x094de6B4c3a3bdD6a5e041E0afAD3d25aC3Ec502"),
		common.HexToAddress("0x0819C43891678F4Fc86ca321b8D7d7dBa81D1FBB"),
		common.HexToAddress("0x00654Ae79c4Ed252999906304bbF062E6aa4c52a"),
	}

	var validators []*Validator
	for i := 0; i < len(addrs); i++ {
		validators = append(validators, NewValidator(addrs[i], stakeAmt[i], common.Address{}))
	}
	validatorList := NewValidatorList(validators)

	for _, vl := range validatorList.Validators {
		validatorList.CalculateAddressRange(vl.Addr, validatorList.StakeBalance(vl.Addr))
	}

	selfAddrs := GetSelfAddr()
	var count int
	for i := 0; i < 200; i++ {
		randomHash := randomHash()
		consensusValidator := validatorList.RandomValidatorV2(11, randomHash)

		var occurCount int

		for j := 0; j < len(selfAddrs); j++ {
			for _, v := range consensusValidator {
				if v.Hex() == selfAddrs[j].Hex() {
					occurCount++
					break
				}
			}
		}
		if occurCount >= 9 {
			count++
		}
	}
	fmt.Println("===Probability of occurrence===", count, "time", time.Now().Unix())

}

func TestMD5(t *testing.T) {
	for i := 0; i < 10; i++ {
		hash := randomHash()

		addr1 := RandomAddr()
		addr2 := RandomAddr()
		addr3 := RandomAddr()

		var buffer bytes.Buffer

		buffer.WriteString(hash.Hex())
		buffer.WriteString(addr1.Hex())
		buffer.WriteString(addr2.Hex())
		buffer.WriteString(addr3.Hex())
		_ = crypto.Keccak256Hash(buffer.Bytes())
	}
}

func TestFAndSize(t *testing.T) {
	validators := NewValidatorList(nil)
	if validators.F() != -1 {
		t.Fatalf("expected -1 , but got %d", validators.F())
	}
	if validators.Size() != 0 {
		t.Fatalf("expected 0 , but got %d", validators.Size())
	}
}

func TestGetValidatorByAddr(t *testing.T) {
	vals := MockValidator()
	cases := []struct {
		Name     string
		Addr     common.Address
		Expected common.Address
	}{
		{"exist addr", common.HexToAddress("0x091DBBa95B26793515cc9aCB9bEb5124c479f27F"), common.HexToAddress("0x091DBBa95B26793515cc9aCB9bEb5124c479f27F")},
		{"not exist addr", common.HexToAddress("ox1"), common.Address{}},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			validator := vals.GetValidatorByAddr(c.Addr)
			if validator.Addr != c.Expected {
				t.Fatalf("expected %s, but %s got", c.Expected, validator.Addr)
			}
		})
	}
}

func TestGetValidatorByIndex(t *testing.T) {
	vals := MockValidator()
	cases := []struct {
		Name     string
		Index    uint64
		Expected common.Address
	}{
		{"exist addr", 0, common.HexToAddress("0xFfAc4cd934f026dcAF0f9d9EEDDcD9af85D8943e")},
		{"not exist addr", 200, common.Address{}},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			validator := vals.GetByIndex(c.Index)
			if validator.Addr != c.Expected {
				t.Fatalf("expected %s, but %s got", c.Expected, validator.Addr)
			}
		})
	}
}

func TestAddValidator(t *testing.T) {
	addrs := []common.Address{
		common.HexToAddress("0x091DBBa95B26793515cc9aCB9bEb5124c479f27F"),
		common.HexToAddress("0x107837Ea83f8f06533DDd3fC39451Cd0AA8DA8BD"),
	}

	c := 140000
	c2 := 70000

	stakeAmt := []*big.Int{
		big.NewInt(int64(c)),
		big.NewInt(int64(c2)),
	}

	var validators []*Validator
	for i := 0; i < len(addrs); i++ {
		validators = append(validators, NewValidator(addrs[i], stakeAmt[i], common.HexToAddress("0xE2FA892CC5CC268a0cC1d924EC907C796351C645")))
	}
	validatorList := NewValidatorList(validators)

	if len(validatorList.Validators) != 2 {
		t.Fatalf("expected %d, but %d got", 2, len(validatorList.Validators))
	}

	if validatorList.Validators[0].Addr != common.HexToAddress("0x091DBBa95B26793515cc9aCB9bEb5124c479f27F") {
		t.Fatalf("expected %s,but %s got", "0x091DBBa95B26793515cc9aCB9bEb5124c479f27F", validatorList.Validators[0].Addr.Hex())
	}

	for _, v := range validatorList.Validators {
		if v.Proxy != common.HexToAddress("0xE2FA892CC5CC268a0cC1d924EC907C796351C645") {
			t.Fatalf("expected %s, but %s got", "0xE2FA892CC5CC268a0cC1d924EC907C796351C645", v.Proxy.Hex())
		}
	}
}

func TestGetByAddress(t *testing.T) {
	addrs := []common.Address{
		common.HexToAddress("0x091DBBa95B26793515cc9aCB9bEb5124c479f27F"),
		common.HexToAddress("0x107837Ea83f8f06533DDd3fC39451Cd0AA8DA8BD"),
	}

	c := 140000
	c2 := 70000

	stakeAmt := []*big.Int{
		big.NewInt(int64(c)),
		big.NewInt(int64(c2)),
	}

	var validators []*Validator
	for i := 0; i < len(addrs); i++ {
		validators = append(validators, NewValidator(addrs[i], stakeAmt[i], common.HexToAddress("0xE2FA892CC5CC268a0cC1d924EC907C796351C645")))
	}
	validatorList := NewValidatorList(validators)

	cases := []struct {
		Name     string
		Addr     common.Address
		Expected int
	}{
		{"exist addr", common.HexToAddress("0x091DBBa95B26793515cc9aCB9bEb5124c479f27F"), 0},
		{"not exist addr", common.HexToAddress("0x091DBBa95B267df515cc9aCB9bEb5124c479f27F"), -1},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			index := validatorList.GetByAddress(c.Addr)

			if index != c.Expected {
				t.Fatalf("expected %d, but %d got", c.Expected, index)
			}
		})
	}
}

// helper func ============================================================== //
func MockValidator() *ValidatorList {
	addrs := AddrList()

	c := 140000
	c2 := 70000
	stakeAmt := []*big.Int{
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),

		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
	}

	var validators []*Validator
	for i := 0; i < len(addrs); i++ {
		validators = append(validators, NewValidator(addrs[i], stakeAmt[i], common.Address{}))
	}
	validatorList := NewValidatorList(validators)

	for _, vl := range validatorList.Validators {
		validatorList.CalculateAddressRange(vl.Addr, validatorList.StakeBalance(vl.Addr))
	}
	return validatorList
}

func GetAddr(count int) []common.Address {
	var addrs []common.Address
	for i := 0; i < count; i++ {
		priKey, _ := crypto.GenerateKey()
		addrs = append(addrs, crypto.PubkeyToAddress(priKey.PublicKey))
	}
	return addrs
}

func randomHash() common.Hash {
	rand.Seed(time.Now().Local().UnixNano())
	var hash common.Hash
	if n, err := rand.Read(hash[:]); n != common.HashLength || err != nil {
		fmt.Println(err)
	}
	return hash
}
