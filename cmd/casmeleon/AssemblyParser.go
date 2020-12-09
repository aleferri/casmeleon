package main

import (
	"errors"
	"fmt"

	"github.com/aleferri/casmeleon/internal/casm"
	"github.com/aleferri/casmeleon/pkg/parser"
	"github.com/aleferri/casmeleon/pkg/text"
)

func IsDirective(s string) bool {
	return s == ".advance" || s == ".org" || s == ".alias" || s == ".db"
}

func ParseSourceLine(lang casm.Language, stream parser.Stream) error {
	parser.Consume(stream, text.EOL)
	name, err := parser.Require(stream, text.Identifier)
	if err != nil {
		return err
	}

	if IsDirective(name.Value()) {
		lastToken := name
		switch name.Value() {
		case ".advance":
			{
				parser.Require(stream, text.Number)
				parser.Consume(stream, text.WHITESPACE)
				lastToken = stream.Next()
			}
		case ".org":
			{
				parser.Require(stream, text.Number)
				parser.Consume(stream, text.WHITESPACE)
				lastToken = stream.Next()
			}
		case ".alias":
			{
				parser.RequireSequence(stream, text.Identifier, text.Number)
				parser.Consume(stream, text.WHITESPACE)
				lastToken = stream.Next()
			}
		case ".db":
			{
				lastToken = name.WithID(text.Comma)
				for lastToken.ID() == text.Comma {
					parser.RequireAny(stream, text.Identifier, text.Number, text.QuotedString)
					lastToken = stream.Next()
				}
			}
		}
		if lastToken.ID() != text.EOL {
			return fmt.Errorf("Expected End Of Line after the directive '%s', found instead '%s'", name.Value(), lastToken.Value())
		}
		fmt.Println("Successfully parsed the directive ", name.Value())
		return nil
	} else {
		if stream.Peek().ID() == text.Colon {
			stream.Next()
			fmt.Println("Found label: ", name.Value())
			return ParseSourceLine(lang, stream)
		} else {
			lastToken := stream.Next()

			args := []text.Symbol{}
			format := []uint32{}

			for lastToken.ID() != text.EOL {
				format = append(format, lastToken.ID())
				args = append(args, lastToken)
				lastToken = stream.Next()
			}

			win := lang.FilterOpcodesByName(name.Value())
			fmt.Println("Candidates number for ", name.Value(), ":", len(win.Candidates()))
			list := win.Candidates()
			for _, opc := range list {
				fmt.Println("Available format: ", opc.StringifyFormat(&lang))
			}
			fmt.Println("Provided Format: ", casm.StringifyFormat(format))
			win = win.FilterByFormat(format)

			op, err := win.PickFirst()
			if err != nil {
				return errors.New("Invalid opcode " + name.Value())
			}

			fmt.Println("Successfully got opcode " + op.Name())

			return nil
		}
	}
}
