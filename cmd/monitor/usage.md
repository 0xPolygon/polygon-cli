![GIF of `polycli monitor`](assets/monitor.gif)

This is a basic tool for monitoring block production on a JSON RPC endpoint.

If you're using the terminal UI and you'd like to be able to select text for copying, you might need to use a modifier key.

```bash
$ polycli monitor https://polygon-rpc.com
```

If you're experiencing missing blocks, try adjusting the `--batch-size` and `--interval` flags so that you poll for more blocks or more frequently.
