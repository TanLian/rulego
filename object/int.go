package object

type Int struct {
	Val int64
}

func (i *Int) object() {}

func (i *Int) Type() Type {
	return TypeInt
}

func (i *Int) GetValue() any {
	return i.Val
}

func (i *Int) Max(b int64) int64 {
	if i.Val > b {
		return i.Val
	}
	return b
}
