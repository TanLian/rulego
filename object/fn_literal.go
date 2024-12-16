package object

type LiteralFn struct {
	Name  string
	Args  []string
	Block Caller
}

func (fl *LiteralFn) object() {}

func (fl *LiteralFn) Type() Type {
	return TypeFnLiteral
}

func (fl *LiteralFn) GetValue() any {
	return fl
}

func (fl *LiteralFn) Call(args []Object) Object {
	return fl.Block.Call(args)
}
