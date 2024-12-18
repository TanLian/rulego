package ast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type Closure struct {
	*object.Empty
	Name      string
	Args      []string
	States    []Statement
	OuterVars map[string]object.Object
}

func (c *Closure) Eval(env *environment.Environment) object.Object {
	// 初始化外部变量映射
	if c.OuterVars == nil {
		c.OuterVars = make(map[string]object.Object)
	}

	// 遍历所有语句，收集外部变量引用
	for _, stmt := range c.States {
		c.collectOuterVars(stmt, env)
	}

	// 将闭包注册到当前环境
	env.SetCurrent(c.Name, c)
	return c
}

// collectOuterVars 收集语句中引用的外部变量
func (c *Closure) collectOuterVars(stmt Statement, env *environment.Environment) {
	switch node := stmt.(type) {
	case *ExpressionStatement:
		c.collectOuterVarsFromExpr(node.Expr, env)
	case *Block:
		for _, s := range node.States {
			c.collectOuterVars(s, env)
		}
	}
}

// collectOuterVarsFromExpr 从表达式中收集外部变量引用
func (c *Closure) collectOuterVarsFromExpr(expr Expression, env *environment.Environment) {
	switch node := expr.(type) {
	case *Ident:
		// 检查是否是外部变量引用
		name := node.Token.Value
		if !c.isLocalVar(name) {
			// 如果不是局部变量，尝试从环境中获取
			if obj, ok := env.Get(name); ok {
				c.OuterVars[name] = obj
			}
		}
	case *Assign:
		// 处理赋值表达式的右侧
		c.collectOuterVarsFromExpr(node.Right, env)
		// 如果左侧是标识符，也需要检查
		if ident, ok := node.Left.(*Ident); ok {
			name := ident.Token.Value
			if !c.isLocalVar(name) {
				if obj, ok := env.Get(name); ok {
					c.OuterVars[name] = obj
				}
			}
		}
	case *Bang:
		c.collectOuterVarsFromExpr(node.Expr, env)
	case *BitwiseAnd:
		c.collectOuterVarsFromExpr(node.Left, env)
		c.collectOuterVarsFromExpr(node.Right, env)
	case *BitwiseNot:
		c.collectOuterVarsFromExpr(node.Expr, env)
	case *BitwiseXOR:
		c.collectOuterVarsFromExpr(node.Left, env)
		c.collectOuterVarsFromExpr(node.Right, env)
	case *Compare:
		c.collectOuterVarsFromExpr(node.Left, env)
		c.collectOuterVarsFromExpr(node.Right, env)
	case *Divide:
		c.collectOuterVarsFromExpr(node.Left, env)
		c.collectOuterVarsFromExpr(node.Right, env)
	case *Dot:
		c.collectOuterVarsFromExpr(node.Left, env)
		c.collectOuterVarsFromExpr(node.Right, env)
	case *For:
		c.collectOuterVars(node.Initial, env)
		c.collectOuterVarsFromExpr(node.Condition, env)
		c.collectOuterVars(node.Post, env)
		for _, sta := range node.Statements {
			c.collectOuterVars(sta, env)
		}
	case *Group:
		c.collectOuterVarsFromExpr(node.Expr, env)
	case *If:
		for _, v := range node.Ifs {
			c.collectOuterVarsFromExpr(v.Expr, env)
			for _, sta := range v.States {
				c.collectOuterVars(sta, env)
			}
		}
		for _, sta := range node.Else {
			c.collectOuterVars(sta, env)
		}
	case *Index:
		c.collectOuterVarsFromExpr(node.Data, env)
		c.collectOuterVarsFromExpr(node.Key, env)
		c.collectOuterVarsFromExpr(node.End, env)
	case *LeftShift:
		c.collectOuterVarsFromExpr(node.Left, env)
		c.collectOuterVarsFromExpr(node.Right, env)
	case *LogicAnd:
		c.collectOuterVarsFromExpr(node.Left, env)
		c.collectOuterVarsFromExpr(node.Right, env)
	case *LogicOr:
		c.collectOuterVarsFromExpr(node.Left, env)
		c.collectOuterVarsFromExpr(node.Right, env)
	case *Map:
		for k, v := range node.KV {
			c.collectOuterVarsFromExpr(k, env)
			c.collectOuterVarsFromExpr(v, env)
		}
	case *Minus:
		c.collectOuterVarsFromExpr(node.Left, env)
		c.collectOuterVarsFromExpr(node.Right, env)
	case *Mod:
		c.collectOuterVarsFromExpr(node.Left, env)
		c.collectOuterVarsFromExpr(node.Right, env)
	case *Negative:
		c.collectOuterVarsFromExpr(node.Expr, env)
	case *Or:
		c.collectOuterVarsFromExpr(node.Left, env)
		c.collectOuterVarsFromExpr(node.Right, env)
	case *Positive:
		c.collectOuterVarsFromExpr(node.Expr, env)
	case *Power:
		c.collectOuterVarsFromExpr(node.Left, env)
		c.collectOuterVarsFromExpr(node.Right, env)
	case *RightShift:
		c.collectOuterVarsFromExpr(node.Left, env)
		c.collectOuterVarsFromExpr(node.Right, env)
	case *Slice:
		for _, v := range node.Data {
			c.collectOuterVarsFromExpr(v, env)
		}
		c.collectOuterVarsFromExpr(node.InitExpr, env)
		c.collectOuterVarsFromExpr(node.LenExpr, env)
	case *Switch:
		c.collectOuterVarsFromExpr(node.Expr, env)
		for _, v := range node.Cases {
			c.collectOuterVarsFromExpr(v.Expr, env)
			for _, sta := range v.States {
				c.collectOuterVars(sta, env)
			}
		}
		for _, sta := range node.Default {
			c.collectOuterVars(sta, env)
		}
	case *Call: // 处理函数调用
		c.collectOuterVarsFromExpr(node.Left, env)

		args, ok := node.Arguments.(*Slice)
		if !ok {
			panic("TypeError: expect slice expression")
		}
		for _, v := range args.Data {
			c.collectOuterVarsFromExpr(v, env)
		}
	case *Plus:
		// 处理加法表达式的两侧
		c.collectOuterVarsFromExpr(node.Left, env)
		c.collectOuterVarsFromExpr(node.Right, env)
		// 可以根据需要添加其他表达式类型的处理...
	case *Return:
		c.collectOuterVarsFromExpr(node.Expr, env)
	case *Times:
		c.collectOuterVarsFromExpr(node.Left, env)
		c.collectOuterVarsFromExpr(node.Right, env)
	default: // TODO
	}
}

// isLocalVar 检查变量是否是局部变量
func (c *Closure) isLocalVar(name string) bool {
	for _, v := range c.Args {
		if v == name {
			return true
		}
	}
	return false
}

func (c *Closure) Call(args []object.Object) object.Object {
	env := environment.New(environment.Root)
	for i := 0; i < len(c.Args); i++ {
		env.SetCurrent(c.Args[i], args[i])
	}
	for k, v := range c.OuterVars {
		env.SetCurrent(k, v)
	}
	for i, v := range c.States {
		if i == len(c.States)-1 {
			obj, _ := v.Exec(env)
			return obj
		}

		if obj, flg := v.Exec(env); flg&RETURN != 0 {
			return obj
		}
	}
	return object.Null
}

func (c *Closure) Type() object.Type {
	return object.TypeClosure
}

func (c *Closure) String() string {
	if c == nil {
		return ""
	}
	var s strings.Builder
	s.WriteString("fn ")
	s.WriteString(c.Name)
	s.WriteString(fmt.Sprintf("(%s) {", strings.Join(c.Args, ",")))
	for _, v := range c.States {
		s.WriteString(v.String() + ";")
	}
	s.WriteString("}")
	return s.String()
}

func (c *Closure) AST(num int) string {
	if c == nil {
		return ""
	}
	var s strings.Builder
	s.WriteString("*ast.Closure {\n")
	s.WriteString(strings.Repeat(". ", num+1) + " Name: " + c.Name + "\n")
	s.WriteString(strings.Repeat(". ", num+1) + " Args: (" + strings.Join(c.Args, ",") + ")\n")
	s.WriteString(strings.Repeat(". ", num+1) + " Statements: {\n")
	for i, v := range c.States {
		s.WriteString(strings.Repeat(". ", num+2) + strconv.Itoa(i) + ": " + v.AST(num+2))
	}
	s.WriteString(strings.Repeat(". ", num+1) + " }\n")
	s.WriteString(strings.Repeat(". ", num) + " }\n")
	return s.String()
}

func (c *Closure) expressionNode() {}
