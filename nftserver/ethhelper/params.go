package ethhelper

import common2 "github.com/nftexchange/nftserver/ethhelper/common"

const (
	erc721Input    = "0x01ffc9a7" + common2.Erc721Interface + "00000000000000000000000000000000000000000000000000000000"
	erc1155Input   = "0x01ffc9a7" + common2.Erc1155Interface + "00000000000000000000000000000000000000000000000000000000"
	nameHash       = "0x06fdde03"
	postUrl        = "http://192.168.1.238:8081/v2"
	newCollections = "/newCollections"

	weth9                = "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"
	weth10               = "0xf4bb2e28688e89fcce3c0580d37d36a7672e8a9f"
	erc721TransferEvent  = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	erc1155TransferEvent = "0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62"
	tokenApproveEvent    = "0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925"
	sigT                 = "0x4a71940655b075316ae19b02457201ed0f719d14f2d20c986b8c16113233e047535d5d1cc4eb293609e79bc60daf622216b190d50a16519d6f826bee05e548051b"
)
const privKey = "564ea566096d3de340fc5ddac98aef672f916624c8b0e4664a908cd2a2d156fe"
const from = "0x077d34394Ed01b3f31fBd9816cF35d4558146066"

type CallParamTemp struct {
	To   string `json:"to"`
	Data string `json:"data"`
}

type Block struct {
	Transactions []Tx   `json:"transactions" `
	Ts           string `json:"timestamp" `
}
type Tx struct {
	Hash  string `json:"hash"`
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
}
type Log struct {
	Address          string   `json:"address"`
	BlockNumber      string   `json:"blockNumber"`
	Topics           []string `json:"topics"`
	TxHash           string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
}
type Receipt struct {
	TransactionHash   string `json:"transactionHash"`
	TransactionIndex  string `json:"transactionIndex"`
	BlockNumber       string `json:"blockNumber"`
	BlockHash         string `json:"blockHash"`
	Logs              []Log  `json:"logs"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	GasUsed           string `json:"gasUsed"`
	ContractAddress   string `json:"contractAddress"`
	LogsBloom         string `json:"logsBloom"`
	Status            string `json:"status"`
}
type CallParam struct {
	From string `json:"from"`
	To   string `json:"to"`
	Data string `json:"data"`
}

type LogFilter struct {
	FromBlock string   `json:"fromBlock"`
	ToBlock   string   `json:"toBlock"`
	Topics    []string `json:"topics"`
}
type RawData struct {
	Data string `json:"data"`
}
type WethTransfer struct {
	From  string `gorm:"column:from;type:varchar(50) " json:"from"`
	To    string `gorm:"column:to;type:varchar(50) " json:"to"`
	Value string `gorm:"column:value;type:varchar(50) " json:"value"`
}
