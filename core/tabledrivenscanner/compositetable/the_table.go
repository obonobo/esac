package compositetable

import (
	"github.com/obonobo/esac/core/scanner"
	t "github.com/obonobo/esac/core/tabledrivenscanner"
	"github.com/obonobo/esac/core/token"
)

// Creates a new table driven scanner using TABLE() as the table
func NewTableDrivenScanner(chars scanner.CharSource) *t.TableDrivenScanner {
	return t.NewScanner(chars, TABLE())
}

// A factory function for creating the default implementation table that is used
// by the ESAC compiler
func TABLE() *CompositeTable {
	return &CompositeTable{
		Start:               t.START,
		UnterminatedComment: t.UNTERMINATEDCOMMENT,

		Transitions: map[Key]t.State{
			// INVALID CHARS
			{1, t.ANY}: 68,

			// SPACES
			{1, ' '}:    1,
			{1, '\n'}:   1,
			{1, '\t'}:   1,
			{1, '\r'}:   1,
			{1, '\x00'}: 1,

			// COMMENTS
			{1, '/'}:    2,
			{2, '/'}:    5,
			{5, t.ANY}:  5,
			{5, '\n'}:   6,
			{2, '*'}:    4,
			{2, t.ANY}:  3,
			{63, t.ANY}: 19,

			// OPERATORS AND PUNCTUATION
			{1, '='}:   9,
			{9, t.ANY}: 10,
			{9, '='}:   11,

			{1, '-'}:    13,
			{13, t.ANY}: 14,
			{13, '>'}:   15,

			{1, '<'}:    20,
			{20, t.ANY}: 21,
			{20, '>'}:   22,
			{20, '='}:   23,

			{1, '>'}:    24,
			{24, t.ANY}: 25,
			{24, '='}:   26,

			{1, ':'}:    36,
			{36, t.ANY}: 37,
			{36, ':'}:   38,

			{1, '+'}: 12,
			{1, '|'}: 16,
			{1, '&'}: 17,
			{1, '!'}: 18,
			{1, '*'}: 63,
			{1, '('}: 27,
			{1, ')'}: 28,
			{1, '{'}: 29,
			{1, '}'}: 30,
			{1, '['}: 31,
			{1, ']'}: 32,
			{1, ';'}: 33,
			{1, ','}: 34,
			{1, '.'}: 35,

			// ID AND RESERVED WORDS
			{1, t.LETTER}:  39,
			{39, t.LETTER}: 39,
			{39, '0'}:      39,
			{39, '1'}:      39,
			{39, '2'}:      39,
			{39, '3'}:      39,
			{39, '4'}:      39,
			{39, '5'}:      39,
			{39, '6'}:      39,
			{39, '7'}:      39,
			{39, '8'}:      39,
			{39, '9'}:      39,
			{39, '_'}:      39,
			{39, t.ANY}:    40,

			// INTS AND FLOATS
			{1, '1'}: 41,
			{1, '2'}: 41,
			{1, '3'}: 41,
			{1, '4'}: 41,
			{1, '5'}: 41,
			{1, '6'}: 41,
			{1, '7'}: 41,
			{1, '8'}: 41,
			{1, '9'}: 41,

			{41, '0'}: 41,
			{41, '1'}: 41,
			{41, '2'}: 41,
			{41, '3'}: 41,
			{41, '4'}: 41,
			{41, '5'}: 41,
			{41, '6'}: 41,
			{41, '7'}: 41,
			{41, '8'}: 41,
			{41, '9'}: 41,

			{1, '0'}:    42,
			{41, t.ANY}: 43,
			{42, t.ANY}: 43,
			{41, '.'}:   44,
			{42, '.'}:   44,

			{41, t.LETTER}: 58,
			{41, '_'}:      58,
			{42, t.LETTER}: 58,
			{42, '_'}:      58,

			{1, '_'}:       58,
			{58, '_'}:      58,
			{58, t.LETTER}: 58,
			{58, '0'}:      58,
			{58, '1'}:      58,
			{58, '2'}:      58,
			{58, '3'}:      58,
			{58, '4'}:      58,
			{58, '5'}:      58,
			{58, '6'}:      58,
			{58, '7'}:      58,
			{58, '8'}:      58,
			{58, '9'}:      58,
			{58, t.ANY}:    59,

			{42, '0'}: 56,
			{42, '1'}: 56,
			{42, '2'}: 56,
			{42, '3'}: 56,
			{42, '4'}: 56,
			{42, '5'}: 56,
			{42, '6'}: 56,
			{42, '7'}: 56,
			{42, '8'}: 56,
			{42, '9'}: 56,

			{56, t.ANY}: 57,
			{56, '0'}:   56,
			{56, '1'}:   56,
			{56, '2'}:   56,
			{56, '3'}:   56,
			{56, '4'}:   56,
			{56, '5'}:   56,
			{56, '6'}:   56,
			{56, '7'}:   56,
			{56, '8'}:   56,
			{56, '9'}:   56,
			{56, '.'}:   56,
			{56, 'e'}:   56,

			{44, '0'}:   45,
			{44, '1'}:   45,
			{44, '2'}:   45,
			{44, '3'}:   45,
			{44, '4'}:   45,
			{44, '5'}:   45,
			{44, '6'}:   45,
			{44, '7'}:   45,
			{44, '8'}:   45,
			{44, '9'}:   45,
			{44, t.ANY}: 60,

			{45, '1'}: 45,
			{45, '2'}: 45,
			{45, '3'}: 45,
			{45, '4'}: 45,
			{45, '5'}: 45,
			{45, '6'}: 45,
			{45, '7'}: 45,
			{45, '8'}: 45,
			{45, '9'}: 45,

			{45, t.ANY}: 46,
			{45, '0'}:   47,

			{47, '1'}: 45,
			{47, '2'}: 45,
			{47, '3'}: 45,
			{47, '4'}: 45,
			{47, '5'}: 45,
			{47, '6'}: 45,
			{47, '7'}: 45,
			{47, '8'}: 45,
			{47, '9'}: 45,

			{47, t.ANY}: 48,

			{45, 'e'}: 49,
			{49, '0'}: 51,
			{49, '-'}: 50,
			{49, '+'}: 50,

			{49, '1'}: 54,
			{49, '2'}: 54,
			{49, '3'}: 54,
			{49, '4'}: 54,
			{49, '5'}: 54,
			{49, '6'}: 54,
			{49, '7'}: 54,
			{49, '8'}: 54,
			{49, '9'}: 54,

			{49, t.ANY}: 62, // DOUBLE BACKTRACK

			{50, '0'}: 51,

			{50, '1'}:   54,
			{50, '2'}:   54,
			{50, '3'}:   54,
			{50, '4'}:   54,
			{50, '5'}:   54,
			{50, '6'}:   54,
			{50, '7'}:   54,
			{50, '8'}:   54,
			{50, '9'}:   54,
			{50, t.ANY}: 61,

			{51, '0'}: 52,
			{51, '1'}: 52,
			{51, '2'}: 52,
			{51, '3'}: 52,
			{51, '4'}: 52,
			{51, '5'}: 52,
			{51, '6'}: 52,
			{51, '7'}: 52,
			{51, '8'}: 52,
			{51, '9'}: 52,

			{51, t.ANY}: 53,

			{52, '0'}: 52,
			{52, '1'}: 52,
			{52, '2'}: 52,
			{52, '3'}: 52,
			{52, '4'}: 52,
			{52, '5'}: 52,
			{52, '6'}: 52,
			{52, '7'}: 52,
			{52, '8'}: 52,
			{52, '9'}: 52,

			{52, t.ANY}: 55,

			{54, '0'}: 54,
			{54, '1'}: 54,
			{54, '2'}: 54,
			{54, '3'}: 54,
			{54, '4'}: 54,
			{54, '5'}: 54,
			{54, '6'}: 54,
			{54, '7'}: 54,
			{54, '8'}: 54,
			{54, '9'}: 54,

			{54, t.ANY}: 53,
		},

		// STATES THAT NEED BACKUP
		NeedBackup: map[t.State]struct{}{
			3:  {},
			10: {},
			14: {},
			19: {},
			21: {},
			25: {},
			37: {},
			40: {},
			43: {},
			46: {},
			48: {},
			53: {},
			55: {},
			57: {},
			59: {},
			60: {},
			61: {},
		},

		NeedDoubleBackup: map[t.State]struct{}{
			62: {},
		},

		PopStates: map[t.State]struct{}{
			63: {},

			// 64: {},

			65: {},
		},

		PushStates: map[t.State]struct{}{
			2: {},
			4: {},

			// 5: {},
		},

		StackTransitions: map[Key]t.State{
			// Exit comments
			// {1, '\n'}:   64,
			{1, '*'}:    63,
			{63, t.ANY}: 1,
			{63, '/'}:   65,

			// Enter new comments
			{1, '/'}:   2,
			{2, t.ANY}: 1,

			// {2, '/'}:   5,
			{2, '*'}: 4,
		},

		Comments: map[token.Kind]t.State{
			token.INLINECMT: 66,
			token.BLOCKCMT:  67,
		},

		// STATE TO TOKEN MAPPING
		Tokens: map[t.State]token.Kind{
			3: token.DIV,

			4: token.OPENBLOCK,

			// 5:  token.OPENINLINE,
			6: token.INLINECMT,

			19: token.MULT,
			// 64: token.CLOSEINLINE,
			65: token.CLOSEBLOCK,
			66: token.INLINECMT,
			67: token.BLOCKCMT,

			10: token.ASSIGN,
			11: token.EQ,
			12: token.PLUS,
			14: token.MINUS,
			15: token.ARROW,
			16: token.OR,
			17: token.AND,
			18: token.NOT,
			21: token.LT,
			22: token.NOTEQ,
			23: token.LEQ,
			25: token.GT,
			26: token.GEQ,
			27: token.OPENPAR,
			28: token.CLOSEPAR,
			29: token.OPENCUBR,
			30: token.CLOSECUBR,
			31: token.OPENSQBR,
			32: token.CLOSESQBR,
			33: token.SEMI,
			34: token.COMMA,
			35: token.DOT,
			37: token.COLON,
			38: token.COLONCOLON,
			40: token.ID,
			43: token.INTNUM,
			46: token.FLOATNUM,
			53: token.FLOATNUM,
			48: token.INVALIDNUM,
			55: token.INVALIDNUM,
			57: token.INVALIDNUM,
			59: token.INVALIDID,
			60: token.INVALIDNUM,
			61: token.INVALIDNUM,
			62: token.FLOATNUM,
			68: token.INVALIDCHAR,

			t.UNTERMINATEDCOMMENT: token.UNTERMINATEDCOMMENT,
		},

		// WHICH SYMBOLS COUNT AS LETTERS
		Letters: map[rune]struct{}{
			'a': {},
			'b': {},
			'c': {},
			'd': {},
			'e': {},
			'f': {},
			'g': {},
			'h': {},
			'i': {},
			'j': {},
			'k': {},
			'l': {},
			'm': {},
			'n': {},
			'o': {},
			'p': {},
			'q': {},
			'r': {},
			's': {},
			't': {},
			'u': {},
			'v': {},
			'w': {},
			'x': {},
			'y': {},
			'z': {},

			'A': {},
			'B': {},
			'C': {},
			'D': {},
			'E': {},
			'F': {},
			'G': {},
			'H': {},
			'I': {},
			'J': {},
			'K': {},
			'L': {},
			'M': {},
			'N': {},
			'O': {},
			'P': {},
			'Q': {},
			'R': {},
			'S': {},
			'T': {},
			'U': {},
			'V': {},
			'W': {},
			'X': {},
			'Y': {},
			'Z': {},
		},

		Whitespace: map[rune]struct{}{
			' ':    {},
			'\n':   {},
			'\t':   {},
			'\r':   {},
			'\x00': {},
		},
	}
}
