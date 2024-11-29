package object

import "reflect"

type MapPointer struct {
	Value reflect.Value
}

func (m *MapPointer) object() {}

func (m *MapPointer) Type() Type {
	return TypeMapPointer
}

func (m *MapPointer) GetValue() any {
	res := make(map[any]any)
	for _, v := range m.Value.MapKeys() {
		res[v.Interface()] = m.Value.MapIndex(v).Interface()
	}
	return res
}

func (m *MapPointer) SetField(k, v any) {
	m.Value.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
}

func (m *MapPointer) GetField(k any) any {
	return m.Value.MapIndex(reflect.ValueOf(k)).Interface()
}
