package scanner

//FastScan scan the whole buffer for easy separators, without attributing the index or checking for errors
func FastScan(buffer []rune, stop bool, settings DelimitersMap) (toks []Token, left []rune) {
	const ACCUMULATE = 0
	const FOLLOW = 1
	state := ACCUMULATE
	followState := uint32(1)
	var followFunc Follow = nil
	var ok bool
	tokens := []Token{}
	lastDelimiter := 0
	len := len(buffer)
	for i, b := range buffer {

		if state == FOLLOW {
			followState = followFunc(b, followState)
			if followState == 0 {
				tokens = append(tokens, Token{buffer[lastDelimiter:i], 0})
				lastDelimiter = i
				state = ACCUMULATE
			}
		}

		if state == ACCUMULATE {
			followFunc, ok = settings.delimiters[b]
			if ok {
				tokens = append(tokens, Token{buffer[lastDelimiter:i], 0})
				lastDelimiter = i
				state = FOLLOW
				followState = followFunc(b, 1)
			}
		}

	}

	if stop && lastDelimiter != len {
		tokens = append(tokens, Token{buffer[lastDelimiter:], 0})
		return tokens, []rune{}
	}
	return tokens, buffer[lastDelimiter:]
}

//Merge matching tokens
func Merge(matching map[int32]int32, tokens []Token, merged *Token, last int32) ([]Token, *Token, int32) {
	valid := []Token{}
	for i, t := range tokens {
		if last != -1 {
			*merged = merged.Merge(t)
			if last == t.basicID {
				valid = append(valid, *merged)
				last = -1
				merged = nil
			}
		} else {
			p, ok := matching[t.basicID]
			if ok {
				merged = &tokens[i]
				last = p
			} else {
				if len(t.slice) > 0 {
					valid = append(valid, t)
				}
			}
		}
	}
	return valid, merged, last
}

//ClassifyMergeableTokens for successive Join
func ClassifyMergeableTokens(tokens []Token) {
	for i := range tokens {
		size := len(tokens[i].slice)
		if size == 2 {
			a := tokens[i].slice[0]
			b := tokens[i].slice[1]
			if a == '/' && b == '*' {
				tokens[i].basicID = 1
			}
			if a == '/' && b == '/' {
				tokens[i].basicID = 4
			}
		} else if size == 1 {
			a := tokens[i].slice[0]
			if a == '"' {
				tokens[i].basicID = 2
			} else if a == '\'' {
				tokens[i].basicID = 3
			} else if a == '\n' {
				tokens[i].basicID = 5
			}
		}
	}
}

//MergeASMLine matching tokens
func MergeASMLine(line []Token) []Token {
	valid := []Token{}
	var merged *Token = nil
	match := int32(-1)
	for i, t := range line {
		if t.basicID == 6 {
			continue
		}
		if t.basicID == 5 {
			if match != -1 {
				valid = append(valid, *merged)
				merged = nil
				match = -1
			}
			valid = append(valid, t)
		} else {
			if match != -1 {
				*merged = merged.Merge(t)
				if match == t.basicID {
					valid = append(valid, *merged)
					match = -1
					merged = nil
				}
			} else {
				if t.basicID == 1 || t.basicID == 2 || t.basicID == 3 {
					merged = &line[i]
					match = t.basicID
					if t.basicID == 3 {
						match = 7
					}
				} else if len(t.slice) > 0 {
					valid = append(valid, t)
				}
			}
		}
	}

	return valid
}

func ClassifyBasicASMTokens(tokens []Token) {
	for i := range tokens {
		size := len(tokens[i].slice)
		if size == 1 {
			a := tokens[i].slice[0]
			if a == '"' {
				tokens[i].basicID = 2
			} else if a == '\'' {
				tokens[i].basicID = 1
			} else if a == '\n' {
				tokens[i].basicID = 5
			} else if a == ';' {
				tokens[i].basicID = 3
			} else if a == '\r' {
				tokens[i].basicID = 6
			}
		}
	}
}
