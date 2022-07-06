module github.com/nftexchange/nftserver/ethhelper/database

go 1.16

require (
	github.com/jinzhu/gorm v1.9.16
	github.com/nftexchange/nftserver/ethhelper/common v0.0.0
)

replace github.com/nftexchange/nftserver/ethhelper/common v0.0.0 => ../common
