package langdef

import (
	"strings"
	"testing"

	"github.com/aleferri/casmeleon/internal/text"
	"github.com/aleferri/casmeleon/internal/ui"
)

func TestLangDef_AssembleOpcode(t *testing.T) {
	textUI := ui.NewConsole(false, false)
	lang := NewLangDef("test_assemble", strings.NewReader(".opcode inx -> { if this_address > 4096 { .error this_address \"too long\"; } else { .db 0xFF; } }"), textUI)
	output, err := lang.AssembleOpcode(0, []text.Token{text.NewInternalToken("this_address")}, []int{0}, 0)
	if err != nil {
		println(err.Error(), err.GetLine())
		t.Fail()
		return
	}
	if len(output) != 1 {
		t.Errorf("Assembled opcode too fat, expected 1 found: %d\n", len(output))
		t.Fail()
		return
	}
	if output[0] != 0xFF {
		t.Errorf("Expected 0xFF, found %d\n", output[0])
		t.Fail()
		return
	}
}

const vsCPULangDef = ".enum regs\n{\nA\n}\n.enum ports\n{\nport_0,port_1}" +
	".opcode jmp mem -> {\n " +
	"if mem > 63 { .error mem \"address out of range, max: 63\"; }" +
	".db 2 << 6 + mem;\n" +
	"}" +
	".opcode add mem -> {\n" +
	".db mem;\n" +
	"}" +
	".opcode and mem -> {\n" +
	".db 1 << 6 + mem;\n" +
	"}\n" +
	".opcode sta mem -> {\n" +
	".db 3 << 6 + mem;\n" +
	"}" +
	".opcode ldp port -> {\n if port .in ports { .db port; } else { .error port \"invalid port name\" ;}}"

const vsCPUOpcodeNumber = 5

func TestLangDef_GetEnumValue(t *testing.T) {
	textUI := ui.NewConsole(false, false)
	lang := NewLangDef("test_enum_value", strings.NewReader(vsCPULangDef), textUI)
	if lang.GetEnumValue("port_1") != 1 {
		t.Errorf("Cannot find value port_1 %d\n", lang.GetEnumValue("port_1"))
		t.Fail()
	}
	if lang.GetEnumValue("A") != 0 {
		t.Errorf("Cannot find value A %d\n", lang.GetEnumValue("A"))
		t.Fail()
	}
}

func TestLangDef_GetEnumValueFail(t *testing.T) {
	textUI := ui.NewConsole(false, false)
	lang := NewLangDef("test_enum_value", strings.NewReader(vsCPULangDef), textUI)
	if lang.GetEnumValue("GHT") != -1 {
		t.Errorf("GHT does not exist, should be -1, %d\n", lang.GetEnumValue("GHT"))
		t.Fail()
	}
	if lang.GetEnumValue("port_4") != -1 {
		t.Errorf("port_4 does not exist, should be -1, %d\n", lang.GetEnumValue("port_4"))
		t.Fail()
	}
}

func TestLangDef_SortOpcodes(t *testing.T) {
	textUI := ui.NewConsole(false, false)
	lang := NewLangDef("test_sort_opcode", strings.NewReader(vsCPULangDef), textUI)
	if vsCPUOpcodeNumber != lang.GetOpcodeNumber() {
		t.Errorf("Different in length: found %d, expected 2\n", lang.GetOpcodeNumber())
		t.Fail()
		return
	}
	lang.SortOpcodes()
	if !strings.EqualFold(lang.opcodes[0].name, "add") || !strings.EqualFold(lang.opcodes[1].name, "and") {
		t.Errorf("Invalid sequence expected add, and found %s %s\n", lang.opcodes[0].name, lang.opcodes[1].name)
		t.Fail()
		return
	}
}

func TestFilterWindow_FilterByName(t *testing.T) {
	textUI := ui.NewConsole(false, false)
	lang := NewLangDef("test_filter", strings.NewReader(vsCPULangDef), textUI)
	if vsCPUOpcodeNumber != lang.GetOpcodeNumber() {
		t.Errorf("Different in length: found %d, expected 2\n", lang.GetOpcodeNumber())
		t.Fail()
		return
	}
	lang.SortOpcodes()
	window := lang.NewFilterWindow()
	byName := window.FilterByName("add")
	byNameStruct := byName.(DefaultFilterWindow)
	if len(byNameStruct.opcodes) != 1 {
		t.Errorf("Expected 1 opcode, found: %d\n", len(byNameStruct.opcodes))
		t.Fail()
	}
}

func TestFilterWindow_FilterByToken(t *testing.T) {
	textUI := ui.NewConsole(false, false)
	lang := NewLangDef("test_filter", strings.NewReader(vsCPULangDef), textUI)
	if vsCPUOpcodeNumber != lang.GetOpcodeNumber() {
		t.Errorf("Different in length: found %d, expected 2\n", lang.GetOpcodeNumber())
		t.Fail()
		return
	}
	lang.SortOpcodes()
	window := lang.NewFilterWindow()
	byName := window.FilterByName("add")
	byNameStruct := byName.(DefaultFilterWindow)
	if len(byNameStruct.opcodes) != 1 {
		t.Errorf("Expected 1 opcode, found: %d\n", len(byNameStruct.opcodes))
		t.Fail()
	}
	byToken, keep := byName.FilterByToken(text.NewInternalToken("mem"))
	byTokenStruct := byToken.(DefaultFilterWindow)
	if !keep {
		t.Error("Sent an identifier, expected true")
		t.Fail()
	}
	if len(byTokenStruct.opcodes) != 1 {
		t.Errorf("Expected 1 opcode, found: %d\n", len(byTokenStruct.opcodes))
		t.Fail()
	}
}

func TestFilterWindow_Harvest(t *testing.T) {
	textUI := ui.NewConsole(false, false)
	lang := NewLangDef("test_filter", strings.NewReader(vsCPULangDef), textUI)
	if vsCPUOpcodeNumber != lang.GetOpcodeNumber() {
		t.Errorf("Different in length: found %d, expected 2\n", lang.GetOpcodeNumber())
		t.Fail()
		return
	}
	lang.SortOpcodes()
	window := lang.NewFilterWindow()
	byName := window.FilterByName("jmp")
	byToken, _ := byName.FilterByToken(text.NewInternalToken("mem"))
	found, index := byToken.Collect()
	if !found {
		t.Error("Not found, expected found\n")
		t.Fail()
	}
	if index != 2 {
		t.Errorf("Expected 2, found %d\n", index)
	}
}
