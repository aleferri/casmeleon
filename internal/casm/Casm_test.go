package casm

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"testing"
	"unicode/utf8"

	"github.com/aleferri/casmeleon/pkg/parser"
	"github.com/aleferri/casmeleon/pkg/text"
)

func TestCasmProcessing(t *testing.T) {
	fileName := "../../tests/parser_test.casm"
	var file, fileErr = os.Open(fileName)
	if fileErr != nil {
		wnd, _ := os.Getwd()
		t.Errorf("Error during opening of file %s from %s\n", fileName, wnd)
		return
	}
	source := bufio.NewReader(file)
	repo := text.BuildSource("parser_test.casm")

	stream := BuildStream(source, &repo)

	root, err := ParseCasm(stream, repo)

	if err != nil {
		switch err.(type) {
		case *parser.MatchError:
			{
				m := err.(*parser.MatchError)
				v := idDescriptor[m.Expected()]
				fileName, lineIndex, column := repo.FindPosition(m.Found())
				fmt.Printf("In file %s, error at %d, %d: ", fileName, lineIndex+1, column+1)
				fmt.Printf(m.Error(), m.Found().Value(), v)
				fmt.Println()
			}
		case *ParserError:
			{
				p := err.(*ParserError)
				m := p.internal
				v := idDescriptor[m.Expected()]
				fileName, lineIndex, column := repo.FindPosition(m.Found())
				fmt.Printf("In file %s, error at %d, %d: ", fileName, lineIndex+1, column+1)
				fmt.Printf(m.Error(), m.Found().Value(), v)
				fmt.Println()
				context := repo.SliceScope(m.Found(), p.context)
				for _, c := range context {
					fmt.Print(c.Value())
				}
				fmt.Println()
				line := repo.SliceLine(m.Found())
				for _, c := range line[1:] {
					char := ""
					if c.Equals(m.Found()) {
						char = "^"
					} else if c.ID() > 2 {
						char = "-"
					}
					runeCount := utf8.RuneCountInString(c.Value())
					for k := 0; k < runeCount; k++ {
						fmt.Print(char)
					}
				}
				fmt.Println()
			}
		default:
			{
				fmt.Println("Cannot identify source")
				fmt.Println(err.Error())
				context := stream.Source().SliceScope(stream.Peek(), "{")
				for _, t := range context {
					fmt.Print(t.Value())
				}
				fmt.Println()
			}
		}
		t.Fail()
	}

	lang, semErr := MakeLanguage(root)
	if semErr != nil {
		fmt.Println("Error " + semErr.Error())
		t.Fail()
	}

	for _, in := range lang.inlines {
		fmt.Println("Inline name: " + in.name)
		for _, exec := range in.runList {
			fmt.Println(reflect.TypeOf(exec))
		}
	}

	for _, in := range lang.opcodes {
		fmt.Println("Opcode name: " + in.name)
		for _, exec := range in.runList {
			fmt.Println(reflect.TypeOf(exec))
		}
	}
	t.Fail()
}
