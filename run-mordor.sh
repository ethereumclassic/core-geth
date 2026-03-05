#!/bin/bash
# Run core-geth on Mordor testnet (ETC, chain ID 63)
# Ports: 8545 (HTTP), 8546 (WS), 30303 (P2P)

exec geth \
  --mordor \
  --datadir=/media/dev/2tb/data/blockchain/core-geth/mordor \
  --http \
  --http.addr=0.0.0.0 \
  --http.port=8545 \
  --http.corsdomain="*" \
  --http.api=eth,net,web3,debug,txpool \
  --ws \
  --ws.addr=0.0.0.0 \
  --ws.port=8546 \
  --port=30303 \
  --cache=1024 \
  --verbosity=3 \
  "$@"
