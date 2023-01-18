# OKCoin Proof-of-Reserves

## Background

OKCoin launches Proof of Reserves (PoR) to improve the security and transparency of user's assets. These tools will allow
you to independently audit OKCoin's Proof of Reserves and verify OKCoin's reserves exceed the exchange's known liabilities to
users, in order to confirm the solvency of OKCoin.

## Introduction

### Building the source

Download the latest build for your operating system and architecture. Also, you can build the source by yourself. 

Building this open source tool requires Go (version >= 1.17).

Install dependencies
```shell
 go mod tidy 
```

Compile
```shell
make all
```

### Executable

Proof-of-Reserves executable are in the cmd directory

|    Command    | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| :-----------: | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
|   `VerifyAddress`    | We have signed a specific message with a private key to each address published by OKCoin. This tool can be used to verify OKCoin's signature and verify OKCoin's ownership of the address.  |
|   `CheckBalance`    | Configure blockchain node RPC or OKLink API to use this tool, you can check the balance on the chain corresponding to the snapshot height of OKCoin, then compare it with the balance published by OKCoin, and query the total assets of OKCoin's wallet address on the chain. |
|   `MerkleValidator`    | OKCoin's PoR uses a Merkle tree, and you can use this tool to check whether your account assets are included in the Merkle tree published by OKCoin. |

## Reserves

Download OKCoin's Proof of Reserves File, verify the ownership of the OKCoin's public address, and check whether the OKCoin
snapshot height balance is consistent with the published balance.  Details here

### VerifyAddress

OKCoin's public file contains address, message "I am an OKCoin address" and signature. You can use VerifyAddress to verify
OKCoin's ownership of published address.

```shell
  ./build/VerifyAddress  --por_csv_filename ./example/okcoin_por_example.csv
```

At the same time, you can use third-party tools to verify the ownership
of [BTC single addresses](https://www.bitcoin.com/tools/verify-message/), [EVM](https://etherscan.io/verifiedsignatures)
, and [TRX addresses](https://tronscan.org/#/tools/verify-sign).

### CheckBalance

You can use CheckBalance to verify the OKCoin wallet address balance with the corresponding block height snapshot. [Details here](./docs/checkbalance.md)

Sum of all address balances

```shell
  ./build/CheckBalance --rpc_json_filename="./example/rpc.json" --por_csv_filename ./example/okcoin_por_example.csv
```

Query the snapshot height balance on the chain

```shell
./build/CheckBalance --mode="single_coin" --coin_name="ETH" --rpc_json_filename="./example/rpc.json" --por_csv_filename="./example/okcoin_por_example.csv"
```

## Liabilities

OKCoin's PoR uses Merkle tree technology to allow each user to independently review OKCoin's digital asset reserve on the
basis of protecting user privacy. Details here

### MerkleValidator

Log in to OKCoin, go to the Audits page to view audit details, download the Merkle tree path data, copy and save it as a
file merkle_proof_file.json, and run the following command to check whether your assets are included in the total user
assets of OKCoin. 

```shell
./build/MerkleValidator --merkle_proof_file ./example/merkle_proof_file.json
```





