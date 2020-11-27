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

//Join matching tokens
func Join(matching map[int32]int32, tokens []Token, merged *Token, last int32) ([]Token, *Token, int32) {
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

//Classify basic id
func Classify(tokens []Token, classifier func(t *Token) int32) {
	for i := range tokens {
		tokens[i].basicID = classifier(&tokens[i])
	}
}
