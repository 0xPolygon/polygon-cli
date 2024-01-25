# `polycli signer import`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Import a private key into the keyring / keystore

```bash
polycli signer import [flags]
```

## Usage

It's possible to import a simple hex encoded private key into a local
keystore or as a crypto key version in GCP KMS.

In order to import into a local keystore, a command like this could be used:

```bash
polycli signer import --keystore /tmp/keystore --private-key cf42d151cec45693f2ac1201e803b056c5f9e2e5d1af627ce41ab3b6faceda25
```

### Importing into GCP KMS

Importing into GCP KMS is a little bit more complicated. In order to run the import, a command like this would be used:

```bash
polycli signer import --private-key 42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa --kms gcp --gcp-project-id prj-polygonlabs-devtools-dev --key-id jhilliard-code-quality --gcp-import-job-id test-import-job
```

There are a few things going on here:

1. We're specifying a `--private-key` that's hex encoded. That's the only format that `import` accepts at this time.
2. We're using `--kms gcp` which tell polycli that we want to use gcp kms as our backend
3. We've specified `--gcp-project-id` which names a test project that we're using. We've left out `--gcp-location` and `--gcp-keyring-id` which means we're using the defaults.
4. We've set `--gcp-import-job-id test-import-job` which names a job that will be used to import the key. Basically GCP will give us a public key that we use to encrypt our key for importing

The `--key-id` is also important. This will set the name of the key that's going to be imported. Note, the key-id must have already been created in order for import to work. When we're doing the import, we're actually importing a new version of the key that already exists. For the time being, the `signer import` command will not create the first version of the key for you.

## Flags

```bash
  -h, --help   help for import
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
