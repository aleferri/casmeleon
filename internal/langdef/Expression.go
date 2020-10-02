package langdef

import (
	"github.com/aleferri/casmeleon/internal/operators"
	"github.com/aleferri/casmeleon/internal/parsing"
	"github.com/aleferri/casmeleon/internal/text"
	"fmt"
	"strconv"
	"strings"
)

const (
	literal = 0 + iota
	argRef
	enumRef
)

//Operand is used to replace OperandUnion
type Operand struct {
	value       int64
	operandType int //operandType is one of Literal, ArgIndex
}

//Operation pop values from expression stack and push the result
type Operation interface {
	Eval(eStack *[]Operand, values []OpcodeArg) //Eval Operation
	String() string                             //String return the equivalent string for the operation
}

//Eval push the operand in the expression stack
func (p Operand) Eval(eStack *[]Operand, values []OpcodeArg) {
	*eStack = append(*eStack, p)
}

//String return the stringed version of the operand
func (p Operand) String() string {
	return strconv.FormatInt(p.value, 10) + "?" + strconv.FormatInt(int64(p.operandType), 10)
}

//Value of the Operand
func (p Operand) Value(values []OpcodeArg) int64 {
	if p.operandType == literal {
		return p.value
	} else if p.operandType == argRef {
		return int64(values[p.value].value)
	}
	return 0
}

//UnaryOperation is the unary operand name
type UnaryOperation string

//Eval unary operation
func (u UnaryOperation) Eval(eStack *[]Operand, values []OpcodeArg) {
	size := len(*eStack)
	a := (*eStack)[size-1]
	if strings.EqualFold(string(u), ".len") {
		(*eStack)[size-1] = Operand{int64(len(values) - 4), literal} // -3 -> .blist, this_address, this_opcode
		return
	}
	r := int64(operators.Eval(string(u), 0, int(a.Value(values))))
	(*eStack)[size-1] = Operand{r, literal}
}

//String return the name of the operator
func (u UnaryOperation) String() string {
	return string(u)
}

//BinaryOperation is the binary operand name
type BinaryOperation string

//Eval unary operation
func (bin BinaryOperation) Eval(eStack *[]Operand, values []OpcodeArg) {
	size := len(*eStack)
	b := (*eStack)[size-1]
	a := (*eStack)[size-2]
	var r int64
	if strings.EqualFold(string(bin), ".get") {
		index := b.Value(values) + 4
		if index >= int64(len(values)) {
			println("Index out of range,", index, "with a len of", len(values), ", skipping '.get'")
			r = 0
		} else {
			r = int64(values[index].value)
		}
	} else {
		r = int64(operators.Eval(string(bin), int(a.Value(values)), int(b.Value(values))))
	}
	(*eStack)[size-2] = Operand{r, literal}
	*eStack = (*eStack)[:size-1]
}

//String return the name of the operator
func (bin BinaryOperation) String() string {
	return string(bin)
}

//InOperation check if the name is in enums
type InOperation struct {
	enums []Enum
}

//Eval '.in' operation
func (i InOperation) Eval(eStack *[]Operand, values []OpcodeArg) {
	size := len(*eStack)
	b := (*eStack)[size-1]
	a := (*eStack)[size-2]
	if i.enums[b.value].Contains(values[a.value].token.Value()) {
		(*eStack)[size-2] = Operand{1, literal}
		*eStack = (*eStack)[:size-1]
	}
	(*eStack)[size-2] = Operand{0, literal}
	*eStack = (*eStack)[:size-1]
}

//String return a string version of the operator
func (i InOperation) String() string {
	return ".in"
}

//Expression postfix representation
type Expression struct {
	ops []Operation
}

func indexOfString(s string, a []string) int {
	for i, e := range a {
		if strings.EqualFold(s, e) {
			return i
		}
	}
	return -1
}

func indexOfEnumName(name string, enums []Enum) int {
	for i, enum := range enums {
		if strings.EqualFold(enum.name, name) {
			return i
		}
	}
	return -1
}

func indexOfEnumValue(name string, enums []Enum) int {
	for _, enum := range enums {
		if enum.Contains(name) {
			return enum.IndexOf(name)
		}
	}
	return -1
}

func getOperand(t text.Token, paramsName []string, enums []Enum) (Operand, bool) {
	valType := argRef
	k := indexOfString(t.Value(), paramsName)
	if k < 0 {
		valType = enumRef
		k = indexOfEnumName(t.Value(), enums)
		if k < 0 {
			valType = literal
			k = indexOfEnumValue(t.Value(), enums)
		}
	}
	if k >= 0 {
		return Operand{int64(k), valType}, true
	}
	return Operand{}, false
}

func parseLiteral(s string) int64 {
	lenS := len(s)
	if lenS < 3 {
		lit, _ := strconv.ParseInt(s, 0, 32)
		return lit
	}
	var lit int64
	if s[1] == 'b' {
		lit, _ = strconv.ParseInt(s[2:lenS], 2, 32)
	} else if s[1] == 'o' {
		lit, _ = strconv.ParseInt(s[2:lenS], 8, 32)
	} else if s[1] == 'x' {
		lit, _ = strconv.ParseInt(s[2:lenS], 16, 32)
	} else {
		lit, _ = strconv.ParseInt(s, 10, 32)
	}
	return lit
}

//NewExpression check and create the expression in the form of postfix operations
func NewExpression(postfix []text.Token, paramsName []string, enums []Enum) (Expression, *parsing.ParsingError) {
	e := Expression{}
	for i, t := range postfix {
		if t.EnumType() == text.Identifier {
			operand, found := getOperand(t, paramsName, enums)
			if !found {
				err := parsing.NewParsingError(postfix[i], parsing.ErrorUnresolvedSymbol, "unexpected identifier")
				return e, &err
			}
			e.ops = append(e.ops, operand)
		} else if t.EnumType() == text.BinaryOperator {
			isIn := strings.EqualFold(t.Value(), ".in")
			if isIn && indexOfEnumName(postfix[i-1].Value(), enums) < 0 {
				err := parsing.NewParsingError(postfix[i], parsing.ErrorUnresolvedSymbol, "required enum name")
				return e, &err
			}
			if isIn {
				e.ops = append(e.ops, InOperation{enums})
			} else {
				e.ops = append(e.ops, BinaryOperation(t.Value()))
			}
		} else if t.EnumType() == text.UnaryOperator {
			e.ops = append(e.ops, UnaryOperation(t.Value()))
		} else if t.EnumType() == text.Number {
			e.ops = append(e.ops, Operand{parseLiteral(t.Value()), literal})
		}
	}
	return e, nil
}

//Eval expression
func (e *Expression) Eval(params []OpcodeArg) uint {
	eStack := make([]Operand, 0, len(e.ops))
	for _, element := range e.ops {
		element.Eval(&eStack, params)
	}
	return uint(eStack[0].Value(params))
}

//EvalDebug eval expression in debug mode
func (e *Expression) EvalDebug(params []OpcodeArg) uint {
	eStack := make([]Operand, 0, len(e.ops))
	for _, element := range e.ops {
		element.Eval(&eStack, params)
		fmt.Printf("%v, ", eStack)
	}
	println(eStack[0].Value(params))
	return uint(eStack[0].Value(params))
}

//DebugPrint print a debug dump of the Expression
func (e Expression) DebugPrint() {
	for _, elem := range e.ops {
		print(elem.String(), " ; ")
	}
	println()
}
