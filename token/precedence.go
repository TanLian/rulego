package token

// 优先级定义，数字越低优先级越低
const (
	PrecedenceLowest int = iota
	PrecedenceAssign     // =
	//PrecedenceEquals             // ==
	PrecedenceOr             // ||
	PrecedenceAnd            // &&
	PrecedenceCompare        // > or <
	PrecedenceAddMinus       // + or -
	PrecedenceMultiplyDivide // * or /
	PrecedencePrefix         // -X or !X
	PrecedenceCall           // myFunction(X)
	PrecedenceIndex          // array[index]
)

var Precedences = map[TokenType]int{
	ASSIGN:        PrecedenceAssign,
	PLUS:          PrecedenceAddMinus,
	MINUS:         PrecedenceAddMinus,
	PLUS_ASSIGN:   PrecedenceAddMinus,
	TIMES:         PrecedenceMultiplyDivide,
	DIVIDE:        PrecedenceMultiplyDivide,
	MOD:           PrecedenceMultiplyDivide,
	BANG:          PrecedencePrefix,
	GREATER:       PrecedenceCompare,
	GREATER_EQUAL: PrecedenceCompare,
	LESS:          PrecedenceCompare,
	LESS_EQUAL:    PrecedenceCompare,
	EQUAL:         PrecedenceCompare,
	NOT_EQUAL:     PrecedenceCompare,
	AND:           PrecedenceAnd,
	OR:            PrecedenceOr,
	LPAREN:        PrecedenceCall,
	LBRACKET:      PrecedenceIndex,
	DOT:           PrecedenceIndex,
}

func GetPrecedence(t TokenType) int {
	if v, ok := Precedences[t]; ok {
		return v
	}
	return PrecedenceLowest
}
