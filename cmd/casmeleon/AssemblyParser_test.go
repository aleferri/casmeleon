package main

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/aleferri/casmeleon/internal/casm"
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

	for asmStream.Peek().ID() != text.EOF {
		asmErr := ParseSourceLine(lang, asmStream)
		if asmErr != nil {
			t.Error(asmErr.Error())
		}
	}
}