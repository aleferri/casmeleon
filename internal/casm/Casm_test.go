package casm

import (
	"bufio"
	"fmt"
	"os"
	"testing"

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
		parseErr, ok := err.(*ParserError)
		if !ok {
			fmt.Println("Unexpected Error")
		} else {
			parseErr.PrettyPrint(&repo)
		}
		t.Fail()
	}

	lang, semErr := MakeLanguage(root)
	if semErr != nil {
		fmt.Println("Error " + semErr.Error())
		t.Fail()
	}

	for i, fnName := range lang.fnNames {
		fmt.Println("Func name: " + fnName)
		lang.fnList[i].Dump()
	}

	for _, in := range lang.opcodes {
		fmt.Println("Opcode name: " + in.name)
		for _, exec := range in.runList {
			fmt.Println(exec.String())
		}
	}
}
