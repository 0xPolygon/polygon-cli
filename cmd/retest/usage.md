The goal of this tool is to take test [fillers](https://github.com/ethereum/tests/tree/develop/src) 
from the standard ethereum test library and convert them into a format that
works well with other tools like [cast](https://book.getfoundry.sh/cast/).


To try this out, first checkout https://github.com/ethereum/tests/

```bash
# Move into the filler directory
cd tests/src

# Convert the yaml based tests to json. There will be some failures depending on the version of yq used
find . -type f -name '*.yml'  | while read yaml_file ; do
    yq '.' $yaml_file > $yaml_file.json
    retval=$?
    if [[ $retval -ne 0 ]]; then
        2>&1 echo "the file $yaml_file could not be converted to json"
    fi
done


# Check for duplicates... There are a few so we should be mindful of that
find . -type f -name '*.json' | xargs cat | jq -r 'to_entries[].key' | uniq -c | sort

# Consolidate everything.. The kzg tests seem to be a different format.. So excluding them with the array check
find . -type f -name '*.json' | xargs cat | jq 'select((. | type) != "array")' | jq -s 'add' > merged.json
# there are some fields like "//comment" that make parsing very difficult
jq 'walk(if type == "object" then with_entries(select(.key | startswith("//") | not)) else . end)' merged.json  > merged.nocomment.json

rm -f single-transactions.txt
find . -type f -name '*.json'  | while read json_file ; do
    jq --arg fname $json_file 'to_entries[].value.transaction | select(.!= null) | .fname = $fname' $json_file >> single-transactions.txt
    jq --arg fname $json_file 'to_entries[].value.blocks | select(. != null) | .[].transactions | select(. != null) | .[] | .fname = $fname' $json_file >> single-transactions.txt
    jq --arg fname $json_file 'to_entries[].value.txbytes | select( . != null) | {txbytes: ., fname: $fname} ' $json_file >> single-transactions.txt
done
cat single-transactions.txt | jq -s '.' > consolidated.json

```

Now we should have a giant file filled with an array of transactions

```bash
cat consolidated.json | jq '.[] | select(( .value | type) == "array") | select((.value | length) > 1)'
cat single-transactions.txt | jq -r '.txbytes | select( . != null)' | xargs -I xxx cast publish --rpc-url http://34.175.214.161:18124 xxx
```


```bash
```
