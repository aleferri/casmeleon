package casm

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/aleferri/casmeleon/pkg/exec"
	"github.com/aleferri/casmeleon/pkg/parser"
)

//Language definition
type Language struct {
	numberBases []NumberBase
	sets        []Set
	opcodes     []Opcode
	inlines     []Inline
	addrUsed    bool
}

//SetOf return the Set of a symbol
func (l *Language) SetOf(s string) (*Set, bool) {
	lenSets := len(l.sets) - 1
	for i := range l.sets {
		set := l.sets[lenSets-i]
		if set.Contains(s) {
			return &set, true
		}
	}
	return nil, false
}

//SetByName return the set by his name
func (l *Language) SetByName(s string) (*Set, bool) {
	for _, set := range l.sets {
		if set.name == s {
			return &set, true
		}
	}
	return nil, false
}

func (l *Language) InlineCode(name string) (exec.Executable, error) {
	for _, i := range l.inlines {
		if i.name == name {
			return &InlineCall{i}, nil
		}
	}
	return nil, errors.New("No inline named " + name + " exists")
}

func (l *Language) MarkAddressUsed() {
	l.addrUsed = true
}

func (l *Language) FilterOpcodesByName(name string) FilterWindow {
	wnd := FilterWindow{[]string{}, []Opcode{}}
	for _, op := range l.opcodes {
		if op.name == name {
			wnd.filtered = append(wnd.filtered, op)
		}
	}
	return wnd
}

func (lang *Language) ParseUint(value string) (uint64, error) {
	for _, base := range lang.numberBases {
		if base.prefix != "" {
			if !strings.HasPrefix(value, base.prefix) {
				continue
			}
		}
		if base.suffix != "" {
			if !strings.HasSuffix(value, base.suffix) {
				continue
			}
		}
		count := len(value)
		cut := value[len(base.prefix) : count-len(base.suffix)]
		return strconv.ParseUint(cut, int(base.n), 64)
	}
	return strconv.ParseUint(value, 10, 64)
}

func MakeLanguage(root parser.CSTNode) (Language, error) {
	labels := Set{"_FormatLabels", 0, func(string) int32 { return 0 }}
	integers := Set{"Ints", 1, func(a string) int32 {
		v, _ := strconv.ParseInt(a, 10, 32)
		return int32(v)
	}}
	lang := Language{[]NumberBase{}, []Set{labels, integers}, []Opcode{}, []Inline{}, false}
	for _, k := range root.Children() {
		switch k.ID() {
		case NUMBER_BASE:
			{
				base, err := PruneToNumBase(k)
				if err != nil {
					return lang, err
				}
				lang.numberBases = append(lang.numberBases, base)
			}
		case SET_NODE:
			{
				set := PruneToSet(k, uint32(len(lang.sets)))
				lang.sets = append(lang.sets, set)
			}
		case INLINE_NODE:
			{
				inline, body, err := PruneToInline(&lang, k)
				if err != nil {
					return lang, err
				}
				lang.inlines = append(lang.inlines, inline)
				list, errBody := ExecutableListFromNode(&lang, inline.params, body, nil)
				if errBody != nil {
					return lang, errBody
				}
				lang.inlines[len(lang.inlines)-1].runList = *list
			}
		case OPCODE_NODE:
			{
				opcode, body, err := PruneToOpcode(&lang, k)
				if err != nil {
					return lang, err
				}
				lang.opcodes = append(lang.opcodes, opcode)
				list, errBody := ExecutableListFromNode(&lang, opcode.params, body, nil)
				if errBody != nil {
					fmt.Printf("Arguments were: %v\n", opcode.params)
					return lang, errors.New("In Opcode " + opcode.name + ":\n" + errBody.Error())
				}
				lang.opcodes[len(lang.opcodes)-1].runList = *list
				lang.opcodes[len(lang.opcodes)-1].useAddr = lang.addrUsed
				lang.addrUsed = false
			}
		}

	}
	return lang, nil
}

func ExecutableListFromNode(lang *Language, params []string, root parser.CSTNode, list *[]exec.Executable) (*[]exec.Executable, error) {
	if list == nil {
		list = &[]exec.Executable{}
	}

	for _, node := range root.Children() {
		switch node.ID() {
		case STMT_BRANCH:
			{
				_, err := CompileExpression(lang, params, list, node.Symbols())
				if err != nil {
					return nil, err
				}
				children := node.Children()
				taken, bodyErr := ExecutableListFromNode(lang, params, children[0], nil)
				if bodyErr != nil {
					return nil, bodyErr
				}

				if len(children) > 1 {
					notTaken, elseErr := ExecutableListFromNode(lang, params, children[1], nil)
					if elseErr != nil {
						return nil, elseErr
					}
					*list = append(*list, exec.MakeBranchCode(*taken, *notTaken))
				} else {
					*list = append(*list, exec.MakeBranchCode(*taken, nil))
				}
			}
		case STMT_ERROR:
			{
				*list = append(*list, exec.EmitErrorOf(0, "error, will implement later"))
			}
		case STMT_RET:
			{
				_, err := CompileExpression(lang, params, list, node.Children()[0].Symbols())
				if err != nil {
					return list, err
				}
				*list = append(*list, exec.MakeReturn())
			}
		case STMT_OUT:
			{
				outList := []exec.Executable{}
				for _, expr := range node.Children() {
					execList := []exec.Executable{}
					_, err := CompileExpression(lang, params, &execList, expr.Symbols())
					if err != nil {
						return list, errors.New("In .out statement:\n" + err.Error())
					}
					outList = append(outList, exec.BuildStackExpression(execList))
				}
				*list = append(*list, exec.MakeOutResult(outList))
			}
		case STMT_WARNING:
			{
				*list = append(*list, exec.EmitWarningOf(0, "error, will implement later"))
			}
		}
	}

	return list, nil
}

//CST Tags
const (
	EXPRESSION   = 0
	NUMBER_BASE  = 1
	SYMBOL_SET   = 2
	SET_NODE     = 3
	INLINE_NODE  = 4
	OPCODE_NODE  = 5
	OPCODEFORMAT = 6
	OPCODE_ARGS  = 7
	WITH_TYPES   = 8
	STMT_WARNING = 9
	STMT_ERROR   = 10
	STMT_BRANCH  = 11
	STMT_BLOCK   = 12
	STMT_OUT     = 13
	STMT_RET     = 14
	EXPR_CALL    = 15
	OPERAND      = 16
	BIN_OPERATOR = 17
	URY_OPERATOR = 18
	ROOT_NODE    = 19
)
