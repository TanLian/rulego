package object

type Map struct {
	Val map[any]Object
}

func (m *Map) object() {}

func (m *Map) Type() Type {
	return TypeMap
}

func (m *Map) GetValue() any {
	res := make(map[any]any)
	for k, v := range m.Val {
		res[k] = v.GetValue()
	}
	return res
}
