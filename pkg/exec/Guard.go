package exec

import "sort"

//Guard is the sensitivity list that guard a code path
type Guard struct {
	names []string
}

//Simplify a guard by removing a set of parameters
func (g Guard) Simplify(set []string) Guard {
	filtered := []string{}
	for _, c := range g.names {
		p := sort.SearchStrings(set, c)
		if p >= len(set) || set[p] != c {
			filtered = append(filtered, c)
		}
	}
	return BuildGuard(filtered)
}

//BuildGuard with a list of name
func BuildGuard(set []string) Guard {
	sort.Strings(set)
	return Guard{set}
}
