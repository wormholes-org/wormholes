module github.com/ethereum/go-ethereum

go 1.15

require (
	github.com/Azure/azure-storage-blob-go v0.7.0
	github.com/Azure/go-autorest/autorest/adal v0.9.21 // indirect
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/VictoriaMetrics/fastcache v1.6.0
	github.com/aws/aws-sdk-go-v2 v1.2.0
	github.com/aws/aws-sdk-go-v2/config v1.1.1
	github.com/aws/aws-sdk-go-v2/credentials v1.1.1
	github.com/aws/aws-sdk-go-v2/service/route53 v1.1.1
	github.com/btcsuite/btcd v0.22.0-beta
	github.com/cespare/cp v0.1.0
	github.com/cloudflare/cloudflare-go v0.14.0
	github.com/consensys/gnark-crypto v0.4.1-0.20210426202927-39ac3d4b3f1f
	github.com/coreos/go-systemd/v22 v22.3.2
	github.com/davecgh/go-spew v1.1.1
	github.com/deckarep/golang-set v0.0.0-20180603214616-504e848d77ea
	github.com/dlclark/regexp2 v1.7.0 // indirect
	github.com/docker/docker v1.4.2-0.20180625184442-8e610b2b55bf
	github.com/dop251/goja v0.0.0-20200721192441-a695b0cdd498
	github.com/edsrzf/mmap-go v1.0.0
	github.com/fatih/color v1.9.0
	github.com/fjl/memsize v0.0.0-20190710130421-bcb5799ab5e5
	github.com/gballet/go-libpcsclite v0.0.0-20190607065134-2772fd86a8ff
	github.com/go-redis/redis/v7 v7.4.0
	github.com/go-sourcemap/sourcemap v2.1.3+incompatible // indirect
	github.com/go-stack/stack v1.8.0
	github.com/golang/protobuf v1.5.2
	github.com/golang/snappy v0.0.4
	github.com/google/gofuzz v1.1.1-0.20200604201612-c04b05f3adfa
	github.com/google/uuid v1.3.0
	github.com/gorilla/websocket v1.4.2
	github.com/gotestyourself/gotestyourself v1.4.0 // indirect
	github.com/graph-gophers/graphql-go v0.0.0-20201113091052-beb923fada29
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/golang-lru v0.5.5-0.20210104140557-80c98217689d
	github.com/holiman/bloomfilter/v2 v2.0.3
	github.com/holiman/uint256 v1.2.0
	github.com/huin/goupnp v1.0.2
	github.com/influxdata/influxdb v1.8.3
	github.com/influxdata/influxdb-client-go/v2 v2.4.0
	github.com/ipfs/go-bitswap v0.5.1 // indirect
	github.com/ipfs/go-blockservice v0.2.1 // indirect
	github.com/ipfs/go-cid v0.1.0
	github.com/ipfs/go-datastore v0.5.0 // indirect
	github.com/ipfs/go-ipfs v0.10.0
	github.com/ipfs/go-ipfs-blockstore v0.2.1 // indirect
	github.com/ipfs/go-ipfs-cmds v0.6.0
	github.com/ipfs/go-ipfs-config v0.17.0
	github.com/ipfs/go-ipfs-exchange-interface v0.1.0 // indirect
	github.com/ipfs/go-ipfs-exchange-offline v0.1.1 // indirect
	github.com/ipfs/go-ipfs-files v0.0.9
	github.com/ipfs/go-ipfs-routing v0.2.1 // indirect
	github.com/ipfs/go-ipfs-util v0.0.2
	github.com/ipfs/go-ipld-format v0.2.0
	github.com/ipfs/go-log v1.0.5
	github.com/ipfs/go-merkledag v0.5.1
	github.com/ipfs/go-metrics-prometheus v0.0.2
	github.com/ipfs/go-path v0.1.2
	github.com/ipfs/go-pinning-service-http-client v0.1.0
	github.com/ipfs/go-unixfs v0.2.5
	github.com/ipfs/interface-go-ipfs-core v0.5.1
	github.com/ipld/go-car v0.3.2 // indirect
	github.com/jackpal/go-nat-pmp v1.0.2
	github.com/jbenet/goprocess v0.1.4
	github.com/jedisct1/go-minisign v0.0.0-20190909160543-45766022959e
	github.com/julienschmidt/httprouter v1.3.0
	github.com/karalabe/usb v0.0.0-20190919080040-51dc0efba356
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/libp2p/go-libp2p-core v0.9.0
	github.com/libp2p/go-libp2p-loggables v0.1.0
	github.com/libp2p/go-socket-activation v0.1.0
	github.com/mattn/go-colorable v0.1.8
	github.com/mattn/go-isatty v0.0.13
	github.com/multiformats/go-multiaddr v0.4.1
	github.com/multiformats/go-multiaddr-dns v0.3.1
	github.com/naoina/go-stringutil v0.1.0 // indirect
	github.com/naoina/toml v0.1.2-0.20170918210437-9fafd6967416
	github.com/olekukonko/tablewriter v0.0.5
	github.com/peterh/liner v1.1.1-0.20190123174540-a2c9a5303de7
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/tsdb v0.7.1
	github.com/rjeczalik/notify v0.9.1
	github.com/rs/cors v1.7.0
	github.com/shirou/gopsutil v3.21.4-0.20210419000835-c7a38de76ee5+incompatible
	github.com/status-im/keycard-go v0.0.0-20190316090335-8537d3370df4
	github.com/stretchr/testify v1.7.0
	github.com/syndtr/goleveldb v1.0.1-0.20210819022825-2ae1ddf74ef7
	github.com/tklauser/go-sysconf v0.3.10 // indirect
	github.com/tyler-smith/go-bip39 v1.0.1-0.20181017060643-dbb3b84ba2ef
	go.uber.org/zap v1.19.0
	golang.org/x/crypto v0.0.0-20220722155217-630584e8d5aa
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20220128215802-99c3d69c2c27
	golang.org/x/text v0.3.7
	golang.org/x/time v0.0.0-20210220033141-f8bda1e9f3ba
	gopkg.in/karalabe/cookiejar.v2 v2.0.0-20150724131613-8dcd6a7f4951
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce
	gopkg.in/olebedev/go-duktape.v3 v3.0.0-20200619000410-60c24ae608a6
	gopkg.in/urfave/cli.v1 v1.20.0
	gotest.tools v1.4.0 // indirect
)

//require github.com/ethereum/go-ethereum/go-ipfs v0.0.0

//require github.com/ethereum/go-ethereum/go-ipfs/cmd/ipfs v0.0.0
//replace github.com/ethereum/go-ethereum/go-ipfs v0.0.0 => ./go-ipfs

replace (
	github.com/ipfs/go-bitswap v0.5.1 => github.com/ipfs/go-bitswap v0.4.0
	github.com/ipfs/go-blockservice v0.2.1 => github.com/ipfs/go-blockservice v0.1.7
	github.com/ipfs/go-cid v0.1.0 => github.com/ipfs/go-cid v0.0.7
	github.com/ipfs/go-datastore v0.5.0 => github.com/ipfs/go-datastore v0.4.6
	github.com/ipfs/go-ipfs-blockstore v0.2.1 => github.com/ipfs/go-ipfs-blockstore v0.1.6
	github.com/ipfs/go-ipfs-config v0.17.0 => github.com/ipfs/go-ipfs-config v0.16.0

	github.com/ipfs/go-ipfs-exchange-interface v0.1.0 => github.com/ipfs/go-ipfs-exchange-interface v0.0.1
	github.com/ipfs/go-ipfs-exchange-offline v0.1.1 => github.com/ipfs/go-ipfs-exchange-offline v0.0.1
	github.com/ipfs/go-ipfs-routing v0.2.1 => github.com/ipfs/go-ipfs-routing v0.1.0
	github.com/ipfs/go-merkledag v0.5.1 => github.com/ipfs/go-merkledag v0.4.0
	github.com/multiformats/go-multiaddr v0.4.1 => github.com/multiformats/go-multiaddr v0.4.0
	//github.com/syndtr/goleveldb v1.0.1-0.20210819022825-2ae1ddf74ef7 => github.com/syndtr/goleveldb v1.0.0
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 => golang.org/x/crypto v0.0.0-20210813211128-0a44fdfbc16e
)
