# `polycli signer create`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Create a new key.

```bash
polycli signer create [flags]
```

## Usage

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

## Flags

```bash
  -h, --help   help for create
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
