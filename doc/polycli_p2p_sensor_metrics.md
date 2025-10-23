
## Sensor Metrics


### sensor_block_range
Difference between head and oldest block numbers

Metric Type: Gauge


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
The number and type of messages the sensor has sent and received

Metric Type: CounterVec

Variable Labels:
- message
- url
- name
- direction


### sensor_oldest_block_number
Oldest block number (floor for parent fetching)

Metric Type: Gauge


### sensor_peers
The number of peers the sensor is connected to

Metric Type: Gauge

