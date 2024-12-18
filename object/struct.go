package object

type Struct struct {
	Name   string
	Values map[string]Object
}

func (rs *Struct) SetFieldValue(field string, value Object) {
	rs.Values[field] = value
}

func (rs *Struct) GetFieldValue(field string) Object {
	if v, ok := rs.Values[field]; ok {
		return v
	}
	return Null
}

func (rs *Struct) object() {}

func (rs *Struct) Type() Type {
	return TypeRgStruct
}

func (rs *Struct) GetValue() any {
	return rs
}
