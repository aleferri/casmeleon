package langdef

import (
	"bitbucket.org/mrpink95/casmeleon/internal/lexing"
	"bitbucket.org/mrpink95/casmeleon/internal/parsing"
	"bitbucket.org/mrpink95/casmeleon/internal/text"
	"bitbucket.org/mrpink95/casmeleon/internal/ui"
	"io"
	"sort"
	"strconv"
	"strings"
)

const byteListKeyword = "...bytes"

//LangDef is the language definition
//it contains the list of defined enums,
//the chosen number format and a list
//of the allowed opcodes
type LangDef struct {
	enums         []Enum
	numberFormats []NumberFormat
	opcodes       []Opcode
}

//LangParsingDefaultOptions return a valid matching options for lexing of a language definition
func LangParsingDefaultOptions() lexing.TokenMatchingOptions {
	const separators = "@$#,;:{}()[]'\"\t "
	const operators = "+ - * / % ^ & | ~ > < ! <= >= == && || != << >> -> "
	const lineComment = "//"
	return lexing.NewMatchingOptions(operators, separators, operators+lineComment, lineComment)
}

func identifyOperators(t text.Token) text.Token {
	if t.EnumType() != text.Identifier {
		return t
	}
	if strings.EqualFold(t.Value(), ".in") {
		return t.WithType(text.BinaryOperator)
	} else if strings.EqualFold(t.Value(), ".get") {
		return t.WithType(text.BinaryOperator)
	} else if strings.EqualFold(t.Value(), ".len") {
		return t.WithType(text.UnaryOperator)
	}
	return t
}

//NewLangDef create a new language definition from a reader
func NewLangDef(name string, reader io.Reader, ui ui.UI) *LangDef {
	src := parsing.NewTokenBuffer(name, reader, LangParsingDefaultOptions())
	report := func(err *parsing.ParsingError) {
		line := src.Lines()[err.GetLine()-1]
		err.Report(ui, line)
	}
	state := NewParserState(ui, src, true, report)
	builder := NewLangBuilder()
	return (&builder).Build(state)
}

//GetOpcodeNumber return the number of the opcode in the assembly language
func (lang *LangDef) GetOpcodeNumber() int {
	return len(lang.opcodes)
}

//GetEnumValue return the value of the specified identifier
func (lang *LangDef) GetEnumValue(identifier string) int {
	for _, enum := range lang.enums {
		index := enum.IndexOf(identifier)
		if index > -1 {
			return index
		}
	}
	return -1
}

//IsEnumValue return true if the identifier is an enum value
func (lang *LangDef) IsEnumValue(identifier string) bool {
	return lang.GetEnumValue(identifier) > -1
}

//IdentifyNumber identify if the token is a number or not
func (lang *LangDef) IdentifyNumber(t text.Token) (text.Token, ui.SourceCodeError) {
	if t.EnumType() != text.Identifier && t.EnumType() != text.Number {
		return t, nil
	}
	for _, n := range lang.numberFormats {
		if n.IsNumber(t.Value()) {
			number := int(n.ParseNumber(t.Value()))
			return text.NewSpecialToken(t, strconv.Itoa(number), text.Number), nil
		}
	}
	if t.EnumType() == text.Identifier {
		return t, nil
	}
	_, decimal := strconv.Atoi(t.Value())
	if decimal != nil {
		return t, parsing.NewParsingError(t, parsing.InvalidNumberFormat, "")
	}
	return t, nil
}

//RemoveSpecialPurposeChars remove char used by numbers format as prefix or suffix
func (lang *LangDef) RemoveSpecialPurposeChars(s string) string {
	for _, num := range lang.numberFormats {
		s = strings.Trim(s, num.str)
	}
	return s
}

//SortOpcodes sort opcode slice for successive searches
func (lang *LangDef) SortOpcodes() {
	var opcodeSlice OpcodeSlice
	opcodeSlice = lang.opcodes
	sort.Sort(opcodeSlice)
}

//StringOpcode return the equivalent string for opcode name and format
func (lang *LangDef) StringOpcode(index int, argsName []text.Token) string {
	if index < 0 || index >= len(lang.opcodes) {
		return ""
	}
	opcode := lang.opcodes[index]
	argIndex := 0
	result := opcode.name + " "
	isByteList := len(opcode.paramsFormat) > 0 && strings.EqualFold(opcode.paramsName[0], byteListKeyword)
	for i, param := range opcode.paramsFormat {
		if param == text.Identifier {
			result += " " + argsName[argIndex].Value()
			argIndex++
		} else {
			result += " " + opcode.paramsString[i]
		}
	}
	if isByteList {
		for argIndex < len(argsName)-2 {
			result += ", " + argsName[argIndex].Value()
			argIndex++
		}
	}
	return result
}

//AssembleOpcode assemble the opcode at index n
func (lang *LangDef) AssembleOpcode(n int, tokens []text.Token, values []int, flags uint) ([]uint, ui.SourceCodeError) {
	if n < 0 || n >= len(lang.opcodes) {
		return []uint{}, nil
	}
	return lang.opcodes[n].Assemble(lang.enums, tokens, values, flags)
}

//FilterWindow filter is a token by token filter
type FilterWindow interface {
	FilterByName(name string) FilterWindow
	FilterByToken(token text.Token) (FilterWindow, bool)
	Collect() (bool, int)
	IsEmpty() bool
	String() string
}

//DefaultFilterWindow is a the common implementation of FilterWindow
type DefaultFilterWindow struct {
	opcodes    []Opcode
	matchIndex int
	firstIndex int
}

//NewFilterWindow return a new filterable data set
func (lang *LangDef) NewFilterWindow() FilterWindow {
	return DefaultFilterWindow{lang.opcodes, 0, 0}
}

//FilterByName filter opcode by name
func (f DefaultFilterWindow) FilterByName(name string) FilterWindow {
	match := false
	var k int
	for k = 0; k < len(f.opcodes) && !match; k++ {
		opcode := f.opcodes[k]
		match = strings.EqualFold(opcode.name, name)
	}
	start := k
	if !match {
		return DefaultFilterWindow{nil, -1, 0}
	}
	start--
	for ; k < len(f.opcodes) && match; k++ {
		opcode := f.opcodes[k]
		match = strings.EqualFold(opcode.name, name)
	}
	if !match {
		k--
	}
	if len(f.opcodes[start].paramsName) > 0 && strings.EqualFold(f.opcodes[start].paramsName[0], byteListKeyword) {
		return ByteListFilterWindow{f.opcodes[start : start+1], false, start + f.firstIndex}
	}
	return DefaultFilterWindow{f.opcodes[start:k], 0, start + f.firstIndex}
}

//FilterByToken filter the window using the token in the current position
//return a new filter window and if the token must be saved
func (f DefaultFilterWindow) FilterByToken(token text.Token) (FilterWindow, bool) {
	match := false
	keepToken := false
	var k int
	for k = 0; k < len(f.opcodes) && !match; k++ {
		opcode := f.opcodes[k]
		if len(opcode.paramsFormat) > f.matchIndex {
			isValue := token.EnumType() == text.Identifier || token.EnumType() == text.Number
			match = opcode.paramsFormat[f.matchIndex] == text.Identifier && isValue
			match = match || (opcode.paramsFormat[f.matchIndex] != text.Identifier && strings.EqualFold(opcode.paramsString[f.matchIndex], token.Value()))
			keepToken = isValue
		}
	}
	start := k - 1
	if !match {
		return DefaultFilterWindow{nil, -1, 0}, false
	}
	for ; k < len(f.opcodes) && match; k++ {
		opcode := f.opcodes[k]
		if len(opcode.paramsFormat) > f.matchIndex {
			isValue := token.EnumType() == text.Identifier || token.EnumType() == text.Number
			match = opcode.paramsFormat[f.matchIndex] == text.Identifier && isValue
			match = match || (opcode.paramsFormat[f.matchIndex] != text.Identifier && strings.EqualFold(opcode.paramsString[f.matchIndex], token.Value()))
		} else {
			match = false
		}
	}
	if !match {
		k--
	}
	return DefaultFilterWindow{f.opcodes[start:k], f.matchIndex + 1, start + f.firstIndex}, keepToken
}

//Collect return the true if there is exactly 1 opcode left and opcode in the first index
func (f DefaultFilterWindow) Collect() (bool, int) {
	if len(f.opcodes) < 1 {
		return false, 0
	}
	return len(f.opcodes[0].paramsFormat) == f.matchIndex, f.firstIndex
}

//IsEmpty check if the window has no more opcodes
func (f DefaultFilterWindow) IsEmpty() bool {
	return len(f.opcodes) == 0
}

func (f DefaultFilterWindow) String() string {
	return "std_filter_window"
}

//ByteListFilterWindow is a filter window for byte list opcodes
type ByteListFilterWindow struct {
	opcodes     []Opcode
	expectComma bool
	index       int
}

//FilterByName in truth makes no sense, useful only to implement the FilterWindow interface
//probably in future should return an error
func (f ByteListFilterWindow) FilterByName(name string) FilterWindow {
	return f
}

//FilterByToken check for match of the byteList or ',' after identifier/value and identifier/value after ','
func (f ByteListFilterWindow) FilterByToken(token text.Token) (FilterWindow, bool) {
	keepToken := !f.expectComma
	noError := f.expectComma && token.EnumType() == text.Comma ||
		!f.expectComma && (token.EnumType() == text.Identifier || token.EnumType() == text.DoubleQuotedString || token.EnumType() == text.Number)
	if noError {
		f.expectComma = !f.expectComma
		return f, keepToken
	}
	f.opcodes = nil
	return f, false
}

//Collect return the true if there is exactly 1 opcode left and opcode in the first index
func (f ByteListFilterWindow) Collect() (bool, int) {
	if len(f.opcodes) < 1 || !f.expectComma {
		return false, 0
	}
	return true, f.index
}

//IsEmpty check if the window has no more opcodes
func (f ByteListFilterWindow) IsEmpty() bool {
	return len(f.opcodes) == 0
}

//String return a stringed version of f
func (f ByteListFilterWindow) String() string {
	if len(f.opcodes) != 0 {
		return "byte_list_filter_window " + f.opcodes[0].name
	}
	return "empty byte_list_filter_window"
}
