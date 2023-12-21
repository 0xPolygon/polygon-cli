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

TODO
## Flags

```bash
  -h, --help   help for create
```

The command also inherits flags from parent commands.

```bash
      --chain-id uint            The chain id for the transactions.
      --config string            config file (default is $HOME/.polygon-cli.yaml)
      --data-file string         File name holding data to be signed
      --gcp-keyring-id string    The GCP Keyring ID to be used (default "polycli-keyring")
      --gcp-location string      The GCP Region to use (default "europe-west2")
      --gcp-project-id string    The GCP Project ID to use
      --key-id string            The id of the key to be used for signing
      --keystore string          Use the keystore in the given folder or file
      --kms string               AWS or GCP if the key is stored in the cloud
      --pretty-logs              Should logs be in pretty format or JSON (default true)
      --private-key string       Use the provided hex encoded private key
      --type string              The type of signer to use: latest, cancun, london, eip2930, eip155 (default "london")
      --unsafe-password string   A non-interactively specified password for unlocking the keystore
  -v, --verbosity int            0 - Silent
                                 100 Fatal
                                 200 Error
                                 300 Warning
                                 400 Info
                                 500 Debug
                                 600 Trace (default 400)
```

## See also

- [polycli signer](polycli_signer.md) - Utilities for security signing transactions
