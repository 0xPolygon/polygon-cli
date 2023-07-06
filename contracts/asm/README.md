These are lower level contract implementations meant to push the limits
of the EVM and network

```bash
./build/bin/evm compile ~/code/polygon-cli/contracts/asm/noop-loop.easm > noop-loop.bin
./build/bin/evm --codefile noop-loop.bin --gas 100000 --debug --json --dump run

./build/bin/evm compile ~/code/polygon-cli/contracts/asm/delegate-call-loop.easm > delegate-call-loop.bin
./build/bin/evm --codefile delegate-call-loop.bin --gas 100000 --debug --json --dump --prestate init.json run

./build/bin/evm compile ~/code/polygon-cli/contracts/asm/sstore-loop.easm > sstore-loop.bin
./build/bin/evm --codefile sstore-loop.bin --gas 100000 --debug --json --dump run

./build/bin/evm compile ~/code/polygon-cli/contracts/asm/fib.easm > fib.bin
./build/bin/evm --codefile fib.bin --gas 100000 --debug --json --dump run


./build/bin/evm compile ~/code/polygon-cli/contracts/asm/fib-nostore.easm > fib-nostore.bin
./build/bin/evm --codefile fib-nostore.bin --gas 100000 --debug --json --dump run



cat noop-loop.bin | tr -d "\n" | wc
./build/bin/evm compile ~/code/polygon-cli/contracts/asm/deploy-header.easm
```

```javascript
eth.coinbase = eth.accounts[0];
eth.sendTransaction({
  from: eth.coinbase,
  to: "0x85da99c8a7c2c95964c8efd687e95e632fc533d6",
  value: web3.toWei(5000, "ether"),
});

loopCode = "0x6014600c60003960146000f360005b6001018062065b9a116300000002575000";
txHash = eth.sendTransaction({ from: eth.coinbase, data: loopCode });
loopCodeReceipt = eth.getTransactionReceipt(txHash);

eth.getCode(loopCodeReceipt.contractAddress);

txHash = eth.sendTransaction({
  from: eth.coinbase,
  to: loopCodeReceipt.contractAddress,
});
eth.getTransaction(txHash);

debug.traceCall(
  { from: eth.coinbase, to: loopCodeReceipt.contractAddress },
  "latest"
);

eth.getTransactionReceipt(loopCodeReceipt);

delegateCode =
  "0x6068600c60003960686000f360005b600101600080808073d2581362bbd7c8ad4ab412068198cde1a8a9bd3b62070000f4508062065b9a116300000002575000";
txHash = eth.sendTransaction({ from: eth.coinbase, data: delegateCode });
delegateCodeReceipt = eth.getTransactionReceipt(txHash);
txHash = eth.sendTransaction({
  from: eth.coinbase,
  to: delegateCodeReceipt.contractAddress,
  gas: 100000,
});

debug.traceCall(
  { from: eth.coinbase, to: delegateCodeReceipt.contractAddress },
  "latest"
);

debug.traceCall(
  { from: eth.coinbase, to: delegateCodeReceipt.contractAddress, gas: 100000 },
  "latest"
);
```
