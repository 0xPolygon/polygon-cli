# `polycli signer`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Utilities for security signing transactions

## Usage

This command is meant to allow for easy creation of signed transactions. A raw transaction can then be published with a call to
[`eth_sendRawTransaction`](https://ethereum.org/en/developers/docs/apis/json-rpc/#eth_sendrawtransaction) or using [`cast publish`](https://book.getfoundry.sh/reference/cast/cast-publish).

## Flags

```bash
      --chain-id uint              The chain id for the transactions.
      --data-file string           File name holding data to be signed
      --gcp-import-job-id string   The GCP Import Job ID to use when importing a key
      --gcp-key-version int        The GCP crypto key version to use (default 1)
      --gcp-keyring-id string      The GCP Keyring ID to be used (default "polycli-keyring")
      --gcp-location string        The GCP Region to use (default "europe-west2")
      --gcp-project-id string      The GCP Project ID to use
  -h, --help                       help for signer
      --key-id string              The id of the key to be used for signing
      --keystore string            Use the keystore in the given folder or file
      --kms string                 AWS or GCP if the key is stored in the cloud
      --private-key string         Use the provided hex encoded private key
      --type string                The type of signer to use: latest, cancun, london, eip2930, eip155 (default "london")
      --unsafe-password string     A non-interactively specified password for unlocking the keystore
```

The command also inherits flags from parent commands.

```bash
      --config string   config file (default is $HOME/.polygon-cli.yaml)
      --pretty-logs     Should logs be in pretty format or JSON (default true)
  -v, --verbosity int   0 - Silent
                        100 Panic
                        200 Fatal
                        300 Error
                        400 Warning
                        500 Info
                        600 Debug
                        700 Trace (default 500)
```

## See also

- [polycli](polycli.md) - A Swiss Army knife of blockchain tools.
- [polycli signer create](polycli_signer_create.md) - Create a new key

- [polycli signer import](polycli_signer_import.md) - Import a private key into the keyring / keystore

- [polycli signer list](polycli_signer_list.md) - List the keys in the keyring / keystore

- [polycli signer sign](polycli_signer_sign.md) - Sign tx data

