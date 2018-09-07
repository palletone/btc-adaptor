## 准备

1.安装go

## 下载源码、编译

### btcd

1.下载 btcd 代码

```
mkdir F:\work\src\golang.org\x
cd F:\work\src\golang.org\x
git clone https://github.com/golang/crypto.git
go get -u github.com/btcsuite/btcd
```

2.编译出 btcd 可执行程序

```
cd F:\work\src\github.com\btcsuite\btcd
go build
```

### btcwallet

1.下载 btcwallet 代码

linux need clone `github.com/golang.org/x/sys`

```
cd /home/pallet/go/src/golang.org/x
git clone https://github.com/golang/sys.git
```

```
mkdir F:\work\src\google.golang.org
cd F:\work\src\google.golang.org
git clone https://github.com/grpc/grpc-go.git
rename grpc-go grpc
git clone https://github.com/google/go-genproto.git
rename go-genproto genproto

cd F:\work\src\golang.org\x
git clone https://github.com/golang/net.git
git clone https://github.com/golang/text.git

go get -u github.com/aead/siphash
go get -u github.com/coreos/bbolt
go get -u github.com/davecgh/go-spew/spew
go get -u github.com/golang/protobuf/proto
go get -u github.com/jessevdk/go-flags
go get -u github.com/jrick/logrotate/rotator
go get -u github.com/kkdai/bstream
go get -u github.com/lightninglabs/gozmq
go get -u github.com/lightninglabs/neutrino

go get -u github.com/btcsuite/btcwallet
```

2.编译出 btcwallet 可执行程序

```
cd F:\work\src\github.com\btcsuite\btcwallet
go build
```


## 启动

+ 启动 btcd

```
cd F:\work\src\github.com\btcsuite\btcd
.\btcd.exe -u test -P 123456 --datadir "d:\\btctest\\" --testnet --txindex --addrindex
```

-u 是 rpcuser ，-P 是 rpcpasswd ， --datadir 是区块数据目录， --testnet 是指定测试链
去掉参数 --testnet 即是正式链
--txindex 是建立交易索引， --addrindex 是建立地址索引，两者顺序不能颠倒

+ 启动 btcwallet

先创建钱包，再启动

```
cd F:\work\src\github.com\btcsuite\btcwallet
.\btcwallet.exe -u test -P 123456 --datadir "d:\\test\\" --create --testnet
.\btcwallet.exe -u test -P 123456 --datadir "d:\\test\\" --testnet
```
-u 是 rpcuser ，-P 是 rpcpasswd ， -datadir 是钱包数据目录， --testnet 是指定测试链
去掉参数 --testnet 即是正式链



## 示例

```
F:\work\src\github.com\btcsuite\btcd>.\btcd.exe -u test -P 123456 --datadir "d:\\btctest\\" --testnet
2018-06-27 13:55:31.585 [INF] BTCD: Version 0.12.0-beta
2018-06-27 13:55:31.610 [INF] BTCD: Loading block database from 'd:\btctest\testnet\blocks_ffldb'
2018-06-27 13:55:31.659 [INF] BTCD: Block database loaded
2018-06-27 13:55:31.686 [INF] INDX: cf index is enabled
2018-06-27 13:55:31.688 [INF] INDX: Catching up indexes from height -1 to 0
2018-06-27 13:55:31.689 [INF] INDX: Indexes caught up to height 0
2018-06-27 13:55:31.689 [INF] CHAN: Chain state (height 0, hash 000000000933ea01ad0ee984209779baaec3ced90fa3f408719526f8d77f4943, totaltx 1, work 4295032833)
2018-06-27 13:55:31.700 [INF] RPCS: RPC server listening on [::1]:18334
2018-06-27 13:55:31.700 [INF] RPCS: RPC server listening on 127.0.0.1:18334
2018-06-27 13:55:31.700 [INF] AMGR: Loaded 0 addresses from file 'd:\btctest\testnet\peers.json'
2018-06-27 13:55:31.701 [INF] CMGR: Server listening on 0.0.0.0:18333
2018-06-27 13:55:31.701 [INF] CMGR: Server listening on [::]:18333
2018-06-27 13:55:31.739 [INF] CMGR: 1 addresses found from DNS seed testnet-seed.bluematt.me
2018-06-27 13:55:31.739 [INF] CMGR: 25 addresses found from DNS seed testnet-seed.bitcoin.schildbach.de
2018-06-27 13:55:31.740 [INF] CMGR: 23 addresses found from DNS seed testnet-seed.bitcoin.jonasschnelli.ch
2018-06-27 13:55:31.742 [INF] CMGR: 24 addresses found from DNS seed seed.tbtc.petertodd.org
2018-06-27 13:55:36.964 [INF] SYNC: New valid peer 13.78.14.162:18333 (outbound) (/Satoshi:0.16.0/)
2018-06-27 13:55:36.965 [INF] SYNC: Syncing to block height 1326466 from peer 13.78.14.162:18333
2018-06-27 13:55:36.968 [INF] SYNC: Downloading headers for blocks 1 to 546 from peer 13.78.14.162:18333
2018-06-27 13:55:37.030 [INF] SYNC: New valid peer 172.105.194.235:18333 (outbound) (/Satoshi:0.16.0(bitcore)/)
```

```
F:\work\src\github.com\btcsuite\btcwallet>.\btcwallet.exe -u test -P 123456 --datadir "d:\\test\\" --create --testnet
datadir option has been replaced by appdata -- please update your config
Enter the private passphrase for your new wallet:
Confirm passphrase:
Do you want to add an additional layer of encryption for public data? (n/no/y/yes) [no]:
Do you have an existing wallet seed you want to use? (n/no/y/yes) [no]:
Your wallet generation seed is:
de37216919d23faf454a37f2efa73e5207487a1526a196715d381662d9acded4
IMPORTANT: Keep the seed in a safe place as you
will NOT be able to restore your wallet without it.
Please keep in mind that anyone who has access
to the seed can also restore your wallet thereby
giving them access to all your funds, so it is
imperative that you keep it in a secure location.
Once you have stored the seed in a safe and secure location, enter "OK" to continue: OK
Creating the wallet...
2018-06-27 13:52:53.799 [INF] WLLT: Opened wallet
The wallet has been created successfully.
F:\work\src\github.com\btcsuite\btcwallet>
```

```
F:\work\src\github.com\btcsuite\btcwallet>.\btcwallet.exe -u test -P 123456 --datadir "d:\\test\\" --testnet
datadir option has been replaced by appdata -- please update your config
2018-06-27 13:53:08.999 [WRN] BTCW: open d:\test\btcwallet.conf: The system cannot find the file specified.
2018-06-27 13:53:09.023 [INF] BTCW: Version 0.7.0-alpha
2018-06-27 13:53:09.023 [INF] BTCW: Generating TLS certificates...
2018-06-27 13:53:09.058 [INF] BTCW: Done generating TLS certificates
2018-06-27 13:53:09.058 [INF] RPCS: Listening on 127.0.0.1:18332
2018-06-27 13:53:09.058 [INF] RPCS: Listening on [::1]:18332
2018-06-27 13:53:09.058 [INF] BTCW: Attempting RPC client connection to localhost:18334
2018-06-27 13:53:09.801 [INF] WLLT: Opened wallet
```
