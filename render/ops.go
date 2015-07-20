package render

// DO NOT USE POINTERS IN THESE TYPES
type (
	// op 0x01
	Color struct{ R, G, B, A byte }

	// op 0x02
	Rect struct{ A, B Point }

	// op 0x03
	Tri struct{ A, B, C Point }

	// op 0x04
	Quad struct{ A, B, C, D Point }

	// op 0x05
	Line struct{ A, B Point }

	// op 0x06
	Bezier struct{ A0, A1, B0, B1 Point }
)
