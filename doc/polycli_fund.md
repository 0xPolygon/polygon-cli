# `polycli fund`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Bulk fund crypto wallets automatically.

```bash
polycli fund [flags]
```

## Usage

Bulk fund 30 wallets using a pre-deployed contract address.

```bash
$ polycli fund --verbosity 700 --wallets 30 --funder-address 0x311981b16534238422d301e95a42f6c27a24f346
4:51PM DBG Starting logger in console mode
4:51PM DBG Input parameters params={"FunderAddress":"0x311981b16534238422d301e95a42f6c27a24f346","OutputFile":"wallets.json","PrivateKey":"0x42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa","RpcUrl":"http://localhost:8545","WalletCount":30,"WalletFundingHexAmount":"0xb1a2bc2ec50000"}
4:51PM TRC Detected chain ID chainID=1337
4:51PM DBG List of wallets to be funded addresses=["0x8098e0092875a89d8db66708ec1dd248d2dd4dac","0x27f93e701cf7e278687fb1fa1cc9a30932e17587","0x9174d0b938c10f89787ebbd37905593cb76a26e4","0x50280161aa7656f57c18be4ef558786e2c5510c1","0xf62947ecacd778d888b3a2d01a80a5608257afd2","0x964ffc8352ce1883fe740b02af1ccefbbba31e0f","0x03ccadd9aededaed3f7d0df934914f36a91a3063","0xa1138efc9f93f709a0a01fe0480b6d0a7088a488","0xac1cd574d98d3b6cb11a2628af49e2d0dc0292ef","0xc6b7a7f0b6a8e15f44e75f8c5d84457968bf540e","0x25d9964b6e05c6741d99f25131b734cc6383e6db","0x3960b01c6d4058ebfb61599e7d366af747a78683","0x2a77bc0048442970878cf81ff65654ac42fc8675","0x805309b6f9dc0d4c90a68d66938c00a8538de856","0x08f574385a01e8ac136ee281ee9f8ed01219a50f","0x5440214ab31969f3dbb6b41c1dcdbfee2ccc46e2","0x610f98ba071eea9e76a2d7ea4fd3db28770da77b","0x562a606458dfce441db4f89c2eee8d06da407442","0x6d39a81d0e63ba4749e49fedd05504a2c1b538e7","0x537a338c6be9fffceb2bfe45b320f14d36311c2f","0x410165c205fccf40ef4a78a80542ce69a71630c1","0x3d0fe3f63bf613a22dac5c0c4c7fb686358d64c2","0x67f7dbd3f1942a783c6ec69968d63985f7d3bee4","0xeba6d5c38cba411978dd979fa83b9a757136b328","0xce6a59a028d26149c1c2d93b09fab3725f96196c","0x48a70f9bb831a14b928a17691a16cec5d3e40bc3","0x68a3bd96c015297f560d6e3319742630ab0e17bc","0x5e1c242b5bfea09c5156e7fb978840388ebe2eeb","0xf94c08fe881a0ea565f4642558eb59f3d2df09c8","0x93d1acff7cfffbe390f4e18797b49404500ac123"]
4:51PM INF Wallet addresses and private keys saved to file fileName=wallets.json
4:51PM INF Accounts funded! 💸
4:51PM INF Total execution time: 20.9145ms
```

Extract from `wallets.json`.

```json
[
  {
    "Address": "0x8098e0092875a89d8db66708eC1dD248D2DD4Dac",
    "PrivateKey": "7f025a5ab0a8699ca79495d8158ddbb9a6b471085a92a20ff39a274235499f22"
  },
  {
    "Address": "0x27F93e701cf7e278687FB1fA1cc9A30932E17587",
    "PrivateKey": "95b2d4484fa219f4ca74df41e04ae8557becf707824d04e4cbaab10c410ac983"
  },
  {
    "Address": "0x9174D0B938C10f89787ebBD37905593cB76a26E4",
    "PrivateKey": "8b6b5032e56dcf95e9c3ef317da5bc41b89538af652d3502973c3ccf3a36fd83"
  },
  {
    "Address": "0x50280161AA7656F57c18Be4ef558786E2c5510C1",
    "PrivateKey": "cdba7d43e36981672d5c73af82298526597a44ef20e4fef3d439283c40c2b8a1"
  },
  ...
]
```

Check the balances of the wallets.

```bash
$ cast balance 0x8098e0092875a89d8db66708eC1dD248D2DD4Dac
50000000000000000

$ cast balance 0x27F93e701cf7e278687FB1fA1cc9A30932E17587
50000000000000000
...
```

## Flags

```bash
  -a, --amount string           The amount of wei to send to each wallet (default "0xb1a2bc2ec50000")
  -f, --file string             The output JSON file path for storing the addresses and private keys of funded wallets (default "wallets.json")
      --funder-address string   The address of a pre-deployed Funder contract
  -h, --help                    help for fund
      --private-key string      The hex encoded private key that we'll use to send transactions (default "0x42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa")
  -r, --rpc-url string          The RPC endpoint url (default "http://localhost:8545")
  -w, --wallets uint            The number of wallets to fund (default 10)
```

The command also inherits flags from parent commands.

```bash
      --config string   config file (default is $HOME/.polygon-cli.yaml)
      --pretty-logs     Should logs be in pretty format or JSON (default true)
  -v, --verbosity int   0 - Silent
                        100 Fatal
                        200 Error
                        300 Warning
                        400 Info
                        500 Debug
                        600 Trace (default 400)
```

## See also

- [polycli](polycli.md) - A Swiss Army knife of blockchain tools.
