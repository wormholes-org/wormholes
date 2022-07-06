一、编辑配置文件(app.conf)
    1、配置mysql相关参数，例如：
        #数据库参数配置
        #"admin:user123456@tcp(192.168.1.238:3306)/"
        #用户名
        dbusername = admin
        #用户密码
        dbuserpassword = user123456
        #数据库服务器地址
        dbserverip = 192.168.56.128
        #数据库服务器端口
        dbserverport = 3306
        #数据库名称
        dbname = nftdb
    2、配置合约相关参数，例如：
        #交易所合约
        TradeAddr = 0xD8D5D49182d7Abf3cFc1694F8Ed17742886dDE82
        #1155合约
        NFT1155Addr = 0xA1e67a33e090Afe696D7317e05c506d7687Bb2E5
        #管理员列表合约
        AdminAddr = 0x56c971ebBC0cD7Ba1f977340140297C0B48b7955
        #合约事件节点接入点
        EthersNode = https://rinkeby.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161
        EthersWsNode = wss://rinkeby.infura.io/ws/v3/97cb2119c79842b7818a7a37df749b2b

二、签名配置文件
    签名前需要把app.conf文件中[time]标签（包括[time]标签）后的数据全部删除。
    ./signappconf -f app.conf -key 不带0x开头的私钥
    把签名后的app.conf文件拷贝到与nftserver执行文件同目录下的conf目录中。

三、启动nftserver服务
    setsid ./nftserver > log
