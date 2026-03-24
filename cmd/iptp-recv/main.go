package main

import (
    "flag"
    "fmt"
    "log"

    "github.com/spicecoder/iptp-pipe/iptp"
)

func main() {
    port := flag.Int("port", 9000, "listen port")
    flag.Parse()

    addr := fmt.Sprintf(":%d", *port)
    log.Printf("iptp-recv listening on %s", addr)

    err := iptp.Listen(addr, func(env iptp.Envelope) error {
        log.Printf("received id=%s intention=%s pulses=%d", env.ID, env.Signal.Intention, len(env.Signal.Pulses))
        return nil
    })
    if err != nil {
        log.Fatal(err)
    }
}
