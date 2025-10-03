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
      --config string      config file (default is $HOME/.polygon-cli.yaml)
      --pretty-logs        output logs in pretty format instead of JSON (default true)
  -v, --verbosity string   log level (string or int):
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

- [polycli signer](polycli_signer.md) - Utilities for security signing transactions
