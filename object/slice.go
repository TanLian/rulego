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

func (s *Slice) Pop() any {
	if len(s.Val) == 0 {
		return nil
	}
	res := s.Val[len(s.Val)-1]
	s.Val = s.Val[:len(s.Val)-1]
	return res
}

func (s *Slice) Len() int {
	return len(s.Val)
}

func (s *Slice) Clone() []any {
	res := make([]any, len(s.Val))
	for i, v := range s.Val {
		res[i] = v
	}
	return res
}
