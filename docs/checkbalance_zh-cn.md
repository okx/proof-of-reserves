# 验证欧易钱包储备金地址余额

您可以通过以下步骤验证：

1. 快照高度时刻的特定币种链上地址余额，与欧易公布的快照文件中地址余额进行对比，从而验证是否一致。
2. 快照高度时刻的特定币种链上地址余额总和，与欧易公布的快照文件中地址余额总和进行对比，验证是否一致。

查询快照高度时刻地址余额时，可以配置节点rpc进行余额查询，或者配置OKLink的open API进行余额查询。以下是详细步骤介绍：

## 配置rpc.json文件

在二进制压缩包中和源代码目录中都有一个rpc.json配置文件，我们使用此文件配置节点rpc、OKLink open
API、token合约地址等信息，以便我们的工具可以查询链上余额，并与快照文件中的余额做对比验证。配置的具体示例如下：

- 如果您自己搭建节点或使用三方节点，请配置对应币种的rpc选项。详细请看[节点RPC获取](#节点rpc获取)

```text
// RPC配置支持的币种：
// 'BTC','ETH','ETH-ARBITRUM','ETH-OPTIMISM','USDT-ERC20','USDT-POLY','USDT-AVAXC','USDT-ARBITRUM','USDT-OPTIMISM'

"rpc": {
    "endpoint": "http://127.0.0.1:8332/",       // 节点rpc
    "jsonPattern": "$.result.total_amount",
    "defaultUnit": "BTC",                       // 显示单位
    "customHeaders": {},
    "authUser": "OKX",                          // 访问节点需要的用户名，如访问btc节点
    "authPassword": "OKXWallet",                // 访问节点需要的密码，如访问btc节点
    "tokenAddress": "",                         // token合约地址
    "enabled": true                             // 此配置开关启用节点查询
}
```

- 如果您使用OKLink open API，请配置对应币种的api选项。（注意：您需要到OKLink网站获取apiKey。详细请看[获取OKLink apiKey](#获取oklink-apikey)）

```text
// API配置支持的币种：
// 'ETH','ETH-ARBITRUM','ETH-OPTIMISM','USDT-ERC20','USDT-TRC20','USDT-POLY','USDT-AVAXC','USDT-ARBITRUM','USDT-OPTIMISM','USDT-OMNI'

"api": {
    "endpoint": "https://www.oklink.com/api/v5/explorer/block/address-balance-history", // OKLink API
    "jsonPattern": "$.data[0].balance",
    "defaultUnit": "",                          // 显示单位
    "tokenAddress": "",                         // token合约地址
    "customHeaders": {
        "Ok-Access-Key": "OKLink apiKey",       // OKLink apiKey
        "Content-Type": "application/x-www-form-urlencoded",
        "Accept": "*/*"
    },
    "enabled": true                             // 此配置开关 走OKLink API查询
}
```

**说明**
1. 当同时配置RPC和API时，会使用RPC进行查询。
2. 快照文件包含质押在Compound平台的USDT，相关配置如下。由于节点不支持查询质押在Compound平台的USDT，当USDT-ERC20 RPC配置开启时，会从快照数据中获取质押在Compound平台上的USDT余额， 您可以使用[Compound API](https://api.compound.finance/api/v2/account?addresses%5B%5D=0xb99cc7e10fe0acc68c50c7829f473d81e23249cc&block_number=16023042)进行验证。

```text
{
    "name": "usdt-erc20",
    "coin": "eth",
    "api": { ... },
    "rpc": { ... },
    "whiteList": [
        {
          "project": "comp",
          "projectFullName": "compound",
          "address": "0xb99cc7e10fe0acc68c50c7829f473d81e23249cc",
          "tokenAddress": "0xf650c3d88d12db855b8bf7d11be6c55a4e07dcc9"
        }
    ]
}
```


## 验证余额

获取了可执行文件和快照文件，并配置了rpc.json后，您可以开始验证余额。具体操作如下：

```bash
# 验证单个地址余额
./CheckBalance --mode="single_address" --coin_name="btc" --address="3A1JRKqfGGxoq2qSHLv85u4zn935VR9ToL" --por_csv_filename=okx_por_20221122.csv

# 验证所有地址余额总和
./CheckBalance --mode="single_coin_total_balance" --coin_name="btc" --por_csv_filename=okx_por_20221122.csv

```

运行以上命令会输出链上余额数据，您可以查看快照文件中的数据进行对比验证。命令详细介绍请看[CheckBalance命令介绍](#checkbalance命令介绍)

## 节点RPC获取

### Bitcoin节点

安装Bitcoin Core客户端，同步到最新高度，然后将区块回滚到欧易快照时的高度。

1. 可在此处下载Bitcoin Core软件：<https://bitcoincore.org/en/download/> ，请下载 0.21 或之后的版本。
2. 需要编辑Bitcoin Core的配置文件，以使节点RPC可访问。创建 ~/.bitcoin/bitcoin.conf 文件并用编辑器打开或运行 `vi ~/.bitcoin/bitcoin.conf` 命令，编辑

    ```bash
    server=1
    rpcuser=OKX
    rpcpassword=OKXWallet
    ```

3. 进入 bin 目录，运行 `./bitcoind` 命令，启动节点。
4. 等待节点同步到最新高度，大约需要12个小时。
5. 同步到最新高度后，需要回滚节点到欧易快照高度，以查询快照高度余额，操作如下：
    1. 上BTC浏览器上查询快照高度的下一个高度的区块hash，复制此区块hash并填入下面命令的hash值部分。
    2. 运行 `./bitcoin-cli invalidateblock 00000000000000000005829017993a7a21e4b7c731c95b9cb979c01294a7bd27` 命令。
    3. 等待节点回滚到快照高度，可以运行 `./bitcoin-cli getblockcount` 命令查看是否回滚完成，也可以查看节点输出日志判断。

### 获取evm系归档节点

- 手动安装归档节点。可能需要一定时间同步，参考<https://geth.ethereum.org/docs/install-and-build/installing-geth>

- 使用[Infura](https://infura.io/) 、[Alchemy](https://alchemy.com/) 等第三方节点，例如：Alchemy提供了ethereum归档节点服务<https://www.alchemy.com/overviews/archive-nodes>

## 获取OKLink apiKey

1. 登录[OKLink](https://www.oklink.com/en/account/login)
2. 点击右上角的人形图标，会显示一个下拉菜单，选择API
3. 点击链上数据部分的创建API按钮创建apiKey

## CheckBalance命令介绍

### Flags

* address: 设置需要验证的某条公链的地址
* coin_name: 设置需要验证的公链名称，默认: ETH
* mode: 设置验证余额的模式
* rpc_json_filename: 设置rpc.json文件路径，默认: rpc.json(根目录)
* por_csv_filename: 设置por csv数据文件路径

### Usage

**coin_name** 支持: 'BTC','ETH','ETH-ARBITRUM','ETH-OPTIMISM','USDT','USDT-ERC20','USDT-TRC20','USDT-POLY','USDT-AVAXC','USDT-ARBITRUM','USDT-OPTIMISM','USDT-OMNI',
'USDC','POLY-USDC','USDC-AVAXC','USDC-ARBITRUM','USDC-OPTIMISM','USDC-OKC20','OKB','OKB-OKC20','OKT','FILK-OKC20','SHIBK-KIP20','DOTK-OKC20','XRPK-KIP20','UNIK-OKC20',
'LINKK-OKC20','TRXK-KIP20','SHIB','UNI','LINK','PEOPLE'

**mode** 支持: 'single_address','single_coin_total_balance','single_coin','all_coin','all_coin_total_balance'

* single_address: 单地址模式，支持验证某个地址的余额，需要配合传参 address 和 coin_name
* single_coin_total_balance: 单币种总余额模式，支持验证单个币种所有地址的总余额，需要配合传参 coin_name
* single_coin: 单币种模式，支持验证单个币种所有地址余额，需要配合传参数 coin_name
* all_coin: 全币种模式，支持验证所有币种的各地址余额
* all_coin_total_balance: 全币种总余额模式，支持验证所有币种的地址总余额

**注意**

* 建议使用 single_address 和 single_coin_total_balance 模式，由于有的币种地址数量庞大，采用节点RPC验证的方式，耗时较长，single_coin、all_coin和all_coin_total_balance 验证时间可能会很长

### Example

single address：

```shell
./CheckBalance --mode="single_address" --coin_name="ETH" --address="0x07e47ed3c5a8ff59fb5d1df4051c34da67fc5547" --rpc_json_filename="rpc.json" --por_csv_filename="por_test.csv"
```

single coin total address：

```shell
./CheckBalance --mode="single_coin_total_balance" --coin_name="ETH" --rpc_json_filename="rpc.json" --por_csv_filename="por_test.csv"
```
