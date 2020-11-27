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
					tokens = append(tokens, Token{buffer[lastDelimiter:i]})
					lastDelimiter = i
					state = FOLLOW
				}
			}
		case FOLLOW:
			{
				followState = followFunc(b, followState)
				if followState == 0 {
					tokens = append(tokens, Token{buffer[lastDelimiter:i]})
					lastDelimiter = i
					state = ACCUMULATE
					followState = 1
				}
			}
		}
	}

	if stop && lastDelimiter != len {
		tokens = append(tokens, Token{buffer[lastDelimiter:]})
		return tokens, []rune{}
	}
	return tokens, buffer[lastDelimiter:]
}
