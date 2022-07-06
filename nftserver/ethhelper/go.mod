module github.com/nftexchange/nftserver/ethhelper

go 1.16

require (
	github.com/ethereum/go-ethereum v1.10.9
	github.com/nftexchange/nftserver/ethhelper/common v0.0.0
	github.com/nftexchange/nftserver/ethhelper/database v0.0.0
)

replace (
	github.com/nftexchange/nftserver/ethhelper v0.0.0 => ../ethhelper
	github.com/nftexchange/nftserver/ethhelper/common v0.0.0 => ../ethhelper/common
	github.com/nftexchange/nftserver/ethhelper/database v0.0.0 => ../ethhelper/database
)
