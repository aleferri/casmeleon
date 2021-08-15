package casm

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/aleferri/casmeleon/pkg/expr"
	"github.com/aleferri/casmeleon/pkg/parser"
	"github.com/aleferri/casmvm/pkg/opcodes"
	"github.com/aleferri/casmvm/pkg/vm"
)

//Language definition
type Language struct {
	numberBases []NumberBase
	sets        []Set
	opcodes     []Opcode
	fnList      []vm.Callable
	fnNames     []string
	endianess   bool // 0 big endian, 1 little endian
}

func (l *Language) FindAddressOf(name string) (uint32, bool) {
	for i, c := range l.fnNames {
		if c == name {
			return uint32(i), true
		}
	}
	return 0, false
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

func (lang *Language) ParseInt(value string) (int64, error) {
	if value[0] == '-' {
		fmt.Println("Negative number found")
		v, e := lang.ParseUint(value[1:])
		return -int64(v), e
	}
	v, e := lang.ParseUint(value)
	return int64(v), e
}

func (lang *Language) AssignFrame(c vm.Callable, name string) int32 {
	n := len(lang.fnList)
	lang.fnList = append(lang.fnList, c)
	lang.fnNames = append(lang.fnNames, name)
	return int32(n)
}

func (lang *Language) Executables() []vm.Callable {
	return lang.fnList
}

func (lang *Language) Endianess() bool {
	return lang.endianess
}

func MakeLanguage(root parser.CSTNode) (Language, error) {
	labels := Set{"_FormatLabels", 0, func(string) int32 { return 0 }}
	integers := Set{"Ints", 1, func(a string) int32 {
		v, _ := strconv.ParseInt(a, 10, 32)
		return int32(v)
	}}
	lang := Language{[]NumberBase{}, []Set{labels, integers}, []Opcode{}, []vm.Callable{}, []string{}, true}
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

				list, _, errBody := CompileListing(&lang, inline.params, body, nil)
				if errBody != nil {
					return lang, errBody
				}

				callable := vm.MakeCallable(*list)

				lang.fnList = append(lang.fnList, callable)
				lang.fnNames = append(lang.fnNames, inline.name)
			}
		case OPCODE_NODE:
			{
				opcode, body, err := PruneToOpcode(&lang, k)
				if err != nil {
					return lang, err
				}
				lang.opcodes = append(lang.opcodes, opcode)
				list, useAddr, errBody := CompileListing(&lang, opcode.params, body, nil)
				if errBody != nil {
					fmt.Printf("Arguments were: %v\n", opcode.params)
					return lang, errors.New("In Opcode " + opcode.name + ":\n" + errBody.Error())
				}
				lang.fnList[opcode.frame] = vm.MakeCallable(*list)
				lang.opcodes[len(lang.opcodes)-1].runList = *list
				lang.opcodes[len(lang.opcodes)-1].useAddr = useAddr
			}
		}

	}
	return lang, nil
}

func CompileListing(lang *Language, params []string, root parser.CSTNode, listing *[]opcodes.Opcode) (*[]opcodes.Opcode, bool, error) {
	if listing == nil {
		listing = &[]opcodes.Opcode{}
	}

	useAddr := false
	nextLocal := uint16(len(params))

	for _, node := range root.Children() {
		switch node.ID() {
		case STMT_BRANCH:
			{
				children := node.Children()
				status := expr.MakeConverter(children[0].Symbols(), nextLocal)

				err := CompileExpression(lang, params, listing, &status)
				if err != nil {
					return nil, false, err
				}

				useAddr = useAddr || status.HasFlag(USE_THIS_ADDR)

				taken, tUseAddr, bodyErr := CompileListing(lang, params, children[1], nil)
				if bodyErr != nil {
					return nil, false, bodyErr
				}

				useAddr = useAddr || tUseAddr

				takenLen := len(*taken)

				if len(children) > 2 {
					notTaken, fUseAddr, elseErr := CompileListing(lang, params, children[2], nil)
					if elseErr != nil {
						return nil, false, elseErr
					}

					useAddr = useAddr || fUseAddr

					notTakenLen := len(*notTaken)

					brElse := opcodes.MakeBranch(0, status.Pop().Local(), int32(takenLen)+1)
					*listing = append(*listing, brElse)
					*listing = append(*listing, *taken...)
					brExit := opcodes.MakeGoto(int32(notTakenLen))
					*listing = append(*listing, brExit)
					*listing = append(*listing, *notTaken...)
				} else {
					brExit := opcodes.MakeBranch(0, status.Pop().Local(), int32(takenLen))
					*listing = append(*listing, brExit)
					*listing = append(*listing, *taken...)
				}
				nextLocal = status.LabelLocal()
			}
		case STMT_ERROR:
			{
				*listing = append(*listing, opcodes.MakeSigError("error, will implement later", 0))
			}
		case STMT_RET:
			{
				status := expr.MakeConverter(node.Children()[0].Symbols(), nextLocal)
				err := CompileExpression(lang, params, listing, &status)
				if err != nil {
					return listing, false, err
				}
				useAddr = useAddr || status.HasFlag(USE_THIS_ADDR)
				*listing = append(*listing, opcodes.MakeLeave(status.Pop().Local()))
				nextLocal = status.LabelLocal()
			}
		case STMT_OUT:
			{
				refs := []uint16{}
				for _, item := range node.Children() {
					itemStatus := expr.MakeConverter(item.Symbols(), nextLocal)
					err := CompileExpression(lang, params, listing, &itemStatus)
					if err != nil {
						return listing, false, errors.New("In .out statement:\n" + err.Error())
					}
					useAddr = useAddr || itemStatus.HasFlag(USE_THIS_ADDR)
					refs = append(refs, itemStatus.Pop().Local())
					nextLocal = itemStatus.LabelLocal()
				}
				*listing = append(*listing, opcodes.MakeLeave(refs...))
			}
		case STMT_OUTR:
			{
				refs := make([]uint16, len(node.Children()))
				index := len(node.Children()) - 1
				for _, item := range node.Children() {
					itemStatus := expr.MakeConverter(item.Symbols(), nextLocal)
					err := CompileExpression(lang, params, listing, &itemStatus)
					if err != nil {
						return listing, false, errors.New("In .outr statement:\n" + err.Error())
					}
					useAddr = useAddr || itemStatus.HasFlag(USE_THIS_ADDR)
					refs[index] = itemStatus.Pop().Local()
					index--
					nextLocal = itemStatus.LabelLocal()
				}
				*listing = append(*listing, opcodes.MakeLeave(refs...))
			}
		case STMT_WARNING:
			{
				*listing = append(*listing, opcodes.MakeSigWarning("error, will implement later", 0))
			}
		}
	}

	return listing, useAddr, nil
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
	STMT_OUTR    = 20
)
