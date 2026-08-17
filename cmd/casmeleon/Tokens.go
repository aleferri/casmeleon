package main

import (
	"unicode"

	"github.com/aleferri/casmeleon/pkg/scanner"
	"github.com/aleferri/casmeleon/pkg/text"
)

var scanDelimiters = map[rune]scanner.Follow{
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
}

var scanFollowMap = scanner.FromMap(scanDelimiters)

// followNumberPrefix keeps accumulating alphanumeric runes after the prefix rune, so a
// literal like %1010 stays a single token instead of being split by the delimiter map.
func followNumberPrefix(r rune, state uint32) uint32 {
	if state == 1 {
		return 2
	}
	if unicode.IsLetter(r) || unicode.IsDigit(r) {
		return 2
	}
	return 0
}

// RelaxNumberPrefixDelimiters must be called after the language file has been parsed. A
// prefix declared by .num whose first rune is also a delimiter would otherwise be split
// off from its digits. Operators are not expression syntax in the assembly source - they
// are literal pattern tokens - so repurposing such a rune costs nothing.
// A rune that is not a declared prefix keeps its original Follow. scanDelimiters is left
// pristine and a copy is installed, so a relaxation cannot leak across languages.
func RelaxNumberPrefixDelimiters(prefixes []string) {
	relaxed := map[rune]scanner.Follow{}
	for r, f := range scanDelimiters {
		relaxed[r] = f
	}
	changed := false
	for _, p := range prefixes {
		if p == "" {
			continue
		}
		first := []rune(p)[0]
		if unicode.IsDigit(first) {
			continue
		}
		if _, isDelimiter := relaxed[first]; isDelimiter {
			relaxed[first] = followNumberPrefix
			changed = true
		}
	}
	if changed {
		scanFollowMap = scanner.FromMap(relaxed)
	}
}

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
