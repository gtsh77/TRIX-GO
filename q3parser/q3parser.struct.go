package q3parser

const (
	MAXFACES  = 6
	ENTMAXVAL = 8
	MAXTEXDUP = 10000
	CMAPDIR = "maps/"
	CMAPEXT = ".cmap"	
)

type CHEAD struct {
	BrushCount  uint32
	EntityCount uint32
	TexelCount  uint32
}

type CBRUSH struct {
	ID            uint32
	FaceCount     uint8
	PlaneCount    uint8
	Faces         [MAXFACES]int
	Vertices      [12 * MAXFACES]float64
	Planes        [9 * MAXFACES]int32
	Texel         [MAXFACES]string
	ShiftX        [MAXFACES]int32
	ShiftY        [MAXFACES]int32
	ScaleX        [MAXFACES]float32
	ScaleY        [MAXFACES]float32
	Width         [MAXFACES]uint32
	Height        [MAXFACES]uint32
	StartX        [MAXFACES]float64
	StartY        [MAXFACES]float64
	StartZ        [MAXFACES]float64
	EndX          [MAXFACES]float64
	EndY          [MAXFACES]float64
	EndZ          [MAXFACES]float64
	DirectionCode [MAXFACES]uint8
}

type CTEX struct {
	Path string
}

type CENTPROP struct {
	Name  string
	Value string
}

type CENT struct {
	ID        uint32
	ClassName string
	ValueCnt  uint8
	Values    [ENTMAXVAL]CENTPROP
}
