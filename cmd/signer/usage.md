```bash
polycli signer create
polycli signer create --keystore /tmp/keystore
polycli signer create --kms GCP --gcp-project-id prj-polygonlabs-devtools-dev --key-id jhilliard-trash
polycli signer sign --kms GCP --gcp-project-id prj-polygonlabs-devtools-dev --key-id jhilliard-trash --data-file foo.json
polycli signer sign  --keystore /tmp/keystore  --key-id 0x58ce4bE73Ee7D0dee75395Ef662e98F91AD2E740 --data-file foo.json
```