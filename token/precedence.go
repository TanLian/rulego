package token

// 优先级定义，数字越低优先级越低
const (
	PrecedenceLowest         int = iota
	PrecedenceAssign             // =
	PrecedenceLogicOr            // ||
	PrecedenceLogicAnd           // &&
	PrecedenceCompare            // > or <
	PrecedenceBitwiseOr          // |
	PrecedenceBitwiseXor         // ^
	PrecedenceBitwiseAnd         // &
	PrecedenceBitShift           // << or >>
	PrecedenceAddMinus           // + or -
	PrecedenceMultiplyDivide     // * or /
	PrecedencePrefix             // -X or !X
	PrecedencePower              // **
	PrecedenceCall               // myFunction(X)
	PrecedenceIndex              // array[index]
)

var Precedences = map[TokenType]int{
	ASSIGN:      PrecedenceAssign,
	PLUS:        PrecedenceAddMinus,
	MINUS:       PrecedenceAddMinus,
	TIMES:       PrecedenceMultiplyDivide,
	DIVIDE:      PrecedenceMultiplyDivide,
	MOD:         PrecedenceMultiplyDivide,
	BANG:        PrecedencePrefix,
	GREATER:     PrecedenceCompare,
	LESS:        PrecedenceCompare,
	EQUAL:       PrecedenceCompare,
	NOT_EQUAL:   PrecedenceCompare,
	AND:         PrecedenceBitwiseAnd,
	OR:          PrecedenceBitwiseOr,
	XOR:         PrecedenceBitwiseXor,
	LEFT_SHIFT:  PrecedenceBitShift,
	RIGHT_SHIFT: PrecedenceBitShift,
	LOGIC_AND:   PrecedenceLogicAnd,
	LOGIC_OR:    PrecedenceLogicOr,
	LPAREN:      PrecedenceCall,
	POWER:       PrecedencePower,
	LBRACKET:    PrecedenceIndex,
	DOT:         PrecedenceIndex,
	LBRACE:      PrecedenceIndex,
}

func GetPrecedence(t TokenType) int {
	if v, ok := Precedences[t]; ok {
		return v
	}
	return PrecedenceLowest
}
