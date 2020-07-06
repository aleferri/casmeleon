package langdef

import (
	"bitbucket.org/mrpink95/casmeleon/internal/parsing"
	"bitbucket.org/mrpink95/casmeleon/internal/text"
)

//Tag values, Ignore and OnlyAppend respect defaults declared in parsing
const (
	Ignore = 0 + iota
	OnlyAppend
	TagNumberFormat
	TagEnumDeclaration
	TagErrorStatement
	TagDepositStatement
	TagElseStatement
	TagBranchStatement
	TagForStatement
	TagExpression
	TagStatement
	TagBlock
	TagOpcode
)

func setSkipEOL(skipEOL bool) parsing.MatchRule {
	return func(tokens []text.Token, pState parsing.ParserState) ([]text.Token, bool) {
		langState := pState.(*LangParserState)
		langState.SetSkipEOL(skipEOL, &tokens)
		return tokens, true
	}
}

func createNumberStatementParser() parsing.MatchRule {
	enableEOL := setSkipEOL(false)
	disableEOL := setSkipEOL(true)
	expectedNumberDir := parsing.NewIncompleteError(parsing.ErrorExpectedToken, ", expected '.number'")
	expectedNumberBase := parsing.NewIncompleteError(parsing.ErrorExpectedToken, ", expected '.hex' or '.bin' or '.dec' or '.oct'")
	expectedNumberPosition := parsing.NewIncompleteError(parsing.ErrorExpectedToken, ", expected 'suffix' or 'prefix'")
	expectedSingleQuoteString := parsing.NewIncompleteError(parsing.ErrorExpectedToken, ", expected single quote string")
	expectedEOLNumber := parsing.NewIncompleteError(parsing.ErrorExpectedEOL, "after '.number' statement")
	matchNumberDir := parsing.MatchToken(expectedNumberDir, text.DirNumberFormat)
	matchNumberBase := parsing.MatchAnyString(expectedNumberBase, ".bin", ".dec", ".hex", ".oct")
	matchNumberPosition := parsing.MatchAnyString(expectedNumberPosition, "suffix", "prefix")
	matchSingleQuoteString := parsing.MatchToken(expectedSingleQuoteString, text.SingleQuotedString)
	matchEOLNumber := parsing.MatchAnyToken(expectedEOLNumber, text.EOL, text.EOF)
	return parsing.MatchAll(TagNumberFormat, enableEOL, matchNumberDir, matchNumberBase, matchNumberPosition, matchSingleQuoteString, matchEOLNumber, disableEOL)
}

func createEnumDeclarationParser() parsing.MatchRule {
	expectedEnumDir := parsing.NewIncompleteError(parsing.ErrorExpectedToken, ", expected '.enum'")
	expectedIdentifier := parsing.NewIncompleteError(parsing.ErrorExpectedToken, ", expected identifier")
	expectedLeftBrace := parsing.NewIncompleteError(parsing.ErrorExpectedToken, ", expected '{'")
	expectedRightBrace := parsing.NewIncompleteError(parsing.ErrorExpectedToken, ", expected '}'")
	testComma := parsing.TryMatchToken(text.Comma, true)
	testNotRightBrace := parsing.TryNotMatchToken(text.RBrace, false)
	matchIdentifier := parsing.MatchToken(expectedIdentifier, text.Identifier)
	matchRepeatIdentifier := parsing.TryMatchRepeat(OnlyAppend, testComma, matchIdentifier)
	matchEnumDir := parsing.MatchToken(expectedEnumDir, text.KeywordEnum)
	matchLeftBrace := parsing.MatchSkipToken(expectedLeftBrace, text.LBrace)
	matchFirstIdentifier := parsing.TryMatch(OnlyAppend, testNotRightBrace, matchIdentifier)
	matchRightBrace := parsing.MatchSkipToken(expectedRightBrace, text.RBrace)
	return parsing.MatchAll(TagEnumDeclaration, matchEnumDir, matchIdentifier, matchLeftBrace, matchFirstIdentifier, matchRepeatIdentifier, matchRightBrace)
}

func createErrorStatementParser() parsing.MatchRule {
	expectedErrorDir := parsing.NewIncompleteError(parsing.ErrorExpectedToken, ", expected '.error'")
	expectedIdentifier := parsing.NewIncompleteError(parsing.ErrorExpectedToken, ", expected arg name after '.error'")
	expectedDoubleQuote := parsing.NewIncompleteError(parsing.ErrorExpectedToken, ", expected double quoted error string")
	matchErrorDir := parsing.MatchToken(expectedErrorDir, text.DirError)
	matchIdentifier := parsing.MatchToken(expectedIdentifier, text.Identifier)
	matchDoubleQuote := parsing.MatchToken(expectedDoubleQuote, text.DoubleQuotedString)
	expectedSemicolon := parsing.NewIncompleteError(parsing.ErrorExpectedToken, ", expected ';'")
	matchSemicolon := parsing.MatchSkipToken(expectedSemicolon, text.Semicolon)
	return parsing.MatchAll(TagErrorStatement, matchErrorDir, matchIdentifier, matchDoubleQuote, matchSemicolon)
}

func createDepositStatementParser(matchExpression parsing.MatchRule) parsing.MatchRule {
	expectedDepositDir := parsing.NewIncompleteError(parsing.ErrorExpectedToken, ", expected '.db' or '.dw' or '.dd'")
	matchDepositDir := parsing.MatchToken(expectedDepositDir, text.DirDeposit)
	expectedSemicolon := parsing.NewIncompleteError(parsing.ErrorExpectedToken, ", expected ';'")
	matchSemicolon := parsing.MatchToken(expectedSemicolon, text.Semicolon)
	return parsing.MatchAll(TagDepositStatement, matchDepositDir, matchExpression, matchSemicolon)
}

func createLoopParser(matchExpression parsing.MatchRule, matchDeposit parsing.MatchRule) parsing.MatchRule {
	expectedFOR := parsing.NewIncompleteError(parsing.ErrorExpectedToken, ", expected 'for' in loop")
	expectedIdentifier := parsing.NewIncompleteError(parsing.ErrorExpectedToken, ", expected identfier name after for")
	expectedUNTIL := parsing.NewIncompleteError(parsing.ErrorExpectedToken, ", expected 'until' keyword")
	matchFOR := parsing.MatchToken(expectedFOR, text.KeywordFOR)
	matchIdentifier := parsing.MatchToken(expectedIdentifier, text.Identifier)
	matchUNTIL := parsing.MatchToken(expectedUNTIL, text.KeywordUNTIL)
	return parsing.MatchAll(TagForStatement, matchFOR, matchIdentifier, matchUNTIL, matchExpression, matchDeposit)
}

func createBranchParser(matchExpression parsing.MatchRule, matchBlock parsing.MatchRule) parsing.MatchRule {
	expectedIF := parsing.NewIncompleteError(parsing.ErrorExpectedToken, ", expected 'if' in conditional")
	matchIF := parsing.MatchSkipToken(expectedIF, text.KeywordIF)
	testELSE := parsing.TryMatchToken(text.KeywordELSE, true)
	matchElse := parsing.TryMatch(TagElseStatement, testELSE, matchBlock)
	return parsing.MatchAll(TagBranchStatement, matchIF, matchExpression, matchBlock, matchElse)
}

func createBlockParser() (parsing.MatchRule, parsing.MatchRule) {
	matchExpression := MatchExpression
	expectedLeftBrace := parsing.NewIncompleteError(parsing.ErrorExpectedToken, ", expected '{'")
	expectedRightBrace := parsing.NewIncompleteError(parsing.ErrorExpectedToken, ", expected '}'")
	matchLeftBrace := parsing.MatchSkipToken(expectedLeftBrace, text.LBrace)
	matchRightBrace := parsing.MatchSkipToken(expectedRightBrace, text.RBrace)
	matchMap := parsing.NewMultiMatch()
	matchErrorStatement := createErrorStatementParser()
	matchMap.AddMatch(text.DirError, matchErrorStatement)
	matchDepositStatement := createDepositStatementParser(matchExpression)
	matchMap.AddMatch(text.DirDeposit, matchDepositStatement)
	matchLoop := createLoopParser(matchExpression, matchDepositStatement)
	matchMap.AddMatch(text.KeywordFOR, matchLoop)
	expectedStatement := parsing.NewIncompleteError(parsing.ErrorExpectedToken, ", expected loop or branch or deposit or error statement")
	matchStatement := matchMap.MatchWithMap(expectedStatement)
	testNotRightBrace := parsing.TryNotMatchToken(text.RBrace, false)
	matchStatements := parsing.TryMatchRepeat(TagStatement, testNotRightBrace, matchStatement)
	matchBlock := parsing.MatchAll(TagBlock, matchLeftBrace, matchStatements, matchRightBrace)
	matchBranch := createBranchParser(matchExpression, matchBlock)
	matchMap.AddMatch(text.KeywordIF, matchBranch)
	return matchBlock, matchBranch
}

func createOpcodeParser(matchBlock parsing.MatchRule) parsing.MatchRule {
	enableEOL := setSkipEOL(false)
	disableEOL := setSkipEOL(true)
	expectedOpcode := parsing.NewIncompleteError(parsing.ErrorExpectedToken, ", expected '.opcode")
	unexpectedEOL := parsing.NewIncompleteError(parsing.ErrorExpectedToken, ", expected '->' before EOL or EOF")
	expectedIS := parsing.NewIncompleteError(parsing.ErrorExpectedToken, ", expected '->' before EOL or EOF")
	matchOpcode := parsing.MatchToken(expectedOpcode, text.KeywordOpcode)
	matchNotEOL := parsing.MatchNotAnyToken(unexpectedEOL, text.EOF, text.EOL)
	testNotIS := parsing.TryNotMatchToken(text.SymbolArrow, false)
	matchArgs := parsing.TryMatchRepeat(OnlyAppend, testNotIS, matchNotEOL)
	matchIS := parsing.MatchToken(expectedIS, text.SymbolArrow)
	return parsing.MatchAll(TagOpcode, enableEOL, matchOpcode, matchArgs, matchIS, disableEOL, matchBlock)
}

func createLangParser() parsing.MatchRule {
	matchBlock, _ := createBlockParser()
	matchOpcode := createOpcodeParser(matchBlock)
	matchNumberStatement := createNumberStatementParser()
	matchEnumDeclaration := createEnumDeclarationParser()
	matchTopLevel := parsing.NewMultiMatch()
	matchTopLevel.AddMatch(text.KeywordOpcode, matchOpcode)
	matchTopLevel.AddMatch(text.DirNumberFormat, matchNumberStatement)
	matchTopLevel.AddMatch(text.KeywordEnum, matchEnumDeclaration)
	unexpectedToken := parsing.NewIncompleteError(parsing.ErrorExpectedToken, "'.opcode' or '.number' or '.enum'")
	testNotEOF := parsing.TryNotMatchToken(text.EOF, false)
	matchLanguage := parsing.TryMatchRepeat(Ignore, testNotEOF, matchTopLevel.MatchWithMap(unexpectedToken))
	return matchLanguage
}
