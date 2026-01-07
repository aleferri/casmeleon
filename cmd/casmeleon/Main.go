package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aleferri/casmeleon/internal/casm"
	"github.com/aleferri/casmeleon/internal/ui"
	"github.com/aleferri/casmeleon/pkg/asm"
	"github.com/aleferri/casmeleon/pkg/parser"
	"github.com/aleferri/casmeleon/pkg/text"
	"github.com/aleferri/casmvm/pkg/vmex"
	"github.com/aleferri/casmvm/pkg/vmio"
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

func exportOutput(originalFileName string, ui ui.UI, output []uint8) {
	lastDot := strings.LastIndex(originalFileName, ".")
	fileNoExtension := originalFileName[0:lastDot]
	out, err := os.Create(fileNoExtension + ".export.txt")
	if err != nil {
		ui.ReportError("Output to file failed: "+err.Error(), true)
		return
	}

	writer := bufio.NewWriter(out)

	if len(output) > 1 {
		for _, v := range output[0 : len(output)-2] {
			bin := strconv.FormatUint(uint64(v), 2)
			writer.WriteString("8'b" + bin)
			writer.WriteRune('\n')
		}
	}

	if len(output) > 0 {
		bin := strconv.FormatUint(uint64(output[len(output)-1]), 2)
		writer.WriteString("8'b" + bin)
		writer.WriteRune('\n')
	}

	writer.Flush()
}

func ParseIncludedASMFile(lang casm.Language, program *AssemblyProgram, symTable *SymbolTable, sourceFile string) error {
	var programfile, programErr = os.Open(sourceFile)
	if programErr != nil {
		wnd, _ := os.Getwd()
		return fmt.Errorf("error during opening of file %s from %s", sourceFile, wnd)
	}

	code := text.BuildSource(sourceFile)
	programCode := bufio.NewReader(programfile)

	stream := MakeRootStream(programCode, &code)

	parser.ConsumeAll(stream, text.EOL)

	for stream.Peek().ID() != text.EOF {
		if stream.Peek().Value() == ".include" {
			stream.Next()
			toInclude, noFile := parser.Require(stream, text.QuotedString)
			if noFile != nil {
				return noFile
			}
			includedFileName := toInclude.Value()
			includedErr := ParseIncludedASMFile(lang, program, symTable, filepath.Dir(sourceFile)+"/"+includedFileName[1:len(includedFileName)-1])
			if includedErr != nil {
				return includedErr
			}
		}
		parseErr := ParseSourceLine(lang, stream, symTable, program)
		if parseErr != nil {
			parseErr, ok := parseErr.(*casm.ParserError)
			if !ok {
				fmt.Println(parseErr.Error())
			} else {
				parseErr.PrettyPrint(&code)
			}
			return errors.New("error during compilation")
		}
		parser.ConsumeAll(stream, text.EOL)
	}
	programfile.Close()
	return nil
}

func ParseASMFile(lang casm.Language, sourceFile string) (*AssemblyProgram, error) {
	var programfile, programErr = os.Open(sourceFile)
	if programErr != nil {
		wnd, _ := os.Getwd()
		return nil, fmt.Errorf("error during opening of file %s from %s", sourceFile, wnd)
	}

	code := text.BuildSource(sourceFile)
	programCode := bufio.NewReader(programfile)

	stream := MakeRootStream(programCode, &code)

	program := MakeAssemblyProgram()
	symTable := MakeSymbolTable()
	parser.ConsumeAll(stream, text.EOL)

	for stream.Peek().ID() != text.EOF {
		if stream.Peek().Value() == ".include" {
			stream.Next()
			toInclude, noFile := parser.Require(stream, text.QuotedString)
			if noFile != nil {
				return nil, noFile
			}
			includedFileName := toInclude.Value()
			includedErr := ParseIncludedASMFile(lang, &program, &symTable, filepath.Dir(sourceFile)+"/"+includedFileName[1:len(includedFileName)-1])
			if includedErr != nil {
				return nil, includedErr
			}
		}
		parseErr := ParseSourceLine(lang, stream, &symTable, &program)
		if parseErr != nil {
			parseErr, ok := parseErr.(*casm.ParserError)
			if !ok {
				fmt.Println(parseErr.Error())
			} else {
				parseErr.PrettyPrint(&code)
			}
			return nil, errors.New("error during compilation")
		}

		parser.ConsumeAll(stream, text.EOL)
	}

	if len(symTable.watchList) > 0 {
		for _, miss := range symTable.watchList {
			fmt.Printf("Missing symbol %s\n", miss.Value())
		}
		return nil, fmt.Errorf("missing %d symbols", len(symTable.watchList))
	}
	programfile.Close()
	return &program, nil
}

func ExportTraces(lang *casm.Language, list []asm.Compilable) {
	file, err := os.Create("dump.trace")
	if err != nil {
		return
	}

	for _, entry := range lang.Executables() {
		fmt.Fprintln(file, "fn", entry.Name())
		for _, op := range entry.Listing() {
			fmt.Fprintln(file, op.String())
		}
		fmt.Fprintln(file)
	}

	fmt.Fprintln(file, "fn file_listing")
	for _, call := range list {
		instance, ok := call.(*OpcodeInstance)
		if ok {
			args := []int64{}
			for _, a := range instance.parameters {
				args = append(args, a.Value())
			}

			args = append(args, 0xFFFFFFFF)
			fmt.Fprintln(file, "fn", instance.name, instance.line)
			fmt.Fprintln(file, "invoke", instance.invokeTarget, args)
		}
	}
}

func main() {
	var langFileName string
	var debugMode bool
	var exportAssembly string
	var dumpTrace bool
	var byteSize uint

	flag.StringVar(&langFileName, "lang", ".", "-lang=langfile")
	flag.BoolVar(&debugMode, "debug", false, "-debug=true|false")
	flag.BoolVar(&dumpTrace, "trace", false, "-trace=true|false")
	flag.StringVar(&exportAssembly, "export", "none", "-export=bin|hex")
	flag.UintVar(&byteSize, "byteSize", 8, "-byteSize=8|16|32")
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

	lang, semErr := casm.MakeLanguage(root, uint32(byteSize))
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

			if debugMode {
				fmt.Println("Parsing " + f)
			}

			program, errAsm := ParseASMFile(lang, f)

			if debugMode {
				fmt.Println("End parsing")
			}

			if errAsm != nil {
				fmt.Println("Error: " + errAsm.Error())
				break
			}

			if dumpTrace {
				ExportTraces(&lang, program.list)
			}

			log := vmio.MakeVMLoggerConsole(vmio.ALL)
			ex := vmex.MakeNaiveVM(lang.Executables(), log, vmex.MakeVMFrame())
			ctx := asm.MakeSourceContext(uint32(byteSize))
			binaryImage, compilingErr := asm.AssembleSource(ex, program.list, ctx)

			if compilingErr != nil {
				fmt.Println(compilingErr.Error())
				return
			}

			dumpOutput(f, tUI, binaryImage)

			if exportAssembly == "bin" {
				exportOutput(f, tUI, binaryImage)
			}
		}
	}
}
