package mp4

import (
	"encoding/json"
	"io"

	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	"github.com/wangyoucao577/multimedia-experiments/medialib/mp4/box"
	"github.com/wangyoucao577/multimedia-experiments/medialib/mp4/box/free"
	"github.com/wangyoucao577/multimedia-experiments/medialib/mp4/box/ftyp"
	"github.com/wangyoucao577/multimedia-experiments/medialib/mp4/box/mdat"
	"github.com/wangyoucao577/multimedia-experiments/medialib/mp4/box/moof"
	"github.com/wangyoucao577/multimedia-experiments/medialib/mp4/box/moov"
)

// MoofMdat represents composition of one moof and one mdat, since they're stored interleavely like this.
type MoofMdat struct {
	Moof moof.Box `json:"moof,omitempty"`
	Mdat mdat.Box `json:"mdat,omitempty"`
}

// Boxes represents mp4 boxes.
type Boxes struct {
	Ftyp     *ftyp.Box  `json:"ftyp,omitempty"`
	Free     []free.Box `json:"free,omitempty"`
	Moov     *moov.Box  `json:"moov,omitempty"`
	MoofMdat []MoofMdat `json:"moof_mdat,omitempty"` // make sure moof,mdat can be pared and stored interleavely

	//TODO: other boxes

	// internal vars for parsing or other handling
	boxesCreator map[string]box.NewFunc `json:"-"`
}

func newBoxes() Boxes {
	return Boxes{
		boxesCreator: map[string]box.NewFunc{
			box.TypeFtyp: ftyp.New,
			box.TypeFree: free.New,
			box.TypeSkip: free.New,
			box.TypeMdat: mdat.New,
			box.TypeMoov: moov.New,
			box.TypeMoof: moof.New,
		},
	}
}

// JSON marshals boxes to JSON representation
func (b Boxes) JSON() ([]byte, error) {
	return json.Marshal(b)
}

// JSONIndent marshals boxes to JSON representation with customized indent.
func (b Boxes) JSONIndent(prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(b, prefix, indent)
}

// YAML formats boxes to YAML representation.
func (b Boxes) YAML() ([]byte, error) {
	j, err := json.Marshal(b)
	if err != nil {
		return j, err
	}
	return yaml.JSONToYAML(j)
}

// CreateSubBox creates directly included box, such as create `mvhd` in `moov`, or create `moov` on top level.
//   return ErrNotImplemented is the box doesn't have any sub box.
func (b *Boxes) CreateSubBox(h box.Header) (box.Box, error) {
	creator, ok := b.boxesCreator[h.Type.String()]
	if !ok {
		glog.V(2).Infof("unknown box type %s, size %d payload %d", h.Type.String(), h.Size, h.PayloadSize())
		return nil, box.ErrUnknownBoxType
	}

	createdBox := creator(h)
	if createdBox == nil {
		glog.Fatalf("create box type %s failed", h.Type.String())
	}

	switch h.Type.String() {
	case box.TypeFtyp:
		b.Ftyp = createdBox.(*ftyp.Box)
	case box.TypeFree, box.TypeSkip:
		b.Free = append(b.Free, *createdBox.(*free.Box))
		createdBox = &b.Free[len(b.Free)-1] // reference to the last empty free box
	case box.TypeMdat:
		if len(b.MoofMdat) > 0 {
			if err := b.MoofMdat[len(b.MoofMdat)-1].Mdat.Validate(); err == nil { // expect error
				glog.Warningf("expect empty mdat but got a valid one %v", b.MoofMdat[len(b.MoofMdat)-1].Mdat)
				b.MoofMdat = append(b.MoofMdat, MoofMdat{}) // append new one to avoid lost mdat
			}
		}
		b.MoofMdat[len(b.MoofMdat)-1].Mdat = *createdBox.(*mdat.Box)
		createdBox = &b.MoofMdat[len(b.MoofMdat)-1].Mdat
	case box.TypeMoov:
		b.Moov = createdBox.(*moov.Box)
	case box.TypeMoof:
		// Moof is required present before Mdat, so always create a new one if moof encountered.
		b.MoofMdat = append(b.MoofMdat, MoofMdat{Moof: *createdBox.(*moof.Box)})
		createdBox = &b.MoofMdat[len(b.MoofMdat)-1].Moof // reference to the last empty moof box
	}

	return createdBox, nil
}

// ParsePayload acts as an root box to parse all sub boxes.
func (b *Boxes) ParsePayload(r io.Reader) error {

	for {
		if _, err := box.ParseBox(r, b); err != nil {
			if err == io.EOF {
				break
			} else if err == box.ErrUnknownBoxType {
				continue
			}
			return err
		}
	}

	return nil
}
