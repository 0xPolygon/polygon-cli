For various reasons, we might want to dump a large range of blocks for analytics or replay purposes. This is a simple util to export over RPC a range of blocks.

The following command would download the first 500K blocks and zip them and then look for blocks with transactions that create an account.

```bash
$ polycli dumpblocks 0 500000 --rpc-url http://172.26.26.12:8545/ | gzip > foo.gz
$ zcat < foo.gz | jq '. | select(.transactions | length > 0) | select(.transactions[].to == null)'
```

Dumpblocks can also output to protobuf format.

If you wish to make changes to the protobuf.

1. Install the protobuf compiler

On GNU/Linux:

```bash
$ sudo apt install protoc-gen-go
```

On a MAC:

```bash
$ brew install protobuf
```

2. Install the protobuf plugin

```bash
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

3. Compile the proto file

```bash
$ make generate
```

4. Depending on what endpoint and chain you're querying, you may be missing some fields in the proto definition.

```json
{
  "level": "error",
  "error": "proto: (line 1:813): unknown field \"nonce\"",
  "time": "2023-01-17T13:35:53-05:00",
  "message": "failed to unmarshal proto message"
}
```

To solve this, add the unknown fields to the `.proto` files and recompile them (step 3).
