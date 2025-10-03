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
      --chain-id uint              chain ID for transactions
      --data-file string           file name holding data to be signed
      --gcp-import-job-id string   GCP import job ID to use when importing key
      --gcp-key-version int        GCP crypto key version to use (default 1)
      --gcp-keyring-id string      GCP keyring ID to be used (default "polycli-keyring")
      --gcp-location string        GCP region to use (default "europe-west2")
      --gcp-project-id string      GCP project ID to use
  -h, --help                       help for signer
      --key-id string              ID of key to be used for signing
      --keystore string            use keystore in given folder or file
      --kms string                 AWS or GCP if key is stored in cloud
      --private-key string         use provided hex encoded private key
      --type string                type of signer to use: latest, cancun, london, eip2930, eip155 (default "london")
      --unsafe-password string     non-interactively specified password for unlocking keystore
```

The command also inherits flags from parent commands.

```bash
      --config string   config file (default is $HOME/.polygon-cli.yaml)
      --pretty-logs     output logs in pretty format instead of JSON (default true)
  -v, --verbosity int   0 - silent
                        100 panic
                        200 fatal
                        300 error
                        400 warning
                        500 info
                        600 debug
                        700 trace (default 500)
```

## See also

- [polycli](polycli.md) - A Swiss Army knife of blockchain tools.
- [polycli signer create](polycli_signer_create.md) - Create a new key

- [polycli signer import](polycli_signer_import.md) - Import a private key into the keyring / keystore

- [polycli signer list](polycli_signer_list.md) - List the keys in the keyring / keystore

- [polycli signer sign](polycli_signer_sign.md) - Sign tx data

