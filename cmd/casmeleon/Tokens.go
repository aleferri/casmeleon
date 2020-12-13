package main

import (
	"github.com/aleferri/casmeleon/pkg/scanner"
	"github.com/aleferri/casmeleon/pkg/text"
)

var scanFollowMap = scanner.FromMap(map[rune]scanner.Follow{
	'\n': scanner.FollowNone,
	'\r': scanner.FollowNone,
	' ':  scanner.FollowSpaces,
	'\t': scanner.FollowNone,
	'&':  scanner.FollowNone,
	'|':  scanner.FollowNone,
	'^':  scanner.FollowNone,
	'!':  scanner.FollowNone,
	'<':  scanner.FollowNone,
	'>':  scanner.FollowNone,
	'=':  scanner.FollowNone,
	'*':  scanner.FollowNone,
	'+':  scanner.FollowNone,
	'/':  scanner.FollowNone,
	'-':  scanner.FollowNone,
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
	'{':  scanner.FollowNone,
	'}':  scanner.FollowNone,
	'"':  scanner.FollowNone,
	'\'': scanner.FollowNone,
})

var identifyMap = map[string]uint32{
	"&": text.OperatorAnd, "&&": text.OperatorLAnd, "|": text.OperatorOr, "||": text.OperatorLOr, "+": text.OperatorPlus,
	"-": text.OperatorMinus, "*": text.OperatorMul, "/": text.OperatorDiv, "%": text.OperatorMod, "^": text.OperatorXor, "!": text.OperatorNot,
	"~": text.OperatorNeg, "<": text.OperatorLess, "<=": text.OperatorLessEqual, "==": text.OperatorEqual, ">=": text.OperatorGreaterEqual,
	">": text.OperatorGreater, "!=": text.OperatorNotEqual, "<<": text.OperatorLeftShift, ">>": text.OperatorRightShift, "->": text.SymbolArrow,
	"#": text.SymbolHash, "@": text.SymbolHash, "{": text.CurlyOpen, "}": text.CurlyClose, "(": text.RoundOpen, ")": text.RoundClose, "[": text.SquareOpen,
	"]": text.SquareClose, ";": text.Semicolon, ":": text.Colon, ",": text.Comma,
}

var idDescriptor = []string{
	"No Token", "End of Line", "End Of File", "Whitespace", "(", ")", "[", "]", "{", "}", "N/D", "N/D", ",", ":", ";",
	"@", "#", "->", "/*", "*/", "//", "Quoted String", "Quoted Char", "+", "-", "*", "/", "%", ">>", "<<", "&", "&&",
	"|", "||", "^", "!", "~", "<", "==", "<=", ">=", ">", "!=", ".if Keyword", ".else Keyword", ".out Keyword", ".outr Keyword", ".set Keyword",
	".num Keyword", ".atom Keyword", ".inline Keyword", ".opcode Keyword", ".with Keyword", ".expr Keyword", ".warning Keyword", ".error Keyword",
	".return Keyword", "Number", "Identifier", "Out of bounds",
}

var temporaryTokenMarks = map[int32]int32{1: 1, 2: 2, 3: 3, 4: 5}
