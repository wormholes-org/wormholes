module github.com/nftexchange/nftserver/common/contracts

go 1.15

replace github.com/nftexchange/nftserver/controllers v0.0.0 => ../../controllers

require github.com/ethereum/go-ethereum v1.10.9

replace github.com/nftexchange/nftserver/ethhelper v0.0.0 => ../../ethhelper

require (
	github.com/nftexchange/nftserver/ethhelper v0.0.0
	github.com/nftexchange/nftserver/models v0.0.0
)

replace github.com/nftexchange/nftserver/models v0.0.0 => ../../models

replace (
	github.com/nftexchange/nftserver/ethhelper/common v0.0.0 => ../../ethhelper/common
	github.com/nftexchange/nftserver/ethhelper/database v0.0.0 => ../../ethhelper/database
)
