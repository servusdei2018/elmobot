package nim

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

const (
	eventDataPrefix = "data: "
	streamDoneEvent = "[DONE]"
)

type StreamDecoder struct {
	scanner *bufio.Scanner
}

func NewStreamDecoder(r io.Reader) *StreamDecoder {
	return &StreamDecoder{
		scanner: bufio.NewScanner(r),
	}
}

func (d *StreamDecoder) Decode() (*StreamEvent, error) {
	for d.scanner.Scan() {
		line := bytes.TrimSpace(d.scanner.Bytes())

		if len(line) == 0 {
			continue
		}

		if !bytes.HasPrefix(line, []byte(eventDataPrefix)) {
			continue
		}

		data := bytes.TrimPrefix(line, []byte(eventDataPrefix))
		data = bytes.TrimSpace(data)

		if bytes.Equal(data, []byte(streamDoneEvent)) {
			return nil, io.EOF
		}

		var event StreamEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, fmt.Errorf("failed to unmarshal stream event: %w", err)
		}

		return &event, nil
	}

	if err := d.scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}

	return nil, io.EOF
}
