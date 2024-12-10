package object

type RgStruct struct {
	Name    string
	Fields  []string
	Methods map[string]Method
	Values  map[string]Object
}

type Method interface {
	Call(args []Object) Object
}

func (rs *RgStruct) CheckFieldExist(field string) bool {
	for _, v := range rs.Fields {
		if field == v {
			return true
		}
	}
	return false
}

func (rs *RgStruct) SetFieldValue(field string, value Object) {
	rs.Values[field] = value
}

func (rs *RgStruct) GetFieldValue(field string) Object {
	return rs.Values[field]
}

func (rs *RgStruct) Call(name string, args []Object) Object {
	m, ok := rs.Methods[name]
	if !ok {
		return Null
	}
	return m.Call(args)
}

func (rs *RgStruct) Clone() *RgStruct {
	res := &RgStruct{
		Name:    rs.Name,
		Methods: make(map[string]Method),
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

func (rs *RgStruct) object() {}

func (rs *RgStruct) Type() Type {
	return TypeRgStruct
}

func (rs *RgStruct) GetValue() any {
	return rs
}
