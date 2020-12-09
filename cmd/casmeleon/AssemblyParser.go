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
			}
		case ".org":
			{
				parser.Require(stream, text.Number)
			}
		case ".alias":
			{
				parser.RequireSequence(stream, text.Identifier, text.Number)
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
			return errors.New("Expected End Of Line after a directive")
		}
		fmt.Println("Successfully parsed a directive")
		return nil
	} else {
		if stream.Peek().ID() == text.Colon {
			stream.Next()
			return ParseSourceLine(lang, stream)
		} else {
			lastToken := name

			args := []text.Symbol{}
			format := []uint32{}

			for lastToken.ID() != text.EOL {
				format = append(format, lastToken.ID())
				args = append(args, lastToken)
				lastToken = stream.Next()
			}

			win := lang.FilterOpcodesByName(name.Value())
			win = win.FilterByFormat(format)

			op, err := win.PickFirst()
			if err != nil {
				return errors.New("Invalid opcode")
			}

			fmt.Println("Successfully got opcode " + op.Name())

			return nil
		}
	}
}
