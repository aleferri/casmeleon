package main

import (
	"bytes"
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

//ParseDepositValues consumes the comma separated value list of a .db or a .dw.
//Numbers and quoted strings become constants, identifiers become symbols
//resolved at assembly time, so that a label can be deposited.
func ParseDepositValues(lang casm.Language, stream parser.Stream, table *SymbolTable, directive string) ([]asm.Symbol, error) {
	values := []asm.Symbol{}

	first := true
	for first || stream.Peek().ID() == text.Comma {
		if !first {
			stream.Next()
		}
		first = false

		negate := ""
		if stream.Peek().ID() == text.OperatorMinus {
			negate = stream.Next().Value()
		}
		tok, err := parser.RequireAny(stream, text.Identifier, text.Number, text.QuotedString)
		if err != nil {
			return values, casm.WrapMatchError(err, directive, "\n")
		}

		switch tok.ID() {
		case text.QuotedString:
			{
				str := strings.TrimSuffix(strings.TrimPrefix(tok.Value(), "\""), "\"")
				for _, c := range bytes.Runes([]byte(str)) {
					values = append(values, asm.MakeConstant(int64(c)))
				}
			}
		case text.Number:
			{
				val, convErr := lang.ParseInt(negate + tok.Value())
				if convErr != nil {
					return values, convErr
				}
				values = append(values, asm.MakeConstant(val))
			}
		default:
			{
				if negate != "" {
					return values, fmt.Errorf("cannot negate the symbol '%s' in %s", tok.Value(), directive)
				}
				name := tok.Value()
				if name[0] == '.' {
					if table.lastGlobalLabel == nil {
						return values, fmt.Errorf("local label '%s' in %s without a global label before it", name, directive)
					}
					tok = tok.WithText(table.lastGlobalLabel.Name() + name)
				}
				lookup, found := table.Search(tok.Value())
				if !found {
					lookup = MakePatchSymbol(tok.Value(), table)
					table.Watch(tok)
				}
				values = append(values, lookup)
			}
		}
	}
	return values, nil
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
			syms, err := parser.RequireSequence(stream, text.Identifier, text.Number)
			if err != nil {
				return casm.WrapMatchError(err, ".alias", "\n")
			}
			name := syms[0].Value()
			if _, exists := table.Search(name); exists {
				return fmt.Errorf("symbol '%s' is already defined", name)
			}
			val, convErr := lang.ParseInt(syms[1].Value())
			if convErr != nil {
				return convErr
			}
			//a named constant is a static symbol: it never moves, so it does not
			//participate in the address fixed point at all
			table.Add(MakeNamedConstant(name, val))
			table.UnWatch(name)
		}
	case ".db":
		{
			values, err := ParseDepositValues(lang, stream, table, ".db")
			if err != nil {
				return err
			}
			prog.Add(asm.MakeDepositSymbols(values, 1, !lang.IsLittleEndian()))
		}
	case ".dw":
		{
			values, err := ParseDepositValues(lang, stream, table, ".dw")
			if err != nil {
				return err
			}
			prog.Add(asm.MakeDepositSymbols(values, 2, !lang.IsLittleEndian()))
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

		prog.Add(MakeOpcodeInstance(op, args, table, lang.ByteSize()/8, lang.IsBigEndian()))

		return nil
	}
}
