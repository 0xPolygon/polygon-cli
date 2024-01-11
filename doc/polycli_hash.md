# `polycli hash`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Provide common crypto hashing functions.

```bash
polycli hash [md4|md5|sha1|sha224|sha256|sha384|sha512|ripemd160|sha3_224|sha3_256|sha3_384|sha3_512|sha512_224|sha512_256|blake2s_256|blake2b_256|blake2b_384|blake2b_512|keccak256|keccak512] [flags]
```

## Usage

```bash
$ echo -n "hello" > hello.txt
$ polycli hash sha1 --file hello.txt
aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d
$ echo -n "hello" | polycli hash sha1
aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d
$ polycli hash sha1 hello
aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d
```

## Flags

```bash
      --file string   Provide a filename to read and hash
  -h, --help          help for hash
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
