# `polycli parseethwallet`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Extract the private key from an eth wallet.

```bash
polycli parseethwallet [flags]
```

## Usage

This function can take a geth style wallet file and extract the private key as hex. It can also do the opposite.

This command takes the private key and imports it into a local keystore with no password.

```bash
$ polycli parseethwallet --hexkey 42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa

$ cat UTC--2023-05-09T22-48-57.582848385Z--85da99c8a7c2c95964c8efd687e95e632fc533d6  | jq '.'
{
  "address": "85da99c8a7c2c95964c8efd687e95e632fc533d6",
  "crypto": {
    "cipher": "aes-128-ctr",
    "ciphertext": "d0b4377a4ae5ebc9a5bef06ce4be99565d10cb0dedc2f7ff5aaa07ea68e7b597",
    "cipherparams": {
      "iv": "8ecd172ff7ace15ed5bc44ea89473d8e"
    },
    "kdf": "scrypt",
    "kdfparams": {
      "dklen": 32,
      "n": 262144,
      "p": 1,
      "r": 8,
      "salt": "cd6ec772dc43225297412809feaae441d578642c6a67cabf4e29bcaf594f575b"
    },
    "mac": "c992128ed466ad15a9648f4112af22929b95f511f065b12a80abcfb7e4d39a79"
  },
  "id": "82af329d-2af5-41a6-ae6b-624f3e1c224b",
  "version": 3
}
```

If we wanted to go the opposite direction, we could run a command like this.

```bash
polycli parseethwallet --file /tmp/keystore/UTC--2023-05-09T22-48-57.582848385Z--85da99c8a7c2c95964c8efd687e95e632fc533d6  | jq '.'
{
  "Address": "0x85da99c8a7c2c95964c8efd687e95e632fc533d6",
  "PublicKey": "507cf9a75e053cda6922467721ddb10412da9bec30620347d9529cc77fca24334a4cf59685be4a2fdeabf4e7753350e42d2d3a20250fd9dc554d226463c8a3d5",
  "PrivateKey": "42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa"
}
```

## Flags

```bash
      --file string       Provide a file with the key information 
  -h, --help              help for parseethwallet
      --hexkey string     An optional hexkey that would be use to generate a geth style key
      --keystore string   The directory where keys would be stored when importing a raw hex (default "/tmp/keystore")
      --password string   An optional password use to unlock the key
```

The command also inherits flags from parent commands.

```bash
      --config string   config file (default is $HOME/.polygon-cli.yaml)
      --pretty-logs     Should logs be in pretty format or JSON (default true)
  -v, --verbosity int   0 - Silent
                        100 Panic
                        200 Fatal
                        300 Error
                        400 Warning
                        500 Info
                        600 Debug
                        700 Trace (default 500)
```

## See also

- [polycli](polycli.md) - A Swiss Army knife of blockchain tools.
