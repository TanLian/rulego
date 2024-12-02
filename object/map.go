package object

type Map struct {
	Val map[any]any
}

func (m *Map) object() {}

func (m *Map) Type() Type {
	return TypeMap
}

func (m *Map) GetValue() any {
	//res := make(map[any]any)
	//for k, v := range m.Val {
	//	res[k] = v
	//}
	return m.Val
}

func (m *Map) Len() int {
	return len(m.Val)
}

func (m *Map) Insert(k, v any) {
	m.Val[k] = v
}

func (m *Map) ContainsKey(k any) bool {
	_, ok := m.Val[k]
	return ok
}

func (m *Map) Get(k any) any {
	return m.Val[k]
}

func (m *Map) Remove(k any) {
	delete(m.Val, k)
}
