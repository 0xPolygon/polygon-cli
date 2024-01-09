# `polycli signer list`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

List the keys in the keyring / keystore

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
                                   100 Panic
                                   200 Fatal
                                   300 Error
                                   400 Warning
                                   500 Info
                                   600 Debug
                                   700 Trace (default 500)
```

## See also

- [polycli signer](polycli_signer.md) - Utilities for security signing transactions
