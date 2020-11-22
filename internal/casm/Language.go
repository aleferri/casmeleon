package language

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
	OPCODE_NODE  = 4
	OPCODE_ARGS  = 5
	WITH_TYPES   = 6
	STMT_WARNING = 7
	STMT_ERROR   = 8
	STMT_BRANCH  = 9
	STMT_BLOCK   = 10
	STMT_OUT     = 11
)
