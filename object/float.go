package object

type Float struct {
	Val float64
}

func (f *Float) object() {}

func (f *Float) Type() Type {
	return TypeFloat
}

func (f *Float) GetValue() any {
	return f.Val
}

func (f *Float) Max(b float64) float64 {
	if f.Val > b {
		return f.Val
	}
	return b
}
