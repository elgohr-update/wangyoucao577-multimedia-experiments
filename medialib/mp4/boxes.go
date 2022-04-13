package mp4

import (
	"github.com/wangyoucao577/multimedia-experiments/medialib/mp4/box/free"
	"github.com/wangyoucao577/multimedia-experiments/medialib/mp4/box/ftyp"
)

// Boxes represents mp4 boxes.
type Boxes struct {
	Ftyp *ftyp.Box
	Free []free.Box

	//TODO: other boxes
}
