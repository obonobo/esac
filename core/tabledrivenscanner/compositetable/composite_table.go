package compositetable

import (
	tds "github.com/obonobo/esac/core/tabledrivenscanner"
	"github.com/obonobo/esac/core/token"
)

const INITIAL tds.State = 1

type Key struct {
	state tds.State // Current state that the scanner is on
	next  rune      // The symbol that is being processed
}

// State transition table. Once initialized, it's contents should never be
// changed. The table should never be written to, only read from. CompositeTable
// has composite Key and Values
type CompositeTable struct {
	Start               tds.State
	UnterminatedComment tds.State

	Transitions map[Key]tds.State
	Tokens      map[tds.State]token.Kind
	Comments    map[token.Kind]tds.State

	NeedBackup       map[tds.State]struct{}
	NeedDoubleBackup map[tds.State]struct{}
	PopStates        map[tds.State]struct{}
	PushStates       map[tds.State]struct{}

	// Transitions that are only chosen if the stack is non-empty
	StackTransitions map[Key]tds.State

	Letters    map[rune]struct{}
	Whitespace map[rune]struct{}

	commentStack []token.Kind
}

// Perform a transition
func (t *CompositeTable) Next(state tds.State, char rune) tds.State {
	// If we are in a comment, try a different set of transitions
	if !t.commentStackIsEmpty() {
		if s, ok := t.StackTransitions[Key{state, char}]; ok {
			// Handle pop commentStack
			if t.isPopState(s) {
				if s, ok := t.handlePopState(s); ok {
					return s
				}
			} else if t.isPushState(s) {
				return s
			}
		}

		// If the symbol is not comment related, then we ignore it (self-loop on
		// start state)
		return t.Initial()
	}

	// Check the symbol itself
	if s, ok := t.Transitions[Key{state, char}]; ok {
		return s
	}

	// Check tabledrivenscanner.LETTER
	if _, isLetter := t.Letters[char]; isLetter {
		if s, ok := t.Transitions[Key{state, tds.LETTER}]; ok {
			return s
		}
	}

	// Check tabledrivenscanner.ANY state
	if s, ok := t.Transitions[Key{state, tds.ANY}]; ok {
		return s
	}

	return tds.NOSTATE
}

// Check if a state requires the scanner to backup
func (t *CompositeTable) NeedsBackup(state tds.State) bool {
	_, ok := t.NeedBackup[state]
	return ok
}

// Check if a state requires the scanner to backup TWICE
func (t *CompositeTable) NeedsDoubleBackup(state tds.State) bool {
	_, ok := t.NeedDoubleBackup[state]
	return ok
}

// The initial state
func (t *CompositeTable) Initial() tds.State {
	return t.Start
}

// Check if a state is a final state
func (t *CompositeTable) IsFinal(state tds.State) bool {
	_, ok := t.Tokens[state]
	return ok
}

// Checks if a symbol is whitespace
func (t *CompositeTable) IsWhiteSpace(char rune) bool {
	_, ok := t.Whitespace[char]
	return ok
}

// Report whether the table is in comment mode
func (t *CompositeTable) InCommentMode() bool {
	return !t.commentStackIsEmpty()
}

// Returns the UNTERMINATEDCOMMENT state of the table
func (t *CompositeTable) UnterminatedCommentState() tds.State {
	return t.UnterminatedComment
}

// Generates a token given a State
func (t *CompositeTable) CreateToken(
	state tds.State,
	lexeme token.Lexeme,
	line, col int,
) (token.Token, error) {
	symbol, ok := t.Tokens[state]
	if !ok {
		return token.Token{}, tds.UnrecognizedStateError(state)
	}

	// Handle push states
	if t.isPushState(state) {
		t.commentStackPush(symbol)
		return token.Token{
			Id:     symbol,
			Lexeme: lexeme,
			Line:   line,
			Column: col,
		}, tds.PartialTokenError{}
	}

	// IDs could actually be RESERVED WORDS
	if symbol == token.ID {
		if res, ok := token.IsReservedWordString(string(lexeme)); ok {
			symbol = res
		}
	}

	return token.Token{
		Id:     symbol,
		Lexeme: lexeme,
		Line:   line,
		Column: col,
	}, nil
}

// Call this function to unset the comment stack of the table
func (t *CompositeTable) ResetComments() {
	t.commentStack = make([]token.Kind, 0, cap(t.commentStack))
}

func (t *CompositeTable) handlePopState(
	state tds.State,
) (tds.State, bool) {
	if t.commentStackIsEmpty() {
		return tds.NOSTATE, false
	}

	// Now we have a `\n` or `*/` token
	token, ok := t.Tokens[state]
	if !ok {
		// We have pop token state, but it is not final
		return state, true
	}

	top := t.commentStackTop()

	// Now we have a full comment token
	comment, ok := t.matchCommentTokens(top, token)
	if !ok {
		return tds.NOSTATE, false
	}

	// Now we have a match. Pop the stack and return a final state, if the stack
	// is emptied
	t.commentStackPop()
	if t.commentStackIsEmpty() {
		return t.Comments[comment], true
	}
	return tds.NOSTATE, false
}

func (t *CompositeTable) matchCommentTokens(open, close token.Kind) (token.Kind, bool) {
	if open == token.OPENINLINE {
		return token.INLINECMT, close == token.CLOSEINLINE
	} else if open == token.OPENBLOCK {
		return token.BLOCKCMT, close == token.CLOSEBLOCK
	}
	return "", false
}

func (t *CompositeTable) isPopState(state tds.State) bool {
	_, ok := t.PopStates[state]
	return ok
}

func (t *CompositeTable) isPushState(state tds.State) bool {
	_, ok := t.PushStates[state]
	return ok
}

func (t *CompositeTable) commentStackIsEmpty() bool {
	return len(t.commentStack) == 0
}

func (t *CompositeTable) commentStackTop() token.Kind {
	if t.commentStackIsEmpty() {
		return ""
	}
	return t.commentStack[len(t.commentStack)-1]
}

func (t *CompositeTable) commentStackPush(symbol token.Kind) {
	t.commentStack = append(t.commentStack, symbol)
}

func (t *CompositeTable) commentStackPop() token.Kind {
	if t.commentStackIsEmpty() {
		return ""
	}
	s := t.commentStack[len(t.commentStack)-1]
	t.commentStack = t.commentStack[:len(t.commentStack)-1]
	return s
}
