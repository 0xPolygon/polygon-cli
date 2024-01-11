# `polycli mnemonic`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Generate a BIP39 mnemonic seed.

```bash
polycli mnemonic [flags]
```

## Flags

```bash
  -h, --help              help for mnemonic
      --language string   Which language to use [ChineseSimplified, ChineseTraditional, Czech, English, French, Italian, Japanese, Korean, Spanish] (default "english")
      --words int         The number of words to use in the mnemonic (default 24)
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
