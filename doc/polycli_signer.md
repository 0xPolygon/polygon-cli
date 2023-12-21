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

TODO
## Flags

```bash
      --chain-id uint            The chain id for the transactions.
      --data-file string         File name holding data to be signed
      --gcp-keyring-id string    The GCP Keyring ID to be used (default "polycli-keyring")
      --gcp-location string      The GCP Region to use (default "europe-west2")
      --gcp-project-id string    The GCP Project ID to use
  -h, --help                     help for signer
      --key-id string            The id of the key to be used for signing
      --keystore string          Use the keystore in the given folder or file
      --kms string               AWS or GCP if the key is stored in the cloud
      --private-key string       Use the provided hex encoded private key
      --type string              The type of signer to use: latest, cancun, london, eip2930, eip155 (default "london")
      --unsafe-password string   A non-interactively specified password for unlocking the keystore
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
- [polycli signer create](polycli_signer_create.md) - Create a new key

- [polycli signer sign](polycli_signer_sign.md) - Sign tx data

