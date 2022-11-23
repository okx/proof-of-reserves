# Verify OKX Wallet Reserve Address Balance

You can verify OKX Wallet Reserve Address Balance with the following steps:
1. Compare the address balance of a specific crypto chain at the snapshot height with the address balance in the snapshot file published by OKX to verify whether they are consistent.
2. Compare the sum of address balances of a specific crypto chain at the snapshot height with the sum of address balances in the snapshot file published by OKX to verify whether they are consistent.

When querying the address balance at the height of the snapshot, you can configure node rpc for balance query, or configure OKLink's open API for balance query. The following is a detailed step-by-step description.

## Configure the rpc.json file

There is an rpc.json configuration file in the binary zip and in the source code directory. We use this file to configure node rpc, OKLink open API, token contract address and other information, so that our tool can query the balance on the chain and compare it with the balance in the snapshot file for verification. A specific example of configuration is as follows:

- If you build a node yourself or use a third-party node, please configure the rpc option of the corresponding crypto. For details, please see [Get Node RPC](#get-node-rpc)

```text
"rpc": {
    "endpoint": "http://127.0.0.1:8332/",       // node rpc
    "jsonPattern": "$.result.total_amount",
    "defaultUnit": "BTC",                       // display unit
    "customHeaders": {},
    "authUser": "OKX",                          // rpcuser
    "authPassword": "OKXWallet",                // rpcpassword
    "tokenAddress": "",                         // token contract address
    "enabled": true                             // this configuration switch enables node queries
}
```

- If you use OKLink open API, please configure the api option of the corresponding crypto. (Note: You need to go to the OKLink website to obtain the apiKey. For details, please see [Get OKLink apiKey](#get-oklink-apikey))

```text
"api": {
    "endpoint": "https://www.oklink.com/api/v5/explorer/block/address-balance-history", // OKLink API
    "jsonPattern": "$.data[0].balance",
    "defaultUnit": "",                          // display unit
    "tokenAddress": "",                         // token contract address
    "customHeaders": {
        "Ok-Access-Key": "OKLink apiKey",       // OKLink apiKey
        "Content-Type": "application/x-www-form-urlencoded",
        "Accept": "*/*"
    },
    "enabled": true                             // this configuration switch enables OKLink API queries
}
```

## Verify Balance

Once you have obtained the executable and snapshot file and configured rpc.json, you can start verifying the balance. Please see below commands:

```bash
# Verify single address balance
./CheckBalance --mode="single_address" --coin_name="btc" --address="3A1JRKqfGGxoq2qSHLv85u4zn935VR9ToL" --por_csv_filename=okx_por_20221122.csv

# Verify the total of all address balances
./CheckBalance --mode="single_coin_total_balance" --coin_name="btc" --por_csv_filename=okx_por_20221122.csv

```

Running the above command will output the balance data on the chain, and you can view the data in the snapshot file for comparison and verification. For a detailed description of the command, please see [CheckBalance Command Introduction](#checkbalance-command-introduction)

## Get Node RPC

### Prepare Bitcoin Core Node

Install the Bitcoin Core client, synchronize to the latest height, and then roll back the block to the height of the OKX snapshot.

1. You can download Bitcoin Core here: <https://bitcoincore.org/en/download/> , Please download version 0.21 or later.
2. You need to edit the configuration file of Bitcoin Core, in order to visit RPC node. Create ~/.bitcoin/bitcoin.conf file and open it with an editor or run the command of `vi ~/.bitcoin/bitcoin.conf`, edit

    ```bash
    server=1
    rpcuser=OKX
    rpcpassword=OKXWallet
    ```

3. Enter the bin directory, run `./bitcoind` command, and start the node.
4. Wait for the node to synchronize to the latest height, it will take about 12 hours.
5. After synchronizing to the latest height, you need to roll back the node to the OKX snapshot height to query the balance of the snapshot height. The steps are as follows:
    1. Go to the BTC browser to query the block hash of the next height of the snapshot height, copy this block hash and fill in the hash value part of the command below.
    2. Run `./bitcoin-cli invalidateblock 00000000000000000005829017993a7a21e4b7c731c95b9cb979c01294a7bd27` command.
    3. Wait for the node to roll back to the snapshot height, you can run `./bitcoin-cli getblockcount` command to check whether the roll back is complete, you can also view the node output log judgment.

### Get EVM Archive Node

- Install the archive node manually. It may take some time to synchronize, reference: <https://geth.ethereum.org/docs/install-and-build/installing-geth>

- Use third-party nodes: [Infura](https://infura.io/), [Alchemy](https://alchemy.com/) e.g: Alchemy provides ethereum archive node service
  <https://www.alchemy.com/overviews/archive-nodes>

## Get OKLink apiKey

1. First login to [OKLink](https://www.oklink.com/en/account/login).
2. Click on the human icon in the upper right corner, a drop-down menu will be displayed, select API.
3. Click the Create API button in the On-chain Data section to create the apiKey.

## CheckBalance Command Introduction

### Flags

* address: Set the address of a blockchain that needs to be verified
* coin_name: Set the name of the blockchain to be verified, default: ETH
* mode: Set the mode to verify the balance
* rpc_json_filename: Set rpc.json file path, default: rpc.json(root directory)
* por_csv_filename: Set por csv data file path

### Usage

**coin_name** Supported: 'BTC','ETH','ETH-ARBITRUM','ETH-OPTIMISM','USDT-ERC20','USDT-TRC20','USDT-POLY','USDT-AVAXC','USDT-ARBITRUM','USDT-OPTIMISM','USDT-OMNI'

**mode** Supported: 'single_address','single_coin_total_balance','single_coin','all_coin','all_coin_total_balance'

* single_address: Single address mode, support to verify the balance of an address, need to pass address and coin_name
* single_coin_total_balance: Single coin total balance mode, support to verify the total balance of all addresses of a single coin, need to pass coin_name
* single_coin: Single coin mode, support to verify all address balance of single coin, need to pass coin_name
* all_coin: All coin mode, support to verify the balance of each address in all coins
* all_coin_total_balance: All coin total balance mode, support to verify the total balance of addresses in all coins

**Note**

* It is recommended to use single_address and single_coin_total_balance mode. Due to the huge number of addresses of some coins, it takes a long time to use node RPC verification, single_coin, all_coin and all_coin_total_balance verification time may be long

### Example

single address:

```shell
./CheckBalance --mode="single_address" --coin_name="ETH" --address="0x07e47ed3c5a8ff59fb5d1df4051c34da67fc5547" --rpc_json_filename="rpc.json" --por_csv_filename="por_test.csv"
```

single coin total address:

```shell
./CheckBalance --mode="single_coin_total_balance" --coin_name="ETH" --rpc_json_filename="rpc.json" --por_csv_filename="por_test.csv"
```
