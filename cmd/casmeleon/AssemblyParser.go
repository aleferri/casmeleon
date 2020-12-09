package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aleferri/casmeleon/internal/casm"
	"github.com/aleferri/casmeleon/pkg/asm"
	"github.com/aleferri/casmeleon/pkg/parser"
	"github.com/aleferri/casmeleon/pkg/text"
)

func IsDirective(s string) bool {
	return s == ".advance" || s == ".org" || s == ".alias" || s == ".db"
}

func ParseDirective(lang casm.Language, stream parser.Stream, table *SymbolTable, prog *AssemblyProgram, name text.Symbol) error {
	switch name.Value() {
	case ".advance":
		{
			target, err := parser.Require(stream, text.Number)
			if err != nil {
				return casm.WrapMatchError(err, ".advance", "\n")
			}
			addr, convErr := lang.ParseUint(target.Value())
			if convErr != nil {
				return convErr
			}
			prog.Add(asm.MakeAdvance(uint32(addr)))
		}
	case ".org":
		{
			target, err := parser.Require(stream, text.Number)
			if err != nil {
				return casm.WrapMatchError(err, ".org", "\n")
			}
			addr, convErr := lang.ParseUint(target.Value())
			if convErr != nil {
				return convErr
			}
			prog.Add(asm.MakeOrg(uint32(addr)))
		}
	case ".alias":
		{
			_, err := parser.RequireSequence(stream, text.Identifier, text.Number)
			if err != nil {
				return casm.WrapMatchError(err, ".alias", "\n")
			}
			return errors.New("Unsupported yet, casmeleon v1 did not support it, anyway")
		}
	case ".db":
		{
			lastToken := name.WithID(text.Comma)
			rawVals := []text.Symbol{}
			for lastToken.ID() == text.Comma {
				stream.Next()
				seq, err := parser.RequireAny(stream, text.Identifier, text.Number, text.QuotedString)
				if err != nil {
					return casm.WrapMatchError(err, ".db", "\n")
				}
				rawVals = append(rawVals, seq)
				lastToken = stream.Peek()
			}
			values := []uint8{}
			for _, p := range rawVals {
				if p.ID() == text.QuotedString {
					str := p.Value()
					values = append(values, []byte(strings.TrimSuffix(strings.TrimPrefix(str, "\""), "\""))...)
				} else if p.ID() == text.Number {
					val, _ := lang.ParseUint(p.Value())
					values = append(values, uint8(val))
				} else if p.ID() == text.Identifier {
					values = append(values, 0)
					fmt.Println("Unsupported deposit of identifiers, using 0 as a value, proper support will be implemented later")
				}
			}
			prog.Add(asm.MakeDeposit(values))
		}
	}
	parser.Consume(stream, text.WHITESPACE)

	if stream.Peek().ID() != text.EOL {
		return fmt.Errorf("Expected End Of Line after the directive '%s', found instead '%s'", name.Value(), stream.Next().Value())
	}
	stream.Next()
	fmt.Println("Successfully parsed the directive ", name.Value())
	return nil
}

func ParseSourceLine(lang casm.Language, stream parser.Stream, table *SymbolTable, prog *AssemblyProgram) error {
	parser.Consume(stream, text.EOL)
	if stream.Peek().ID() == text.EOF {
		return nil
	}
	name, err := parser.Require(stream, text.Identifier)
	if err != nil {
		return casm.WrapMatchError(err, "\n", "\n")
	}

	if IsDirective(name.Value()) {
		return ParseDirective(lang, stream, table, prog, name)
	} else {
		if stream.Peek().ID() == text.Colon {
			stream.Next()
			fmt.Println("Found label: ", name.Value())
			return ParseSourceLine(lang, stream, table, prog)
		} else {
			lastToken := stream.Next()

			args := []text.Symbol{}
			tokensFormat := []text.Symbol{}

			for lastToken.ID() != text.EOL {
				tokensFormat = append(tokensFormat, lastToken)
				args = append(args, lastToken)
				lastToken = stream.Next()
			}

			win := lang.FilterOpcodesByName(name.Value())
			fmt.Println("Candidates number for ", name.Value(), ":", len(win.Candidates()))
			list := win.Candidates()
			for _, opc := range list {
				fmt.Println("Available format: ", opc.StringifyFormat(&lang))
			}
			types := []uint32{}
			format := []uint32{}
			numSet, _ := lang.SetByName("Ints")
			for _, tok := range tokensFormat {
				if tok.ID() == text.Number {
					types = append(types, numSet.ID())
					format = append(format, text.Identifier)
				} else if tok.ID() == text.Identifier {
					setName, found := lang.SetOf(tok.Value())
					if found {
						types = append(types, setName.ID())
					} else {
						fmt.Printf("Probable Label found, require patchup later")
						table.Add(asm.MakeLabel(tok.Value(), nil))
						types = append(types, numSet.ID())
					}
					format = append(format, text.Identifier)
				} else {
					format = append(format, tok.ID())
				}
			}
			fmt.Println("Provided Format: ", casm.StringifyFormat(&lang, format, types))
			win = win.FilterByFormat(format, types)

			op, err := win.PickFirst()
			if err != nil {
				fmt.Println(types)
				return errors.New("Invalid opcode " + name.Value())
			}

			fmt.Println("Successfully got opcode " + op.Name())

			return nil
		}
	}
}
