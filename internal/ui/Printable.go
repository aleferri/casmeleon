package ui

//Printable is an interface of a user printable object
type Printable interface {
	Print(ui UI)                            //Print the object using ui
	RuneIndex(ui UI, indexOfWord uint) uint //RuneIndex calculate the index of the first rune of the specified word
	LineNumber() uint                       //LineNumber return the lineNumber
	StringAt(i uint) string                 //StringAt return the i-th string
}
