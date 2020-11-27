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

		switch state {
		case ACCUMULATE:
			{
				followFunc, ok = settings.delimiters[b]
				if ok {
					tokens = append(tokens, Token{buffer[lastDelimiter:i], -1})
					lastDelimiter = i
					state = FOLLOW
				}
			}
		case FOLLOW:
			{
				followState = followFunc(b, followState)
				if followState == 0 {
					tokens = append(tokens, Token{buffer[lastDelimiter:i], -1})
					lastDelimiter = i
					state = ACCUMULATE
					followState = 1
				}
			}
		}
	}

	if stop && lastDelimiter != len {
		tokens = append(tokens, Token{buffer[lastDelimiter:], -1})
		return tokens, []rune{}
	}
	return tokens, buffer[lastDelimiter:]
}

//Join matching tokens
func Join(matching map[int32]int32, tokens []Token, merged *Token, last int32) ([]Token, *Token, int32) {
	valid := []Token{}
	for _, t := range tokens {
		p, ok := matching[t.basicID]
		if ok {
			if last == p {
				merged = merged.Merge(t)
				valid = append(valid, *merged)
				merged = nil
				last = -1
			} else {
				merged = &t
			}
		} else {
			if last != -1 {
				merged = merged.Merge(t)
			} else {
				valid = append(valid, t)
			}
		}
	}
	return valid, merged, last
}

//Classify basic id
func Classify(tokens []Token, classifier func(t *Token) int32) {
	for _, t := range tokens {
		t.basicID = classifier(&t)
	}
}
