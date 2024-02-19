#!/bin/bash
# This script takes an argument of the dest directory for the rpc types json file
# Usage: ./rpctypes.sh rpctypes/schemas/
readonly url=https://github.com/ethereum/execution-apis.git
readonly dest=tmp/execution-apis
readonly schema_dest=schemas

rm -rf ./$dest
git clone --depth=1 $url $dest

pushd $dest
npm install
npm run build

methods=($(cat openrpc.json | jq -r '.methods[].name' | sort))

mkdir $schema_dest
echo "Methods:"
for method in "${methods[@]}"; do
  echo "Generating schemas for: $method"
  cat openrpc.json | jq --arg methodName $method '.methods[] | select(.name == $methodName) | .result.schema' > "$schema_dest/$method.json"
done
popd

mkdir -p ./$1
echo "Copying schemas to $1..."
cp -f $dest/$schema_dest/* $1
