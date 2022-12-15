// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package params

import "github.com/ethereum/go-ethereum/common"

// MainnetBootnodes are the enode URLs of the P2P bootstrap nodes running on
// the main Ethereum network.
var MainnetBootnodes = []string{
	//Ethereum Foundation Go Bootnodes
	//"enode://d860a01f9722d78051619d1e2351aba3f43f943f6f00718d1b9baa4101932a1f5011f16bb2b1bb35db20d6fe28fa0bf09636d26a87d31de9ec6203eeedb1f666@18.138.108.67:30303",   // bootnode-aws-ap-southeast-1-001
	//"enode://22a8232c3abc76a16ae9d6c3b164f98775fe226f0917b0ca871128a74a8e9630b458460865bab457221f1d448dd9791d24c4e5d88786180ac185df813a68d4de@3.209.45.79:30303",     // bootnode-aws-us-east-1-001
	//"enode://ca6de62fce278f96aea6ec5a2daadb877e51651247cb96ee310a318def462913b653963c155a0ef6c7d50048bba6e6cea881130857413d9f50a621546b590758@34.255.23.113:30303",   // bootnode-aws-eu-west-1-001
	//"enode://279944d8dcd428dffaa7436f25ca0ca43ae19e7bcf94a8fb7d1641651f92d121e972ac2e8f381414b80cc8e5555811c2ec6e1a99bb009b3f53c4c69923e11bd8@35.158.244.151:30303",  // bootnode-aws-eu-central-1-001
	//"enode://8499da03c47d637b20eee24eec3c356c9a2e6148d6fe25ca195c7949ab8ec2c03e3556126b0d7ed644675e78c4318b08691b7b57de10e5f0d40d05b09238fa0a@52.187.207.27:30303",   // bootnode-azure-australiaeast-001
	//"enode://103858bdb88756c71f15e9b5e09b56dc1be52f0a5021d46301dbbfb7e130029cc9d0d6f73f693bc29b665770fff7da4d34f3c6379fe12721b5d7a0bcb5ca1fc1@191.234.162.198:30303", // bootnode-azure-brazilsouth-001
	//"enode://715171f50508aba88aecd1250af392a45a330af91d7b90701c436b618c86aaa1589c9184561907bebbb56439b8f8787bc01f49a7c77276c58c1b09822d75e8e8@52.231.165.108:30303",  // bootnode-azure-koreasouth-001
	//"enode://5d6d7cd20d6da4bb83a1d28cadb5d409b64edf314c0335df658c1a54e32c7c4a7ab7823d57c39b6a757556e68ff1df17c748b698544a55cb488b52479a92b60f@104.42.217.25:30303",   // bootnode-azure-westus-001
}

// RopstenBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Ropsten test network.
var RopstenBootnodes = []string{
	//"enode://30b7ab30a01c124a6cceca36863ece12c4f5fa68e3ba9b0b51407ccc002eeed3b3102d20a88f1c1d3c3154e2449317b8ef95090e77b312d5cc39354f86d5d606@52.176.7.10:30303",    // US-Azure geth
	//"enode://865a63255b3bb68023b6bffd5095118fcc13e79dcf014fe4e47e065c350c7cc72af2e53eff895f11ba1bbb6a2b33271c1116ee870f266618eadfc2e78aa7349c@52.176.100.77:30303",  // US-Azure parity
	//"enode://6332792c4a00e3e4ee0926ed89e0d27ef985424d97b6a45bf0f23e51f0dcb5e66b875777506458aea7af6f9e4ffb69f43f3778ee73c81ed9d34c51c4b16b0b0f@52.232.243.152:30303", // Parity
	//"enode://94c15d1b9e2fe7ce56e458b9a3b672ef11894ddedd0c6f247e0f1d3487f52b66208fb4aeb8179fce6e3a749ea93ed147c37976d67af557508d199d9594c35f09@192.81.208.223:30303", // @gpip
}

// RinkebyBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Rinkeby test network.
var RinkebyBootnodes = []string{
	//"enode://a24ac7c5484ef4ed0c5eb2d36620ba4e4aa13b8c84684e1b4aab0cebea2ae45cb4d375b77eab56516d34bfbd3c1a833fc51296ff084b770b94fb9028c4d25ccf@52.169.42.101:30303", // IE
	//"enode://343149e4feefa15d882d9fe4ac7d88f885bd05ebb735e547f12e12080a9fa07c8014ca6fd7f373123488102fe5e34111f8509cf0b7de3f5b44339c9f25e87cb8@52.3.158.184:30303",  // INFURA
	//"enode://b6b28890b006743680c52e64e0d16db57f28124885595fa03a562be1d2bf0f3a1da297d56b13da25fb992888fd556d4c1a27b1f39d531bde7de1921c90061cc6@159.89.28.211:30303", // AKASHA
}

// GoerliBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Görli test network.
var GoerliBootnodes = []string{
	////Upstream bootnodes
	//"enode://011f758e6552d105183b1761c5e2dea0111bc20fd5f6422bc7f91e0fabbec9a6595caf6239b37feb773dddd3f87240d99d859431891e4a642cf2a0a9e6cbb98a@51.141.78.53:30303",
	//"enode://176b9417f511d05b6b2cf3e34b756cf0a7096b3094572a8f6ef4cdcb9d1f9d00683bf0f83347eebdf3b81c3521c2332086d9592802230bf528eaf606a1d9677b@13.93.54.137:30303",
	//"enode://46add44b9f13965f7b9875ac6b85f016f341012d84f975377573800a863526f4da19ae2c620ec73d11591fa9510e992ecc03ad0751f53cc02f7c7ed6d55c7291@94.237.54.114:30313",
	//"enode://b5948a2d3e9d486c4d75bf32713221c2bd6cf86463302339299bd227dc2e276cd5a1c7ca4f43a0e9122fe9af884efed563bd2a1fd28661f3b5f5ad7bf1de5949@18.218.250.66:30303",
	//
	//// Ethereum Foundation bootnode
	//"enode://a61215641fb8714a373c80edbfa0ea8878243193f57c96eeb44d0bc019ef295abd4e044fd619bfc4c59731a73fb79afe84e9ab6da0c743ceb479cbb6d263fa91@3.11.147.67:30303",
	//
	//// Goerli Initiative bootnodes
	//"enode://a869b02cec167211fb4815a82941db2e7ed2936fd90e78619c53eb17753fcf0207463e3419c264e2a1dd8786de0df7e68cf99571ab8aeb7c4e51367ef186b1dd@51.15.116.226:30303",
	//"enode://807b37ee4816ecf407e9112224494b74dd5933625f655962d892f2f0f02d7fbbb3e2a94cf87a96609526f30c998fd71e93e2f53015c558ffc8b03eceaf30ee33@51.15.119.157:30303",
	//"enode://a59e33ccd2b3e52d578f1fbd70c6f9babda2650f0760d6ff3b37742fdcdfdb3defba5d56d315b40c46b70198c7621e63ffa3f987389c7118634b0fefbbdfa7fd@51.15.119.157:40303",
}

var V5Bootnodes = []string{
	//// Teku team's bootnode
	//"enr:-KG4QOtcP9X1FbIMOe17QNMKqDxCpm14jcX5tiOE4_TyMrFqbmhPZHK_ZPG2Gxb1GE2xdtodOfx9-cgvNtxnRyHEmC0ghGV0aDKQ9aX9QgAAAAD__________4JpZIJ2NIJpcIQDE8KdiXNlY3AyNTZrMaEDhpehBDbZjM_L9ek699Y7vhUJ-eAdMyQW_Fil522Y0fODdGNwgiMog3VkcIIjKA",
	//"enr:-KG4QDyytgmE4f7AnvW-ZaUOIi9i79qX4JwjRAiXBZCU65wOfBu-3Nb5I7b_Rmg3KCOcZM_C3y5pg7EBU5XGrcLTduQEhGV0aDKQ9aX9QgAAAAD__________4JpZIJ2NIJpcIQ2_DUbiXNlY3AyNTZrMaEDKnz_-ps3UUOfHWVYaskI5kWYO_vtYMGYCQRAR3gHDouDdGNwgiMog3VkcIIjKA",
	//// Prylab team's bootnodes
	//"enr:-Ku4QImhMc1z8yCiNJ1TyUxdcfNucje3BGwEHzodEZUan8PherEo4sF7pPHPSIB1NNuSg5fZy7qFsjmUKs2ea1Whi0EBh2F0dG5ldHOIAAAAAAAAAACEZXRoMpD1pf1CAAAAAP__________gmlkgnY0gmlwhBLf22SJc2VjcDI1NmsxoQOVphkDqal4QzPMksc5wnpuC3gvSC8AfbFOnZY_On34wIN1ZHCCIyg",
	//"enr:-Ku4QP2xDnEtUXIjzJ_DhlCRN9SN99RYQPJL92TMlSv7U5C1YnYLjwOQHgZIUXw6c-BvRg2Yc2QsZxxoS_pPRVe0yK8Bh2F0dG5ldHOIAAAAAAAAAACEZXRoMpD1pf1CAAAAAP__________gmlkgnY0gmlwhBLf22SJc2VjcDI1NmsxoQMeFF5GrS7UZpAH2Ly84aLK-TyvH-dRo0JM1i8yygH50YN1ZHCCJxA",
	//"enr:-Ku4QPp9z1W4tAO8Ber_NQierYaOStqhDqQdOPY3bB3jDgkjcbk6YrEnVYIiCBbTxuar3CzS528d2iE7TdJsrL-dEKoBh2F0dG5ldHOIAAAAAAAAAACEZXRoMpD1pf1CAAAAAP__________gmlkgnY0gmlwhBLf22SJc2VjcDI1NmsxoQMw5fqqkw2hHC4F5HZZDPsNmPdB1Gi8JPQK7pRc9XHh-oN1ZHCCKvg",
	//// Lighthouse team's bootnodes
	//"enr:-IS4QLkKqDMy_ExrpOEWa59NiClemOnor-krjp4qoeZwIw2QduPC-q7Kz4u1IOWf3DDbdxqQIgC4fejavBOuUPy-HE4BgmlkgnY0gmlwhCLzAHqJc2VjcDI1NmsxoQLQSJfEAHZApkm5edTCZ_4qps_1k_ub2CxHFxi-gr2JMIN1ZHCCIyg",
	//"enr:-IS4QDAyibHCzYZmIYZCjXwU9BqpotWmv2BsFlIq1V31BwDDMJPFEbox1ijT5c2Ou3kvieOKejxuaCqIcjxBjJ_3j_cBgmlkgnY0gmlwhAMaHiCJc2VjcDI1NmsxoQJIdpj_foZ02MXz4It8xKD7yUHTBx7lVFn3oeRP21KRV4N1ZHCCIyg",
	//// EF bootnodes
	//"enr:-Ku4QHqVeJ8PPICcWk1vSn_XcSkjOkNiTg6Fmii5j6vUQgvzMc9L1goFnLKgXqBJspJjIsB91LTOleFmyWWrFVATGngBh2F0dG5ldHOIAAAAAAAAAACEZXRoMpC1MD8qAAAAAP__________gmlkgnY0gmlwhAMRHkWJc2VjcDI1NmsxoQKLVXFOhp2uX6jeT0DvvDpPcU8FWMjQdR4wMuORMhpX24N1ZHCCIyg",
	//"enr:-Ku4QG-2_Md3sZIAUebGYT6g0SMskIml77l6yR-M_JXc-UdNHCmHQeOiMLbylPejyJsdAPsTHJyjJB2sYGDLe0dn8uYBh2F0dG5ldHOIAAAAAAAAAACEZXRoMpC1MD8qAAAAAP__________gmlkgnY0gmlwhBLY-NyJc2VjcDI1NmsxoQORcM6e19T1T9gi7jxEZjk_sjVLGFscUNqAY9obgZaxbIN1ZHCCIyg",
	//"enr:-Ku4QPn5eVhcoF1opaFEvg1b6JNFD2rqVkHQ8HApOKK61OIcIXD127bKWgAtbwI7pnxx6cDyk_nI88TrZKQaGMZj0q0Bh2F0dG5ldHOIAAAAAAAAAACEZXRoMpC1MD8qAAAAAP__________gmlkgnY0gmlwhDayLMaJc2VjcDI1NmsxoQK2sBOLGcUb4AwuYzFuAVCaNHA-dy24UuEKkeFNgCVCsIN1ZHCCIyg",
	//"enr:-Ku4QEWzdnVtXc2Q0ZVigfCGggOVB2Vc1ZCPEc6j21NIFLODSJbvNaef1g4PxhPwl_3kax86YPheFUSLXPRs98vvYsoBh2F0dG5ldHOIAAAAAAAAAACEZXRoMpC1MD8qAAAAAP__________gmlkgnY0gmlwhDZBrP2Jc2VjcDI1NmsxoQM6jr8Rb1ktLEsVcKAPa08wCsKUmvoQ8khiOl_SLozf9IN1ZHCCIyg",
}

//var TestnetBootnodes = []string{
//	"enode://af73c8b5a0e8ce5ac563e784051aec62ab16f72b900a384942648380cdf11920512d81fe9ae792aaa90f391505a5c2af1db560d4054da3d230d2faf3b8f35f24@127.0.0.1:30320",
//	"enode://28faa3402ab50153aa5ff2ba3cd7ebbbf96475df2a3f809ace43cf50ba91ed4e676aa94e69c42ea7fab89e5260a6873f7a7781edbd6c18c22609778dd64e14d6@127.0.0.1:30321",
//	"enode://412a18472f07bbb25092451615d5dfea35c9560a6f1cfbe7e8e02e5914747c2a2cc670f82c9c8ac7f162219527c2305c99db1b1b643b9c823023f46fb9bf619d@127.0.0.1:30322",
//	"enode://da922f8adefacfa10096e561cbcb5b310cb34ddd8ddccd880cf77db1379701e0642b32105d5edf54f8cdd279efd0cd54ea15fdbf9e53cefd099a5832404831eb@127.0.0.1:30323",
//	"enode://0266f251db2cb1ae455e58e4c77ec243e2742aefe5bdef90bdb4778a189a809710a9e006e25cf084f3709c372a64f700a8253935bdeca94894f954ec5eb7c363@127.0.0.1:30324",
//	"enode://241f38eb1a40a38941fee481380af8b3dbfff45c96ea33917c7b228d2c1bd867e0880b80daa170bf0b9ef1e52a5408f176b89d8a48684f939d6193d2118308e4@127.0.0.1:30325",
//	"enode://c16b87b005cab802a1a579d3599c2f8d5b9ea8920b53cfe5eec0ef934279e16b55f30e6725200b975a4401146d3a780877304e39d01866ac690063766f6459fd@127.0.0.1:30326",
//	"enode://f06d5b328b8d122c5b8b010359c5fecc7dbb8221c3a8c6efa05010cb81eb65523a6fadbf8fc3e3bbb1b0443b12aa9cf42aa01e989e5ed76da04a2a54a6758c2c@127.0.0.1:30327",
//	"enode://0ba07bd762c6f9c4d08ed10c1e72f18a40116c102f062149c4f796ac7e38757e2cad12b8ccb46be9e0e3686b85b5f757cdbc26ec92d8cc8ef1fe463283c89718@127.0.0.1:30328",
//	"enode://44bccb90d4b1eb92187f58967505e001c7768c66a8026cdb96aa85f923587fdb2b10afc410ffcdbd9e55052237fc735e1ef9a5e9813f8dea3e84ba24181aace5@127.0.0.1:30329",
//	"enode://d3ba34ef68876c2661a782d7fd5402b9629eaa64f7b45ec59a6c03eb403fd79808b9e95dd91d8004f19a6d837f34805ba0881b58cf7232322098781e9643ac1e@127.0.0.1:30330",
//	"enode://9bd5378222f34e869582c712cd42edee6aa60a9dc0f6850ef2f0d1a6893490e34697fac2330eaa7f343230d2e014a073b5e9a75f20eab1d8337075018183abc7@127.0.0.1:30331",
//	"enode://635a0a547fe7fbfcd498774a46e55e849cc9c1fc1e6e90059cc718392f4a896a929a67e30bec7fafc99f4ee3468667faffd66b090a0f1445013870fdcb6b0c14@127.0.0.1:30332",
//	"enode://3da5539593717cfbf9147f02183346a72dfe1702e14fc1fde96df07c6e4cfad06d25fd0289c0ee53b33df02242740226dc9d2e99335044941f3c343fd1da0684@127.0.0.1:30333",
//	"enode://1ae41110c2c30cd9995abcb17f6255ca25f62b89e35e9a820cf274b2e6651e7a133a8b0e85d2ff79b1d594915027900b6fac0ff464f96f947e55d8ed597a7eb4@127.0.0.1:30334",
//}
var TestnetBootnodes = []string{
	"enode://8dfa963b177833d7d181324e9d260d88338274fb7e2855975efc6a9adddc9c335c01566b241c2fac2281cffd07055a8e63d5ef9f5110f1fa29e640dc529618e5@127.0.0.1:30320",
	"enode://0b8f42d969ea2c9001ee83622e5db38e7df7b6036ec3cf63f9d62b2a7aa99558426a2aad8248e0ceed95b6921d11432ce804dd3c038e59f125b3fbdda214afe2@127.0.0.1:30321",
	"enode://8c2d79a0a0a8dc1debcc00ce05563eea9ebb341c7e865518eeee3f341f2a41e8b35ecda518ee651de7438977b671378b8a4df0f632d31a813542368269144367@127.0.0.1:30322",
	"enode://a145065c04ed7e4898f0c724c946c29496cb3a746bdaed88a6885936f07c48429dc93605d89ad65d708eb2a2ddbed793d80ff3e4c68be84ae69e6afc98a502d0@127.0.0.1:30323",
	"enode://2c91e6771802fcea57c4ec67fc1ab5a2e13a006cbd206bb852e7b7a02d3a0dfb17fb50c1cab91998a55ecb02a73dabe90f147ea8522bda40e81ec1058fe5b218@127.0.0.1:30324",
	"enode://fcaee3896fc898385a7b4e564236d807e160e491d35863fdc7c80157361726b9b27694d86bf38a8b0cfd9abb0720b2bda5a2f4386eb7525f2975cdecb62132c8@127.0.0.1:30325",
	"enode://3061ee83846c47170cf2539209ba6d4849b0490224dea4311785fc3322db9eabb6de4684fab4d6d5331b9707a6f64b6882dce9545c03678646a8f5b896f4ef25@127.0.0.1:30326",
	"enode://90dc742cba1bfef400ddcfd7aa9f400419cce740b99284bd3e04eb58a2c40b9883f8c7eab23c615f5c97ad8349b65634cade63f3831f2a843c65c9bfdc86b340@127.0.0.1:30327",
	"enode://081b117da0d60514c2c247bfe47ee9c73ad3bb547890c722c7e1e85ad2e01204361f7400b36c77d0a17b3deef7c213c84d91e0edb8f7f64e35f76ab58a487fd4@127.0.0.1:30328",
	"enode://17d82c989a517ee065081577039da71bdbbd07fc326e1f53a86d1ec5c9c420bb2ba83dc87e0cf3206af77ceced8ae545bcf76462e66f3e3c94fcbb0836a552ff@127.0.0.1:30329",
	"enode://370d9a40b73e6f1f0795214cfc4f7703d2231998ecd9f26d9038fe266a22b63e54e83ae223b8282ad0dc917fb8067b42065ac83010c218bf689df0dc4f3738fc@127.0.0.1:30330",
	"enode://9574dfb9e130e46a34e574da9b6a360cbfc9f254f19887c1d3f49e0b4847d6fa42d3440ea85983824a2d13f3d09c4a735e8c73d212f37eb90bda64f6eb115b6a@127.0.0.1:30331",
	"enode://775a0973bc39ac999734eed0edefa5f325ac0aaf4e1bd1ba2e65a160f7018abae7f24b4199ef3a8368f3f1cb2fd2ee9f07ddd72903fb6f22f37f47a7d5b46277@127.0.0.1:30332",
	"enode://422eaddae6804ee2fa1dc97fe954baab018e8151f7e794bf8d8783218353bcc8928a18d2d46a72a309d14e09c278647be6db2f9e26a8d97da56e44a30c685ccc@127.0.0.1:30333",
	"enode://a9939ea4e177ffd54f119c7241a72ed8273bbd199ed8972d2a4461641d8f6a50da6c9d9d9587c8c31790fd3f3b3dccb9d51f9bd9c05cfaa819a3061170893b3d@127.0.0.1:30334",
	"enode://03cd10c01f1b7a4979c9f9cd042ee6d3c73283966911c43ed7d6b1c420b2858e51714943f481eb3e52af16b875f9741898a7cec7da251787b9bb34c312edaf8e@127.0.0.1:30335",
}

//var DevnetBootnodes = []string{
//	"enode://8dfa963b177833d7d181324e9d260d88338274fb7e2855975efc6a9adddc9c335c01566b241c2fac2281cffd07055a8e63d5ef9f5110f1fa29e640dc529618e5@129.226.154.223:30320",
//	"enode://ae51de0a0a0deb2ff61c967b2cac74687275b0b77e82d5ecaa5bd1270057057923e217a357cf9a961239c932d53ab1a88327f95fd037755185e8212e7ee97770@54.177.66.169:30320",
//	"enode://8c2d79a0a0a8dc1debcc00ce05563eea9ebb341c7e865518eeee3f341f2a41e8b35ecda518ee651de7438977b671378b8a4df0f632d31a813542368269144367@18.169.133.50:30320",
//	"enode://f7e85ded754f75331568b2351382f7b20e5b9b51982311b032c6ee39a71327807fe268ee488b28795b70fa185509eebd33c6325dd6319128a2d38da24c9d9506@54.155.55.37:30320",
//	"enode://2c91e6771802fcea57c4ec67fc1ab5a2e13a006cbd206bb852e7b7a02d3a0dfb17fb50c1cab91998a55ecb02a73dabe90f147ea8522bda40e81ec1058fe5b218@15.236.202.37:30320",
//	"enode://2dbcdc8addd533ac488d1c2ae3456550699c675f8db0f9c04dd308050c03eb327e421cc3fd39e84fa7feca167f7d15f4377555dae22f2afd10658dffab9bebbc@13.211.121.131:30320",
//	"enode://3061ee83846c47170cf2539209ba6d4849b0490224dea4311785fc3322db9eabb6de4684fab4d6d5331b9707a6f64b6882dce9545c03678646a8f5b896f4ef25@15.160.94.149:30320",
//	"enode://ba3834fd7b0e9d5b87ed15d1e4291528adff0b5a58c9ae54b9a6d3bb47b57d9b40f039d9195a63f5ba8176f395e9d09ca8e109e0a86a30b5091843a07784edb7@15.152.100.177:30320",
//	"enode://081b117da0d60514c2c247bfe47ee9c73ad3bb547890c722c7e1e85ad2e01204361f7400b36c77d0a17b3deef7c213c84d91e0edb8f7f64e35f76ab58a487fd4@43.130.68.86:30320",
//	"enode://b84442564b7ac83f7cfad8eeb938b13e3a6c3bb14ac2a21aead907362c9055e421aadcc583b96c1628d803c63db2ad6e043414f666fbfa37b06dfbfc34f1c993@43.130.3.69:30320",
//	"enode://370d9a40b73e6f1f0795214cfc4f7703d2231998ecd9f26d9038fe266a22b63e54e83ae223b8282ad0dc917fb8067b42065ac83010c218bf689df0dc4f3738fc@49.51.91.16:30320",
//	"enode://1ea818a87fdac02e211297cf9db1223f93f1885901060e08e2886e83ef029bcbe9c17257de260375e45ea7b09f56e2605039d8100755c101367fd5e879f78dc2@162.62.225.142:30320",
//	"enode://775a0973bc39ac999734eed0edefa5f325ac0aaf4e1bd1ba2e65a160f7018abae7f24b4199ef3a8368f3f1cb2fd2ee9f07ddd72903fb6f22f37f47a7d5b46277@43.133.195.3:30320",
//	"enode://503fa1c7fb94d1381425ae4eb3c8a62ea82118475349e03e1b4d93ead42609d6295957dd5c106c4c81b84f3139f4ab5819d8449af7c2dc2e6f7ac82369f695f2@119.28.151.131:30320",
//	"enode://a9939ea4e177ffd54f119c7241a72ed8273bbd199ed8972d2a4461641d8f6a50da6c9d9d9587c8c31790fd3f3b3dccb9d51f9bd9c05cfaa819a3061170893b3d@150.109.149.220:30320",
//}

//var DevnetBootnodes = []string{
//	"enode://8dfa963b177833d7d181324e9d260d88338274fb7e2855975efc6a9adddc9c335c01566b241c2fac2281cffd07055a8e63d5ef9f5110f1fa29e640dc529618e5@129.226.154.223:30320",
//	"enode://0b8f42d969ea2c9001ee83622e5db38e7df7b6036ec3cf63f9d62b2a7aa99558426a2aad8248e0ceed95b6921d11432ce804dd3c038e59f125b3fbdda214afe2@43.128.189.18:30320",
//	"enode://8c2d79a0a0a8dc1debcc00ce05563eea9ebb341c7e865518eeee3f341f2a41e8b35ecda518ee651de7438977b671378b8a4df0f632d31a813542368269144367@43.128.203.129:30320",
//	"enode://a145065c04ed7e4898f0c724c946c29496cb3a746bdaed88a6885936f07c48429dc93605d89ad65d708eb2a2ddbed793d80ff3e4c68be84ae69e6afc98a502d0@54.155.55.37:30320",
//	"enode://2c91e6771802fcea57c4ec67fc1ab5a2e13a006cbd206bb852e7b7a02d3a0dfb17fb50c1cab91998a55ecb02a73dabe90f147ea8522bda40e81ec1058fe5b218@15.236.202.37:30320",
//	"enode://fcaee3896fc898385a7b4e564236d807e160e491d35863fdc7c80157361726b9b27694d86bf38a8b0cfd9abb0720b2bda5a2f4386eb7525f2975cdecb62132c8@13.211.121.131:30320",
//	"enode://3061ee83846c47170cf2539209ba6d4849b0490224dea4311785fc3322db9eabb6de4684fab4d6d5331b9707a6f64b6882dce9545c03678646a8f5b896f4ef25@15.160.94.149:30320",
//	"enode://90dc742cba1bfef400ddcfd7aa9f400419cce740b99284bd3e04eb58a2c40b9883f8c7eab23c615f5c97ad8349b65634cade63f3831f2a843c65c9bfdc86b340@15.152.100.177:30320",
//	"enode://081b117da0d60514c2c247bfe47ee9c73ad3bb547890c722c7e1e85ad2e01204361f7400b36c77d0a17b3deef7c213c84d91e0edb8f7f64e35f76ab58a487fd4@43.130.68.86:30320",
//	"enode://17d82c989a517ee065081577039da71bdbbd07fc326e1f53a86d1ec5c9c420bb2ba83dc87e0cf3206af77ceced8ae545bcf76462e66f3e3c94fcbb0836a552ff@43.130.3.69:30320",
//	"enode://370d9a40b73e6f1f0795214cfc4f7703d2231998ecd9f26d9038fe266a22b63e54e83ae223b8282ad0dc917fb8067b42065ac83010c218bf689df0dc4f3738fc@49.51.91.16:30320",
//	"enode://9574dfb9e130e46a34e574da9b6a360cbfc9f254f19887c1d3f49e0b4847d6fa42d3440ea85983824a2d13f3d09c4a735e8c73d212f37eb90bda64f6eb115b6a@162.62.225.142:30320",
//	"enode://775a0973bc39ac999734eed0edefa5f325ac0aaf4e1bd1ba2e65a160f7018abae7f24b4199ef3a8368f3f1cb2fd2ee9f07ddd72903fb6f22f37f47a7d5b46277@43.133.195.3:30320",
//	"enode://422eaddae6804ee2fa1dc97fe954baab018e8151f7e794bf8d8783218353bcc8928a18d2d46a72a309d14e09c278647be6db2f9e26a8d97da56e44a30c685ccc@119.28.151.131:30320",
//	"enode://a9939ea4e177ffd54f119c7241a72ed8273bbd199ed8972d2a4461641d8f6a50da6c9d9d9587c8c31790fd3f3b3dccb9d51f9bd9c05cfaa819a3061170893b3d@150.109.149.220:30320",
//	"enode://03cd10c01f1b7a4979c9f9cd042ee6d3c73283966911c43ed7d6b1c420b2858e51714943f481eb3e52af16b875f9741898a7cec7da251787b9bb34c312edaf8e@18.169.133.50:30320",
//	"enode://50c4499a6c03c8f6545a9cf43692d8d169d5b00a72ba1666a9b51d1f0af52506802189cfff8f577a73f0bc40041e3657137bd3d9dfed7be0031783c912e27f75@47.91.123.160:30320",
//	"enode://5c0f71a344ffabff3a3866233ee016170d24e932101bf65c56f2327afdf900bdcbcadbdce8cf21cb85ad50c42600d3d85ab295423b185370cb02eafd5b6f6080@47.254.195.229:30320",
//	"enode://1da57711973dc120be4129771156cca903b703974d6a3f1f05004f499f1aaf703314c3b2d6be700ea9f234b8fa9f1c4f2c0029aaf3e33a2126027e02a320f143@54.177.66.169:30320",
//	"enode://a922890802ff51f65c745a66c06745a026245376cbabbb99c3e832fc9292be7112f96da279e9654c0a1c1d141473ec98d110da7687e3a105cd28f4a7bff1ccee@43.135.202.91:30320",
//}

var DevnetBootnodes = []string{
	"enode://8dfa963b177833d7d181324e9d260d88338274fb7e2855975efc6a9adddc9c335c01566b241c2fac2281cffd07055a8e63d5ef9f5110f1fa29e640dc529618e5@129.226.154.223:40320",
	"enode://0b8f42d969ea2c9001ee83622e5db38e7df7b6036ec3cf63f9d62b2a7aa99558426a2aad8248e0ceed95b6921d11432ce804dd3c038e59f125b3fbdda214afe2@43.128.189.18:40320",
	"enode://8c2d79a0a0a8dc1debcc00ce05563eea9ebb341c7e865518eeee3f341f2a41e8b35ecda518ee651de7438977b671378b8a4df0f632d31a813542368269144367@43.128.203.129:40320",
	"enode://a145065c04ed7e4898f0c724c946c29496cb3a746bdaed88a6885936f07c48429dc93605d89ad65d708eb2a2ddbed793d80ff3e4c68be84ae69e6afc98a502d0@54.155.55.37:40320",
	"enode://2c91e6771802fcea57c4ec67fc1ab5a2e13a006cbd206bb852e7b7a02d3a0dfb17fb50c1cab91998a55ecb02a73dabe90f147ea8522bda40e81ec1058fe5b218@15.236.202.37:40320",
	"enode://fcaee3896fc898385a7b4e564236d807e160e491d35863fdc7c80157361726b9b27694d86bf38a8b0cfd9abb0720b2bda5a2f4386eb7525f2975cdecb62132c8@13.211.121.131:40320",
	"enode://3061ee83846c47170cf2539209ba6d4849b0490224dea4311785fc3322db9eabb6de4684fab4d6d5331b9707a6f64b6882dce9545c03678646a8f5b896f4ef25@15.160.94.149:40320",
	"enode://90dc742cba1bfef400ddcfd7aa9f400419cce740b99284bd3e04eb58a2c40b9883f8c7eab23c615f5c97ad8349b65634cade63f3831f2a843c65c9bfdc86b340@15.152.100.177:40320",
	"enode://081b117da0d60514c2c247bfe47ee9c73ad3bb547890c722c7e1e85ad2e01204361f7400b36c77d0a17b3deef7c213c84d91e0edb8f7f64e35f76ab58a487fd4@43.130.68.86:40320",
	"enode://17d82c989a517ee065081577039da71bdbbd07fc326e1f53a86d1ec5c9c420bb2ba83dc87e0cf3206af77ceced8ae545bcf76462e66f3e3c94fcbb0836a552ff@43.130.3.69:40320",
	"enode://370d9a40b73e6f1f0795214cfc4f7703d2231998ecd9f26d9038fe266a22b63e54e83ae223b8282ad0dc917fb8067b42065ac83010c218bf689df0dc4f3738fc@49.51.91.16:40320",
	"enode://9574dfb9e130e46a34e574da9b6a360cbfc9f254f19887c1d3f49e0b4847d6fa42d3440ea85983824a2d13f3d09c4a735e8c73d212f37eb90bda64f6eb115b6a@162.62.225.142:40320",
	"enode://775a0973bc39ac999734eed0edefa5f325ac0aaf4e1bd1ba2e65a160f7018abae7f24b4199ef3a8368f3f1cb2fd2ee9f07ddd72903fb6f22f37f47a7d5b46277@43.133.195.3:40320",
	"enode://422eaddae6804ee2fa1dc97fe954baab018e8151f7e794bf8d8783218353bcc8928a18d2d46a72a309d14e09c278647be6db2f9e26a8d97da56e44a30c685ccc@119.28.151.131:40320",
	"enode://a9939ea4e177ffd54f119c7241a72ed8273bbd199ed8972d2a4461641d8f6a50da6c9d9d9587c8c31790fd3f3b3dccb9d51f9bd9c05cfaa819a3061170893b3d@150.109.149.220:40320",
	"enode://03cd10c01f1b7a4979c9f9cd042ee6d3c73283966911c43ed7d6b1c420b2858e51714943f481eb3e52af16b875f9741898a7cec7da251787b9bb34c312edaf8e@18.169.133.50:40320",
	"enode://50c4499a6c03c8f6545a9cf43692d8d169d5b00a72ba1666a9b51d1f0af52506802189cfff8f577a73f0bc40041e3657137bd3d9dfed7be0031783c912e27f75@47.91.123.160:40320",
	"enode://5c0f71a344ffabff3a3866233ee016170d24e932101bf65c56f2327afdf900bdcbcadbdce8cf21cb85ad50c42600d3d85ab295423b185370cb02eafd5b6f6080@47.254.195.229:40320",
	"enode://1da57711973dc120be4129771156cca903b703974d6a3f1f05004f499f1aaf703314c3b2d6be700ea9f234b8fa9f1c4f2c0029aaf3e33a2126027e02a320f143@54.177.66.169:40320",
	"enode://a922890802ff51f65c745a66c06745a026245376cbabbb99c3e832fc9292be7112f96da279e9654c0a1c1d141473ec98d110da7687e3a105cd28f4a7bff1ccee@43.135.202.91:40320",
}

//var DevnetBootnodes = []string{
//	"enode://7759e4908b2e70159cd7d096b24e9494642362a4c775a41f5d9656726d3936bfa51c80c0d57768ccd1c951bba30e0f896a075fae493d5261227c52003cf4136b@54.177.66.169:0?discport=30303",
//	"enode://b6cae1802bb27570362703491539636784bd1212e52bef6a2298446f61c1c03e00011c12eb46c889588b8976a2a9e1759d6b1da75eb852f1bf04029cc56f3f31@43.135.202.91:0?discport=30303",
//}

const dnsPrefix = "enrtree://AKA3AM6LPBYEUDMVNU3BSVQJ5AD45Y7YPOHJLEF6W26QOE4VTUDPE@"

// KnownDNSNetwork returns the address of a public DNS-based node list for the given
// genesis hash and protocol. See https://github.com/ethereum/discv4-dns-lists for more
// information.
func KnownDNSNetwork(genesis common.Hash, protocol string) string {
	var net string
	switch genesis {
	case MainnetGenesisHash:
		net = "mainnet"
	case RopstenGenesisHash:
		net = "ropsten"
	case RinkebyGenesisHash:
		net = "rinkeby"
	case GoerliGenesisHash:
		net = "goerli"
	case TestNetGenesisHash:
		net = "testnet"
	case DevNetGenesisHash:
		net = "devnet"
	default:
		return ""
	}
	return dnsPrefix + protocol + "." + net + ".ethdisco.net"
}
