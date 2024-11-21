package object

import "strings"

type String struct {
	Val []rune
}

func (s *String) object() {}

func (s *String) Type() Type {
	return TypeString
}

func (s *String) GetValue() any {
	return string(s.Val)
}

func (s *String) Trim(cutSet string) string {
	return strings.Trim(string(s.Val), cutSet)
}

func (s *String) TrimSpace() string {
	return strings.TrimSpace(string(s.Val))
}

func (s *String) Split(sep string) []string {
	return strings.Split(string(s.Val), sep)
}

func (s *String) Reverse() string {
	n := len(s.Val)
	res := make([]rune, n)
	for i, v := range s.Val {
		res[n-i-1] = v
	}
	return string(res)
}

func (s *String) Contains(substr string) bool {
	return strings.Contains(string(s.Val), substr)
}
