package parser

//ExpectedOptions encapsulate alternative list of acceptable/expected tokens
type ExpectedOptions interface {
	StringFromMap(defs map[uint32]string) string
	StringFromArray(defs []string) string
}

//ExpectedKind is an implementation of ExpectedOptions for single possibilities
type ExpectedKind struct {
	symID uint32
}

func (e *ExpectedKind) StringFromMap(defs map[uint32]string) string {
	return defs[e.symID]
}

func (e *ExpectedKind) StringFromArray(defs []string) string {
	return defs[e.symID]
}

//MakeExpectedKind create a single kind list of expected/allowed tokens
func MakeExpectedKind(a uint32) ExpectedOptions {
	return &ExpectedKind{a}
}

//ExpectedAny is an implementation of ExpectedOptions
type ExpectedAny struct {
	symsID []uint32
}

func (e *ExpectedAny) StringFromMap(defs map[uint32]string) string {
	buf := ""
	for _, t := range e.symsID {
		buf += ", " + defs[t]
	}
	return buf
}

func (e *ExpectedAny) StringFromArray(defs []string) string {
	buf := ""
	for _, t := range e.symsID {
		buf += ", " + defs[t]
	}
	return buf
}

//MakeExpectedAny is an implementation of ExpectedOptions
func MakeExpectedAny(list ...uint32) ExpectedOptions {
	return &ExpectedAny{list}
}
