package main

import (
    "flag"
    "fmt"
    "log"
    "sync"
    "time"

    "github.com/spicecoder/iptp-pipe/iptp"
)

func main() {
    to := flag.String("to", "localhost:9000", "receiver host:port")
    id := flag.String("id", "msg-1", "message id")
    intention := flag.String("intention", "order:create", "signal intention")
    requireAck := flag.Bool("ack", false, "require protocol ack")
    multi := flag.Int("multi", 1, "number of concurrent sends over one connection")
    flag.Parse()

    client, err := iptp.Dial(*to)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    opts := iptp.SendOptions{
        RequireProtocolAck: *requireAck,
        AckTimeout:         2 * time.Second,
        MaxRetries:         2,
    }

    if *multi <= 1 {
        env := iptp.Envelope{
            ID: *id,
            Signal: iptp.Signal{
                Intention: *intention,
                Pulses: []iptp.Pulse{
                    {Name: "example", Value: "Y"},
                },
            },
        }
        if err := client.Send(env, opts); err != nil {
            log.Fatal(err)
        }
        fmt.Printf("sent id=%s intention=%s ack=%v\n", env.ID, env.Signal.Intention, *requireAck)
        return
    }

    var wg sync.WaitGroup
    errs := make(chan error, *multi)

    for i := 0; i < *multi; i++ {
        i := i
        wg.Add(1)
        go func() {
            defer wg.Done()
            env := iptp.Envelope{
                ID: fmt.Sprintf("%s-%d", *id, i),
                Signal: iptp.Signal{
                    Intention: *intention,
                    Pulses: []iptp.Pulse{
                        {Name: "example", Value: "Y"},
                        {Name: "index", Value: fmt.Sprintf("%d", i)},
                    },
                },
            }
            errs <- client.Send(env, opts)
        }()
    }

    wg.Wait()
    close(errs)

    failed := false
    for err := range errs {
        if err != nil {
            failed = true
            fmt.Println("send error:", err)
        }
    }
    if failed {
        log.Fatal("one or more multiplexed sends failed")
    }
    fmt.Printf("sent %d concurrent signals over one connection, ack=%v\n", *multi, *requireAck)
}
