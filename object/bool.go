package object

type Bool struct {
	Val bool
}

func (b *Bool) object() {}

func (b *Bool) Type() Type {
	return TypeBool
}

func (b *Bool) GetValue() any {
	return b.Val
}
