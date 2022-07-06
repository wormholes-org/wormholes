module github.com/nftexchange/nftserver/controllers/nftexchangev2

go 1.16

require (
	github.com/beego/beego/v2 v2.0.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/nftexchange/nftserver/common/signature v0.0.0
	github.com/nftexchange/nftserver/controllers v0.0.0
	github.com/nftexchange/nftserver/models v0.0.0
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	gorm.io/gorm v1.21.15
)

replace (
	github.com/nftexchange/nftserver/common/signature v0.0.0 => ../../common/signature
	github.com/nftexchange/nftserver/controllers v0.0.0 => ../../controllers
	github.com/nftexchange/nftserver/ethhelper v0.0.0 => ../../ethhelper
	github.com/nftexchange/nftserver/ethhelper/common v0.0.0 => ../../ethhelper/common
	github.com/nftexchange/nftserver/ethhelper/database v0.0.0 => ../../ethhelper/database
	github.com/nftexchange/nftserver/models v0.0.0 => ../../models

)
