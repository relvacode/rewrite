package rewrite

import (
	"bufio"
	"bytes"
	"io"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// New creates a new stream.
// Replace all occurrences of `what` with the contents of `with`
func New(r io.Reader, what, with []byte) *Stream {
	return &Stream{
		From:    bufio.NewReaderSize(r, len(with)),
		What:    what,
		With:    with,
		replace: bytes.NewReader(with),
	}
}

// Stream is a content replacing proxy for a byte stream
type Stream struct {
	From *bufio.Reader
	What []byte
	With []byte

	doReplace bool
	replace   *bytes.Reader
}

func (sr *Stream) Read(p []byte) (int, error) {
	var n int
	for {
		if n == len(p) {
			return n, nil
		}

		if sr.doReplace {
			rn, err := sr.replace.Read(p[n:])
			n += rn
			if err == nil {
				continue
			}
			sr.replace.Seek(0, io.SeekStart)
			sr.doReplace = false
		}

		rn, err := sr.read(p[n:])
		n += rn
		if err != nil {
			return n, err
		}
	}
}

func (sr *Stream) read(p []byte) (int, error) {
	var n = len(sr.What)
	// If the next n bytes are exactly the match condition
	peek, err := sr.From.Peek(n)
	if err != nil && err != io.EOF {
		return 0, err
	}

	if bytes.Equal(sr.What, peek) {
		sr.From.Discard(n)
		sr.doReplace = true
		return 0, nil
	}

	if len(peek) > 1 {
		i := bytes.IndexByte(peek[1:], sr.What[0])
		if i > 0 {
			return sr.From.Read(p[:min(len(p), i+1)])
		}

	}

	return sr.From.Read(p[:min(len(p), n)])
}
