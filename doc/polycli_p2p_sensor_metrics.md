
## Sensor Metrics


### sensor_block_range
Difference between head and oldest block numbers

Metric Type: Gauge


### sensor_broadcast_batch_size
Number of transactions per broadcast batch

Metric Type: Histogram


### sensor_broadcast_queue_depth
Number of transaction batches in broadcast queue

Metric Type: Gauge


### sensor_broadcast_send_errors
Number of failed broadcast sends

Metric Type: Counter


### sensor_broadcast_txs_rejected
Number of transactions dropped from the broadcast path, by rejection reason

Metric Type: CounterVec

Variable Labels:
- reason


### sensor_broadcast_txs_validated
Number of unique transactions checked before rebroadcast, by outcome

Metric Type: CounterVec

Variable Labels:
- result


### sensor_head_block_age
Time since head block was received (in seconds)

Metric Type: Gauge


### sensor_head_block_number
Current head block number

Metric Type: Gauge


### sensor_head_block_timestamp
Head block timestamp in Unix epoch seconds

Metric Type: Gauge


### sensor_messages
Number and type of messages the sensor has sent and received

Metric Type: CounterVec

Variable Labels:
- message
- direction


### sensor_oldest_block_number
Oldest block number (floor for parent fetching)

Metric Type: Gauge


### sensor_peers
Number of peers the sensor is connected to

Metric Type: Gauge


### sensor_rpc_requests
Number of RPC requests made

Metric Type: CounterVec

Variable Labels:
- method
- proxied


### sensor_tx_decode_errors
Number of transactions received from peers that failed to decode

Metric Type: Counter

