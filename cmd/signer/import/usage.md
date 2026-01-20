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
