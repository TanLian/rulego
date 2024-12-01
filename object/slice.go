package object

type Slice struct {
	Val []any
}

func (s *Slice) object() {}

func (s *Slice) Type() Type {
	return TypeSlice
}

func (s *Slice) GetValue() any {
	var res []any
	for _, v := range s.Val {
		if v == Null || v == nil {
			continue
		}
		res = append(res, v)
	}
	return res
}

func (s *Slice) Push(item any) {
	s.Val = append(s.Val, item)
}

func (s *Slice) Len() int {
	return len(s.Val)
}
