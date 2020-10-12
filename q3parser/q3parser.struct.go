package q3parser

type CHEAD struct {
	BrushCount uint32
	EntityCount uint32
	TexelCount uint32
}

type CTEX struct {
	Path string
}

type CBRUSH struct {
	ID uint32
	FaceCount uint8
	PlaneCount uint8
	Faces uint8
	Vertices int32
	Planes int32
	Texel []string
	ShiftX int32
	ShiftY int32
	ScaleX float32
	ScaleY float32
	Width uint32
	Height uint32
	StartX float64
	StartY float64
	StartZ float64
	EndX float64
	EndY float64
	EndZ float64
	DirectionCode uint8
}

type CENTPROP struct {
	Name string
	Value string
}

type CENT struct {
	ID uint32
	ClassName string
	ValueCnt uint8
	Values []CENTPROP
}