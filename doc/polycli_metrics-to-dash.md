# `polycli metrics-to-dash`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Create a dashboard from an Openmetrics / Prometheus response.

```bash
polycli metrics-to-dash [flags]
```

## Usage

Here is how you can use this command.

```bash
$ polycli metrics-to-dash -i avail-metrics.txt -p avail. -t "Avail Devnet Dashboard" -T basedn -D devnet01.avail.polygon.private -T host -D validator-001 -s substrate_ -s sub_ -P true -S true

$ polycli metrics-to-dash -i avail-light-metrics.txt -p avail_light. -t "Avail Light Devnet Dashboard" -T basedn -D devnet01.avail.polygon.private -T host -D validator-001 -s substrate_ -s sub_ -P true -S true
```

## Flags

```bash
  -d, --desc string                         description for the dashboard (default "Polycli Dashboard")
  -H, --height int                          widget height (default 3)
  -h, --help                                help for metrics-to-dash
  -i, --input-file string                   the metrics file to be used
  -p, --prefix string                       prefix to use before all metrics
  -P, --pretty-name                         Should the metric names be prettified (default true)
  -S, --show-help                           Should we show the help text for each metric
  -s, --strip-prefix stringArray            A prefix that can be removed from the metrics
  -D, --template-var-defaults stringArray   The defaults to use for the template variables
  -T, --template-vars stringArray           The template variables to use for the dashboard
  -t, --title string                        title for the dashboard (default "Polycli Dashboard")
  -W, --width int                           widget width (default 4)
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
