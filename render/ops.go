package render

// DO NOT USE POINTERS IN THESE TYPES
type (
	// op 0x01
	Begin struct{}

	// op 0x02
	End struct{}

	// op 0x03
	Color struct{ R, G, B byte }

	// op 0x04
	Point struct{ X, Y float32 }

	// op 0x05
	Rect struct{ A, B Point }

	// op 0x06
	Tri struct{ A, B, C Point }

	// op 0x07
	Quad struct{ A, B, C, D Point }
)
