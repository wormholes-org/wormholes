module github.com/nftexchange/nftserver/controllers

go 1.16

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/nftexchange/nftserver/models v0.0.0
)

replace (
	github.com/nftexchange/nftserver/ethhelper v0.0.0 => ../ethhelper
	github.com/nftexchange/nftserver/ethhelper/common v0.0.0 => ../ethhelper/common
	github.com/nftexchange/nftserver/ethhelper/database v0.0.0 => ../ethhelper/database
	github.com/nftexchange/nftserver/models v0.0.0 => ../models
)
