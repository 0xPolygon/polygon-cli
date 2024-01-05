After creating a few keys in the keystore or KMS, it's useful just to be able to list the keys. If you're using a keystore, the accounts can be listed using this command:

```bash
polycli signer list --keystore /tmp/keystore
```

In the case of GCP KMS, the keyring will need to be provided and the keys can be listed with this command:

```bash
polycli signer list --kms GCP --gcp-project-id prj-polygonlabs-devtools-dev --gcp-keyring-id polycli-keyring
```
