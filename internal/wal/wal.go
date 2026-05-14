package wal

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

const (
	maxSize int = 64 << 20
)

var table *crc32.Table = crc32.MakeTable(crc32.IEEE)

type entry struct {
	key   []byte
	value []byte
}

type wal struct {
	mu     sync.Mutex
	fd     *os.File
	logger log.Logger
	path   string
}

func Create(dirPath string, seqNum int) (*wal, error) {
	w := &wal{
		mu:   sync.Mutex{},
		path: filepath.Join(dirPath, "wal", fmt.Sprintf("wal-%06d.log", seqNum)),
	}
	fd, err := os.Create(w.path)
	if err != nil {
		return nil, err
	}
	w.fd = fd
	return w, nil

}

// write []byte to w.path, it returns the total number of bytes written and error
func (w *wal) Write(entries ...entry) (int, error) {

	w.mu.Lock()
	defer w.mu.Unlock()

	totalBytes := 0
	f := w.fd

	for _, e := range entries {

		buf := encodePayload(e)
		totalBytes += len(buf)

		n, err := f.Write(buf)
		if err != nil || n != len(buf) {
			return 0, err
		}

		err = f.Sync()
		if err != nil {
			return 0, err
		}

	}

	return totalBytes, nil

}

func (w *wal) Read() ([]entry, error) {

	w.mu.Lock()
	defer w.mu.Unlock()

	f := w.fd
	if _, err := f.Seek(0, 0); err != nil {
		return nil, err
	}

	var entries []entry

	for {

		header := make([]byte, 8)
		_, err := io.ReadFull(f, header)

		if err == io.EOF {
			break
		}

		if err == io.ErrUnexpectedEOF {
			break
		}

		if err != nil {
			return entries, err
		}

		length := binary.LittleEndian.Uint32(header[0:4])

		buf := make([]byte, 8+length)
		copy(buf[:8], header)

		if _, err := io.ReadFull(f, buf[8:]); err != nil {
			break
		}

		_, entryTmp := decodePayload(buf)

		entries = append(entries, entryTmp)

	}

	return entries, nil

}

// shaped as [length checksum payload]
func encodePayload(e entry) []byte {

	payload := make([]byte, 8+len(e.key)+len(e.value))
	binary.LittleEndian.PutUint32(payload[0:4], uint32(len(e.key)))
	copy(payload[4:4+len(e.key)], e.key)
	binary.LittleEndian.PutUint32(payload[4+len(e.key):8+len(e.key)], uint32(len(e.value)))
	copy(payload[8+len(e.key):], e.value)

	buf := make([]byte, 8+len(payload))

	checkSum := crc32.Checksum(payload, table)

	binary.LittleEndian.PutUint32(buf[0:4], uint32(len(payload)))
	binary.LittleEndian.PutUint32(buf[4:8], checkSum)
	copy(buf[8:], payload)

	return buf
}

// shaped as [length checksum payload]
func decodePayload(buf []byte) (int, entry) {

	length := int(binary.LittleEndian.Uint32(buf[0:4]))
	checkSum := binary.LittleEndian.Uint32(buf[4:8])
	payload := buf[8:]

	if crc32.Checksum(payload, table) != checkSum {
		return 0, entry{}
	}

	keyLength := binary.LittleEndian.Uint32(payload[0:4])
	valueLength := binary.LittleEndian.Uint32(payload[4+keyLength : 8+keyLength])
	keys := make([]byte, keyLength)
	copy(keys, payload[4:4+keyLength])
	values := make([]byte, valueLength)
	copy(values, payload[8+keyLength:8+keyLength+valueLength])

	e := entry{
		key:   keys,
		value: values,
	}

	return length, e
}

// swap a new log
func rotate() {

}

// api for outer to swap a new log
func (w *wal) Rotate() {

}
