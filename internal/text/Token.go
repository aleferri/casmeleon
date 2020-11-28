package text

//TokenType is the token type of the token
type TokenType int

//TokenType values
const (
	LParen = 0 + iota
	RParen
	LBrace
	RBrace
	LBracket
	RBracket
	Comma
	Semicolon
	Colon
	SymbolHash
	SymbolAt
	SymbolDollar
	SymbolArrow
	EOL
	EOF
	UnaryOperator
	BinaryOperator
	Identifier
	Number
	SingleQuotedString
	DoubleQuotedString
	KeywordIF
	KeywordELSE
	KeywordFOR
	KeywordUNTIL
	KeywordEnum
	KeywordOpcode
	DirDeposit
	DirNumberFormat
	DirError
	KeywordInclude
	GenericSeparator
	TokenCommand
)

//String method of TokenType
func (tt TokenType) String() string {
	var tokenTypeStrings = [...]string{"LPAREN", "RPAREN", "LBRACE", "RBRACE", "LBRACKET", "RBRACKET", "COMMA", "SEMICOLON",
		"HASH", "AT", "DOLLAR", "EOL", "EOF", "UNARY", "BINARY", "ARROW", "IDENTIFIER", "NUMBER", "SINGLE_QUOTED", "DOUBLE_QUOTED",
		"IF", "ELSE", "FOR", "UNTIL", "DEPOSIT", "BASE_FORMAT", "ENUM", "OPCODE", "ERROR_DIR", "INCLUDE_DIR", "SEPARATOR", "COMMAND"}
	if tt <= TokenCommand {
		return tokenTypeStrings[tt]
	}
	return ""
}

//Token of parsing, extends lexing.Token adding a TokenType to the token
type Token struct {
	value     string
	tType     TokenType
	lineIndex uint
	line      *SourceLine
}

//NewToken create a new Token
func NewToken(value string, tType TokenType, lineIndex uint, line *SourceLine) Token {
	return Token{value, tType, lineIndex, line}
}

//NewInternalToken create new token with internal identifier
//Testing only
func NewInternalToken(value string) Token {
	return NewToken(value, Identifier, 0, NewSourceLine(nil, 0, "internal"))
}

//NewSpecialToken create a special token that have the same position of another one
func NewSpecialToken(derived Token, value string, tType TokenType) Token {
	return NewToken(value, tType, derived.lineIndex, derived.line)
}

//Position return the position of this token in the source (as collection of sources)
func (t Token) Position() (line uint, pos uint, file string) {
	return t.line.LineNumber(), t.lineIndex, t.line.SourceName()
}

//EnumType return the type assigned by parser
func (t Token) EnumType() TokenType {
	return t.tType
}

//Value return the value assigned by parser
func (t Token) Value() string {
	return t.value
}

//WithType return a copy of the token with specified type
func (t Token) WithType(tType TokenType) Token {
	if tType == t.tType {
		return t
	}
	return Token{t.value, tType, t.lineIndex, t.line}
}

//WithValue return a copy of the token with specified value
func (t Token) WithValue(value string) Token {
	return Token{value, t.tType, t.lineIndex, t.line}
}

//TokenTypeMap contains separators and keyword tokens
var TokenTypeMap = map[string]TokenType{
	"(": LParen, ")": RParen, "[": LBracket, "]": RBracket, "{": LBrace, "}": RBrace,
	",": Comma, ";": Semicolon, "if": KeywordIF, "else": KeywordELSE, ".db": DirDeposit, ".dw": DirDeposit,
	".dd": DirDeposit, ".opcode": KeywordOpcode, ".enum": KeywordEnum, ".number": DirNumberFormat, ".error": DirError,
	"->": SymbolArrow, "$": SymbolDollar, "#": SymbolHash, "@": SymbolAt, ".include": KeywordInclude, "for": KeywordFOR,
	"until": KeywordUNTIL,
}
