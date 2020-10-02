package langdef

import (
	"github.com/aleferri/casmeleon/internal/text"
	"github.com/aleferri/casmeleon/internal/ui"
	"fmt"
	"strconv"
	"strings"
)

//NumberFormat tells how number of the base are prefixed or suffixed
type NumberFormat struct {
	str    string
	prefix bool
	base   int
}

//NewNumberFormat return a new NumberFormat with the specified parameters, also unquote the string if it is quoted
func NewNumberFormat(str string, prefix bool, base int) NumberFormat {
	flag := strings.Trim(str, "'")
	return NumberFormat{flag, prefix, base}
}

//IsNumber check if a string is a number for that base
func (numFormat NumberFormat) IsNumber(n string) bool {
	var number string
	var isDeclaredNumber bool
	if numFormat.prefix {
		number = strings.TrimPrefix(n, numFormat.str)
		isDeclaredNumber = strings.HasPrefix(n, numFormat.str)
	} else {
		number = strings.TrimSuffix(n, numFormat.str)
		isDeclaredNumber = strings.HasSuffix(n, numFormat.str)
	}
	if !isDeclaredNumber {
		return false
	}
	_, err := strconv.ParseInt(number, numFormat.base, 64)
	return err == nil
}

//ParseNumber parse a number in that format
func (numFormat NumberFormat) ParseNumber(n string) int64 {
	number := strings.TrimPrefix(n, numFormat.str)
	if !numFormat.prefix {
		number = strings.TrimSuffix(n, numFormat.str)
	}
	val, _ := strconv.ParseInt(number, numFormat.base, 64)
	return val
}

func (numFormat NumberFormat) String() string {
	position := "postfix"
	if numFormat.prefix {
		position = "prefix"
	}
	return strconv.Itoa(numFormat.base) + " - " + position + " \"" + numFormat.str + "\""
}

//Enum is a set of names
type Enum struct {
	name     string
	nameList []string
}

//String return a stringed version of the enum
func (e Enum) String() string {
	return fmt.Sprintf("%s %v", e.name, e.nameList)
}

//IndexOf return th index of the name in the enums
func (e Enum) IndexOf(name string) int {
	for i, value := range e.nameList {
		if strings.EqualFold(value, name) {
			return i
		}
	}
	return -1
}

//Contains return true if the enum contains the name
func (e Enum) Contains(name string) bool {
	return e.IndexOf(name) > -1
}

//OpcodeArg is the argument sent to the opcode
type OpcodeArg struct {
	name  string
	value int
	token text.Token
}

//Opcode contains name and argument format (if present)
//<opcode declaration> ::= <name> [<args format>]
//<name> ::= <identifier>
//<args format> ::= <arg> | <symbol> | <number>
//<arg> ::= <identifier>
//<number> ::= <base prefix> <base digit> [ {<base digit>} ]
//<identifier> is a valid C identifier
//<symbol> anythings else but comment marker, '->',
type Opcode struct {
	name         string           //opcode name
	paramsFormat []text.TokenType //opcode parameters format
	paramsString []string         //opcode parameters exact string
	paramsName   []string         //opcode parameters name
	nodes        Block            //opcode nodes
}

//NewOpcode return a new opcode
func NewOpcode(name string, paramsFormat []text.TokenType, paramsString []string, nodes Block) Opcode {
	var paramsName []string
	for i, t := range paramsFormat {
		if t == text.Identifier {
			paramsName = append(paramsName, paramsString[i])
		}
	}
	return Opcode{name, paramsFormat, paramsString, paramsName, nodes}
}

func splitByteList(args []text.Token, values []int) []OpcodeArg {
	const maxByte = 4
	size := len(args) - 2
	params := []OpcodeArg{{byteListKeyword, 0, args[0]}}
	params = append(params, OpcodeArg{"this_address", values[size], args[size]})
	params = append(params, OpcodeArg{"this_opcode", values[size+1], args[size+1]})
	params = append(params, OpcodeArg{"", 0, args[0]})
	j := 0
	for i := 0; i < size; i++ {
		arg := args[i]
		if arg.EnumType() == text.Number || arg.EnumType() == text.Identifier {
			number := values[i]
			var k uint
			for k = 0; k < maxByte; k++ {
				shifted := number >> ((maxByte - k - 1) * 8)
				if shifted != 0 || k == maxByte-1 { //least significant byte must be inserted at least
					params = append(params, OpcodeArg{strconv.Itoa(j), shifted, arg})
					j++
				}
			}
		} else if arg.EnumType() == text.DoubleQuotedString {
			str, _ := strconv.Unquote(arg.Value())
			for k := 0; k < len(str); k++ {
				params = append(params, OpcodeArg{strconv.Itoa(j), int(str[k]), arg})
				j++
			}
		}
	}
	return params
}

//Assemble the opcode given the arguments
func (p *Opcode) Assemble(enums []Enum, args []text.Token, values []int, flags uint) ([]uint, ui.SourceCodeError) {
	if len(p.paramsName) == 1 && strings.EqualFold(p.paramsName[0], byteListKeyword) {
		params := splitByteList(args, values)
		return p.nodes.Assemble(params, flags)
	}
	extendedParamsName := append(p.paramsName, "this_address", "this_opcode")
	params := make([]OpcodeArg, 0, len(extendedParamsName))
	for i, a := range args {
		param := OpcodeArg{extendedParamsName[i], values[i], a}
		params = append(params, param)
	}
	return p.nodes.Assemble(params, flags)
}

//OpcodeSlice is a slice of opcodes
type OpcodeSlice []Opcode

//Len of the opcode slice
func (p OpcodeSlice) Len() int {
	return len(p)
}

//Less return true if p[i] come first than p[j]
func (p OpcodeSlice) Less(i, j int) bool {
	compare := strings.Compare(p[i].name, p[j].name)
	if compare < 0 {
		return true
	} else if compare > 0 {
		return false
	}
	min := len(p[i].paramsFormat)
	iMin := true
	if min > len(p[j].paramsFormat) {
		min = len(p[j].paramsFormat)
		iMin = false
	}
	for k := 0; k < min && compare == 0; k++ {
		compare = strings.Compare(p[i].paramsFormat[k].String(), p[j].paramsFormat[k].String())
	}
	return compare < 0 || compare == 0 && iMin
}

//Swap swap two opcode
func (p OpcodeSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
