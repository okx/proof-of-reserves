# OKX Proof-of-Reserves

## Background

OKX launches [Proof of Reserves (PoR)](https://www.okx.com/proof-of-reserves) to improve the security and transparency
of user's assets. These tools will allow you to independently audit OKX's Proof of Reserves and verify OKX's reserves
exceed the exchange's known liabilities to users, in order to confirm the solvency of OKX.

## Introduction

### Building the source

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

### Executable

Proof-of-Reserves executable are in the cmd directory

|    Command    | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| :-----------: | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
|   `VerifyAddress`    | We have signed a specific message with a private key to each address published by OKX. This tool can be used to verify OKX's signature and verify OKX's ownership of the address.  |
|   `CheckBalance`    | Configure blockchain node RPC or OKLink API to use this tool, you can check the balance on the chain corresponding to the snapshot height of OKX, then compare it with the balance published by OKX, and query the total assets of OKX's wallet address on the chain. |
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
privacy. [Details here](https://www.okx.com/support/hc/en-us/articles/10660988139661-How-to-verify-if-your-assets-are-included-in-the-OKX-Merkle-tree-)

### MerkleValidator

**Currently, version v2 has been released, which is compatible with v1. The following is the operation flow of the two
versions**

Verification process for v1 version

1. Log in to OKX, go to the Audits page to view audit details, download the Merkle tree path data, copy and save it as a
   file merkle_proof_file.json, and run the following command to check whether your assets are included in the total
   user assets of OKX.

```shell
./build/MerkleValidator --merkle_proof_file ./example/merkle_proof_file.json
```

Verification process for v2 version

1. Visit [Full merkle tree file download page](https://okx.com/proof-of-reserves/download?tab=liabilities) to download
   the full merkle tree file
2. Login to OKX and visit Audit page to copy and save the data as user_info_file.json
3. Run the following command to check whether your assets are included in the total user assets of OKX.

```shell
./build/MerkleValidator --merkle_file ./example/full-liabilities-merkle-tree.txt --user_info_file ./example/user_info_file.json
```

