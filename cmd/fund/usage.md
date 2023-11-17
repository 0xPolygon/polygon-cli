```bash
$ polycli fund \
  --wallet-count=5 \
  --funding-wallet-pk="REPLACE" \
  --chain-id=100 \
  --concurrency=5 \
  --rpc-url="https://rootchain-devnetsub.zkevmdev.net"  \
  --wallet-funding-amt=0.00015 \
  --wallet-funding-gas=50000 \
  --output-file="/opt/generated_keys.json"
  --verbosity=true
```