package tabledrivenscanner

import (
	"errors"
	"fmt"

	"github.com/obonobo/esac/core/scanner"
	"github.com/obonobo/esac/core/token"
)

type lexemeSpec struct {
	s    token.Lexeme
	line int
	col  int
}

type TableDrivenScanner struct {
	chars  scanner.CharSource // A source of characters
	table  Table              // A table for performing transitions
	lexeme lexemeSpec         // The lexeme that is being built
	err    error              // The most recent error returned by the charsource
}

func NewScanner(chars scanner.CharSource, table Table) *TableDrivenScanner {
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
func (t *TableDrivenScanner) NextToken() (token.Token, error) {
	var tok *token.Token
	state := t.table.Initial()
	for {
		if t.err != nil {
			return token.Token{}, t.err
		}

		lookup, err := t.nextChar()

		if err != nil {
			// We are out of input
			t.err = fmt.Errorf("TableDrivenScanner.NextToken: %w", err)

			// If we are still in comment mode, then create and return an
			// UNTERMINATEDCOMMENT token
			if t.table.InCommentMode() {
				tt, _ := t.createToken(t.table.UnterminatedCommentState())
				return *tt, nil
			}

			// If there is an ANY transition available, then we can take it,
			// otherwise return the error
			if len(t.lexeme.s) > 0 {
				state = t.table.Next(state, ANY)
				if state == NOSTATE {
					return token.Token{}, t.err
				}
			} else {
				return token.Token{}, t.err
			}
		} else {
			state = t.table.Next(state, lookup)
		}

		// This branch is never supposed to be hit - it is here to reveal any
		// bugs in the transition table. If the transition table that is
		// provided to the TableDrivenScanner returns NOSTATE, then there is a
		// case unaccounted for.
		if state == NOSTATE {
			return token.Token{}, NoStateError{State: state, Lookup: lookup}
		}

		if t.table.IsFinal(state) {
			doubleBacktrack := t.table.NeedsDoubleBackup(state)
			backtrack := t.table.NeedsBackup(state)
			if !backtrack && !doubleBacktrack {
				t.pushLexeme(lookup)
			} else if doubleBacktrack {
				t.popLexeme()
			}

			tt, err := t.createToken(state)
			if e := new(PartialTokenError); errors.As(err, e) {
				state = t.table.Initial()
				continue
			}
			tok = tt

			if err != nil {
				return token.Token{}, fmt.Errorf("TableDrivenScanner: %w", err)
			}

			if backtrack {
				t.backup()
			} else if doubleBacktrack {
				t.backup()
				t.backup()
			}
		}

		if tok != nil {
			break
		}
		t.pushLexeme(lookup)
	}
	return *tok, nil
}

func (t *TableDrivenScanner) pushLexeme(char rune) {
	whitespace := t.table.IsWhiteSpace(char) && len(t.lexeme.s) == 0
	if whitespace {
		t.resetLexeme()
	} else {
		t.lexeme.s += token.Lexeme(char)
	}
}

func (t *TableDrivenScanner) popLexeme() {
	if len(t.lexeme.s) > 0 {
		t.lexeme.s = t.lexeme.s[:len(t.lexeme.s)-1]
	}
}

func (t *TableDrivenScanner) resetLexeme() {
	t.lexeme.s = ""
	t.lexeme.col = t.chars.Column()
	t.lexeme.line = t.chars.Line()
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
) (*token.Token, error) {
	tt, err := t.table.CreateToken(state, t.lexeme.s, t.lexeme.line, t.lexeme.col)
	if e := new(PartialTokenError); errors.As(err, e) {
		return &tt, fmt.Errorf("TableDrivenScanner: %w", err)
	}
	t.resetLexeme()
	return &tt, err
}
