package casm

import (
	"github.com/aleferri/casmeleon/pkg/scanner"
	"github.com/aleferri/casmeleon/pkg/text"
)

var scanFollowMap = scanner.FromMap(map[rune]scanner.Follow{
	'\n': scanner.FollowNone,
	'\r': scanner.FollowNone,
	' ':  scanner.FollowSpaces,
	'\t': scanner.FollowNone,
	'&':  scanner.FollowSequence('&', '&'),
	'|':  scanner.FollowSequence('|', '|'),
	'^':  scanner.FollowNone,
	'!':  scanner.FollowSequence('!', '='),
	'<':  scanner.FollowComparison,
	'>':  scanner.FollowComparison,
	'=':  scanner.FollowSequence('=', '='),
	'*':  scanner.FollowCommentClose,
	'+':  scanner.FollowNone,
	'/':  scanner.FollowCommentOpen,
	'-':  scanner.FollowSequence('-', '>'),
	'%':  scanner.FollowNone,
	'@':  scanner.FollowNone,
	'#':  scanner.FollowNone,
	',':  scanner.FollowNone,
	';':  scanner.FollowNone,
	':':  scanner.FollowNone,
	'(':  scanner.FollowNone,
	')':  scanner.FollowNone,
	'[':  scanner.FollowNone,
	']':  scanner.FollowNone,
	'{':  scanner.FollowSequence('{', '{'),
	'}':  scanner.FollowSequence('}', '}'),
	'"':  scanner.FollowNone,
	'\'': scanner.FollowNone,
})

var identifyMap = map[string]uint32{
	".if": text.KeywordIF, ".else": text.KeywordELSE, ".opcode": text.KeywordOpcode, ".with": text.KeywordWith, ".num": text.KeywordNum,
	".set": text.KeywordSet, ".out": text.KeywordOut, ".outr": text.KeywordOutR, ".expr": text.KeywordExpr, ".error": text.KeywordError,
	".warning": text.KeywordWarning, ".inline": text.KeywordInline, "&": text.OperatorAnd, "&&": text.OperatorLAnd, "|": text.OperatorOr,
	"||": text.OperatorLOr, "+": text.OperatorPlus, "-": text.OperatorMinus, "*": text.OperatorMul, "/": text.OperatorDiv, "%": text.OperatorMod,
	"^": text.OperatorXor, "!": text.OperatorNot, "~": text.OperatorNeg, "<": text.OperatorLess, "<=": text.OperatorLessEqual,
	"==": text.OperatorEqual, ">=": text.OperatorGreaterEqual, ">": text.OperatorGreater, "!=": text.OperatorNotEqual, ".atom": text.KeywordAtom,
	"<<": text.OperatorLeftShift, ">>": text.OperatorRightShift, "->": text.SymbolArrow, "#": text.SymbolHash, "@": text.SymbolHash,
	"{{": text.DoubleCurlyOpen, "}}": text.DoubleCurlyClose, ".return": text.KeywordReturn, "{": text.CurlyOpen, "}": text.CurlyClose,
	"(": text.RoundOpen, ")": text.RoundClose, "[": text.SquareOpen, "]": text.SquareClose, ";": text.Semicolon, ":": text.Colon, ",": text.Comma,
}

var idDescriptor = []string{
	"No Token", "End of Line", "End Of File", "Whitespace", "(", ")", "[", "]", "{", "}", "{{", "}}", ",", ":", ";",
	"@", "#", "->", "/*", "*/", "//", "Quoted String", "Quoted Char", "Unary +", "+", "Unary -", "-", "*", "/", "%", ">>", "<<", "&", "&&",
	"|", "||", "^", "!", "~", "<", "==", "<=", ">=", ">", "!=", ".if Keyword", ".else Keyword", ".out Keyword", ".outr Keyword", ".set Keyword",
	".num Keyword", ".atom Keyword", ".inline Keyword", ".opcode Keyword", ".with Keyword", ".expr Keyword", ".warning Keyword", ".error Keyword",
	".return Keyword", "number", "identifier", "text label", "Errore di fuori indice",
}

var temporaryTokenMarks = map[int32]int32{1: 1, 2: 2, 3: 3, 4: 5}
