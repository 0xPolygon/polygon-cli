#!/bin/bash
# This script takes an argument of the dest directory for the rpc types json file
# Usage: ./rpctypes.sh rpctypes/schemas/
readonly url="https://github.com/ethereum/execution-apis.git"
readonly commit_id="0c18fb0"
readonly dest="tmp/execution-apis"
readonly schema_dest="schemas"

rm -rf "./$dest"
git clone "$url" "$dest"
pushd "$dest" || exit
git checkout "$commit_id"

npm install
npm run build

# shellcheck disable=SC2207
methods=($(jq -r '.methods[].name' openrpc.json | sort))

mkdir "$schema_dest"
echo "Methods:"
for method in "${methods[@]}"; do
  echo "Generating schemas for: $method"
  jq --arg methodName "$method" '.methods[] | select(.name == $methodName) | .result.schema' openrpc.json > "$schema_dest/$method.json"
done
popd || exit

mkdir -p "./$1"
echo "Copying schemas to $1..."
cp -f $dest/$schema_dest/* "$1"
