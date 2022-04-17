package box

// Box types
const (
	TypeUUID = "uuid"

	TypeFtyp = "ftyp"
	TypeFree = "free"
	TypeSkip = "skip"
	TypeMdat = "mdat"
	TypeMoov = "moov"
	TypeMvhd = "mvhd"
	TypeUdta = "udta"
	TypeCprt = "cprt"
	TypeMeta = "meta"
	TypeHdlr = "hdlr"
	TypeIlst = "ilst"
	TypeTrak = "trak"
	TypeTkhd = "tkhd"
	TypeMdia = "mdia"
	TypeMdhd = "mdhd"
	TypeMinf = "minf"
	TypeStbl = "stbl"
	TypeDinf = "dinf"
	TypeSmhd = "smhd"
	TypeVmhd = "vmhd"
	TypeStsd = "stsd"
	TypeStts = "stts"
	TypeStss = "stss"
	TypeStsc = "stsc"
	TypeStsz = "stsz"
	TypeStco = "stco"
	TypeCtts = "ctts"
)

var boxTypes = map[string]struct{}{
	TypeUUID: {},

	TypeFtyp: {},
	TypeFree: {},
	TypeSkip: {},
	TypeMdat: {},
	TypeMoov: {},
	TypeMvhd: {},
	TypeUdta: {},
	TypeCprt: {},
	TypeMeta: {},
	TypeHdlr: {},
	TypeIlst: {},
	TypeTrak: {},
	TypeTkhd: {},
	TypeMdia: {},
	TypeMdhd: {},
	TypeMinf: {},
	TypeStbl: {},
	TypeDinf: {},
	TypeSmhd: {},
	TypeVmhd: {},
	TypeStsd: {},
	TypeStts: {},
	TypeStss: {},
	TypeStsc: {},
	TypeStsz: {},
	TypeStco: {},
	TypeCtts: {},
}

// BoxTypes returns box types map.
func BoxTypes() map[string]struct{} {
	return boxTypes
}

// IsValidBoxType checks whether input box type is valid or not.
func IsValidBoxType(t string) bool {
	_, ok := boxTypes[t]
	return ok
}
