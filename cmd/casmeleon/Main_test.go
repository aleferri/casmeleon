package main

import (
	"github.com/aleferri/casmeleon/internal/langdef"
	"github.com/aleferri/casmeleon/internal/ui"
	"bufio"
	"fmt"
	"os"
	"testing"
)

func TestParseASM(t *testing.T) {
	textUI := ui.NewConsoleUI(false, false)
	langFileName := "../testing/vs_cpu.casm"
	langFile, err := os.Open(langFileName)
	if err != nil {
		textUI.ReportError("failed open of file "+langFileName+", "+err.Error(), true)
		return
	}
	lang := langdef.NewLangDef("vs_cpu.casm", bufio.NewReader(langFile), textUI)
	if textUI.GetErrorCount() > 0 {
		t.Fail()
		return
	}
	settings := NewSettings(lang, 3, textUI)
	output, success := ParseASM("../testing/test.txt", settings)
	if len(output) == 0 || !success {
		t.Fail()
	}
	fmt.Printf("%v\n", output)
}

func TestParseASMFail(t *testing.T) {
	textUI := ui.NewConsoleUI(false, false)
	langFileName := "../testing/microcode.casm"
	langFile, err := os.Open(langFileName)
	if err != nil {
		textUI.ReportError("failed open of file "+langFileName+", "+err.Error(), true)
		return
	}
	_ = langdef.NewLangDef("microcode.casm", bufio.NewReader(langFile), textUI)
	if textUI.GetErrorCount() < 1 {
		t.Fail()
		return
	}
}

func TestInclude(t *testing.T) {
	textUI := ui.NewConsoleUI(false, false)
	langFileName := "../testing/vs_cpu.casm"
	langFile, err := os.Open(langFileName)
	if err != nil {
		textUI.ReportError("failed open of file "+langFileName+", "+err.Error(), true)
		return
	}
	lang := langdef.NewLangDef("vs_cpu.casm", bufio.NewReader(langFile), textUI)
	if textUI.GetErrorCount() > 0 {
		t.Fail()
		return
	}
	settings := NewSettings(lang, 3, textUI)
	output, success := ParseASM("../testing/test_include.txt", settings)
	if len(output) == 0 || !success {
		t.Fail()
	}
}

func TestIncludeFail(t *testing.T) {
	textUI := ui.NewConsoleUI(false, false)
	langFileName := "../testing/vs_cpu.casm"
	langFile, err := os.Open(langFileName)
	if err != nil {
		textUI.ReportError("failed open of file "+langFileName+", "+err.Error(), true)
		return
	}
	lang := langdef.NewLangDef("vs_cpu.casm", bufio.NewReader(langFile), textUI)
	if textUI.GetErrorCount() > 0 {
		t.Fail()
		return
	}
	settings := NewSettings(lang, 3, textUI)
	ParseASM("../testing/test_include_fail.txt", settings)
	if textUI.GetErrorCount() < 1 {
		t.Fail()
		return
	}
}

func BenchmarkParseASM(b *testing.B) {
	textUI := ui.NewConsoleUI(false, false)
	langFileName := "../testing/vs_cpu.casm"
	langFile, err := os.Open(langFileName)
	if err != nil {
		textUI.ReportError("failed open of file "+langFileName+", "+err.Error(), true)
		return
	}
	lang := langdef.NewLangDef("vs_cpu", bufio.NewReader(langFile), textUI)
	if textUI.GetErrorCount() > 0 {
		b.Fail()
		return
	}
	settings := NewSettings(lang, 3, textUI)
	settings.SetDebugMode(true)
	for i := 0; i < b.N; i++ {
		output, success := ParseASM("../testing/test.txt", settings)
		if len(output) == 0 || !success {
			b.Fail()
		}
	}
}
