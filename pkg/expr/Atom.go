package expr

import "math"

//Atom struct contains necessary requirements for structure handling
type Atom struct {
	raw   string //raw atom
	val   int64  //value of the atom, if it has one, default is 0
	tag   int16  //tag of the atom (0: literal, 1: parameter, 2: set member)
	local uint16 //local index of this atom, default is MAX_UINT16
}

func (a Atom) Tag() int16 {
	return a.tag
}

func (a Atom) Local() uint16 {
	return a.local
}

func (a Atom) Value() int64 {
	return a.val
}

func MakeAtom(raw string, val int64, tag int16) Atom {
	return Atom{raw, val, tag, math.MaxUint16}
}

func MakeLiteral(raw string, val int64) Atom {
	return Atom{raw, val, 0, math.MaxUint16}
}

func MakeParameter(raw string, val int64, local uint16) Atom {
	return Atom{raw, val, 1, local}
}

func MakeMember(raw string, val int64) Atom {
	return Atom{raw, val, 2, math.MaxUint16}
}

func MakeLocal(raw string, val int64, local uint16) Atom {
	return Atom{raw, val, 3, local}
}
