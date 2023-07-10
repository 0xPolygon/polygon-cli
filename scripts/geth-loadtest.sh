#!/bin/bash

# Start geth node
echo Running geth in dev mode in background
make geth &

# Waiting for the node to start
sleep 5

# Start load testing
make geth-loadtest
