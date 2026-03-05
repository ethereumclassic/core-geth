#!/bin/bash
# Run core-geth on ETC mainnet (chain ID 61)
# Ports: 8545 (HTTP), 8546 (WS), 30303 (P2P)

exec geth \
  --classic \
  --datadir=/media/dev/2tb/data/blockchain/core-geth/classic \
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
