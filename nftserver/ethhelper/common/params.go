package common

const (
	//mainPoint = "https://mainnet.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161"
	MainPoint = "http://192.168.1.235:8546"
	//MainPoint        = "https://rinkeby.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161"
	InfraPoint       = "https://rinkeby.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161"
	erc721           = "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"
	erc1155          = "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"
	Erc721Interface  = "80ac58cd"
	Erc1155Interface = "d9b67a26"
)
const (
	WETH          = "0xf4bb2e28688e89fcce3c0580d37d36a7672e8a9f"
	BalanceOfHash = "0x70a08231"
	AllowanceHash = "0xdd62ed3e"
)
const (
	Admin         = "0x56c971ebBC0cD7Ba1f977340140297C0B48b7955"
	AdminListHash = "0x0f560cd7"
	AddIndexHash  = "0x151b01f9"
)

const (
	TradeCore   = "0xcec263F075945d904bE747266B4750Aab87B4813"
	AuctionHash = "0x151b01f9"
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
	Data             string   `json:"data"`
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
