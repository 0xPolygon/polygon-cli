# `polycli signer list`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

List the keys in the keyring / keystore.

```bash
polycli signer list [flags]
```

## Usage

After creating a few keys in the keystore or KMS, it's useful just to be able to list the keys. If you're using a keystore, the accounts can be listed using this command:

```bash
polycli signer list --keystore /tmp/keystore
```

In the case of GCP KMS, the keyring will need to be provided and the keys can be listed with this command:

```bash
polycli signer list --kms GCP --gcp-project-id prj-polygonlabs-devtools-dev --gcp-keyring-id polycli-keyring
```

## Flags

```bash
  -h, --help   help for list
```

The command also inherits flags from parent commands.

```bash
      --chain-id uint              chain ID for transactions
      --config string              config file (default is $HOME/.polygon-cli.yaml)
      --data-file string           file name holding data to be signed
      --gcp-import-job-id string   GCP import job ID to use when importing key
      --gcp-key-version int        GCP crypto key version to use (default 1)
      --gcp-keyring-id string      GCP keyring ID to be used (default "polycli-keyring")
      --gcp-location string        GCP region to use (default "europe-west2")
      --gcp-project-id string      GCP project ID to use
      --key-id string              ID of key to be used for signing
      --keystore string            use keystore in given folder or file
      --kms string                 AWS or GCP if key is stored in cloud
      --pretty-logs                output logs in pretty format instead of JSON (default true)
      --private-key string         use provided hex encoded private key
      --type string                type of signer to use: latest, cancun, london, eip2930, eip155 (default "london")
      --unsafe-password string     non-interactively specified password for unlocking keystore
  -v, --verbosity string           log level (string or int):
                                     0   - silent
                                     100 - panic
                                     200 - fatal
                                     300 - error
                                     400 - warn
                                     500 - info (default)
                                     600 - debug
                                     700 - trace (default "info")
```

## See also

- [polycli signer](polycli_signer.md) - Utilities for security signing transactions.
