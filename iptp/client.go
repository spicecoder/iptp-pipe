package iptp

import (
    "bufio"
    "encoding/json"
    "errors"
    "fmt"
    "net"
    "sync"
    "time"
)

type Client struct {
    conn    net.Conn
    reader  *bufio.Reader
    writeMu sync.Mutex

    pendingMu sync.Mutex
    pending   map[string]chan Envelope

    closeOnce sync.Once
    closed    chan struct{}
}

func Dial(addr string) (*Client, error) {
    conn, err := net.Dial("tcp", addr)
    if err != nil {
        return nil, err
    }
    c := &Client{
        conn:    conn,
        reader:  bufio.NewReader(conn),
        pending: make(map[string]chan Envelope),
        closed:  make(chan struct{}),
    }
    go c.readLoop()
    return c, nil
}

func (c *Client) Close() error {
    var err error
    c.closeOnce.Do(func() {
        close(c.closed)
        err = c.conn.Close()
        c.failAllPending(errors.New("client_closed"))
    })
    return err
}

func (c *Client) Send(env Envelope, opts SendOptions) error {
    if opts.AckTimeout <= 0 {
        opts.AckTimeout = 2 * time.Second
    }
    if opts.MaxRetries < 0 {
        opts.MaxRetries = 0
    }
    env.SentAt = time.Now().UTC()

    for attempt := 0; attempt <= opts.MaxRetries; attempt++ {
        var ackCh chan Envelope
        if opts.RequireProtocolAck {
            ackCh = make(chan Envelope, 1)
            c.registerPending(env.ID, ackCh)
        }

        payload, err := MarshalWithChecksum(&env)
        if err != nil {
            if opts.RequireProtocolAck {
                c.unregisterPending(env.ID)
            }
            return err
        }

        c.writeMu.Lock()
        err = WriteFrame(c.conn, payload)
        c.writeMu.Unlock()
        if err != nil {
            if opts.RequireProtocolAck {
                c.unregisterPending(env.ID)
            }
            return err
        }

        if !opts.RequireProtocolAck {
            return nil
        }

        select {
        case resp, ok := <-ackCh:
            c.unregisterPending(env.ID)
            if !ok {
                continue
            }
            switch resp.Signal.Intention {
            case "iptp:ack":
                return nil
            case "iptp:nack":
                if attempt == opts.MaxRetries {
                    reason, _ := PulseValue(resp.Signal, "reason")
                    if reason == "" {
                        reason = "protocol_nack"
                    }
                    return fmt.Errorf("send rejected: %s", reason)
                }
                continue
            default:
                continue
            }
        case <-time.After(opts.AckTimeout):
            c.unregisterPending(env.ID)
            continue
        case <-c.closed:
            c.unregisterPending(env.ID)
            return errors.New("client_closed")
        }
    }

    return errors.New("no_protocol_ack")
}

func (c *Client) readLoop() {
    for {
        payload, err := ReadFrame(c.reader)
        if err != nil {
            c.failAllPending(err)
            _ = c.Close()
            return
        }

        var env Envelope
        if err := json.Unmarshal(payload, &env); err != nil {
            continue
        }
        if err := VerifyEnvelope(&env); err != nil {
            continue
        }

        if env.Signal.Intention != "iptp:ack" && env.Signal.Intention != "iptp:nack" {
            continue
        }
        ackFor, ok := PulseValue(env.Signal, "ack_for")
        if !ok || ackFor == "" {
            continue
        }

        c.pendingMu.Lock()
        ch := c.pending[ackFor]
        c.pendingMu.Unlock()
        if ch != nil {
            select {
            case ch <- env:
            default:
            }
        }
    }
}

func (c *Client) registerPending(id string, ch chan Envelope) {
    c.pendingMu.Lock()
    c.pending[id] = ch
    c.pendingMu.Unlock()
}

func (c *Client) unregisterPending(id string) {
    c.pendingMu.Lock()
    ch, ok := c.pending[id]
    if ok {
        delete(c.pending, id)
        close(ch)
    }
    c.pendingMu.Unlock()
}

func (c *Client) failAllPending(_ error) {
    c.pendingMu.Lock()
    for id, ch := range c.pending {
        delete(c.pending, id)
        close(ch)
    }
    c.pendingMu.Unlock()
}
