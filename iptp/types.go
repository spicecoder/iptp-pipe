package iptp

import "time"

type Pulse struct {
    Name  string `json:"name"`
    Value string `json:"value"`
}

type Signal struct {
    Intention string  `json:"intention"`
    Pulses    []Pulse `json:"pulses"`
}

type Envelope struct {
    ID        string    `json:"id"`
    SentAt    time.Time `json:"sentAt,omitempty"`
    Signal    Signal    `json:"signal"`
    Checksum  uint32    `json:"checksum"`
}

type SendOptions struct {
    RequireProtocolAck bool
    AckTimeout         time.Duration
    MaxRetries         int
}
