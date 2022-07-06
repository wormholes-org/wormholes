module github.com/nftexchange/nftserver/controllers/nftexchangev1

go 1.16

require (
	github.com/beego/beego/v2 v2.0.1
	github.com/nftexchange/nftserver/controllers v0.0.0
	github.com/nftexchange/nftserver/models v0.0.0
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
)

replace (
	github.com/nftexchange/nftserver/common/signature v0.0.0 => ../../common/signature
	github.com/nftexchange/nftserver/controllers v0.0.0 => ../../controllers
	github.com/nftexchange/nftserver/ethhelper v0.0.0 => ../../ethhelper
	github.com/nftexchange/nftserver/ethhelper/common v0.0.0 => ../../ethhelper/common
	github.com/nftexchange/nftserver/ethhelper/database v0.0.0 => ../../ethhelper/database
	github.com/nftexchange/nftserver/models v0.0.0 => ../../models
)
