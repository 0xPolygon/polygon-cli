# `polycli retest`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Convert the standard ETH test fillers into something to be replayed against an RPC

```bash
polycli retest [flags]
```

## Usage

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
```

Now we should have a giant file filled with an array of transactions. We can take that output and process it witht the `retest` command now

```bash
go run . retest -v 500 --file merged.nocomment.json > simple.json
```

## LLLC

This project will depend on an installation of `solc` (specifically
0.8.20) and `lllc`. Installing solidity is pretty easy, but LLLC can
be a little tricky.

Since the version is pretty old, it might not build well on your host
os. Building within docker might make your life easier:

```bash
docker run -it debian:buster /bin/bash
```

From within the docker shell some steps like this should get you in
the right direction:

```bash
apt update
apt install --yes libboost-filesystem-dev libboost-system-dev libboost-program-options-dev libboost-test-dev git cmake g++
git clone --depth 1 -b master https://github.com/winsvega/solidity.git /solidity
mkdir /build && cd /build
cmake /solidity -DCMAKE_BUILD_TYPE=Release -DLLL=1 && make lllc
```

Assuming that all worked, we should be able to copy the binary out of
docker and into our host OS:

```bash
docker cp 95511e9d0996:/build/lllc/lllc /usr/local/bin/
```

## Flags

```bash
      --file string   Provide a file that's filed with test transaction fillers
  -h, --help          help for retest
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
