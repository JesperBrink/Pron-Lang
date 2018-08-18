package ast

import (
	"Pron-Lang/token"
	"bytes"
	"strings"
)

type Node interface {
	TokenLiteral() string // Metod used for debugging and testing
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// Root of the ast
type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type VarStatement struct {
	Token token.Token // the token.VAR token
	Name  *Identifier
	Value Expression
}

func (ls *VarStatement) statementNode()       {}
func (ls *VarStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *VarStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

type ClassStatement struct {
	Token      token.Token // the token.Class token
	Name       *Identifier
	Fields     []*VarStatement
	Functions  []*DirectFunctionStatement
	InitParams []*InitParam
	InitBody   *BlockStatement
}

func (cs *ClassStatement) statementNode()       {}
func (cs *ClassStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *ClassStatement) String() string {
	var out bytes.Buffer

	fields := []string{}
	for _, field := range cs.Fields {
		fields = append(fields, "var "+field.Name.Value+" = "+field.Value.String())
	}

	params := []string{}
	for _, param := range cs.InitParams {
		params = append(params, param.Parameter.String())
	}

	functions := []string{}
	for _, function := range cs.Functions {
		params = append(functions, function.String())
	}

	out.WriteString("class " + cs.Name.Value + " {")
	out.WriteString(strings.Join(fields, "\n"))
	out.WriteString("init(" + strings.Join(params, ", ") + ") {")
	out.WriteString(cs.InitBody.String())
	out.WriteString("}")
	out.WriteString(strings.Join(functions, "\n"))
	out.WriteString("}")

	return out.String()
}

type InitParam struct {
	Token       token.Token
	Parameter   *Identifier
	IsThisParam bool // true if it is a 'this.paramName'
}

func (ip *InitParam) expressionNode()      {}
func (ip *InitParam) TokenLiteral() string { return ip.Token.Literal }
func (ip *InitParam) String() string       { return ip.Parameter.Value }

type Identifier struct {
	Token         token.Token // the token.IDENT token
	Value         string
	HasThisPrefix bool // is prefixed with 'this.'
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type Null struct{}

func (n *Null) expressionNode()      {}
func (n *Null) TokenLiteral() string { return "Null" }
func (n *Null) String() string       { return "Null" }

type ReturnStatement struct {
	Token       token.Token // the 'return' statement
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type DirectFunctionStatement struct {
	Token    token.Token // the 'func' token
	Name     *Identifier
	Function FunctionLiteral
}

func (dfs *DirectFunctionStatement) statementNode()       {}
func (dfs *DirectFunctionStatement) TokenLiteral() string { return dfs.Token.Literal }
func (dfs *DirectFunctionStatement) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range dfs.Function.Parameters {
		//fmt.Print(p.String())
		params = append(params, p.String())
	}

	exps := []string{}
	for _, e := range dfs.Function.Body.Statements {
		exps = append(exps, e.String())
	}

	out.WriteString(dfs.TokenLiteral() + " ")
	out.WriteString(dfs.Name.String())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString("{")
	out.WriteString(strings.Join(exps, ", "))
	out.WriteString("}")

	return out.String()
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntegerLiteral) expressionNode()      {}
func (i *IntegerLiteral) TokenLiteral() string { return i.Token.Literal }
func (i *IntegerLiteral) String() string       { return i.Token.Literal }

type RealLiteral struct {
	Token token.Token
	Value float64
}

func (r *RealLiteral) expressionNode()      {}
func (r *RealLiteral) TokenLiteral() string { return r.Token.Literal }
func (r *RealLiteral) String() string       { return r.Token.Literal }

type PrefixExpression struct {
	Token    token.Token // The prefix token e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Token //The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie InfixExpression) expressionNode()      {}
func (ie InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

type ElseIfExpression struct {
	Token                          token.Token
	ConditionAndBlockstatementList []*ConditionAndBlockstatementExpression
	Alternative                    *BlockStatement
}

func (ei *ElseIfExpression) expressionNode()      {}
func (ei *ElseIfExpression) TokenLiteral() string { return ei.Token.Literal }
func (ei *ElseIfExpression) String() string {
	var out bytes.Buffer

	elifs := []string{}
	for _, elif := range ei.ConditionAndBlockstatementList[1:] {
		printStr := elif.String()
		elifs = append(elifs, printStr)
	}

	out.WriteString("if (")
	out.WriteString(ei.ConditionAndBlockstatementList[0].Condition.String())
	out.WriteString(") {")
	out.WriteString(ei.ConditionAndBlockstatementList[0].Consequence.String())
	out.WriteString("}")
	out.WriteString(strings.Join(elifs, "\n"))

	if ei.Alternative != nil {
		out.WriteString("else {")
		out.WriteString(ei.Alternative.String())
		out.WriteString("}")
	}

	return out.String()
}

type ConditionAndBlockstatementExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
}

func (cb *ConditionAndBlockstatementExpression) expressionNode()      {}
func (cb *ConditionAndBlockstatementExpression) TokenLiteral() string { return cb.Token.Literal }
func (cb *ConditionAndBlockstatementExpression) String() string {
	var out bytes.Buffer

	out.WriteString("elif (")
	out.WriteString(cb.Condition.String())
	out.WriteString(") {")
	out.WriteString(cb.Consequence.String())
	out.WriteString("}")

	return out.String()
}

type BlockStatement struct {
	Token      token.Token //The  { token
	Statements []Statement
}

func (bs *BlockStatement) expressionNode()      {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token //the 'func' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")

	return out.String()
}

type CallExpression struct {
	Token     token.Token //the '(' token
	Function  Expression  //Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

type ArrayLiteral struct {
	Token    token.Token // the '[' token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type IndexExpression struct {
	Token token.Token //the '[' token
	Left  Expression  //the object being accessed: myArr[2], returnsArray()[1], etc.
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")
	return out.String()
}

type HashLiteral struct {
	Token token.Token //The ´{´ token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode()      {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

type IncrementForloopExpression struct {
	Token    token.Token // The 'for' token
	LocalVar Expression
	From     Expression
	To       Expression
	Body     *BlockStatement
}

func (ic *IncrementForloopExpression) expressionNode()      {}
func (ic *IncrementForloopExpression) TokenLiteral() string { return ic.Token.Literal }
func (ic *IncrementForloopExpression) String() string {
	var out bytes.Buffer

	out.WriteString("for ")
	out.WriteString("( " + ic.LocalVar.String() + " ")
	out.WriteString("from ")
	out.WriteString(ic.From.String() + " ")
	out.WriteString("to ")
	out.WriteString(ic.To.String() + " ) ")
	out.WriteString("{")
	out.WriteString(ic.Body.String())
	out.WriteString("}")

	return out.String()
}

type ArrayForloopExpression struct {
	Token     token.Token // The 'for' token
	LocalVar  Expression
	ArrayName Expression
	Body      *BlockStatement
}

func (af *ArrayForloopExpression) expressionNode()      {}
func (af *ArrayForloopExpression) TokenLiteral() string { return af.Token.Literal }
func (af *ArrayForloopExpression) String() string {
	var out bytes.Buffer

	out.WriteString("for ")
	out.WriteString("( " + af.LocalVar.String() + " ")
	out.WriteString("in ")
	out.WriteString(af.ArrayName.String() + " ) ")
	out.WriteString("{")
	out.WriteString(af.Body.String())
	out.WriteString("}")

	return out.String()
}

type ObjectInitialization struct {
	Token     token.Token // the 'new' token
	Name      *Identifier
	Arguments []Expression
}

func (oi *ObjectInitialization) expressionNode()      {}
func (oi *ObjectInitialization) TokenLiteral() string { return oi.Token.Literal }
func (oi *ObjectInitialization) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, arg := range oi.Arguments {
		args = append(args, arg.String())
	}

	out.WriteString("new " + oi.Name.Value + "(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

type CallObjectFunction struct {
	Token        token.Token // the DOT token
	ObjectName   *Identifier
	FunctionName *Identifier
	Arguments    []Expression
}

func (cof *CallObjectFunction) expressionNode()      {}
func (cof *CallObjectFunction) TokenLiteral() string { return cof.Token.Literal }
func (cof *CallObjectFunction) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, arg := range cof.Arguments {
		args = append(args, arg.String())
	}

	out.WriteString(cof.ObjectName.String())
	out.WriteString(".")
	out.WriteString(cof.FunctionName.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

type Increment struct {
	Token token.Token // the ++
	Name  Identifier
}

func (i *Increment) expressionNode()      {}
func (i *Increment) TokenLiteral() string { return i.Token.Literal }
func (i *Increment) String() string       { return i.Name.Value + "++" }

type Decrement struct {
	Token token.Token // the ++
	Name  Identifier
}

func (d *Decrement) expressionNode()      {}
func (d *Decrement) TokenLiteral() string { return d.Token.Literal }
func (d *Decrement) String() string       { return d.Name.Value + "--" }
