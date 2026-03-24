package iptp

import (
    "encoding/json"
    "fmt"
    "sync"
    "testing"
    "time"
)

func TestMarshalVerifyRoundTrip(t *testing.T) {
    env := Envelope{
        ID: "msg-1",
        Signal: Signal{
            Intention: "demo",
            Pulses: []Pulse{{Name: "p", Value: "Y"}},
        },
    }
    raw, err := MarshalWithChecksum(&env)
    if err != nil {
        t.Fatal(err)
    }

    var got Envelope
    if err := json.Unmarshal(raw, &got); err != nil {
        t.Fatal(err)
    }
    if err := VerifyEnvelope(&got); err != nil {
        t.Fatal(err)
    }
}

func TestMultiplexProtocolAck(t *testing.T) {
    addr := "127.0.0.1:19109"

    go func() {
        _ = Listen(addr, func(env Envelope) error { return nil })
    }()
    time.Sleep(150 * time.Millisecond)

    client, err := Dial(addr)
    if err != nil {
        t.Fatal(err)
    }
    defer client.Close()

    const n = 5
    var wg sync.WaitGroup
    errs := make(chan error, n)

    for i := 0; i < n; i++ {
        i := i
        wg.Add(1)
        go func() {
            defer wg.Done()
            env := Envelope{
                ID: fmt.Sprintf("msg-%d", i),
                Signal: Signal{
                    Intention: "order:create",
                    Pulses: []Pulse{{Name: "idx", Value: fmt.Sprintf("%d", i)}},
                },
            }
            errs <- client.Send(env, SendOptions{
                RequireProtocolAck: true,
                AckTimeout:         2 * time.Second,
                MaxRetries:         1,
            })
        }()
    }

    wg.Wait()
    close(errs)

    for err := range errs {
        if err != nil {
            t.Fatalf("unexpected send error: %v", err)
        }
    }
}
