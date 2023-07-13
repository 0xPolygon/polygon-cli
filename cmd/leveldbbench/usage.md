This command is meant to give us a sense of the system level performance for leveldb.

```bash
go run main.go leveldbbench --degree-of-parallelism 2 | jq '.' > result.json
```


