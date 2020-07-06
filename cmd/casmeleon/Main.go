package main

import (
	"bitbucket.org/mrpink95/casmeleon/internal/langdef"
	"bitbucket.org/mrpink95/casmeleon/internal/lexing"
	"bitbucket.org/mrpink95/casmeleon/internal/parsing"
	"bitbucket.org/mrpink95/casmeleon/internal/text"
	"bitbucket.org/mrpink95/casmeleon/internal/ui"
	"bufio"
	"flag"
	"os"
	"path/filepath"
	"strings"
)

//ASMDefaultOptions return the default options for assembly files
func ASMDefaultOptions(lang *langdef.LangDef) lexing.TokenMatchingOptions {
	const separators = "@#$,:{}()[]'\"\t "
	const operators = "+ - * / % ^ & | ~ > < ! <= >= == && || != << >> ->"
	const lineComment = ";"
	return lexing.NewMatchingOptions(operators, lang.RemoveSpecialPurposeChars(separators), operators+lineComment, lineComment)
}

//Settings are the global Casmeleon settings
type Settings struct {
	lang          *langdef.LangDef
	userInterface ui.UI
	passLimit     int
	debugMode     bool //print the output to video in extended form
}

//NewSettings build a new set of settings
func NewSettings(lang *langdef.LangDef, passLimit int, userInterface ui.UI) Settings {
	return Settings{lang, userInterface, passLimit, false}
}

func (p *Settings) SetDebugMode(b bool) {
	p.debugMode = b
}

var sourceLibrary = parsing.NewSourceLibrary()

func GetSyncOfFile(sourceName string) func([]*text.SourceLine) {
	return func(lines []*text.SourceLine) {
		sourceLibrary.AddSource(sourceName, true, lines)
	}
}

func ParseASM(fileName string, settings Settings) ([]uint, bool) {
	buffer := parsing.NewTokenBufferFromFile(fileName, ASMDefaultOptions(settings.lang))
	buffer.SyncLines(GetSyncOfFile(fileName))
	parent := NewASMSourceRoot(fileName, settings.lang)
	child := ASMSourceLeaf{parent, []OpcodeInfo{}, []ASMLabel{}}
	parent.AddChild(&child)
	pState := NewParserState(settings.userInterface, buffer)
	pState.SetCustomIdentification(func(t text.Token) (text.Token, ui.SourceCodeError) { return settings.lang.IdentifyNumber(t) })
	builder := NewASMReader(pState, &child)
	fullName, _ := filepath.Abs(fileName)
	includeList := make([]ASMSourceInclude, 0)
	noErr := true
	noMoreFiles := false
	for noErr && !noMoreFiles {
		noErr = builder.Parse(settings.lang, &includeList)
		noMoreFiles = true
		if len(includeList) > 0 {
			toParse := includeList[0]
			builder = toParse.SetupBuilder(builder, settings)
			includeList = includeList[1:]
			noMoreFiles = false
		}
	}
	builder.ReportUnresolvedSymbols()
	if !noErr || settings.userInterface.GetErrorCount() > 0 {
		return nil, false
	}
	return AssembleSource(fullName, settings, builder.pass, parent)
}

func AssembleSource(fullName string, settings Settings, firstPass PassResult, parent ASMSource) ([]uint, bool) {
	oldPass := firstPass
	suppressLoops := uint(1)
	suppressErrors := uint(1 << 2)
	for i := 0; i < settings.passLimit; i++ {
		pass := PassResult{oldPass.labels, []uint{}, []int{}}
		diff, err := pass.makePass(&oldPass, parent, suppressLoops|suppressErrors)
		if err != nil {
			err.Report(settings.userInterface, sourceLibrary.GetSource(fullName)[err.GetLine()-1])
			return nil, false
		}
		if diff == 0 {
			if settings.debugMode {
				pass.PrintResult(parent)
			}
			return pass.output, true
		}
		oldPass = pass
		suppressLoops = 0
		if i == 1 {
			suppressErrors = 0
		}
	}
	settings.userInterface.ReportError("Unstable output", true)
	return nil, false
}

func dumpOutput(originalFileName string, settings Settings, output []uint) {
	lastDot := strings.LastIndex(originalFileName, ".")
	fileNoExtension := originalFileName[0:lastDot]
	out, err := os.Create(fileNoExtension + ".bin")
	if err != nil {
		settings.userInterface.ReportError("Output to file failed: "+err.Error(), true)
		return
	}
	writer := bufio.NewWriter(out)
	for _, b := range output {
		val := uint8(b)
		writer.WriteByte(val)
	}
	writer.Flush()
	err = out.Close()
	if err != nil {
		settings.userInterface.ReportError(err.Error(), true)
	}
}

func main() {
	var langFileName string
	flag.StringVar(&langFileName, "lang", ".", "-lang=langfile")
	var debugMode bool
	flag.BoolVar(&debugMode, "debug", false, "-debug=true|false")
	var userPass int
	flag.IntVar(&userPass, "pass", 3, "-pass=x")
	flag.Parse()
	tUI := ui.NewConsoleUI(false, false)
	if strings.EqualFold(langFileName, ".") {
		tUI.ReportError("missing -lang=langfile", true)
		return
	}
	langFile, err := os.Open(langFileName)
	if err != nil {
		tUI.ReportError("failed open of file "+langFileName+", "+err.Error(), true)
		return
	}
	lang := langdef.NewLangDef(langFileName, bufio.NewReader(langFile), tUI)
	langFile.Close()
	if tUI.GetErrorCount() > 0 {
		return
	}
	lang.SortOpcodes()
	settings := NewSettings(lang, userPass, tUI)
	settings.SetDebugMode(debugMode)
	for _, f := range flag.Args() {
		if !strings.HasPrefix(f, "-") {
			output, success := ParseASM(f, settings)
			if success {
				dumpOutput(f, settings, output)
			}
		}
	}
}
