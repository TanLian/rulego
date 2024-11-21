package object

type Rune struct {
	Val rune
}

func (r *Rune) object() {}

func (r *Rune) Type() Type {
	return TypeRune
}

func (r *Rune) GetValue() any {
	return r.Val
}
