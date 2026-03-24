# IPTP Pipe (TCP + framing + checksum + optional protocol ACK + multiplexing)

IPTP Pipe transports **one Signal per envelope**, where a Signal is:

> **Intention + Pulse set**

This repo keeps IPTP independent of Intention Space.
It does **not** evaluate field state, gate conditions, or design-time rules.

## What this version adds

- framed TCP transport
- CRC32 checksum
- optional `RequireProtocolAck`
- multiplexed in-flight sends over one TCP connection
- async ACK/NACK correlation by `Envelope.ID`

## Repo layout

```text
iptp-pipe/
├── iptp/
│   ├── client.go
│   ├── pipe.go
│   ├── types.go
│   └── iptp_test.go
└── cmd/
    ├── iptp-recv/
    │   └── main.go
    └── iptp-send/
        └── main.go
```

## Protocol meaning

### TCP already gives
- ordered bytes
- retransmission
- transport integrity

### IPTP adds
- message boundaries
- one complete signal per envelope
- protocol-level ACK/NACK if requested

### Protocol ACK means only
- frame parsed
- checksum valid
- envelope structurally accepted

It does **not** mean:
- workflow accepted
- field conditions satisfied
- Intention Space semantics passed

## Run receiver

```bash
go run ./cmd/iptp-recv -port 9000
```

## Send one signal without protocol ACK

```bash
go run ./cmd/iptp-send -to localhost:9000 -id msg-1 -intention order:create
```

## Send one signal with protocol ACK

```bash
go run ./cmd/iptp-send -to localhost:9000 -id msg-2 -intention order:create -ack
```

## Multiplex demo

This sends several signals concurrently over one connection, each with its own ID:

```bash
go run ./cmd/iptp-send -to localhost:9000 -intention order:create -ack -multi 5
```

## ACK/NACK examples

ACK:
```json
{
  "id": "ack-msg-2",
  "signal": {
    "intention": "iptp:ack",
    "pulses": [
      {"name":"ack_for","value":"msg-2"},
      {"name":"status","value":"accepted"}
    ]
  }
}
```

NACK:
```json
{
  "id": "nack-msg-2",
  "signal": {
    "intention": "iptp:nack",
    "pulses": [
      {"name":"ack_for","value":"msg-2"},
      {"name":"reason","value":"checksum_failed"}
    ]
  }
}
```

## Why multiplexing still fits the model

Multiplexing does **not** change IPTP semantics because:

- each envelope still carries exactly **one Intention**
- each envelope still carries its own **Pulse set**
- concurrency is at the transport layer
- interpretation remains outside IPTP

So one connection may carry many complete semantic units, but IPTP never merges intentions.
# iptp-pipe
