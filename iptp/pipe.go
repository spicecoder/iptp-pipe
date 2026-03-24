package iptp

import (
    "bufio"
    "bytes"
    "encoding/binary"
    "encoding/json"
    "errors"
    "hash/crc32"
    "io"
)

var magic = []byte{'I', 'P', 'T', 'P'}
const version byte = 1

func checksum(b []byte) uint32 {
    return crc32.ChecksumIEEE(b)
}

func MarshalWithChecksum(env *Envelope) ([]byte, error) {
    env.Checksum = 0
    raw, err := json.Marshal(env)
    if err != nil {
        return nil, err
    }
    env.Checksum = checksum(raw)
    return json.Marshal(env)
}

func VerifyEnvelope(env *Envelope) error {
    original := env.Checksum
    env.Checksum = 0
    raw, err := json.Marshal(env)
    if err != nil {
        return err
    }
    if checksum(raw) != original {
        return errors.New("checksum_mismatch")
    }
    env.Checksum = original
    return nil
}

// Frame format:
// MAGIC(4) | VERSION(1) | FLAGS(1) | PAYLOAD_LENGTH(4) | PAYLOAD
func WriteFrame(w io.Writer, payload []byte) error {
    header := bytes.NewBuffer(nil)
    header.Write(magic)
    header.WriteByte(version)
    header.WriteByte(0) // flags reserved
    if err := binary.Write(header, binary.BigEndian, uint32(len(payload))); err != nil {
        return err
    }
    if _, err := w.Write(header.Bytes()); err != nil {
        return err
    }
    _, err := w.Write(payload)
    return err
}

func ReadFrame(r *bufio.Reader) ([]byte, error) {
    header := make([]byte, 10)
    if _, err := io.ReadFull(r, header); err != nil {
        return nil, err
    }
    if !bytes.Equal(header[:4], magic) {
        return nil, errors.New("bad_magic")
    }
    if header[4] != version {
        return nil, errors.New("bad_version")
    }
    n := binary.BigEndian.Uint32(header[6:10])
    payload := make([]byte, n)
    if _, err := io.ReadFull(r, payload); err != nil {
        return nil, err
    }
    return payload, nil
}

func PulseValue(sig Signal, name string) (string, bool) {
    for _, p := range sig.Pulses {
        if p.Name == name {
            return p.Value, true
        }
    }
    return "", false
}
