This function is meant to help handle ENR data. Given an input ENR it will output the parsed enode and other values that are part of the payload.

The command below will take an ENR and process it:
```bash
echo 'enr:-IS4QHCYrYZbAKWCBRlAy5zzaDZXJBGkcnh4MHcBFZntXNFrdvJjX04jRzjzCBOonrkTfj499SZuOh8R33Ls8RRcy5wBgmlkgnY0gmlwhH8AAAGJc2VjcDI1NmsxoQPKY0yuDUmstAHYpMa2_oxVtw0RW_QAdpzBQA8yWM0xOIN1ZHCCdl8' | \
    polycli enr | jq '.'
```

This is the output:
```json
{
  "enode": "enode://ca634cae0d49acb401d8a4c6b6fe8c55b70d115bf400769cc1400f3258cd31387574077f301b421bc84df7266c44e9e6d569fc56be00812904767bf5ccd1fc7f@127.0.0.1:0?discport=30303",
  "enr": "enr:-IS4QHCYrYZbAKWCBRlAy5zzaDZXJBGkcnh4MHcBFZntXNFrdvJjX04jRzjzCBOonrkTfj499SZuOh8R33Ls8RRcy5wBgmlkgnY0gmlwhH8AAAGJc2VjcDI1NmsxoQPKY0yuDUmstAHYpMa2_oxVtw0RW_QAdpzBQA8yWM0xOIN1ZHCCdl8",
  "id": "a448f24c6d18e575453db13171562b71999873db5b286df957af199ec94617f7",
  "ip": "127.0.0.1",
  "tcp": "0",
  "udp": "30303"
}
```

This command can be used a few different ways
```bash
enr_data="enr:-IS4QHCYrYZbAKWCBRlAy5zzaDZXJBGkcnh4MHcBFZntXNFrdvJjX04jRzjzCBOonrkTfj499SZuOh8R33Ls8RRcy5wBgmlkgnY0gmlwhH8AAAGJc2VjcDI1NmsxoQPKY0yuDUmstAHYpMa2_oxVtw0RW_QAdpzBQA8yWM0xOIN1ZHCCdl8"

# First form - reading from stdin
echo "$enr_data" | polycli enr

# Second form - reading from file
tmp_file="$(mktemp)"
echo "$enr_data" > "$tmp_file" 
polycli enr --file "$tmp_file"

# Third form - command line args
polycli enr "$enr_data" 
```

All three forms support multiple lines. Each line will be convert into a JSON object and printed.