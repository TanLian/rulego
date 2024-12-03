package object

type Slice struct {
	Val []any
}

func (s *Slice) object() {}

func (s *Slice) Type() Type {
	return TypeSlice
}

func (s *Slice) GetValue() any {
	return s.Val
}

func (s *Slice) Push(item any) {
	s.Val = append(s.Val, item)
}

func (s *Slice) Len() int {
	return len(s.Val)
}
