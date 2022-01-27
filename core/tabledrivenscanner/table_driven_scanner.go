package tabledrivenscanner

import (
	"fmt"

	"github.com/obonobo/compiler/core/scanner"
)

type lexemeSpec struct {
	s         scanner.Lexeme
	line, col int
}

type TableDrivenScanner struct {
	chars scanner.CharSource // A source of characters
	table Table              // A table for performing transitions

	// The lexeme that is being built
	lexeme lexemeSpec
}

func NewTableDrivenScanner(chars scanner.CharSource, table Table) *TableDrivenScanner {
	return &TableDrivenScanner{
		chars: chars,
		table: table,
		lexeme: lexemeSpec{
			line: chars.Line(),
			col:  chars.Column(),
		},
	}
}

// Scans for the next token present in the character source
func (t *TableDrivenScanner) NextToken() (scanner.Token, error) {
	var token *scanner.Token
	state := t.table.Initial()
	for {
		lookup, err := t.nextChar()
		if err != nil {
			return scanner.Token{}, fmt.Errorf("TableDrivenScanner: %w", err)
		}

		state = t.table.Next(state, lookup)
		if state == 0 {
			return scanner.Token{},
				fmt.Errorf("TableDrivenScanner: no possible transition (state = 0)")
		}

		if t.table.IsFinal(state) {
			backtrack := t.table.NeedsBackup(state)
			if !backtrack {
				t.pushLexeme(lookup)
			}

			token, err = t.createToken(state)
			if err != nil {
				return scanner.Token{}, fmt.Errorf("TableDrivenScanner: %w", err)
			}

			if backtrack {
				t.backup()
			}
		}

		if token != nil {
			break
		}
		t.pushLexeme(lookup)
	}
	return *token, nil
}

func (t *TableDrivenScanner) pushLexeme(char rune) {
	t.lexeme.s += scanner.Lexeme(char)
}

func (t *TableDrivenScanner) backup() error {
	_, err := t.chars.BackupChar()
	t.lexeme.col = t.chars.Column()
	t.lexeme.line = t.chars.Line()
	return err
}

func (t *TableDrivenScanner) nextChar() (rune, error) {
	r, err := t.chars.NextChar()
	return r, err
}

func (t *TableDrivenScanner) createToken(
	state State,
) (*scanner.Token, error) {
	tt, err := t.table.CreateToken(state, t.lexeme.s, t.lexeme.line, t.lexeme.col)
	t.lexeme.s = ""
	t.lexeme.col = t.chars.Column()
	t.lexeme.line = t.chars.Line()
	return &tt, err
}
