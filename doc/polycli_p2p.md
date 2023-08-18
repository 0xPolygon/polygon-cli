# `polycli p2p`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Set of commands related to devp2p.

## Usage

Pinging a peer is useful to determine information about the peer and retrieving the `Hello` and `Status` messages. By default, it will listen to the peer after the status exchange for blocks and transactions. To disable this behavior, set the `--listen` flag.

```bash
$ polycli p2p ping <enode/enr or nodes.json file>
```

Running the sensor will do peer discovery and continue to watch for blocks and transactions from those peers. This is useful for observing the network for forks and reorgs without the need to run the entire full node infrastructure.

```bash
$ polycli p2p sensor nodes.json --bootnodes enode://0cb82b395094ee4a2915e9714894627de9ed8498fb881cec6db7c65e8b9a5bd7f2f25cc84e71e89d0947e51c76e85d0847de848c7782b13c0255247a6758178c@44.232.55.71:30303,enode://88116f4295f5a31538ae409e4d44ad40d22e44ee9342869e7d68bdec55b0f83c1530355ce8b41fbec0928a7d75a5745d528450d30aec92066ab6ba1ee351d710@159.203.9.164:30303,enode://4be7248c3a12c5f95d4ef5fff37f7c44ad1072fdb59701b2e5987c5f3846ef448ce7eabc941c5575b13db0fb016552c1fa5cca0dda1a8008cf6d63874c0f3eb7@3.93.224.197:30303,enode://32dd20eaf75513cf84ffc9940972ab17a62e88ea753b0780ea5eca9f40f9254064dacb99508337043d944c2a41b561a17deaad45c53ea0be02663e55e6a302b2@3.212.183.151:30303 --network-id 137 --sensor-id "sensor" --project-id "devtools-sandbox"
```

To crawl the network for nodes and write the output json to a file. This will not engage in block or transaction propagation, but it can give a good indicator of network size, and the output json can be used to quick start other nodes.

```bash
$ polycli p2p crawl nodes.json --bootnodes enode://0cb82b395094ee4a2915e9714894627de9ed8498fb881cec6db7c65e8b9a5bd7f2f25cc84e71e89d0947e51c76e85d0847de848c7782b13c0255247a6758178c@44.232.55.71:30303,enode://88116f4295f5a31538ae409e4d44ad40d22e44ee9342869e7d68bdec55b0f83c1530355ce8b41fbec0928a7d75a5745d528450d30aec92066ab6ba1ee351d710@159.203.9.164:30303,enode://4be7248c3a12c5f95d4ef5fff37f7c44ad1072fdb59701b2e5987c5f3846ef448ce7eabc941c5575b13db0fb016552c1fa5cca0dda1a8008cf6d63874c0f3eb7@3.93.224.197:30303,enode://32dd20eaf75513cf84ffc9940972ab17a62e88ea753b0780ea5eca9f40f9254064dacb99508337043d944c2a41b561a17deaad45c53ea0be02663e55e6a302b2@3.212.183.151:30303 --network-id 137
```

## Flags

```bash
  -h, --help   help for p2p
```

The command also inherits flags from parent commands.

```bash
      --config string   config file (default is $HOME/.polygon-cli.yaml)
      --pretty-logs     Should logs be in pretty format or JSON (default true)
  -v, --verbosity int   0 - Silent
                        100 Fatal
                        200 Error
                        300 Warning
                        400 Info
                        500 Debug
                        600 Trace (default 400)
```

## See also

- [polycli](polycli.md) - A Swiss Army knife of blockchain tools.
- [polycli p2p crawl](polycli_p2p_crawl.md) - Crawl a network on the devp2p layer and generate a nodes JSON file.

- [polycli p2p nodelist](polycli_p2p_nodelist.md) - Generate a node list to seed a node

- [polycli p2p ping](polycli_p2p_ping.md) - Ping node(s) and return the output.

- [polycli p2p sensor](polycli_p2p_sensor.md) - Start a devp2p sensor that discovers other peers and will receive blocks and transactions.

