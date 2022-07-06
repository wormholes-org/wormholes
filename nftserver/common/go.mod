module github.com/nftexchange/nftserver/common

go 1.15

require github.com/nftexchange/nftserver/controllers v0.0.0
replace github.com/nftexchange/nftserver/controllers v0.0.0 => ../controllers
require github.com/ethereum/go-ethereum v1.10.8
replace github.com/nftexchange/nftserver/ethhelper v0.0.0 => ../ethhelper
require github.com/nftexchange/nftserver/models v0.0.0
replace github.com/nftexchange/nftserver/models v0.0.0 => ../models
replace github.com/nftexchange/nftserver/ethhelper/common v0.0.0 => ../ethhelper/common
replace github.com/nftexchange/nftserver/ethhelper/database v0.0.0 => ../ethhelper/database