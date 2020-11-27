package scanner

//DelimitersMap is the map of the delimiters and their follow
type DelimitersMap struct {
	delimiters map[rune]Follow
}

//FromMap create DelimitersMap
func FromMap(delimiters map[rune]Follow) DelimitersMap {
	return DelimitersMap{delimiters}
}

//FollowCommentOpen delimiter
func FollowCommentOpen(r rune, state uint32) uint32 {
	switch state {
	case 1:
		{
			return 2
		}
	case 2:
		{
			if r == '*' || r == '/' {
				return 3
			}
			return 0
		}
	case 3:
		{
			return 0
		}
	default:
		{
			return 0
		}
	}
}

//FollowCommentClose delimiter
func FollowCommentClose(r rune, state uint32) uint32 {
	switch state {
	case 1:
		{
			return 2
		}
	case 2:
		{
			if r == '/' {
				return 3
			}
			return 0
		}
	case 3:
		{
			return 0
		}
	default:
		{
			return 0
		}
	}
}

//FollowNone of the next chars, example: end of the line
func FollowNone(r rune, state uint32) uint32 {
	return 0
}

//FollowSpaces follow whitespace to accumulate
func FollowSpaces(r rune, state uint32) uint32 {
	if r == ' ' {
		return 1
	}
	return 0
}

//FollowSame rune until something different is found
func FollowSame(s rune) Follow {
	return func(r rune, state uint32) uint32 {
		switch state {
		case 1:
			{
				if r == s {
					return 1
				}
				return 0
			}
		default:
			{
				return 0
			}
		}
	}
}

//FollowSequence of runes
func FollowSequence(seq ...rune) Follow {
	return func(r rune, state uint32) uint32 {
		i := state - 1
		if r == seq[i] {
			return state + 1
		}
		return 0
	}
}

//FollowComparison operators
func FollowComparison(r rune, state uint32) uint32 {
	switch state {
	case 1:
		{
			if r == '<' {
				return 3
			}
			return 5
		}
	case 2:
		{
			if r == '=' {
				return 3
			} else if r == '<' {
				return 4
			}
			return 0
		}
	case 3:
		{
			return 0
		}
	case 4:
		{
			if r == '<' {
				return 3
			}
			return 0
		}
	case 5:
		{
			if r == '=' {
				return 6
			} else if r == '>' {
				return 7
			}
			return 0
		}
	case 6:
		{
			return 0
		}
	case 7:
		{
			if r == '>' {
				return 3
			}
			return 0
		}
	default:
		{
			return 0
		}
	}
}
