package main

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/aleferri/casmeleon/internal/casm"
	"github.com/aleferri/casmeleon/pkg/asm"
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

	stream := casm.BuildStream(source, &repo)

	root, err := casm.ParseCasm(stream, repo)

	if err != nil {
		parseErr, ok := err.(*casm.ParserError)
		if !ok {
			fmt.Println("Unexpected Error")
		} else {
			parseErr.PrettyPrint(&repo)
		}

		t.Fail()
	}

	lang, semErr := casm.MakeLanguage(root)
	if semErr != nil {
		fmt.Println("Error " + semErr.Error())
		t.Fail()
	}

	programFileName := "../../tests/example_test.s"

	var programfile, programErr = os.Open(programFileName)
	if programErr != nil {
		wnd, _ := os.Getwd()
		t.Errorf("Error during opening of file %s from %s\n", fileName, wnd)
		return
	}

	program := bufio.NewReader(programfile)

	asmSource := text.BuildSource("example_test.s")

	asmStream := MakeRootStream(program, &asmSource)

	asmProgram := MakeAssemblyProgram()
	asmSymbolTable := MakeSymbolTable()

	for asmStream.Peek().ID() != text.EOF {
		asmErr := ParseSourceLine(lang, asmStream, &asmSymbolTable, &asmProgram)
		if asmErr != nil {
			parseErr, ok := asmErr.(*casm.ParserError)
			if !ok {
				fmt.Println(asmErr.Error())
			} else {
				parseErr.PrettyPrint(&asmSource)
			}
			t.Fail()
			break
		}
	}

	if len(asmSymbolTable.watchList) > 0 {
		t.Errorf("Missing %d symbols:\n", len(asmSymbolTable.watchList))
		for _, miss := range asmSymbolTable.watchList {
			t.Errorf("Missing symbol %s\n", miss.Value())
		}
	}

	ctx := asm.MakeSourceContext()
	_, compilingErr := asm.AssembleSource(asmProgram.list, ctx)
	if compilingErr != nil {
		t.Error(compilingErr.Error())
	}
}
