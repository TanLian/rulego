package object

type Struct struct {
	Name    string
	Fields  []string
	Methods map[string]Caller
	Values  map[string]Object
}

type Caller interface {
	Call(args []Object) Object
}

func (rs *Struct) CheckFieldExist(field string) bool {
	for _, v := range rs.Fields {
		if field == v {
			return true
		}
	}
	return false
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

func (rs *Struct) Call(name string, args []Object) Object {
	m, ok := rs.Methods[name]
	if !ok {
		return Null
	}
	return m.Call(args)
}

func (rs *Struct) Clone() *Struct {
	res := &Struct{
		Name:    rs.Name,
		Methods: make(map[string]Caller),
		Values:  make(map[string]Object),
	}
	for _, v := range rs.Fields {
		res.Fields = append(res.Fields, v)
	}
	for k, v := range rs.Methods {
		res.Methods[k] = v
	}
	for k, v := range rs.Values {
		res.Values[k] = v
	}
	return res
}

func (rs *Struct) object() {}

func (rs *Struct) Type() Type {
	return TypeRgStruct
}

func (rs *Struct) GetValue() any {
	return rs
}
