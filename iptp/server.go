package iptp

import (
    "bufio"
    "encoding/json"
    "net"
)

type HandlerFunc func(Envelope) error

func Listen(addr string, onSignal HandlerFunc) error {
    ln, err := net.Listen("tcp", addr)
    if err != nil {
        return err
    }
    for {
        conn, err := ln.Accept()
        if err != nil {
            continue
        }
        go handleConn(conn, onSignal)
    }
}

func handleConn(conn net.Conn, onSignal HandlerFunc) {
    defer conn.Close()
    r := bufio.NewReader(conn)

    for {
        payload, err := ReadFrame(r)
        if err != nil {
            return
        }

        var env Envelope
        if err := json.Unmarshal(payload, &env); err != nil {
            sendNack(conn, "", "invalid_json")
            continue
        }
        if env.ID == "" {
            sendNack(conn, "", "missing_id")
            continue
        }
        if env.Signal.Intention == "" {
            sendNack(conn, env.ID, "missing_intention")
            continue
        }
        if err := VerifyEnvelope(&env); err != nil {
            sendNack(conn, env.ID, "checksum_failed")
            continue
        }

        if err := onSignal(env); err != nil {
            sendNack(conn, env.ID, err.Error())
            continue
        }
        sendAck(conn, env.ID)
    }
}

func sendAck(conn net.Conn, id string) {
    env := Envelope{
        ID: "ack-" + id,
        Signal: Signal{
            Intention: "iptp:ack",
            Pulses: []Pulse{
                {Name: "ack_for", Value: id},
                {Name: "status", Value: "accepted"},
            },
        },
    }
    payload, _ := MarshalWithChecksum(&env)
    _ = WriteFrame(conn, payload)
}

func sendNack(conn net.Conn, id, reason string) {
    env := Envelope{
        ID: "nack-" + id,
        Signal: Signal{
            Intention: "iptp:nack",
            Pulses: []Pulse{
                {Name: "ack_for", Value: id},
                {Name: "reason", Value: reason},
            },
        },
    }
    payload, _ := MarshalWithChecksum(&env)
    _ = WriteFrame(conn, payload)
}
