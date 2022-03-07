package reporting

import (
	"fmt"
	"log"
	"strings"

	"github.com/obonobo/esac/core/scanner"
	"github.com/obonobo/esac/core/tabledrivenparser"
	"github.com/obonobo/esac/core/token"
	"github.com/obonobo/esac/util"
)

// Spools errors reported by the parser, logs them on the provided logger
func ErrSpool(logger *log.Logger) chan<- tabledrivenparser.ParserError {
	errc := make(chan tabledrivenparser.ParserError, 1024)
	go func() {
		for err := range errc {
			if logger != nil {
				logger.Printf(
					"Syntax error on line %v, column %v: %v",
					err.Tok.Line, err.Tok.Column, err.Err)
			}
		}
	}()
	return errc
}

// Spools rules reported by the parser, logs them on the provided logger
func RuleSpool(logger *log.Logger) chan<- token.Rule {
	rulec := make(chan token.Rule, 1024)
	go func() {
		for err := range rulec {
			if logger != nil {
				logger.Println(err)
			}
		}
	}()
	return rulec
}

// TODO: finish this function when you finish the CLI part
// Consumes the scanner and outputs the derivation as well as the ast produced
// from parsing (syntactical analysis).
func Parse(
	scnr scanner.Scanner,
) (
	ast <-chan token.AST,
	derivations <-chan string,
	errors <-chan error,
) {
	var (
		bufsize     = 1 << 16
		astc        = make(chan token.AST, 1)
		derivationc = make(chan string, bufsize)
		errc        = make(chan error, bufsize)
	)

	return astc, derivationc, errc
}

func StreamLinesSplitErrors(
	scnr scanner.Scanner,
	bufSize int,
) (tokens, errors <-chan string) {
	return StreamLinesOptionallySplitErrors(scnr, bufSize, true)
}

// Groups tokens by line and prints them on the output chan
func StreamLines(scnr scanner.Scanner, bufSize int) (lines <-chan string) {
	lines, _ = StreamLinesOptionallySplitErrors(scnr, bufSize, false)
	return lines
}

// Groups tokens by line and prints them on the output chan, optionally prints
// errors to the error chan
func StreamLinesOptionallySplitErrors(
	scnr scanner.Scanner,
	bufsize int,
	splitErrors bool,
) (
	tokens <-chan string,
	errors <-chan string,
) {
	bufsize = intOr1024(bufsize)
	out := make(chan string, bufsize)

	var errs chan string
	if splitErrors {
		errs = make(chan string, bufsize)
	}

	go func() {
		line := struct {
			n     int
			print string
		}{1, ""}

		resetLine := func(next token.Token) {
			if len(line.print) > 0 {
				out <- strings.TrimLeft(line.print, " ")
			}
			line.n = next.Line
			line.print = next.Report()
		}

		for {
			t, err := scnr.NextToken()
			if err != nil {
				break
			}

			// Errors tokens may optionally be split into a separate stream
			if splitErrors && token.IsError(t.Id) {
				errs <- errorify(t)
			} else {
				if t.Line != line.n {
					resetLine(t)
				} else {
					line.print += " " + t.Report()
				}
			}
		}
		resetLine(token.Token{})
		close(out)
		if splitErrors {
			close(errs)
		}
	}()

	return out, errs
}

func errorify(tok token.Token) string {
	if !token.IsError(tok.Id) {
		return ""
	}

	errTypes := map[token.Kind]string{
		token.INVALIDID:           "Invalid identifier",
		token.INVALIDCHAR:         "Invalid character",
		token.INVALIDNUM:          "Invalid number",
		token.UNTERMINATEDCOMMENT: "Unterminated comment",
	}

	return fmt.Sprintf(""+
		"Lexical error: %v: \"%v\": line %v.",
		errTypes[tok.Id], util.SingleLinify(string(tok.Lexeme)), tok.Line)
}

func intOr1024(i int) int {
	if i < 0 {
		return 1024
	}
	return i
}
