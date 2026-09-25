package mode

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/rs/zerolog/log"
)

// ErrInputExhausted is returned by a mode when the external input source
// that drives its transactions has run dry. The runner treats it as a clean
// stop signal rather than a failed request.
var ErrInputExhausted = errors.New("input exhausted")

// ChunkFeeder splits an io.Reader into fixed-size chunks and hands them out
// to concurrent consumers. A single reader goroutine performs the reads so
// the underlying stream is never accessed concurrently, and a small buffered
// channel lets the producer stay slightly ahead of the consumers.
//
// A blocking Read on a pipe cannot be interrupted, so if the producer stalls
// forever the reader goroutine outlives the context. It exits as soon as the
// read returns, and holds nothing but its own buffer in the meantime.
type ChunkFeeder struct {
	ch chan []byte
	// err is the read error that ended the stream, if it was not a plain
	// EOF. Written by the reader goroutine before it closes ch, and read by
	// consumers only after they observe the close, so the channel close
	// orders the accesses.
	err error
}

// NewChunkFeeder starts reading size-byte chunks from r. Chunks are delivered
// in stream order until EOF, at which point Next returns ErrInputExhausted.
// A trailing partial chunk is discarded with a warning. Any other read error
// ends the stream and is returned from Next so the caller can fail the run
// rather than mistake a broken producer for a clean EOF.
func NewChunkFeeder(ctx context.Context, r io.Reader, size int, buffer int) (*ChunkFeeder, error) {
	if size <= 0 {
		return nil, errors.New("chunk size must be positive")
	}
	if buffer < 0 {
		return nil, errors.New("buffer size must not be negative")
	}
	f := &ChunkFeeder{ch: make(chan []byte, buffer)}
	go f.read(ctx, r, size)
	return f, nil
}

func (f *ChunkFeeder) read(ctx context.Context, r io.Reader, size int) {
	defer close(f.ch)
	var chunks uint64
	for {
		buf := make([]byte, size)
		n, err := io.ReadFull(r, buf)
		switch {
		case err == nil:
		case errors.Is(err, io.EOF):
			log.Info().Uint64("chunks", chunks).Msg("Calldata input reached EOF")
			return
		case errors.Is(err, io.ErrUnexpectedEOF):
			log.Warn().
				Uint64("chunks", chunks).
				Int("strayBytes", n).
				Int("chunkSize", size).
				Msg("Calldata input ended with a partial chunk, discarding it")
			return
		default:
			f.err = err
			return
		}

		select {
		case f.ch <- buf:
			chunks++
		case <-ctx.Done():
			return
		}
	}
}

// Next returns the next chunk. It returns ctx.Err() if the context is done,
// ErrInputExhausted once the stream has ended cleanly and all buffered chunks
// have been consumed, and the wrapped read error if the stream ended because
// the read failed.
func (f *ChunkFeeder) Next(ctx context.Context) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	select {
	case buf, ok := <-f.ch:
		if !ok {
			if f.err != nil {
				return nil, fmt.Errorf("reading calldata input: %w", f.err)
			}
			return nil, ErrInputExhausted
		}
		return buf, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
