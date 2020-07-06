package lexing

import (
	"container/list"
	"strings"
)

// SplitKeepSeparators split a string considering single char separators like StringTokenizer in java,
// return a list with the tokens and delimiters
func SplitKeepSeparators(s, sep string) *list.List {
	var k = 0
	var result = list.New()
	var end = 0
	for i, c := range s {
		if strings.ContainsRune(sep, c) {
			if i <= k {
				result.PushBack(string(c))
				k++
			} else {
				result.PushBack(s[k:i])
				result.PushBack(string(c))
				k = i + 1
			}
		}
		end = i + 1
	}
	if k < end {
		result.PushBack(s[k:])
	}
	return result
}

func readUntilInvalid(l *list.List, symbols string, elem string, i int) (int, string) {
	var parts []string
	size := l.Len() + i
	symbol := elem
	indexBegin := strings.Index(symbols, " "+symbol)
	trySymbol := symbol
	for k := i; k < size; k++ {
		val := l.Front().Value.(string)
		trySymbol += val
		indexBegin = strings.Index(symbols, " "+trySymbol)
		if strings.EqualFold(val, " ") || indexBegin == -1 {
			break
		}
		parts = append(parts, val)
		l.Remove(l.Front())
		symbol = trySymbol
	}
	indexBegin = strings.Index(symbols, symbol)
	found := indexBegin != -1
	rightIsWhite := found && (indexBegin+len(symbol) == len(symbols) || symbols[indexBegin+len(symbol)] == ' ')
	if rightIsWhite {
		return i + len(parts), symbol
	}
	for k := len(parts) - 1; k > -1; k-- {
		l.PushFront(parts[k])
	}
	return i, elem
}

//RegroupSymbols regroup symbols split after NextLine()
//multiMatch format: "<symbol><space>...".
//example: "+ - ++ -- & |"
func RegroupSymbols(l *list.List, symbols string) []string {
	var size = l.Len()
	if size < 1 {
		return []string{}
	}
	symbols = " " + symbols
	realSize := 0
	result := make([]string, size)
	for i := 0; i < size; i++ {
		var elem = l.Front().Value.(string)
		l.Remove(l.Front())
		var indexBegin = strings.Index(symbols, " "+elem)
		if indexBegin > -1 {
			k, found := readUntilInvalid(l, symbols, elem, i)
			i = k
			result[realSize] = found
		} else {
			result[realSize] = elem
		}
		realSize++
	}
	return result[:realSize]
}

//JoinMultiCharSymbols join symbols with more than 1 char like && or >=.
//Max len of symbol is 2
//<symbol padded> = <symbol> <space to next even position>,
//multiMatch format: "<symbol padded><symbol padded>...".
//example: "+ - ++  --  & |"
func JoinMultiCharSymbols(l *list.List, multiChar string) *list.List {
	var listLen = l.Len()
	var i = 0
	for i < listLen {
		var elem = l.Front()
		var first = elem.Value.(string)
		var index = strings.Index(multiChar, first)
		if index > -1 && index%2 == 0 && i+1 < listLen {
			l.Remove(elem)
			var next = l.Front()
			var second = next.Value.(string)
			if strings.Contains(multiChar, first+second) {
				l.Remove(next)
				l.PushBack(first + second)
				i++ // increment another time, one has been skipped
			} else {
				l.PushBack(first)
			}
		} else {
			l.MoveToBack(elem)
		}
		i++
	}
	return l
}

//updateQuote is the companion function of JoinQuote
func updateQuote(quote string, elem *list.Element, inQuote int, l *list.List) (string, int) {
	var quotePriority = strings.IndexAny("'\"", elem.Value.(string))
	var isQuote = quotePriority >= inQuote && quotePriority > -1
	if inQuote < 0 && !isQuote {
		l.MoveToBack(elem)
		return "", -1
	}
	quote += l.Remove(elem).(string)
	if inQuote > -1 && isQuote {
		l.PushBack(quote)
		return "", -1
	} else if isQuote {
		inQuote = quotePriority
	}
	return quote, inQuote
}

//JoinQuote join quote char and form a quoted string
func JoinQuote(l *list.List) *list.List {
	var quote string
	var inQuote = -1
	var listLen = l.Len()
	for i := 0; i < listLen; i++ {
		var elem = l.Front()
		var s = elem.Value.(string)
		if strings.Contains(s, "\\") {
			if inQuote > -1 {
				quote += l.Remove(elem).(string)
				quote += l.Remove(l.Front()).(string)
				i++
			} else {
				l.MoveToBack(elem)
			}
		} else {
			quote, inQuote = updateQuote(quote, elem, inQuote, l)
		}
	}
	return l
}
