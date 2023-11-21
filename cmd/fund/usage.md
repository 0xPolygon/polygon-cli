```bash
$ polycli fund \
  --rpc-url="https://rootchain-devnetsub.zkevmdev.net"  \
  --funding-wallet-pk="REPLACE" \
  --wallet-count=5 \
  --wallet-funding-amt=0.00015 \
  --wallet-funding-gas=50000 \
  --concurrency=5 \
  --output-file="/opt/funded_wallets.json"
```