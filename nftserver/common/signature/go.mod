module github.com/nftexchange/nftserver/common/signature

go 1.16

require (
	github.com/ethereum/go-ethereum v1.10.8
	github.com/nftexchange/nftserver/ethhelper v0.0.0
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2
)

replace (
	github.com/nftexchange/nftserver/ethhelper v0.0.0 => ../../ethhelper
	github.com/nftexchange/nftserver/ethhelper/common v0.0.0 => ../../ethhelper/common
    github.com/nftexchange/nftserver/ethhelper/database v0.0.0 => ../../ethhelper/database

)