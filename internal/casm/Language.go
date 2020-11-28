package casm

//Language definition
type Language struct {
	numberBases []NumberBase
	sets        []Set
}

//CST Tags
const (
	EXPRESSION   = 0
	NUMBER_BASE  = 1
	SYMBOL_SET   = 2
	SET_NODE     = 3
	INLINE_NODE  = 4
	OPCODE_NODE  = 5
	OPCODE_ARGS  = 6
	WITH_TYPES   = 7
	STMT_WARNING = 8
	STMT_ERROR   = 9
	STMT_BRANCH  = 10
	STMT_BLOCK   = 11
	STMT_OUT     = 12
	STMT_RET     = 13
)
