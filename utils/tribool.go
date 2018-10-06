package utils

type TriBool uint8

const (
	TriBoolUndefined TriBool = 0
	TriBoolTrue      TriBool = 1
	TriBoolFalse     TriBool = 2
)

func BoolToTri(b bool) TriBool {
	if b {
		return TriBoolTrue
	}
	return TriBoolFalse
}

func BoolPtrToTri(b *bool) TriBool {
	if b == nil {
		return TriBoolUndefined
	}
	return BoolToTri(*b)
}

func (tb TriBool) IsUndefined() bool {
	return tb == TriBoolUndefined
}

func (tb TriBool) IsTrue() bool {
	return tb == TriBoolTrue
}

func (tb TriBool) IsFalse() bool {
	return tb == TriBoolFalse
}

func (tb TriBool) Not() TriBool {
	switch tb {
	case TriBoolUndefined:
		return TriBoolUndefined
	case TriBoolTrue:
		return TriBoolFalse
	case TriBoolFalse:
		return TriBoolTrue
	}
	panic("unreachable")
}

func (tb TriBool) And(another TriBool) TriBool {
	switch tb {
	case TriBoolUndefined:
		if another.IsFalse() {
			return TriBoolFalse
		}
		return TriBoolUndefined
	case TriBoolTrue:
		return another
	case TriBoolFalse:
		return TriBoolFalse
	}
	panic("unreachable")
}

func (tb TriBool) Or(another TriBool) TriBool {
	switch tb {
	case TriBoolUndefined:
		if another.IsTrue() {
			return TriBoolTrue
		}
		return TriBoolUndefined
	case TriBoolTrue:
		return TriBoolTrue
	case TriBoolFalse:
		return another
	}
	panic("unreachable")
}

func (tb TriBool) UnwrapOr(b bool) bool {
	if tb.IsUndefined() {
		return b
	}
	return tb.IsTrue()
}
