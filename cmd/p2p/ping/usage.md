Pinging a peer is useful to determine information about the peer and retrieving
the `Hello` and `Status` messages. By default, it will listen to the peer after
the status exchange for blocks and transactions. To disable this behavior, set
the `--listen` flag.

## Example

```bash
polycli p2p ping <enode/enr or nodes.json file>
```
