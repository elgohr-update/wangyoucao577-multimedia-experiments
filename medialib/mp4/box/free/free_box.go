// Package free represents free-space boxes which may has type `free` or `skip`.
package free

import (
	"fmt"
	"io"

	"github.com/wangyoucao577/multimedia-experiments/medialib/mp4/box"
	"github.com/wangyoucao577/multimedia-experiments/medialib/util"
)

// Box represents a ftyp box.
type Box struct {
	box.Header

	Data []byte
}

// New creates a new Box.
func New(h box.Header) box.Box {
	return &Box{
		Header: h,
	}
}

// String serializes Box.
func (b Box) String() string {
	return fmt.Sprintf("Header:{%v} Data:%s", b.Header, string(b.Data))
}

// ParsePayload parse payload which requires basic box already exist.
func (b *Box) ParsePayload(r io.Reader) error {
	if b.PayloadSize() > 0 {
		b.Data = make([]byte, b.PayloadSize())
		if err := util.ReadOrError(r, b.Data[:]); err != nil {
			return err
		}
	}
	return nil
}
