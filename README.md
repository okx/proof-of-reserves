# OKX Proof-of-Reserves

## Background

OKX launches [Proof of Reserves (PoR)](https://www.okx.com/proof-of-reserves) to improve the security and transparency
of user's assets. These tools will allow you to independently audit OKX's Proof of Reserves and verify OKX's reserves
exceed the exchange's known liabilities to users, in order to confirm the solvency of OKX.

## Introduction

### Building the source

#### Reserves part

Download the [latest build](https://github.com/okx/proof-of-reserves/releases/latest) for your operating system and
architecture. Also, you can build the source by yourself.

Building this open source tool requires Go (version >= 1.17).

Install dependencies

```shell
 go mod tidy 
```

Compile

```shell
make all
```
#### Liabilities part
Building this open source tool requires Python(version >=3.9) 

Install dependencies
```shell
  pip install pycryptodome
  pip install pyinstaller
```

Compile
```shell
  pyinstaller -F zk_STARK_Validator.py
```

### Executable

Proof-of-Reserves executable are in the cmd directory

|    Command    | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| :-----------: | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
|   `VerifyAddress`    | We have signed a specific message with a private key to each address published by OKX. This tool can be used to verify OKX's signature and verify OKX's ownership of the address.  |
|   `CheckBalance`    | Configure blockchain node RPC or OKLink API to use this tool, you can check the balance on the chain corresponding to the snapshot height of OKX, then compare it with the balance published by OKX, and query the total assets of OKX's wallet address on the chain. |
|   `zkSTARKValidator`    | Current OKX's PoR uses zk-STARK(Zero-Knowledge Scalable Transparent Argument of Knowledge), a cryptographic proof technology, to verify data and prove the authenticity of our audits. |
|   `MerkleValidator`    | OKX's PoR uses a Merkle tree, and you can use this tool to check whether your account assets are included in the Merkle tree published by OKX. |

## Reserves

Download OKX's [Proof of Reserves File](https://www.okx.com/proof-of-reserves/download), verify the ownership of the
OKX's public address, and check whether the OKX snapshot height balance is consistent with the published
balance.  [Details here](https://www.okx.com/support/hc/en-us/articles/10781041719437-How-to-verify-OKX-s-ownership-and-balance-of-the-wallet-address-)

### VerifyAddress

OKX's public file contains address, message "I am an OKX address" and signature. You can use VerifyAddress to verify
OKX's ownership of published address.

```shell
  ./build/VerifyAddress  --por_csv_filename ./example/okx_por_example.csv
```

At the same time, you can use third-party tools to verify the ownership
of [BTC single addresses](https://www.bitcoin.com/tools/verify-message/), [EVM](https://etherscan.io/verifiedsignatures)
, and [TRX addresses](https://tronscan.org/#/tools/verify-sign).

### CheckBalance

You can use CheckBalance to verify the OKX wallet address balance with the corresponding block height
snapshot. [Details here](./docs/checkbalance.md)

Sum of all address balances

```shell
  ./build/CheckBalance --rpc_json_filename="./example/rpc.json" --por_csv_filename ./example/okx_por_example.csv
```

Query the snapshot height balance on the chain

```shell
./build/CheckBalance --mode="single_coin" --coin_name="ETH" --rpc_json_filename="./example/rpc.json" --por_csv_filename="./example/okx_por_example.csv"
```

## Liabilities

OKX's PoR uses Merkle tree technology to allow each user to independently review OKX's digital asset reserve on the
basis of protecting user
privacy. [Details here](https://www.okx.com/support/hc/en-us/articles/15165477403917)

### zk-STARK Validator(From 04/27)
**Currently, zk-STARK PoR has been released. The following is the operation flow of the zk-STARK version.**
Verify the inclusion constraint

1. To verify if your asset balance has been included as a Merkle leaf, navigate to ["Audits"](https://www.okx.com/balance/audit) and click **Details** to access your audit data.
2. Get the data you need for verification by clicking **Copy data** and pasting the JSON string as a file in a new folder. 
3. Download [zk-STARKValidator](https://github.com/okx/proof-of-reserves/releases/tag/v3.0.0), the OKX open-source verification tool, and save it to the same folder containing the JSON file. 
4. Open zk-STARKValidator to auto-run the JSON file you saved to check whether the inclusion constraint is satisfied. 

Verify the total balance and non-negative constraints
1. Under ["Audit files"](https://www.okx.com/proof-of-reserves/download?tab=liabilities), download the zk-STARK file from the "Liability report" tab.
2. Unzip the file to reveal a "sum proof data" folder with branch and trunk folders containing "sum_proof.json," "sum_value.json" files.
3. Download [zk-STARKValidator](https://github.com/okx/proof-of-reserves/releases/tag/v3.0.0), the OKX open-source verification tool, and place it in the same root folder as the "sum proof data" folder.
4. Open zk-STARKValidator to auto-run the unzipped zk-STARK file to check whether the total balance and non-negative constraints are satisfied. 



### MerkleValidator(From 11/22/2022 to 03/16)

**On 03/21, version v2 was released, which is compatible with v1. The following is the operation flow of the two versions**

Verification process for v1 version [Detail here](https://www.okx.com/support/hc/en-us/articles/10660988139661-How-to-verify-if-your-assets-are-included-in-the-OKX-Merkle-tree-)

1. Log in to OKX, go to the Audits page to view audit details, download the Merkle tree path data, copy and save it as a
   file merkle_proof_file.json, and run the following command to check whether your assets are included in the total
   user assets of OKX.

```shell
./build/MerkleValidator --merkle_proof_file ./example/merkle_proof_file.json
```

Verification process for v2 version [Detail here](https://www.okx.com/support/hc/en-us/articles/13747778159245-How-to-verify-if-your-assets-are-included-in-the-OKX-Merkle-tree-Merkle-Tree-V2-)

1. Visit [Full merkle tree file download page](https://okx.com/proof-of-reserves/download?tab=liabilities) to download
   the full merkle tree file
2. Login to OKX and visit Audit page to copy and save the data as user_info_file.json
3. Run the following command to check whether your assets are included in the total user assets of OKX.

```shell
./build/MerkleValidator --merkle_file ./example/full-liabilities-merkle-tree.txt --user_info_file ./example/user_info_file.json
```

