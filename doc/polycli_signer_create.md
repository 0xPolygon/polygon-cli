# `polycli signer create`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Create a new key

```bash
polycli signer create [flags]
```

## Usage

This command is meant to allow for easy creation of signed transactions. A raw transaction can then be published with a call to
[`eth_sendRawTransaction`](https://ethereum.org/en/developers/docs/apis/json-rpc/#eth_sendrawtransaction) or using [`cast publish`](https://book.getfoundry.sh/reference/cast/cast-publish).

## Creating Keys

The create subcommand will create a new key pair. By default, a hex private key will be written to `stdout`.
```bash
polycli signer create > private-key.txt
```

If you need to work with a go-ethereum style keystore, a key can be added by setting a `--keystore` directory. When you run this command, you'll need to specify a password to encrypt the private key.

```bash
polycli signer create --keystore /tmp/keystore
```

Polycli also has basic support for KMS with GCP. Creating a new key in the cloud can be accomplished with a command like this

```bash
# polycli assumes that there is default login that's been done already
gcloud auth application-default login
polycli signer create --kms GCP --gcp-project-id prj-polygonlabs-devtools-dev --key-id jhilliard-trash
```

## Listing Keys

After creating a few keys in the keystore or KMS, it's useful just to be able to list the keys. If you're using a keystore, the accounts can be listed using this command:

```bash
polycli signer list --keystore /tmp/keystore
```

In the case of GCP KMS, the keyring will need to be provided and the keys can be listed with this command:

```bash
polycli signer list --kms GCP --gcp-project-id prj-polygonlabs-devtools-dev --gcp-keyring-id polycli-keyring
```

## Signing

Currently, polycli only supports signing transactions. Message and arbitrary signing might happen in the future. In order to sign as message we first need to create some transaction data in a JSON file

```json
{
    "from": "0xB41C20404dffA411fd3F5453a9EA4432Da64e70b",
    "to": "0xB41C20404dffA411fd3F5453a9EA4432Da64e70b",
    "gas": "0x8000",
    "gasPrice": "0x30000000",
    "maxFeePerGas": "0x30000000",
    "maxPriorityFeePerGas": "0x30000000",
    "value": "0x1",
    "nonce": "0x0",

    "input": "",

    "chainId": "0x539",
    "accessList": []
}
```

The file format here is defined by [`SendTxArgs`](https://pkg.go.dev/github.com/ethereum/go-ethereum@v1.13.7/signer/core/apitypes#SendTxArgs) in go-ethereum. This is a lower level transaction format and it will require you to manually specify things like `nonce` and `gasPrice`. In other libraries these are often computed for you.
The benefit here is that it gives us the ability to specify lower level options like `accessList` which aren't usually available in other libraries.

Assuming we have valid transaction data in `tx.json` we can sign the transaction three different ways depending on where our private key is stored

### Signing with Hex Key

This is the easiest, but least secure way to sign. In this case, we're providing a private key as a command line argument to polycli and using that to sign the transaction data. The signed transaction is written to `stdout`
```bash
polycli signer sign --private-key $(cat private-key.txt) --data-file tx.json  --chain-id 1337 | jq '.'
```

This is the output that is generated. `signedTx` is the JSON formatted transaction which is readable but not readily usable. The `rawSignedTx` can be directly published.

```json
{
  "rawSignedTx": "02f86c820539808430000000843000000082800094b41c20404dffa411fd3f5453a9ea4432da64e70b0180c080a0978b7e99d4941fddcbfc792632a53bd4ac4b690ae4395d8203ecec9836e53dd8a00e32626e8456afb6e59f1fb2a8835bd647a97fec4d9da6a46ecadbf310b345d6",
  "signedTx": {
    "type": "0x2",
    "chainId": "0x539",
    "nonce": "0x0",
    "to": "0xb41c20404dffa411fd3f5453a9ea4432da64e70b",
    "gas": "0x8000",
    "gasPrice": null,
    "maxPriorityFeePerGas": "0x30000000",
    "maxFeePerGas": "0x30000000",
    "value": "0x1",
    "input": "0x",
    "accessList": [],
    "v": "0x0",
    "r": "0x978b7e99d4941fddcbfc792632a53bd4ac4b690ae4395d8203ecec9836e53dd8",
    "s": "0xe32626e8456afb6e59f1fb2a8835bd647a97fec4d9da6a46ecadbf310b345d6",
    "yParity": "0x0",
    "hash": "0xc03d220111f2d10b6a2b6b22c98e0e7e728a869cac0d3730e33b8bff683d677d"
  }
}
```

### Signing with keystore

Signing with a keystore requires that you specify the `--keystore` location and the `--key-id` which in this case is the address of the key that you'd like to use for signing

```bash
polycli signer sign --keystore /tmp/keystore --key-id 0x58ce4bE73Ee7D0dee75395Ef662e98F91AD2E740 --data-file tx.json --chain-id 1337
```

### Signing with GCP KMS

The syntax for signing with KMS should look familiar.

```bash
# polycli assumes that there is default login that's been done already
gcloud auth application-default login
polycli signer sign --kms GCP --gcp-project-id prj-polygonlabs-devtools-dev --key-id jhilliard-trash --data-file tx.json --chain-id 1337
```

## Flags

```bash
  -h, --help   help for create
```

The command also inherits flags from parent commands.

```bash
      --chain-id uint              The chain id for the transactions.
      --config string              config file (default is $HOME/.polygon-cli.yaml)
      --data-file string           File name holding data to be signed
      --gcp-import-job-id string   The GCP Import Job ID to use when importing a key
      --gcp-key-version int        The GCP crypto key version to use (default 1)
      --gcp-keyring-id string      The GCP Keyring ID to be used (default "polycli-keyring")
      --gcp-location string        The GCP Region to use (default "europe-west2")
      --gcp-project-id string      The GCP Project ID to use
      --key-id string              The id of the key to be used for signing
      --keystore string            Use the keystore in the given folder or file
      --kms string                 AWS or GCP if the key is stored in the cloud
      --pretty-logs                Should logs be in pretty format or JSON (default true)
      --private-key string         Use the provided hex encoded private key
      --type string                The type of signer to use: latest, cancun, london, eip2930, eip155 (default "london")
      --unsafe-password string     A non-interactively specified password for unlocking the keystore
  -v, --verbosity int              0 - Silent
                                   100 Fatal
                                   200 Error
                                   300 Warning
                                   400 Info
                                   500 Debug
                                   600 Trace (default 400)
```

## See also

- [polycli signer](polycli_signer.md) - Utilities for security signing transactions