package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/aleferri/casmeleon/internal/casm"
	"github.com/aleferri/casmeleon/internal/ui"
	"github.com/aleferri/casmeleon/pkg/asm"
	"github.com/aleferri/casmeleon/pkg/text"
)

func dumpOutput(originalFileName string, ui ui.UI, output []uint8) {
	lastDot := strings.LastIndex(originalFileName, ".")
	fileNoExtension := originalFileName[0:lastDot]
	out, err := os.Create(fileNoExtension + ".bin")
	if err != nil {
		ui.ReportError("Output to file failed: "+err.Error(), true)
		return
	}
	writer := bufio.NewWriter(out)
	for _, b := range output {
		writer.WriteByte(b)
	}
	writer.Flush()
	err = out.Close()
	if err != nil {
		ui.ReportError(err.Error(), true)
	}
}

func main() {
	var langFileName string
	flag.StringVar(&langFileName, "lang", ".", "-lang=langfile")
	var debugMode bool
	flag.BoolVar(&debugMode, "debug", false, "-debug=true|false")
	flag.Parse()
	tUI := ui.NewConsole(false, false)
	if strings.EqualFold(langFileName, ".") {
		tUI.ReportError("missing -lang=langfile", true)
		return
	}
	langFile, err := os.Open(langFileName)
	if err != nil {
		tUI.ReportError("failed open of file "+langFileName+", "+err.Error(), true)
		return
	}
	source := bufio.NewReader(langFile)
	repo := text.BuildSource(langFileName)

	stream := casm.BuildStream(source, &repo)

	root, err := casm.ParseCasm(stream, repo)

	if err != nil {
		parseErr, ok := err.(*casm.ParserError)
		if !ok {
			fmt.Println("Unexpected Error")
		} else {
			parseErr.PrettyPrint(&repo)
		}
		return
	}

	lang, semErr := casm.MakeLanguage(root)
	if semErr != nil {
		fmt.Println("Error " + semErr.Error())
		return
	}

	langFile.Close()
	if tUI.GetErrorCount() > 0 {
		return
	}

	for _, f := range flag.Args() {
		if !strings.HasPrefix(f, "-") {
			programFileName := "../../tests/example_test.s"

			var programfile, programErr = os.Open(programFileName)
			if programErr != nil {
				wnd, _ := os.Getwd()
				fmt.Printf("Error during opening of file %s from %s\n", langFileName, wnd)
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
					return
				}
			}

			if len(asmSymbolTable.watchList) > 0 {
				fmt.Printf("Missing %d symbols:\n", len(asmSymbolTable.watchList))
				for _, miss := range asmSymbolTable.watchList {
					fmt.Printf("Missing symbol %s\n", miss.Value())
				}
				return
			}

			ctx := asm.MakeSourceContext()
			binaryImage, compilingErr := asm.AssembleSource(asmProgram.list, ctx)
			if compilingErr != nil {
				fmt.Println(compilingErr.Error())
				return
			}

			dumpOutput(f, tUI, binaryImage)
		}
	}
}
