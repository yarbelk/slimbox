package sh

import (
	"bufio"
	"io"
	"strings"
)

type frame struct {
	i            int
	s            string
	line, column int
}
type Lexer struct {
	// The lexer runs in its own goroutine, and communicates via channel 'ch'.
	ch      chan frame
	ch_stop chan bool
	// We record the level of nesting because the action could return, and a
	// subsequent call expects to pick up where it left off. In other words,
	// we're simulating a coroutine.
	// TODO: Support a channel-based variant that compatible with Go's yacc.
	stack []frame
	stale bool

	// The 'l' and 'c' fields were added for
	// https://github.com/wagerlabs/docker/blob/65694e801a7b80930961d70c69cba9f2465459be/buildfile.nex
	// Since then, I introduced the built-in Line() and Column() functions.
	l, c int

	parseResult interface{}

	// The following line makes it easy for scripts to insert fields in the
	// generated code.
	// [NEX_END_OF_LEXER_STRUCT]
}

// NewLexerWithInit creates a new Lexer object, runs the given callback on it,
// then returns it.
func NewLexerWithInit(in io.Reader, initFun func(*Lexer)) *Lexer {
	Shlex := new(Lexer)
	if initFun != nil {
		initFun(Shlex)
	}
	Shlex.ch = make(chan frame)
	Shlex.ch_stop = make(chan bool, 1)
	var scan func(in *bufio.Reader, ch chan frame, ch_stop chan bool, family []dfa, line, column int)
	scan = func(in *bufio.Reader, ch chan frame, ch_stop chan bool, family []dfa, line, column int) {
		// Index of DFA and length of highest-precedence match so far.
		matchi, matchn := 0, -1
		var buf []rune
		n := 0
		checkAccept := func(i int, st int) bool {
			// Higher precedence match? DFAs are run in parallel, so matchn is at most len(buf), hence we may omit the length equality check.
			if family[i].acc[st] && (matchn < n || matchi > i) {
				matchi, matchn = i, n
				return true
			}
			return false
		}
		var state [][2]int
		for i := 0; i < len(family); i++ {
			mark := make([]bool, len(family[i].startf))
			// Every DFA starts at state 0.
			st := 0
			for {
				state = append(state, [2]int{i, st})
				mark[st] = true
				// As we're at the start of input, follow all ^ transitions and append to our list of start states.
				st = family[i].startf[st]
				if -1 == st || mark[st] {
					break
				}
				// We only check for a match after at least one transition.
				checkAccept(i, st)
			}
		}
		atEOF := false
		stopped := false
		for {
			if n == len(buf) && !atEOF {
				r, _, err := in.ReadRune()
				switch err {
				case io.EOF:
					atEOF = true
				case nil:
					buf = append(buf, r)
				default:
					panic(err)
				}
			}
			if !atEOF {
				r := buf[n]
				n++
				var nextState [][2]int
				for _, x := range state {
					x[1] = family[x[0]].f[x[1]](r)
					if -1 == x[1] {
						continue
					}
					nextState = append(nextState, x)
					checkAccept(x[0], x[1])
				}
				state = nextState
			} else {
			dollar: // Handle $.
				for _, x := range state {
					mark := make([]bool, len(family[x[0]].endf))
					for {
						mark[x[1]] = true
						x[1] = family[x[0]].endf[x[1]]
						if -1 == x[1] || mark[x[1]] {
							break
						}
						if checkAccept(x[0], x[1]) {
							// Unlike before, we can break off the search. Now that we're at the end, there's no need to maintain the state of each DFA.
							break dollar
						}
					}
				}
				state = nil
			}

			if state == nil {
				lcUpdate := func(r rune) {
					if r == '\n' {
						line++
						column = 0
					} else {
						column++
					}
				}
				// All DFAs stuck. Return last match if it exists, otherwise advance by one rune and restart all DFAs.
				if matchn == -1 {
					if len(buf) == 0 { // This can only happen at the end of input.
						break
					}
					lcUpdate(buf[0])
					buf = buf[1:]
				} else {
					text := string(buf[:matchn])
					buf = buf[matchn:]
					matchn = -1
					select {
					case ch <- frame{matchi, text, line, column}:
						{
						}
					case stopped = <-ch_stop:
						{
						}
					}
					if stopped {
						break
					}
					if len(family[matchi].nest) > 0 {
						scan(bufio.NewReader(strings.NewReader(text)), ch, ch_stop, family[matchi].nest, line, column)
					}
					if atEOF {
						break
					}
					for _, r := range text {
						lcUpdate(r)
					}
				}
				n = 0
				for i := 0; i < len(family); i++ {
					state = append(state, [2]int{i, 0})
				}
			}
		}
		ch <- frame{-1, "", line, column}
	}
	go scan(bufio.NewReader(in), Shlex.ch, Shlex.ch_stop, dfas, 0, 0)
	return Shlex
}

type dfa struct {
	acc          []bool           // Accepting states.
	f            []func(rune) int // Transitions.
	startf, endf []int            // Transitions at start and end of input.
	nest         []dfa
}

var dfas = []dfa{
	//  \t
	{[]bool{false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 9:
				return -1
			case 32:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 9:
				return 2
			case 32:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 9:
				return -1
			case 32:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1}, nil},

	// (&&|\|\||;;|<<|>>|<&|>&|<>|<<-|>\|)
	{[]bool{false, false, false, false, false, false, true, true, true, true, true, true, true, true, true, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 38:
				return 1
			case 45:
				return -1
			case 59:
				return 2
			case 60:
				return 3
			case 62:
				return 4
			case 124:
				return 5
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 38:
				return 15
			case 45:
				return -1
			case 59:
				return -1
			case 60:
				return -1
			case 62:
				return -1
			case 124:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 38:
				return -1
			case 45:
				return -1
			case 59:
				return 14
			case 60:
				return -1
			case 62:
				return -1
			case 124:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 38:
				return 10
			case 45:
				return -1
			case 59:
				return -1
			case 60:
				return 11
			case 62:
				return 12
			case 124:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 38:
				return 7
			case 45:
				return -1
			case 59:
				return -1
			case 60:
				return -1
			case 62:
				return 8
			case 124:
				return 9
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 38:
				return -1
			case 45:
				return -1
			case 59:
				return -1
			case 60:
				return -1
			case 62:
				return -1
			case 124:
				return 6
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 38:
				return -1
			case 45:
				return -1
			case 59:
				return -1
			case 60:
				return -1
			case 62:
				return -1
			case 124:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 38:
				return -1
			case 45:
				return -1
			case 59:
				return -1
			case 60:
				return -1
			case 62:
				return -1
			case 124:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 38:
				return -1
			case 45:
				return -1
			case 59:
				return -1
			case 60:
				return -1
			case 62:
				return -1
			case 124:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 38:
				return -1
			case 45:
				return -1
			case 59:
				return -1
			case 60:
				return -1
			case 62:
				return -1
			case 124:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 38:
				return -1
			case 45:
				return -1
			case 59:
				return -1
			case 60:
				return -1
			case 62:
				return -1
			case 124:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 38:
				return -1
			case 45:
				return 13
			case 59:
				return -1
			case 60:
				return -1
			case 62:
				return -1
			case 124:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 38:
				return -1
			case 45:
				return -1
			case 59:
				return -1
			case 60:
				return -1
			case 62:
				return -1
			case 124:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 38:
				return -1
			case 45:
				return -1
			case 59:
				return -1
			case 60:
				return -1
			case 62:
				return -1
			case 124:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 38:
				return -1
			case 45:
				return -1
			case 59:
				return -1
			case 60:
				return -1
			case 62:
				return -1
			case 124:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 38:
				return -1
			case 45:
				return -1
			case 59:
				return -1
			case 60:
				return -1
			case 62:
				return -1
			case 124:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// (if|then|else|elif|fi|do|done|case|esac|while|until|for|{|}|!|in)
	{[]bool{false, true, false, false, false, false, false, false, false, false, true, true, false, false, false, true, false, false, false, true, false, false, true, true, true, true, false, true, false, false, false, true, false, false, true, true, true, false, true, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 33:
				return 1
			case 97:
				return -1
			case 99:
				return 2
			case 100:
				return 3
			case 101:
				return 4
			case 102:
				return 5
			case 104:
				return -1
			case 105:
				return 6
			case 108:
				return -1
			case 110:
				return -1
			case 111:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return 7
			case 117:
				return 8
			case 119:
				return 9
			case 123:
				return 10
			case 125:
				return 11
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return -1
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 102:
				return -1
			case 104:
				return -1
			case 105:
				return -1
			case 108:
				return -1
			case 110:
				return -1
			case 111:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return 39
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 102:
				return -1
			case 104:
				return -1
			case 105:
				return -1
			case 108:
				return -1
			case 110:
				return -1
			case 111:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return -1
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 102:
				return -1
			case 104:
				return -1
			case 105:
				return -1
			case 108:
				return -1
			case 110:
				return -1
			case 111:
				return 36
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return -1
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 102:
				return -1
			case 104:
				return -1
			case 105:
				return -1
			case 108:
				return 28
			case 110:
				return -1
			case 111:
				return -1
			case 114:
				return -1
			case 115:
				return 29
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return -1
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 102:
				return -1
			case 104:
				return -1
			case 105:
				return 25
			case 108:
				return -1
			case 110:
				return -1
			case 111:
				return 26
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return -1
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 102:
				return 23
			case 104:
				return -1
			case 105:
				return -1
			case 108:
				return -1
			case 110:
				return 24
			case 111:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return -1
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 102:
				return -1
			case 104:
				return 20
			case 105:
				return -1
			case 108:
				return -1
			case 110:
				return -1
			case 111:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return -1
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 102:
				return -1
			case 104:
				return -1
			case 105:
				return -1
			case 108:
				return -1
			case 110:
				return 16
			case 111:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return -1
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 102:
				return -1
			case 104:
				return 12
			case 105:
				return -1
			case 108:
				return -1
			case 110:
				return -1
			case 111:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return -1
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 102:
				return -1
			case 104:
				return -1
			case 105:
				return -1
			case 108:
				return -1
			case 110:
				return -1
			case 111:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return -1
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 102:
				return -1
			case 104:
				return -1
			case 105:
				return -1
			case 108:
				return -1
			case 110:
				return -1
			case 111:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return -1
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 102:
				return -1
			case 104:
				return -1
			case 105:
				return 13
			case 108:
				return -1
			case 110:
				return -1
			case 111:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return -1
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 102:
				return -1
			case 104:
				return -1
			case 105:
				return -1
			case 108:
				return 14
			case 110:
				return -1
			case 111:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return -1
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return 15
			case 102:
				return -1
			case 104:
				return -1
			case 105:
				return -1
			case 108:
				return -1
			case 110:
				return -1
			case 111:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return -1
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 102:
				return -1
			case 104:
				return -1
			case 105:
				return -1
			case 108:
				return -1
			case 110:
				return -1
			case 111:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return -1
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 102:
				return -1
			case 104:
				return -1
			case 105:
				return -1
			case 108:
				return -1
			case 110:
				return -1
			case 111:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return 17
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return -1
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 102:
				return -1
			case 104:
				return -1
			case 105:
				return 18
			case 108:
				return -1
			case 110:
				return -1
			case 111:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return -1
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 102:
				return -1
			case 104:
				return -1
			case 105:
				return -1
			case 108:
				return 19
			case 110:
				return -1
			case 111:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return -1
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 102:
				return -1
			case 104:
				return -1
			case 105:
				return -1
			case 108:
				return -1
			case 110:
				return -1
			case 111:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return -1
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return 21
			case 102:
				return -1
			case 104:
				return -1
			case 105:
				return -1
			case 108:
				return -1
			case 110:
				return -1
			case 111:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return -1
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 102:
				return -1
			case 104:
				return -1
			case 105:
				return -1
			case 108:
				return -1
			case 110:
				return 22
			case 111:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return -1
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 102:
				return -1
			case 104:
				return -1
			case 105:
				return -1
			case 108:
				return -1
			case 110:
				return -1
			case 111:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return -1
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 102:
				return -1
			case 104:
				return -1
			case 105:
				return -1
			case 108:
				return -1
			case 110:
				return -1
			case 111:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return -1
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 102:
				return -1
			case 104:
				return -1
			case 105:
				return -1
			case 108:
				return -1
			case 110:
				return -1
			case 111:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return -1
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 102:
				return -1
			case 104:
				return -1
			case 105:
				return -1
			case 108:
				return -1
			case 110:
				return -1
			case 111:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return -1
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 102:
				return -1
			case 104:
				return -1
			case 105:
				return -1
			case 108:
				return -1
			case 110:
				return -1
			case 111:
				return -1
			case 114:
				return 27
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return -1
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 102:
				return -1
			case 104:
				return -1
			case 105:
				return -1
			case 108:
				return -1
			case 110:
				return -1
			case 111:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return -1
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 102:
				return -1
			case 104:
				return -1
			case 105:
				return 32
			case 108:
				return -1
			case 110:
				return -1
			case 111:
				return -1
			case 114:
				return -1
			case 115:
				return 33
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return 30
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 102:
				return -1
			case 104:
				return -1
			case 105:
				return -1
			case 108:
				return -1
			case 110:
				return -1
			case 111:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return -1
			case 99:
				return 31
			case 100:
				return -1
			case 101:
				return -1
			case 102:
				return -1
			case 104:
				return -1
			case 105:
				return -1
			case 108:
				return -1
			case 110:
				return -1
			case 111:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return -1
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 102:
				return -1
			case 104:
				return -1
			case 105:
				return -1
			case 108:
				return -1
			case 110:
				return -1
			case 111:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return -1
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 102:
				return 35
			case 104:
				return -1
			case 105:
				return -1
			case 108:
				return -1
			case 110:
				return -1
			case 111:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return -1
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return 34
			case 102:
				return -1
			case 104:
				return -1
			case 105:
				return -1
			case 108:
				return -1
			case 110:
				return -1
			case 111:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return -1
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 102:
				return -1
			case 104:
				return -1
			case 105:
				return -1
			case 108:
				return -1
			case 110:
				return -1
			case 111:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return -1
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 102:
				return -1
			case 104:
				return -1
			case 105:
				return -1
			case 108:
				return -1
			case 110:
				return -1
			case 111:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return -1
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 102:
				return -1
			case 104:
				return -1
			case 105:
				return -1
			case 108:
				return -1
			case 110:
				return 37
			case 111:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return -1
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return 38
			case 102:
				return -1
			case 104:
				return -1
			case 105:
				return -1
			case 108:
				return -1
			case 110:
				return -1
			case 111:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return -1
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 102:
				return -1
			case 104:
				return -1
			case 105:
				return -1
			case 108:
				return -1
			case 110:
				return -1
			case 111:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return -1
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 102:
				return -1
			case 104:
				return -1
			case 105:
				return -1
			case 108:
				return -1
			case 110:
				return -1
			case 111:
				return -1
			case 114:
				return -1
			case 115:
				return 40
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return -1
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return 41
			case 102:
				return -1
			case 104:
				return -1
			case 105:
				return -1
			case 108:
				return -1
			case 110:
				return -1
			case 111:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 33:
				return -1
			case 97:
				return -1
			case 99:
				return -1
			case 100:
				return -1
			case 101:
				return -1
			case 102:
				return -1
			case 104:
				return -1
			case 105:
				return -1
			case 108:
				return -1
			case 110:
				return -1
			case 111:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			case 119:
				return -1
			case 123:
				return -1
			case 125:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, nil},
}

func NewLexer(in io.Reader) *Lexer {
	return NewLexerWithInit(in, nil)
}

func (ShLex *Lexer) Stop() {
	ShLex.ch_stop <- true
}

// Text returns the matched text.
func (Shlex *Lexer) Text() string {
	return Shlex.stack[len(Shlex.stack)-1].s
}

// Line returns the current line number.
// The first line is 0.
func (Shlex *Lexer) Line() int {
	if len(Shlex.stack) == 0 {
		return 0
	}
	return Shlex.stack[len(Shlex.stack)-1].line
}

// Column returns the current column number.
// The first column is 0.
func (Shlex *Lexer) Column() int {
	if len(Shlex.stack) == 0 {
		return 0
	}
	return Shlex.stack[len(Shlex.stack)-1].column
}

func (Shlex *Lexer) next(lvl int) int {
	if lvl == len(Shlex.stack) {
		l, c := 0, 0
		if lvl > 0 {
			l, c = Shlex.stack[lvl-1].line, Shlex.stack[lvl-1].column
		}
		Shlex.stack = append(Shlex.stack, frame{0, "", l, c})
	}
	if lvl == len(Shlex.stack)-1 {
		p := &Shlex.stack[lvl]
		*p = <-Shlex.ch
		Shlex.stale = false
	} else {
		Shlex.stale = true
	}
	return Shlex.stack[lvl].i
}
func (Shlex *Lexer) pop() {
	Shlex.stack = Shlex.stack[:len(Shlex.stack)-1]
}
func (Shlex Lexer) Error(e string) {
	panic(e)
}

// Lex runs the lexer. Always returns 0.
// When the -s option is given, this function is not generated;
// instead, the NN_FUN macro runs the lexer.
func (Shlex *Lexer) Lex(lval *ShSymType) int {
OUTER0:
	for {
		switch Shlex.next(0) {
		case 0:
			{ /*skip*/
			}
		case 1:
			{
				lval.empty = struct{}{}
				oText := Shlex.Text()
				return getOperator(oText)
			}
		case 2:
			{
				lval.empty = struct{}{}
				oText := Shlex.Text()
				return getReserved(oText)
			}
		default:
			break OUTER0
		}
		continue
	}
	Shlex.pop()

	return 0
}
