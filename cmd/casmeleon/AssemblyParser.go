package main

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/aleferri/casmeleon/internal/casm"
	"github.com/aleferri/casmeleon/pkg/asm"
	"github.com/aleferri/casmeleon/pkg/parser"
	"github.com/aleferri/casmeleon/pkg/text"
)

func IsDirective(s string) bool {
	return s == ".advance" || s == ".org" || s == ".alias" || s == ".db" || s == ".dw"
}

func ParseDirective(lang casm.Language, stream parser.Stream, table *SymbolTable, prog *AssemblyProgram, directive text.Symbol) error {
	switch directive.Value() {
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
			return errors.New("unsupported yet, casmeleon v1 did not support it, anyway")
		}
	case ".db":
		{
			rawVals := []text.Symbol{}

			sym, err := parser.RequireAny(stream, text.Identifier, text.Number, text.QuotedString)
			rawVals = append(rawVals, sym)
			if err != nil {
				return casm.WrapMatchError(err, ".db", "\n")
			}
			for stream.Peek().ID() == text.Comma {
				stream.Next()
				toSum := ""
				if stream.Peek().ID() == text.OperatorMinus {
					toSum = stream.Next().Value()
				}
				sym, err = parser.RequireAny(stream, text.Identifier, text.Number, text.QuotedString)
				if err != nil {
					return casm.WrapMatchError(err, ".db", "\n")
				}
				rawVals = append(rawVals, sym.WithText(toSum+sym.Value()))
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
	case ".dw":
		{
			rawVals := []text.Symbol{}

			sym, err := parser.RequireAny(stream, text.Identifier, text.Number, text.QuotedString)
			rawVals = append(rawVals, sym)
			if err != nil {
				return casm.WrapMatchError(err, ".dw", "\n")
			}
			for stream.Peek().ID() == text.Comma {
				stream.Next()
				toSum := ""
				if stream.Peek().ID() == text.OperatorMinus {
					toSum = stream.Next().Value()
				}
				sym, err = parser.RequireAny(stream, text.Identifier, text.Number, text.QuotedString)
				if err != nil {
					return casm.WrapMatchError(err, ".dw", "\n")
				}
				rawVals = append(rawVals, sym.WithText(toSum+sym.Value()))
			}
			values := []uint8{}
			for _, p := range rawVals {
				if p.ID() == text.QuotedString {
					str := p.Value()
					trimmed := strings.TrimSuffix(strings.TrimPrefix(str, "\""), "\"")
					for _, c := range bytes.Runes([]byte(trimmed)) {
						if lang.IsLittleEndian() {
							values = append(values, uint8(c&255))
							values = append(values, uint8(c>>8))
						} else {
							values = append(values, uint8(c>>8))
							values = append(values, uint8(c&255))
						}

					}
				} else if p.ID() == text.Number {
					val, _ := lang.ParseUint(p.Value())
					if lang.IsLittleEndian() {
						values = append(values, uint8(val&255))
						values = append(values, uint8(val>>8))
					} else {
						values = append(values, uint8(val>>8))
						values = append(values, uint8(val&255))
					}

				} else if p.ID() == text.Identifier {
					values = append(values, 0)
					values = append(values, 0)
					fmt.Println("Unsupported deposit of identifiers, using 0 as a value, proper support will be implemented later")
				}
			}
			prog.Add(asm.MakeDeposit(values))
		}
	}
	parser.Consume(stream, text.WHITESPACE)

	if stream.Peek().ID() != text.EOL {
		return fmt.Errorf("expected End Of Line after the directive '%s', found instead '%s'", directive.Value(), stream.Next().Value())
	}
	stream.Next()
	return nil
}

func ParseLabel(lang casm.Language, stream parser.Stream, table *SymbolTable, prog *AssemblyProgram, labelToken text.Symbol) error {
	labelName := labelToken.Value()
	fqln := labelName
	isLocalLabel := labelName[0] == '.'
	if isLocalLabel {
		if table.lastGlobalLabel == nil {
			matchErr := parser.ExpectedAnyOf(labelToken, "Unexpected a local label %s: expected global label '%s'", text.Identifier)
			parseErr := casm.WrapMatchError(matchErr, "\n", "\n")
			return parseErr
		}
		fqln = table.lastGlobalLabel.Name() + labelName
	}
	label := asm.MakeLabel(fqln, nil, lang.ByteSize())
	if !isLocalLabel {
		table.lastGlobalLabel = label
	}
	table.Add(label)
	table.UnWatch(label.Name())
	prog.Add(label)

	return ParseSourceLine(lang, stream, table, prog)
}

func TokensToFormat(lang casm.Language, symTable *SymbolTable, tokens []text.Symbol) (ArgumentFormat, error) {
	args := MakeFormat()
	numSet, _ := lang.SetByName("Ints")
	for _, tok := range tokens {
		if tok.ID() == text.Number {
			args.types = append(args.types, numSet.ID())
			args.format = append(args.format, text.Identifier)
			numVal, err := lang.ParseInt(tok.Value())
			if err != nil {
				matchErr := parser.ExpectedSymbol(tok, "Unexpected '%s' found, expecting a valid %s", text.Number)
				return args, casm.WrapMatchError(matchErr, "\n", "\n")
			}
			args.parameters = append(args.parameters, asm.MakeConstant(numVal))
		} else if tok.ID() == text.Identifier {
			setName, found := lang.SetOf(tok.Value())
			if found && setName.ID() > 1 {
				args.types = append(args.types, setName.ID())
				setValue, _ := setName.Value(tok.Value())
				args.parameters = append(args.parameters, asm.MakeConstant(int64(setValue)))
			} else {
				name := tok.Value()
				if name[0] == '.' {
					tok = tok.WithText(symTable.lastGlobalLabel.Name() + name)
				}
				lookup, found := symTable.Search(tok.Value())
				if !found {
					lookup = MakePatchSymbol(tok.Value(), symTable)
					symTable.Watch(tok)
				}
				args.parameters = append(args.parameters, lookup)
				args.types = append(args.types, numSet.ID())
			}
			args.format = append(args.format, text.Identifier)
		} else {
			args.format = append(args.format, tok.ID())
		}
	}
	return args, nil
}

func ParseSourceLine(lang casm.Language, stream parser.Stream, table *SymbolTable, prog *AssemblyProgram) error {
	parser.ConsumeAll(stream, text.EOL)
	if stream.Peek().ID() == text.EOF {
		return nil
	}
	name, err := parser.Require(stream, text.Identifier)
	if err != nil {
		return casm.WrapMatchError(err, "\n", "\n")
	}

	if IsDirective(name.Value()) {
		return ParseDirective(lang, stream, table, prog, name)
	} else if stream.Peek().ID() == text.Colon {
		stream.Next()
		return ParseLabel(lang, stream, table, prog, name)
	} else {
		lastToken := stream.Next()

		tokensFormat := []text.Symbol{}

		for lastToken.ID() != text.EOL {
			tokensFormat = append(tokensFormat, lastToken)
			lastToken = stream.Next()
			if lastToken.ID() == text.OperatorMinus {
				lastToken = stream.Next()
				lastToken = lastToken.WithText("-" + lastToken.Value())
			}
		}

		win := lang.FilterOpcodesByName(name.Value())

		args, literalErrs := TokensToFormat(lang, table, tokensFormat)

		if literalErrs != nil {
			return literalErrs
		}

		win = win.FilterByFormat(args.format, args.types)

		op, err := win.PickFirst()
		if err != nil {
			matchErr := parser.ExpectedAnyOf(name, "Expected valid opcode, but %s was found, unrecognized %s", text.Identifier)
			return casm.WrapMatchError(matchErr, name.Value(), "\n")
		}

		prog.Add(MakeOpcodeInstance(op, args, table, 1))

		return nil
	}
}
