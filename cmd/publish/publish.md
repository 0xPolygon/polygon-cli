This command publish transactions with high-throughput.

The command accepts a list of rlp hex encoded transactions that can be provided via a file, 
command line or stdin.

Internally it uses a worker pool strategy that can be dimensioned via flag, so it can be adjusted 
for optimal performance depending on the hardware available.

Since this command focus on high-throughput, please ensure the RPC will not rate-limit the requests.

Below are some example of how to use it

File: to use a file, set the file path using the --file flag
```bash
polycli publis --rpc-url https://sepolia.drpc.org --file /home/tclemos/txs
```

Command Line: to use command line args, set as many args you need when calling the command
```bash
polycli publis --rpc-url https://sepolia.drpc.org 0x000...001 0x000...002 0x000...003 0x000...004 ...
```

Stdin: to use std int, run the command without file or 0x args and then type one tx rlp per line
```bash
polycli cdk rollup monitor --rpc-url https://sepolia.drpc.org


```
