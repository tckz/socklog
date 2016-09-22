package socklog

import (
	"bufio"
	"io"
)

type MaskingWriter struct {
	buffer *bufio.Writer
}

func (m *MaskingWriter) Write(p []byte) (n int, err error) {
	var written int = 0

	for _, e := range p {
		var c byte = e
		if e <= 0x1f && e != 0x0a && e != 0x0d && e != 0x09 && e != 0x0c {
			c = '.'
		}
		err := m.buffer.WriteByte(c)
		if err != nil {
			return written, err
		}
		written = written + 1
	}

	err = m.buffer.Flush()
	if err != nil {
		return written, err
	}

	return written, nil
}

func NewMaskingWriter(dest io.Writer) *MaskingWriter {
	ret := &MaskingWriter{
		buffer: bufio.NewWriter(dest),
	}
	return ret
}
