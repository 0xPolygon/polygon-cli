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

- [polycli signer](polycli_signer.md) - Utilities for security signing transactions.
